package evm

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"time"

	_ "embed"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"

	"gitlab.com/thorchain/thornode/bifrost/pkg/chainclients/shared/evm"
	"gitlab.com/thorchain/thornode/bifrost/thorclient"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/tokenlist"

	. "gitlab.com/thorchain/thornode/test/simulation/pkg/types"
)

////////////////////////////////////////////////////////////////////////////////////////
// Config
////////////////////////////////////////////////////////////////////////////////////////

const ContractGasLimit = 3000000

////////////////////////////////////////////////////////////////////////////////////////
// Init
////////////////////////////////////////////////////////////////////////////////////////

//go:embed abi/router.json
var routerABIJson string

//go:embed abi/erc20.json
var erc20ABIJson string

var routerABI, erc20ABI abi.ABI

func init() {
	var err error
	routerABI, err = abi.JSON(strings.NewReader(routerABIJson))
	if err != nil {
		panic(fmt.Errorf("failed to parse router contract abi: %w", err))
	}

	erc20ABI, err = abi.JSON(strings.NewReader(erc20ABIJson))
	if err != nil {
		panic(fmt.Errorf("failed to parse erc20 contract abi: %w", err))
	}
}

////////////////////////////////////////////////////////////////////////////////////////
// ABI
////////////////////////////////////////////////////////////////////////////////////////
//

func RouterABI() abi.ABI {
	return routerABI
}

func ERC20ABI() abi.ABI {
	return erc20ABI
}

////////////////////////////////////////////////////////////////////////////////////////
// Tokens
////////////////////////////////////////////////////////////////////////////////////////

// Tokens returns the list of tokens that are used in the simulation. All tokens will be
// looked up and included in the GetAccount response. This method can be replicated with
// build tag scope for testing different environments in the future.
func Tokens(chain common.Chain) map[common.Asset]tokenlist.ERC20Token {
	tokenMap := make(map[common.Asset]tokenlist.ERC20Token)

	// gather the available tokens
	var tokens []tokenlist.ERC20Token
	switch chain {
	case common.ETHChain, common.BSCChain:
		tokens = []tokenlist.ERC20Token{
			{
				Address:  "0x52C84043CD9c865236f11d9Fc9F56aa003c1f922",
				Symbol:   "TKN",
				Decimals: 18,
			},
		}
	case common.AVAXChain:
		tokens = []tokenlist.ERC20Token{
			{
				Address:  "0x17aB05351fC94a1a67Bf3f56DdbB941aE6c63E25",
				Symbol:   "TKN",
				Decimals: 18,
			},
		}
	}

	// create mapping of asset to token
	for _, token := range tokens {
		tokenMap[token.Asset(chain)] = token
	}

	return tokenMap
}

////////////////////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////////////////////

func ctx() context.Context {
	return context.Background()
}

////////////////////////////////////////////////////////////////////////////////////////
// Client
////////////////////////////////////////////////////////////////////////////////////////

type Client struct {
	chain common.Chain
	rpc   *ethclient.Client

	keys    *thorclient.Keys
	privKey *ecdsa.PrivateKey
	signer  etypes.EIP155Signer
	pubKey  common.PubKey
	address common.Address
}

var _ LiteChainClient = &Client{}

func NewConstructor(host string) LiteChainClientConstructor {
	return func(chain common.Chain, keys *thorclient.Keys) (LiteChainClient, error) {
		return NewClient(chain, host, keys)
	}
}

