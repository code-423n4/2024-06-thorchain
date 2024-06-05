package thorchain

import (
	"gitlab.com/thorchain/thornode/common"
	. "gopkg.in/check.v1"
)

type SwapperCloutInitTHORNamesTestSuite struct{}

var _ = Suite(&SwapperCloutInitTHORNamesTestSuite{})

func (s *StoreManagerTestSuite) TestParseInitCloutTHORNames(c *C) {
	for _, item := range getInitCloutTHORNames() {
		_, err := common.NewAddress(item.address)
		c.Assert(err, IsNil)
	}
}
