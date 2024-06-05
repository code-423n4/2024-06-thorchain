package keeperv1

import (
	"fmt"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper/types"
)

func (k KVStore) setStreamingSwap(ctx cosmos.Context, key string, record StreamingSwap) {
	store := ctx.KVStore(k.storeKey)
	buf := k.cdc.MustMarshal(&record)
	if buf == nil {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), buf)
	}
}

func (k KVStore) getStreamingSwap(ctx cosmos.Context, key string, record *StreamingSwap) (bool, error) {
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

// GetStreamingSwapIterator iterate streaming swaps
func (k KVStore) GetStreamingSwapIterator(ctx cosmos.Context) cosmos.Iterator {
	key := k.GetKey(ctx, prefixStreamingSwap, "")
	return k.getIterator(ctx, types.DbPrefix(key))
}

// GetStreamingSwap retrieve streaming swap from the data store
func (k KVStore) GetStreamingSwap(ctx cosmos.Context, hash common.TxID) (StreamingSwap, error) {
	record := NewStreamingSwap(hash, 0, 0, cosmos.ZeroUint(), cosmos.ZeroUint())
	_, err := k.getStreamingSwap(ctx, k.GetKey(ctx, prefixStreamingSwap, hash.String()), &record)
	return record, err
}

// StreamingSwapExists check whether the given hash is associated with a swap
func (k KVStore) StreamingSwapExists(ctx cosmos.Context, hash common.TxID) bool {
	return k.has(ctx, k.GetKey(ctx, prefixStreamingSwap, hash.String()))
}

// SetStreamingSwap save the streaming swap to kv store
func (k KVStore) SetStreamingSwap(ctx cosmos.Context, swp StreamingSwap) {
	k.setStreamingSwap(ctx, k.GetKey(ctx, prefixStreamingSwap, swp.TxID.String()), swp)
}

// RemoveStreamingSwap remove the loan to kv store
func (k KVStore) RemoveStreamingSwap(ctx cosmos.Context, hash common.TxID) {
	k.del(ctx, k.GetKey(ctx, prefixStreamingSwap, hash.String()))
}
