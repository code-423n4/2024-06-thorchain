package thorchain

import (
	"strconv"
	"strings"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
)

// SwapQueueV117 is going to manage the swaps queue
type SwapQueueV117 struct {
	k       keeper.Keeper
	handler func(mgr Manager) cosmos.Handler
}

// newSwapQueueV117 create a new vault manager
func newSwapQueueV117(k keeper.Keeper) *SwapQueueV117 {
	return &SwapQueueV117{
		k:       k,
		handler: NewInternalHandler,
	}
}

// FetchQueue - grabs all swap queue items from the kvstore and returns them
func (vm *SwapQueueV117) FetchQueue(ctx cosmos.Context) (swapItems, error) { // nolint
	items := make(swapItems, 0)
	iterator := vm.k.GetSwapQueueIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var msg MsgSwap
		if err := vm.k.Cdc().Unmarshal(iterator.Value(), &msg); err != nil {
			ctx.Logger().Error("fail to fetch swap msg from queue", "error", err)
			continue
		}

		ss := strings.Split(string(iterator.Key()), "-")
		i, err := strconv.Atoi(ss[len(ss)-1])
		if err != nil {
			ctx.Logger().Error("fail to parse swap queue msg index", "key", iterator.Key(), "error", err)
			continue
		}

		// exclude streaming swaps when its not "their time". Always want to
		// allow the first sub-swap immediately (ie no LastHeight yet)
		if msg.IsStreaming() {
			pausedStreaming := vm.k.GetConfigInt64(ctx, constants.StreamingSwapPause)
			if pausedStreaming > 0 {
				continue
			}
			swp := msg.GetStreamingSwap()
			if vm.k.StreamingSwapExists(ctx, msg.Tx.ID) {
				var err error
				swp, err = vm.k.GetStreamingSwap(ctx, msg.Tx.ID)
				if err != nil {
					ctx.Logger().Error("fail to fetch streaming swap", "error", err)
					continue
				}
			}
			if swp.LastHeight > 0 { // if we don't have a height, do first swap attempt now
				if swp.LastHeight >= ctx.BlockHeight() {
					// last swap must be in the past
					continue // skip
				}
				if (ctx.BlockHeight()-swp.LastHeight)%int64(swp.Interval) != 0 {
					continue // skip
				}
			}
		}

		items = append(items, swapItem{
			msg:   msg,
			index: i,
			fee:   cosmos.ZeroUint(),
			slip:  cosmos.ZeroUint(),
		})
	}

	return items, nil
}

