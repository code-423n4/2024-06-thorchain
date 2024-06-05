package types

import (
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

func (m *StreamingSwap) NextSizeV116() (cosmos.Uint, cosmos.Uint) {
	swapSize := m.DefaultSwapSize()

	// sanity check, ensure we never exceed the deposit amount
	// Also, if this is the last swap, just do the remainder
	if m.Deposit.LT(m.In.Add(swapSize)) || m.Count+1 >= m.Quantity {
		// use remainder of `m.Depost - m.In` instead
		swapSize = common.SafeSub(m.Deposit, m.In)
	}

	// calculate trade target for this sub-swap
	remainingIn := common.SafeSub(m.Deposit, m.In)       // remaining inbound
	remainingOut := common.SafeSub(m.TradeTarget, m.Out) // remaining outbound
	target := common.GetSafeShare(swapSize, remainingIn, remainingOut)

	return swapSize, target
}

func (m *StreamingSwap) NextSizeV115() (cosmos.Uint, cosmos.Uint) {
	swapSize := m.DefaultSwapSize()

	// sanity check, ensure we never exceed the deposit amount
	if m.Deposit.LT(m.In.Add(swapSize)) {
		// use remainder of `m.Depost - m.In` instead
		swapSize = common.SafeSub(m.Deposit, m.In)
	}

	// calculate trade target for this sub-swap
	remainingIn := common.SafeSub(m.Deposit, m.In)       // remaining inbound
	remainingOut := common.SafeSub(m.TradeTarget, m.Out) // remaining outbound
	target := common.GetSafeShare(swapSize, remainingIn, remainingOut)

	return swapSize, target
}
