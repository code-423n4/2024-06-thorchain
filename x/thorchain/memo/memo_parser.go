package thorchain

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

type parser struct {
	memo           string
	txType         TxType
	ctx            cosmos.Context
	keeper         keeper.Keeper
	parts          []string
	errs           []error
	version        semver.Version
	requiredFields int
}

func newParser(ctx cosmos.Context, keeper keeper.Keeper, version semver.Version, memo string) (parser, error) {
	if len(memo) == 0 {
		return parser{}, fmt.Errorf("memo can't be empty")
	}
	if version.GTE(semver.MustParse("1.115.0")) {
		memo = strings.Split(memo, "|")[0]
	}
	parts := strings.Split(memo, ":")
	memoType, err := StringToTxType(parts[0])
	if err != nil {
		return parser{}, err
	}

	return parser{
		memo:    memo,
		txType:  memoType,
		ctx:     ctx,
		keeper:  keeper,
		version: version,
		parts:   parts,
		errs:    make([]error, 0),
	}, nil
}

func (p *parser) getType() TxType {
	return p.txType
}

func (p *parser) incRequired(required bool) {
	if required {
		p.requiredFields += 1
	}
}

func (p *parser) parse() (mem Memo, err error) {
	defer func() {
		if err == nil {
			err = p.Error()
		}
	}()
	switch p.getType() {
	case TxLeave:
		return p.ParseLeaveMemo()
	case TxDonate:
		return p.ParseDonateMemo()
	case TxAdd:
		return p.ParseAddLiquidityMemo()
	case TxWithdraw:
		return p.ParseWithdrawLiquidityMemo()
	case TxSwap, TxLimitOrder:
		if p.getType() == TxLimitOrder && p.version.LT(semver.MustParse("1.98.0")) {
			return EmptyMemo, fmt.Errorf("TxType not supported: %s", p.getType().String())
		}
		return p.ParseSwapMemo()
	case TxOutbound:
		return p.ParseOutboundMemo()
	case TxRefund:
		return p.ParseRefundMemo()
	case TxBond:
		return p.ParseBondMemo()
	case TxUnbond:
		return p.ParseUnbondMemo()
	case TxReserve:
		return p.ParseReserveMemo()
	case TxMigrate:
		return p.ParseMigrateMemo()
	case TxRagnarok:
		return p.ParseRagnarokMemo()
	case TxNoOp:
		return p.ParseNoOpMemo()
	case TxConsolidate:
		return p.ParseConsolidateMemo()
	case TxTHORName:
		return p.ParseManageTHORNameMemo()
	case TxLoanOpen:
		return p.ParseLoanOpenMemo()
	case TxLoanRepayment:
		return p.ParseLoanRepaymentMemo()
	case TxTradeAccountDeposit:
		return p.ParseTradeAccountDeposit()
	case TxTradeAccountWithdrawal:
		return p.ParseTradeAccountWithdrawal()

	case TxSwitch: // TODO remove on hard fork
		if p.version.GTE(semver.MustParse("1.117.0")) {
			return EmptyMemo, fmt.Errorf("TxType not supported: %s", p.getType().String())
		}
		return p.ParseSwitchMemo()
	case TxYggdrasilFund: // TODO remove on hard fork
		if p.version.GTE(semver.MustParse("1.124.0")) {
			return EmptyMemo, fmt.Errorf("TxType not supported: %s", p.getType().String())
		}
		return p.ParseYggdrasilFundMemo()
	case TxYggdrasilReturn: // TODO remove on hard fork
		if p.version.GTE(semver.MustParse("1.124.0")) {
			return EmptyMemo, fmt.Errorf("TxType not supported: %s", p.getType().String())
		}
		return p.ParseYggdrasilReturnMemo()

	default:
		return EmptyMemo, fmt.Errorf("TxType not supported: %s", p.getType().String())
	}
}

func (p *parser) addErr(err error) {
	p.errs = append(p.errs, err)
}

func (p *parser) Error() error {
	p.hasMinParams(p.requiredFields + 1)
	if len(p.errs) == 0 {
		return nil
	}
	errStrs := make([]string, len(p.errs))
	for i, err := range p.errs {
		errStrs[i] = err.Error()
	}
	err := fmt.Errorf("MEMO: %s\nPARSE FAILURE(S): %s", p.memo, strings.Join(errStrs, "-"))
	return err
}

// check if memo has enough parameters
func (p *parser) hasMinParams(count int) {
	if len(p.parts) < count {
		p.addErr(fmt.Errorf("not enough parameters: %d/%d", len(p.parts), count))
	}
}

// Safe accessor for split memo parts - always returns empty
// string for indices that are out of bounds.
func (p *parser) get(idx int) string {
	if idx < 0 || len(p.parts) <= idx {
		return ""
	}
	return p.parts[idx]
}

func (p *parser) getInt64(idx int, required bool, def int64) int64 {
	p.incRequired(required)
	value, err := strconv.ParseInt(p.get(idx), 10, 64)
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as an int64: %w", p.get(idx), err))
		}
		return def
	}
	return value
}

