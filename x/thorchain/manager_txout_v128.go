package thorchain

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/blang/semver"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/types"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
	"gitlab.com/thorchain/thornode/mimir"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
)

// TxOutStorageV128 is going to manage all the outgoing tx
type TxOutStorageV128 struct {
	keeper        keeper.Keeper
	constAccessor constants.ConstantValues
	eventMgr      EventManager
	gasManager    GasManager
}

// newTxOutStorageV128 will create a new instance of TxOutStore.
func newTxOutStorageV128(keeper keeper.Keeper, constAccessor constants.ConstantValues, eventMgr EventManager, gasManager GasManager) *TxOutStorageV128 {
	return &TxOutStorageV128{
		keeper:        keeper,
		eventMgr:      eventMgr,
		constAccessor: constAccessor,
		gasManager:    gasManager,
	}
}

func (tos *TxOutStorageV128) EndBlock(ctx cosmos.Context, mgr Manager) error {
	// update the max gas for all outbounds in this block. This can be useful
	// if an outbound transaction was scheduled into the future, and the gas
	// for that blockchain changes in that time span. This avoids the need to
	// reschedule the transaction to Asgard.
	txOut, err := tos.GetBlockOut(ctx)
	if err != nil {
		return err
	}

	maxGasCache := make(map[common.Chain]common.Coin)
	gasRateCache := make(map[common.Chain]int64)

	for i, tx := range txOut.TxArray {
		voter, err := tos.keeper.GetObservedTxInVoter(ctx, tx.InHash)
		if err != nil {
			ctx.Logger().Error("fail to get observe tx in voter", "error", err)
			continue
		}

		// if the outbound height exists and is in the past, then no need to calculate new max gas
		if voter.OutboundHeight > 0 && voter.OutboundHeight < ctx.BlockHeight() {
			continue
		}

		// update max gas, take the larger of the current gas, or the last gas used

		// update cache if needed
		if _, ok := maxGasCache[tx.Chain]; !ok {
			maxGasCache[tx.Chain], _ = mgr.GasMgr().GetMaxGas(ctx, tx.Chain)
		}
		if _, ok := gasRateCache[tx.Chain]; !ok {
			gasRateCache[tx.Chain] = int64(mgr.GasMgr().GetGasRate(ctx, tx.Chain).Uint64())
		}

		maxGas := maxGasCache[tx.Chain]
		gasRate := gasRateCache[tx.Chain]
		if len(tx.MaxGas) == 0 || maxGas.Amount.GT(tx.MaxGas[0].Amount) {
			txOut.TxArray[i].MaxGas = common.Gas{maxGas}
			// Update MaxGas in ObservedTxVoter action as well
			err := updateTxOutGas(ctx, tos.keeper, tx, common.Gas{maxGas})
			if err != nil {
				ctx.Logger().Error("Failed to update MaxGas of action in ObservedTxVoter", "hash", tx.InHash, "error", err)
			}
		}
		// Equals checks GasRate so update actions GasRate too (before updating in the queue item)
		// for future updates of MaxGas, which must match for matchActionItem in AddOutTx.
		if err := updateTxOutGasRate(ctx, tos.keeper, tx, gasRate); err != nil {
			ctx.Logger().Error("Failed to update GasRate of action in ObservedTxVoter", "hash", tx.InHash, "error", err)
		}
		txOut.TxArray[i].GasRate = gasRate
	}

	if err := tos.keeper.SetTxOut(ctx, txOut); err != nil {
		return fmt.Errorf("fail to save tx out : %w", err)
	}
	return nil
}

// GetBlockOut read the TxOut from kv store
func (tos *TxOutStorageV128) GetBlockOut(ctx cosmos.Context) (*TxOut, error) {
	return tos.keeper.GetTxOut(ctx, ctx.BlockHeight())
}

// GetOutboundItems read all the outbound item from kv store
func (tos *TxOutStorageV128) GetOutboundItems(ctx cosmos.Context) ([]TxOutItem, error) {
	block, err := tos.keeper.GetTxOut(ctx, ctx.BlockHeight())
	if block == nil {
		return nil, nil
	}
	return block.TxArray, err
}

// GetOutboundItemByToAddress read all the outbound items filter by the given to address
func (tos *TxOutStorageV128) GetOutboundItemByToAddress(ctx cosmos.Context, to common.Address) []TxOutItem {
	filterItems := make([]TxOutItem, 0)
	items, _ := tos.GetOutboundItems(ctx)
	for _, item := range items {
		if item.ToAddress.Equals(to) {
			filterItems = append(filterItems, item)
		}
	}
	return filterItems
}

// ClearOutboundItems remove all the tx out items , mostly used for test
func (tos *TxOutStorageV128) ClearOutboundItems(ctx cosmos.Context) {
	_ = tos.keeper.ClearTxOut(ctx, ctx.BlockHeight())
}

// When TryAddTxOutItem returns an error, there should be no state changes from it,
// including funds movements or fee events from prepareTxOutItem.
// So, use CacheContext to only commit state changes when cachedTryAddTxOutItem doesn't return an error.
func (tos *TxOutStorageV128) TryAddTxOutItem(ctx cosmos.Context, mgr Manager, toi TxOutItem, minOut cosmos.Uint) (bool, error) {
	if toi.ToAddress.IsNoop() {
		return true, nil
	}
	cacheCtx, commit := ctx.CacheContext()
	success, err := tos.cachedTryAddTxOutItem(cacheCtx, mgr, toi, minOut)
	if err == nil {
		commit()
		ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
	}
	return success, err
}

