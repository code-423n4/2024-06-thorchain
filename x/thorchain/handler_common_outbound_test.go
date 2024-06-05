package thorchain

import (
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
	"gopkg.in/check.v1"
	. "gopkg.in/check.v1"
)

type HandlerCommonOutboundSuite struct{}

var _ = Suite(&HandlerCommonOutboundSuite{})

func (s *HandlerCommonOutboundSuite) TestIsOutboundFakeGasTX(c *C) {
	coins := common.Coins{
		common.NewCoin(common.ETHAsset, cosmos.NewUint(1)),
	}
	gas := common.Gas{
		{Asset: common.ETHAsset, Amount: cosmos.NewUint(1)},
	}
	fakeGasTx := types.ObservedTx{
		Tx: common.NewTx("123", "0xabc", "0x123", coins, gas, "=:AVAX.AVAX:0x123"),
	}

	c.Assert(isOutboundFakeGasTX(fakeGasTx), Equals, true)

	coins = common.Coins{
		common.NewCoin(common.ETHAsset, cosmos.NewUint(100000)),
	}
	theftTx := types.ObservedTx{
		Tx: common.NewTx("123", "0xabc", "0x123", coins, gas, "=:AVAX.AVAX:0x123"),
	}
	c.Assert(isOutboundFakeGasTX(theftTx), Equals, false)

	coins = common.Coins{
		common.NewCoin(common.BTCAsset, cosmos.NewUint(1)),
	}
	theftTx2 := types.ObservedTx{
		Tx: common.NewTx("123", "0xabc", "0x123", coins, gas, "=:AVAX.AVAX:0x123"),
	}
	c.Assert(isOutboundFakeGasTX(theftTx2), Equals, false)
}

func (s *HandlerCommonOutboundSuite) TestSplitCloutEvenDistribution(c *check.C) {
	clout1 := cosmos.NewUint(50)
	clout2 := cosmos.NewUint(50)
	spent := cosmos.NewUint(60)

	split1, split2 := calcReclaim(clout1, clout2, spent)

	c.Assert(split1.String(), check.Equals, "30")
	c.Assert(split2.String(), check.Equals, "30")
}

func (s *HandlerCommonOutboundSuite) TestSplitCloutExcessSpent(c *check.C) {
	clout1 := cosmos.NewUint(50)
	clout2 := cosmos.NewUint(50)
	spent := cosmos.NewUint(120)

	split1, split2 := calcReclaim(clout1, clout2, spent)

	c.Assert(split1.String(), check.Equals, "50")
	c.Assert(split2.String(), check.Equals, "50")
}

func (s *HandlerCommonOutboundSuite) TestSplitCloutInsufficientFirstClout(c *check.C) {
	clout1 := cosmos.NewUint(20)
	clout2 := cosmos.NewUint(80)
	spent := cosmos.NewUint(60)

	split1, split2 := calcReclaim(clout1, clout2, spent)

	c.Assert(split1.String(), check.Equals, "20")
	c.Assert(split2.String(), check.Equals, "40")
}

func (s *HandlerCommonOutboundSuite) TestSplitCloutInsufficientSecondClout(c *check.C) {
	clout1 := cosmos.NewUint(80)
	clout2 := cosmos.NewUint(20)
	spent := cosmos.NewUint(60)

	split1, split2 := calcReclaim(clout1, clout2, spent)

	c.Assert(split1.String(), check.Equals, "40")
	c.Assert(split2.String(), check.Equals, "20")
}

func (s *HandlerCommonOutboundSuite) TestSplitCloutSpentIsZero(c *check.C) {
	clout1 := cosmos.NewUint(50)
	clout2 := cosmos.NewUint(50)
	spent := cosmos.NewUint(0)

	split1, split2 := calcReclaim(clout1, clout2, spent)

	c.Assert(split1.IsZero(), check.Equals, true)
	c.Assert(split2.IsZero(), check.Equals, true)
}

func (s *HandlerCommonOutboundSuite) TestSplitCloutOneSideIsZero(c *check.C) {
	clout1 := cosmos.NewUint(0)
	clout2 := cosmos.NewUint(100)
	spent := cosmos.NewUint(60)

	split1, split2 := calcReclaim(clout1, clout2, spent)

	c.Assert(split1.IsZero(), check.Equals, true)
	c.Assert(split2.String(), check.Equals, "60")
}

func (s *HandlerCommonOutboundSuite) TestSplitCloutBoundaryCondition(c *check.C) {
	clout1 := cosmos.NewUint(1)
	clout2 := cosmos.NewUint(100000000000)
	spent := cosmos.NewUint(2)

	split1, split2 := calcReclaim(clout1, clout2, spent)

	c.Assert(split1.String(), check.Equals, "1")
	c.Assert(split2.String(), check.Equals, "1")
}

func (s *HandlerCommonOutboundSuite) TestSplitCloutBothCloutsZero(c *check.C) {
	clout1 := cosmos.ZeroUint()
	clout2 := cosmos.ZeroUint()
	spent := cosmos.NewUint(60)

	split1, split2 := calcReclaim(clout1, clout2, spent)

	c.Assert(split1.IsZero(), check.Equals, true)
	c.Assert(split2.IsZero(), check.Equals, true)
}
