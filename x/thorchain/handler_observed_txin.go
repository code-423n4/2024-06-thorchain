package thorchain

import (
	"context"
	"fmt"

	"github.com/armon/go-metrics"
	"github.com/blang/semver"
	"github.com/cosmos/cosmos-sdk/telemetry"
	se "github.com/cosmos/cosmos-sdk/types/errors"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

// ObservedTxInHandler to handle MsgObservedTxIn
type ObservedTxInHandler struct {
	mgr Manager
}

// NewObservedTxInHandler create a new instance of ObservedTxInHandler
func NewObservedTxInHandler(mgr Manager) ObservedTxInHandler {
	return ObservedTxInHandler{
		mgr: mgr,
	}
}

// Run is the main entry point of ObservedTxInHandler
func (h ObservedTxInHandler) Run(ctx cosmos.Context, m cosmos.Msg) (*cosmos.Result, error) {
	msg, ok := m.(*MsgObservedTxIn)
	if !ok {
		return nil, errInvalidMessage
	}
	err := h.validate(ctx, *msg)
	if err != nil {
		ctx.Logger().Error("MsgObservedTxIn failed validation", "error", err)
		return nil, err
	}

	result, err := h.handle(ctx, *msg)
	if err != nil {
		ctx.Logger().Error("fail to handle MsgObservedTxIn message", "error", err)
	}
	return result, err
}

func (h ObservedTxInHandler) validate(ctx cosmos.Context, msg MsgObservedTxIn) error {
	version := h.mgr.GetVersion()
	if version.GTE(semver.MustParse("0.1.0")) {
		return h.validateV1(ctx, msg)
	}
	return errInvalidVersion
}

func (h ObservedTxInHandler) validateV1(ctx cosmos.Context, msg MsgObservedTxIn) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	if !isSignedByActiveNodeAccounts(ctx, h.mgr.Keeper(), msg.GetSigners()) {
		return cosmos.ErrUnauthorized(fmt.Sprintf("%+v are not authorized", msg.GetSigners()))
	}

	return nil
}

func (h ObservedTxInHandler) handle(ctx cosmos.Context, msg MsgObservedTxIn) (*cosmos.Result, error) {
	version := h.mgr.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.131.0")):
		return h.handleV131(ctx, msg)
	case version.GTE(semver.MustParse("1.129.0")):
		return h.handleV129(ctx, msg)
	case version.GTE(semver.MustParse("1.128.0")):
		return h.handleV128(ctx, msg)
	case version.GTE(semver.MustParse("1.124.0")):
		return h.handleV124(ctx, msg)
	case version.GTE(semver.MustParse("1.116.0")):
		return h.handleV116(ctx, msg)
	case version.GTE(semver.MustParse("1.113.0")):
		return h.handleV113(ctx, msg)
	case version.GTE(semver.MustParse("1.112.0")):
		return h.handleV112(ctx, msg)
	case version.GTE(semver.MustParse("1.107.0")):
		return h.handleV107(ctx, msg)
	case version.GTE(semver.MustParse("1.89.0")):
		return h.handleV89(ctx, msg)
	case version.GTE(semver.MustParse("0.78.0")):
		return h.handleV78(ctx, msg)
	}
	return nil, errBadVersion
}

func (h ObservedTxInHandler) preflight(ctx cosmos.Context, voter ObservedTxVoter, nas NodeAccounts, tx ObservedTx, signer cosmos.AccAddress) (ObservedTxVoter, bool) {
	version := h.mgr.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.123.0")):
		return h.preflightV123(ctx, voter, nas, tx, signer)
	case version.GTE(semver.MustParse("1.119.0")):
		return h.preflightV119(ctx, voter, nas, tx, signer)
	case version.GTE(semver.MustParse("1.116.0")):
		return h.preflightV116(ctx, voter, nas, tx, signer)
	default:
		return h.preflightV1(ctx, voter, nas, tx, signer)
	}
}

