//go:build !stagenet && !mocknet
// +build !stagenet,!mocknet

package aggregators

import (
	"gitlab.com/thorchain/thornode/common"
)

func DexAggregatorsV116() []Aggregator {
	return []Aggregator{
		// TSAggregatorVTHOR Ethereum V2
		{common.ETHChain, `0x0581a9aB98c467dCA614C940104E6dD102BE5C7d`, 400_000},
		// TSAggregatorPancakeSwap Ethereum V2
		{common.ETHChain, `0x35CF22003c90126528fbe95b21bB3ADB2ca8c53D`, 400_000},
		// TSAggregatorWoofi Avalanche V2
		{common.AVAXChain, `0x5505BE604dFA8A1ad402A71f8A357fba47F9bf5a`, 400_000},
		// TSSwapGeneric Avalanche V2
		{common.AVAXChain, `0x77b34A3340eDdD56799749be4Be2c322547E2428`, 400_000},
		// TSAggregatorGeneric Avalanche V2
		{common.AVAXChain, `0x94a852F0a21E473078846cf88382dd8d15bD1Dfb`, 400_000},
		// TSAggregatorTraderJoe Avalanche V2
		{common.AVAXChain, `0xce5d236164D2Bc0B2f65351f23B617c2A7D5Cc28`, 400_000},
		// TSAggregatorPangolin Avalanche V2
		{common.AVAXChain, `0x9aC752Ed433f7E038Be4070544858cB3d83cC0d7`, 400_000},
		// TSSwapGeneric Ethereum V2
		{common.ETHChain, `0x213255345a740324cbCE0242e32076Ab735906e2`, 400_000},
		// TSAggregatorGeneric Ethereum V2
		{common.ETHChain, `0x0ccD5Dd5BcF1Af77dc358d1E2F06eE880EF63C3c`, 400_000},
		// TSAggregatorUniswapV2 Ethereum V2
		{common.ETHChain, `0x14D52a5709743C9563a2C36842B3Fe7Db1fCf5bc`, 400_000},
		// TSAggregatorSushiswap Ethereum V2
		{common.ETHChain, `0x7334543783a6A87BDD028C902f7c87AFB703cCbC`, 400_000},
		// TSAggregatorUniswapV3_500 Ethereum V2
		{common.ETHChain, `0xBcd954803163094590AF749377c082619014acD5`, 400_000},
		// TSAggregatorUniswapV3_3000 Ethereum V2
		{common.ETHChain, `0xd785Eb8D8cf2adC99b742C4E7C77d39f1bC604F1`, 400_000},
		// TSAggregatorUniswapV3_10000 Ethereum V2
		{common.ETHChain, `0xDE3205dc90336C916CbBAD21383eA95F418a7cbA`, 400_000},
		// TSAggregatorGeneric
		{common.ETHChain, `0xd31f7e39afECEc4855fecc51b693F9A0Cec49fd2`, 400_000},
		// TSAggregator2LegUniswapV2 USDC
		{common.ETHChain, `0x3660dE6C56cFD31998397652941ECe42118375DA`, 400_000},
		// RangoThorchainOutputAggUniV2
		{common.ETHChain, `0x2a7813412b8da8d18Ce56FE763B9eb264D8e28a8`, 400_000},
		// RangoThorchainOutputAggUniV3
		{common.ETHChain, `0xbB8De86F3b041B3C084431dcf3159fE4827c5F0D`, 400_000},
		// PangolinAggregator
		{common.AVAXChain, `0x7a68c37D8AFA3078f3Ad51D98eA23Fe57a8Ae21a`, 400_000},
		// TSAggregatorUniswapV2 - short notation
		{common.ETHChain, `0x86904eb2b3c743400d03f929f2246efa80b91215`, 400_000},
		// TSAggregatorSushiswap - short notation
		{common.ETHChain, `0xbf365e79aa44a2164da135100c57fdb6635ae870`, 400_000},
		// TSAggregatorUniswapV3 100 - short notation
		{common.ETHChain, `0xbd68cbe6c247e2c3a0e36b8f0e24964914f26ee8`, 400_000},
		// TSAggregatorUniswapV3 500 - short notation
		{common.ETHChain, `0xe4ddca21881bac219af7f217703db0475d2a9f02`, 400_000},
		// TSAggregatorUniswapV3 3000 - short notation
		{common.ETHChain, `0x11733abf0cdb43298f7e949c930188451a9a9ef2`, 400_000},
		// TSAggregatorUniswapV3 10000 - short notation
		{common.ETHChain, `0xb33874810e5395eb49d8bd7e912631db115d5a03`, 400_000},
		// TSAggregatorPangolin
		{common.AVAXChain, `0x942c6dA485FD6cEf255853ef83a149d43A73F18a`, 400_000},
		// TSAggregatorTraderJoe
		{common.AVAXChain, `0x3b7DbdD635B99cEa39D3d95Dbd0217F05e55B212`, 400_000},
		// TSAggregatorAvaxGeneric
		{common.AVAXChain, `0x7C38b8B2efF28511ECc14a621e263857Fb5771d3`, 400_000},
		// XDEFIAggregatorEthGeneric
		{common.ETHChain, `0x53E4DD4072A9a8ed56289e048f5BD5AA51c9Bf6E`, 400_000},
		// XDEFIAggregatorEthUniswapV2
		{common.ETHChain, `0xeEe520b0DA1F8a9e4a0480F92CC4c5f6C027ef1E`, 400_000},
		// XDEFIAggregatorAvaxGeneric
		{common.AVAXChain, `0xd0269244A876F7Bc600D1f38B03a9916864b73C6`, 400_000},
		// XDEFIAggregatorAvaxTraderJoe
		{common.AVAXChain, `0x4ab34123A077aE294A39844f3e8df418d2A3D8c4`, 400_000},
		// XDEFIAggregatorUniswapV3 100 - short notation
		{common.ETHChain, `0x88100E08e5287bA3445F95d448ABfF3113d82a4C`, 400_000},
		// XDEFIAggregatorUniswapV3 500 - short notation
		{common.ETHChain, `0xC1faA12981160945903E0725888828E2d6a15821`, 400_000},
		// XDEFIAggregatorUniswapV3 3000 - short notation
		{common.ETHChain, `0x7E019988299cd8038091D8d7fe38f7a1dd3f90F1`, 400_000},
		// XDEFIAggregatorUniswapV3 10000 - short notation
		{common.ETHChain, `0x95B6b888a9fCc5BCA4A3004Df5E9498B63195F48`, 400_000},
		// TSAggregatorGeneric
		{common.BSCChain, `0xB6fA6f1DcD686F4A573Fd243a6FABb4ba36Ba98c`, 400_000},
		// TSAggregatorPancakeV2 BinanceSmartChain
		{common.BSCChain, `0x30912B38618D3D37De3191A4FFE982C65a9aEC2E`, 400_000},
		// TSAggregatorStargate Ethereum gen2
		{common.ETHChain, `0x1204b5Bf0D6d48E718B1d9753A4166A7258B8432`, 800_000},
		// LayerZero Executor Ethereum
		{common.ETHChain, `0xe93685f3bBA03016F02bD1828BaDD6195988D950`, 800_000},
		// LayerZero Executor Avalanche
		{common.AVAXChain, `0xe93685f3bBA03016F02bD1828BaDD6195988D950`, 800_000},
		// LayerZero Executor BinanceSmartChain
		{common.BSCChain, `0xe93685f3bBA03016F02bD1828BaDD6195988D950`, 800_000},
	}
}
