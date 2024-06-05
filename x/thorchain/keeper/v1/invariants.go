package keeperv1

import (
	"fmt"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

// InvariantRoutes return the keeper's invariant routes
func (k KVStore) InvariantRoutes() []common.InvariantRoute {
	return []common.InvariantRoute{
		common.NewInvariantRoute("asgard", AsgardInvariant(k)),
		common.NewInvariantRoute("bond", BondInvariant(k)),
		common.NewInvariantRoute("thorchain", THORChainInvariant(k)),
		common.NewInvariantRoute("affiliate_collector", AffilliateCollectorInvariant(k)),
		common.NewInvariantRoute("pools", PoolsInvariant(k)),
		common.NewInvariantRoute("streaming_swaps", StreamingSwapsInvariant(k)),
	}
}

// AsgardInvariant the asgard module backs pool rune, savers synths, and native
// coins in queued swaps
func AsgardInvariant(k KVStore) common.Invariant {
	return func(ctx cosmos.Context) (msg []string, broken bool) {
		// sum all rune liquidity on pools, including pending
		var poolCoins common.Coins
		pools, _ := k.GetPools(ctx)
		for _, pool := range pools {
			switch {
			case pool.Asset.IsSyntheticAsset():
				coin := common.NewCoin(
					pool.Asset,
					pool.BalanceAsset,
				)
				poolCoins = poolCoins.Add_deprecated(coin)
			case !pool.Asset.IsDerivedAsset():
				coin := common.NewCoin(
					common.RuneAsset(),
					pool.BalanceRune.Add(pool.PendingInboundRune),
				)
				poolCoins = poolCoins.Add_deprecated(coin)
			}
		}

		// sum all rune in pending swaps
		var swapCoins common.Coins
		swapIter := k.GetSwapQueueIterator(ctx)
		defer swapIter.Close()
		for ; swapIter.Valid(); swapIter.Next() {
			var swap MsgSwap
			k.Cdc().MustUnmarshal(swapIter.Value(), &swap)

			if len(swap.Tx.Coins) != 1 {
				broken = true
				msg = append(msg, fmt.Sprintf("wrong number of coins for swap: %d, %s", len(swap.Tx.Coins), swap.Tx.ID))
				continue
			}

			coin := swap.Tx.Coins[0]
			if !coin.IsNative() && !swap.TargetAsset.IsNative() {
				continue // only verifying native coins in this invariant
			}

			// adjust for streaming swaps
			ss := swap.GetStreamingSwap() // GetStreamingSwap() rather than var so In.IsZero() doesn't panic
			// A non-streaming affiliate swap and streaming main swap could have the same TxID,
			// so explicitly check IsStreaming to not double-count the main swap's In and Out amounts.
			if swap.IsStreaming() {
				var err error
				ss, err = k.GetStreamingSwap(ctx, swap.Tx.ID)
				if err != nil {
					ctx.Logger().Error("error getting streaming swap", "error", err)
					continue // should never happen
				}
			}

			if coin.IsNative() {
				if !ss.In.IsZero() {
					// adjust for stream swap amount, the amount In has been added
					// to the pool but not deducted from the tx or module, so deduct
					// that In amount from the tx coin
					coin.Amount = coin.Amount.Sub(ss.In)
				}
				swapCoins = swapCoins.Add_deprecated(coin)
			}

			if swap.TargetAsset.IsNative() && !ss.Out.IsZero() {
				swapCoins = swapCoins.Add_deprecated(common.NewCoin(swap.TargetAsset, ss.Out))
			}
		}

		// get asgard module balance
		asgardAddr := k.GetModuleAccAddress(AsgardName)
		asgardCoins := k.GetBalance(ctx, asgardAddr)

		// asgard balance is expected to equal sum of pool and swap coins
		expNative, _ := poolCoins.Adds_deprecated(swapCoins).Native()

		// note: coins must be sorted for SafeSub
		diffCoins, _ := asgardCoins.SafeSub(expNative.Sort())
		if !diffCoins.IsZero() {
			broken = true
			for _, coin := range diffCoins {
				if coin.IsPositive() {
					msg = append(msg, fmt.Sprintf("oversolvent: %s", coin))
				} else {
					coin.Amount = coin.Amount.Neg()
					msg = append(msg, fmt.Sprintf("insolvent: %s", coin))
				}
			}
		}

		return msg, broken
	}
}

// BondInvariant the bond module backs node bond and pending reward bond
func BondInvariant(k KVStore) common.Invariant {
	return func(ctx cosmos.Context) (msg []string, broken bool) {
		// sum all rune bonded to nodes
		bondedRune := cosmos.ZeroUint()
		naIter := k.GetNodeAccountIterator(ctx)
		defer naIter.Close()
		for ; naIter.Valid(); naIter.Next() {
			var na NodeAccount
			k.Cdc().MustUnmarshal(naIter.Value(), &na)
			bondedRune = bondedRune.Add(na.Bond)
		}

		// get pending bond reward rune
		network, _ := k.GetNetwork(ctx)
		bondRewardRune := network.BondRewardRune

		// get rune balance of bond module
		bondModuleRune := k.GetBalanceOfModule(ctx, BondName, common.RuneAsset().Native())

		// bond module is expected to equal bonded rune and pending rewards
		expectedRune := bondedRune.Add(bondRewardRune)
		if expectedRune.GT(bondModuleRune) {
			broken = true
			diff := expectedRune.Sub(bondModuleRune)
			coin, _ := common.NewCoin(common.RuneAsset(), diff).Native()
			msg = append(msg, fmt.Sprintf("insolvent: %s", coin))

		} else if expectedRune.LT(bondModuleRune) {
			broken = true
			diff := bondModuleRune.Sub(expectedRune)
			coin, _ := common.NewCoin(common.RuneAsset(), diff).Native()
			msg = append(msg, fmt.Sprintf("oversolvent: %s", coin))
		}

		return msg, broken
	}
}

// THORChainInvariant the thorchain module should never hold a balance
func THORChainInvariant(k KVStore) common.Invariant {
	return func(ctx cosmos.Context) (msg []string, broken bool) {
		// module balance of theorchain
		tcAddr := k.GetModuleAccAddress(ModuleName)
		tcCoins := k.GetBalance(ctx, tcAddr)

		// thorchain module should never carry a balance
		if !tcCoins.Empty() {
			broken = true
			for _, coin := range tcCoins {
				msg = append(msg, fmt.Sprintf("oversolvent: %s", coin))
			}
		}

		return msg, broken
	}
}

// AffilliateCollectorInvariant the affiliate_collector module backs accrued affiliate
// rewards
func AffilliateCollectorInvariant(k KVStore) common.Invariant {
	return func(ctx cosmos.Context) (msg []string, broken bool) {
		affColModuleRune := k.GetBalanceOfModule(ctx, AffiliateCollectorName, common.RuneAsset().Native())
		affCols, err := k.GetAffiliateCollectors(ctx)
		if err != nil {
			if affColModuleRune.IsZero() {
				return nil, false
			}
			msg = append(msg, err.Error())
			return msg, true
		}

		totalAffRune := cosmos.ZeroUint()
		for _, ac := range affCols {
			totalAffRune = totalAffRune.Add(ac.RuneAmount)
		}

		if totalAffRune.GT(affColModuleRune) {
			broken = true
			diff := totalAffRune.Sub(affColModuleRune)
			coin, _ := common.NewCoin(common.RuneAsset(), diff).Native()
			msg = append(msg, fmt.Sprintf("insolvent: %s", coin))
		} else if totalAffRune.LT(affColModuleRune) {
			broken = true
			diff := affColModuleRune.Sub(totalAffRune)
			coin, _ := common.NewCoin(common.RuneAsset(), diff).Native()
			msg = append(msg, fmt.Sprintf("oversolvent: %s", coin))
		}

		return msg, broken
	}
}

// PoolsInvariant pool units and pending rune/asset should match the sum
// of units and pending rune/asset for all lps
func PoolsInvariant(k KVStore) common.Invariant {
	return func(ctx cosmos.Context) (msg []string, broken bool) {
		pools, _ := k.GetPools(ctx)
		for _, pool := range pools {
			if pool.Asset.IsNative() {
				continue // only looking at layer-one pools
			}

			lpUnits := cosmos.ZeroUint()
			lpPendingRune := cosmos.ZeroUint()
			lpPendingAsset := cosmos.ZeroUint()

			lpIter := k.GetLiquidityProviderIterator(ctx, pool.Asset)
			defer lpIter.Close()
			for ; lpIter.Valid(); lpIter.Next() {
				var lp LiquidityProvider
				k.Cdc().MustUnmarshal(lpIter.Value(), &lp)
				lpUnits = lpUnits.Add(lp.Units)
				lpPendingRune = lpPendingRune.Add(lp.PendingRune)
				lpPendingAsset = lpPendingAsset.Add(lp.PendingAsset)
			}

			check := func(poolValue, lpValue cosmos.Uint, valueType string) {
				if poolValue.GT(lpValue) {
					diff := poolValue.Sub(lpValue)
					msg = append(msg, fmt.Sprintf("%s oversolvent: %s %s", pool.Asset, diff.String(), valueType))
					broken = true
				} else if poolValue.LT(lpValue) {
					diff := lpValue.Sub(poolValue)
					msg = append(msg, fmt.Sprintf("%s insolvent: %s %s", pool.Asset, diff.String(), valueType))
					broken = true
				}
			}

			check(pool.LPUnits, lpUnits, "units")
			check(pool.PendingInboundRune, lpPendingRune, "pending rune")
			check(pool.PendingInboundAsset, lpPendingAsset, "pending asset")
		}

		return msg, broken
	}
}

// StreamingSwapsInvariant every streaming swap should have a corresponding
// queued swap, stream deposit should equal the queued swap's source coin,
// and the stream should be internally consistent
func StreamingSwapsInvariant(k KVStore) common.Invariant {
	return func(ctx cosmos.Context) (msg []string, broken bool) {
		// fetch all streaming swaps from the swap queue
		var swaps []MsgSwap
		swapIter := k.GetSwapQueueIterator(ctx)
		defer swapIter.Close()
		for ; swapIter.Valid(); swapIter.Next() {
			var swap MsgSwap
			k.Cdc().MustUnmarshal(swapIter.Value(), &swap)
			if swap.IsStreaming() {
				swaps = append(swaps, swap)
			}
		}

		// fetch all stream swap records
		var streams []StreamingSwap
		ssIter := k.GetStreamingSwapIterator(ctx)
		defer ssIter.Close()
		for ; ssIter.Valid(); ssIter.Next() {
			var stream StreamingSwap
			k.Cdc().MustUnmarshal(ssIter.Value(), &stream)
			streams = append(streams, stream)
		}

		for _, stream := range streams {
			found := false
			for _, swap := range swaps {
				if !swap.Tx.ID.Equals(stream.TxID) {
					continue
				}
				found = true
				if !swap.Tx.Coins[0].Amount.Equal(stream.Deposit) {
					broken = true
					msg = append(msg, fmt.Sprintf(
						"%s: swap.coin %s != stream.deposit %s",
						stream.TxID.String(),
						swap.Tx.Coins[0].Amount,
						stream.Deposit.String()))
				}
				if stream.Count > stream.Quantity {
					broken = true
					msg = append(msg, fmt.Sprintf(
						"%s: stream.count %d > stream.quantity %d",
						stream.TxID.String(),
						stream.Count,
						stream.Quantity))
				}
				if stream.In.GT(stream.Deposit) {
					broken = true
					msg = append(msg, fmt.Sprintf(
						"%s: stream.in %s > stream.deposit %s",
						stream.TxID.String(),
						stream.In.String(),
						stream.Deposit.String()))
				}
			}
			if !found {
				broken = true
				msg = append(msg, fmt.Sprintf("swap not found for stream: %s", stream.TxID.String()))
			}
		}

		return msg, broken
	}
}
