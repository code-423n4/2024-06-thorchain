package common

import (
	. "gopkg.in/check.v1"
)

type ChainNetworkSuite struct{}

var _ = Suite(&ChainNetworkSuite{})

func (s *ChainNetworkSuite) TestSoftEquals(c *C) {
	c.Assert(MainNet.SoftEquals(MainNet), Equals, true)
	c.Assert(MockNet.SoftEquals(MockNet), Equals, true)
	c.Assert(MainNet.SoftEquals(MockNet), Equals, false)
}
