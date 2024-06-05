package thorchain

import (
	"errors"
	"fmt"

	"github.com/blang/semver"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
)

// NewSendHandler create a new instance of SendHandler
func NewSendHandler(mgr Manager) BaseHandler[*MsgSend] {
	return BaseHandler[*MsgSend]{
		mgr:    mgr,
		logger: MsgSendLogger,
		validators: NewValidators[*MsgSend]().
			Register("1.130.0", MsgSendValidateV130).
			Register("1.121.0", MsgSendValidateV121).
			Register("1.87.0", MsgSendValidateV87).
			Register("0.1.0", MsgSendValidateV1),
		handlers: NewHandlers[*MsgSend]().
			Register("1.116.0", MsgSendHandleV116).
			Register("1.115.0", MsgSendHandleV115).
			Register("1.112.0", MsgSendHandleV112).
			Register("1.108.0", MsgSendHandleV108).
			Register("0.1.0", MsgSendHandleV1),
	}
}

func MsgSendValidateV130(ctx cosmos.Context, mgr Manager, msg *MsgSend) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	// disallow sends to modules, they should only be interacted with via deposit messages
	if IsModuleAccAddress(mgr.Keeper(), msg.ToAddress) {
		return fmt.Errorf("cannot use MsgSend for Module transactions, use MsgDeposit instead")
	}

	// TODO on hard fork move to ValidateBasic
	// send only allowed with one coin
	if len(msg.Amount) != 1 {
		return errors.New("only one coin is allowed")
	}

	return nil
}

func MsgSendLogger(ctx cosmos.Context, msg *MsgSend) {
	ctx.Logger().Info("receive MsgSend", "from", msg.FromAddress, "to", msg.ToAddress, "coins", msg.Amount)
}

func MsgSendHandleV116(ctx cosmos.Context, mgr Manager, msg *MsgSend) (*cosmos.Result, error) {
	if mgr.Keeper().IsChainHalted(ctx, common.THORChain) {
		return nil, fmt.Errorf("unable to use MsgSend while THORChain is halted")
	}

	err := mgr.Keeper().SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &cosmos.Result{}, nil
}

// SendAnteHandler called by the ante handler to gate mempool entry
// and also during deliver. Store changes will persist if this function
// succeeds, regardless of the success of the transaction.
func SendAnteHandler(ctx cosmos.Context, v semver.Version, k keeper.Keeper, msg MsgSend) error {
	// TODO remove on hard fork
	if v.LT(semver.MustParse("1.115.0")) {
		nativeTxFee := k.GetNativeTxFee(ctx)
		gas := common.NewCoin(common.RuneNative, nativeTxFee)
		gasFee, err := gas.Native()
		if err != nil {
			return fmt.Errorf("fail to get gas fee: %w", err)
		}
		totalCoins := cosmos.NewCoins(gasFee)
		if !k.HasCoins(ctx, msg.GetSigners()[0], totalCoins) {
			return cosmos.ErrInsufficientCoins(err, "insufficient funds")
		}
		return nil
	}

	return k.DeductNativeTxFeeFromAccount(ctx, msg.GetSigners()[0])
}
