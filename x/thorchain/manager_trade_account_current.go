package thorchain

import (
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
)

// TradeMgrVCUR is VCUR implementation of slasher
type TradeMgrVCUR struct {
	keeper   keeper.Keeper
	eventMgr EventManager
}

// newTradeMgrVCUR create a new instance of Slasher
func newTradeMgrVCUR(keeper keeper.Keeper, eventMgr EventManager) *TradeMgrVCUR {
	return &TradeMgrVCUR{
		keeper:   keeper,
		eventMgr: eventMgr,
	}
}

func (s *TradeMgrVCUR) EndBlock(ctx cosmos.Context, keeper keeper.Keeper) error {
	// TODO: implement liquidation
	return nil
}

func (s *TradeMgrVCUR) BalanceOf(ctx cosmos.Context, asset common.Asset, addr cosmos.AccAddress) cosmos.Uint {
	asset = asset.GetTradeAsset()
	tu, err := s.keeper.GetTradeUnit(ctx, asset)
	if err != nil {
		return cosmos.ZeroUint()
	}

	tr, err := s.keeper.GetTradeAccount(ctx, addr, asset)
	if err != nil {
		return cosmos.ZeroUint()
	}

	// Proportion of total Depth that the account's Units entitle it to:
	return common.GetSafeShare(tr.Units, tu.Units, tu.Depth)
}

func (s *TradeMgrVCUR) Deposit(ctx cosmos.Context, asset common.Asset, amount cosmos.Uint, owner cosmos.AccAddress, assetAddr common.Address, txID common.TxID) (cosmos.Uint, error) {
	asset = asset.GetTradeAsset()
	tu, err := s.keeper.GetTradeUnit(ctx, asset)
	if err != nil {
		return cosmos.ZeroUint(), err
	}

	tr, err := s.keeper.GetTradeAccount(ctx, owner, asset)
	if err != nil {
		return cosmos.ZeroUint(), err
	}
	tr.LastAddHeight = ctx.BlockHeight()

	units := s.calcDepositUnits(tu.Units, tu.Depth, amount)
	tu.Units = tu.Units.Add(units)
	tr.Units = tr.Units.Add(units)
	tu.Depth = tu.Depth.Add(amount)

	s.keeper.SetTradeUnit(ctx, tu)
	s.keeper.SetTradeAccount(ctx, tr)

	depositEvent := NewEventTradeAccountDeposit(amount, asset, assetAddr, common.Address(owner.String()), txID)
	if err := s.eventMgr.EmitEvent(ctx, depositEvent); err != nil {
		ctx.Logger().Error("fail to emit trade account deposit event", "error", err)
	}

	return amount, nil
}

func (s *TradeMgrVCUR) calcDepositUnits(oldUnits, depth, add cosmos.Uint) cosmos.Uint {
	if oldUnits.IsZero() || depth.IsZero() {
		return add
	}
	if add.IsZero() {
		return cosmos.ZeroUint()
	}
	return common.GetUncappedShare(add, depth, oldUnits)
}

func (s *TradeMgrVCUR) Withdrawal(ctx cosmos.Context, asset common.Asset, amount cosmos.Uint, owner cosmos.AccAddress, assetAddr common.Address, txID common.TxID) (cosmos.Uint, error) {
	asset = asset.GetTradeAsset()
	tu, err := s.keeper.GetTradeUnit(ctx, asset)
	if err != nil {
		return cosmos.ZeroUint(), err
	}

	tr, err := s.keeper.GetTradeAccount(ctx, owner, asset)
	if err != nil {
		return cosmos.ZeroUint(), err
	}
	tr.LastWithdrawHeight = ctx.BlockHeight()

	// assetAvailable is the same as BalanceOf:
	// Proportion of total Depth that the account's Units entitle it to:
	assetAvailable := common.GetSafeShare(tr.Units, tu.Units, tu.Depth)

	// unitsToClaim is the account's units for the specified amount to be withdrawn from assetAvailable,
	// capped at the accounts total Units.
	unitsToClaim := common.GetSafeShare(amount, assetAvailable, tr.Units)

	// tokensToClaim is the exact amount to be withdrawn from those unitsToClaim,
	// capped at the account's assetAvailable.
	tokensToClaim := common.GetSafeShare(unitsToClaim, tr.Units, assetAvailable)

	tu.Units = common.SafeSub(tu.Units, unitsToClaim)
	tr.Units = common.SafeSub(tr.Units, unitsToClaim)
	tu.Depth = common.SafeSub(tu.Depth, tokensToClaim)

	s.keeper.SetTradeUnit(ctx, tu)
	s.keeper.SetTradeAccount(ctx, tr)

	withdrawEvent := NewEventTradeAccountWithdraw(tokensToClaim, asset, assetAddr, common.Address(owner.String()), txID)
	if err := s.eventMgr.EmitEvent(ctx, withdrawEvent); err != nil {
		ctx.Logger().Error("fail to emit trade account withdraw event", "error", err)
	}

	return tokensToClaim, nil
}
