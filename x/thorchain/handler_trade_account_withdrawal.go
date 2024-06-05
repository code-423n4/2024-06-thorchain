package thorchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/blang/semver"
	"github.com/hashicorp/go-multierror"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/mimir"
)

// TradeAccountWithdrawalHandler is handler to process MsgTradeAccountWithdrawal
type TradeAccountWithdrawalHandler struct {
	mgr Manager
}

// NewTradeAccountWithdrawalHandler create a new instance of TradeAccountWithdrawalHandler
func NewTradeAccountWithdrawalHandler(mgr Manager) TradeAccountWithdrawalHandler {
	return TradeAccountWithdrawalHandler{
		mgr: mgr,
	}
}

// Run is the main entry point for TradeAccountWithdrawalHandler
func (h TradeAccountWithdrawalHandler) Run(ctx cosmos.Context, m cosmos.Msg) (*cosmos.Result, error) {
	msg, ok := m.(*MsgTradeAccountWithdrawal)
	if !ok {
		return nil, errInvalidMessage
	}
	if err := h.validate(ctx, *msg); err != nil {
		ctx.Logger().Error("MsgTradeAccountWithdrawal failed validation", "error", err)
		return nil, err
	}
	err := h.handle(ctx, *msg)
	if err != nil {
		ctx.Logger().Error("fail to process MsgTradeAccountWithdrawal", "error", err)
	}
	return &cosmos.Result{}, err
}

func (h TradeAccountWithdrawalHandler) validate(ctx cosmos.Context, msg MsgTradeAccountWithdrawal) error {
	version := h.mgr.GetVersion()
	if version.GTE(semver.MustParse("0.1.0")) {
		return h.validateV1(ctx, msg)
	}
	return errBadVersion
}

func (h TradeAccountWithdrawalHandler) validateV1(ctx cosmos.Context, msg MsgTradeAccountWithdrawal) error {
	if mimir.NewTradeAccountsEnabled().IsOff(ctx, h.mgr.Keeper()) {
		return fmt.Errorf("trade account is disabled")
	}
	return msg.ValidateBasic()
}

func (h TradeAccountWithdrawalHandler) handle(ctx cosmos.Context, msg MsgTradeAccountWithdrawal) error {
	version := h.mgr.GetVersion()
	if version.GTE(semver.MustParse("0.1.0")) {
		return h.handleV1(ctx, msg)
	}
	return errBadVersion
}

// handle process MsgTradeAccountWithdrawal
func (h TradeAccountWithdrawalHandler) handleV1(ctx cosmos.Context, msg MsgTradeAccountWithdrawal) error {
	withdraw, err := h.mgr.TradeAccountManager().Withdrawal(ctx, msg.Asset, msg.Amount, msg.Signer, msg.AssetAddress, msg.Tx.ID)
	if err != nil {
		return err
	}

	var ok bool
	layer1Asset := msg.Asset.GetLayer1Asset()

	rawHash := sha256.Sum256(ctx.TxBytes())
	hash := hex.EncodeToString(rawHash[:])
	txID, err := common.NewTxID(hash)
	if err != nil {
		return err
	}
	toi := TxOutItem{
		Chain:     layer1Asset.GetChain(),
		InHash:    txID,
		ToAddress: msg.AssetAddress,
		Coin:      common.NewCoin(layer1Asset, withdraw),
	}

	ok, err = h.mgr.TxOutStore().TryAddTxOutItem(ctx, h.mgr, toi, cosmos.ZeroUint())
	if err != nil {
		return multierror.Append(errFailAddOutboundTx, err)
	}
	if !ok {
		return errFailAddOutboundTx
	}

	return nil
}
