package keeperv1

import (
	"fmt"
	"strings"

	"github.com/blang/semver"

	"gitlab.com/thorchain/thornode/common/cosmos"
)

// TODO: Remove Kraken checks on next hard fork.
const KRAKEN string = "ReleaseTheKraken"

// GetMimir get a mimir value from key value store
func (k KVStore) GetMimir(ctx cosmos.Context, key string) (int64, error) {
	if k.GetVersion().LT(semver.MustParse("1.92.0")) {
		// Once node mimir reach super majority  , admin mimir will be set automatically to lock in the vote
		// thus , there is no need to check node mimir every time
		nodeMimirs, err := k.GetNodeMimirs(ctx, key)
		if err != nil {
			return -1, err
		}
		activeNodes, err := k.ListActiveValidators(ctx)
		if err != nil {
			return -1, err
		}
		if i, ok := nodeMimirs.HasSuperMajority(key, activeNodes.GetNodeAddresses()); ok {
			return i, nil
		}
	}

	// TODO: Remove Kraken checks on next hard fork.
	// The Kraken functionality has been removed, but we need to maintain
	// the same KVStore cost accounting. So always check, but ignore the
	// value.
	ignored := k.haveKraken(ctx)
	_ = ignored

	record := int64(-1)
	_, err := k.getInt64(ctx, k.GetKey(ctx, prefixMimir, key), &record)
	return record, err
}

// GetMimirWithRef is a helper function to more readably insert references (such as Asset MimirString or Chain) into Mimir key templates.
func (k KVStore) GetMimirWithRef(ctx cosmos.Context, template, ref string) (int64, error) {
	// 'template' should be something like "Halt%sChain" (to halt an arbitrary specified chain)
	// or "Ragnarok-%s" (to halt the pool of an arbitrary specified Asset (MimirString used for Assets to join Chain and Symbol with a hyphen).
	key := fmt.Sprintf(template, ref)
	return k.GetMimir(ctx, key)
}

// TODO: Remove Kraken checks on next hard fork.
// haveKraken - check to see if we have "released the kraken"
func (k KVStore) haveKraken(ctx cosmos.Context) bool {
	record := int64(-1)
	_, _ = k.getInt64(ctx, k.GetKey(ctx, prefixMimir, KRAKEN), &record)
	return record >= 0
}

// SetMimir save a mimir value to key value store
func (k KVStore) SetMimir(ctx cosmos.Context, key string, value int64) {
	// TODO: Remove Kraken checks on next hard fork.
	if ignored := k.haveKraken(ctx); ignored {
		// The Kraken functionality has been removed, but we need to maintain
		// the same KVStore cost accounting. So always check, but ignore the
		// value.
		_ = ignored
	}
	k.setInt64(ctx, k.GetKey(ctx, prefixMimir, key), value)
}

// GetNodeMimirs get node mimirs value from key value store
func (k KVStore) GetNodeMimirs(ctx cosmos.Context, key string) (NodeMimirs, error) {
	key = strings.ToUpper(key)
	record := NodeMimirs{}
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(k.GetKey(ctx, prefixNodeMimir, key))) {
		return record, nil
	}
	bz := store.Get([]byte(k.GetKey(ctx, prefixNodeMimir, key)))
	if err := k.cdc.Unmarshal(bz, &record); err != nil {
		return NodeMimirs{}, dbError(ctx, fmt.Sprintf("Unmarshal kvstore: (%T) %s", record, key), err)
	}
	return record, nil
}

// SetNodeMimir save a mimir value to key value store for a specific node
func (k KVStore) SetNodeMimir(ctx cosmos.Context, key string, value int64, acc cosmos.AccAddress) error {
	key = strings.ToUpper(key)
	kvkey := k.GetKey(ctx, prefixNodeMimir, key)
	record, err := k.GetNodeMimirs(ctx, key)
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

// DeleteNodeMimirs deletes all node mimir votes for a given key
func (k KVStore) DeleteNodeMimirs(ctx cosmos.Context, key string) {
	k.del(ctx, k.GetKey(ctx, prefixNodeMimir, key))
}

func (k KVStore) PurgeOperationalNodeMimirs(ctx cosmos.Context) {
	iterNode := k.GetNodeMimirIterator(ctx)
	defer iterNode.Close()
	for ; iterNode.Valid(); iterNode.Next() {
		key := strings.TrimPrefix(string(iterNode.Key()), string(prefixNodeMimir)+"/")
		if k.IsOperationalMimir(key) {
			k.DeleteNodeMimirs(ctx, key)
		}
	}
}

// GetMimirIterator iterate gas units
func (k KVStore) GetMimirIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixMimir)
}

// GetNodeMimirIterator iterate gas units
func (k KVStore) GetNodeMimirIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixNodeMimir)
}

func (k KVStore) DeleteMimir(ctx cosmos.Context, key string) error {
	k.del(ctx, k.GetKey(ctx, prefixMimir, key))
	return nil
}

func (k KVStore) GetNodePauseChain(ctx cosmos.Context, acc cosmos.AccAddress) int64 {
	record := int64(-1)
	_, _ = k.getInt64(ctx, k.GetKey(ctx, prefixNodePauseChain, acc.String()), &record)
	return record
}

func (k KVStore) SetNodePauseChain(ctx cosmos.Context, acc cosmos.AccAddress) {
	k.setInt64(ctx, k.GetKey(ctx, prefixNodePauseChain, acc.String()), ctx.BlockHeight())
}

func (k KVStore) IsOperationalMimir(key string) bool {
	version := k.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.133.0")):
		return isOperationalMimirV133(key)
	default:
		return false
	}
}

func isOperationalMimirV133(key string) bool {
	exactMatches := []string{
		"BurnSynths",
		"MintSynths",
	}
	for i := range exactMatches {
		if strings.EqualFold(key, exactMatches[i]) {
			return true
		}
	}

	// Past this point, compare only upper-case strings due to case sensitivity.
	key = strings.ToUpper(key)
	partialMatches := []string{
		"HALT",
		"PAUSE",
		"STOPSOLVENCYCHECK",
	}
	for i := range partialMatches {
		// Contains rather than HasPrefix to include cases like StreamingSwapPause.
		if strings.Contains(key, partialMatches[i]) {
			return true
		}
	}

	return false
}