// (cached)TryAddTxOutItem add an outbound tx to block
// return bool indicate whether the transaction had been added successful or not
// return error indicate error
func (tos *TxOutStorageV128) cachedTryAddTxOutItem(ctx cosmos.Context, mgr Manager, toi TxOutItem, minOut cosmos.Uint) (bool, error) {
	if toi.Coin.Asset.IsTradeAsset() {
		// no outbound needed for trade assets
		return true, nil
	}

	outputs, totalOutboundFeeRune, err := tos.prepareTxOutItem(ctx, toi)
	if err != nil {
		return false, fmt.Errorf("fail to prepare outbound tx: %w", err)
	}
	if len(outputs) == 0 {
		return false, ErrNotEnoughToPayFee
	}

	sumOut := cosmos.ZeroUint()
	for _, o := range outputs {
		sumOut = sumOut.Add(o.Coin.Amount)
	}
	if sumOut.LT(minOut) {
		// **NOTE** this error string is utilized by the order book manager to
		// catch the error. DO NOT change this error string without updating
		// the order book manager as well
		return false, fmt.Errorf("outbound amount does not meet requirements (%d/%d)", sumOut.Uint64(), minOut.Uint64())
	}

	// calculate the single block height to send all of these txout items,
	// using the summed amount
	outboundHeight := ctx.BlockHeight()
	cloutApplied := cosmos.ZeroUint()
	if !toi.Chain.IsTHORChain() && !toi.InHash.IsEmpty() && !toi.InHash.Equals(common.BlankTxID) {
		toi.Memo = outputs[0].Memo
		voter, err := tos.keeper.GetObservedTxInVoter(ctx, toi.InHash)
		if err != nil {
			ctx.Logger().Error("fail to get observe tx in voter", "error", err)
			return false, fmt.Errorf("fail to get observe tx in voter,err:%w", err)
		}

		var targetHeight int64
		targetHeight, cloutApplied, err = tos.CalcTxOutHeight(ctx, mgr.GetVersion(), toi)
		if err != nil {
			ctx.Logger().Error("failed to calc target block height for txout item", "error", err)
		}

		// adjust delay to include streaming swap time since inbound consensus
		if voter.Height > 0 {
			targetHeight = (targetHeight - ctx.BlockHeight()) + voter.Height
		}

		if targetHeight > outboundHeight {
			outboundHeight = targetHeight
		}

		// When the inbound transaction already has an outbound , the make sure the outbound will be scheduled on the same block
		if voter.OutboundHeight > 0 {
			outboundHeight = voter.OutboundHeight
		} else {
			voter.OutboundHeight = outboundHeight
			tos.keeper.SetObservedTxInVoter(ctx, voter)
		}
	}

	// sum total output asset
	sumOutput := cosmos.ZeroUint()
	for _, output := range outputs {
		sumOutput = sumOutput.Add(output.Coin.Amount)
	}

	// add tx to block out
	totalCloutShare := cosmos.ZeroUint()
	for i, output := range outputs {
		cloutShare := cosmos.ZeroUint()
		if i < len(outputs)-1 {
			cloutShare = common.GetSafeShare(output.Coin.Amount, sumOutput, cloutApplied)
			totalCloutShare = totalCloutShare.Add(cloutShare)
		} else {
			cloutShare = common.SafeSub(cloutApplied, totalCloutShare) // remainder
		}
		output.CloutSpent = &cloutShare
		if err := tos.addToBlockOut(ctx, mgr, output, outboundHeight); err != nil {
			return false, err
		}
	}

	// Add total outbound fee to the OutboundGasWithheldRune. totalOutboundFeeRune will be 0 if these are Migration outbounds
	// Don't count outbounds on THORChain ($RUNE and Synths)
	if !totalOutboundFeeRune.IsZero() && !toi.Chain.IsTHORChain() {
		network, err := tos.keeper.GetNetwork(ctx)
		if err != nil {
			ctx.Logger().Error("fail to get network data", "error", err)
		} else {
			network.OutboundGasWithheldRune += totalOutboundFeeRune.Uint64()
			if err := tos.keeper.SetNetwork(ctx, network); err != nil {
				ctx.Logger().Error("fail to set network data", "error", err)
			}
		}
	}

	return true, nil
}

// UnSafeAddTxOutItem - blindly adds a tx out, skipping vault selection, transaction
// fee deduction, etc
func (tos *TxOutStorageV128) UnSafeAddTxOutItem(ctx cosmos.Context, mgr Manager, toi TxOutItem, height int64) error {
	if toi.ToAddress.IsNoop() {
		return nil
	}
	// BCH chain will convert legacy address to new format automatically , thus when observe it back can't be associated with the original inbound
	// so here convert the legacy address to new format
	if toi.Chain.Equals(common.BCHChain) {
		newBCHAddress, err := common.ConvertToNewBCHAddressFormatV83(toi.ToAddress)
		if err != nil {
			return fmt.Errorf("fail to convert BCH address to new format: %w", err)
		}
		if newBCHAddress.IsEmpty() {
			return fmt.Errorf("empty to address , can't send out")
		}
		toi.ToAddress = newBCHAddress
	}
	return tos.addToBlockOut(ctx, mgr, toi, height)
}

