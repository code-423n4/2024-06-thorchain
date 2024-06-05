package thorchain

import (
	"errors"
	"fmt"
	"strconv"
)

func ParseRagnarokMemoV1(parts []string) (RagnarokMemo, error) {
	if len(parts) < 2 {
		return RagnarokMemo{}, errors.New("not enough parameters")
	}
	blockHeight, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return RagnarokMemo{}, fmt.Errorf("fail to convert (%s) to a valid block height: %w", parts[1], err)
	}
	return NewRagnarokMemo(blockHeight), nil
}
