package utxo

import (
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/atomic"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"

	"gitlab.com/thorchain/thornode/bifrost/blockscanner"
	btypes "gitlab.com/thorchain/thornode/bifrost/blockscanner/types"
	"gitlab.com/thorchain/thornode/bifrost/metrics"
	"gitlab.com/thorchain/thornode/bifrost/pkg/chainclients/shared/runners"
	"gitlab.com/thorchain/thornode/bifrost/pkg/chainclients/shared/signercache"
	"gitlab.com/thorchain/thornode/bifrost/pkg/chainclients/shared/utxo"
	"gitlab.com/thorchain/thornode/bifrost/pkg/chainclients/utxo/rpc"
	"gitlab.com/thorchain/thornode/bifrost/thorclient"
	"gitlab.com/thorchain/thornode/bifrost/thorclient/types"
	"gitlab.com/thorchain/thornode/bifrost/tss"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/config"
	"gitlab.com/thorchain/thornode/constants"
	mem "gitlab.com/thorchain/thornode/x/thorchain/memo"
	gotss "gitlab.com/thorchain/tss/go-tss/tss"
)

////////////////////////////////////////////////////////////////////////////////////////
// Generate
////////////////////////////////////////////////////////////////////////////////////////

//go:generate go run generate.go

////////////////////////////////////////////////////////////////////////////////////////
// Client - Base
////////////////////////////////////////////////////////////////////////////////////////

// Client defines a generic UTXO client. Since there are differences in addresses, RPCs,
// and txscript between chains, chain additions should audit switches on chain type and
// extend where appropriate.
type Client struct {
	cfg config.BifrostChainConfiguration
	log zerolog.Logger
	m   *metrics.Metrics
	rpc *rpc.Client

	// ---------- signing ----------
	nodePubKey         common.PubKey
	nodePrivKey        *btcec.PrivateKey
	tssKeySigner       tss.ThorchainKeyManager
	signerCacheManager *signercache.CacheManager

	// ---------- sync ----------
	wg                    *sync.WaitGroup
	signerLock            *sync.Mutex
	vaultSignerLocks      map[string]*sync.Mutex
	consolidateInProgress *atomic.Bool

	// ---------- scanner ----------
	blockScanner    *blockscanner.BlockScanner
	temporalStorage *utxo.TemporalStorage

	// ---------- control ----------
	globalErrataQueue   chan<- types.ErrataBlock
	globalSolvencyQueue chan<- types.Solvency
	stopchan            chan struct{}
	currentBlockHeight  *atomic.Int64

	// ---------- thornode state ----------
	bridge          thorclient.ThorchainBridge
	asgardAddresses []common.Address
	lastAsgard      time.Time

	// TODO: The current block height uses an atomic int64, but the following values are
	// accessed without a lock. Update everything to use atomic values and switch to the
	// standard library sync/atomic package (recent versions include all necessary types).

	// ---------- fees / solvency ----------
	minRelayFeeSats         uint64
	lastFeeRate             uint64
	feeRateCache            []uint64
	lastSolvencyCheckHeight int64

	// ---------- testing ----------
	disableVinZeroBatch bool
}

