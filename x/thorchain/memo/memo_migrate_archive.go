package thorchain

import (
	"errors"
	"fmt"
	"strconv"
)

func ParseMigrateMemoV1(parts []string) (MigrateMemo, error) {
	if len(parts) < 2 {
		return MigrateMemo{}, errors.New("not enough parameters")
	}
	blockHeight, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return MigrateMemo{}, fmt.Errorf("fail to convert (%s) to a valid block height: %w", parts[1], err)
	}
	return NewMigrateMemo(blockHeight), nil
}
