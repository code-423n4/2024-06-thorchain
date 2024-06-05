package keeperv1

import (
	"fmt"
	"strings"

	"gitlab.com/thorchain/thornode/common/cosmos"
)

// SetMimir save a mimir value to key value store
func (k KVStore) SetMimirV2(ctx cosmos.Context, key string, value int64) {
	k.setInt64(ctx, k.GetKey(ctx, prefixMimirV2, key), value)
}

// GetMimir get a mimir value from key value store
func (k KVStore) GetMimirV2(ctx cosmos.Context, key string) (int64, error) {
	record := int64(-1)
	_, err := k.getInt64(ctx, k.GetKey(ctx, prefixMimirV2, key), &record)
	return record, err
}

// GetNodeMimirs get node mimirs value from key value store
func (k KVStore) GetNodeMimirsV2(ctx cosmos.Context, key string) (NodeMimirs, error) {
	record := NodeMimirs{}
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(k.GetKey(ctx, prefixNodeMimirV2, key)))
	if bz == nil {
		// not found
		return record, nil
	}
	if err := k.cdc.Unmarshal(bz, &record); err != nil {
		return NodeMimirs{}, dbError(ctx, fmt.Sprintf("Unmarshal kvstore: (%T) %s", record, key), err)
	}
	return record, nil
}

// SetNodeMimir save a mimir value to key value store for a specific node
func (k KVStore) SetNodeMimirV2(ctx cosmos.Context, key string, value int64, acc cosmos.AccAddress) error {
	key = strings.ToUpper(key) // ensure uppercase
	kvkey := k.GetKey(ctx, prefixNodeMimirV2, key)
	record, err := k.GetNodeMimirsV2(ctx, key)
	if err != nil {
		return err
	}
	record.Set(key, value, acc)
	store := ctx.KVStore(k.storeKey)
	buf := k.cdc.MustMarshal(&record)
	if buf == nil || len(record.Mimirs) == 0 {
		store.Delete([]byte(kvkey))
	} else {
		store.Set([]byte(kvkey), buf)
	}
	return err
}

// GetNodeMimirIterator iterate node mimirs
func (k KVStore) GetNodeMimirIteratorV2(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixNodeMimirV2)
}
