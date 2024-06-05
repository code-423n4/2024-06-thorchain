package thorchain

type ConsolidateMemo struct {
	MemoBase
}

// String implement fmt.Stringer
func (m ConsolidateMemo) String() string {
	return "consolidate"
}

// NewConsolidateMemo create a new instance of ConsolidateMemo
func NewConsolidateMemo() ConsolidateMemo {
	return ConsolidateMemo{
		MemoBase: MemoBase{TxType: TxConsolidate},
	}
}

// ParseConsolidateMemo try to parse the memo
func (p *parser) ParseConsolidateMemo() (ConsolidateMemo, error) {
	return NewConsolidateMemo(), nil
}