func (tos *TxOutStorageV128) DiscoverOutbounds(ctx cosmos.Context, transactionFeeAsset cosmos.Uint, maxGasAsset common.Coin, toi TxOutItem, vaults Vaults) ([]TxOutItem, cosmos.Uint) {
	var outputs []TxOutItem

	// When there is more than one vault, sort the vaults by
	// (as an integer) how many vaults of that size
	// would be necessary to fulfill the outbound (smallest number first).
	// Having already been sorted by security, for a given vaults-necessary
	// the lowest security ones will still be ordered first.
	// The greater a vault's vaults-necessary, the less its security would be
	// decreased by taking part in the outbound;
	// also, outbounds from negligible-amount vaults (other than wasting gas) risk creating
	// duplicate txout items of which all but one would be stuck in the outbound queue.
	// Note that for vaults of equal (integer) vaults-necessary, any previous sort order remains.
	if len(vaults) > 1 {
		type VaultsNecessary struct {
			Vault    Vault
			Estimate uint64
		}

		vaultsNecessary := make([]VaultsNecessary, 0)

		for _, vault := range vaults {
			// Avoid a divide-by-zero by ignoring vaults with zero of the asset.
			if vault.GetCoin(toi.Coin.Asset).Amount.IsZero() {
				continue
			}

			// if vault is frozen, don't send more txns to sign, as they may be
			// delayed. Once a txn is skipped here, it will not be rescheduled again.
			if len(vault.Frozen) > 0 {
				chains, err := common.NewChains(vault.Frozen)
				if err != nil {
					ctx.Logger().Error("failed to convert chains", "error", err)
				}
				if chains.Has(maxGasAsset.Asset.GetChain()) {
					continue
				}
			}

			vaultsNecessary = append(vaultsNecessary, VaultsNecessary{
				Vault:    vault,
				Estimate: toi.Coin.Amount.Quo(vault.GetCoin(toi.Coin.Asset).Amount).Uint64(),
			})
		}

		// If more than one vault remains, sort by vaults-necessary ascending.
		if len(vaultsNecessary) > 1 {
			sort.SliceStable(vaultsNecessary, func(i, j int) bool {
				return vaultsNecessary[i].Estimate < vaultsNecessary[j].Estimate
			})
		}

		// Set 'vaults' to the sorted order.
		vaults = make(Vaults, len(vaultsNecessary))
		for i, v := range vaultsNecessary {
			vaults[i] = v.Vault
		}
	}

	for _, vault := range vaults {
		// Ensure THORNode are not sending from and to the same address
		fromAddr, err := vault.PubKey.GetAddress(toi.Chain)
		if err != nil || fromAddr.IsEmpty() || toi.ToAddress.Equals(fromAddr) {
			continue
		}
		// if the asset in the vault is not enough to pay for the fee , then skip it
		if vault.GetCoin(toi.Coin.Asset).Amount.LTE(transactionFeeAsset) {
			continue
		}
		// if the vault doesn't have gas asset in it , or it doesn't have enough to pay for gas
		gasAsset := vault.GetCoin(toi.Chain.GetGasAsset())
		if gasAsset.IsEmpty() || gasAsset.Amount.LT(maxGasAsset.Amount) {
			continue
		}

		toi.VaultPubKey = vault.PubKey
		if toi.Coin.Amount.LTE(vault.GetCoin(toi.Coin.Asset).Amount) {
			outputs = append(outputs, toi)
			toi.Coin.Amount = cosmos.ZeroUint()
			break
		} else {
			remainingAmount := common.SafeSub(toi.Coin.Amount, vault.GetCoin(toi.Coin.Asset).Amount)
			toi.Coin.Amount = common.SafeSub(toi.Coin.Amount, remainingAmount)
			outputs = append(outputs, toi)
			toi.Coin.Amount = remainingAmount
		}
	}
	return outputs, toi.Coin.Amount
}