// NewClient generates a new Client
func NewClient(
	thorKeys *thorclient.Keys,
	cfg config.BifrostChainConfiguration,
	server *gotss.TssServer,
	bridge thorclient.ThorchainBridge,
	m *metrics.Metrics,
) (*Client, error) {
	// verify the chain is supported
	supported := map[common.Chain]bool{
		common.DOGEChain: true,
		common.BCHChain:  true,
		common.LTCChain:  true,
		common.BTCChain:  true,
	}
	if !supported[cfg.ChainID] {
		return nil, fmt.Errorf("unsupported utxo chain: %s", cfg.ChainID)
	}

	logger := log.Logger.With().Stringer("chain", cfg.ChainID).Logger()

	// create rpc client
	rpcClient, err := rpc.NewClient(cfg.RPCHost, cfg.UserName, cfg.Password, cfg.MaxRPCRetries, logger)
	if err != nil {
		return nil, fmt.Errorf("fail to create rpc client: %w", err)
	}

	// node key setup
	tssKeysign, err := tss.NewKeySign(server, bridge)
	if err != nil {
		return nil, fmt.Errorf("fail to create tss signer: %w", err)
	}
	thorPrivateKey, err := thorKeys.GetPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("fail to get THORChain private key: %w", err)
	}
	nodePrivKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), thorPrivateKey.Bytes())
	nodePubKey, err := bech32AccountPubKey(nodePrivKey)
	if err != nil {
		return nil, fmt.Errorf("fail to get node account public key: %w", err)
	}

	// create base client
	c := &Client{
		cfg:                   cfg,
		log:                   logger,
		m:                     m,
		rpc:                   rpcClient,
		nodePubKey:            nodePubKey,
		nodePrivKey:           nodePrivKey,
		tssKeySigner:          tssKeysign,
		wg:                    &sync.WaitGroup{},
		signerLock:            &sync.Mutex{},
		vaultSignerLocks:      make(map[string]*sync.Mutex),
		consolidateInProgress: atomic.NewBool(false),
		stopchan:              make(chan struct{}),
		currentBlockHeight:    atomic.NewInt64(0),
		bridge:                bridge,
	}

	// import the node local address in the daemon wallet
	if err = c.RegisterPublicKey(c.nodePubKey); err != nil {
		return nil, fmt.Errorf("fail to register (%s): %w", c.nodePubKey, err)
	}

	var path string // fallback to in memory storage if unset
	if len(c.cfg.BlockScanner.DBPath) > 0 {
		path = fmt.Sprintf("%s/%s", c.cfg.BlockScanner.DBPath, c.cfg.BlockScanner.ChainID)
	}
	storage, err := blockscanner.NewBlockScannerStorage(path, c.cfg.ScannerLevelDB)
	if err != nil {
		return c, fmt.Errorf("fail to create blockscanner storage: %w", err)
	}

	c.blockScanner, err = blockscanner.NewBlockScanner(c.cfg.BlockScanner, storage, m, bridge, c)
	if err != nil {
		return c, fmt.Errorf("fail to create block scanner: %w", err)
	}

	c.temporalStorage, err = utxo.NewTemporalStorage(storage.GetInternalDb(), c.cfg.MempoolTxIDCacheSize)
	if err != nil {
		return c, fmt.Errorf("fail to create utxo storage: %w", err)
	}

	signerCacheManager, err := signercache.NewSignerCacheManager(storage.GetInternalDb())
	if err != nil {
		return nil, fmt.Errorf("fail to create signer cache manager, err: %w", err)
	}
	c.signerCacheManager = signerCacheManager
	c.updateNetworkInfo()

	return c, nil
}

// GetConfig returns the chain configuration.
func (c *Client) GetConfig() config.BifrostChainConfiguration {
	return c.cfg
}

// GetChain returns the chain ID.
func (c *Client) GetChain() common.Chain {
	return c.cfg.ChainID
}

// IsBlockScannerHealthy returns true if the block scanner is healthy.
func (c *Client) IsBlockScannerHealthy() bool {
	return c.blockScanner.IsHealthy()
}

// GetHeight returns current chain (not scanner) height.
func (c *Client) GetHeight() (int64, error) {
	return c.rpc.GetBlockCount()
}

// GetBlockScannerHeight returns blockscanner height
func (c *Client) GetBlockScannerHeight() (int64, error) {
	return c.blockScanner.PreviousHeight(), nil
}

func (c *Client) GetLatestTxForVault(vault string) (string, string, error) {
	lastObserved, err := c.signerCacheManager.GetLatestRecordedTx(types.InboundCacheKey(vault, c.GetChain().String()))
	if err != nil {
		return "", "", err
	}
	lastBroadCasted, err := c.signerCacheManager.GetLatestRecordedTx(types.BroadcastCacheKey(vault, c.GetChain().String()))
	return lastObserved, lastBroadCasted, err
}

// GetAddress returns chain address for the given public key.
func (c *Client) GetAddress(pubkey common.PubKey) string {
	addr, err := pubkey.GetAddress(c.cfg.ChainID)
	if err != nil {
		c.log.Error().Err(err).Str("pubkey", pubkey.String()).Msg("fail to get pool address")
		return ""
	}
	return addr.String()
}

////////////////////////////////////////////////////////////////////////////////////////
// Client - Control
////////////////////////////////////////////////////////////////////////////////////////

