package thorchain

import (
	"fmt"
	"strings"

	"github.com/blang/semver"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

func (r *RouterUpgradeController) recallYggdrasilFund(ctx cosmos.Context, version semver.Version) error {
	switch {
	case version.GTE(semver.MustParse("1.94.0")):
		return r.recallYggdrasilFundV94(ctx, version)
	case version.GTE(semver.MustParse("0.1.0")):
		return r.recallYggdrasilFundV1(ctx)
	default:
		return fmt.Errorf("invalid version %s", version.String())
	}
}

func (r *RouterUpgradeController) recallYggdrasilFundV94(ctx cosmos.Context, version semver.Version) error {
	chains, err := r.getRouterChains(version)
	if err != nil {
		return fmt.Errorf("fail to get router chains: %w", err)
	}

	for _, chain := range chains {
		mimirKey := fmt.Sprintf(MimirRecallFundTemplate, chain)
		recallFund, err := r.mgr.Keeper().GetMimir(ctx, mimirKey)
		if err != nil {
			ctx.Logger().Error("fail to get recall funds mimir", "chain", chain.String(), "error", err)
			continue
		}

		if recallFund <= 0 {
			// mimir has not been set, continue
			continue
		}
		networkMgr := r.mgr.NetworkMgr()
		if err := networkMgr.RecallChainFunds(ctx, chain, r.mgr, common.PubKeys{}); err != nil {
			ctx.Logger().Error("fail to recall chain funds", "chain", chain.String(), "error", err)
			continue
		}
		err = r.mgr.Keeper().DeleteMimir(ctx, mimirKey)
		if err != nil {
			ctx.Logger().Error("fail to unset recall funds mimir", "chain", chain.String(), "error", err)
		}
	}

	return nil
}

func (r *RouterUpgradeController) recallYggdrasilFundV1(ctx cosmos.Context) error {
	recallFund, err := r.mgr.Keeper().GetMimir(ctx, MimirRecallFund)
	if err != nil {
		return fmt.Errorf("fail to get mimir: %w", err)
	}

	if recallFund <= 0 {
		// mimir has not been set , return
		return nil
	}
	networkMgr := r.mgr.NetworkMgr()
	if err := networkMgr.RecallChainFunds(ctx, common.ETHChain, r.mgr, common.PubKeys{}); err != nil {
		return fmt.Errorf("fail to recall chain funds, err:%w", err)
	}
	return r.mgr.Keeper().DeleteMimir(ctx, MimirRecallFund)
}

func (r *RouterUpgradeController) upgradeContractV94(ctx cosmos.Context, version semver.Version) error {
	chains, err := r.getRouterChains(version)
	if err != nil {
		return fmt.Errorf("fail to get router chains: %w", err)
	}

	// Iterate through all the chains with routers, see if any need their contracts updated
	for _, chain := range chains {
		mimirKey := fmt.Sprintf(MimirUpgradeContractTemplate, chain)
		mimirValue, err := r.mgr.Keeper().GetMimir(ctx, mimirKey)
		if err != nil {
			ctx.Logger().Error("fail to get router upgrade mimir", "chain", chain.String(), "error", err)
			continue
		}

		if mimirValue <= 0 {
			// mimir not set, skip
			continue
		}

		oldRouter, newRouter, err := r.getChainOldAndNewRouters(chain)
		if err != nil {
			ctx.Logger().Error("fail to get old and new router", "chain", chain.String(), "error", err)
			continue
		}

		currentChainContract, err := r.mgr.Keeper().GetChainContract(ctx, chain)
		if err != nil {
			ctx.Logger().Error("fail to get existing contract", "chain", chain.String(), "error", err)
			continue
		}

		// old router should be current router
		if !strings.EqualFold(currentChainContract.Router.String(), oldRouter) {
			ctx.Logger().Error("old router not current router", "chain", chain.String())
			continue
		}

		// new router should not be current router
		if strings.EqualFold(currentChainContract.Router.String(), newRouter) {
			ctx.Logger().Info("new router already set", "chain", chain.String())
			continue
		}

		// Update ChainContract
		// TODO: make this non-EVM agnostic (should not need to be an address)
		newRouterAddr, err := common.NewAddress(newRouter)
		if err != nil {
			ctx.Logger().Error("fail to parse new contract address", "chain", chain.String(), "addr", newRouter, "error", err)
			continue
		}
		newChainContract := ChainContract{
			Chain:  chain,
			Router: newRouterAddr,
		}
		r.mgr.Keeper().SetChainContract(ctx, newChainContract)

		// Update Vaults
		vaultIter := r.mgr.Keeper().GetVaultIterator(ctx)
		defer vaultIter.Close()
		for ; vaultIter.Valid(); vaultIter.Next() {
			var vault Vault
			if err := r.mgr.Keeper().Cdc().Unmarshal(vaultIter.Value(), &vault); err != nil {
				ctx.Logger().Error("fail to unmarshal vault", "error", err)
				continue
			}
			// vault is empty, ignore
			if vault.IsEmpty() {
				continue
			}

			if vault.IsType(YggdrasilVault) {
				// update yggdrasil vault to use new router contract
				vault.UpdateContract(newChainContract)
			}
			if err := r.mgr.Keeper().SetVault(ctx, vault); err != nil {
				ctx.Logger().Error("fail to save vault", "error", err)
			}

		}

		// Unset upgrade router mimir
		err = r.mgr.Keeper().DeleteMimir(ctx, mimirKey)
		if err != nil {
			ctx.Logger().Error("fail to unset router upgrade mimir", "chain", chain.String(), "error", err)
		}
	}

	return nil
}

func (r *RouterUpgradeController) upgradeContractV1(ctx cosmos.Context) error {
	upgrade, err := r.mgr.Keeper().GetMimir(ctx, MimirUpgradeContract)
	if err != nil {
		return fmt.Errorf("fail to get mimir: %w", err)
	}

	if upgrade <= 0 {
		// mimir has not been set , return
		return nil
	}

	newRouterAddr, err := common.NewAddress(ethNewRouter)
	if err != nil {
		return fmt.Errorf("fail to parse router address, err: %w", err)
	}
	cc, err := r.mgr.Keeper().GetChainContract(ctx, common.ETHChain)
	if err != nil {
		return fmt.Errorf("fail to get existing chain contract,err:%w", err)
	}
	// ensure it is upgrading the current router contract
	if !strings.EqualFold(cc.Router.String(), ethOldRouter) {
		return fmt.Errorf("old router is not %s, no need to upgrade", ethOldRouter)
	}
	chainContract := ChainContract{
		Chain:  common.ETHChain,
		Router: newRouterAddr,
	}
	// Update the contract address
	r.mgr.Keeper().SetChainContract(ctx, chainContract)

	vaultIter := r.mgr.Keeper().GetVaultIterator(ctx)
	defer vaultIter.Close()
	for ; vaultIter.Valid(); vaultIter.Next() {
		var vault Vault
		if err := r.mgr.Keeper().Cdc().Unmarshal(vaultIter.Value(), &vault); err != nil {
			ctx.Logger().Error("fail to unmarshal vault", "error", err)
			continue
		}
		// vault is empty , ignore
		if vault.IsEmpty() {
			continue
		}

		if vault.IsType(YggdrasilVault) {
			// update yggdrasil vault to use new router contract
			vault.UpdateContract(chainContract)
		}
		if err := r.mgr.Keeper().SetVault(ctx, vault); err != nil {
			ctx.Logger().Error("fail to save vault", "error", err)
		}

	}

	return r.mgr.Keeper().DeleteMimir(ctx, MimirUpgradeContract)
}