// prepareTxOutItem will do some data validation which include the following
// 1. Make sure it has a legitimate memo
// 2. choose an appropriate vault(s) to send from (active asgard, then retiring asgard)
// 3. deduct transaction fee, keep in mind, only take transaction fee when active nodes are  more then minimumBFT
// return list of outbound transactions
func (tos *TxOutStorageV128) prepareTxOutItem(ctx cosmos.Context, toi TxOutItem) ([]TxOutItem, types.Uint, error) {
	var outputs []TxOutItem
	var remaining cosmos.Uint

	// Default the memo to the standard outbound memo
	if toi.Memo == "" {
		toi.Memo = NewOutboundMemo(toi.InHash).String()
	}

	// Ensure the InHash is set
	if toi.InHash.IsEmpty() {
		toi.InHash = common.BlankTxID
	} else {
		// fetch inbound txn memo, and append arbitrary data (if applicable)
		inboundVoter, err := tos.keeper.GetObservedTxInVoter(ctx, toi.InHash)
		if err == nil {
			parts := strings.SplitN(inboundVoter.Tx.Tx.Memo, "|", 2)
			if len(parts) == 2 {
				toi.Memo = fmt.Sprintf("%s|%s", toi.Memo, parts[1])
				if len(toi.Memo) > constants.MaxMemoSize {
					toi.Memo = toi.Memo[:constants.MaxMemoSize]
				}
			}
		}
	}
	if toi.ToAddress.IsEmpty() {
		return outputs, cosmos.ZeroUint(), fmt.Errorf("empty to address, can't send out")
	}
	if !toi.ToAddress.IsChain(toi.Chain) {
		return outputs, cosmos.ZeroUint(), fmt.Errorf("to address(%s), is not of chain(%s)", toi.ToAddress, toi.Chain)
	}

	// BCH chain will convert legacy address to new format automatically , thus when observe it back can't be associated with the original inbound
	// so here convert the legacy address to new format
	if toi.Chain.Equals(common.BCHChain) {
		newBCHAddress, err := common.ConvertToNewBCHAddressFormatV83(toi.ToAddress)
		if err != nil {
			return outputs, cosmos.ZeroUint(), fmt.Errorf("fail to convert BCH address to new format: %w", err)
		}
		if newBCHAddress.IsEmpty() {
			return outputs, cosmos.ZeroUint(), fmt.Errorf("empty to address , can't send out")
		}
		toi.ToAddress = newBCHAddress
	}

	// ensure amount is rounded to appropriate decimals
	toiPool, err := tos.keeper.GetPool(ctx, toi.Coin.Asset.GetLayer1Asset())
	if err != nil {
		return nil, cosmos.ZeroUint(), fmt.Errorf("fail to get pool for txout manager: %w", err)
	}

	signingTransactionPeriod := tos.constAccessor.GetInt64Value(constants.SigningTransactionPeriod)
	transactionFeeRune := tos.gasManager.GetFee(ctx, toi.Chain, common.RuneAsset())
	transactionFeeAsset := tos.gasManager.GetFee(ctx, toi.Chain, toi.Coin.Asset)
	maxGasAsset, err := tos.gasManager.GetMaxGas(ctx, toi.Chain)
	if err != nil {
		ctx.Logger().Error("fail to get max gas asset", "error", err)
	}
	if toi.Chain.IsTHORChain() {
		outputs = append(outputs, toi)
	} else {
		if !toi.VaultPubKey.IsEmpty() {
			// a vault is already manually selected, blindly go forth with that
			outputs = append(outputs, toi)
		} else {
			// THORNode don't have a vault already selected to send from, discover one.
			// List all pending outbounds for the asset, this will be used
			// to deduct balances of vaults that have outstanding txs assigned
			pendingOutbounds := tos.keeper.GetPendingOutbounds(ctx, toi.Coin.Asset)

			// ///////////// COLLECT ACTIVE ASGARD VAULTS ///////////////////
			activeAsgards, err := tos.keeper.GetAsgardVaultsByStatus(ctx, ActiveVault)
			if err != nil {
				ctx.Logger().Error("fail to get active vaults", "error", err)
			}

			// All else being equal, prefer lower-security vaults for outbounds.
			activeAsgards = tos.keeper.SortBySecurity(ctx, activeAsgards, signingTransactionPeriod)

			for i := range activeAsgards {
				// having sorted by security, deduct the value of any assigned pending outbounds
				activeAsgards[i].DeductVaultPendingOutbounds(pendingOutbounds)
			}
			// //////////////////////////////////////////////////////////////

			// ///////////// COLLECT RETIRING ASGARD VAULTS /////////////////
			retiringAsgards, err := tos.keeper.GetAsgardVaultsByStatus(ctx, RetiringVault)
			if err != nil {
				ctx.Logger().Error("fail to get retiring vaults", "error", err)
			}

			// All else being equal, prefer lower-security vaults for outbounds.
			retiringAsgards = tos.keeper.SortBySecurity(ctx, retiringAsgards, signingTransactionPeriod)

			for i := range retiringAsgards {
				// having sorted by security, deduct the value of any assigned pending outbounds
				retiringAsgards[i].DeductVaultPendingOutbounds(pendingOutbounds)
			}
			// //////////////////////////////////////////////////////////////

			// iterate over discovered vaults and find vaults to send funds from

			// All else being equal, choose active Asgards over retiring Asgards.
			outputs, remaining = tos.DiscoverOutbounds(ctx, transactionFeeAsset, maxGasAsset, toi, append(activeAsgards, retiringAsgards...))

			// Check we found enough funds to satisfy the request, error if we didn't
			if !remaining.IsZero() {
				return nil, cosmos.ZeroUint(), fmt.Errorf("insufficient funds for outbound request: %s %s remaining", toi.ToAddress.String(), remaining.String())
			}
		}
	}
	var finalOutput []TxOutItem
	var pool Pool
	var feeEvents []*EventFee
	finalRuneFee := cosmos.ZeroUint()
	for i := range outputs {
		if outputs[i].MaxGas.IsEmpty() {
			maxGasCoin, err := tos.gasManager.GetMaxGas(ctx, outputs[i].Chain)
			if err != nil {
				return nil, cosmos.ZeroUint(), fmt.Errorf("fail to get max gas coin: %w", err)
			}
			outputs[i].MaxGas = common.Gas{
				maxGasCoin,
			}
			// THOR Chain doesn't need to have max gas
			if outputs[i].MaxGas.IsEmpty() && !outputs[i].Chain.Equals(common.THORChain) {
				return nil, cosmos.ZeroUint(), fmt.Errorf("max gas cannot be empty: %s", outputs[i].MaxGas)
			}

			outputs[i].GasRate = int64(tos.gasManager.GetGasRate(ctx, outputs[i].Chain).Uint64())
		}

		runeFee := transactionFeeRune // Fee is the prescribed fee

		// get the lending address to avoid deducting the outbound fee
		lendAddr, err := tos.keeper.GetModuleAddress(LendingName)
		if err != nil {
			return nil, cosmos.ZeroUint(), fmt.Errorf("fail to get lending address: %w", err)
		}

		// Deduct OutboundTransactionFee from TOI and add to Reserve
		memo, err := ParseMemoWithTHORNames(ctx, tos.keeper, outputs[i].Memo)
		if err == nil && !memo.IsType(TxMigrate) && !memo.IsType(TxRagnarok) && !toi.ToAddress.Equals(lendAddr) {
			if outputs[i].Coin.Asset.IsRune() {
				if outputs[i].Coin.Amount.LTE(transactionFeeRune) {
					runeFee = outputs[i].Coin.Amount // Fee is the full amount
				}
				finalRuneFee = finalRuneFee.Add(runeFee)
				outputs[i].Coin.Amount = common.SafeSub(outputs[i].Coin.Amount, runeFee)
				fee := common.NewFee(common.Coins{common.NewCoin(outputs[i].Coin.Asset, runeFee)}, cosmos.ZeroUint())
				feeEvents = append(feeEvents, NewEventFee(outputs[i].InHash, fee, cosmos.ZeroUint()))
			} else {
				if pool.IsEmpty() {
					var err error
					pool, err = tos.keeper.GetPool(ctx, toi.Coin.Asset.GetLayer1Asset()) // Get pool
					if err != nil {
						// the error is already logged within kvstore
						return nil, cosmos.ZeroUint(), fmt.Errorf("fail to get pool: %w", err)
					}
				}

				// if pool units is zero, no asset fee is taken
				if !pool.GetPoolUnits().IsZero() {
					assetFee := transactionFeeAsset
					if outputs[i].Coin.Amount.LTE(assetFee) {
						assetFee = outputs[i].Coin.Amount // Fee is the full amount
					}

					outputs[i].Coin.Amount = common.SafeSub(outputs[i].Coin.Amount, assetFee) // Deduct Asset fee
					if outputs[i].Coin.Asset.IsNative() {
						// burn the native asset which used to pay for fee, that's only required when sending from asgard
						if outputs[i].GetModuleName() == AsgardName {
							if err := tos.keeper.SendFromModuleToModule(ctx,
								AsgardName,
								ModuleName,
								common.NewCoins(common.NewCoin(outputs[i].Coin.Asset, assetFee))); err != nil {
								ctx.Logger().Error("fail to move native asset fee from asgard to Module", "error", err)
							} else {
								if err := tos.keeper.BurnFromModule(ctx, ModuleName, common.NewCoin(outputs[i].Coin.Asset, assetFee)); err != nil {
									ctx.Logger().Error("fail to burn native asset", "error", err)
								} else {
									burnEvt := NewEventMintBurn(BurnSupplyType, outputs[i].Coin.Asset.Native(), assetFee, "burn_native_fee")
									if err := tos.eventMgr.EmitEvent(ctx, burnEvt); err != nil {
										ctx.Logger().Error("fail to emit burn event", "error", err)
									}
								}
							}
						}
					}
					var poolDeduct cosmos.Uint
					runeFee = pool.RuneDisbursementForAssetAdd(assetFee)
					if runeFee.GT(pool.BalanceRune) {
						poolDeduct = pool.BalanceRune
					} else {
						poolDeduct = runeFee
					}
					finalRuneFee = finalRuneFee.Add(poolDeduct)
					if !outputs[i].Coin.Asset.IsSyntheticAsset() {
						pool.BalanceAsset = pool.BalanceAsset.Add(assetFee) // Add Asset fee to Pool
					}
					pool.BalanceRune = common.SafeSub(pool.BalanceRune, poolDeduct) // Deduct Rune from Pool
					fee := common.NewFee(common.Coins{common.NewCoin(outputs[i].Coin.Asset, assetFee)}, poolDeduct)
					feeEvents = append(feeEvents, NewEventFee(outputs[i].InHash, fee, cosmos.ZeroUint()))
				}
			}
		}

		vault, err := tos.keeper.GetVault(ctx, outputs[i].VaultPubKey)
		if err != nil && !outputs[i].Chain.IsTHORChain() {
			// For THORChain outputs (since having an empty VaultPubKey)
			// GetVault is expected to fail, so do not log the error.
			ctx.Logger().Error("fail to get vault", "error", err)
		}
		// when it is ragnarok , the network doesn't charge fee , however if the output asset is gas asset,
		// then the amount of max gas need to be taken away from the customer , otherwise the vault will be insolvent and doesn't
		// have enough to fulfill outbound
		// Also the MaxGas has not put back to pool ,so there is no need to subside pool when ragnarok is in progress
		// OR, if the vault is inactive, subtract maxgas from amount so we have gas to spend to refund the txn
		if (memo.IsType(TxRagnarok) || vault.Status == InactiveVault) && outputs[i].Coin.Asset.IsGasAsset() {
			gasAmt := outputs[i].MaxGas.ToCoins().GetCoin(outputs[i].Coin.Asset).Amount
			outputs[i].Coin.Amount = common.SafeSub(outputs[i].Coin.Amount, gasAmt)
		}
		if outputs[i].Coin.IsEmpty() {
			ctx.Logger().Info("tx out item has zero coin", "tx_out", outputs[i].String())

			// Need to determinate whether the outbound is triggered by a withdrawal request
			// if the outbound is trigger by withdrawal request, and emit asset is not enough to pay for the fee
			// this need to return with an error , thus handler_withdraw can restore LP's LPUnits
			// and also the fee event will not be emitted
			if !outputs[i].InHash.IsEmpty() && !outputs[i].InHash.Equals(common.BlankTxID) {
				inboundVoter, err := tos.keeper.GetObservedTxInVoter(ctx, outputs[i].InHash)
				if err != nil {
					ctx.Logger().Error("fail to get observed txin voter", "error", err)
					continue
				}
				if inboundVoter.Tx.IsEmpty() {
					continue
				}
				inboundMemo, err := ParseMemoWithTHORNames(ctx, tos.keeper, inboundVoter.Tx.Tx.Memo)
				if err != nil {
					ctx.Logger().Error("fail to parse inbound transaction memo", "error", err)
					continue
				}
				if inboundMemo.IsType(TxWithdraw) {
					return nil, cosmos.ZeroUint(), errors.New("tx out item has zero coin")
				}
			}
			continue
		}

		// If the outbound coin is synthetic, respecting decimals is unnecessary
		// and leaves unburnt synths in the Pool Module
		if !outputs[i].Coin.Asset.IsSyntheticAsset() {
			// sanity check: ensure outbound amount respect asset decimals
			outputs[i].Coin.Amount = cosmos.RoundToDecimal(outputs[i].Coin.Amount, toiPool.Decimals)
		}

		if !outputs[i].InHash.Equals(common.BlankTxID) {
			// increment out number of out tx for this in tx
			voter, err := tos.keeper.GetObservedTxInVoter(ctx, outputs[i].InHash)
			if err != nil {
				return nil, cosmos.ZeroUint(), fmt.Errorf("fail to get observed tx voter: %w", err)
			}
			voter.FinalisedHeight = ctx.BlockHeight()
			voter.Actions = append(voter.Actions, outputs[i])
			tos.keeper.SetObservedTxInVoter(ctx, voter)
		}

		finalOutput = append(finalOutput, outputs[i])
	}

	if !pool.IsEmpty() {
		if err := tos.keeper.SetPool(ctx, pool); err != nil { // Set Pool
			return nil, cosmos.ZeroUint(), fmt.Errorf("fail to save pool: %w", err)
		}
	}
	for _, feeEvent := range feeEvents {
		if err := tos.eventMgr.EmitFeeEvent(ctx, feeEvent); err != nil {
			ctx.Logger().Error("fail to emit fee event", "error", err)
		}
	}
	if !finalRuneFee.IsZero() {
		if toi.Coin.Asset.IsNativeRune() {
			// If the source module is the Reserve, leave the fee in the Reserve without a transfer.
			if toi.ModuleName != ReserveName {
				sourceModule := toi.GetModuleName() // Ensure that non-"".
				coin := common.NewCoin(common.RuneAsset(), finalRuneFee)
				err := tos.keeper.SendFromModuleToModule(ctx, sourceModule, ReserveName, common.NewCoins(coin))
				if err != nil {
					ctx.Logger().Error("fail to send fee to reserve", "error", err, "module", sourceModule)
				}
			}
		} else {
			// Layer 1 or Synth Asset is implicitly swapped in a pool
			// whether in vault or burnt from another network module,
			// but Derived Asset has no outbound fee taken
			// so that the emitted amount passed to the loan handler
			// and the amount transferred to the Lending module are the same.
			// (If a fee were taken, then being for a Derived Asset pool
			//  it would contribute to Lending breathing room
			//  rather than affecting Pool Module RUNE.)
			if !toi.Coin.Asset.IsDerivedAsset() {
				if err := tos.keeper.AddPoolFeeToReserve(ctx, finalRuneFee); err != nil {
					ctx.Logger().Error("fail to add pool fee to reserve", "error", err)
				}
			}
		}
	}

	return finalOutput, finalRuneFee, nil
}

