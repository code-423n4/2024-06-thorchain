package thorchain

import (
	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

type LeaveMemo struct {
	MemoBase
	NodeAddress cosmos.AccAddress
}

func (m LeaveMemo) GetAccAddress() cosmos.AccAddress { return m.NodeAddress }

func NewLeaveMemo(addr cosmos.AccAddress) LeaveMemo {
	return LeaveMemo{
		MemoBase:    MemoBase{TxType: TxLeave},
		NodeAddress: addr,
	}
}

func (p *parser) ParseLeaveMemo() (LeaveMemo, error) {
	switch {
	case p.version.GTE(semver.MustParse("1.116.0")):
		return p.ParseLeaveMemoV116()
	default:
		return ParseLeaveMemoV1(p.parts)
	}
}

func (p *parser) ParseLeaveMemoV116() (LeaveMemo, error) {
	addr := p.getAccAddress(1, true, nil)
	return NewLeaveMemo(addr), p.Error()
}
