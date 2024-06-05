package thorchain

import (
	"errors"

	"github.com/blang/semver"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
)

type TestRagnarokChainKeeper struct {
	keeper.KVStoreDummy
	activeVault Vault
	retireVault Vault
	yggVault    Vault
	pools       Pools
	lps         LiquidityProviders
	na          NodeAccount
	err         error
}

func (k *TestRagnarokChainKeeper) ListValidatorsWithBond(_ cosmos.Context) (NodeAccounts, error) {
	return NodeAccounts{k.na}, k.err
}

func (k *TestRagnarokChainKeeper) ListActiveValidators(_ cosmos.Context) (NodeAccounts, error) {
	return NodeAccounts{k.na}, k.err
}

func (k *TestRagnarokChainKeeper) GetNodeAccount(ctx cosmos.Context, signer cosmos.AccAddress) (NodeAccount, error) {
	if k.na.NodeAddress.Equals(signer) {
		return k.na, nil
	}
	return NodeAccount{}, nil
}

func (k *TestRagnarokChainKeeper) GetAsgardVaultsByStatus(_ cosmos.Context, vt VaultStatus) (Vaults, error) {
	if vt == ActiveVault {
		return Vaults{k.activeVault}, k.err
	}
	return Vaults{k.retireVault}, k.err
}

func (k *TestRagnarokChainKeeper) VaultExists(_ cosmos.Context, _ common.PubKey) bool {
	return true
}

func (k *TestRagnarokChainKeeper) GetVault(_ cosmos.Context, _ common.PubKey) (Vault, error) {
	return k.yggVault, k.err
}

func (k *TestRagnarokChainKeeper) GetMostSecure(ctx cosmos.Context, vaults Vaults, signingTransPeriod int64) Vault {
	return vaults[0]
}

func (k *TestRagnarokChainKeeper) GetLeastSecure(ctx cosmos.Context, vaults Vaults, signingTransPeriod int64) Vault {
	return vaults[0]
}

func (k *TestRagnarokChainKeeper) GetPools(_ cosmos.Context) (Pools, error) {
	return k.pools, k.err
}

func (k *TestRagnarokChainKeeper) GetPool(_ cosmos.Context, asset common.Asset) (Pool, error) {
	for _, pool := range k.pools {
		if pool.Asset.Equals(asset) {
			return pool, nil
		}
	}
	return Pool{}, errors.New("pool not found")
}

func (k *TestRagnarokChainKeeper) SetPool(_ cosmos.Context, pool Pool) error {
	for i, p := range k.pools {
		if p.Asset.Equals(pool.Asset) {
			k.pools[i] = pool
		}
	}
	return k.err
}

func (k *TestRagnarokChainKeeper) PoolExist(_ cosmos.Context, _ common.Asset) bool {
	return true
}

func (k *TestRagnarokChainKeeper) GetModuleAddress(_ string) (common.Address, error) {
	return common.NoAddress, nil
}

func (k *TestRagnarokChainKeeper) SetPOL(_ cosmos.Context, pol ProtocolOwnedLiquidity) error {
	return nil
}

func (k *TestRagnarokChainKeeper) GetPOL(_ cosmos.Context) (ProtocolOwnedLiquidity, error) {
	return NewProtocolOwnedLiquidity(), nil
}

func (k *TestRagnarokChainKeeper) GetLiquidityProviderIterator(ctx cosmos.Context, _ common.Asset) cosmos.Iterator {
	cdc := makeTestCodec()
	iter := keeper.NewDummyIterator()
	for _, lp := range k.lps {
		iter.AddItem([]byte("key"), cdc.MustMarshal(lp))
	}
	return iter
}

func (k *TestRagnarokChainKeeper) AddOwnership(ctx cosmos.Context, coin common.Coin, addr cosmos.AccAddress) error {
	lp, _ := common.NewAddress(addr.String())
	for i, skr := range k.lps {
		if lp.Equals(skr.RuneAddress) {
			k.lps[i].Units = k.lps[i].Units.Add(coin.Amount)
		}
	}
	return nil
}

func (k *TestRagnarokChainKeeper) RemoveOwnership(ctx cosmos.Context, coin common.Coin, addr cosmos.AccAddress) error {
	lp, _ := common.NewAddress(addr.String())
	for i, skr := range k.lps {
		if lp.Equals(skr.RuneAddress) {
			k.lps[i].Units = k.lps[i].Units.Sub(coin.Amount)
		}
	}
	return nil
}

