package thorchain

import (
	"fmt"

	"gitlab.com/thorchain/thornode/common"
)

func ParseOutboundMemoV1(parts []string) (OutboundMemo, error) {
	if len(parts) < 2 {
		return OutboundMemo{}, fmt.Errorf("not enough parameters")
	}
	txID, err := common.NewTxID(parts[1])
	return NewOutboundMemo(txID), err
}
