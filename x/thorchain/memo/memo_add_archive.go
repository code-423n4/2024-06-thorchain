package thorchain

// trunk-ignore-all(golangci-lint/govet)

import (
	"strconv"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
)

func (p *parser) ParseAddLiquidityMemoV116() (AddLiquidityMemo, error) {
	asset := p.getAsset(1, true, common.EmptyAsset)
	addr := p.getAddressWithKeeper(2, false, common.NoAddress, asset.Chain)
	affAddr := p.getAddressWithKeeper(3, false, common.NoAddress, common.THORChain)
	affPts := p.getUintWithMaxValue(4, false, 0, constants.MaxBasisPts)
	return NewAddLiquidityMemo(asset, addr, affAddr, affPts), p.Error()
}

func ParseAddLiquidityMemoV104(ctx cosmos.Context, keeper keeper.Keeper, asset common.Asset, parts []string) (AddLiquidityMemo, error) {
	var err error
	addr := common.NoAddress
	affAddr := common.NoAddress
	affPts := cosmos.ZeroUint()
	if addrStr := GetPart(parts, 2); addrStr != "" {
		if keeper == nil {
			addr, err = common.NewAddress(addrStr)
		} else {
			addr, err = FetchAddress(ctx, keeper, addrStr, asset.Chain)
		}
		if err != nil {
			return AddLiquidityMemo{}, err
		}
	}

	affAddrStr := GetPart(parts, 3)
	affPtsStr := GetPart(parts, 4)
	if affAddrStr != "" && affPtsStr != "" {
		if keeper == nil {
			affAddr, err = common.NewAddress(affAddrStr)
		} else {
			affAddr, err = FetchAddress(ctx, keeper, affAddrStr, common.THORChain)
		}
		if err != nil {
			return AddLiquidityMemo{}, err
		}
		affPts, err = ParseAffiliateBasisPoints(ctx, keeper, affPtsStr)
		if err != nil {
			return AddLiquidityMemo{}, err
		}
	}
	return NewAddLiquidityMemo(asset, addr, affAddr, affPts), nil
}

func ParseAddLiquidityMemoV1(ctx cosmos.Context, keeper keeper.Keeper, asset common.Asset, parts []string) (AddLiquidityMemo, error) {
	var err error
	addr := common.NoAddress
	affAddr := common.NoAddress
	affPts := cosmos.ZeroUint()
	if len(parts) >= 3 && len(parts[2]) > 0 {
		if keeper == nil {
			addr, err = common.NewAddress(parts[2])
		} else {
			addr, err = FetchAddress(ctx, keeper, parts[2], asset.Chain)
		}
		if err != nil {
			return AddLiquidityMemo{}, err
		}
	}

	if len(parts) > 4 && len(parts[3]) > 0 && len(parts[4]) > 0 {
		if keeper == nil {
			affAddr, err = common.NewAddress(parts[3])
		} else {
			affAddr, err = FetchAddress(ctx, keeper, parts[3], common.THORChain)
		}
		if err != nil {
			return AddLiquidityMemo{}, err
		}
		pts, err := strconv.ParseUint(parts[4], 10, 64)
		if err != nil {
			return AddLiquidityMemo{}, err
		}
		affPts = cosmos.NewUint(pts)
	}
	return NewAddLiquidityMemo(asset, addr, affAddr, affPts), nil
}
