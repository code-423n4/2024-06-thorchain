package static

import (
	"math/rand"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/test/simulation/actors"
	"gitlab.com/thorchain/thornode/test/simulation/pkg/evm"
	. "gitlab.com/thorchain/thornode/test/simulation/pkg/types"
)

////////////////////////////////////////////////////////////////////////////////////////
// Swaps
////////////////////////////////////////////////////////////////////////////////////////

func Swaps() *Actor {
	a := NewActor("Swaps")

	// gather all pools we expect to swap through
	swapPools := []common.Asset{}
	for _, chain := range common.AllChains {
		// skip thorchain and deprecated chains
		switch chain {
		case common.THORChain, common.BNBChain, common.TERRAChain:
			continue
		}

		swapPools = append(swapPools, chain.GetGasAsset())

		// add tokens to swap pools
		if !chain.IsEVM() {
			continue
		}
		for asset := range evm.Tokens(chain) {
			swapPools = append(swapPools, asset)
		}
	}

	// copy pools to new slice for shuffling target pools
	shufflePools := make([]common.Asset, len(swapPools))
	copy(shufflePools, swapPools)

	// check every gas asset swap route
	for _, source := range swapPools {
		// shuffle the pools
		rand.Shuffle(len(shufflePools), func(i, j int) {
			shufflePools[i], shufflePools[j] = shufflePools[j], shufflePools[i]
		})

		// swap to a random half of the pools
		for _, target := range shufflePools[:len(shufflePools)/2] {
			// skip swap to self
			if source.Equals(target) {
				continue
			}
			a.Children[actors.NewSwapActor(source, target)] = true
		}
	}

	return a
}
