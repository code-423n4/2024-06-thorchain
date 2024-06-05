package thorchain

import (
	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
)

// using int64 so this can also represent deltas
type ModuleBalances struct {
	Asgard  int64
	Bond    int64
	Reserve int64
	Module  int64
}

func getModuleBalances(c *C, ctx cosmos.Context, k keeper.Keeper) ModuleBalances {
	return ModuleBalances{
		Asgard:  int64(k.GetRuneBalanceOfModule(ctx, AsgardName).Uint64()),
		Bond:    int64(k.GetRuneBalanceOfModule(ctx, BondName).Uint64()),
		Reserve: int64(k.GetRuneBalanceOfModule(ctx, ReserveName).Uint64()),
		Module:  int64(k.GetRuneBalanceOfModule(ctx, ModuleName).Uint64()),
	}
}

func testAndCheckModuleBalances(c *C, ctx cosmos.Context, k keeper.Keeper, runTest func(), expDeltas ModuleBalances) {
	before := getModuleBalances(c, ctx, k)
	runTest()
	after := getModuleBalances(c, ctx, k)

	c.Assert(expDeltas.Asgard, Equals, after.Asgard-before.Asgard)
	c.Assert(expDeltas.Bond, Equals, after.Bond-before.Bond)
	c.Assert(expDeltas.Reserve, Equals, after.Reserve-before.Reserve)
	c.Assert(expDeltas.Module, Equals, after.Module-before.Module)
}
