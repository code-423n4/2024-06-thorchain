package ethereum

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	ecore "github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
	etypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/thornode/bifrost/pkg/chainclients/shared/evm"
	"gitlab.com/thorchain/thornode/bifrost/pkg/chainclients/shared/runners"
	"gitlab.com/thorchain/thornode/bifrost/pkg/chainclients/shared/signercache"
	"gitlab.com/thorchain/thornode/bifrost/pkg/chainclients/shared/utxo"

	tssp "gitlab.com/thorchain/tss/go-tss/tss"

	"gitlab.com/thorchain/thornode/bifrost/blockscanner"
	"gitlab.com/thorchain/thornode/bifrost/metrics"
	"gitlab.com/thorchain/thornode/bifrost/pubkeymanager"
	"gitlab.com/thorchain/thornode/bifrost/thorclient"
	stypes "gitlab.com/thorchain/thornode/bifrost/thorclient/types"
	"gitlab.com/thorchain/thornode/bifrost/tss"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/config"
	"gitlab.com/thorchain/thornode/constants"
	"gitlab.com/thorchain/thornode/x/thorchain/aggregators"
	mem "gitlab.com/thorchain/thornode/x/thorchain/memo"
)

const (
	ethBlockRewardAndFee = 3 * 1e18
)

// Client is a structure to sign and broadcast tx to Ethereum chain used by signer mostly
type Client struct {
	logger                  zerolog.Logger
	cfg                     config.BifrostChainConfiguration
	localPubKey             common.PubKey
	client                  *ethclient.Client
	chainID                 *big.Int
	kw                      *evm.KeySignWrapper
	ethScanner              *ETHScanner
	bridge                  thorclient.ThorchainBridge
	blockScanner            *blockscanner.BlockScanner
	vaultABI                *abi.ABI
	pubkeyMgr               pubkeymanager.PubKeyValidator
	poolMgr                 thorclient.PoolManager
	asgardAddresses         []common.Address
	lastAsgard              time.Time
	tssKeySigner            *tss.KeySign
	wg                      *sync.WaitGroup
	stopchan                chan struct{}
	globalSolvencyQueue     chan stypes.Solvency
	signerCacheManager      *signercache.CacheManager
	lastSolvencyCheckHeight int64
}

// NewClient create new instance of Ethereum client
func NewClient(thorKeys *thorclient.Keys,
	cfg config.BifrostChainConfiguration,
	server *tssp.TssServer,
	bridge thorclient.ThorchainBridge,
	m *metrics.Metrics,
	pubkeyMgr pubkeymanager.PubKeyValidator,
	poolMgr thorclient.PoolManager,
) (*Client, error) {
	if thorKeys == nil {
		return nil, fmt.Errorf("fail to create ETH client,thor keys is empty")
	}
	tssKm, err := tss.NewKeySign(server, bridge)
	if err != nil {
		return nil, fmt.Errorf("fail to create tss signer: %w", err)
	}

	priv, err := thorKeys.GetPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("fail to get private key: %w", err)
	}

	temp, err := codec.ToTmPubKeyInterface(priv.PubKey())
	if err != nil {
		return nil, fmt.Errorf("fail to get tm pub key: %w", err)
	}
	pk, err := common.NewPubKeyFromCrypto(temp)
	if err != nil {
		return nil, fmt.Errorf("fail to get pub key: %w", err)
	}

	if bridge == nil {
		return nil, errors.New("THORChain bridge is nil")
	}
	if pubkeyMgr == nil {
		return nil, errors.New("pubkey manager is nil")
	}
	if poolMgr == nil {
		return nil, errors.New("pool manager is nil")
	}
	ethPrivateKey, err := evm.GetPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	ethClient, err := ethclient.Dial(cfg.RPCHost)
	if err != nil {
		return nil, fmt.Errorf("fail to dial ETH rpc host(%s): %w", cfg.RPCHost, err)
	}
	chainID, err := getChainID(ethClient, cfg.BlockScanner.HTTPRequestTimeout)
	if err != nil {
		return nil, err
	}

	keysignWrapper, err := evm.NewKeySignWrapper(ethPrivateKey, pk, tssKm, chainID, "ETH")
	if err != nil {
		return nil, fmt.Errorf("fail to create ETH key sign wrapper: %w", err)
	}
	vaultABI, _, err := evm.GetContractABI(routerContractABI, erc20ContractABI)
	if err != nil {
		return nil, fmt.Errorf("fail to get contract abi: %w", err)
	}
	pubkeyMgr.GetPubKeys()
	c := &Client{
		logger:       log.With().Str("module", "ethereum").Logger(),
		cfg:          cfg,
		client:       ethClient,
		chainID:      chainID,
		localPubKey:  pk,
		kw:           keysignWrapper,
		bridge:       bridge,
		vaultABI:     vaultABI,
		pubkeyMgr:    pubkeyMgr,
		poolMgr:      poolMgr,
		tssKeySigner: tssKm,
		wg:           &sync.WaitGroup{},
		stopchan:     make(chan struct{}),
	}

	c.logger.Info().Msgf("current chain id: %d", chainID.Uint64())
	if chainID.Uint64() == 0 {
		return nil, fmt.Errorf("chain id is: %d , invalid", chainID.Uint64())
	}
	var path string // if not set later, will in memory storage
	if len(c.cfg.BlockScanner.DBPath) > 0 {
		path = fmt.Sprintf("%s/%s", c.cfg.BlockScanner.DBPath, c.cfg.BlockScanner.ChainID)
	}
	storage, err := blockscanner.NewBlockScannerStorage(path, c.cfg.ScannerLevelDB)
	if err != nil {
		return c, fmt.Errorf("fail to create blockscanner storage: %w", err)
	}
	signerCacheManager, err := signercache.NewSignerCacheManager(storage.GetInternalDb())
	if err != nil {
		return nil, fmt.Errorf("fail to create signer cache manager")
	}

	c.signerCacheManager = signerCacheManager
	c.ethScanner, err = NewETHScanner(c.cfg.BlockScanner, storage, chainID, c.client, c.bridge, m, pubkeyMgr, c.ReportSolvency, signerCacheManager)
	if err != nil {
		return c, fmt.Errorf("fail to create eth block scanner: %w", err)
	}

	c.blockScanner, err = blockscanner.NewBlockScanner(c.cfg.BlockScanner, storage, m, c.bridge, c.ethScanner)
	if err != nil {
		return c, fmt.Errorf("fail to create block scanner: %w", err)
	}
	localNodeETHAddress, err := c.localPubKey.GetAddress(common.ETHChain)
	if err != nil {
		c.logger.Err(err).Msg("fail to get local node's ETH address")
	}
	c.logger.Info().Msgf("local node ETH address %s", localNodeETHAddress)

	return c, nil
}

