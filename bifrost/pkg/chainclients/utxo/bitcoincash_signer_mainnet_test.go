//go:build !stagenet && !mocknet
// +build !stagenet,!mocknet

package utxo

import (
	"github.com/gcash/bchd/chaincfg"
	. "gopkg.in/check.v1"
)

func (s *BitcoinCashSignerSuite) TestGetChainCfg(c *C) {
	param := s.client.getChainCfgBCH()
	c.Assert(param, Equals, &chaincfg.MainNetParams)
}
