//go:build !stagenet && !mocknet
// +build !stagenet,!mocknet

package aggregators

import (
	"gitlab.com/thorchain/thornode/common"
)

func DexAggregatorsV93() []Aggregator {
	return []Aggregator{
		// TSAggregatorGeneric
		{common.ETHChain, `0xd31f7e39afECEc4855fecc51b693F9A0Cec49fd2`, 400_000},
		// TSAggregatorUniswapV2
		{common.ETHChain, `0x7C38b8B2efF28511ECc14a621e263857Fb5771d3`, 400_000},
		// TSAggregatorUniswapV3 500
		{common.ETHChain, `0x1C0Ee4030f771a1BB8f72C86150730d063f6b3ff`, 400_000},
		// TSAggregatorUniswapV3 3000
		{common.ETHChain, `0x96ab925EFb957069507894CD941F40734f0288ad`, 400_000},
		// TSAggregatorUniswapV3 10000
		{common.ETHChain, `0xE308B9562de7689B2d31C76a41649933F38ab761`, 400_000},
		// TSAggregator2LegUniswapV2 USDC
		{common.ETHChain, `0x3660dE6C56cFD31998397652941ECe42118375DA`, 400_000},
		// TSAggregator SUSHIswap
		{common.ETHChain, `0x0F2CD5dF82959e00BE7AfeeF8245900FC4414199`, 400_000},
		// RangoThorchainOutputAggUniV2
		{common.ETHChain, `0x2a7813412b8da8d18Ce56FE763B9eb264D8e28a8`, 400_000},
		// RangoThorchainOutputAggUniV3
		{common.ETHChain, `0xbB8De86F3b041B3C084431dcf3159fE4827c5F0D`, 400_000},
	}
}