// IsETH return true if the token address equals to ethToken address
func IsETH(token string) bool {
	return strings.EqualFold(token, ethToken)
}

// Start to monitor Ethereum block chain
func (c *Client) Start(globalTxsQueue chan stypes.TxIn, globalErrataQueue chan stypes.ErrataBlock, globalSolvencyQueue chan stypes.Solvency) {
	c.ethScanner.globalErrataQueue = globalErrataQueue
	c.globalSolvencyQueue = globalSolvencyQueue
	c.tssKeySigner.Start()
	c.blockScanner.Start(globalTxsQueue)
	c.wg.Add(1)
	go c.unstuck()
	c.wg.Add(1)
	go runners.SolvencyCheckRunner(c.GetChain(), c, c.bridge, c.stopchan, c.wg, constants.ThorchainBlockTime)
}

// Stop ETH client
func (c *Client) Stop() {
	c.tssKeySigner.Stop()
	c.blockScanner.Stop()
	c.client.Close()
	close(c.stopchan)
	c.wg.Wait()
}

func (c *Client) IsBlockScannerHealthy() bool {
	return c.blockScanner.IsHealthy()
}

// GetConfig return the configurations used by ETH chain
func (c *Client) GetConfig() config.BifrostChainConfiguration {
	return c.cfg
}

func (c *Client) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), c.cfg.BlockScanner.HTTPRequestTimeout)
}

// getChainID retrieve the chain id from ETH node, and determinate whether we are running on test net by checking the status
// when it failed to get chain id , it will assume LocalNet
func getChainID(client *ethclient.Client, timeout time.Duration) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get chain id ,err: %w", err)
	}
	return chainID, err
}

// GetChain get chain
func (c *Client) GetChain() common.Chain {
	return common.ETHChain
}

// GetHeight gets height from eth scanner
func (c *Client) GetHeight() (int64, error) {
	return c.ethScanner.GetHeight()
}

// GetBlockScannerHeight returns blockscanner height
func (c *Client) GetBlockScannerHeight() (int64, error) {
	return c.blockScanner.PreviousHeight(), nil
}

func (c *Client) GetLatestTxForVault(vault string) (string, string, error) {
	lastObserved, err := c.signerCacheManager.GetLatestRecordedTx(stypes.InboundCacheKey(vault, c.GetChain().String()))
	if err != nil {
		return "", "", err
	}
	lastBroadCasted, err := c.signerCacheManager.GetLatestRecordedTx(stypes.BroadcastCacheKey(vault, c.GetChain().String()))
	return lastObserved, lastBroadCasted, err
}

// GetAddress return current signer address, it will be bech32 encoded address
func (c *Client) GetAddress(poolPubKey common.PubKey) string {
	addr, err := poolPubKey.GetAddress(common.ETHChain)
	if err != nil {
		c.logger.Error().Err(err).Str("pool_pub_key", poolPubKey.String()).Msg("fail to get pool address")
		return ""
	}
	return addr.String()
}

// GetGasFee gets gas fee
func (c *Client) GetGasFee(gas uint64) common.Gas {
	return common.GetEVMGasFee(common.ETHChain, c.GetGasPrice(), gas)
}

// GetGasPrice gets gas price from eth scanner
func (c *Client) GetGasPrice() *big.Int {
	gasPrice := c.ethScanner.GetGasPrice()
	return gasPrice
}

// estimateGas estimates gas for tx
func (c *Client) estimateGas(from string, tx *etypes.Transaction) (uint64, error) {
	ctx, cancel := c.getContext()
	defer cancel()
	return c.client.EstimateGas(ctx, ethereum.CallMsg{
		From:     ecommon.HexToAddress(from),
		To:       tx.To(),
		GasPrice: tx.GasPrice(),
		// Gas:      tx.Gas(),
		Value: tx.Value(),
		Data:  tx.Data(),
	})
}

// GetNonce returns the nonce (including pending) for the given address.
func (c *Client) GetNonce(addr string) (uint64, error) {
	ctx, cancel := c.getContext()
	defer cancel()
	nonce, err := c.client.PendingNonceAt(ctx, ecommon.HexToAddress(addr))
	if err != nil {
		return 0, fmt.Errorf("fail to get account nonce: %w", err)
	}
	return nonce, nil
}

