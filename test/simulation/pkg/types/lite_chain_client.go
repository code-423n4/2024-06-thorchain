package types

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"gitlab.com/thorchain/thornode/bifrost/thorclient"
	"gitlab.com/thorchain/thornode/common"
)

// SimTx is a struct used for simulation transactions.
type SimTx struct {
	Chain     common.Chain
	ToAddress common.Address
	Coin      common.Coin
	Memo      string
}

// SimContractTx is a struct used for simulation contract call transactions.
type SimContractTx struct {
	Chain    common.Chain
	Contract common.Address
	ABI      abi.ABI
	Method   string
	Args     []interface{}
}

// LiteChainClient is a subset of the ChainClient interface used for simulation tests.
type LiteChainClient interface {
	// GetAccount returns the account for the given public key. If the key is nil, it
	// returns the account for the client's configured key.
	GetAccount(_ *common.PubKey) (*common.Account, error)

	// SignTx returns the signed transaction.
	SignTx(tx SimTx) ([]byte, error)

	// BroadcastTx broadcasts the transaction and returns the hash.
	BroadcastTx([]byte) (string, error)
}

// LiteChainClientConstructor is a function that creates a new LiteChainClient.
type LiteChainClientConstructor func(_ common.Chain, _ *thorclient.Keys) (LiteChainClient, error)
