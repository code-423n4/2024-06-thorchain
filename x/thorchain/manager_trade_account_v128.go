package thorchain

import (
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
)

// TradeMgrV128 is V128 implementation of slasher
type TradeMgrV128 struct {
	keeper   keeper.Keeper
	eventMgr EventManager
}

// newTradeMgrV128 create a new instance of Slasher
func newTradeMgrV128(keeper keeper.Keeper, eventMgr EventManager) *TradeMgrV128 {
	return &TradeMgrV128{
		keeper:   keeper,
		eventMgr: eventMgr,
	}
}

func (s *TradeMgrV128) EndBlock(ctx cosmos.Context, keeper keeper.Keeper) error {
	// TODO: implement liquidation
	return nil
}

func (s *TradeMgrV128) BalanceOf(ctx cosmos.Context, asset common.Asset, addr cosmos.AccAddress) cosmos.Uint {
	asset = asset.GetTradeAsset()
	tu, err := s.keeper.GetTradeUnit(ctx, asset)
	if err != nil {
		return cosmos.ZeroUint()
	}

	tr, err := s.keeper.GetTradeAccount(ctx, addr, asset)
	if err != nil {
		return cosmos.ZeroUint()
	}

	return common.GetSafeShare(tu.Units, tu.Depth, tr.Units)
}

func (s *TradeMgrV128) Deposit(ctx cosmos.Context, asset common.Asset, amount cosmos.Uint, owner cosmos.AccAddress, assetAddr common.Address, txID common.TxID) (cosmos.Uint, error) {
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

func (s *TradeMgrV128) calcDepositUnits(oldUnits, depth, add cosmos.Uint) cosmos.Uint {
	if oldUnits.IsZero() || depth.IsZero() {
		return add
	}
	if add.IsZero() {
		return cosmos.ZeroUint()
	}
	return common.GetUncappedShare(add, depth, oldUnits)
}

func (s *TradeMgrV128) Withdrawal(ctx cosmos.Context, asset common.Asset, amount cosmos.Uint, owner cosmos.AccAddress, assetAddr common.Address, txID common.TxID) (cosmos.Uint, error) {
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

	assetAvailable := common.GetSafeShare(tu.Units, tu.Depth, tr.Units)
	unitsToClaim := common.GetSafeShare(amount, assetAvailable, tr.Units)

	tokensToClaim := common.GetSafeShare(unitsToClaim, tu.Units, tu.Depth)
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