func (k *TestRagnarokChainKeeper) GetLiquidityProvider(_ cosmos.Context, asset common.Asset, addr common.Address) (LiquidityProvider, error) {
	if asset.Equals(common.BTCAsset) {
		for i, lp := range k.lps {
			if addr.Equals(lp.RuneAddress) {
				return k.lps[i], k.err
			}
		}
	}
	return LiquidityProvider{}, k.err
}

func (k *TestRagnarokChainKeeper) SetLiquidityProvider(_ cosmos.Context, lp LiquidityProvider) {
	for i, skr := range k.lps {
		if lp.RuneAddress.Equals(skr.RuneAddress) {
			lp.Units = k.lps[i].Units
			k.lps[i] = lp
		}
	}
}

func (k *TestRagnarokChainKeeper) RemoveLiquidityProvider(_ cosmos.Context, lp LiquidityProvider) {
	for i, skr := range k.lps {
		if lp.RuneAddress.Equals(skr.RuneAddress) {
			k.lps[i] = lp
		}
	}
}

func (k *TestRagnarokChainKeeper) GetGas(_ cosmos.Context, _ common.Asset) ([]cosmos.Uint, error) {
	return []cosmos.Uint{cosmos.NewUint(10)}, k.err
}

func (k *TestRagnarokChainKeeper) GetLowestActiveVersion(_ cosmos.Context) semver.Version {
	return GetCurrentVersion()
}

func (k *TestRagnarokChainKeeper) AddPoolFeeToReserve(_ cosmos.Context, _ cosmos.Uint) error {
	return k.err
}

func (k *TestRagnarokChainKeeper) IsActiveObserver(_ cosmos.Context, _ cosmos.AccAddress) bool {
	return true
}

type VaultManagerTestHelpKeeper struct {
	keeper.Keeper
	failToGetAsgardVaults      bool
	failToListActiveAccounts   bool
	failToSetVault             bool
	failGetRetiringAsgardVault bool
	failGetActiveAsgardVault   bool
	failToSetPool              bool
	failGetNetwork             bool
	failGetTotalLiquidityFee   bool
	failGetPools               bool
}

func NewVaultGenesisSetupTestHelper(k keeper.Keeper) *VaultManagerTestHelpKeeper {
	return &VaultManagerTestHelpKeeper{
		Keeper: k,
	}
}

func (h *VaultManagerTestHelpKeeper) GetNetwork(ctx cosmos.Context) (Network, error) {
	if h.failGetNetwork {
		return Network{}, errKaboom
	}
	return h.Keeper.GetNetwork(ctx)
}

func (h *VaultManagerTestHelpKeeper) GetAsgardVaults(ctx cosmos.Context) (Vaults, error) {
	if h.failToGetAsgardVaults {
		return Vaults{}, errKaboom
	}
	return h.Keeper.GetAsgardVaults(ctx)
}

func (h *VaultManagerTestHelpKeeper) ListActiveValidators(ctx cosmos.Context) (NodeAccounts, error) {
	if h.failToListActiveAccounts {
		return NodeAccounts{}, errKaboom
	}
	return h.Keeper.ListActiveValidators(ctx)
}

func (h *VaultManagerTestHelpKeeper) SetVault(ctx cosmos.Context, v Vault) error {
	if h.failToSetVault {
		return errKaboom
	}
	return h.Keeper.SetVault(ctx, v)
}

func (h *VaultManagerTestHelpKeeper) GetAsgardVaultsByStatus(ctx cosmos.Context, vs VaultStatus) (Vaults, error) {
	if h.failGetRetiringAsgardVault && vs == RetiringVault {
		return Vaults{}, errKaboom
	}
	if h.failGetActiveAsgardVault && vs == ActiveVault {
		return Vaults{}, errKaboom
	}
	return h.Keeper.GetAsgardVaultsByStatus(ctx, vs)
}

func (h *VaultManagerTestHelpKeeper) SetPool(ctx cosmos.Context, p Pool) error {
	if h.failToSetPool {
		return errKaboom
	}
	return h.Keeper.SetPool(ctx, p)
}

func (h *VaultManagerTestHelpKeeper) GetTotalLiquidityFees(ctx cosmos.Context, height uint64) (cosmos.Uint, error) {
	if h.failGetTotalLiquidityFee {
		return cosmos.ZeroUint(), errKaboom
	}
	return h.Keeper.GetTotalLiquidityFees(ctx, height)
}

func (h *VaultManagerTestHelpKeeper) GetPools(ctx cosmos.Context) (Pools, error) {
	if h.failGetPools {
		return Pools{}, errKaboom
	}
	return h.Keeper.GetPools(ctx)
}