func (tos *TxOutStorageV128) addToBlockOut(ctx cosmos.Context, mgr Manager, item TxOutItem, outboundHeight int64) error {
	// if we're sending native assets, transfer them now and return
	if item.Chain.IsTHORChain() {
		return tos.nativeTxOut(ctx, mgr, item)
	}

	vault, err := tos.keeper.GetVault(ctx, item.VaultPubKey)
	if err != nil {
		ctx.Logger().Error("fail to get vault", "error", err)
	}
	memo, _ := ParseMemo(mgr.GetVersion(), item.Memo) // ignore err
	labels := []metrics.Label{
		telemetry.NewLabel("vault_type", vault.Type.String()),
		telemetry.NewLabel("pubkey", item.VaultPubKey.String()),
		telemetry.NewLabel("memo_type", memo.GetType().String()),
	}
	telemetry.SetGaugeWithLabels([]string{"thornode", "vault", "out_txn"}, float32(1), labels)

	if err := tos.eventMgr.EmitEvent(ctx, NewEventScheduledOutbound(item)); err != nil {
		ctx.Logger().Error("fail to emit scheduled outbound event", "error", err)
	}

	return tos.keeper.AppendTxOut(ctx, outboundHeight, item)
}

func (tos *TxOutStorageV128) calcClout(ctx cosmos.Context, runeValue cosmos.Uint, toi TxOutItem) (cosmos.Uint, cosmos.Uint) {
	cloutOut, err := tos.keeper.GetSwapperClout(ctx, toi.ToAddress)
	if err != nil {
		ctx.Logger().Error("fail to get swapper clout destination address", "error", err)
	}
	voter, err := tos.keeper.GetObservedTxInVoter(ctx, toi.InHash)
	if err != nil {
		ctx.Logger().Error("fail to get txin for clout calculation", "error", err)
	}
	cloutIn, err := tos.keeper.GetSwapperClout(ctx, voter.Tx.Tx.FromAddress)
	if err != nil {
		ctx.Logger().Error("fail to get swapper clout destination address", "error", err)
	}

	swapperCloutReset := mimir.NewSwapperCloutReset().FetchValue(ctx, tos.keeper)
	swapperCloutLimit := mimir.NewSwapperCloutLimit().FetchValue(ctx, tos.keeper)

	// if last clout spend was over an hour ago, restore clout available to
	// 100%
	cloutOut.Restore(ctx.BlockHeight(), swapperCloutReset)
	cloutIn.Restore(ctx.BlockHeight(), swapperCloutReset)

	clout1, clout2, newValue := tos.splitClout(
		ctx,
		cosmos.SafeUintFromInt64(swapperCloutLimit),
		cloutIn.Available(),
		cloutOut.Available(),
		runeValue,
	)

	// sanity check, newValue + clout1 + clout2 should equal runeValue
	if !newValue.Add(clout1).Add(clout2).Equal(runeValue) {
		return runeValue, cosmos.ZeroUint()
	}

	if !clout1.IsZero() {
		cloutIn.Spent = cloutIn.Spent.Add(clout1)
		cloutIn.LastSpentHeight = ctx.BlockHeight()
		if err := tos.keeper.SetSwapperClout(ctx, cloutIn); err != nil {
			ctx.Logger().Error("fail to save swapper clout", "error", err)
		}
	}

	if !clout2.IsZero() {
		cloutOut.Spent = cloutOut.Spent.Add(clout2)
		cloutOut.LastSpentHeight = ctx.BlockHeight()
		if err := tos.keeper.SetSwapperClout(ctx, cloutOut); err != nil {
			ctx.Logger().Error("fail to save swapper clout", "error", err)
		}
	}

	return newValue, clout1.Add(clout2)
}

