package thorchain

import (
	"fmt"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

func ParseWithdrawLiquidityMemoV1(asset common.Asset, parts []string) (WithdrawLiquidityMemo, error) {
	var err error
	if len(parts) < 2 {
		return WithdrawLiquidityMemo{}, fmt.Errorf("not enough parameters")
	}
	withdrawalBasisPts := cosmos.ZeroUint()
	withdrawalAsset := common.EmptyAsset
	if len(parts) > 2 {
		withdrawalBasisPts, err = cosmos.ParseUint(parts[2])
		if err != nil {
			return WithdrawLiquidityMemo{}, err
		}
		if withdrawalBasisPts.IsZero() || withdrawalBasisPts.GT(cosmos.NewUint(types.MaxWithdrawBasisPoints)) {
			return WithdrawLiquidityMemo{}, fmt.Errorf("withdraw amount %s is invalid", parts[2])
		}
	}
	if len(parts) > 3 {
		withdrawalAsset, err = common.NewAsset(parts[3])
		if err != nil {
			return WithdrawLiquidityMemo{}, err
		}
	}
	return NewWithdrawLiquidityMemo(asset, withdrawalBasisPts, withdrawalAsset), nil
}
