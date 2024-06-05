package thorchain

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/blang/semver"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

type SwapMemo struct {
	MemoBase
	Destination          common.Address
	SlipLimit            cosmos.Uint
	AffiliateAddress     common.Address
	AffiliateBasisPoints cosmos.Uint
	DexAggregator        string
	DexTargetAddress     string
	DexTargetLimit       *cosmos.Uint
	OrderType            types.OrderType
	StreamInterval       uint64
	StreamQuantity       uint64
	AffiliateTHORName    *types.THORName
	RefundAddress        common.Address
}

func (m SwapMemo) GetDestination() common.Address        { return m.Destination }
func (m SwapMemo) GetSlipLimit() cosmos.Uint             { return m.SlipLimit }
func (m SwapMemo) GetAffiliateAddress() common.Address   { return m.AffiliateAddress }
func (m SwapMemo) GetAffiliateBasisPoints() cosmos.Uint  { return m.AffiliateBasisPoints }
func (m SwapMemo) GetDexAggregator() string              { return m.DexAggregator }
func (m SwapMemo) GetDexTargetAddress() string           { return m.DexTargetAddress }
func (m SwapMemo) GetDexTargetLimit() *cosmos.Uint       { return m.DexTargetLimit }
func (m SwapMemo) GetOrderType() types.OrderType         { return m.OrderType }
func (m SwapMemo) GetStreamQuantity() uint64             { return m.StreamQuantity }
func (m SwapMemo) GetStreamInterval() uint64             { return m.StreamInterval }
func (m SwapMemo) GetAffiliateTHORName() *types.THORName { return m.AffiliateTHORName }
func (m SwapMemo) GetRefundAddress() common.Address      { return m.RefundAddress }

func (m SwapMemo) String() string {
	return m.string(false)
}

func (m SwapMemo) ShortString() string {
	return m.string(true)
}

func (m SwapMemo) string(short bool) string {
	slipLimit := m.SlipLimit.String()
	if m.SlipLimit.IsZero() {
		slipLimit = ""
	}
	if m.StreamInterval > 0 || m.StreamQuantity > 1 {
		slipLimit = fmt.Sprintf("%s/%d/%d", m.SlipLimit.String(), m.StreamInterval, m.StreamQuantity)
	}

	// prefer short notation for generate swap memo
	txType := m.TxType.String()
	if m.TxType == TxSwap {
		txType = "="
	}

	var assetString string
	if short && len(m.Asset.ShortCode()) > 0 {
		assetString = m.Asset.ShortCode()
	} else {
		assetString = m.Asset.String()
	}

	// destination + custom refund addr
	destString := m.Destination.String()
	if !m.RefundAddress.IsEmpty() {
		destString = m.Destination.String() + "/" + m.RefundAddress.String()
	}

	args := []string{
		txType,
		assetString,
		destString,
		slipLimit,
		m.AffiliateAddress.String(),
		m.AffiliateBasisPoints.String(),
		m.DexAggregator,
		m.DexTargetAddress,
	}

	last := 3
	if !m.SlipLimit.IsZero() || m.StreamInterval > 0 || m.StreamQuantity > 1 {
		last = 4
	}

	if !m.AffiliateAddress.IsEmpty() {
		last = 6
	}

	if m.DexAggregator != "" {
		last = 8
	}

	if m.DexTargetLimit != nil && !m.DexTargetLimit.IsZero() {
		args = append(args, m.DexTargetLimit.String())
		last = 9
	}

	return strings.Join(args[:last], ":")
}

func NewSwapMemo(asset common.Asset, dest common.Address, slip cosmos.Uint, affAddr common.Address, affPts cosmos.Uint, dexAgg, dexTargetAddress string, dexTargetLimit cosmos.Uint, orderType types.OrderType, quan, interval uint64, tn types.THORName, refundAddress common.Address) SwapMemo {
	swapMemo := SwapMemo{
		MemoBase:             MemoBase{TxType: TxSwap, Asset: asset},
		Destination:          dest,
		SlipLimit:            slip,
		AffiliateAddress:     affAddr,
		AffiliateBasisPoints: affPts,
		DexAggregator:        dexAgg,
		DexTargetAddress:     dexTargetAddress,
		OrderType:            orderType,
		StreamQuantity:       quan,
		StreamInterval:       interval,
		RefundAddress:        refundAddress,
	}
	if !dexTargetLimit.IsZero() {
		swapMemo.DexTargetLimit = &dexTargetLimit
	}
	if !tn.Owner.Empty() {
		swapMemo.AffiliateTHORName = &tn
	}
	return swapMemo
}

