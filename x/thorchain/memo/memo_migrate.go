package thorchain

import (
	"fmt"

	"github.com/blang/semver"
)

type MigrateMemo struct {
	MemoBase
	BlockHeight int64
}

func (m MigrateMemo) String() string {
	return fmt.Sprintf("MIGRATE:%d", m.BlockHeight)
}

func (m MigrateMemo) GetBlockHeight() int64 {
	return m.BlockHeight
}

func NewMigrateMemo(blockHeight int64) MigrateMemo {
	return MigrateMemo{
		MemoBase:    MemoBase{TxType: TxMigrate},
		BlockHeight: blockHeight,
	}
}

func (p *parser) ParseMigrateMemo() (MigrateMemo, error) {
	switch {
	case p.version.GTE(semver.MustParse("1.116.0")):
		return p.ParseMigrateMemoV116()
	default:
		return ParseMigrateMemoV1(p.parts)
	}
}

func (p *parser) ParseMigrateMemoV116() (memo MigrateMemo, err error) {
	blockHeight := p.getInt64(1, true, 0)
	return NewMigrateMemo(blockHeight), p.Error()
}
