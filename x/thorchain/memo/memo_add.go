package thorchain

import (
	"strings"

	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common"
	cosmos "gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
)

type AddLiquidityMemo struct {
	MemoBase
	Address              common.Address
	AffiliateAddress     common.Address
	AffiliateBasisPoints cosmos.Uint
}

func (m AddLiquidityMemo) GetDestination() common.Address { return m.Address }

func (m AddLiquidityMemo) String() string {
	txType := m.TxType.String()
	if m.TxType == TxAdd {
		txType = "+"
	}

	args := []string{
		txType,
		m.Asset.String(),
		m.Address.String(),
		m.AffiliateAddress.String(),
		m.AffiliateBasisPoints.String(),
	}

	last := 2
	if !m.Address.IsEmpty() {
		last = 3
	}
	if !m.AffiliateAddress.IsEmpty() {
		last = 5
	}

	return strings.Join(args[:last], ":")
}

func NewAddLiquidityMemo(asset common.Asset, addr, affAddr common.Address, affPts cosmos.Uint) AddLiquidityMemo {
	return AddLiquidityMemo{
		MemoBase:             MemoBase{TxType: TxAdd, Asset: asset},
		Address:              addr,
		AffiliateAddress:     affAddr,
		AffiliateBasisPoints: affPts,
	}
}

func (p *parser) ParseAddLiquidityMemo() (AddLiquidityMemo, error) {
	if p.keeper == nil {
		return ParseAddLiquidityMemoV1(p.ctx, p.keeper, p.getAsset(1, true, common.EmptyAsset), p.parts)
	}
	switch {
	case p.version.GTE(semver.MustParse("1.128.0")):
		return p.ParseAddLiquidityMemoV128()
	case p.version.GTE(semver.MustParse("1.116.0")):
		return p.ParseAddLiquidityMemoV116()
	case p.version.GTE(semver.MustParse("1.104.0")):
		return ParseAddLiquidityMemoV104(p.ctx, p.keeper, p.getAsset(1, true, common.EmptyAsset), p.parts)
	default:
		return ParseAddLiquidityMemoV1(p.ctx, p.keeper, p.getAsset(1, true, common.EmptyAsset), p.parts)
	}
}

func (p *parser) ParseAddLiquidityMemoV128() (AddLiquidityMemo, error) {
	asset := p.getAsset(1, true, common.EmptyAsset)
	addr := p.getAddressWithKeeper(2, false, common.NoAddress, asset.Chain)
	affChain := common.THORChain
	if asset.IsSyntheticAsset() {
		// For a Savers add, an Affiliate THORName must be resolved
		// to an address for the Layer 1 Chain of the synth to succeed.
		affChain = asset.GetLayer1Asset().GetChain()
	}
	affAddr := p.getAddressWithKeeper(3, false, common.NoAddress, affChain)
	affPts := p.getUintWithMaxValue(4, false, 0, constants.MaxBasisPts)
	return NewAddLiquidityMemo(asset, addr, affAddr, affPts), p.Error()
}
