package thorchain

import (
	"context"
	"fmt"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/blang/semver"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/cosmos/cosmos-sdk/types"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
)

// CommonOutboundTxHandler is the place where those common logic can be shared
// between multiple different kind of outbound tx handler
// at the moment, handler_refund, and handler_outbound_tx are largely the same
// , only some small difference
type CommonOutboundTxHandler struct {
	mgr Manager
}

// NewCommonOutboundTxHandler create a new instance of the CommonOutboundTxHandler
func NewCommonOutboundTxHandler(mgr Manager) CommonOutboundTxHandler {
	return CommonOutboundTxHandler{
		mgr: mgr,
	}
}

func (h CommonOutboundTxHandler) slashV96(ctx cosmos.Context, tx ObservedTx) error {
	toSlash := make(common.Coins, len(tx.Tx.Coins))
	copy(toSlash, tx.Tx.Coins)
	toSlash = toSlash.Adds_deprecated(tx.Tx.Gas.ToCoins())

	ctx = ctx.WithContext(context.WithValue(ctx.Context(), constants.CtxMetricLabels, []metrics.Label{ // nolint
		telemetry.NewLabel("reason", "failed_outbound"),
		telemetry.NewLabel("chain", string(tx.Tx.Chain)),
	}))

	return h.mgr.Slasher().SlashVault(ctx, tx.ObservedPubKey, toSlash, h.mgr)
}

func (h CommonOutboundTxHandler) handle(ctx cosmos.Context, tx ObservedTx, inTxID common.TxID) (*cosmos.Result, error) {
	version := h.mgr.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.127.0")):
		return h.handleV127(ctx, tx, inTxID)
	case version.GTE(semver.MustParse("1.125.0")):
		return h.handleV125(ctx, tx, inTxID)
	case version.GTE(semver.MustParse("1.118.0")):
		return h.handleV118(ctx, tx, inTxID)
	case version.GTE(semver.MustParse("1.98.0")):
		return h.handleV98(ctx, tx, inTxID)
	case version.GTE(semver.MustParse("1.96.0")):
		return h.handleV96(ctx, tx, inTxID)
	case version.GTE(semver.MustParse("1.94.0")):
		return h.handleV94(ctx, tx, inTxID)
	case version.GTE(semver.MustParse("1.92.0")):
		return h.handleV92(ctx, tx, inTxID)
	case version.GTE(semver.MustParse("1.88.0")):
		return h.handleV88(ctx, tx, inTxID)
	case version.GTE(semver.MustParse("1.87.0")):
		return h.handleV87(ctx, tx, inTxID)
	case version.GTE(semver.MustParse("1.85.0")):
		return h.handleV85(ctx, tx, inTxID)
	case version.GTE(semver.MustParse("0.69.0")):
		return h.handleV69(ctx, tx, inTxID)
	}
	return nil, errBadVersion
}

