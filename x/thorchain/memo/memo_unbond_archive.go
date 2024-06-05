package thorchain

import (
	"fmt"

	"gitlab.com/thorchain/thornode/common/cosmos"
)

func ParseUnbondMemoV81(parts []string) (UnbondMemo, error) {
	additional := cosmos.AccAddress{}
	if len(parts) < 3 {
		return UnbondMemo{}, fmt.Errorf("not enough parameters")
	}
	addr, err := cosmos.AccAddressFromBech32(parts[1])
	if err != nil {
		return UnbondMemo{}, fmt.Errorf("%s is an invalid thorchain address: %w", parts[1], err)
	}
	amt, err := cosmos.ParseUint(parts[2])
	if err != nil {
		return UnbondMemo{}, fmt.Errorf("fail to parse amount (%s): %w", parts[2], err)
	}
	if len(parts) >= 4 {
		additional, err = cosmos.AccAddressFromBech32(parts[3])
		if err != nil {
			return UnbondMemo{}, fmt.Errorf("%s is an invalid thorchain address: %w", parts[3], err)
		}
	}
	return NewUnbondMemo(addr, additional, amt), nil
}
