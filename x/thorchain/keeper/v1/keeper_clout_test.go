package keeperv1

import (
	"gitlab.com/thorchain/thornode/common/cosmos"
	. "gopkg.in/check.v1"
)

type KeeperCloutSuite struct{}

var _ = Suite(&KeeperCloutSuite{})

func (s *KeeperCloutSuite) TestSwapperClout(c *C) {
	ctx, k := setupKeeperForTest(c)
	var err error

	addr := GetRandomTHORAddress()
	clout := NewSwapperClout(addr)
	clout.Score = cosmos.NewUint(1000)

	c.Assert(k.SetSwapperClout(ctx, clout), IsNil)

	clout, err = k.GetSwapperClout(ctx, addr)
	c.Assert(err, IsNil)
	c.Assert(clout.Address.String(), Equals, addr.String())
	c.Assert(clout.Score.String(), Equals, "1000")
}
