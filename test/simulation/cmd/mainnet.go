//go:build !mocknet && !stagenet
// +build !mocknet,!stagenet

package main

import (
	"gitlab.com/thorchain/thornode/common"
	. "gitlab.com/thorchain/thornode/test/simulation/pkg/types"
)

////////////////////////////////////////////////////////////////////////////////////////
// Init
////////////////////////////////////////////////////////////////////////////////////////

// trunk-ignore(golangci-lint/unused)
var chainRPCs = map[common.Chain]string{}

func InitConfig(parallelism int) *OpConfig {
	return nil
}
