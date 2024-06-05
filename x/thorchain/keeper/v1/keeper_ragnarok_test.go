package keeperv1

import (
	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/thornode/common"
)

type KeeperRagnarokSuite struct{}

var _ = Suite(&KeeperRagnarokSuite{})

func (s *KeeperRagnarokSuite) TestVault(c *C) {
	ctx, k := setupKeeperForTest(c)

	height, err := k.GetRagnarokBlockHeight(ctx)
	c.Check(err, IsNil)
	c.Check(height, Equals, int64(0))
	k.SetRagnarokBlockHeight(ctx, 12)
	height, err = k.GetRagnarokBlockHeight(ctx)
	c.Assert(err, IsNil)
	c.Assert(height, Equals, int64(12))

	n, err := k.GetRagnarokNth(ctx)
	c.Check(err, IsNil)
	c.Check(n, Equals, int64(0))

	k.SetRagnarokNth(ctx, 2)
	nth, err := k.GetRagnarokNth(ctx)
	c.Assert(err, IsNil)
	c.Assert(nth, Equals, int64(2))

	p, err := k.GetRagnarokPending(ctx)
	c.Check(err, IsNil)
	c.Check(p, Equals, int64(0))

	k.SetRagnarokPending(ctx, 4)
	pending, err := k.GetRagnarokPending(ctx)
	c.Assert(err, IsNil)
	c.Assert(pending, Equals, int64(4))

	c.Check(k.RagnarokInProgress(ctx), Equals, true)
}

func (s *KeeperRagnarokSuite) TestIsRagnarok(c *C) {
	ctx, k := setupKeeperForTest(c)

	gasAsset := common.ETHAsset
	tokenAsset, err := common.NewAsset("ETH.TOKEN")
	c.Assert(err, IsNil)

	// no ragnarok
	c.Check(k.IsRagnarok(ctx, nil), Equals, false)
	c.Check(k.IsRagnarok(ctx, []common.Asset{gasAsset}), Equals, false)
	c.Check(k.IsRagnarok(ctx, []common.Asset{tokenAsset}), Equals, false)

	// non-gas ragnarok
	k.SetMimir(ctx, "RAGNAROK-ETH-TOKEN", 1)
	c.Check(k.IsRagnarok(ctx, nil), Equals, false)
	c.Check(k.IsRagnarok(ctx, []common.Asset{gasAsset}), Equals, false)
	c.Check(k.IsRagnarok(ctx, []common.Asset{tokenAsset}), Equals, true)

	_ = k.DeleteMimir(ctx, "RAGNAROK-ETH-TOKEN")

	// gas ragnarok
	k.SetMimir(ctx, "RAGNAROK-ETH-ETH", 1)
	c.Check(k.IsRagnarok(ctx, nil), Equals, false)
	c.Check(k.IsRagnarok(ctx, []common.Asset{gasAsset}), Equals, true)
	c.Check(k.IsRagnarok(ctx, []common.Asset{tokenAsset}), Equals, true)

	_ = k.DeleteMimir(ctx, "RAGNAROK-ETH-ETH")

	// no ragnarok
	c.Check(k.IsRagnarok(ctx, nil), Equals, false)
	c.Check(k.IsRagnarok(ctx, []common.Asset{gasAsset}), Equals, false)
	c.Check(k.IsRagnarok(ctx, []common.Asset{tokenAsset}), Equals, false)
}
