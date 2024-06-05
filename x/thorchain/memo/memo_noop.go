package thorchain

import (
	"fmt"

	"github.com/blang/semver"
)

type NoOpMemo struct {
	MemoBase
	Action string
}

// String implement fmt.Stringer
func (m NoOpMemo) String() string {
	if len(m.Action) == 0 {
		return "noop"
	}
	return fmt.Sprintf("noop:%s", m.Action)
}

// NewNoOpMemo create a new instance of NoOpMemo
func NewNoOpMemo(action string) NoOpMemo {
	return NoOpMemo{
		MemoBase: MemoBase{TxType: TxNoOp},
		Action:   action,
	}
}

func (p *parser) ParseNoOpMemo() (NoOpMemo, error) {
	switch {
	case p.version.GTE(semver.MustParse("1.116.0")):
		return p.ParseNoOpMemoV116()
	default:
		return ParseNoOpMemoV1(p.parts)
	}
}

// ParseNoOpMemo try to parse the memo
func (p *parser) ParseNoOpMemoV116() (NoOpMemo, error) {
	return NewNoOpMemo(p.get(1)), p.Error()
}
