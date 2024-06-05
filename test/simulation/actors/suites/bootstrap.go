package suites

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/test/simulation/actors"
	"gitlab.com/thorchain/thornode/test/simulation/pkg/evm"
	"gitlab.com/thorchain/thornode/test/simulation/pkg/thornode"
	. "gitlab.com/thorchain/thornode/test/simulation/pkg/types"
)

////////////////////////////////////////////////////////////////////////////////////////
// Bootstrap
////////////////////////////////////////////////////////////////////////////////////////

func Bootstrap() *Actor {
	a := NewActor("Bootstrap")

	pools, err := thornode.GetPools()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get pools")
	}

	// bootstrap pools for all chains
	count := 0
	for _, chain := range common.AllChains {
		// skip thorchain and deprecated chains
		switch chain {
		case common.THORChain, common.BNBChain, common.TERRAChain:
			continue
		}
		count++

		// skip bootstrapping existing pools
		found := false
		for _, pool := range pools {
			if pool.Asset == chain.GetGasAsset().String() {
				found = true
				break
			}
		}
		if found {
			log.Info().Str("chain", chain.GetGasAsset().String()).Msg("skip existing pool bootstrap")
			continue
		}

		a.Children[actors.NewDualLPActor(chain.GetGasAsset())] = true
	}

	// create token pools
	tokenPools := NewActor("Bootstrap-TokenPools")
	for _, chain := range common.AllChains {
		if !chain.IsEVM() {
			continue
		}
		count++

		for asset := range evm.Tokens(chain) {
			// skip bootstrapping existing pools
			found := false
			for _, pool := range pools {
				if pool.Asset == asset.String() {
					found = true
					break
				}
			}
			if found {
				log.Info().Str("chain", chain.GetGasAsset().String()).Msg("skip existing pool bootstrap")
				continue
			}

			tokenPools.Children[actors.NewDualLPActor(asset)] = true
		}
	}
	a.Append(tokenPools)

	// verify pools
	verify := NewActor("Bootstrap-Verify")
	verify.Ops = append(verify.Ops, func(config *OpConfig) OpResult {
		pools, err := thornode.GetPools()
		if err != nil {
			return OpResult{Finish: true, Error: err}
		}

		// all pools should be available
		for _, pool := range pools {
			if pool.Status != "Available" {
				return OpResult{
					Finish: true,
					Error:  fmt.Errorf("pool %s not available", pool.Asset),
				}
			}
		}

		// all chains should have pools
		if len(pools) != count {
			return OpResult{
				Finish: true,
				Error:  fmt.Errorf("expected %d pools, got %d", count, len(pools)),
			}
		}

		return OpResult{Finish: true}
	},
	)
	a.Append(verify)

	return a
}
