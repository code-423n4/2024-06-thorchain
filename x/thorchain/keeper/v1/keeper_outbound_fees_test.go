package keeperv1

import (
	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

type KeeperOutboundFeesSuite struct{}

var _ = Suite(&KeeperLiquidityFeesSuite{})

func (s *KeeperOutboundFeesSuite) TestOutboundRuneRecords(c *C) {
	ctx, k := setupKeeperForTest(c)

	// Nothing set returns 0.
	feeWithheldRune, err := k.GetOutboundFeeWithheldRune(ctx, common.BTCAsset)
	c.Assert(err, IsNil)
	feeSpentRune, err := k.GetOutboundFeeSpentRune(ctx, common.BTCAsset)
	c.Assert(err, IsNil)
	c.Check(feeWithheldRune.String(), Equals, "0")
	c.Check(feeSpentRune.String(), Equals, "0")

	// Adding sets.
	err = k.AddToOutboundFeeWithheldRune(ctx, common.BTCAsset, cosmos.NewUint(uint64(200)))
	c.Assert(err, IsNil)
	err = k.AddToOutboundFeeSpentRune(ctx, common.BTCAsset, cosmos.NewUint(uint64(100)))
	c.Assert(err, IsNil)

	feeWithheldRune, err = k.GetOutboundFeeWithheldRune(ctx, common.BTCAsset)
	c.Assert(err, IsNil)
	feeSpentRune, err = k.GetOutboundFeeSpentRune(ctx, common.BTCAsset)
	c.Assert(err, IsNil)
	c.Check(feeWithheldRune.String(), Equals, "200")
	c.Check(feeSpentRune.String(), Equals, "100")

	// Adding again adds.
	err = k.AddToOutboundFeeWithheldRune(ctx, common.BTCAsset, cosmos.NewUint(uint64(400)))
	c.Assert(err, IsNil)
	err = k.AddToOutboundFeeSpentRune(ctx, common.BTCAsset, cosmos.NewUint(uint64(300)))
	c.Assert(err, IsNil)

	feeWithheldRune, err = k.GetOutboundFeeWithheldRune(ctx, common.BTCAsset)
	c.Assert(err, IsNil)
	feeSpentRune, err = k.GetOutboundFeeSpentRune(ctx, common.BTCAsset)
	c.Assert(err, IsNil)
	c.Check(feeWithheldRune.String(), Equals, "600")
	c.Check(feeSpentRune.String(), Equals, "400")

	// Set values are distinct by Asset.
	feeWithheldRune, err = k.GetOutboundFeeWithheldRune(ctx, common.BNBAsset)
	c.Assert(err, IsNil)
	feeSpentRune, err = k.GetOutboundFeeSpentRune(ctx, common.BNBAsset)
	c.Assert(err, IsNil)
	c.Check(feeWithheldRune.String(), Equals, "0")
	c.Check(feeSpentRune.String(), Equals, "0")

	err = k.AddToOutboundFeeWithheldRune(ctx, common.BNBAsset, cosmos.NewUint(uint64(50)))
	c.Assert(err, IsNil)
	err = k.AddToOutboundFeeSpentRune(ctx, common.BTCAsset, cosmos.NewUint(uint64(30)))
	c.Assert(err, IsNil)

	feeWithheldRune, err = k.GetOutboundFeeWithheldRune(ctx, common.BTCAsset)
	c.Assert(err, IsNil)
	feeSpentRune, err = k.GetOutboundFeeSpentRune(ctx, common.BTCAsset)
	c.Assert(err, IsNil)
	c.Check(feeWithheldRune.String(), Equals, "600")
	c.Check(feeSpentRune.String(), Equals, "400")

	feeWithheldRune, err = k.GetOutboundFeeWithheldRune(ctx, common.BNBAsset)
	c.Assert(err, IsNil)
	feeSpentRune, err = k.GetOutboundFeeSpentRune(ctx, common.BNBAsset)
	c.Assert(err, IsNil)
	c.Check(feeWithheldRune.String(), Equals, "50")
	c.Check(feeSpentRune.String(), Equals, "30")
}
