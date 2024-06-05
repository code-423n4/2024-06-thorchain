package mimir

import (
	"strconv"
	"strings"
)

//go:generate stringer --type=Id
type Id int32

const (
	Unknown Id = iota
	AffiliateFeeBasisPointsMax
	BondPause
	ConfMultiplierBasisPoints // https://gitlab.com/thorchain/thornode/-/issues/1599
	MaxConfirmations          // https://gitlab.com/thorchain/thornode/-/issues/1761
	CloutSwapperLimit
	CloutSwapperReset
	SwapSlipBasisPointsMin
	TradeAccountEnabled
)

// GetMimir fetches a mimir by id number
func GetMimir(id Id, ref string) (Mimir, bool) {
	switch id {
	case AffiliateFeeBasisPointsMax:
		return NewAffiliateFeeBasisPointsMax(ref), true
	case BondPause:
		return NewBondPause(ref), true
	case ConfMultiplierBasisPoints:
		return NewConfBasisPointValue(ref), true
	case MaxConfirmations:
		return NewMaxConfValue(ref), true
	case CloutSwapperLimit:
		return NewSwapperCloutLimit(ref), true
	case CloutSwapperReset:
		return NewSwapperCloutReset(ref), true
	case SwapSlipBasisPointsMin:
		return NewSwapSlipBasisPointsMin(ref), true
	case TradeAccountEnabled:
		return NewTradeAccountsEnabled(ref), true
	default:
		return nil, false
	}
}

// GetMimirByKey fetches a mimir by key
func GetMimirByKey(key string) (Mimir, bool) {
	idAndRef := strings.Split(key, "-")
	if len(idAndRef) != 2 {
		return nil, false
	}
	id, err := strconv.Atoi(idAndRef[0])
	if err != nil {
		return nil, false
	}
	return GetMimir(Id(id), idAndRef[1])
}
