package keeperv1

import (
	. "gopkg.in/check.v1"
)

type KeeperMimirV2Suite struct{}

var _ = Suite(&KeeperMimirV2Suite{})

func (s *KeeperMimirV2Suite) TestMimirV2(c *C) {
	ctx, k := setupKeeperForTest(c)
	acc := GetRandomBech32Addr()

	c.Assert(k.SetNodeMimirV2(ctx, "test", 15, acc), IsNil)
	mimirs, err := k.GetNodeMimirsV2(ctx, "test")
	c.Assert(err, IsNil)
	c.Assert(mimirs.Mimirs, HasLen, 1)
	m := mimirs.Mimirs[0]
	c.Check(m.Key, Equals, "TEST")
	c.Check(m.Value, Equals, int64(15))
	c.Check(m.Signer.Equals(acc), Equals, true)
}