// GetNonceFinalized returns the nonce for the given address.
func (c *Client) GetNonceFinalized(addr string) (uint64, error) {
	ctx, cancel := c.getContext()
	defer cancel()
	return c.client.NonceAt(ctx, ecommon.HexToAddress(addr), nil)
}

func getTokenAddressFromAsset(asset common.Asset) string {
	if asset.Equals(common.ETHAsset) {
		return ethToken
	}
	allParts := strings.Split(asset.Symbol.String(), "-")
	return allParts[len(allParts)-1]
}

func (c *Client) getSmartContractAddr(pubkey common.PubKey) common.Address {
	return c.pubkeyMgr.GetContract(common.ETHChain, pubkey)
}

func (c *Client) getSmartContractByAddress(addr common.Address) common.Address {
	for _, pk := range c.pubkeyMgr.GetPubKeys() {
		ethAddr, err := pk.GetAddress(common.ETHChain)
		if err != nil {
			return common.NoAddress
		}
		if ethAddr.Equals(addr) {
			return c.pubkeyMgr.GetContract(common.ETHChain, pk)
		}
	}
	return common.NoAddress
}

func (c *Client) convertSigningAmount(amt *big.Int, token string) *big.Int {
	// convert 1e8 to 1e18
	amt = c.convertThorchainAmountToWei(amt)
	if IsETH(token) {
		return amt
	}
	tm, err := c.ethScanner.getTokenMeta(token)
	if err != nil {
		c.logger.Err(err).Msgf("fail to get token meta for token: %s", token)
		return amt
	}

	if tm.Decimal == defaultDecimals {
		// when the smart contract is using 1e18 as decimals , that means is based on WEI
		// thus the input amt is correct amount to send out
		return amt
	}
	var value big.Int
	amt = amt.Mul(amt, value.Exp(big.NewInt(10), big.NewInt(int64(tm.Decimal)), nil))
	amt = amt.Div(amt, value.Exp(big.NewInt(10), big.NewInt(defaultDecimals), nil))
	return amt
}

func (c *Client) convertThorchainAmountToWei(amt *big.Int) *big.Int {
	return big.NewInt(0).Mul(amt, big.NewInt(common.One*100))
}

