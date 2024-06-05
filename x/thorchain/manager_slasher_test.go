package thorchain

import (
	"errors"

	"gitlab.com/thorchain/thornode/x/thorchain/keeper"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

type TestSlashingLackKeeper struct {
	keeper.KVStoreDummy
	txOut                      *TxOut
	na                         NodeAccount
	vaults                     Vaults
	voter                      ObservedTxVoter
	failGetTxOut               bool
	failGetVault               bool
	failGetNodeAccountByPubKey bool
	failSetNodeAccount         bool
	failGetAsgardByStatus      bool
	failGetObservedTxVoter     bool
	failSetTxOut               bool
	slashPts                   map[string]int64
}

func (k *TestSlashingLackKeeper) PoolExist(ctx cosmos.Context, asset common.Asset) bool {
	return true
}

func (k *TestSlashingLackKeeper) GetObservedTxInVoter(_ cosmos.Context, _ common.TxID) (ObservedTxVoter, error) {
	if k.failGetObservedTxVoter {
		return ObservedTxVoter{}, errKaboom
	}
	return k.voter, nil
}

func (k *TestSlashingLackKeeper) SetObservedTxInVoter(_ cosmos.Context, voter ObservedTxVoter) {
	k.voter = voter
}

func (k *TestSlashingLackKeeper) GetVault(_ cosmos.Context, pk common.PubKey) (Vault, error) {
	if k.failGetVault {
		return Vault{}, errKaboom
	}
	return k.vaults[0], nil
}

func (k *TestSlashingLackKeeper) GetAsgardVaultsByStatus(_ cosmos.Context, _ VaultStatus) (Vaults, error) {
	if k.failGetAsgardByStatus {
		return nil, errKaboom
	}
	return k.vaults, nil
}

func (k *TestSlashingLackKeeper) GetTxOut(_ cosmos.Context, _ int64) (*TxOut, error) {
	if k.failGetTxOut {
		return nil, errKaboom
	}
	return k.txOut, nil
}

func (k *TestSlashingLackKeeper) SetTxOut(_ cosmos.Context, tx *TxOut) error {
	if k.failSetTxOut {
		return errKaboom
	}
	k.txOut = tx
	return nil
}

func (k *TestSlashingLackKeeper) IncNodeAccountSlashPoints(_ cosmos.Context, addr cosmos.AccAddress, pts int64) error {
	if _, ok := k.slashPts[addr.String()]; !ok {
		k.slashPts[addr.String()] = 0
	}
	k.slashPts[addr.String()] += pts
	return nil
}

func (k *TestSlashingLackKeeper) GetNodeAccountByPubKey(_ cosmos.Context, _ common.PubKey) (NodeAccount, error) {
	if k.failGetNodeAccountByPubKey {
		return NodeAccount{}, errKaboom
	}
	return k.na, nil
}

func (k *TestSlashingLackKeeper) SetNodeAccount(_ cosmos.Context, na NodeAccount) error {
	if k.failSetNodeAccount {
		return errKaboom
	}
	k.na = na
	return nil
}

type TestSlashObservingKeeper struct {
	keeper.KVStoreDummy
	addrs                     []cosmos.AccAddress
	nas                       NodeAccounts
	failGetObservingAddress   bool
	failListActiveNodeAccount bool
	failSetNodeAccount        bool
	slashPts                  map[string]int64
}

func (k *TestSlashObservingKeeper) GetObservingAddresses(_ cosmos.Context) ([]cosmos.AccAddress, error) {
	if k.failGetObservingAddress {
		return nil, errKaboom
	}
	return k.addrs, nil
}

func (k *TestSlashObservingKeeper) ClearObservingAddresses(_ cosmos.Context) {
	k.addrs = nil
}

func (k *TestSlashObservingKeeper) IncNodeAccountSlashPoints(_ cosmos.Context, addr cosmos.AccAddress, pts int64) error {
	if _, ok := k.slashPts[addr.String()]; !ok {
		k.slashPts[addr.String()] = 0
	}
	k.slashPts[addr.String()] += pts
	return nil
}

func (k *TestSlashObservingKeeper) ListActiveValidators(_ cosmos.Context) (NodeAccounts, error) {
	if k.failListActiveNodeAccount {
		return nil, errKaboom
	}
	return k.nas, nil
}

func (k *TestSlashObservingKeeper) SetNodeAccount(_ cosmos.Context, na NodeAccount) error {
	if k.failSetNodeAccount {
		return errKaboom
	}
	for i := range k.nas {
		if k.nas[i].NodeAddress.Equals(na.NodeAddress) {
			k.nas[i] = na
			return nil
		}
	}
	return errors.New("node account not found")
}

type TestDoubleSlashKeeper struct {
	keeper.KVStoreDummy
	na          NodeAccount
	network     Network
	slashPoints map[string]int64
	modules     map[string]int64
}

func (k *TestDoubleSlashKeeper) SendFromModuleToModule(_ cosmos.Context, from, to string, coins common.Coins) error {
	k.modules[from] -= int64(coins[0].Amount.Uint64())
	k.modules[to] += int64(coins[0].Amount.Uint64())
	return nil
}

func (k *TestDoubleSlashKeeper) ListActiveValidators(ctx cosmos.Context) (NodeAccounts, error) {
	return NodeAccounts{k.na}, nil
}

func (k *TestDoubleSlashKeeper) GetNodeAccount(ctx cosmos.Context, nodeAddress cosmos.AccAddress) (NodeAccount, error) {
	if nodeAddress.String() == k.na.NodeAddress.String() {
		return k.na, nil
	}
	return NodeAccount{}, errors.New("kaboom")
}

func (k *TestDoubleSlashKeeper) SetNodeAccount(ctx cosmos.Context, na NodeAccount) error {
	k.na = na
	return nil
}

func (k *TestDoubleSlashKeeper) GetNetwork(ctx cosmos.Context) (Network, error) {
	return k.network, nil
}

func (k *TestDoubleSlashKeeper) SetNetwork(ctx cosmos.Context, data Network) error {
	k.network = data
	return nil
}

func (k *TestDoubleSlashKeeper) IncNodeAccountSlashPoints(ctx cosmos.Context, addr cosmos.AccAddress, pts int64) error {
	k.slashPoints[addr.String()] += pts
	return nil
}

func (k *TestDoubleSlashKeeper) DecNodeAccountSlashPoints(ctx cosmos.Context, addr cosmos.AccAddress, pts int64) error {
	k.slashPoints[addr.String()] -= pts
	return nil
}
