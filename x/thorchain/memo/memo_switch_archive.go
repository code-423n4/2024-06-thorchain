package thorchain

import (
	"errors"

	"github.com/blang/semver"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
)

func ParseSwitchMemoV1(ctx cosmos.Context, keeper keeper.Keeper, parts []string) (SwitchMemo, error) {
	if len(parts) < 2 {
		return SwitchMemo{}, errors.New("not enough parameters")
	}
	var destination common.Address
	var err error
	if keeper == nil {
		destination, err = common.NewAddress(parts[1])
	} else {
		destination, err = FetchAddress(ctx, keeper, parts[1], common.THORChain)
	}
	if err != nil {
		return SwitchMemo{}, err
	}
	if destination.IsEmpty() {
		return SwitchMemo{}, errors.New("address cannot be empty")
	}
	return NewSwitchMemo(destination), nil
}

type SwitchMemo struct {
	MemoBase
	Destination common.Address
}

func (m SwitchMemo) GetDestination() common.Address {
	return m.Destination
}

func NewSwitchMemo(addr common.Address) SwitchMemo {
	return SwitchMemo{
		MemoBase:    MemoBase{TxType: TxSwitch},
		Destination: addr,
	}
}

func (p *parser) ParseSwitchMemo() (SwitchMemo, error) {
	switch {
	case p.version.GTE(semver.MustParse("1.116.0")):
		return p.ParseSwitchMemoV116()
	default:
		return ParseSwitchMemoV1(p.ctx, p.keeper, p.parts)
	}
}

func (p *parser) ParseSwitchMemoV116() (SwitchMemo, error) {
	destination := p.getAddressWithKeeper(1, true, common.NoAddress, common.THORChain)
	return NewSwitchMemo(destination), p.Error()
}