// SignTx sign the the given TxArrayItem
func (c *Client) SignTx(tx stypes.TxOutItem, height int64) ([]byte, []byte, *stypes.TxInItem, error) {
	if !tx.Chain.Equals(common.ETHChain) {
		return nil, nil, nil, fmt.Errorf("chain %s is not support by ETH chain client", tx.Chain)
	}

	if c.signerCacheManager.HasSigned(tx.CacheHash()) {
		c.logger.Info().Msgf("transaction(%+v), signed before , ignore", tx)
		return nil, nil, nil, nil
	}

	if tx.ToAddress.IsEmpty() {
		return nil, nil, nil, fmt.Errorf("to address is empty")
	}
	if tx.VaultPubKey.IsEmpty() {
		return nil, nil, nil, fmt.Errorf("vault public key is empty")
	}

	if len(tx.Memo) == 0 {
		return nil, nil, nil, fmt.Errorf("can't sign tx when it doesn't have memo")
	}

	memo, err := mem.ParseMemo(common.LatestVersion, tx.Memo)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("fail to parse memo(%s):%w", tx.Memo, err)
	}

	if memo.IsInbound() {
		return nil, nil, nil, fmt.Errorf("inbound memo should not be used for outbound tx")
	}

	contractAddr := c.getSmartContractAddr(tx.VaultPubKey)
	if contractAddr.IsEmpty() {
		return nil, nil, nil, fmt.Errorf("can't sign tx , fail to get smart contract address")
	}

	value := big.NewInt(0)
	ethValue := big.NewInt(0)
	var tokenAddr string
	if len(tx.Coins) == 1 {
		coin := tx.Coins[0]
		tokenAddr = getTokenAddressFromAsset(coin.Asset)
		value = value.Add(value, coin.Amount.BigInt())
		value = c.convertSigningAmount(value, tokenAddr)
		if IsETH(tokenAddr) {
			ethValue = value
		}
	}

	fromAddr, err := tx.VaultPubKey.GetAddress(common.ETHChain)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("fail to get ETH address for pub key(%s): %w", tx.VaultPubKey, err)
	}

	dest := ecommon.HexToAddress(tx.ToAddress.String())
	var data []byte

	switch memo.GetType() {
	case mem.TxOutbound, mem.TxRefund, mem.TxRagnarok:
		if tx.Aggregator == "" {
			data, err = c.vaultABI.Pack("transferOut", dest, ecommon.HexToAddress(tokenAddr), value, tx.Memo)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("fail to create data to call smart contract(transferOut): %w", err)
			}
		} else {
			memoType := memo.GetType()
			if memoType == mem.TxRefund || memoType == mem.TxRagnarok {
				return nil, nil, nil, fmt.Errorf("%s can't use transferOutAndCall", memoType)
			}
			c.logger.Info().Msgf("aggregator target address: %s", tx.AggregatorTargetAsset)
			if ethValue.Uint64() == 0 {
				return nil, nil, nil, fmt.Errorf("transferOutAndCall can only be used when outbound asset is ETH")
			}
			targetLimit := tx.AggregatorTargetLimit
			if targetLimit == nil {
				zeroLimit := cosmos.ZeroUint()
				targetLimit = &zeroLimit
			}
			aggAddr := ecommon.HexToAddress(tx.Aggregator)
			targetAddr := ecommon.HexToAddress(tx.AggregatorTargetAsset)
			// when address can't be round trip , the tx out item will be dropped
			if !strings.EqualFold(aggAddr.String(), tx.Aggregator) {
				c.logger.Error().Msgf("aggregator address can't roundtrip , ignore tx (%s != %s)", tx.Aggregator, aggAddr.String())
				return nil, nil, nil, nil
			}
			if !strings.EqualFold(targetAddr.String(), tx.AggregatorTargetAsset) {
				c.logger.Error().Msgf("aggregator target asset address can't roundtrip , ignore tx (%s != %s)", tx.AggregatorTargetAsset, targetAddr.String())
				return nil, nil, nil, nil
			}
			data, err = c.vaultABI.Pack("transferOutAndCall", aggAddr, targetAddr, dest, targetLimit.BigInt(), tx.Memo)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("fail to create data to call smart contract(transferOutAndCall): %w", err)
			}
		}
	case mem.TxMigrate:
		if tx.Aggregator != "" || tx.AggregatorTargetAsset != "" {
			return nil, nil, nil, fmt.Errorf("migration can't use aggregator")
		}
		if IsETH(tokenAddr) {
			data, err = c.vaultABI.Pack("transferOut", dest, ecommon.HexToAddress(tokenAddr), value, tx.Memo)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("fail to create data to call smart contract(transferOut): %w", err)
			}
		} else {
			newSmartContractAddr := c.getSmartContractByAddress(tx.ToAddress)
			if newSmartContractAddr.IsEmpty() {
				return nil, nil, nil, fmt.Errorf("fail to get new smart contract address")
			}
			data, err = c.vaultABI.Pack("transferAllowance", ecommon.HexToAddress(newSmartContractAddr.String()), dest, ecommon.HexToAddress(tokenAddr), value, tx.Memo)
			if err != nil {
				return nil, nil, nil, fmt.Errorf("fail to create data to call smart contract(transferAllowance): %w", err)
			}
		}
	}

	// the nonce is stored as the transaction checkpoint, if it is set deserialize it
	// so we only retry with the same nonce to avoid double spend
	var nonce uint64
	if tx.Checkpoint != nil {
		if err = json.Unmarshal(tx.Checkpoint, &nonce); err != nil {
			return nil, nil, nil, fmt.Errorf("fail to deserialize checkpoint: %w", err)
		}
	} else {
		nonce, err = c.GetNonce(fromAddr.String())
		if err != nil {
			return nil, nil, nil, fmt.Errorf("fail to fetch account(%s) nonce : %w", fromAddr, err)
		}

		// abort signing if the pending nonce is too far in the future
		var finalizedNonce uint64
		finalizedNonce, err = c.GetNonceFinalized(fromAddr.String())
		if err != nil {
			return nil, nil, nil, fmt.Errorf("fail to fetch account(%s) finalized nonce: %w", fromAddr, err)
		}
		if (nonce - finalizedNonce) > c.cfg.MaxPendingNonces {
			c.logger.Warn().
				Uint64("nonce", nonce).
				Uint64("finalizedNonce", finalizedNonce).
				Msg("pending nonce too far in future")
			return nil, nil, nil, fmt.Errorf("pending nonce too far in future")
		}
	}
	c.logger.Info().Uint64("nonce", nonce).Msg("account info")

	// serialize nonce for later
	nonceBytes, err := json.Marshal(nonce)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("fail to marshal nonce: %w", err)
	}

	gasRate := c.GetGasPrice()
	if c.cfg.BlockScanner.FixedGasRate > 0 || gasRate.Cmp(big.NewInt(0)) == 0 {
		// if chain gas is zero we are still filling our gas price buffer, use outbound rate
		gasRate = c.convertThorchainAmountToWei(big.NewInt(tx.GasRate))
	} else {
		// Thornode uses a gas rate 1.5x the reported network fee for the rate and computed
		// max gas to ensure the rate is sufficient when it is signed later. Since we now know
		// the more recent rate, we will use our current rate with a lower bound on 2/3 the
		// outbound rate (the original rate we reported to Thornode in the network fee).
		lowerBound := c.convertThorchainAmountToWei(big.NewInt(tx.GasRate))
		lowerBound.Mul(lowerBound, big.NewInt(2))
		lowerBound.Div(lowerBound, big.NewInt(3))

		// round current rate to avoid consensus trouble, same rounding implied in outbound
		gasRate.Div(gasRate, big.NewInt(common.One*100))
		if gasRate.Cmp(big.NewInt(0)) == 0 { // floor at 1 like in network fee reporting
			gasRate = big.NewInt(1)
		}
		gasRate.Mul(gasRate, big.NewInt(common.One*100))

		// if the gas rate is less than the lower bound, use the lower bound
		if gasRate.Cmp(lowerBound) < 0 {
			gasRate = lowerBound
		}
	}

	// tip cap at configured percentage of max fee
	tipCap := new(big.Int).Mul(gasRate, big.NewInt(int64(c.cfg.MaxGasTipPercentage)))
	tipCap.Div(tipCap, big.NewInt(100))

	c.logger.Info().
		Stringer("inHash", tx.InHash).
		Str("outboundRate", c.convertThorchainAmountToWei(big.NewInt(tx.GasRate)).String()).
		Str("currentRate", c.GetGasPrice().String()).
		Str("effectiveRate", gasRate.String()).
		Msg("gas rate")

	// outbound tx always send to smart contract address
	estimatedETHValue := big.NewInt(0)
	if ethValue.Uint64() > 0 {
		// when the ETH value is none zero , here override it with a fix value for estimate gas purpose
		// when ETH value is none zero , if we send the real value for estimate gas , some times it will fail , for many reasons, a few I saw during test
		// 1. insufficient fund
		// 2. gas required exceeds allowance
		// as long as we pass in an ETH value , which we almost guarantee it will not exceed the ETH balance , so we can avoid the above two errors
		estimatedETHValue = estimatedETHValue.SetInt64(21000)
	}

	var createdTx *etypes.Transaction
	if c.cfg.BlockScanner.FixedGasRate == 0 {
		to := ecommon.HexToAddress(contractAddr.String())
		createdTx = etypes.NewTx(&etypes.DynamicFeeTx{
			ChainID:   c.chainID,
			Nonce:     nonce,
			To:        &to,
			Value:     estimatedETHValue,
			GasFeeCap: gasRate, // maxFeePerGas
			GasTipCap: tipCap,  // maxPriorityFeePerGas
			Data:      data,

			// gas is ignored in estimate gas call
			// Gas: c.cfg.BlockScanner.MaxGasLimit,
		})
	} else {
		createdTx = etypes.NewTransaction(nonce, ecommon.HexToAddress(contractAddr.String()), estimatedETHValue, c.cfg.BlockScanner.MaxGasLimit, gasRate, data)
	}

	estimatedGas, err := c.estimateGas(fromAddr.String(), createdTx)
	if err != nil {
		// in an edge case that vault doesn't have enough fund to fulfill an outbound transaction , it will fail to estimate gas
		// the returned error is `execution reverted`
		// when this fail , chain client should skip the outbound and move on to the next. The network will reschedule the outbound
		// after 300 blocks
		c.logger.Err(err).Msgf("fail to estimate gas")
		return nil, nil, nil, nil
	}
	c.logger.Info().Msgf("memo:%s estimated gas unit: %d", tx.Memo, estimatedGas)

	scheduledMaxFee := big.NewInt(0)
	for _, coin := range tx.MaxGas {
		scheduledMaxFee.Add(scheduledMaxFee, c.convertThorchainAmountToWei(coin.Amount.BigInt()))
	}

	if tx.Aggregator != "" {
		var gasLimitForAggregator uint64
		gasLimitForAggregator, err = aggregators.FetchDexAggregatorGasLimit(
			common.LatestVersion, c.cfg.ChainID, tx.Aggregator,
		)
		if err != nil {
			c.logger.Err(err).
				Str("aggregator", tx.Aggregator).
				Msg("fail to get aggregator gas limit, aborting to let thornode reschedule")
			return nil, nil, nil, nil
		}

		// if the estimate gas is over the max, abort and let thornode reschedule for now
		if estimatedGas > gasLimitForAggregator {
			c.logger.Warn().
				Stringer("in_hash", tx.InHash).
				Uint64("estimated_gas", estimatedGas).
				Uint64("aggregator_gas_limit", gasLimitForAggregator).
				Msg("aggregator gas limit exceeded, aborting to let thornode reschedule")
			return nil, nil, nil, nil
		}

		// set limit to aggregator gas limit
		estimatedGas = gasLimitForAggregator

		scheduledMaxFee = scheduledMaxFee.Mul(scheduledMaxFee, big.NewInt(c.cfg.AggregatorMaxGasMultiplier))
	} else if !tx.Coins[0].Asset.IsGasAsset() {
		scheduledMaxFee = scheduledMaxFee.Mul(scheduledMaxFee, big.NewInt(c.cfg.TokenMaxGasMultiplier))
	}

	if c.cfg.BlockScanner.FixedGasRate == 0 {
		// determine max gas units based on scheduled max gas (fee) and current rate
		maxGasUnits := new(big.Int).Div(scheduledMaxFee, gasRate).Uint64()

		// if estimated gas is more than the planned gas, abort and let thornode reschedule
		if estimatedGas > maxGasUnits {
			c.logger.Warn().
				Stringer("in_hash", tx.InHash).
				Stringer("rate", gasRate).
				Uint64("estimated_gas_units", estimatedGas).
				Uint64("max_gas_units", maxGasUnits).
				Str("scheduled_max_fee", scheduledMaxFee.String()).
				Msg("max gas exceeded, aborting to let thornode reschedule")
			return nil, nil, nil, nil
		}

		to := ecommon.HexToAddress(contractAddr.String())
		createdTx = etypes.NewTx(&etypes.DynamicFeeTx{
			ChainID:   c.chainID,
			Nonce:     nonce,
			To:        &to,
			Value:     ethValue,
			Gas:       maxGasUnits,
			GasFeeCap: gasRate,
			GasTipCap: tipCap,
			Data:      data,
		})
	} else {

		// if over max scheduled gas, abort and let thornode reschedule
		estimatedFee := big.NewInt(int64(estimatedGas) * gasRate.Int64())
		if scheduledMaxFee.Cmp(estimatedFee) < 0 {
			c.logger.Warn().
				Stringer("in_hash", tx.InHash).
				Stringer("rate", gasRate).
				Uint64("estimated_gas", estimatedGas).
				Str("estimated_fee", estimatedFee.String()).
				Str("scheduled_max_fee", scheduledMaxFee.String()).
				Msg("max gas exceeded, aborting to let thornode reschedule")
			return nil, nil, nil, nil
		}

		createdTx = etypes.NewTransaction(
			nonce, ecommon.HexToAddress(contractAddr.String()), ethValue, estimatedGas, gasRate, data,
		)
	}

	rawTx, err := c.sign(createdTx, tx.VaultPubKey, height, tx)
	if err != nil || len(rawTx) == 0 {
		return nil, nonceBytes, nil, fmt.Errorf("fail to sign message: %w", err)
	}

	// create the observation to be sent by the signer before broadcast
	chainHeight, err := c.GetHeight()
	if err != nil { // fall back to the scanner height, thornode voter does not use height
		chainHeight = c.ethScanner.currentBlockHeight
	}
	coin := tx.Coins[0]
	gas := common.MakeEVMGas(c.GetChain(), createdTx.GasPrice(), createdTx.Gas())
	// This is the maximum gas, using the gas limit for instant-observation
	// rather than the GasUsed which can only be gotten from the receipt when scanning.

	signedTx := &etypes.Transaction{}
	if err = signedTx.UnmarshalJSON(rawTx); err != nil {
		return nil, rawTx, nil, fmt.Errorf("fail to unmarshal signed tx: %w", err)
	}

	var txIn *stypes.TxInItem

	if err == nil {
		txIn = stypes.NewTxInItem(
			chainHeight+1,
			signedTx.Hash().Hex()[2:],
			tx.Memo,
			fromAddr.String(),
			tx.ToAddress.String(),
			common.NewCoins(
				coin,
			),
			gas,
			tx.VaultPubKey,
			"",
			"",
			nil,
		)
	}

	return rawTx, nil, txIn, nil
}

