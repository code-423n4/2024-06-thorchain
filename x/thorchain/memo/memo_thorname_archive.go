package thorchain

import (
	"fmt"
	"strconv"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

func ParseManageTHORNameMemoV1(parts []string) (ManageTHORNameMemo, error) {
	var err error
	var name string
	var owner cosmos.AccAddress
	preferredAsset := common.EmptyAsset
	expire := int64(0)

	if len(parts) < 4 {
		return ManageTHORNameMemo{}, fmt.Errorf("not enough parameters")
	}

	name = parts[1]
	chain, err := common.NewChain(parts[2])
	if err != nil {
		return ManageTHORNameMemo{}, err
	}

	addr, err := common.NewAddress(parts[3])
	if err != nil {
		return ManageTHORNameMemo{}, err
	}

	if len(parts) >= 5 {
		owner, err = cosmos.AccAddressFromBech32(parts[4])
		if err != nil {
			return ManageTHORNameMemo{}, err
		}
	}

	if len(parts) >= 6 {
		preferredAsset, err = common.NewAsset(parts[5])
		if err != nil {
			return ManageTHORNameMemo{}, err
		}
	}

	if len(parts) >= 7 {
		expire, err = strconv.ParseInt(parts[6], 10, 64)
		if err != nil {
			return ManageTHORNameMemo{}, err
		}
	}

	return NewManageTHORNameMemo(name, chain, addr, expire, preferredAsset, owner), nil
}
