package keeperv1

import (
	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/thornode/common"
)

type KeeperHaltSuite struct{}

var _ = Suite(&KeeperHaltSuite{})

func (s *KeeperHaltSuite) TestIsTradingHalt(c *C) {
	ctx, k := setupKeeperForTest(c)

	tx := common.Tx{Coins: common.Coins{common.Coin{Asset: common.BTCAsset}}}
	swapMsg := &MsgSwap{Tx: tx, TargetAsset: common.ETHAsset}
	addMsg := &MsgAddLiquidity{Asset: common.ETHAsset}
	withdrawMsg := &MsgWithdrawLiquidity{Asset: common.ETHAsset}

	// no halts
	c.Check(k.IsTradingHalt(ctx, swapMsg), Equals, false)
	c.Check(k.IsTradingHalt(ctx, addMsg), Equals, false)
	c.Check(k.IsTradingHalt(ctx, withdrawMsg), Equals, false)

	// eth ragnarok
	k.SetMimir(ctx, "RAGNAROK-ETH-ETH", 1)
	c.Check(k.IsTradingHalt(ctx, swapMsg), Equals, true) // target asset
	c.Check(k.IsTradingHalt(ctx, addMsg), Equals, true)
	c.Check(k.IsTradingHalt(ctx, withdrawMsg), Equals, false)

	// synth to l1 bypasses ragnarok check for swaps
	swapMsg.Tx.Coins[0].Asset = common.ETHAsset.GetSyntheticAsset()
	c.Check(k.IsTradingHalt(ctx, swapMsg), Equals, false)
	c.Check(k.IsTradingHalt(ctx, addMsg), Equals, true)
	c.Check(k.IsTradingHalt(ctx, withdrawMsg), Equals, false)

	swapMsg.Tx.Coins[0].Asset = common.BTCAsset
	_ = k.DeleteMimir(ctx, "RAGNAROK-ETH-ETH")

	// btc ragnarok
	k.SetMimir(ctx, "RAGNAROK-BTC-BTC", 1)
	c.Check(k.IsTradingHalt(ctx, swapMsg), Equals, true) // source asset
	c.Check(k.IsTradingHalt(ctx, addMsg), Equals, false)
	c.Check(k.IsTradingHalt(ctx, withdrawMsg), Equals, false)

	_ = k.DeleteMimir(ctx, "RAGNAROK-BTC-BTC")

	// btc chain trading halt
	k.SetMimir(ctx, "HaltBTCTrading", 1)
	c.Check(k.IsTradingHalt(ctx, swapMsg), Equals, true) // source asset
	c.Check(k.IsTradingHalt(ctx, addMsg), Equals, false)
	c.Check(k.IsTradingHalt(ctx, withdrawMsg), Equals, false)

	_ = k.DeleteMimir(ctx, "HaltBTCTrading")

	// eth chain trading halt
	k.SetMimir(ctx, "HaltETHTrading", 1)
	c.Check(k.IsTradingHalt(ctx, swapMsg), Equals, true) // target asset
	c.Check(k.IsTradingHalt(ctx, addMsg), Equals, true)
	c.Check(k.IsTradingHalt(ctx, withdrawMsg), Equals, false)

	_ = k.DeleteMimir(ctx, "HaltETHTrading")

	// global trading halt
	k.SetMimir(ctx, "HaltTrading", 1)
	c.Check(k.IsTradingHalt(ctx, swapMsg), Equals, true)
	c.Check(k.IsTradingHalt(ctx, addMsg), Equals, true)
	c.Check(k.IsTradingHalt(ctx, withdrawMsg), Equals, true)

	_ = k.DeleteMimir(ctx, "HaltTrading")

	// no halts
	c.Check(k.IsTradingHalt(ctx, swapMsg), Equals, false)
	c.Check(k.IsTradingHalt(ctx, addMsg), Equals, false)
	c.Check(k.IsTradingHalt(ctx, withdrawMsg), Equals, false)
}

func (s *KeeperHaltSuite) TestIsGlobalTradingHalted(c *C) {
	ctx, k := setupKeeperForTest(c)
	ctx = ctx.WithBlockHeight(10)

	// no halts
	c.Check(k.IsGlobalTradingHalted(ctx), Equals, false)

	// expired global trading halt
	k.SetMimir(ctx, "HaltTrading", 10)
	c.Check(k.IsGlobalTradingHalted(ctx), Equals, false)

	// current global trading halt
	k.SetMimir(ctx, "HaltTrading", 1)
	c.Check(k.IsGlobalTradingHalted(ctx), Equals, true)

	_ = k.DeleteMimir(ctx, "HaltTrading")

	// no halts
	c.Check(k.IsGlobalTradingHalted(ctx), Equals, false)
}