// Start starts the scanner, signer, and solvency check.
func (c *Client) Start(
	globalTxsQueue chan types.TxIn,
	globalErrataQueue chan types.ErrataBlock,
	globalSolvencyQueue chan types.Solvency,
) {
	c.globalErrataQueue = globalErrataQueue
	c.globalSolvencyQueue = globalSolvencyQueue
	c.tssKeySigner.Start()
	c.blockScanner.Start(globalTxsQueue)
	c.wg.Add(1)
	go runners.SolvencyCheckRunner(
		c.GetChain(), c, c.bridge, c.stopchan, c.wg, constants.ThorchainBlockTime,
	)
}

// Stop stops the scanner, signer, and solvency check.
func (c *Client) Stop() {
	c.blockScanner.Stop()
	c.tssKeySigner.Stop()
	close(c.stopchan)
	c.wg.Wait()
}

////////////////////////////////////////////////////////////////////////////////////////
// Client - Accounts
////////////////////////////////////////////////////////////////////////////////////////

// RegisterPublicKey imports the provided public key in the chain daemon.
func (c *Client) RegisterPublicKey(pubkey common.PubKey) error {
	addr, err := pubkey.GetAddress(c.cfg.ChainID)
	if err != nil {
		return fmt.Errorf("fail to get address from pubkey(%s): %w", pubkey, err)
	}

	// litecoin does not have a default wallet so we need to create one
	switch c.cfg.ChainID {
	case common.LTCChain, common.BTCChain:
		err = c.rpc.CreateWallet("")
		if err != nil {
			c.log.Info().Err(err).Msg("fail to create wallet")
			return err
		}
	}

	err = c.rpc.ImportAddress(addr.String())
	if err != nil {
		c.log.Error().Err(err).
			Str("pubkey", pubkey.String()).
			Str("addr", addr.String()).
			Msg("fail to import address")
	}
	return err
}

// GetAccount returns the account details for the given public key.
func (c *Client) GetAccount(pubkey common.PubKey, height *big.Int) (common.Account, error) {
	acct := common.Account{}
	if pubkey.IsEmpty() {
		return acct, errors.New("pubkey can't be empty")
	}

	// get all unspent utxos
	addr, err := pubkey.GetAddress(c.cfg.ChainID)
	if err != nil {
		return acct, fmt.Errorf("fail to get address from pubkey(%s): %w", pubkey, err)
	}
	utxos, err := c.rpc.ListUnspent(addr.String())
	if err != nil {
		return acct, fmt.Errorf("fail to get UTXOs: %w", err)
	}

	total := 0.0
	for _, item := range utxos {
		if !c.isValidUTXO(item.ScriptPubKey) {
			continue
		}
		if item.Confirmations == 0 {
			// pending tx in mempool, only count sends from asgard
			if !c.isSelfTransaction(item.TxID) && !c.isFromAsgard(item.TxID) {
				continue
			}
		}
		total += item.Amount
	}
	totalAmt, err := btcutil.NewAmount(total)
	if err != nil {
		return acct, fmt.Errorf("fail to convert total amount: %w", err)
	}
	return common.NewAccount(0, 0,
		common.Coins{
			common.NewCoin(c.cfg.ChainID.GetGasAsset(), cosmos.NewUint(uint64(totalAmt))),
		}, false), nil
}

