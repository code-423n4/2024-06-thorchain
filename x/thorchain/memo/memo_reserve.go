package thorchain

type ReserveMemo struct {
	MemoBase
}

func (p *parser) ParseReserveMemo() (ReserveMemo, error) {
	return ReserveMemo{
		MemoBase: MemoBase{TxType: TxReserve},
	}, nil
}