// splitClout tries to split runeValue into two Uints, ensuring that it doesn't exceed the given clout1 and clout2.
func (tos *TxOutStorageV128) splitClout(ctx cosmos.Context, swapperCloutLimit, clout1, clout2, runeValue cosmos.Uint) (cosmos.Uint, cosmos.Uint, cosmos.Uint) {
	if clout1.Add(clout2).GT(swapperCloutLimit) {
		halfLimit := swapperCloutLimit.QuoUint64(2)
		switch {
		case clout1.GT(halfLimit) && clout2.GT(halfLimit):
			clout1 = halfLimit
			clout2 = halfLimit
		case clout1.GT(clout2):
			clout1 = common.SafeSub(swapperCloutLimit, clout2)
		case clout2.GT(clout1):
			clout2 = common.SafeSub(swapperCloutLimit, clout1)
		}
	}

	// sanity check - ensure total available clout does not exceed our limit
	if clout1.Add(clout2).GT(swapperCloutLimit) {
		ctx.Logger().Error("dev error: clout1 + clout2 cannot exceed clout limit", "clout1", clout1, "clout2", clout2, "clout limit", swapperCloutLimit)
		return cosmos.ZeroUint(), cosmos.ZeroUint(), runeValue
	}

	totalClout := clout1.Add(clout2)
	if totalClout.IsZero() {
		return cosmos.ZeroUint(), cosmos.ZeroUint(), runeValue
	}

	appliedClout := cosmos.MinUint(totalClout, runeValue)
	amountFromClout1 := appliedClout.Mul(clout1).Quo(totalClout)
	amountFromClout2 := appliedClout.Sub(amountFromClout1)

	return amountFromClout1, amountFromClout2, common.SafeSub(runeValue, amountFromClout1.Add(amountFromClout2))
}