func (h ObservedTxInHandler) preflightV123(ctx cosmos.Context, voter ObservedTxVoter, nas NodeAccounts, tx ObservedTx, signer cosmos.AccAddress) (ObservedTxVoter, bool) {
	observeSlashPoints := h.mgr.GetConstants().GetInt64Value(constants.ObserveSlashPoints)
	observeFlex := h.mgr.Keeper().GetConfigInt64(ctx, constants.ObservationDelayFlexibility)

	slashCtx := ctx.WithContext(context.WithValue(ctx.Context(), constants.CtxMetricLabels, []metrics.Label{
		telemetry.NewLabel("reason", "failed_observe_txin"),
		telemetry.NewLabel("chain", string(tx.Tx.Chain)),
	}))
	h.mgr.Slasher().IncSlashPoints(slashCtx, observeSlashPoints, signer)

	ok := false
	if err := h.mgr.Keeper().SetLastObserveHeight(ctx, tx.Tx.Chain, signer, tx.BlockHeight); err != nil {
		ctx.Logger().Error("fail to save last observe height", "error", err, "signer", signer, "chain", tx.Tx.Chain)
	}
	if !voter.Add(tx, signer) {
		return voter, ok
	}
	if voter.HasFinalised(nas) {
		if voter.FinalisedHeight == 0 {
			ok = true
			voter.Height = ctx.BlockHeight() // Always record the consensus height of the finalised Tx
			voter.FinalisedHeight = ctx.BlockHeight()
			voter.Tx = voter.GetTx(nas)
			// tx has consensus now, so decrease the slashing points for all the signers whom had voted for it
			h.mgr.Slasher().DecSlashPoints(slashCtx, observeSlashPoints, voter.GetConsensusSigners()...)
		} else if ctx.BlockHeight() <= (voter.FinalisedHeight+observeFlex) && voter.Tx.Equals(tx) {
			// event the tx had been processed , given the signer just a bit late , so still take away their slash points
			// but only when the tx signer are voting is the tx that already reached consensus
			h.mgr.Slasher().DecSlashPoints(slashCtx, observeSlashPoints, signer)
		}
	}
	if !ok && voter.HasConsensus(nas) && !tx.IsFinal() && voter.FinalisedHeight == 0 {
		if voter.Height == 0 {
			ok = true
			voter.Height = ctx.BlockHeight()
			// this is the tx that has consensus
			voter.Tx = voter.GetTx(nas)

			// tx has consensus now, so decrease the slashing points for all the signers whom had voted for it
			h.mgr.Slasher().DecSlashPoints(slashCtx, observeSlashPoints, voter.Tx.GetSigners()...)
		} else if ctx.BlockHeight() <= (voter.Height+observeFlex) && voter.Tx.Equals(tx) {
			// event the tx had been processed , given the signer just a bit late , so still take away their slash points
			// but only when the tx signer are voting is the tx that already reached consensus
			h.mgr.Slasher().DecSlashPoints(slashCtx, observeSlashPoints, signer)
		}
	}

	h.mgr.Keeper().SetObservedTxInVoter(ctx, voter)

	// Check to see if we have enough identical observations to process the transaction
	return voter, ok
}

