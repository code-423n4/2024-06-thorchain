package thorchain

import (
	mem "gitlab.com/thorchain/thornode/x/thorchain/memo"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

const (
	YggdrasilVault    = types.VaultType_YggdrasilVault
	TxYggdrasilFund   = mem.TxYggdrasilFund
	TxYggdrasilReturn = mem.TxYggdrasilReturn
)

var (
	NewEventSwitch     = types.NewEventSwitch
	NewEventSwitchV87  = types.NewEventSwitchV87
	NewMsgSwitch       = types.NewMsgSwitch
	NewMsgYggdrasil    = types.NewMsgYggdrasil
	GetRandomYggVault  = types.GetRandomYggVault
	NewYggdrasilReturn = mem.NewYggdrasilReturn
	NewYggdrasilFund   = mem.NewYggdrasilFund
)

type (
	MsgSwitch           = types.MsgSwitch
	SwitchMemo          = mem.SwitchMemo
	MsgYggdrasil        = types.MsgYggdrasil
	YggdrasilFundMemo   = mem.YggdrasilFundMemo
	YggdrasilReturnMemo = mem.YggdrasilReturnMemo
)
