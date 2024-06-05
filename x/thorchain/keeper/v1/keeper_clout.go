package keeperv1

import (
	"fmt"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

func (k KVStore) setSwapperClout(ctx cosmos.Context, key string, record SwapperClout) {
	store := ctx.KVStore(k.storeKey)
	buf := k.cdc.MustMarshal(&record)
	if buf == nil {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), buf)
	}
}

func (k KVStore) getSwapperClout(ctx cosmos.Context, key string, record *SwapperClout) (bool, error) {
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

func (k KVStore) SetSwapperClout(ctx cosmos.Context, record SwapperClout) error {
	k.setSwapperClout(ctx, k.GetKey(ctx, prefixSwapperClout, record.Address.String()), record)
	return nil
}

func (k KVStore) GetSwapperClout(ctx cosmos.Context, addr common.Address) (SwapperClout, error) {
	record := NewSwapperClout(addr)
	if addr.IsEmpty() {
		return record, nil
	}
	_, err := k.getSwapperClout(ctx, k.GetKey(ctx, prefixSwapperClout, addr.String()), &record)
	return record, err
}

func (k KVStore) GetSwapperCloutIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixSwapperClout)
}