func NewClient(chain common.Chain, host string, keys *thorclient.Keys) (LiteChainClient, error) {
	// extract the private key
	privateKey, err := keys.GetPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("fail to get private key: %w", err)
	}
	privKey, err := evm.GetPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	// derive the public key
	pk, err := cryptocodec.ToTmPubKeyInterface(privateKey.PubKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get tm pub key: %w", err)
	}
	pubKey, err := common.NewPubKeyFromCrypto(pk)
	if err != nil {
		return nil, fmt.Errorf("fail to create pubkey: %w", err)
	}

	// get pubkey address for the chain
	address, err := pubKey.GetAddress(chain)
	if err != nil {
		return nil, fmt.Errorf("fail to get address from pubkey(%s): %w", pk, err)
	}

	// dial the rpc host
	rpc, err := ethclient.Dial(host)
	if err != nil {
		return nil, fmt.Errorf("fail to dial ETH rpc host(%s): %w", host, err)
	}

	// get the chain id
	chainID, err := rpc.ChainID(ctx())
	if err != nil {
		return nil, fmt.Errorf("fail to get chain id: %w", err)
	}

	// create the signer
	signer := etypes.NewEIP155Signer(chainID)

	return &Client{
		chain:   chain,
		rpc:     rpc,
		keys:    keys,
		privKey: privKey,
		signer:  signer,
		pubKey:  pubKey,
		address: address,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// GetAccount
////////////////////////////////////////////////////////////////////////////////////////

func (c *Client) GetAccount(pk *common.PubKey) (*common.Account, error) {
	// get nonce
	nonce, err := c.rpc.PendingNonceAt(ctx(), ecommon.HexToAddress(c.address.String()))
	if err != nil {
		return nil, fmt.Errorf("fail to get account nonce: %w", err)
	}

	// get balance
	balance, err := c.rpc.BalanceAt(ctx(), ecommon.HexToAddress(c.address.String()), nil)
	if err != nil {
		return nil, fmt.Errorf("fail to get account balance: %w", err)
	}

	// get amount
	amount := sdk.NewUintFromBigInt(balance)
	amount = amount.Quo(sdk.NewUint(1e10)) // 1e18 -> 1e8

	// add gas asset to coins
	coins := common.Coins{
		common.NewCoin(c.chain.GetGasAsset(), amount),
	}

	// lookup any other tokens
	for asset, token := range Tokens(c.chain) {
		// get balance
		abi := ERC20ABI()
		data, err := abi.Pack("balanceOf", ecommon.HexToAddress(c.address.String()))
		if err != nil {
			log.Error().Err(err).Msg("error packing balanceOf")
			continue
		}
		to := ecommon.HexToAddress(token.Address)
		result, err := c.rpc.CallContract(ctx(), ethereum.CallMsg{
			To:   &to,
			Data: data,
		}, nil)
		if err != nil {
			log.Error().Err(err).Msg("error calling contract")
			continue
		}
		balance := new(big.Int)
		balance.SetBytes(result)

		// convert balance from decimals to 1e8
		balance.Mul(balance, big.NewInt(common.One))
		balance.Div(balance, big.NewInt(1).Exp(big.NewInt(10), big.NewInt(int64(token.Decimals)), nil))

		// add to coins
		coins = append(coins, common.NewCoin(asset, sdk.NewUintFromBigInt(balance)))
	}

	// create account
	return &common.Account{
		Sequence: int64(nonce),
		Coins:    coins,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// SignTx
////////////////////////////////////////////////////////////////////////////////////////

func (c *Client) SignTx(tx SimTx) ([]byte, error) {
	// to address to evm address
	toAddress := ecommon.HexToAddress(tx.ToAddress.String())

	// create a standard transfer tx
	txData := &etypes.LegacyTx{
		To:    &toAddress,
		Data:  []byte(tx.Memo),
		Gas:   21000 + 3000,                                   // standard transfer + memo
		Value: tx.Coin.Amount.Mul(sdk.NewUint(1e10)).BigInt(), // 1e8 -> 1e18,
	}

	return c.signTx(txData)
}

// signTx is shared by the base SignTx and EVM custom SignContractTx methods.
func (c *Client) signTx(txData *etypes.LegacyTx) ([]byte, error) {
	// get nonce
	nonce, err := c.rpc.PendingNonceAt(ctx(), ecommon.HexToAddress(c.address.String()))
	if err != nil {
		return nil, fmt.Errorf("fail to get account nonce: %w", err)
	}

	// get gas price
	gasPrice, err := c.rpc.SuggestGasPrice(ctx())
	if err != nil {
		return nil, fmt.Errorf("fail to get gas price: %w", err)
	}

	// set nonce and gas price
	txData.Nonce = nonce
	txData.GasPrice = gasPrice

	// create signable tx
	signable := etypes.NewTx(txData)

	// sign the tx
	hash := c.signer.Hash(signable)
	sig, err := ecrypto.Sign(hash[:], c.privKey)
	if err != nil {
		return nil, fmt.Errorf("fail to sign tx: %w", err)
	}

	// apply the signature
	newTx, err := signable.WithSignature(c.signer, sig)
	if err != nil {
		return nil, fmt.Errorf("fail to apply signature to tx: %w", err)
	}

	// marshal and return
	enc, err := newTx.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("fail to marshal tx to json: %w", err)
	}

	return enc, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// BroadcastTx
////////////////////////////////////////////////////////////////////////////////////////

func (c *Client) BroadcastTx(signed []byte) (string, error) {
	tx := &etypes.Transaction{}
	if err := tx.UnmarshalJSON(signed); err != nil {
		return "", err
	}
	txid := tx.Hash().String()

	// remove 0x prefix
	txid = strings.TrimPrefix(txid, "0x")

	// send the transaction
	err := c.rpc.SendTransaction(ctx(), tx)
	if err != nil {
		return txid, err
	}

	// wait for the transaction receipt
	var receipt *etypes.Receipt
	for i := 0; i < 10; i++ {
		receipt, err = c.rpc.TransactionReceipt(ctx(), tx.Hash())
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return txid, fmt.Errorf("fail to get transaction receipt: %w", err)
	}

	// check the status and return logs in error
	if receipt.Status != 1 {
		return txid, fmt.Errorf("transaction failed with status %d: %v", receipt.Status, receipt.Logs)
	}

	return txid, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// Custom EVM Client Methods
////////////////////////////////////////////////////////////////////////////////////////

func (c *Client) SignContractTx(tx SimContractTx) ([]byte, error) {
	// contract address to evm address
	contractAddress := ecommon.HexToAddress(tx.Contract.String())
	data, err := tx.ABI.Pack(tx.Method, tx.Args...)
	if err != nil {
		return nil, fmt.Errorf("fail to pack contract call: %w", err)
	}

	// create the tx
	txData := &etypes.LegacyTx{
		To:   &contractAddress,
		Data: data,
		Gas:  ContractGasLimit,
	}

	return c.signTx(txData)
}

func (c *Client) GetTokenDecimals(address string) (int, error) {
	// build contract read call
	addr := ecommon.HexToAddress(address)
	abi := ERC20ABI()
	data, err := abi.Pack("decimals")
	if err != nil {
		return 0, fmt.Errorf("fail to pack decimals call: %w", err)
	}

	// read the contract
	result, err := c.rpc.CallContract(ctx(), ethereum.CallMsg{
		To:   &addr,
		Data: data,
	}, nil)
	if err != nil {
		return 0, fmt.Errorf("fail to call contract: %w", err)
	}

	// extract the decimals
	decimals := new(big.Int)
	decimals.SetBytes(result)
	return int(decimals.Uint64()), nil
}