// sign is design to sign a given message with keysign party and keysign wrapper
func (c *Client) sign(tx *etypes.Transaction, poolPubKey common.PubKey, height int64, txOutItem stypes.TxOutItem) ([]byte, error) {
	rawBytes, err := c.kw.Sign(tx, poolPubKey)
	if err == nil && rawBytes != nil {
		return rawBytes, nil
	}
	var keysignError tss.KeysignError
	if errors.As(err, &keysignError) {
		if len(keysignError.Blame.BlameNodes) == 0 {
			// TSS doesn't know which node to blame
			return nil, fmt.Errorf("fail to sign tx: %w", err)
		}
		// key sign error forward the keysign blame to thorchain
		txID, errPostKeysignFail := c.bridge.PostKeysignFailure(keysignError.Blame, height, txOutItem.Memo, txOutItem.Coins, txOutItem.VaultPubKey)
		if errPostKeysignFail != nil {
			return nil, multierror.Append(err, errPostKeysignFail)
		}
		c.logger.Info().Str("tx_id", txID.String()).Msgf("post keysign failure to thorchain")
	}
	return nil, fmt.Errorf("fail to sign tx: %w", err)
}

// GetBalance call smart contract to find out the balance of the given address and token
func (c *Client) GetBalance(addr, token string, height *big.Int) (*big.Int, error) {
	ctx, cancel := c.getContext()
	defer cancel()
	if IsETH(token) {
		return c.client.BalanceAt(ctx, ecommon.HexToAddress(addr), height)
	}
	contractAddresses := c.pubkeyMgr.GetContracts(common.ETHChain)
	if len(contractAddresses) == 0 {
		return nil, fmt.Errorf("fail to get contract address")
	}
	input, err := c.vaultABI.Pack("vaultAllowance", ecommon.HexToAddress(addr), ecommon.HexToAddress(token))
	if err != nil {
		return nil, fmt.Errorf("fail to create vaultAllowance data to call smart contract")
	}
	c.logger.Debug().Msgf("query contract:%s for balance", contractAddresses[0].String())
	toAddr := ecommon.HexToAddress(contractAddresses[0].String())
	res, err := c.client.CallContract(ctx, ethereum.CallMsg{
		From: ecommon.HexToAddress(addr),
		To:   &toAddr,
		Data: input,
	}, height)
	if err != nil {
		return nil, err
	}
	output, err := c.vaultABI.Unpack("vaultAllowance", res)
	if err != nil {
		return nil, err
	}
	value, ok := abi.ConvertType(output[0], new(*big.Int)).(**big.Int)
	if !ok {
		return *value, fmt.Errorf("dev error: unable to get big.Int")
	}
	return *value, nil
}