func (h CommonOutboundTxHandler) handleV127(ctx cosmos.Context, tx ObservedTx, inTxID common.TxID) (*cosmos.Result, error) {
	// note: Outbound tx usually it is related to an inbound tx except migration
	// thus here try to get the ObservedTxInVoter,  and set the tx out hash accordingly
	voter, err := h.mgr.Keeper().GetObservedTxInVoter(ctx, inTxID)
	if err != nil {
		return nil, ErrInternal(err, "fail to get observed tx voter")
	}
	if voter.AddOutTx(h.mgr.GetVersion(), tx.Tx) {
		if err := h.mgr.EventMgr().EmitEvent(ctx, NewEventOutbound(inTxID, tx.Tx)); err != nil {
			return nil, ErrInternal(err, "fail to emit outbound event")
		}
	}
	h.mgr.Keeper().SetObservedTxInVoter(ctx, voter)

	if tx.Tx.Chain.Equals(common.THORChain) {
		return &cosmos.Result{}, nil
	}

	shouldSlash := true
	signingTransPeriod := h.mgr.GetConstants().GetInt64Value(constants.SigningTransactionPeriod)
	// every Signing Transaction Period , THORNode will check whether a
	// TxOutItem had been sent by signer or not
	// if a txout item that is older than SigningTransactionPeriod, but has not
	// been sent out by signer , LackSigning will create a new TxOutItem
	// and mark the previous TxOutItem as complete.
	//
	// Check the blocks backwards
	// (assuming that most signed outbounds will have been scheduled more rather than less recently)
	// to the current height starting from inbound consensus height
	// or one signingTransPeriod ago, whichever is later.
	earliestHeight := ctx.BlockHeight() - signingTransPeriod
	if voter.Height > earliestHeight {
		earliestHeight = voter.Height
	}
	for height := ctx.BlockHeight(); height >= earliestHeight; height-- {
		// update txOut record with our TxID that sent funds out of the pool
		txOut, err := h.mgr.Keeper().GetTxOut(ctx, height)
		if err != nil {
			ctx.Logger().Error("unable to get txOut record", "error", err)
			return nil, cosmos.ErrUnknownRequest(err.Error())
		}

		// Save TxOut back with the TxID only when the TxOut on the block height is
		// not empty
		for i, txOutItem := range txOut.TxArray {
			// withdraw , refund etc, one inbound tx might result two outbound
			// txes, THORNode have to correlate outbound tx back to the
			// inbound, and also txitem , thus THORNode could record both
			// outbound tx hash correctly given every tx item will only have
			// one coin in it , THORNode could use that to identify which tx it
			// is

			if txOutItem.InHash.Equals(inTxID) &&
				txOutItem.OutHash.IsEmpty() &&
				tx.Tx.Chain.Equals(txOutItem.Chain) &&
				tx.Tx.ToAddress.Equals(txOutItem.ToAddress) &&
				strings.EqualFold(tx.Aggregator, txOutItem.Aggregator) &&
				strings.EqualFold(tx.AggregatorTarget, txOutItem.AggregatorTargetAsset) &&
				tx.ObservedPubKey.Equals(txOutItem.VaultPubKey) {

				matchCoin := tx.Tx.Coins.EqualsEx(common.Coins{txOutItem.Coin})
				if !matchCoin {
					// In case the mismatch is caused by decimals , round the tx out item's amount , and compare it again
					p, err := h.mgr.Keeper().GetPool(ctx, txOutItem.Coin.Asset)
					if err != nil {
						ctx.Logger().Error("fail to get pool", "error", err)
					}
					if !p.IsEmpty() {
						matchCoin = tx.Tx.Coins.EqualsEx(common.Coins{
							common.NewCoin(txOutItem.Coin.Asset, cosmos.RoundToDecimal(txOutItem.Coin.Amount, p.Decimals)),
						})
					}
				}
				// when outbound is gas asset
				if !matchCoin && txOutItem.Coin.Asset.Equals(txOutItem.Chain.GetGasAsset()) {
					asset := txOutItem.Chain.GetGasAsset()
					intendToSpend := txOutItem.Coin.Amount.Add(txOutItem.MaxGas.ToCoins().GetCoin(asset).Amount)
					actualSpend := tx.Tx.Coins.GetCoin(asset).Amount.Add(tx.Tx.Gas.ToCoins().GetCoin(asset).Amount)
					if intendToSpend.Equal(actualSpend) {
						matchCoin = true
						maxGasAmt := txOutItem.MaxGas.ToCoins().GetCoin(asset).Amount
						realGasAmt := tx.Tx.Gas.ToCoins().GetCoin(asset).Amount
						ctx.Logger().Info("override match coin", "intend to spend", intendToSpend, "actual spend", actualSpend, "max_gas", maxGasAmt, "actual gas", realGasAmt)
						if maxGasAmt.GT(realGasAmt) {
							// Don't reimburse gas difference if the outbound is from an InactiveVault.
							vault, err := h.mgr.Keeper().GetVault(ctx, tx.ObservedPubKey)
							if err != nil {
								ctx.Logger().Error("fail to get vault", "error", err)
							}
							if vault.Status != InactiveVault {
								// the outbound spend less than MaxGas
								diffGas := maxGasAmt.Sub(realGasAmt)
								h.mgr.GasMgr().AddGasAsset(txOutItem.Coin.Asset, common.Gas{
									common.NewCoin(asset, diffGas),
								}, false)
							}
						} else if maxGasAmt.LT(realGasAmt) {
							// signer spend more than the maximum gas prescribed by THORChain , slash it
							ctx.Logger().Info("slash node", "max gas", maxGasAmt, "real gas spend", realGasAmt, "gap", common.SafeSub(realGasAmt, maxGasAmt).String())
							matchCoin = false
						}
					}
				}
				if txOutItem.Chain.Equals(common.ETHChain) {
					maxGasAmount := txOutItem.MaxGas.ToCoins().GetCoin(common.ETHAsset).Amount
					gasAmount := tx.Tx.Gas.ToCoins().GetCoin(common.ETHAsset).Amount
					// If thornode instruct bifrost to spend more than MaxETHGas , then it should not be slashed.
					if gasAmount.GTE(cosmos.NewUint(constants.MaxETHGas)) && maxGasAmount.LT(cosmos.NewUint(constants.MaxETHGas)) {
						ctx.Logger().Info("ETH chain transaction spend more than MaxETHGas should be slashed", "gas", gasAmount.String(), "max gas", constants.MaxETHGas)
						matchCoin = false
					}
				}

				if !matchCoin {
					continue
				}
				txOut.TxArray[i].OutHash = tx.Tx.ID
				shouldSlash = false
				if err := h.mgr.Keeper().SetTxOut(ctx, txOut); err != nil {
					ctx.Logger().Error("fail to save tx out", "error", err)
				}

				// reclaim clout spent
				outTxn := txOut.TxArray[i]
				spent := outTxn.CloutSpent
				if spent != nil && !spent.IsZero() {
					cloutOut, err := h.mgr.Keeper().GetSwapperClout(ctx, outTxn.ToAddress)
					if err != nil {
						ctx.Logger().Error("fail to get swapper clout destination address", "error", err)
					}
					voter, err := h.mgr.Keeper().GetObservedTxInVoter(ctx, outTxn.InHash)
					if err != nil {
						ctx.Logger().Error("fail to get txin for clout calculation", "error", err)
					}
					cloutIn, err := h.mgr.Keeper().GetSwapperClout(ctx, voter.Tx.Tx.FromAddress)
					if err != nil {
						ctx.Logger().Error("fail to get swapper clout destination address", "error", err)
					}

					clout1, clout2 := calcReclaim(cloutIn.Claimable(), cloutOut.Claimable(), *spent)

					cloutIn.Reclaim(clout1)
					cloutIn.LastReclaimHeight = ctx.BlockHeight()
					if err := h.mgr.Keeper().SetSwapperClout(ctx, cloutIn); err != nil {
						ctx.Logger().Error("fail to save swapper clout in", "error", err)
					}

					cloutOut.Reclaim(clout2)
					cloutOut.LastReclaimHeight = ctx.BlockHeight()
					if err := h.mgr.Keeper().SetSwapperClout(ctx, cloutOut); err != nil {
						ctx.Logger().Error("fail to save swapper clout out", "error", err)
					}
				}

				break

			}
		}
		// If the TxOutItem matching the observed outbound has been found,
		// do not check other blocks.
		if !shouldSlash {
			break
		}
	}

	if shouldSlash {
		ctx.Logger().Info("slash node account, no matched tx out item", "inbound txid", inTxID, "outbound tx", tx.Tx)

		// send security alert for events that are not evm burn
		if !isOutboundFakeGasTX(tx) {
			msg := fmt.Sprintf("missing tx out in=%s", inTxID)
			if err := h.mgr.EventMgr().EmitEvent(ctx, NewEventSecurity(tx.Tx, msg)); err != nil {
				ctx.Logger().Error("fail to emit security event", "error", err)
			}
		}

		if err := h.slashV96(ctx, tx); err != nil {
			return nil, ErrInternal(err, "fail to slash account")
		}
	}

	if err := h.mgr.Keeper().SetLastSignedHeight(ctx, voter.FinalisedHeight); err != nil {
		ctx.Logger().Info("fail to update last signed height", "error", err)
	}

	return &cosmos.Result{}, nil
}

