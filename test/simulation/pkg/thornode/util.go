package thornode

import (
	"fmt"

	"gitlab.com/thorchain/thornode/common/cosmos"
	openapi "gitlab.com/thorchain/thornode/openapi/gen"
)

////////////////////////////////////////////////////////////////////////////////////////
// Exported
////////////////////////////////////////////////////////////////////////////////////////

// ConvertAssetAmount converts the given coin to the target asset and returns the amount.
func ConvertAssetAmount(coin openapi.Coin, asset string) (cosmos.Uint, error) {
	pools, err := GetPools()
	if err != nil {
		return cosmos.ZeroUint(), err
	}

	// find pools for the conversion rate
	var sourcePool, targetPool openapi.Pool
	for _, pool := range pools {
		if pool.Asset == coin.Asset {
			sourcePool = pool
		}
		if pool.Asset == asset {
			targetPool = pool
		}
	}

	// ensure we found both pools
	if sourcePool.Asset == "" {
		return cosmos.ZeroUint(), fmt.Errorf("source asset not found")
	}
	if targetPool.Asset == "" {
		return cosmos.ZeroUint(), fmt.Errorf("target asset not found")
	}

	// convert the amount
	converted := cosmos.NewUintFromString(coin.Amount).
		Mul(cosmos.NewUintFromString(sourcePool.BalanceRune)).
		Quo(cosmos.NewUintFromString(sourcePool.BalanceAsset)).
		Mul(cosmos.NewUintFromString(targetPool.BalanceAsset)).
		Quo(cosmos.NewUintFromString(targetPool.BalanceRune))

	return converted, nil
}
