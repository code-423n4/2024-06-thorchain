package thorchain

import (
	"fmt"

	"gitlab.com/thorchain/thornode/common/cosmos"

	"github.com/blang/semver"
)

type BondMemo struct {
	MemoBase
	NodeAddress         cosmos.AccAddress
	BondProviderAddress cosmos.AccAddress
	NodeOperatorFee     int64
}

func (m BondMemo) GetAccAddress() cosmos.AccAddress { return m.NodeAddress }

func NewBondMemo(addr, additional cosmos.AccAddress, operatorFee int64) BondMemo {
	return BondMemo{
		MemoBase:            MemoBase{TxType: TxBond},
		NodeAddress:         addr,
		BondProviderAddress: additional,
		NodeOperatorFee:     operatorFee,
	}
}

func (p *parser) ParseBondMemo() (BondMemo, error) {
	switch {
	case p.version.GTE(semver.MustParse("1.116.0")):
		return p.ParseBondMemoV116()
	case p.version.GTE(semver.MustParse("1.88.0")):
		return ParseBondMemoV88(p.parts)
	case p.version.GTE(semver.MustParse("0.81.0")):
		return ParseBondMemoV81(p.parts)
	default:
		return BondMemo{}, fmt.Errorf("invalid version(%s)", p.version.String())
	}
}

func (p *parser) ParseBondMemoV116() (BondMemo, error) {
	addr := p.getAccAddress(1, true, nil)
	additional := p.getAccAddress(2, false, nil)
	operatorFee := p.getInt64(3, false, -1)
	return NewBondMemo(addr, additional, operatorFee), p.Error()
}