func (tos *TxOutStorageV128) CalcTxOutHeight(ctx cosmos.Context, version semver.Version, toi TxOutItem) (int64, cosmos.Uint, error) {
	// non-outbound transactions are skipped. This is so this code does not
	// affect internal transactions (ie consolidation and migrate txs)
	memo, _ := ParseMemo(version, toi.Memo) // ignore err
	if !memo.IsType(TxRefund) && !memo.IsType(TxOutbound) {
		return ctx.BlockHeight(), cosmos.ZeroUint(), nil
	}

	minTxOutVolumeThreshold := tos.keeper.GetConfigInt64(ctx, constants.MinTxOutVolumeThreshold)
	txOutDelayRate := tos.keeper.GetConfigInt64(ctx, constants.TxOutDelayRate)
	txOutDelayMax := tos.keeper.GetConfigInt64(ctx, constants.TxOutDelayMax)
	maxTxOutOffset := tos.keeper.GetConfigInt64(ctx, constants.MaxTxOutOffset)

	// only delay if volume threshold and delay rate are positive
	if minTxOutVolumeThreshold <= 0 || txOutDelayRate <= 0 || maxTxOutOffset <= 0 {
		return ctx.BlockHeight(), cosmos.ZeroUint(), nil
	}

	// convert to big ints for safer math
	minVolumeThreshold := cosmos.NewUint(uint64(minTxOutVolumeThreshold))
	delayRate := cosmos.NewUint(uint64(txOutDelayRate))
	maxOffset := cosmos.NewUint(uint64(maxTxOutOffset))

	// get txout item value in rune
	runeValue, _ := tos.keeper.GetTOIsValue(ctx, toi)

	// reduce rune value based on clout
	runeValue, cloutApplied := tos.calcClout(ctx, runeValue, toi)
	// if clout was large enough to cover the outbound value, no delay applied
	if runeValue.IsZero() {
		return ctx.BlockHeight(), cloutApplied, nil
	}

	// sum value of scheduled txns (including this one)
	sumValue := runeValue
	cloutValue := cosmos.ZeroUint()
	for height := ctx.BlockHeight() + 1; height <= ctx.BlockHeight()+txOutDelayMax; height++ {
		value, clout, err := tos.keeper.GetTxOutValue(ctx, height)
		if err != nil {
			ctx.Logger().Error("fail to get tx out array from key value store", "error", err)
			continue
		}
		if height > ctx.BlockHeight()+maxTxOutOffset && value.IsZero() {
			// we've hit our max offset, and an empty block, we can assume the
			// rest will be empty as well
			break
		}
		sumValue = sumValue.Add(value)
		cloutValue = cloutValue.Add(clout)
	}

	// reduce delay rate relative to the total scheduled value. In high volume
	// scenarios, this causes the network to send outbound transactions slower,
	// giving the community & NOs time to analyze and react. In an attack
	// scenario, the attacker is likely going to move as much value as possible
	// (as we've seen in the past). The act of doing this will slow down their
	// own transaction(s), reducing the attack's effectiveness.
	// The common.One is because delayRate, sumValue, and minVolumeThreshold
	// all have the same number of decimals (which cancel otherwise).
	rateReduction := cosmos.NewUint(common.One).Mul(common.SafeSub(sumValue, cloutValue)).Quo(minVolumeThreshold)
	if rateReduction.GTE(delayRate) {
		delayRate = cosmos.NewUint(1)
	} else {
		delayRate = delayRate.Sub(rateReduction)
	}

	// calculate the minimum number of blocks in the future the txn has to be.
	// min shouldn't be anything longer than the max txout offset
	minBlocks := runeValue.Quo(delayRate)
	if minBlocks.GT(maxOffset) {
		minBlocks = maxOffset
	}
	targetBlock := ctx.BlockHeight() + int64(minBlocks.Uint64())

	// find targetBlock that has space for new txout item.
	count := int64(0)
	for count < txOutDelayMax { // max set 1 day into the future
		txOutValue, _, err := tos.keeper.GetTxOutValue(ctx, targetBlock)
		if err != nil {
			ctx.Logger().Error("fail to get txOutValue for block height", "error", err)
			break
		}
		if txOutValue.IsZero() {
			// the txout has no outbound txns, let's use this one
			break
		}
		if txOutValue.Add(runeValue).LTE(minVolumeThreshold) {
			// the txout + this txout item has enough space to fit, lets use this one
			break
		}
		targetBlock++
		count++
	}

	return targetBlock, cloutApplied, nil
}