// GetAccountByAddress is unimplemented for UTXO chains.
func (c *Client) GetAccountByAddress(address string, height *big.Int) (common.Account, error) {
	return common.Account{}, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// Client - Observations
////////////////////////////////////////////////////////////////////////////////////////

// OnObservedTxIn is called by the observer when a transaction is observed.
func (c *Client) OnObservedTxIn(txIn types.TxInItem, blockHeight int64) {
	// sanity check the transaction has a valid hash
	_, err := chainhash.NewHashFromStr(txIn.Tx)
	if err != nil {
		c.log.Error().Err(err).Str("txID", txIn.Tx).Msg("fail to add spendable utxo to storage")
		return
	}

	blockMeta, err := c.temporalStorage.GetBlockMeta(blockHeight)
	if err != nil {
		c.log.Err(err).Int64("height", blockHeight).Msgf("fail to get block meta")
		return
	}
	if blockMeta == nil {
		blockMeta = utxo.NewBlockMeta("", blockHeight, "")
	}
	if _, err = c.temporalStorage.TrackObservedTx(txIn.Tx); err != nil {
		c.log.Err(err).Msgf("fail to add hash (%s) to observed tx cache", txIn.Tx)
	}
	if c.isAsgardAddress(txIn.Sender) {
		c.log.Debug().Int64("height", blockHeight).Msgf("add hash %s as self transaction", txIn.Tx)
		blockMeta.AddSelfTransaction(txIn.Tx)
	} else {
		// add the transaction to block meta
		blockMeta.AddCustomerTransaction(txIn.Tx)
	}
	if err = c.temporalStorage.SaveBlockMeta(blockHeight, blockMeta); err != nil {
		c.log.Err(err).Int64("height", blockHeight).Msgf("fail to save block meta to storage")
	}
	// update the signer cache
	var m mem.Memo
	m, err = mem.ParseMemo(common.LatestVersion, txIn.Memo)
	if err != nil {
		// Debug log only as ParseMemo error is expected for THORName inbounds.
		c.log.Debug().Err(err).Msgf("fail to parse memo: %s", txIn.Memo)
		return
	}
	if !m.IsOutbound() {
		return
	}
	if m.GetTxID().IsEmpty() {
		return
	}
	if err = c.signerCacheManager.SetSigned(txIn.CacheHash(c.GetChain(), m.GetTxID().String()), txIn.CacheVault(c.GetChain()), txIn.Tx); err != nil {
		c.log.Err(err).Msg("fail to update signer cache")
	}
}

////////////////////////////////////////////////////////////////////////////////////////
// Client - Fetch Block
////////////////////////////////////////////////////////////////////////////////////////

// FetchTxs retrieves txs for a block height. The first argument is the block height to
// fetch, the second argument is the current chain tip.
func (c *Client) FetchTxs(height, chainHeight int64) (types.TxIn, error) {
	txIn := types.TxIn{
		Chain:   c.cfg.ChainID,
		TxArray: nil,
	}

	block, err := c.getBlock(height)
	if err != nil {
		if rpcErr, ok := err.(*btcjson.RPCError); ok && rpcErr.Code == btcjson.ErrRPCInvalidParameter {
			// this means the tx had been broadcast to chain, it must be another signer finished quicker then us
			return txIn, btypes.ErrUnavailableBlock
		}
		return txIn, fmt.Errorf("fail to get block: %w", err)
	}
	// if somehow the block is not valid
	if block.Hash == "" && block.PreviousHash == "" {
		return txIn, fmt.Errorf("fail to get block: %w", err)
	}

	c.currentBlockHeight.Store(height)
	reScannedTxs, err := c.processReorg(block)
	if err != nil {
		c.log.Err(err).Msg("fail to process re-org")
	}
	if len(reScannedTxs) > 0 {
		for _, item := range reScannedTxs {
			if len(item.TxArray) == 0 {
				continue
			}
			txIn.TxArray = append(txIn.TxArray, item.TxArray...)
		}
	}

	blockMeta, err := c.temporalStorage.GetBlockMeta(block.Height)
	if err != nil {
		return txIn, fmt.Errorf("fail to get block meta from storage: %w", err)
	}
	if blockMeta == nil {
		blockMeta = utxo.NewBlockMeta(block.PreviousHash, block.Height, block.Hash)
	} else {
		blockMeta.PreviousHash = block.PreviousHash
		blockMeta.BlockHash = block.Hash
	}

	if err = c.temporalStorage.SaveBlockMeta(block.Height, blockMeta); err != nil {
		return txIn, fmt.Errorf("fail to save block meta into storage: %w", err)
	}
	pruneHeight := height - int64(c.cfg.UTXO.BlockCacheCount)
	if pruneHeight > 0 {
		defer func() {
			if err = c.temporalStorage.PruneBlockMeta(pruneHeight, c.canDeleteBlock); err != nil {
				c.log.Err(err).Int64("height", pruneHeight).Msg("fail to prune block meta")
			}
		}()
	}

	var txInBlock types.TxIn
	txInBlock, err = c.extractTxs(block)
	if err != nil {
		return types.TxIn{}, fmt.Errorf("fail to extract txIn from block: %w", err)
	}
	if len(txInBlock.TxArray) > 0 {
		txIn.TxArray = append(txIn.TxArray, txInBlock.TxArray...)
	}

	c.updateNetworkInfo()

	// report network fee and solvency if within flexibility blocks of tip
	if chainHeight-height <= c.cfg.BlockScanner.ObservationFlexibilityBlocks {
		switch c.cfg.ChainID {
		case common.DOGEChain:
			err = c.sendNetworkFeeFromBlock(block)
		case common.BCHChain, common.LTCChain, common.BTCChain:
			err = c.sendNetworkFee(height)
		default:
			c.log.Fatal().Msg("unsupported chain")
		}
		if err != nil {
			c.log.Err(err).Msg("fail to send network fee")
		}
		// when block scanner is not healthy, only report from auto-unhalt SolvencyCheckRunner
		if c.IsBlockScannerHealthy() {
			if err = c.ReportSolvency(height); err != nil {
				c.log.Err(err).Msg("fail to report solvency info")
			}
		}
	}

	// consolidate UTXOs if there is not one in progress
	if !c.consolidateInProgress.Load() {
		c.wg.Add(1)
		c.consolidateInProgress.Store(true)
		go c.consolidateUTXOs()
	}

	txIn.Count = strconv.Itoa(len(txIn.TxArray))
	return txIn, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// Client - Fetch Mempool
////////////////////////////////////////////////////////////////////////////////////////

// FetchMemPool retrieves txs from mempool
func (c *Client) FetchMemPool(height int64) (types.TxIn, error) {
	hashes, err := c.rpc.GetRawMempool()
	if err != nil {
		return types.TxIn{}, fmt.Errorf("fail to get tx hashes from mempool: %w", err)
	}
	txIn := types.TxIn{
		Chain:   c.GetChain(),
		MemPool: true,
	}

	// shuffle to distribute observations when mempool is large
	rand.Shuffle(len(hashes), func(i, j int) {
		hashes[i], hashes[j] = hashes[j], hashes[i]
	})

	// create batches
	batches := [][]string{}
	batch := []string{}
	for _, h := range hashes {
		// skip transactions we have already processed
		if !c.tryAddToMemPoolCache(h) {
			c.log.Debug().Msgf("ignoring processed tx %s", h)
			continue
		}

		// only process up to the batch size at once
		batch = append(batch, h)
		if len(batch) >= c.cfg.UTXO.TransactionBatchSize {
			batches = append(batches, batch)
			batch = []string{}
		}

		// if we have enough batches, stop
		if len(batches) >= c.cfg.UTXO.MaxMempoolBatches {
			break
		}
	}
	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	// clear mempool cache for batches i or later in case of error
	clearMemPoolCache := func(i int) {
		for j := i; j < len(batches); j++ {
			for _, h := range batches[j] {
				c.removeFromMemPoolCache(h)
			}
		}
	}

	var returnErr error
	errCount := 0
	for i, batch := range batches {
		// fetch the batch of results
		var results []*btcjson.TxRawResult
		var errs []error
		results, errs, err = c.rpc.BatchGetRawTransactionVerbose(batch)
		if err != nil { // clear mempool cache for unprocessed batches and return error
			clearMemPoolCache(i)
			returnErr = fmt.Errorf("fail to get raw transactions from mempool: %w", err)
			break
		}

		// process the batch results
		for i := range results {
			result := results[i]
			err = errs[i]
			// the transaction could have been removed, regardless safe to continue
			if err != nil {
				errCount++
				c.removeFromMemPoolCache(batch[i]) // remove from mempool cache so it will retry
				c.log.Err(err).Str("hash", batch[i]).Msg("fail to get raw transaction verbose")
				continue
			}

			// filter transactions
			var txInItem types.TxInItem
			txInItem, err = c.getTxIn(result, height, true, nil)
			if err != nil {
				c.log.Debug().Err(err).Msg("fail to get TxInItem")
				continue
			}
			if txInItem.IsEmpty() {
				continue
			}
			if txInItem.Coins.IsEmpty() {
				continue
			}

			txIn.TxArray = append(txIn.TxArray, txInItem)
		}
	}

	// log some info if we observed or had errors
	if len(txIn.TxArray) > 0 || errCount > 0 {
		c.log.Info().
			Int("batch", len(batch)).
			Int("txins", len(txIn.TxArray)).
			Int("errors", errCount).
			Msg("retrieved mempool batch")
	}

	txIn.Count = strconv.Itoa(len(txIn.TxArray))
	return txIn, returnErr
}

////////////////////////////////////////////////////////////////////////////////////////
// Client - Confirmation Counting
////////////////////////////////////////////////////////////////////////////////////////

// GetConfirmationCount returns the number of blocks required before processing in
// Thorchain.
func (c *Client) GetConfirmationCount(txIn types.TxIn) int64 {
	// if there are no txs, nothing will be reported
	if len(txIn.TxArray) == 0 {
		return 0
	}

	// transactions tagged as mempool do not need confirmation
	if txIn.MemPool {
		return 0
	}

	// get the block height and confirmation required
	height := txIn.TxArray[0].BlockHeight
	confirm, err := c.getBlockRequiredConfirmation(txIn, height)
	if err != nil {
		c.log.Err(err).Int64("height", height).Msg("fail to get required confirmation")
		return 0
	}

	c.log.Info().Int64("height", height).Msgf("confirmation required: %d", confirm)
	return confirm
}

// ConfirmationCountReady will be called by the observer before sending the txIn to
// Thorchain. It will return true if the scanner height is greater than or equal to the
// observed block height + confirmation required.
// https://medium.com/coinmonks/1confvalue-a-simple-pow-confirmation-rule-of-thumb-a8d9c6c483dd
func (c *Client) ConfirmationCountReady(txIn types.TxIn) bool {
	// if there are no txs, nothing will be reported
	if len(txIn.TxArray) == 0 {
		return true
	}

	// transactions tagged as mempool do not need confirmation
	if txIn.MemPool {
		return true
	}

	// check if we have the necessary number of confirmations
	height := txIn.TxArray[0].BlockHeight
	confirm := txIn.ConfirmationRequired
	ready := (c.currentBlockHeight.Load() - height) >= confirm // every tx already has 1
	c.log.Info().
		Int64("height", height).
		Int64("required", confirm).
		Bool("ready", ready).
		Msg("confirmation count check")

	return ready
}

////////////////////////////////////////////////////////////////////////////////////////
// Client - Solvency
////////////////////////////////////////////////////////////////////////////////////////

// ShouldReportSolvency returns true if solvency should be reported at the given height.
func (c *Client) ShouldReportSolvency(height int64) bool {
	if height-c.lastSolvencyCheckHeight <= 1 {
		return false
	}

	switch c.cfg.ChainID {
	case common.DOGEChain:
		return height%10 == 0
	case common.BCHChain, common.BTCChain:
		return true
	case common.LTCChain:
		return height-c.lastSolvencyCheckHeight > 5 && height%5 == 0
	default:
		c.log.Fatal().Msg("unsupported chain")
		return false
	}
}

// ReportSolvency reports solvency for all asgard vaults.
func (c *Client) ReportSolvency(height int64) error {
	if !c.ShouldReportSolvency(height) {
		return nil
	}

	// fetch all asgard vaults
	asgardVaults, err := c.bridge.GetAsgards()
	if err != nil {
		return fmt.Errorf("fail to get asgards: %w", err)
	}

	currentGasFee := cosmos.NewUint(3 * c.cfg.UTXO.EstimatedAverageTxSize * c.lastFeeRate)

	// report insolvent asgard vaults,
	// or else all if the chain is halted and all are solvent
	msgs := make([]types.Solvency, 0, len(asgardVaults))
	solventMsgs := make([]types.Solvency, 0, len(asgardVaults))
	for i := range asgardVaults {
		var acct common.Account
		acct, err = c.GetAccount(asgardVaults[i].PubKey, nil)
		if err != nil {
			c.log.Err(err).Msgf("fail to get account balance")
			continue
		}

		msg := types.Solvency{
			Height: height,
			Chain:  c.cfg.ChainID,
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
		c.log.Info().
			Stringer("asgard", msgs[i].PubKey).
			Interface("coins", msgs[i].Coins).
			Bool("solvent", solvent).
			Msg("reporting solvency")

		// send solvency to thorchain via global queue consumed by the observer
		select {
		case c.globalSolvencyQueue <- msgs[i]:
		case <-time.After(constants.ThorchainBlockTime):
			c.log.Warn().Msgf("timeout sending solvency to thorchain")
		}
	}

	c.lastSolvencyCheckHeight = height
	return nil
}