func (h ObservedTxInHandler) handleV131(ctx cosmos.Context, msg MsgObservedTxIn) (*cosmos.Result, error) {
	activeNodeAccounts, err := h.mgr.Keeper().ListActiveValidators(ctx)
	if err != nil {
		return nil, wrapError(ctx, err, "fail to get list of active node accounts")
	}
	handler := NewInternalHandler(h.mgr)
	for _, tx := range msg.Txs {
		// check we are sending to a valid vault
		if !h.mgr.Keeper().VaultExists(ctx, tx.ObservedPubKey) {
			ctx.Logger().Info("Not valid Observed Pubkey", "observed pub key", tx.ObservedPubKey)
			continue
		}

		voter, err := h.mgr.Keeper().GetObservedTxInVoter(ctx, tx.Tx.ID)
		if err != nil {
			ctx.Logger().Error("fail to get tx in voter", "error", err)
			continue
		}

		voter, isConsensus := h.preflight(ctx, voter, activeNodeAccounts, tx, msg.Signer)
		if !isConsensus {
			if voter.Height == ctx.BlockHeight() || voter.FinalisedHeight == ctx.BlockHeight() {
				// we've already process the transaction, but we should still
				// update the observing addresses
				h.mgr.ObMgr().AppendObserver(tx.Tx.Chain, msg.GetSigners())
			}
			continue
		}

		// all logic after this is upon consensus

		ctx.Logger().Info("handleMsgObservedTxIn request", "Tx:", tx.String())
		if voter.Reverted {
			ctx.Logger().Info("tx had been reverted", "Tx", tx.String())
			continue
		}

		vault, err := h.mgr.Keeper().GetVault(ctx, tx.ObservedPubKey)
		if err != nil {
			ctx.Logger().Error("fail to get vault", "error", err)
			continue
		}

		voter.Tx.Tx.Memo = tx.Tx.Memo

		hasFinalised := voter.HasFinalised(activeNodeAccounts)
		// memo errors are ignored here and will be caught later in processing,
		// after vault update, voter setup, etc and the coin will be refunded
		memo, _ := ParseMemoWithTHORNames(ctx, h.mgr.Keeper(), tx.Tx.Memo)

		// Update vault balances from inbounds with Migrate memos immediately,
		// to minimise any gap between outbound and inbound observations.
		// TODO: In future somehow update both balances in a single action,
		// so the ActiveVault balance increase is guaranteed to never be early nor late?
		if hasFinalised || memo.IsType(TxMigrate) {
			if vault.IsAsgard() && !voter.UpdatedVault {
				if !tx.Tx.FromAddress.Equals(tx.Tx.ToAddress) {
					// Don't add to or subtract from vault balances when the sender and recipient are the same
					// (particularly avoid Consolidate SafeSub zeroing of vault balances).
					vault.AddFunds(tx.Tx.Coins)
					vault.InboundTxCount++
				}
				voter.UpdatedVault = true
			}
		}
		if err := h.mgr.Keeper().SetLastChainHeight(ctx, tx.Tx.Chain, tx.BlockHeight); err != nil {
			ctx.Logger().Error("fail to set last chain height", "error", err)
		}

		// save the changes in Tx Voter to key value store
		h.mgr.Keeper().SetObservedTxInVoter(ctx, voter)
		if err := h.mgr.Keeper().SetVault(ctx, vault); err != nil {
			ctx.Logger().Error("fail to set vault", "error", err)
			continue
		}

		if !vault.IsAsgard() {
			ctx.Logger().Info("Vault is not an Asgard vault, transaction ignored.")
			continue
		}

		if memo.IsOutbound() || memo.IsInternal() {
			// do not process outbound handlers here, or internal handlers
			continue
		}

		// add addresses to observing addresses. This is used to detect
		// active/inactive observing node accounts

		h.mgr.ObMgr().AppendObserver(tx.Tx.Chain, voter.Tx.GetSigners())

		if !hasFinalised {
			ctx.Logger().Info("Tx has not been finalised yet , waiting for confirmation counting", "hash", voter.TxID)
			continue
		}

		if vault.Status == InactiveVault {
			ctx.Logger().Error("observed tx on inactive vault", "tx", tx.String())
			if newErr := refundTx(ctx, tx, h.mgr, CodeInvalidVault, "observed inbound tx to an inactive vault", ""); newErr != nil {
				ctx.Logger().Error("fail to refund", "error", newErr)
			}
			continue
		}

		// construct msg from memo
		m, txErr := processOneTxIn(ctx, h.mgr.GetVersion(), h.mgr.Keeper(), voter.Tx, msg.Signer)
		if txErr != nil {
			ctx.Logger().Error("fail to process inbound tx", "error", txErr.Error(), "tx hash", tx.Tx.ID.String())
			if newErr := refundTx(ctx, tx, h.mgr, CodeInvalidMemo, txErr.Error(), ""); nil != newErr {
				ctx.Logger().Error("fail to refund", "error", err)
			}
			continue
		}

		// check if we've halted trading
		swapMsg, isSwap := m.(*MsgSwap)
		_, isAddLiquidity := m.(*MsgAddLiquidity)

		if isSwap || isAddLiquidity {
			if h.mgr.Keeper().IsTradingHalt(ctx, m) || h.mgr.Keeper().RagnarokInProgress(ctx) {
				if newErr := refundTx(ctx, tx, h.mgr, se.ErrUnauthorized.ABCICode(), "trading halted", ""); nil != newErr {
					ctx.Logger().Error("fail to refund for halted trading", "error", err)
				}
				continue
			}
		}

		// if its a swap, send it to our queue for processing later
		if isSwap {
			h.addSwap(ctx, *swapMsg)
			continue
		}

		// if it is a loan, inject the observed TxID and ToAddress into the context
		_, isLoanOpen := m.(*MsgLoanOpen)
		_, isLoanRepayment := m.(*MsgLoanRepayment)
		mCtx := ctx
		if isLoanOpen || isLoanRepayment {
			mCtx = ctx.WithValue(constants.CtxLoanTxID, tx.Tx.ID)
			mCtx = mCtx.WithValue(constants.CtxLoanToAddress, tx.Tx.ToAddress)
		}

		_, err = handler(mCtx, m)
		if err != nil {
			if err := refundTx(ctx, tx, h.mgr, CodeTxFail, err.Error(), ""); err != nil {
				ctx.Logger().Error("fail to refund", "error", err)
			}
			continue
		}

		// if an outbound is not expected, mark the voter as done
		if !memo.GetType().HasOutbound() {
			// retrieve the voter from store in case the handler caused a change
			voter, err := h.mgr.Keeper().GetObservedTxInVoter(ctx, tx.Tx.ID)
			if err != nil {
				return nil, fmt.Errorf("fail to get voter")
			}
			voter.SetDone()
			h.mgr.Keeper().SetObservedTxInVoter(ctx, voter)
		}
	}
	return &cosmos.Result{}, nil
}

