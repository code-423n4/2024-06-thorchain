package thorchain

import (
	"fmt"

	"gitlab.com/thorchain/thornode/common"
)

func ParseRefundMemoV1(parts []string) (RefundMemo, error) {
	if len(parts) < 2 {
		return RefundMemo{}, fmt.Errorf("not enough parameters")
	}
	txID, err := common.NewTxID(parts[1])
	return NewRefundMemo(txID), err
}
