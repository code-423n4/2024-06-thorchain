package types

import "gitlab.com/thorchain/thornode/common/cosmos"

// ValidateBasic runs stateless checks on the message
func (m *MsgUnBond) ValidateBasic() error {
	if m.NodeAddress.Empty() {
		return cosmos.ErrInvalidAddress("node address cannot be empty")
	}
	if m.Amount.IsZero() {
		return cosmos.ErrUnknownRequest("unbond amount cannot be zero")
	}
	if m.BondAddress.IsEmpty() {
		return cosmos.ErrInvalidAddress("bond address cannot be empty")
	}
	// here we can't call m.TxIn.Valid , because we allow user to send unbond request without any coins in it
	// m.TxIn.Valid will reject this kind request , which result unbond to fail
	if m.TxIn.ID.IsEmpty() {
		return cosmos.ErrUnknownRequest("tx id cannot be empty")
	}
	if m.TxIn.FromAddress.IsEmpty() {
		return cosmos.ErrInvalidAddress("tx from address cannot be empty")
	}
	if m.Signer.Empty() {
		return cosmos.ErrInvalidAddress("empty signer address")
	}
	return nil
}