// GetBalances gets all the balances of the given address
func (c *Client) GetBalances(addr string, height *big.Int) (common.Coins, error) {
	// for all the tokens , this chain client have deal with before
	tokens, err := c.ethScanner.GetTokens()
	if err != nil {
		return nil, fmt.Errorf("fail to get all the tokens: %w", err)
	}
	coins := common.Coins{}
	for _, token := range tokens {
		var balance *big.Int
		balance, err = c.GetBalance(addr, token.Address, height)
		if err != nil {
			c.logger.Err(err).Msgf("fail to get balance for token:%s", token.Address)
			continue
		}
		asset := common.ETHAsset
		if !IsETH(token.Address) {
			asset, err = common.NewAsset(fmt.Sprintf("ETH.%s-%s", token.Symbol, token.Address))
			if err != nil {
				return nil, err
			}
		}
		bal := c.ethScanner.convertAmount(token.Address, balance)
		coins = append(coins, common.NewCoin(asset, bal))
	}

	return coins.Distinct(), nil
}

// GetAccount gets account by address in eth client
func (c *Client) GetAccount(pk common.PubKey, height *big.Int) (common.Account, error) {
	addr := c.GetAddress(pk)
	nonce, err := c.GetNonce(addr)
	if err != nil {
		return common.Account{}, err
	}
	coins, err := c.GetBalances(addr, height)
	if err != nil {
		return common.Account{}, err
	}
	account := common.NewAccount(int64(nonce), 0, coins, false)
	return account, nil
}

// GetAccountByAddress return account information
func (c *Client) GetAccountByAddress(address string, height *big.Int) (common.Account, error) {
	nonce, err := c.GetNonce(address)
	if err != nil {
		return common.Account{}, err
	}
	coins, err := c.GetBalances(address, height)
	if err != nil {
		return common.Account{}, err
	}
	account := common.NewAccount(int64(nonce), 0, coins, false)
	return account, nil
}

