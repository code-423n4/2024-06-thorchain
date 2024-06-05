package thorchain

import (
	"fmt"

	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common"
	cosmos "gitlab.com/thorchain/thornode/common/cosmos"
)

type TradeAccountDepositMemo struct {
	MemoBase
	Address cosmos.AccAddress
}

func (m TradeAccountDepositMemo) GetAccAddress() cosmos.AccAddress { return m.Address }

func NewTradeAccountDepositMemo(addr cosmos.AccAddress) TradeAccountDepositMemo {
	return TradeAccountDepositMemo{
		MemoBase: MemoBase{TxType: TxTradeAccountDeposit},
		Address:  addr,
	}
}

func (p *parser) ParseTradeAccountDeposit() (TradeAccountDepositMemo, error) {
	switch {
	case p.version.GTE(semver.MustParse("1.128.0")):
		return p.ParseTradeAccountDepositV128()
	default:
		return TradeAccountDepositMemo{}, fmt.Errorf("version %s is not supported", p.version)
	}
}

func (p *parser) ParseTradeAccountDepositV128() (TradeAccountDepositMemo, error) {
	addr := p.getAccAddress(1, true, nil)
	return NewTradeAccountDepositMemo(addr), p.Error()
}

type TradeAccountWithdrawalMemo struct {
	MemoBase
	Address common.Address
	Amount  cosmos.Uint
}

func (m TradeAccountWithdrawalMemo) GetAddress() common.Address { return m.Address }
func (m TradeAccountWithdrawalMemo) GetAmount() cosmos.Uint     { return m.Amount }

func NewTradeAccountWithdrawalMemo(addr common.Address) TradeAccountWithdrawalMemo {
	return TradeAccountWithdrawalMemo{
		MemoBase: MemoBase{TxType: TxTradeAccountWithdrawal},
		Address:  addr,
	}
}

func (p *parser) ParseTradeAccountWithdrawal() (TradeAccountWithdrawalMemo, error) {
	switch {
	case p.version.GTE(semver.MustParse("1.128.0")):
		return p.ParseTradeAccountWithdrawalV128()
	default:
		return TradeAccountWithdrawalMemo{}, fmt.Errorf("version %s is not supported", p.version)
	}
}

func (p *parser) ParseTradeAccountWithdrawalV128() (TradeAccountWithdrawalMemo, error) {
	addr := p.getAddress(1, true, common.NoAddress)
	return NewTradeAccountWithdrawalMemo(addr), p.Error()
}
