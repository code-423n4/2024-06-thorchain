package keeperv1

import (
	"fmt"
	"strconv"

	"github.com/blang/semver"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

func (k KVStore) setTxOut(ctx cosmos.Context, key string, record TxOut) {
	store := ctx.KVStore(k.storeKey)
	buf := k.cdc.MustMarshal(&record)
	if buf == nil {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), buf)
	}
}

func (k KVStore) getTxOut(ctx cosmos.Context, key string, record *TxOut) (bool, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return false, nil
	}

	bz := store.Get([]byte(key))
	if err := k.cdc.Unmarshal(bz, record); err != nil {
		return true, dbError(ctx, fmt.Sprintf("Unmarshal kvstore: (%T) %s", record, key), err)
	}
	return true, nil
}

// AppendTxOut - append the given item to txOut
func (k KVStore) AppendTxOut(ctx cosmos.Context, height int64, item TxOutItem) error {
	block, err := k.GetTxOut(ctx, height)
	if err != nil {
		return err
	}
	block.TxArray = append(block.TxArray, item)
	return k.SetTxOut(ctx, block)
}

// ClearTxOut - remove the txout of the given height from key value  store
func (k KVStore) ClearTxOut(ctx cosmos.Context, height int64) error {
	k.del(ctx, k.GetKey(ctx, prefixTxOut, strconv.FormatInt(height, 10)))
	return nil
}

// SetTxOut - write the given txout information to key value store
func (k KVStore) SetTxOut(ctx cosmos.Context, blockOut *TxOut) error {
	if blockOut == nil || blockOut.IsEmpty() {
		return nil
	}
	k.setTxOut(ctx, k.GetKey(ctx, prefixTxOut, strconv.FormatInt(blockOut.Height, 10)), *blockOut)
	return nil
}

// GetTxOutIterator iterate tx out
func (k KVStore) GetTxOutIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixTxOut)
}

// GetTxOut - write the given txout information to key values tore
func (k KVStore) GetTxOut(ctx cosmos.Context, height int64) (*TxOut, error) {
	record := NewTxOut(height)
	_, err := k.getTxOut(ctx, k.GetKey(ctx, prefixTxOut, strconv.FormatInt(height, 10)), record)
	return record, err
}

func (k KVStore) GetTxOutValue(ctx cosmos.Context, height int64) (cosmos.Uint, cosmos.Uint, error) {
	version := k.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.128.0")):
		return k.GetTxOutValueV128(ctx, height)
	default:
		return k.GetTxOutValueV1(ctx, height)
	}
}

func (k KVStore) GetTxOutValueV128(ctx cosmos.Context, height int64) (cosmos.Uint, cosmos.Uint, error) {
	txout, err := k.GetTxOut(ctx, height)
	if err != nil {
		return cosmos.ZeroUint(), cosmos.ZeroUint(), err
	}
	runeValue, cloutValue := k.GetTOIsValue(ctx, txout.TxArray...)
	return runeValue, cloutValue, nil
}

func (k KVStore) GetTOIsValue(ctx cosmos.Context, tois ...TxOutItem) (cosmos.Uint, cosmos.Uint) {
	runeValue := cosmos.ZeroUint()
	cloutValue := cosmos.ZeroUint()
	poolCache := map[common.Asset]Pool{} // Cache the pools to avoid duplicated GetPool calls
	for i := range tois {
		for _, coin := range append(common.Coins{tois[i].Coin}, tois[i].MaxGas...) {
			if coin.Asset.IsRune() {
				runeValue = runeValue.Add(coin.Amount)
				continue
			}

			pool, ok := poolCache[coin.Asset]
			if !ok {
				var err error
				pool, err = k.GetPool(ctx, coin.Asset.GetLayer1Asset())
				if err != nil {
					_ = dbError(ctx, fmt.Sprintf("unable to get pool : %s", coin.Asset), err)
					continue
				}
				poolCache[coin.Asset] = pool
			}
			runeValue = runeValue.Add(pool.AssetValueInRune(coin.Amount))
		}

		if tois[i].CloutSpent != nil {
			cloutValue = cloutValue.Add(*tois[i].CloutSpent)
		}
	}

	return runeValue, cloutValue
}
