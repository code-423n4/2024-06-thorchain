package thorchain

import (
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
)

type gasManagerTestHelper struct {
	keeper.Keeper
	failGetNetwork bool
	failGetPool    bool
	failSetPool    bool
}

func newGasManagerTestHelper(k keeper.Keeper) *gasManagerTestHelper {
	return &gasManagerTestHelper{
		Keeper: k,
	}
}

func (g *gasManagerTestHelper) GetNetwork(ctx cosmos.Context) (Network, error) {
	if g.failGetNetwork {
		return Network{}, errKaboom
	}
	return g.Keeper.GetNetwork(ctx)
}

func (g *gasManagerTestHelper) GetPool(ctx cosmos.Context, asset common.Asset) (Pool, error) {
	if g.failGetPool {
		return NewPool(), errKaboom
	}
	return g.Keeper.GetPool(ctx, asset)
}

func (g *gasManagerTestHelper) SetPool(ctx cosmos.Context, p Pool) error {
	if g.failSetPool {
		return errKaboom
	}
	return g.Keeper.SetPool(ctx, p)
}
