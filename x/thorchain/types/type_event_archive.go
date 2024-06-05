package types

import (
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

const (
	SwitchEventType = "switch"
)

// NewEventSwitch create a new instance of EventSwitch
func NewEventSwitch(from common.Address, to cosmos.AccAddress, coin common.Coin, hash common.TxID) *EventSwitch {
	return &EventSwitch{
		TxID:        hash,
		ToAddress:   to,
		FromAddress: from,
		Burn:        coin,
	}
}

// Type return a string which represent the type of this event
func (m *EventSwitch) Type() string {
	return SwitchEventType
}

// Events return cosmos sdk events
func (m *EventSwitch) Events() (cosmos.Events, error) {
	evt := cosmos.NewEvent(m.Type(),
		cosmos.NewAttribute("txid", m.TxID.String()),
		cosmos.NewAttribute("from", m.FromAddress.String()),
		cosmos.NewAttribute("to", m.ToAddress.String()),
		cosmos.NewAttribute("burn", m.Burn.String()))
	return cosmos.Events{evt}, nil
}

// NewEventSwitchV87 create a new instance of EventSwitch
func NewEventSwitchV87(from common.Address, to cosmos.AccAddress, coin common.Coin, hash common.TxID, mint cosmos.Uint) *EventSwitchV87 {
	return &EventSwitchV87{
		TxID:        hash,
		ToAddress:   to,
		FromAddress: from,
		Burn:        coin,
		Mint:        mint,
	}
}

// Type return a string which represent the type of this event
func (m *EventSwitchV87) Type() string {
	return SwitchEventType
}

// Events return cosmos sdk events
func (m *EventSwitchV87) Events() (cosmos.Events, error) {
	evt := cosmos.NewEvent(m.Type(),
		cosmos.NewAttribute("txid", m.TxID.String()),
		cosmos.NewAttribute("from", m.FromAddress.String()),
		cosmos.NewAttribute("to", m.ToAddress.String()),
		cosmos.NewAttribute("burn", m.Burn.String()),
		cosmos.NewAttribute("mint", m.Mint.String()))
	return cosmos.Events{evt}, nil
}