// EndBlock trigger the real swap to be processed
func (vm *SwapQueueV117) EndBlock(ctx cosmos.Context, mgr Manager) error {
	handler := vm.handler(mgr)

	minSwapsPerBlock, err := vm.k.GetMimir(ctx, constants.MinSwapsPerBlock.String())
	if minSwapsPerBlock < 0 || err != nil {
		minSwapsPerBlock = mgr.GetConstants().GetInt64Value(constants.MinSwapsPerBlock)
	}
	maxSwapsPerBlock, err := vm.k.GetMimir(ctx, constants.MaxSwapsPerBlock.String())
	if maxSwapsPerBlock < 0 || err != nil {
		maxSwapsPerBlock = mgr.GetConstants().GetInt64Value(constants.MaxSwapsPerBlock)
	}
	synthVirtualDepthMult, err := vm.k.GetMimir(ctx, constants.VirtualMultSynthsBasisPoints.String())
	if synthVirtualDepthMult < 1 || err != nil {
		synthVirtualDepthMult = mgr.GetConstants().GetInt64Value(constants.VirtualMultSynthsBasisPoints)
	}

	swaps, err := vm.FetchQueue(ctx)
	if err != nil {
		ctx.Logger().Error("fail to fetch swap queue from store", "error", err)
		return err
	}
	swaps, err = vm.scoreMsgs(ctx, swaps, synthVirtualDepthMult)
	if err != nil {
		ctx.Logger().Error("fail to fetch swap items", "error", err)
		// continue, don't exit, just do them out of order (instead of not at all)
	}
	swaps = swaps.Sort()

	for i := int64(0); i < vm.getTodoNum(int64(len(swaps)), minSwapsPerBlock, maxSwapsPerBlock); i++ {
		pick := swaps[i]
		// grab swp BEFORE a streaming swap modified the msg.Tx.Coins[0].Amount
		// value. This is used later to refund the correct amount
		swp := pick.msg.GetStreamingSwap()

		triggerRefund := false
		_, handleErr := handler(ctx, &pick.msg)
		if handleErr != nil {
			ctx.Logger().Error("fail to swap", "msg", pick.msg.Tx.String(), "error", handleErr)

			var refundErr error
			triggerRefund = !pick.msg.IsStreaming()

			if pick.msg.IsStreaming() {
				if vm.k.StreamingSwapExists(ctx, pick.msg.Tx.ID) {
					var getErr error
					swp, getErr = vm.k.GetStreamingSwap(ctx, pick.msg.Tx.ID)
					if getErr != nil {
						ctx.Logger().Error("fail to fetch streaming swap", "error", getErr)
						return getErr
					}
				}

				// if we haven't made any swaps yet, its safe to do a regular
				// refund. Otherwise allow later code to do partial refunds
				triggerRefund = swp.In.IsZero() && swp.Out.IsZero()
				if triggerRefund {
					// revert the tx amount to the be original deposit amount
					pick.msg.Tx.Coins[0].Amount = swp.Deposit
					vm.k.RemoveStreamingSwap(ctx, pick.msg.Tx.ID)
					vm.k.RemoveSwapQueueItem(ctx, pick.msg.Tx.ID, pick.index)
				}
			} else {
				// Should not refund a failed preferred asset swap, the RUNE is still in the AffiliateCollector module
				affColAddress, err := mgr.Keeper().GetModuleAddress(AffiliateCollectorName)
				if err != nil {
					ctx.Logger().Error("failed to retrieve AffiliateCollector module address", "error", err)
				}

				if strings.HasPrefix(pick.msg.Tx.Memo, "THOR-PREFERRED-ASSET") && pick.msg.Tx.FromAddress.Equals(affColAddress) {
					triggerRefund = false
				}
			}

			if triggerRefund {
				// Get the full ObservedTx from the TxID, for the vault ObservedPubKey to first try to refund from.
				voter, voterErr := mgr.Keeper().GetObservedTxInVoter(ctx, pick.msg.Tx.ID)
				if voterErr == nil && !voter.Tx.IsEmpty() {
					refundErr = refundTx(ctx, ObservedTx{Tx: pick.msg.Tx, ObservedPubKey: voter.Tx.ObservedPubKey}, mgr, CodeSwapFail, handleErr.Error(), "")
				} else {
					// If the full ObservedTx could not be retrieved, proceed with just the MsgSwap's Tx (no ObservedPubKey).
					ctx.Logger().Error("fail to get non-empty observed tx", "error", voterErr)
					refundErr = refundTx(ctx, ObservedTx{Tx: pick.msg.Tx}, mgr, CodeSwapFail, handleErr.Error(), "")
				}

				if refundErr != nil {
					ctx.Logger().Error("fail to refund swap", "error", refundErr)
				}
			}
		}

		if pick.msg.IsStreaming() {
			swp, err := vm.k.GetStreamingSwap(ctx, pick.msg.Tx.ID)
			if err != nil {
				ctx.Logger().Error("fail to fetch streaming swap", "error", err)
				return err
			}
			swp.Count += 1
			if handleErr != nil {
				swp.FailedSwaps = append(swp.FailedSwaps, swp.Count)
				swp.FailedSwapReasons = append(swp.FailedSwapReasons, handleErr.Error())
			}
			swp.LastHeight = ctx.BlockHeight()
			if !triggerRefund {
				mgr.Keeper().SetStreamingSwap(ctx, swp)
			}
			if swp.Valid() == nil && swp.IsDone() {
				vm.k.RemoveSwapQueueItem(ctx, pick.msg.Tx.ID, pick.index)
				vm.k.RemoveStreamingSwap(ctx, pick.msg.Tx.ID)

				tois := make([]TxOutItem, 0)
				if !swp.Out.IsZero() {
					dexAgg := ""
					if len(pick.msg.Aggregator) > 0 {
						dexAgg, err = FetchDexAggregator(
							mgr.GetVersion(),
							pick.msg.TargetAsset.GetChain(),
							pick.msg.Aggregator,
						)
						if err != nil {
							return err
						}
					}
					dexAggTargetAsset := pick.msg.AggregatorTargetAddress

					tois = append(tois, TxOutItem{
						Chain:                 pick.msg.TargetAsset.GetChain(),
						InHash:                pick.msg.Tx.ID,
						ToAddress:             pick.msg.Destination,
						Coin:                  common.NewCoin(pick.msg.TargetAsset, swp.Out),
						Aggregator:            dexAgg,
						AggregatorTargetAsset: dexAggTargetAsset,
						AggregatorTargetLimit: pick.msg.AggregatorTargetLimit,
					})
				}
				if swp.Deposit.GT(swp.In) {
					remainder := common.SafeSub(swp.Deposit, swp.In)
					source := pick.msg.Tx.Coins[0].Asset
					toi := TxOutItem{
						Chain:     source.GetChain(),
						InHash:    pick.msg.Tx.ID,
						ToAddress: pick.msg.Tx.FromAddress,
						Coin:      common.NewCoin(source, remainder),
					}
					tois = append(tois, toi)
				}

				for _, item := range tois {
					// let the txout manager mint our outbound asset if it is a synthetic asset
					if item.Chain.IsTHORChain() && (item.Coin.Asset.IsSyntheticAsset() || item.Coin.Asset.IsDerivedAsset()) {
						item.ModuleName = ModuleName
					}
					ok, err := mgr.TxOutStore().TryAddTxOutItem(ctx, mgr, item, cosmos.ZeroUint())
					if err != nil {
						return ErrInternal(err, "fail to add outbound tx")
					}
					if !ok {
						return errFailAddOutboundTx
					}
				}
				evt := NewEventStreamingSwap(pick.msg.Tx.Coins[0].Asset, pick.msg.TargetAsset, swp)
				if err := mgr.EventMgr().EmitEvent(ctx, evt); err != nil {
					ctx.Logger().Error("fail to emit streaming swap event", "error", err)
				}
			}
		} else {
			vm.k.RemoveSwapQueueItem(ctx, pick.msg.Tx.ID, pick.index)
		}
	}
	return nil
}

