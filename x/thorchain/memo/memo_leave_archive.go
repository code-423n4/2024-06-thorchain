package thorchain

import (
	"fmt"

	"gitlab.com/thorchain/thornode/common/cosmos"
)

func ParseLeaveMemoV1(parts []string) (LeaveMemo, error) {
	if len(parts) < 2 {
		return LeaveMemo{}, fmt.Errorf("not enough parameters")
	}
	addr, err := cosmos.AccAddressFromBech32(parts[1])
	if err != nil {
		return LeaveMemo{}, fmt.Errorf("%s is an invalid thorchain address: %w", parts[1], err)
	}

	return NewLeaveMemo(addr), nil
}
