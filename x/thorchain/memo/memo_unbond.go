package thorchain

import (
	"fmt"

	"github.com/blang/semver"

	"gitlab.com/thorchain/thornode/common/cosmos"
)

type UnbondMemo struct {
	MemoBase
	NodeAddress         cosmos.AccAddress
	Amount              cosmos.Uint
	BondProviderAddress cosmos.AccAddress
}

func (m UnbondMemo) GetAccAddress() cosmos.AccAddress { return m.NodeAddress }
func (m UnbondMemo) GetAmount() cosmos.Uint           { return m.Amount }

func NewUnbondMemo(addr, additional cosmos.AccAddress, amt cosmos.Uint) UnbondMemo {
	return UnbondMemo{
		MemoBase:            MemoBase{TxType: TxUnbond},
		NodeAddress:         addr,
		Amount:              amt,
		BondProviderAddress: additional,
	}
}

func (p *parser) ParseUnbondMemo() (UnbondMemo, error) {
	switch {
	case p.version.GTE(semver.MustParse("1.116.0")):
		return p.ParseUnbondMemoV116()
	case p.version.GTE(semver.MustParse("0.81.0")):
		return ParseUnbondMemoV81(p.parts)
	default:
		return UnbondMemo{}, fmt.Errorf("invalid version(%s)", p.version.String())
	}
}

func (p *parser) ParseUnbondMemoV116() (UnbondMemo, error) {
	addr := p.getAccAddress(1, true, nil)
	amt := p.getUint(2, true, 0)
	additional := p.getAccAddress(3, false, nil)
	return NewUnbondMemo(addr, additional, amt), p.Error()
}