// getTodoNum - determine how many swaps to do.
func (vm *SwapQueueV117) getTodoNum(queueLen, minSwapsPerBlock, maxSwapsPerBlock int64) int64 {
	// Do half the length of the queue. Unless...
	//	1. The queue length is greater than maxSwapsPerBlock
	//  2. The queue length is less than minSwapsPerBlock
	todo := queueLen / 2
	if minSwapsPerBlock >= queueLen {
		todo = queueLen
	}
	if maxSwapsPerBlock < todo {
		todo = maxSwapsPerBlock
	}
	return todo
}

// scoreMsgs - this takes a list of MsgSwap, and converts them to a scored
// swapItem list
func (vm *SwapQueueV117) scoreMsgs(ctx cosmos.Context, items swapItems, synthVirtualDepthMult int64) (swapItems, error) {
	pools := make(map[common.Asset]Pool)

	for i, item := range items {
		// the asset customer send
		sourceAsset := item.msg.Tx.Coins[0].Asset
		// the asset customer want
		targetAsset := item.msg.TargetAsset

		for _, a := range []common.Asset{sourceAsset, targetAsset} {
			if a.IsRune() {
				continue
			}

			if _, ok := pools[a]; !ok {
				var err error
				pools[a], err = vm.k.GetPool(ctx, a.GetLayer1Asset())
				if err != nil {
					ctx.Logger().Error("fail to get pool", "pool", a, "error", err)
					continue
				}
			}
		}

		nonRuneAsset := sourceAsset
		if nonRuneAsset.IsRune() {
			nonRuneAsset = targetAsset
		}
		pool := pools[nonRuneAsset]
		if pool.IsEmpty() || pool.BalanceRune.IsZero() || pool.BalanceAsset.IsZero() {
			continue
		}
		// synths may be redeemed on unavailable pools, score them
		if !pool.IsAvailable() && !sourceAsset.IsSyntheticAsset() {
			continue
		}
		virtualDepthMult := int64(10_000)
		if nonRuneAsset.IsSyntheticAsset() {
			virtualDepthMult = synthVirtualDepthMult
		}
		vm.getLiquidityFeeAndSlip(ctx, pool, item.msg.Tx.Coins[0], &items[i], virtualDepthMult)

		if sourceAsset.IsRune() || targetAsset.IsRune() {
			// single swap , stop here
			continue
		}
		// double swap , thus need to convert source coin to RUNE and calculate fee and slip again
		runeCoin := common.NewCoin(common.RuneAsset(), pool.AssetValueInRune(item.msg.Tx.Coins[0].Amount))
		nonRuneAsset = targetAsset
		pool = pools[nonRuneAsset]
		if pool.IsEmpty() || !pool.IsAvailable() || pool.BalanceRune.IsZero() || pool.BalanceAsset.IsZero() {
			continue
		}
		virtualDepthMult = int64(10_000)
		if targetAsset.IsSyntheticAsset() {
			virtualDepthMult = synthVirtualDepthMult
		}
		vm.getLiquidityFeeAndSlip(ctx, pool, runeCoin, &items[i], virtualDepthMult)
	}

	return items, nil
}

// getLiquidityFeeAndSlip calculate liquidity fee and slip, fee is in RUNE
func (vm *SwapQueueV117) getLiquidityFeeAndSlip(ctx cosmos.Context, pool Pool, sourceCoin common.Coin, item *swapItem, virtualDepthMult int64) {
	// Get our X, x, Y values
	var X, x, Y cosmos.Uint
	x = sourceCoin.Amount
	if sourceCoin.Asset.IsRune() {
		X = pool.BalanceRune
		Y = pool.BalanceAsset
	} else {
		Y = pool.BalanceRune
		X = pool.BalanceAsset
	}

	X = common.GetUncappedShare(cosmos.NewUint(uint64(virtualDepthMult)), cosmos.NewUint(10_000), X)
	Y = common.GetUncappedShare(cosmos.NewUint(uint64(virtualDepthMult)), cosmos.NewUint(10_000), Y)

	swapper, err := GetSwapper(vm.k.GetVersion())
	if err != nil {
		ctx.Logger().Error("fail to fetch swapper", "error", err)
		swapper = newSwapperV92()
	}
	fee := swapper.CalcLiquidityFee(X, x, Y)
	if sourceCoin.Asset.IsRune() {
		fee = pool.AssetValueInRune(fee)
	}
	slip := swapper.CalcSwapSlip(X, x)
	item.fee = item.fee.Add(fee)
	item.slip = item.slip.Add(slip)
}
