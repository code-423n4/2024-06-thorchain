package thorchain

import (
	"fmt"

	"github.com/blang/semver"

	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/mimir"
)

// TradeAccountDepositHandler is handler to process MsgTradeAccountDeposit
type TradeAccountDepositHandler struct {
	mgr Manager
}

// NewTradeAccountDepositHandler create a new instance of TradeAccountDepositHandler
func NewTradeAccountDepositHandler(mgr Manager) TradeAccountDepositHandler {
	return TradeAccountDepositHandler{
		mgr: mgr,
	}
}

// Run is the main entry point for TradeAccountDepositHandler
func (h TradeAccountDepositHandler) Run(ctx cosmos.Context, m cosmos.Msg) (*cosmos.Result, error) {
	msg, ok := m.(*MsgTradeAccountDeposit)
	if !ok {
		return nil, errInvalidMessage
	}
	if err := h.validate(ctx, *msg); err != nil {
		ctx.Logger().Error("MsgTradeAccountDeposit failed validation", "error", err)
		return nil, err
	}
	err := h.handle(ctx, *msg)
	if err != nil {
		ctx.Logger().Error("fail to process MsgTradeAccountDeposit", "error", err)
		return nil, err
	}
	return &cosmos.Result{}, err
}

func (h TradeAccountDepositHandler) validate(ctx cosmos.Context, msg MsgTradeAccountDeposit) error {
	version := h.mgr.GetVersion()
	if version.GTE(semver.MustParse("0.1.0")) {
		return h.validateV1(ctx, msg)
	}
	return errBadVersion
}

func (h TradeAccountDepositHandler) validateV1(ctx cosmos.Context, msg MsgTradeAccountDeposit) error {
	if mimir.NewTradeAccountsEnabled().IsOff(ctx, h.mgr.Keeper()) {
		return fmt.Errorf("trade account is disabled")
	}
	return msg.ValidateBasic()
}

func (h TradeAccountDepositHandler) handle(ctx cosmos.Context, msg MsgTradeAccountDeposit) error {
	version := h.mgr.GetVersion()
	if version.GTE(semver.MustParse("0.1.0")) {
		return h.handleV1(ctx, msg)
	}
	return errBadVersion
}

// handle process MsgTradeAccountDeposit
func (h TradeAccountDepositHandler) handleV1(ctx cosmos.Context, msg MsgTradeAccountDeposit) error {
	_, err := h.mgr.TradeAccountManager().Deposit(ctx, msg.Asset, msg.Amount, msg.Address, msg.Tx.FromAddress, msg.Tx.ID)
	if err != nil {
		ctx.Logger().Error("fail to handle Deposit", "error", err)
		return err
	}
	return nil
}
