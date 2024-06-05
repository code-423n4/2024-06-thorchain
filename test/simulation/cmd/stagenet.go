//go:build stagenet
// +build stagenet

package main

import (
	"gitlab.com/thorchain/thornode/common"
	. "gitlab.com/thorchain/thornode/test/simulation/pkg/types"
)

////////////////////////////////////////////////////////////////////////////////////////
// Init
////////////////////////////////////////////////////////////////////////////////////////

var chainRPCs = map[common.Chain]string{}

func InitConfig(parallelism int) *OpConfig {
	return nil
}
