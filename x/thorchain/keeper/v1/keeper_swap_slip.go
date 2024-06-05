package keeperv1

import (
	"fmt"

	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper/types"
)

// AddToSwapSlip - add swap slip to block
func (k KVStore) AddToSwapSlip(ctx cosmos.Context, asset common.Asset, amt cosmos.Int) error {
	currentHeight := ctx.BlockHeight()

	poolSlip, err := k.GetPoolSwapSlip(ctx, currentHeight, asset)
	if err != nil {
		return err
	}

	poolSlip = poolSlip.Add(amt)

	// update pool slip
	k.setInt64(ctx, k.GetKey(ctx, prefixPoolSwapSlip, fmt.Sprintf("%d-%s", currentHeight, asset.String())), poolSlip.Int64())
	return nil
}

func (k KVStore) DeletePoolSwapSlip(ctx cosmos.Context, height int64, asset common.Asset) {
	key := k.GetKey(ctx, prefixPoolSwapSlip, fmt.Sprintf("%d-%s", height, asset.String()))
	k.del(ctx, key)
}

func (k KVStore) getSwapSlip(ctx cosmos.Context, key string) (cosmos.Int, error) {
	var record int64
	_, err := k.getInt64(ctx, key, &record)
	return cosmos.NewInt(record), err
}

// GetPoolSwapSlip - total of slip in each block per pool
func (k KVStore) GetPoolSwapSlip(ctx cosmos.Context, height int64, asset common.Asset) (cosmos.Int, error) {
	key := k.GetKey(ctx, prefixPoolSwapSlip, fmt.Sprintf("%d-%s", height, asset.String()))
	return k.getSwapSlip(ctx, key)
}

func (k KVStore) GetCurrentRollup(ctx cosmos.Context, asset common.Asset) (int64, error) {
	var currRollup int64
	currRollupKey := k.GetKey(ctx, prefixPoolSwapSlip, fmt.Sprintf("rollup/%s", asset.String()))
	_, err := k.getInt64(ctx, currRollupKey, &currRollup)
	if err != nil {
		return 0, err
	}
	return currRollup, nil
}

func (k KVStore) SetCurrentRollup(ctx cosmos.Context, asset common.Asset, currRollup int64) {
	currRollupKey := k.GetKey(ctx, prefixPoolSwapSlip, fmt.Sprintf("rollup/%s", asset.String()))
	k.setInt64(ctx, currRollupKey, currRollup)
}

// GetSwapSlipSnapShotIterator
func (k KVStore) GetSwapSlipSnapShotIterator(ctx cosmos.Context, asset common.Asset) cosmos.Iterator {
	key := k.GetKey(ctx, prefixPoolSwapSnapShot, asset.String())
	return k.getIterator(ctx, types.DbPrefix(key))
}

func (k KVStore) GetSwapSlipSnapShot(ctx cosmos.Context, asset common.Asset, height int64) (int64, error) {
	snapshotKey := k.GetKey(ctx, prefixPoolSwapSnapShot, fmt.Sprintf("%s/%d", asset.String(), height))
	var record int64
	_, err := k.getInt64(ctx, snapshotKey, &record)
	if err != nil {
		return 0, err
	}
	return record, nil
}

func (k KVStore) SetSwapSlipSnapShot(ctx cosmos.Context, asset common.Asset, height, currRollup int64) {
	snapshotKey := k.GetKey(ctx, prefixPoolSwapSnapShot, fmt.Sprintf("%s/%d", asset.String(), height))
	k.setInt64(ctx, snapshotKey, currRollup)
}

func (k KVStore) GetRollupCount(ctx cosmos.Context, asset common.Asset) (int64, error) {
	var currCount int64
	currCountKey := k.GetKey(ctx, prefixPoolSwapSlip, fmt.Sprintf("rollup-count/%s", asset.String()))
	_, err := k.getInt64(ctx, currCountKey, &currCount)
	if err != nil {
		return 0, err
	}
	return currCount, nil
}

// RollupSwapSlip - sums the amount of slip in a given pool in the last targetCount blocks
func (k KVStore) RollupSwapSlip(ctx cosmos.Context, targetCount int64, asset common.Asset) (cosmos.Int, error) {
	currCount, err := k.GetRollupCount(ctx, asset)
	if err != nil {
		return cosmos.ZeroInt(), err
	}

	currRollup, err := k.GetCurrentRollup(ctx, asset)
	if err != nil {
		return cosmos.ZeroInt(), err
	}

	currCountKey := k.GetKey(ctx, prefixPoolSwapSlip, fmt.Sprintf("rollup-count/%s", asset.String()))
	currRollupKey := k.GetKey(ctx, prefixPoolSwapSlip, fmt.Sprintf("rollup/%s", asset.String()))
	reset := func(err error) (cosmos.Int, error) {
		if err != nil {
			ctx.Logger().Error("resetting pool swap slip rollup", "asset", asset.String(), "err", err)
		}

		k.setInt64(ctx, currCountKey, 0)
		k.setInt64(ctx, currRollupKey, 0)
		return cosmos.ZeroInt(), err
	}

	if currCount > targetCount {
		// we need to reset, likely the target count was changed to a lower
		// number than it was before
		ctx.Logger().Info("resetting pool swap rollup", "asset", asset.String())
		return reset(nil)
	}

	// add the swap slip from the previous block to the rollup
	prevBlockSlip, err := k.GetPoolSwapSlip(ctx, ctx.BlockHeight()-1, asset)
	if err != nil {
		return reset(err)
	}
	currRollup += prevBlockSlip.Int64()
	currCount++

	if currCount > targetCount {
		// remove the oldest swap slip block from the count
		var oldBlockSlip cosmos.Int
		oldBlockSlip, err = k.GetPoolSwapSlip(ctx, ctx.BlockHeight()-targetCount, asset)
		if err != nil {
			return reset(err)
		}
		currRollup -= oldBlockSlip.Int64()
		currCount--
		k.DeletePoolSwapSlip(ctx, ctx.BlockHeight()-targetCount, asset)

		if k.GetVersion().GTE(semver.MustParse("1.121.0")) {
			if targetCount > 0 && ctx.BlockHeight()%targetCount == 0 {
				k.SetSwapSlipSnapShot(ctx, asset, ctx.BlockHeight(), currRollup)
			}
		}
	}

	if k.GetVersion().GTE(semver.MustParse("1.129.0")) {
		// slip rollup should never be negative
		if currRollup < 0 {
			currRollup = 0
		}
	}

	k.setInt64(ctx, currCountKey, currCount)
	k.setInt64(ctx, currRollupKey, currRollup)

	return cosmos.NewInt(currRollup), nil
}

func (k KVStore) GetLongRollup(ctx cosmos.Context, asset common.Asset) (int64, error) {
	var record int64
	key := k.GetKey(ctx, prefixPoolSwapSlipLong, asset.String())
	_, err := k.getInt64(ctx, key, &record)
	return record, err
}

func (k KVStore) SetLongRollup(ctx cosmos.Context, asset common.Asset, slip int64) {
	key := k.GetKey(ctx, prefixPoolSwapSlipLong, asset.String())
	k.setInt64(ctx, key, slip)
}
