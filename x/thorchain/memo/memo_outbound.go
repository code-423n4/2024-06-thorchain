package thorchain

import (
	"fmt"

	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common"
)

type OutboundMemo struct {
	MemoBase
	TxID common.TxID
}

func (m OutboundMemo) GetTxID() common.TxID { return m.TxID }
func (m OutboundMemo) String() string {
	return fmt.Sprintf("OUT:%s", m.TxID.String())
}

func NewOutboundMemo(txID common.TxID) OutboundMemo {
	return OutboundMemo{
		MemoBase: MemoBase{TxType: TxOutbound},
		TxID:     txID,
	}
}

func (p *parser) ParseOutboundMemo() (OutboundMemo, error) {
	switch {
	case p.version.GTE(semver.MustParse("1.116.0")):
		return p.ParseOutboundMemoV116()
	default:
		return ParseOutboundMemoV1(p.parts)
	}
}

func (p *parser) ParseOutboundMemoV116() (OutboundMemo, error) {
	txID := p.getTxID(1, true, common.BlankTxID)
	return NewOutboundMemo(txID), p.Error()
}
