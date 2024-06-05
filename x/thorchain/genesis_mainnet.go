//go:build !regtest
// +build !regtest

package thorchain

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
)

func InitGenesis(ctx cosmos.Context, keeper keeper.Keeper, data GenesisState) []abci.ValidatorUpdate {
	return initGenesis(ctx, keeper, data)
}