// calcReclaim attempts to split spent clout between two reclaimable clouts as equally as possible.
func calcReclaim(reclaimable1, reclaimable2, spent cosmos.Uint) (reclaim1, reclaim2 cosmos.Uint) {
	// Ensure that the spent clout doesn't exceed the total reclaimable clout
	totalReclaimable := reclaimable1.Add(reclaimable2)
	if spent.GT(totalReclaimable) {
		return reclaimable1, reclaimable2
	}

	// Split the spent clout in half
	halfSpent := spent.Quo(types.NewUint(2))

	// If either clout is less than half the spent amount, allocate all to that clout
	if reclaimable1.LT(halfSpent) {
		return reclaimable1, spent.Sub(reclaimable1)
	} else if reclaimable2.LT(halfSpent) {
		return spent.Sub(reclaimable2), reclaimable2
	}

	// Otherwise, split the spent clout equally
	return halfSpent, spent.Sub(halfSpent)
}

// isOutboundFakeGasTX returns true if the observed outbound which is missing an inbound is a "fake" tx sent purposely by bifrost
// this occurs on EVM chains where an outbound cannot be sent generally due to out of gas failures. In these cases, bifrost
// signs an outbound with the gas spent on the failure so that it can be accounted for.
// Checks are:
// - must only have one coin in outbound
// - chain of coin must be an EVM chain
// - coin asset must be the gas asset
// - coin amount must be 1
func isOutboundFakeGasTX(tx ObservedTx) bool {
	isLenCoins1 := len(tx.Tx.Coins) == 1
	if !isLenCoins1 {
		return false
	}
	asset := tx.Tx.Coins[0].Asset
	isChainEVM := asset.Chain.IsEVM()
	if !isChainEVM {
		return false
	}
	gasAsset := asset.Chain.GetGasAsset()
	isAssetGasAsset := asset.Equals(gasAsset)
	if !isAssetGasAsset {
		return false
	}

	if !tx.Tx.Coins[0].Amount.Equal(types.NewUint(1)) {
		return false
	}

	return true
}
