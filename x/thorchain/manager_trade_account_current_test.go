package thorchain

import (
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	. "gopkg.in/check.v1"
)

type TradeManagerVCURSuite struct{}

var _ = Suite(&TradeManagerVCURSuite{})

func (s *TradeManagerVCURSuite) SetUpSuite(_ *C) {
	SetupConfigForTest()
}

func (s *TradeManagerVCURSuite) TestDepositAndWithdrawal(c *C) {
	ctx, k := setupKeeperForTest(c)
	eventMgr, err := GetEventManager(GetCurrentVersion())
	c.Assert(err, IsNil)
	mgr := newTradeMgrVCUR(k, eventMgr)

	asset := common.BTCAsset.GetTradeAsset()
	addr1 := GetRandomBech32Addr()
	addr2 := GetRandomBech32Addr()
	// addr3 := GetRandomBech32Addr()

	amt, err := mgr.Deposit(ctx, asset, cosmos.NewUint(100*common.One), addr1, common.NoAddress, common.BlankTxID)
	c.Assert(err, IsNil)
	c.Check(amt.String(), Equals, cosmos.NewUint(100*common.One).String())

	bal := mgr.BalanceOf(ctx, asset, addr1)
	c.Check(bal.String(), Equals, cosmos.NewUint(100*common.One).String())

	amt, err = mgr.Deposit(ctx, asset, cosmos.NewUint(50*common.One), addr2, common.NoAddress, common.BlankTxID)
	c.Assert(err, IsNil)
	c.Check(amt.String(), Equals, cosmos.NewUint(50*common.One).String())

	bal = mgr.BalanceOf(ctx, asset, addr2)
	c.Check(bal.String(), Equals, cosmos.NewUint(50*common.One).String())
	bal = mgr.BalanceOf(ctx, asset, addr1)
	c.Check(bal.String(), Equals, cosmos.NewUint(100*common.One).String())

	// withdrawal
	amt, err = mgr.Withdrawal(ctx, asset, cosmos.NewUint(30*common.One), addr2, common.NoAddress, common.BlankTxID)
	c.Assert(err, IsNil)
	c.Check(amt.String(), Equals, cosmos.NewUint(30*common.One).String())
	bal = mgr.BalanceOf(ctx, asset, addr2)
	c.Check(bal.String(), Equals, cosmos.NewUint(20*common.One).String())
	amt, err = mgr.Withdrawal(ctx, asset, cosmos.NewUint(30*common.One), addr2, common.NoAddress, common.BlankTxID)
	c.Assert(err, IsNil)
	c.Check(amt.String(), Equals, cosmos.NewUint(20*common.One).String())
	bal = mgr.BalanceOf(ctx, asset, addr2)
	c.Check(bal.String(), Equals, cosmos.NewUint(0).String())
}

func (s *TradeManagerVCURSuite) TestReducedDepth(c *C) {
	ctx, k := setupKeeperForTest(c)
	eventMgr, err := GetEventManager(GetCurrentVersion())
	c.Assert(err, IsNil)
	mgr := newTradeMgrVCUR(k, eventMgr)

	asset := common.BTCAsset.GetTradeAsset()
	addr1 := GetRandomBech32Addr()

	depositDepth := cosmos.NewUint(100 * common.One)

	amt, err := mgr.Deposit(ctx, asset, depositDepth, addr1, common.NoAddress, common.BlankTxID)
	c.Assert(err, IsNil)
	c.Check(amt.String(), Equals, depositDepth.String())

	bal := mgr.BalanceOf(ctx, asset, addr1)
	c.Check(bal.String(), Equals, depositDepth.String())

	tu, err := k.GetTradeUnit(ctx, asset)
	c.Assert(err, IsNil)
	c.Check(tu.Depth.String(), Equals, depositDepth.String())
	c.Check(tu.Units.String(), Equals, depositDepth.String())

	// Halve the depth, as though experiencing negative interest.
	tu.Depth = tu.Depth.QuoUint64(2)
	k.SetTradeUnit(ctx, tu)

	bal = mgr.BalanceOf(ctx, asset, addr1)
	c.Check(bal.String(), Equals, tu.Depth.String())

	// Attempt to withdraw half of the remaining depth
	// which should also reduce the units by half.
	previousDepth := tu.Depth
	previousUnits := tu.Units
	amt, err = mgr.Withdrawal(ctx, asset, previousDepth.QuoUint64(2), addr1, common.NoAddress, common.BlankTxID)
	c.Assert(err, IsNil)
	c.Check(amt.String(), Equals, previousDepth.QuoUint64(2).String())

	tu, err = k.GetTradeUnit(ctx, asset)
	c.Assert(err, IsNil)
	c.Check(tu.Depth.String(), Equals, previousDepth.QuoUint64(2).String())
	c.Check(tu.Units.String(), Equals, previousUnits.QuoUint64(2).String())

	// Now deposit to double the tu.Depth, which should double the tu.Units .
	previousDepth = tu.Depth
	previousUnits = tu.Units
	amt, err = mgr.Deposit(ctx, asset, previousDepth, addr1, common.NoAddress, common.BlankTxID)
	c.Assert(err, IsNil)
	c.Check(amt.String(), Equals, previousDepth.String())
	// 'amt' is always returned by Deposit as-is with no modification.
	tu, err = k.GetTradeUnit(ctx, asset)
	c.Assert(err, IsNil)
	c.Check(tu.Depth.String(), Equals, previousDepth.MulUint64(2).String())
	c.Check(tu.Units.String(), Equals, previousUnits.MulUint64(2).String())

	// Confirm that the address's units are consistent with the total units.
	tr, err := k.GetTradeAccount(ctx, addr1, asset)
	c.Assert(err, IsNil)
	c.Check(tr.Units.String(), Equals, tu.Units.String())
}
