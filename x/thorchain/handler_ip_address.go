package thorchain

import (
	"fmt"

	"github.com/blang/semver"

	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
)

// IPAddressHandler is to handle ip address message
type IPAddressHandler struct {
	mgr Manager
}

// NewIPAddressHandler create new instance of IPAddressHandler
func NewIPAddressHandler(mgr Manager) IPAddressHandler {
	return IPAddressHandler{
		mgr: mgr,
	}
}

// Run it the main entry point to execute ip address logic
func (h IPAddressHandler) Run(ctx cosmos.Context, m cosmos.Msg) (*cosmos.Result, error) {
	msg, ok := m.(*MsgSetIPAddress)
	if !ok {
		return nil, errInvalidMessage
	}
	ctx.Logger().Info("receive ip address", "address", msg.IPAddress)
	if err := h.validate(ctx, *msg); err != nil {
		ctx.Logger().Error("msg set version failed validation", "error", err)
		return nil, err
	}
	if err := h.handle(ctx, *msg); err != nil {
		ctx.Logger().Error("fail to process msg set version", "error", err)
		return nil, err
	}

	return &cosmos.Result{}, nil
}

func (h IPAddressHandler) validate(ctx cosmos.Context, msg MsgSetIPAddress) error {
	version := h.mgr.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.114.0")):
		return h.validateV114(ctx, msg)
	case version.GTE(semver.MustParse("1.112.0")):
		return h.validateV112(ctx, msg)
	case version.GTE(semver.MustParse("0.1.0")):
		return h.validateV1(ctx, msg)
	}
	return errBadVersion
}

func (h IPAddressHandler) validateV114(ctx cosmos.Context, msg MsgSetIPAddress) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	if err := validateIPAddressAuth(ctx, h.mgr.Keeper(), msg.Signer); err != nil {
		return err
	}
	return nil
}

func (h IPAddressHandler) handle(ctx cosmos.Context, msg MsgSetIPAddress) error {
	ctx.Logger().Info("handleMsgSetIPAddress request", "ip address", msg.IPAddress)
	version := h.mgr.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.115.0")):
		return h.handleV115(ctx, msg)
	case version.GTE(semver.MustParse("1.112.0")):
		return h.handleV112(ctx, msg)
	case version.GTE(semver.MustParse("0.57.0")):
		return h.handleV57(ctx, msg)
	}
	ctx.Logger().Error(errInvalidVersion.Error())
	return errBadVersion
}

func (h IPAddressHandler) handleV115(ctx cosmos.Context, msg MsgSetIPAddress) error {
	nodeAccount, err := h.mgr.Keeper().GetNodeAccount(ctx, msg.Signer)
	if err != nil {
		ctx.Logger().Error("fail to get node account", "error", err, "address", msg.Signer.String())
		return cosmos.ErrUnauthorized(fmt.Sprintf("unable to find account: %s", msg.Signer))
	}

	nodeAccount.IPAddress = msg.IPAddress
	if err := h.mgr.Keeper().SetNodeAccount(ctx, nodeAccount); err != nil {
		return fmt.Errorf("fail to save node account: %w", err)
	}

	ctx.EventManager().EmitEvent(
		cosmos.NewEvent("set_ip_address",
			cosmos.NewAttribute("thor_address", msg.Signer.String()),
			cosmos.NewAttribute("address", msg.IPAddress)))

	return nil
}

func validateIPAddressAuth(ctx cosmos.Context, k keeper.Keeper, signer cosmos.AccAddress) error {
	version, _ := k.GetVersionWithCtx(ctx)
	switch {
	case version.GTE(semver.MustParse("1.115.0")):
		return validateIPAddressAuthV115(ctx, k, signer)
	case version.GTE(semver.MustParse("1.114.0")):
		return validateIPAddressAuthV114(ctx, k, signer)
	default:
		return errBadVersion
	}
}

func validateIPAddressAuthV115(ctx cosmos.Context, k keeper.Keeper, signer cosmos.AccAddress) error {
	nodeAccount, err := k.GetNodeAccount(ctx, signer)
	if err != nil {
		ctx.Logger().Error("fail to get node account", "error", err, "address", signer.String())
		return cosmos.ErrUnauthorized(fmt.Sprintf("%s is not authorized", signer))
	}
	if nodeAccount.IsEmpty() {
		ctx.Logger().Error("unauthorized account", "address", signer.String())

		return cosmos.ErrUnauthorized(fmt.Sprintf("%s is not authorized", signer))
	}
	if nodeAccount.Type != NodeTypeValidator {
		ctx.Logger().Error("unauthorized account, node account must be a validator", "address", signer.String(), "type", nodeAccount.Type)
		return cosmos.ErrUnauthorized(fmt.Sprintf("%s is not authorized", signer))
	}
	return nil
}

// IPAddressAnteHandler called by the ante handler to gate mempool entry
// and also during deliver. Store changes will persist if this function
// succeeds, regardless of the success of the transaction.
func IPAddressAnteHandler(ctx cosmos.Context, v semver.Version, k keeper.Keeper, msg MsgSetIPAddress) error {
	if err := validateIPAddressAuth(ctx, k, msg.Signer); err != nil {
		return err
	}
	// TODO on hard fork remove version check
	if v.GTE(semver.MustParse("1.115.0")) {
		return k.DeductNativeTxFeeFromBond(ctx, msg.Signer)
	}
	return nil
}