func (s *KeeperHaltSuite) TestIsChainTradingHalted(c *C) {
	ctx, k := setupKeeperForTest(c)
	ctx = ctx.WithBlockHeight(10)

	// no halts
	c.Check(k.IsChainTradingHalted(ctx, common.BTCChain), Equals, false)
	c.Check(k.IsChainTradingHalted(ctx, common.ETHChain), Equals, false)

	// expired btc trading halt
	k.SetMimir(ctx, "HaltBTCTrading", 10)
	c.Check(k.IsChainTradingHalted(ctx, common.BTCChain), Equals, false)
	c.Check(k.IsChainTradingHalted(ctx, common.ETHChain), Equals, false)

	// current btc trading halt
	k.SetMimir(ctx, "HaltBTCTrading", 1)
	c.Check(k.IsChainTradingHalted(ctx, common.BTCChain), Equals, true)
	c.Check(k.IsChainTradingHalted(ctx, common.ETHChain), Equals, false)

	_ = k.DeleteMimir(ctx, "HaltBTCTrading")

	// current btc chain halt
	k.SetMimir(ctx, "HaltBTCChain", 1)
	c.Check(k.IsChainTradingHalted(ctx, common.BTCChain), Equals, true)
	c.Check(k.IsChainTradingHalted(ctx, common.ETHChain), Equals, false)

	_ = k.DeleteMimir(ctx, "HaltBTCChain")

	// no halts
	c.Check(k.IsChainTradingHalted(ctx, common.BTCChain), Equals, false)
	c.Check(k.IsChainTradingHalted(ctx, common.ETHChain), Equals, false)
}

func (s *KeeperHaltSuite) TestIsChainHalted(c *C) {
	ctx, k := setupKeeperForTest(c)
	ctx = ctx.WithBlockHeight(10)

	// no halts
	c.Check(k.IsChainHalted(ctx, common.BTCChain), Equals, false)
	c.Check(k.IsChainHalted(ctx, common.ETHChain), Equals, false)

	// expired global halt
	k.SetMimir(ctx, "HaltChainGlobal", 10)
	c.Check(k.IsChainHalted(ctx, common.BTCChain), Equals, false)
	c.Check(k.IsChainHalted(ctx, common.ETHChain), Equals, false)

	// current global halt
	k.SetMimir(ctx, "HaltChainGlobal", 1)
	c.Check(k.IsChainHalted(ctx, common.BTCChain), Equals, true)
	c.Check(k.IsChainHalted(ctx, common.ETHChain), Equals, true)

	_ = k.DeleteMimir(ctx, "HaltChainGlobal")

	// expired node pause
	k.SetMimir(ctx, "NodePauseChainGlobal", 1)
	c.Check(k.IsChainHalted(ctx, common.BTCChain), Equals, false)
	c.Check(k.IsChainHalted(ctx, common.ETHChain), Equals, false)

	// current node pause
	k.SetMimir(ctx, "NodePauseChainGlobal", 11)
	c.Check(k.IsChainHalted(ctx, common.BTCChain), Equals, true)
	c.Check(k.IsChainHalted(ctx, common.ETHChain), Equals, true)

	_ = k.DeleteMimir(ctx, "NodePauseChainGlobal")

	// expired btc halt
	k.SetMimir(ctx, "HaltBTCChain", 10)
	c.Check(k.IsChainHalted(ctx, common.BTCChain), Equals, false)
	c.Check(k.IsChainHalted(ctx, common.ETHChain), Equals, false)

	// current btc halt
	k.SetMimir(ctx, "HaltBTCChain", 1)
	c.Check(k.IsChainHalted(ctx, common.BTCChain), Equals, true)
	c.Check(k.IsChainHalted(ctx, common.ETHChain), Equals, false)

	_ = k.DeleteMimir(ctx, "HaltBTCChain")

	// expired btc solvency halt
	k.SetMimir(ctx, "SolvencyHaltBTCChain", 10)
	c.Check(k.IsChainHalted(ctx, common.BTCChain), Equals, false)
	c.Check(k.IsChainHalted(ctx, common.ETHChain), Equals, false)

	// current btc solvency halt
	k.SetMimir(ctx, "SolvencyHaltBTCChain", 1)
	c.Check(k.IsChainHalted(ctx, common.BTCChain), Equals, true)
	c.Check(k.IsChainHalted(ctx, common.ETHChain), Equals, false)

	_ = k.DeleteMimir(ctx, "SolvencyHaltBTCChain")

	// no halts
	c.Check(k.IsChainHalted(ctx, common.BTCChain), Equals, false)
	c.Check(k.IsChainHalted(ctx, common.ETHChain), Equals, false)
}

func (s *KeeperHaltSuite) TestIsLPPaused(c *C) {
	ctx, k := setupKeeperForTest(c)
	ctx = ctx.WithBlockHeight(10)

	// no pauses
	c.Check(k.IsLPPaused(ctx, common.BTCChain), Equals, false)
	c.Check(k.IsLPPaused(ctx, common.ETHChain), Equals, false)

	// expired btc pause
	k.SetMimir(ctx, "PauseLPBTC", 10)
	c.Check(k.IsLPPaused(ctx, common.BTCChain), Equals, false)
	c.Check(k.IsLPPaused(ctx, common.ETHChain), Equals, false)

	// current btc pause
	k.SetMimir(ctx, "PauseLPBTC", 1)
	c.Check(k.IsLPPaused(ctx, common.BTCChain), Equals, true)
	c.Check(k.IsLPPaused(ctx, common.ETHChain), Equals, false)

	_ = k.DeleteMimir(ctx, "PauseLPBTC")

	// expired global pause
	k.SetMimir(ctx, "PauseLP", 10)
	c.Check(k.IsLPPaused(ctx, common.BTCChain), Equals, false)
	c.Check(k.IsLPPaused(ctx, common.ETHChain), Equals, false)

	// current global pause
	k.SetMimir(ctx, "PauseLP", 1)
	c.Check(k.IsLPPaused(ctx, common.BTCChain), Equals, true)
	c.Check(k.IsLPPaused(ctx, common.ETHChain), Equals, true)

	_ = k.DeleteMimir(ctx, "PauseLP")

	// no pauses
	c.Check(k.IsLPPaused(ctx, common.BTCChain), Equals, false)
}