func (h ObservedTxInHandler) addSwap(ctx cosmos.Context, msg MsgSwap) {
	version := h.mgr.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.116.0")):
		h.addSwapV116(ctx, msg)
	case version.GTE(semver.MustParse("1.98.0")):
		h.addSwapV98(ctx, msg)
	default:
		h.addSwapV63(ctx, msg)
	}
}

func (h ObservedTxInHandler) addSwapV116(ctx cosmos.Context, msg MsgSwap) {
	if h.mgr.Keeper().OrderBooksEnabled(ctx) {
		// TODO: swap to synth if layer1 asset (follow on PR)
		// TODO: create handler to modify/cancel an order (follow on PR)

		source := msg.Tx.Coins[0]
		target := common.NewCoin(msg.TargetAsset, msg.TradeTarget)
		evt := NewEventLimitOrder(source, target, msg.Tx.ID)
		if err := h.mgr.EventMgr().EmitEvent(ctx, evt); err != nil {
			ctx.Logger().Error("fail to emit swap event", "error", err)
		}
		if err := h.mgr.Keeper().SetOrderBookItem(ctx, msg); err != nil {
			ctx.Logger().Error("fail to add swap to queue", "error", err)
		}
	} else {
		h.addSwapDirect(ctx, msg)
	}
}

// addSwapDirect adds the swap directly to the swap queue (no order book) - segmented
// out into its own function to allow easier maintenance of original behavior vs order
// book behavior.
func (h ObservedTxInHandler) addSwapDirect(ctx cosmos.Context, msg MsgSwap) {
	version := h.mgr.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.132.0")):
		h.addSwapDirectV132(ctx, msg)
	case version.GTE(semver.MustParse("1.116.0")):
		h.addSwapDirectV116(ctx, msg)
	default:
		h.addSwapV63(ctx, msg)
	}
}