func (p *parser) getUint(idx int, required bool, def uint64) cosmos.Uint {
	p.incRequired(required)
	value, err := cosmos.ParseUint(p.get(idx))
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as an uint: %w", p.get(idx), err))
		}
		return cosmos.NewUint(def)
	}
	return value
}

func (p *parser) getUintWithScientificNotation(idx int, required bool, def uint64) cosmos.Uint {
	p.incRequired(required)
	f, _, err := big.ParseFloat(p.get(idx), 10, 0, big.ToZero)
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as an uint with sci notation: %w", p.get(idx), err))
		}
		return cosmos.NewUint(def)
	}
	i := new(big.Int)
	f.Int(i) // Note: fractional part will be discarded
	result := cosmos.NewUintFromBigInt(i)
	return result
}

func (p *parser) getUintWithMaxValue(idx int, required bool, def, max uint64) cosmos.Uint {
	value := p.getUint(idx, required, def)
	if value.GT(cosmos.NewUint(max)) {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("%s cannot exceed '%d'", p.get(idx), max))
		}
		return cosmos.NewUint(max)
	}
	return value
}

func (p *parser) getAccAddress(idx int, required bool, def cosmos.AccAddress) cosmos.AccAddress {
	p.incRequired(required)
	value, err := cosmos.AccAddressFromBech32(p.get(idx))
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as an AccAddress: %w", p.get(idx), err))
		}
		return def
	}
	return value
}

func (p *parser) getAddress(idx int, required bool, def common.Address) common.Address {
	p.incRequired(required)
	value, err := common.NewAddress(p.get(idx))
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as an Address: %w", p.get(idx), err))
		}
		return def
	}
	return value
}

func (p *parser) getAddressWithKeeper(idx int, required bool, def common.Address, chain common.Chain) common.Address {
	p.incRequired(required)
	if p.keeper == nil {
		return p.getAddress(2, required, common.NoAddress)
	}
	addr, err := FetchAddress(p.ctx, p.keeper, p.get(idx), chain)
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as an Address: %w", p.get(idx), err))
		}
	}
	return addr
}

func (p *parser) getAddressAndRefundAddressWithKeeper(idx int, required bool, def common.Address, chain common.Chain) (common.Address, common.Address) {
	p.incRequired(required)

	//nolint:ineffassign
	destination := common.NoAddress
	refundAddress := common.NoAddress
	addresses := p.get(idx)

	if strings.Contains(addresses, "/") {
		parts := strings.SplitN(addresses, "/", 2)
		if p.keeper == nil {
			dest, err := common.NewAddress(parts[0])
			if err != nil {
				if required || parts[0] != "" {
					p.addErr(fmt.Errorf("cannot parse '%s' as an Address: %w", parts[0], err))
				}
			}
			destination = dest
		} else {
			destination = p.getAddressFromString(parts[0], chain, required)
		}
		if len(parts) > 1 {
			refundAddress, _ = common.NewAddress(parts[1])
		}
	} else {
		destination = p.getAddressWithKeeper(idx, false, common.NoAddress, chain)
	}

	if destination.IsEmpty() && !refundAddress.IsEmpty() {
		p.addErr(fmt.Errorf("refund address is set but destination address is empty"))
	}

	return destination, refundAddress
}

func (p *parser) getAddressFromString(val string, chain common.Chain, required bool) common.Address {
	addr, err := FetchAddress(p.ctx, p.keeper, val, chain)
	if err != nil {
		if required || val != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as an Address: %w", val, err))
		}
	}
	return addr
}

func (p *parser) getChain(idx int, required bool, def common.Chain) common.Chain {
	p.incRequired(required)
	value, err := common.NewChain(p.get(idx))
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as a chain: %w", p.get(idx), err))
		}
		return def
	}
	return value
}

func (p *parser) getAsset(idx int, required bool, def common.Asset) common.Asset {
	p.incRequired(required)
	value, err := common.NewAssetWithShortCodes(p.version, p.get(idx))
	if err != nil && (required || p.get(idx) != "") {
		p.addErr(fmt.Errorf("cannot parse '%s' as an asset: %w", p.get(idx), err))
		return def
	}
	if value.IsTradeAsset() && p.version.LT(semver.MustParse("1.128.0")) {
		p.addErr(fmt.Errorf("trade assets are not yet supported: %s", p.get(idx)))
		return common.EmptyAsset
	}
	return value
}

func (p *parser) getTxID(idx int, required bool, def common.TxID) common.TxID {
	p.incRequired(required)
	value, err := common.NewTxID(p.get(idx))
	if err != nil {
		if required || p.get(idx) != "" {
			p.addErr(fmt.Errorf("cannot parse '%s' as tx hash: %w", p.get(idx), err))
		}
		return def
	}
	return value
}

func (p *parser) getTHORName(idx int, required bool, def types.THORName) types.THORName {
	p.incRequired(required)
	name := p.get(idx)
	if p.keeper == nil {
		return def
	}
	if p.keeper.THORNameExists(p.ctx, name) {
		tn, err := p.keeper.GetTHORName(p.ctx, name)
		if err != nil {
			if required || p.get(idx) != "" {
				p.addErr(fmt.Errorf("cannot parse '%s' as a THORName: %w", p.get(idx), err))
			}
		}
		return tn
	}
	return def
}
