package keeperv1

import (
	"fmt"

	"gitlab.com/thorchain/thornode/common/cosmos"
)

func (k KVStore) GetTxOutValueV1(ctx cosmos.Context, height int64) (cosmos.Uint, cosmos.Uint, error) {
	txout, err := k.GetTxOut(ctx, height)
	if err != nil {
		return cosmos.ZeroUint(), cosmos.ZeroUint(), err
	}

	runeValue := cosmos.ZeroUint()
	cloutValue := cosmos.ZeroUint()
	for _, item := range txout.TxArray {
		if item.Coin.Asset.IsRune() {
			runeValue = runeValue.Add(item.Coin.Amount)
		} else {
			var pool Pool
			pool, err = k.GetPool(ctx, item.Coin.Asset)
			if err != nil {
				_ = dbError(ctx, fmt.Sprintf("unable to get pool : %s", item.Coin.Asset), err)
				continue
			}
			runeValue = runeValue.Add(pool.AssetValueInRune(item.Coin.Amount))
		}
		if item.CloutSpent != nil {
			cloutValue = cloutValue.Add(*item.CloutSpent)
		}
	}

	return runeValue, cloutValue, nil
}