func (h ObservedTxInHandler) addSwapDirectV132(ctx cosmos.Context, msg MsgSwap) {
	if msg.Tx.Coins.IsEmpty() {
		return
	}
	amt := cosmos.ZeroUint()

	// Check if affiliate fee should be paid out
	if !msg.AffiliateBasisPoints.IsZero() && msg.AffiliateAddress.IsChain(common.THORChain) {
		amt = common.GetSafeShare(
			msg.AffiliateBasisPoints,
			cosmos.NewUint(10000),
			msg.Tx.Coins[0].Amount,
		)
		msg.Tx.Coins[0].Amount = common.SafeSub(msg.Tx.Coins[0].Amount, amt)
	}

	// Queue the main swap
	if err := h.mgr.Keeper().SetSwapQueueItem(ctx, msg, 0); err != nil {
		ctx.Logger().Error("fail to add swap to queue", "error", err)
	}

	// Affiliate fee flow
	if !amt.IsZero() {
		affiliateSwap := NewMsgSwap(
			msg.Tx,
			common.RuneAsset(),
			msg.AffiliateAddress,
			cosmos.ZeroUint(),
			common.NoAddress,
			cosmos.ZeroUint(),
			"",
			"", nil,
			MarketOrder,
			0, 0,
			msg.Signer,
		)

		var affThorname *types.THORName
		memo, err := ParseMemoWithTHORNames(ctx, h.mgr.Keeper(), msg.Tx.Memo)
		if err != nil {
			ctx.Logger().Error("fail to parse swap memo", "memo", msg.Tx.Memo, "error", err)
		} else {
			affThorname = memo.GetAffiliateTHORName()
		}

		// PreferredAsset set, swap to the AffiliateCollector Module + check if the
		// preferred asset swap should be triggered
		if affThorname != nil && !affThorname.PreferredAsset.IsEmpty() {
			affcol, err := h.mgr.Keeper().GetAffiliateCollector(ctx, affThorname.Owner)
			if err != nil {
				ctx.Logger().Error("failed to get affiliate collector for thorname", "thorname", affThorname.Name, "error", err)
				return
			}

			affColAddress, err := h.mgr.Keeper().GetModuleAddress(AffiliateCollectorName)
			if err != nil {
				ctx.Logger().Error("failed to retrieve the affiliate collector module address", "error", err)
				return
			}

			// Set AffiliateCollector Module as destination and populate the AffiliateAddress
			// so that the swap handler can increment the emitted RUNE for the affiliate in
			// the AffiliateCollector KVStore.
			affiliateSwap.Destination = affColAddress
			affiliateSwap.AffiliateAddress = msg.AffiliateAddress

			// Check if accrued RUNE is 100x current outbound fee of preferred asset chain, if
			// so trigger the preferred asset swap
			ofRune, err := h.mgr.GasMgr().GetAssetOutboundFee(ctx, affThorname.PreferredAsset, true)
			if err != nil {
				ctx.Logger().Error("failed to get outbound fee for preferred asset, skipping preferred asset swap", "name", affThorname.Name, "asset", affThorname.PreferredAsset, "error", err)
			}
			multiplier := h.mgr.Keeper().GetConfigInt64(ctx, constants.PreferredAssetOutboundFeeMultiplier)
			threshold := ofRune.Mul(cosmos.NewUint(uint64(multiplier)))
			if err == nil && affcol.RuneAmount.GT(threshold) {
				if err = triggerPreferredAssetSwap(ctx, h.mgr, msg.AffiliateAddress, msg.Tx.ID, *affThorname, affcol, 2); err != nil {
					ctx.Logger().Error("fail to swap to preferred asset", "thorname", affThorname.Name, "err", err)
				}
			}
		}

		if affiliateSwap.Tx.Coins[0].Amount.GTE(amt) {
			affiliateSwap.Tx.Coins[0].Amount = amt
		}

		if err := h.mgr.Keeper().SetSwapQueueItem(ctx, *affiliateSwap, 1); err != nil {
			ctx.Logger().Error("fail to add swap to queue", "error", err)
		}
	}
}

func (h ObservedTxInHandler) isFromAsgard(ctx cosmos.Context, tx ObservedTx) (bool, error) {
	asgardVaults, err := h.mgr.Keeper().GetAsgardVaults(ctx)
	if err != nil {
		return false, err
	}
	return asgardVaults.HasAddress(tx.Tx.Chain, tx.Tx.FromAddress)
}

// ObservedTxInAnteHandler called by the ante handler to gate mempool entry
// and also during deliver. Store changes will persist if this function
// succeeds, regardless of the success of the transaction.
func ObservedTxInAnteHandler(ctx cosmos.Context, v semver.Version, k keeper.Keeper, msg MsgObservedTxIn) error {
	if !isSignedByActiveNodeAccounts(ctx, k, msg.GetSigners()) {
		return cosmos.ErrUnauthorized(fmt.Sprintf("%+v are not authorized", msg.GetSigners()))
	}
	return nil
}
