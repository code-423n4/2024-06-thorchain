package types

import (
	"fmt"

	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

type StreamingSwaps []StreamingSwap

func NewStreamingSwap(hash common.TxID, quan, interval uint64, target, deposit cosmos.Uint) StreamingSwap {
	return StreamingSwap{
		TxID:        hash,
		Quantity:    quan,
		Interval:    interval,
		TradeTarget: target,
		Deposit:     deposit,
		In:          cosmos.ZeroUint(),
		Out:         cosmos.ZeroUint(),
	}
}

func (m *StreamingSwap) Valid() error {
	if m.Quantity < 1 {
		return fmt.Errorf("quantity cannot be less than 1")
	}
	if m.Interval < 1 {
		return fmt.Errorf("interval cannot be less than 1")
	}
	if m.Deposit.IsZero() {
		return fmt.Errorf("deposit amount cannot be zero")
	}
	return nil
}

func (m *StreamingSwap) NextSize(version semver.Version) (cosmos.Uint, cosmos.Uint) {
	switch {
	case version.GTE(semver.MustParse("1.121.0")):
		return m.NextSizeV121()
	case version.GTE(semver.MustParse("1.116.0")):
		return m.NextSizeV116()
	default:
		return m.NextSizeV115()
	}
}

func (m *StreamingSwap) NextSizeV121() (cosmos.Uint, cosmos.Uint) {
	swapSize := m.DefaultSwapSize()

	// sanity check, ensure we never exceed the deposit amount
	// Also, if this is the last swap, just do the remainder
	if m.Deposit.LT(m.In.Add(swapSize)) || m.IsLastSwap() {
		// use remainder of `m.Depost - m.In` instead
		swapSize = common.SafeSub(m.Deposit, m.In)
	}

	// calculate trade target for this sub-swap
	remainingIn := common.SafeSub(m.Deposit, m.In)       // remaining inbound
	remainingOut := common.SafeSub(m.TradeTarget, m.Out) // remaining outbound
	target := common.GetSafeShare(swapSize, remainingIn, remainingOut)

	return swapSize, target
}

func (m *StreamingSwap) DefaultSwapSize() cosmos.Uint {
	if m.Quantity == 0 {
		return cosmos.ZeroUint()
	}
	return m.Deposit.QuoUint64(m.Quantity)
}

func (m *StreamingSwap) IsDone() bool {
	return m.Count >= m.Quantity
}

func (m *StreamingSwap) IsLastSwap() bool {
	return m.Count+1 >= m.Quantity
}