// BroadcastTx decodes tx using rlp and broadcasts too Ethereum chain
func (c *Client) BroadcastTx(txOutItem stypes.TxOutItem, hexTx []byte) (string, error) {
	tx := &etypes.Transaction{}
	if err := tx.UnmarshalJSON(hexTx); err != nil {
		return "", err
	}
	ctx, cancel := c.getContext()
	defer cancel()
	if err := c.client.SendTransaction(ctx, tx); err != nil && err.Error() != txpool.ErrAlreadyKnown.Error() && err.Error() != ecore.ErrNonceTooLow.Error() {
		return "", err
	}
	txID := tx.Hash().String()
	c.logger.Info().Msgf("broadcast tx with memo: %s to ETH chain , hash: %s", txOutItem.Memo, txID)

	if err := c.signerCacheManager.SetSigned(txOutItem.CacheHash(), txOutItem.CacheVault(c.GetChain()), txID); err != nil {
		c.logger.Err(err).Msgf("fail to mark tx out item (%+v) as signed", txOutItem)
	}

	blockHeight, err := c.bridge.GetBlockHeight()
	if err != nil {
		c.logger.Err(err).Msgf("fail to get current THORChain block height")
		// at this point , the tx already broadcast successfully , don't return an error
		// otherwise will cause the same tx to retry
	} else if err = c.AddSignedTxItem(txID, blockHeight, txOutItem.VaultPubKey.String(), &txOutItem); err != nil {
		c.logger.Err(err).Msgf("fail to add signed tx item,hash:%s", txID)
	}

	return txID, nil
}

// ConfirmationCountReady check whether the given txIn is ready to be send to THORChain
func (c *Client) ConfirmationCountReady(txIn stypes.TxIn) bool {
	if len(txIn.TxArray) == 0 {
		return true
	}
	// MemPool items doesn't need confirmation
	if txIn.MemPool {
		return true
	}
	blockHeight := txIn.TxArray[0].BlockHeight
	confirm := txIn.ConfirmationRequired
	c.logger.Info().Msgf("confirmation required: %d", confirm)
	// every tx in txIn already have at least 1 confirmation
	return (c.ethScanner.currentBlockHeight - blockHeight) >= confirm
}

func (c *Client) getBlockReward(height int64) (*big.Int, error) {
	return big.NewInt(ethBlockRewardAndFee), nil
}

func (c *Client) getTotalTransactionValue(txIn stypes.TxIn, excludeFrom []common.Address) cosmos.Uint {
	total := cosmos.ZeroUint()
	if len(txIn.TxArray) == 0 {
		return total
	}
	for _, item := range txIn.TxArray {
		fromAsgard := false
		for _, fromAddress := range excludeFrom {
			if strings.EqualFold(fromAddress.String(), item.Sender) {
				fromAsgard = true
				break
			}
		}
		if fromAsgard {
			continue
		}
		for _, coin := range item.Coins {
			if coin.IsEmpty() {
				continue
			}
			amount := coin.Amount
			if !coin.Asset.Equals(common.ETHAsset) {
				var err error
				amount, err = c.poolMgr.GetValue(coin.Asset, common.ETHAsset, coin.Amount)
				if err != nil {
					c.logger.Err(err).Msgf("fail to get value for %s", coin.Asset)
					continue
				}

			}
			total = total.Add(amount)
		}
	}
	return total
}

// getBlockRequiredConfirmation find out how many confirmation the given txIn need to have before it can be send to THORChain
func (c *Client) getBlockRequiredConfirmation(txIn stypes.TxIn, height int64) (int64, error) {
	asgards, err := c.getAsgardAddress()
	if err != nil {
		c.logger.Err(err).Msg("fail to get asgard addresses")
		asgards = c.asgardAddresses
	}
	c.logger.Debug().Msgf("asgards: %+v", asgards)
	totalTxValue := c.getTotalTransactionValue(txIn, asgards)
	totalTxValueInWei := c.convertThorchainAmountToWei(totalTxValue.BigInt())
	confMul, err := utxo.GetConfMulBasisPoint(c.GetChain().String(), c.bridge)
	if err != nil {
		c.logger.Err(err).Msgf("failed to get conf multiplier mimir value for %s", c.GetChain().String())
	}
	totalFeeAndSubsidy, err := c.getBlockReward(height)
	confValue := common.GetUncappedShare(confMul, cosmos.NewUint(constants.MaxBasisPts), cosmos.NewUintFromBigInt(totalFeeAndSubsidy))
	if err != nil {
		return 0, fmt.Errorf("fail to get coinbase value: %w", err)
	}
	confirm := cosmos.NewUintFromBigInt(totalTxValueInWei).MulUint64(2).Quo(confValue).Uint64()
	confirm, err = utxo.MaxConfAdjustment(confirm, c.GetChain().String(), c.bridge)
	if err != nil {
		c.logger.Err(err).Msgf("fail to get max conf value adjustment for %s", c.GetChain().String())
	}
	c.logger.Info().Msgf("totalTxValue:%s,total fee and Subsidy:%d,confirmation:%d", totalTxValueInWei, totalFeeAndSubsidy, confirm)
	if confirm < 2 {
		// in ETH PoS (post merge) reorgs are harder to do but can occur. In
		// looking at 1k reorg blocks, 10 were reorg'ed at a height of 2, and
		// the rest were one (none were three or larger). While the odds of
		// getting reorg'ed are small (as it can only happen for very small
		// trades), the additional delay to swappers is also small (12 secs or
		// so). Thus, the determination by thorsec, 9R and devs were to set the
		// new min conf is 2.
		return 2, nil
	}
	return int64(confirm), nil
}

