package keeperv1

import (
	. "gopkg.in/check.v1"

	cosmos "gitlab.com/thorchain/thornode/common/cosmos"
)

type KeeperStreamingSwapSuite struct{}

var _ = Suite(&KeeperStreamingSwapSuite{})

func (mas *KeeperStreamingSwapSuite) SetUpSuite(c *C) {
	SetupConfigForTest()
}

func (s *KeeperStreamingSwapSuite) TestStreamingSwap(c *C) {
	ctx, k := setupKeeperForTest(c)
	txID := GetRandomTxHash()

	_, err := k.GetStreamingSwap(ctx, txID)
	c.Assert(err, IsNil)

	swp := NewStreamingSwap(txID, 10, 20, cosmos.NewUint(13), cosmos.NewUint(1000))
	k.SetStreamingSwap(ctx, swp)
	swp, err = k.GetStreamingSwap(ctx, txID)
	c.Assert(err, IsNil)
	c.Check(swp.Quantity, Equals, uint64(10))
	c.Check(swp.Interval, Equals, uint64(20))
	c.Check(swp.TradeTarget.String(), Equals, "13")
	c.Check(swp.Deposit.String(), Equals, "1000")
	iter := k.GetStreamingSwapIterator(ctx)
	c.Check(iter, NotNil)
	iter.Close()
	k.RemoveStreamingSwap(ctx, txID)
}
