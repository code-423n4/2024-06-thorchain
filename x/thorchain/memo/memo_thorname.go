package thorchain

import (
	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

type ManageTHORNameMemo struct {
	MemoBase
	Name           string
	Chain          common.Chain
	Address        common.Address
	PreferredAsset common.Asset
	Expire         int64
	Owner          cosmos.AccAddress
}

func (m ManageTHORNameMemo) GetName() string            { return m.Name }
func (m ManageTHORNameMemo) GetChain() common.Chain     { return m.Chain }
func (m ManageTHORNameMemo) GetAddress() common.Address { return m.Address }
func (m ManageTHORNameMemo) GetBlockExpire() int64      { return m.Expire }

func NewManageTHORNameMemo(name string, chain common.Chain, addr common.Address, expire int64, asset common.Asset, owner cosmos.AccAddress) ManageTHORNameMemo {
	return ManageTHORNameMemo{
		MemoBase:       MemoBase{TxType: TxTHORName},
		Name:           name,
		Chain:          chain,
		Address:        addr,
		PreferredAsset: asset,
		Expire:         expire,
		Owner:          owner,
	}
}

func (p *parser) ParseManageTHORNameMemo() (ManageTHORNameMemo, error) {
	switch {
	case p.version.GTE(semver.MustParse("1.116.0")):
		return p.ParseManageTHORNameMemoV116()
	default:
		return ParseManageTHORNameMemoV1(p.parts)
	}
}

func (p *parser) ParseManageTHORNameMemoV116() (ManageTHORNameMemo, error) {
	chain := p.getChain(2, true, common.EmptyChain)
	addr := p.getAddress(3, true, common.NoAddress)
	owner := p.getAccAddress(4, false, nil)
	preferredAsset := p.getAsset(5, false, common.EmptyAsset)
	expire := p.getInt64(6, false, 0)
	return NewManageTHORNameMemo(p.get(1), chain, addr, expire, preferredAsset, owner), p.Error()
}
