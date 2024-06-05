package types

import (
	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

type StreamingSwapSuite struct{}

var _ = Suite(&StreamingSwapSuite{})

func (s *StreamingSwapSuite) TestNextSize(c *C) {
	v := GetCurrentVersion()

	swp := NewStreamingSwap(common.BlankTxID, 2, 10, cosmos.NewUint(10), cosmos.NewUint(10))
	size, target := swp.NextSize(v)
	c.Check(size.String(), Equals, "5")
	c.Check(target.String(), Equals, "5")
	swp.In = cosmos.NewUint(5)
	swp.Out = cosmos.NewUint(5)
	size, target = swp.NextSize(v)
	c.Check(size.String(), Equals, "5")
	c.Check(target.String(), Equals, "5")
	swp.In = cosmos.NewUint(10)
	swp.Out = cosmos.NewUint(10)
	size, target = swp.NextSize(v)
	c.Check(size.String(), Equals, "0")
	c.Check(target.String(), Equals, "0")

	swp.Quantity = 10
	swp.Deposit = cosmos.NewUint(100)
	swp.TradeTarget = cosmos.NewUint(100)
	swp.In = cosmos.NewUint(0)
	swp.Out = cosmos.NewUint(0)
	size, target = swp.NextSize(v)
	c.Check(size.String(), Equals, "10")
	c.Check(target.String(), Equals, "10")

	swp.In = cosmos.NewUint(10)
	swp.Out = cosmos.NewUint(20)
	size, target = swp.NextSize(v)
	c.Check(size.String(), Equals, "10")
	c.Check(target.String(), Equals, "9")

	swp.In = cosmos.NewUint(20)
	swp.Out = cosmos.NewUint(40)
	size, target = swp.NextSize(v)
	c.Check(size.String(), Equals, "10")
	c.Check(target.String(), Equals, "8")

	swp.In = cosmos.NewUint(30)
	swp.Out = cosmos.NewUint(60)
	size, target = swp.NextSize(v)
	c.Check(size.String(), Equals, "10")
	c.Check(target.String(), Equals, "6")

	// test no remainder is left behind
	swp = NewStreamingSwap(common.BlankTxID, 5, 10, cosmos.NewUint(2345), cosmos.NewUint(472659))
	total := cosmos.ZeroUint()
	for i := 1; i <= 5; i++ {
		size, _ = swp.NextSize(v)
		total = total.Add(size)
		swp.Count += 1
		swp.In = swp.In.Add(size)
	}
	c.Check(total.String(), Equals, "472659")
}

func (s *StreamingSwapSuite) TestValidate(c *C) {
	// happy path
	swp := NewStreamingSwap(common.BlankTxID, 2, 10, cosmos.NewUint(10), cosmos.NewUint(10))
	c.Assert(swp.Valid(), IsNil)

	// non-happy path
	swp.Quantity = 0
	c.Assert(swp.Valid(), NotNil)
	swp.Quantity = 2
	swp.Interval = 0
	c.Assert(swp.Valid(), NotNil)
	swp.Interval = 10
	swp.Deposit = cosmos.ZeroUint()
	c.Assert(swp.Valid(), NotNil)
}

func (s *StreamingSwapSuite) TestIsDone(c *C) {
	// happy path
	swp := NewStreamingSwap(common.BlankTxID, 10, 10, cosmos.NewUint(10), cosmos.NewUint(10))

	swp.Count = 5
	c.Check(swp.IsDone(), Equals, false)
	swp.Count = 10
	c.Check(swp.IsDone(), Equals, true)
	swp.Count = 20
	c.Check(swp.IsDone(), Equals, true)
}