// GetConfirmationCount decide the given txIn how many confirmation it requires
func (c *Client) GetConfirmationCount(txIn stypes.TxIn) int64 {
	if len(txIn.TxArray) == 0 {
		return 0
	}
	// MemPool items doesn't need confirmation
	if txIn.MemPool {
		return 0
	}
	blockHeight := txIn.TxArray[0].BlockHeight
	confirm, err := c.getBlockRequiredConfirmation(txIn, blockHeight)
	c.logger.Debug().Msgf("confirmation required: %d", confirm)
	if err != nil {
		c.logger.Err(err).Msg("fail to get block confirmation ")
		return 0
	}
	return confirm
}

func (c *Client) getAsgardAddress() ([]common.Address, error) {
	if time.Since(c.lastAsgard) < constants.ThorchainBlockTime && c.asgardAddresses != nil {
		return c.asgardAddresses, nil
	}
	newAddresses, err := utxo.GetAsgardAddress(common.ETHChain, c.bridge)
	if err != nil {
		return nil, fmt.Errorf("fail to get asgards : %w", err)
	}
	if len(newAddresses) > 0 { // ensure we don't overwrite with empty list
		c.asgardAddresses = newAddresses
	}
	c.lastAsgard = time.Now()
	return c.asgardAddresses, nil
}

// OnObservedTxIn gets called from observer when we have a valid observation
func (c *Client) OnObservedTxIn(txIn stypes.TxInItem, blockHeight int64) {
	c.ethScanner.onObservedTxIn(txIn, blockHeight)
	m, err := mem.ParseMemo(common.LatestVersion, txIn.Memo)
	if err != nil {
		// Debug log only as ParseMemo error is expected for THORName inbounds.
		c.logger.Debug().Err(err).Msgf("fail to parse memo: %s", txIn.Memo)
		return
	}
	if !m.IsOutbound() {
		return
	}
	if m.GetTxID().IsEmpty() {
		return
	}
	if err = c.signerCacheManager.SetSigned(txIn.CacheHash(c.GetChain(), m.GetTxID().String()), txIn.CacheVault(c.GetChain()), txIn.Tx); err != nil {
		c.logger.Err(err).Msg("fail to update signer cache")
	}
}

func (c *Client) ReportSolvency(ethBlockHeight int64) error {
	if !c.ShouldReportSolvency(ethBlockHeight) {
		return nil
	}

	// when block scanner is not healthy, only report from auto-unhalt SolvencyCheckRunner
	// (FetchTxs passes currentBlockHeight, while SolvencyCheckRunner passes chainHeight)
	if !c.IsBlockScannerHealthy() && ethBlockHeight == c.ethScanner.currentBlockHeight {
		return nil
	}

	// fetch all asgard vaults
	asgardVaults, err := c.bridge.GetAsgards()
	if err != nil {
		return fmt.Errorf("fail to get asgards,err: %w", err)
	}

	currentGasFee := cosmos.NewUint(3 * c.cfg.BlockScanner.MaxGasLimit * c.ethScanner.lastReportedGasPrice)

	// report insolvent asgard vaults,
	// or else all if the chain is halted and all are solvent
	msgs := make([]stypes.Solvency, 0, len(asgardVaults))
	solventMsgs := make([]stypes.Solvency, 0, len(asgardVaults))
	for i := range asgardVaults {
		var acct common.Account
		acct, err = c.GetAccount(asgardVaults[i].PubKey, new(big.Int).SetInt64(ethBlockHeight))
		if err != nil {
			c.logger.Err(err).Msgf("fail to get account balance")
			continue
		}

		msg := stypes.Solvency{
			Height: ethBlockHeight,
			Chain:  common.ETHChain,
			PubKey: asgardVaults[i].PubKey,
			Coins:  acct.Coins,
		}

		if runners.IsVaultSolvent(acct, asgardVaults[i], currentGasFee) {
			solventMsgs = append(solventMsgs, msg) // Solvent-vault message
			continue
		}
		msgs = append(msgs, msg) // Insolvent-vault message
	}

	// Only if the block scanner is unhealthy (e.g. solvency-halted) and all vaults are solvent,
	// report that all the vaults are solvent.
	// If there are any insolvent vaults, report only them.
	// Not reporting both solvent and insolvent vaults is to avoid noise (spam):
	// Reporting both could halt-and-unhalt SolvencyHalt in the same THOR block
	// (resetting its height), plus making it harder to know at a glance from solvency reports which vaults were insolvent.
	solvent := false
	if !c.IsBlockScannerHealthy() && len(solventMsgs) == len(asgardVaults) {
		msgs = solventMsgs
		solvent = true
	}

	for i := range msgs {
		c.logger.Info().
			Stringer("asgard", msgs[i].PubKey).
			Interface("coins", msgs[i].Coins).
			Bool("solvent", solvent).
			Msg("reporting solvency")

		// send solvency to thorchain via global queue consumed by the observer
		select {
		case c.globalSolvencyQueue <- msgs[i]:
		case <-time.After(constants.ThorchainBlockTime):
			c.logger.Info().Msgf("fail to send solvency info to THORChain, timeout")
		}
	}
	c.lastSolvencyCheckHeight = ethBlockHeight
	return nil
}

// ShouldReportSolvency with given block height , should chain client report Solvency to THORNode?
func (c *Client) ShouldReportSolvency(height int64) bool {
	return height%20 == 0
}