func (tos *TxOutStorageV128) nativeTxOut(ctx cosmos.Context, mgr Manager, toi TxOutItem) error {
	addr, err := cosmos.AccAddressFromBech32(toi.ToAddress.String())
	if err != nil {
		return err
	}

	toi.ModuleName = toi.GetModuleName() // Ensure that non-"".

	// mint if we're sending from THORChain module
	if toi.ModuleName == ModuleName {
		if err := tos.keeper.MintToModule(ctx, toi.ModuleName, toi.Coin); err != nil {
			return fmt.Errorf("fail to mint coins during txout: %w", err)
		}
		mintEvt := NewEventMintBurn(MintSupplyType, toi.Coin.Asset.Native(), toi.Coin.Amount, "native_tx_out")
		if err := tos.eventMgr.EmitEvent(ctx, mintEvt); err != nil {
			ctx.Logger().Error("fail to emit mint event", "error", err)
		}
	}

	polAddress, err := tos.keeper.GetModuleAddress(ReserveName)
	if err != nil {
		ctx.Logger().Error("fail to get from address", "err", err)
		return err
	}

	affColAddress, err := tos.keeper.GetModuleAddress(AffiliateCollectorName)
	if err != nil {
		ctx.Logger().Error("fail to get from address", "err", err)
		return err
	}

	// send funds to/from modules
	var sdkErr error
	switch {
	case polAddress.Equals(toi.ToAddress):
		sdkErr = tos.keeper.SendFromModuleToModule(ctx, toi.ModuleName, ReserveName, common.NewCoins(toi.Coin))
	case affColAddress.Equals(toi.ToAddress):
		sdkErr = tos.keeper.SendFromModuleToModule(ctx, toi.ModuleName, AffiliateCollectorName, common.NewCoins(toi.Coin))
	default:
		sdkErr = tos.keeper.SendFromModuleToAccount(ctx, toi.ModuleName, addr, common.NewCoins(toi.Coin))
	}

	if sdkErr != nil {
		return errors.New(sdkErr.Error())
	}

	from, err := tos.keeper.GetModuleAddress(toi.ModuleName)
	if err != nil {
		ctx.Logger().Error("fail to get from address", "err", err)
		return err
	}
	outboundTxFee := tos.keeper.GetOutboundTxFee(ctx)

	tx := common.NewTx(
		common.BlankTxID,
		from,
		toi.ToAddress,
		common.Coins{toi.Coin},
		common.Gas{common.NewCoin(common.RuneAsset(), outboundTxFee)},
		toi.Memo,
	)

	active, err := tos.keeper.GetAsgardVaultsByStatus(ctx, ActiveVault)
	if err != nil {
		ctx.Logger().Error("fail to get active vaults", "err", err)
		return err
	}

	if len(active) == 0 {
		return fmt.Errorf("dev error: no pubkey for native txn")
	}

	observedTx := ObservedTx{
		ObservedPubKey: active[0].PubKey,
		BlockHeight:    ctx.BlockHeight(),
		Tx:             tx,
		FinaliseHeight: ctx.BlockHeight(),
	}
	m, err := processOneTxIn(ctx, mgr.GetVersion(), tos.keeper, observedTx, tos.keeper.GetModuleAccAddress(AsgardName))
	if err != nil {
		ctx.Logger().Error("fail to process txOut", "error", err, "tx", tx.String())
		return err
	}

	handler := NewInternalHandler(mgr)

	_, err = handler(ctx, m)
	if err != nil {
		ctx.Logger().Error("TxOut Handler failed:", "error", err)
		return err
	}

	return nil
}