func (p *parser) ParseSwapMemo() (SwapMemo, error) {
	if p.keeper == nil {
		return ParseSwapMemoV1(p.ctx, p.keeper, p.getAsset(1, true, common.EmptyAsset), p.parts)
	}

	// TODO: remove me on hard fork
	var err error
	if len(p.parts) > 1 {
		_, err = common.NewAssetWithShortCodes(p.version, GetPart(p.parts, 1))
	}

	switch {
	case p.version.GTE(semver.MustParse("1.131.0")):
		return p.ParseSwapMemoV131()
	case p.version.GTE(semver.MustParse("1.123.0")):
		return p.ParseSwapMemoV123()
	case p.version.GTE(semver.MustParse("1.116.0")):
		return p.ParseSwapMemoV116()
	case err != nil:
		return SwapMemo{}, err // To resolve block 6130730 sync failure
	case p.version.GTE(semver.MustParse("1.115.0")):
		return ParseSwapMemoV115(p.ctx, p.keeper, p.getAsset(1, true, common.EmptyAsset), p.parts)
	case p.version.GTE(semver.MustParse("1.112.0")):
		return ParseSwapMemoV112(p.ctx, p.keeper, p.getAsset(1, true, common.EmptyAsset), p.parts)
	case p.version.GTE(semver.MustParse("1.104.0")):
		return ParseSwapMemoV104(p.ctx, p.keeper, p.getAsset(1, true, common.EmptyAsset), p.parts)
	case p.version.GTE(semver.MustParse("1.98.0")):
		return ParseSwapMemoV98(p.ctx, p.keeper, p.getAsset(1, true, common.EmptyAsset), p.parts)
	case p.version.GTE(semver.MustParse("1.92.0")):
		return ParseSwapMemoV92(p.ctx, p.keeper, p.getAsset(1, true, common.EmptyAsset), p.parts)
	default:
		return ParseSwapMemoV1(p.ctx, p.keeper, p.getAsset(1, true, common.EmptyAsset), p.parts)
	}
}

func (p *parser) ParseSwapMemoV131() (SwapMemo, error) {
	var err error
	asset := p.getAsset(1, true, common.EmptyAsset)
	var order types.OrderType
	if strings.EqualFold(p.parts[0], "limito") || strings.EqualFold(p.parts[0], "lo") {
		order = types.OrderType_limit
	}

	// DESTADDR can be empty , if it is empty , it will swap to the sender address
	destination, refundAddress := p.getAddressAndRefundAddressWithKeeper(2, false, common.NoAddress, asset.Chain)

	// price limit can be empty , when it is empty , there is no price protection
	var slip cosmos.Uint
	var streamInterval, streamQuantity uint64
	if strings.Contains(p.get(3), "/") {
		parts := strings.SplitN(p.get(3), "/", 3)
		for i := range parts {
			if parts[i] == "" {
				parts[i] = "0"
			}
		}
		if len(parts) < 1 {
			return SwapMemo{}, fmt.Errorf("invalid streaming swap format: %s", p.get(3))
		}
		slip, err = parseTradeTarget(parts[0])
		if err != nil {
			return SwapMemo{}, fmt.Errorf("swap price limit:%s is invalid: %s", parts[0], err)
		}
		if len(parts) > 1 {
			streamInterval, err = strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				return SwapMemo{}, fmt.Errorf("failed to parse stream frequency: %s: %s", parts[1], err)
			}
		}
		if len(parts) > 2 {
			streamQuantity, err = strconv.ParseUint(parts[2], 10, 64)
			if err != nil {
				return SwapMemo{}, fmt.Errorf("failed to parse stream quantity: %s: %s", parts[2], err)
			}
		}
	} else {
		slip = p.getUintWithScientificNotation(3, false, 0)
	}

	affAddr := p.getAddressWithKeeper(4, false, common.NoAddress, common.THORChain)
	affPts := p.getUintWithMaxValue(5, false, 0, constants.MaxBasisPts)

	dexAgg := p.get(6)
	dexTargetAddress := p.get(7)
	dexTargetLimit := p.getUintWithScientificNotation(8, false, 0)

	tn := p.getTHORName(4, false, types.NewTHORName("", 0, nil))

	return NewSwapMemo(asset, destination, slip, affAddr, affPts, dexAgg, dexTargetAddress, dexTargetLimit, order, streamQuantity, streamInterval, tn, refundAddress), p.Error()
}
