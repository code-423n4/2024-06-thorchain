package thorchain

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/blang/semver"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog/log"
	abci "github.com/tendermint/tendermint/abci/types"
	tmhttp "github.com/tendermint/tendermint/rpc/client/http"

	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/config"
	"gitlab.com/thorchain/thornode/constants"
	"gitlab.com/thorchain/thornode/mimir"
	openapi "gitlab.com/thorchain/thornode/openapi/gen"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
	q "gitlab.com/thorchain/thornode/x/thorchain/query"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
	"gitlab.com/thorchain/tss/go-tss/conversion"
)

var (
	initManager   = func(mgr *Mgrs, ctx cosmos.Context) {}
	optionalQuery = func(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
		return nil, cosmos.ErrUnknownRequest(
			fmt.Sprintf("unknown thorchain query endpoint: %s", path[0]),
		)
	}
	tendermintClient   *tmhttp.HTTP
	initTendermintOnce = sync.Once{}
)

func initTendermint() {
	// get tendermint port from config
	portSplit := strings.Split(config.GetThornode().Tendermint.RPC.ListenAddress, ":")
	port := portSplit[len(portSplit)-1]

	// setup tendermint client
	var err error
	tendermintClient, err = tmhttp.New(fmt.Sprintf("tcp://localhost:%s", port), "/websocket")
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create tendermint client")
	}
}

// NewQuerier is the module level router for state queries
func NewQuerier(mgr *Mgrs, kbs cosmos.KeybaseStore) cosmos.Querier {
	return func(ctx cosmos.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		initManager(mgr, ctx) // NOOP except regtest

		defer telemetry.MeasureSince(time.Now(), path[0])
		switch path[0] {
		case q.QueryPool.Key:
			return queryPool(ctx, path[1:], req, mgr)
		case q.QueryPools.Key:
			return queryPools(ctx, req, mgr)
		case q.QueryDerivedPool.Key:
			return queryDerivedPool(ctx, path[1:], req, mgr)
		case q.QueryDerivedPools.Key:
			return queryDerivedPools(ctx, req, mgr)
		case q.QuerySavers.Key:
			return queryLiquidityProviders(ctx, path[1:], req, mgr, true)
		case q.QuerySaver.Key:
			return queryLiquidityProvider(ctx, path[1:], req, mgr, true)
		case q.QueryBorrowers.Key:
			return queryBorrowers(ctx, path[1:], req, mgr)
		case q.QueryBorrower.Key:
			return queryBorrower(ctx, path[1:], req, mgr)
		case q.QueryLiquidityProviders.Key:
			return queryLiquidityProviders(ctx, path[1:], req, mgr, false)
		case q.QueryLiquidityProvider.Key:
			return queryLiquidityProvider(ctx, path[1:], req, mgr, false)
		case q.QueryTradeUnit.Key:
			return queryTradeUnit(ctx, path[1:], req, mgr)
		case q.QueryTradeUnits.Key:
			return queryTradeUnits(ctx, path[1:], req, mgr)
		case q.QueryTradeAccount.Key:
			return queryTradeAccount(ctx, path[1:], req, mgr)
		case q.QueryTradeAccounts.Key:
			return queryTradeAccounts(ctx, path[1:], req, mgr)
		case q.QueryTxStages.Key:
			return queryTxStages(ctx, path[1:], req, mgr)
		case q.QueryTxStatus.Key:
			return queryTxStatus(ctx, path[1:], req, mgr)
		case q.QueryTxVoter.Key:
			return queryTxVoters(ctx, path[1:], req, mgr)
		case q.QueryTxVoterOld.Key:
			return queryTxVoters(ctx, path[1:], req, mgr)
		case q.QueryTx.Key:
			return queryTx(ctx, path[1:], req, mgr)
		case q.QueryKeysignArray.Key:
			return queryKeysign(ctx, kbs, path[1:], req, mgr)
		case q.QueryKeysignArrayPubkey.Key:
			return queryKeysign(ctx, kbs, path[1:], req, mgr)
		case q.QueryKeygensPubkey.Key:
			return queryKeygen(ctx, kbs, path[1:], req, mgr)
		case q.QueryQueue.Key:
			return queryQueue(ctx, path[1:], req, mgr)
		case q.QueryHeights.Key:
			return queryLastBlockHeights(ctx, path[1:], req, mgr)
		case q.QueryChainHeights.Key:
			return queryLastBlockHeights(ctx, path[1:], req, mgr)
		case q.QueryNode.Key:
			return queryNode(ctx, path[1:], req, mgr)
		case q.QueryNodes.Key:
			return queryNodes(ctx, path[1:], req, mgr)
		case q.QueryInboundAddresses.Key:
			return queryInboundAddresses(ctx, path[1:], req, mgr)
		case q.QueryNetwork.Key:
			return queryNetwork(ctx, mgr)
		case q.QueryPOL.Key:
			return queryPOL(ctx, mgr)
		case q.QueryBalanceModule.Key:
			return queryBalanceModule(ctx, path[1:], mgr)
		case q.QueryVaultsAsgard.Key:
			return queryAsgardVaults(ctx, mgr)
		case q.QueryVaultsYggdrasil.Key: // TODO remove on hard fork
			return queryYggdrasilVaults(ctx, mgr)
		case q.QueryVault.Key:
			return queryVault(ctx, path[1:], mgr)
		case q.QueryVaultPubkeys.Key:
			return queryVaultsPubkeys(ctx, mgr)
		case q.QueryConstantValues.Key:
			return queryConstantValues(ctx, path[1:], req, mgr)
		case q.QueryVersion.Key:
			return queryVersion(ctx, path[1:], req, mgr)
		case q.QueryMimirValues.Key:
			return queryMimirValues(ctx, path[1:], req, mgr)
		case q.QueryMimirV2Values.Key:
			return queryMimirV2Values(ctx, path[1:], req, mgr)
		case q.QueryMimirV2NodesAll.Key:
			return queryMimirV2NodesAll(ctx, path[1:], req, mgr)
		case q.QueryMimirV2IDs.Key:
			return queryMimirV2IDs(ctx, path[1:], req, mgr)
		case q.QueryMimirWithKey.Key:
			return queryMimirWithKey(ctx, path[1:], req, mgr)
		case q.QueryMimirAdminValues.Key:
			return queryMimirAdminValues(ctx, path[1:], req, mgr)
		case q.QueryMimirNodesAllValues.Key:
			return queryMimirNodesAllValues(ctx, path[1:], req, mgr)
		case q.QueryMimirNodesValues.Key:
			return queryMimirNodesValues(ctx, path[1:], req, mgr)
		case q.QueryMimirNodeValues.Key:
			return queryMimirNodeValues(ctx, path[1:], req, mgr)
		case q.QueryOutboundFees.Key:
			return queryOutboundFees(ctx, path[1:], req, mgr)
		case q.QueryOutboundFee.Key:
			return queryOutboundFees(ctx, path[1:], req, mgr)
		case q.QueryBan.Key:
			return queryBan(ctx, path[1:], req, mgr)
		case q.QueryRagnarok.Key:
			return queryRagnarok(ctx, mgr)
		case q.QueryPendingOutbound.Key:
			return queryPendingOutbound(ctx, mgr)
		case q.QueryScheduledOutbound.Key:
			return queryScheduledOutbound(ctx, mgr)
		case q.QuerySwapQueue.Key:
			return querySwapQueue(ctx, mgr)
		case q.QueryPoolSlip.Key:
			return queryPoolSlips(ctx, path[1:], req, mgr)
		case q.QueryPoolSlips.Key:
			return queryPoolSlips(ctx, path[1:], req, mgr)
		case q.QuerySwapperClout.Key:
			return querySwapperClout(ctx, path[1:], mgr)
		case q.QueryStreamingSwap.Key:
			return queryStreamingSwap(ctx, path[1:], mgr)
		case q.QueryStreamingSwaps.Key:
			return queryStreamingSwaps(ctx, mgr)
		case q.QueryTssKeygenMetrics.Key:
			return queryTssKeygenMetric(ctx, path[1:], req, mgr)
		case q.QueryTssMetrics.Key:
			return queryTssMetric(ctx, path[1:], req, mgr)
		case q.QueryTHORName.Key:
			return queryTHORName(ctx, path[1:], req, mgr)
		case q.QueryQuoteSwap.Key:
			return queryQuoteSwap(ctx, path[1:], req, mgr)
		case q.QueryQuoteSaverDeposit.Key:
			return queryQuoteSaverDeposit(ctx, path[1:], req, mgr)
		case q.QueryQuoteSaverWithdraw.Key:
			return queryQuoteSaverWithdraw(ctx, path[1:], req, mgr)
		case q.QueryQuoteLoanOpen.Key:
			return queryQuoteLoanOpen(ctx, path[1:], req, mgr)
		case q.QueryQuoteLoanClose.Key:
			return queryQuoteLoanClose(ctx, path[1:], req, mgr)
		case q.QueryInvariants.Key:
			return queryInvariants(ctx, mgr)
		case q.QueryInvariant.Key:
			return queryInvariant(ctx, path[1:], mgr)
		case q.QueryBlock.Key:
			return queryBlock(ctx, mgr)
		default:
			return optionalQuery(ctx, path, req, mgr)
		}
	}
}

func getPeerIDFromPubKey(pubkey common.PubKey) string {
	peerID, err := conversion.GetPeerIDFromPubKey(pubkey.String())
	if err != nil {
		// Don't break the entire endpoint if something goes wrong with the Peer ID derivation.
		return err.Error()
	}

	return peerID.String()
}

func jsonify(ctx cosmos.Context, r any) ([]byte, error) {
	res, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		ctx.Logger().Error("fail to marshal response to json", "error", err)
		return nil, fmt.Errorf("fail to marshal response to json: %w", err)
	}
	return res, nil
}

func queryRagnarok(ctx cosmos.Context, mgr *Mgrs) ([]byte, error) {
	ragnarokInProgress := mgr.Keeper().RagnarokInProgress(ctx)
	return jsonify(ctx, ragnarokInProgress)
}

func queryBalanceModule(ctx cosmos.Context, path []string, mgr *Mgrs) ([]byte, error) {
	moduleName := path[0]
	if len(moduleName) == 0 {
		moduleName = AsgardName
	}

	modAddr := mgr.Keeper().GetModuleAccAddress(moduleName)
	bal := mgr.Keeper().GetBalance(ctx, modAddr)
	balance := struct {
		Name    string            `json:"name"`
		Address cosmos.AccAddress `json:"address"`
		Coins   sdk.Coins         `json:"coins"`
	}{
		Name:    moduleName,
		Address: modAddr,
		Coins:   bal,
	}
	return jsonify(ctx, balance)
}

func queryTHORName(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	name, err := mgr.Keeper().GetTHORName(ctx, path[0])
	if err != nil {
		return nil, ErrInternal(err, "fail to fetch THORName")
	}

	affRune := cosmos.ZeroUint()
	affCol, err := mgr.Keeper().GetAffiliateCollector(ctx, name.Owner)
	if err == nil {
		affRune = affCol.RuneAmount
	}

	// convert to openapi types
	aliases := []openapi.ThornameAlias{}
	for _, alias := range name.Aliases {
		aliases = append(aliases, openapi.ThornameAlias{
			Chain:   wrapString(alias.Chain.String()),
			Address: wrapString(alias.Address.String()),
		})
	}

	resp := openapi.Thorname{
		Name:                   wrapString(name.Name),
		ExpireBlockHeight:      wrapInt64(name.ExpireBlockHeight),
		Owner:                  wrapString(name.Owner.String()),
		PreferredAsset:         name.PreferredAsset.String(),
		Aliases:                aliases,
		AffiliateCollectorRune: wrapString(affRune.String()),
	}

	return jsonify(ctx, resp)
}

func queryVault(ctx cosmos.Context, path []string, mgr *Mgrs) ([]byte, error) {
	if len(path) < 1 {
		return nil, errors.New("not enough parameters")
	}
	pubkey, err := common.NewPubKey(path[0])
	if err != nil {
		return nil, fmt.Errorf("%s is invalid pubkey", path[0])
	}
	v, err := mgr.Keeper().GetVault(ctx, pubkey)
	if err != nil {
		return nil, fmt.Errorf("fail to get vault with pubkey(%s),err:%w", pubkey, err)
	}
	if v.IsEmpty() {
		return nil, errors.New("vault not found")
	}

	resp := openapi.Vault{
		BlockHeight:           wrapInt64(v.BlockHeight),
		PubKey:                wrapString(v.PubKey.String()),
		Coins:                 castCoins(v.Coins...),
		Type:                  wrapString(v.Type.String()),
		Status:                v.Status.String(),
		StatusSince:           wrapInt64(v.StatusSince),
		Membership:            v.Membership,
		Chains:                v.Chains,
		InboundTxCount:        wrapInt64(v.InboundTxCount),
		OutboundTxCount:       wrapInt64(v.OutboundTxCount),
		PendingTxBlockHeights: v.PendingTxBlockHeights,
		Routers:               castVaultRouters(v.Routers),
		Addresses:             getVaultChainAddresses(ctx, v),
		Frozen:                v.Frozen,
	}
	return jsonify(ctx, resp)
}

func queryAsgardVaults(ctx cosmos.Context, mgr *Mgrs) ([]byte, error) {
	vaults, err := mgr.Keeper().GetAsgardVaults(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get asgard vaults: %w", err)
	}

	var vaultsWithFunds []openapi.Vault
	for _, vault := range vaults {
		if vault.Status == InactiveVault {
			continue
		}
		if !vault.IsAsgard() {
			continue
		}
		// Being in a RetiringVault blocks a node from unbonding, so display them even if having no funds.
		if vault.HasFunds() || vault.Status == ActiveVault || vault.Status == RetiringVault {
			vaultsWithFunds = append(vaultsWithFunds, openapi.Vault{
				BlockHeight:           wrapInt64(vault.BlockHeight),
				PubKey:                wrapString(vault.PubKey.String()),
				Coins:                 castCoins(vault.Coins...),
				Type:                  wrapString(vault.Type.String()),
				Status:                vault.Status.String(),
				StatusSince:           wrapInt64(vault.StatusSince),
				Membership:            vault.Membership,
				Chains:                vault.Chains,
				InboundTxCount:        wrapInt64(vault.InboundTxCount),
				OutboundTxCount:       wrapInt64(vault.OutboundTxCount),
				PendingTxBlockHeights: vault.PendingTxBlockHeights,
				Routers:               castVaultRouters(vault.Routers),
				Frozen:                vault.Frozen,
				Addresses:             getVaultChainAddresses(ctx, vault),
			})
		}
	}

	return jsonify(ctx, vaultsWithFunds)
}

func getVaultChainAddresses(ctx cosmos.Context, vault Vault) []openapi.VaultAddress {
	var result []openapi.VaultAddress
	allChains := append(vault.GetChains(), common.THORChain)
	for _, c := range allChains.Distinct() {
		addr, err := vault.PubKey.GetAddress(c)
		if err != nil {
			ctx.Logger().Error("fail to get address for %s:%w", c.String(), err)
			continue
		}
		result = append(result,
			openapi.VaultAddress{
				Chain:   c.String(),
				Address: addr.String(),
			})
	}
	return result
}

// TODO remove on hard fork
func queryYggdrasilVaults(ctx cosmos.Context, mgr *Mgrs) ([]byte, error) {
	vaults := make(Vaults, 0)
	iter := mgr.Keeper().GetVaultIterator(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var vault Vault
		if err := mgr.Keeper().Cdc().Unmarshal(iter.Value(), &vault); err != nil {
			ctx.Logger().Error("fail to unmarshal yggdrasil", "error", err)
			return nil, fmt.Errorf("fail to unmarshal yggdrasil: %w", err)
		}
		if vault.IsYggdrasil() && vault.HasFunds() {
			vaults = append(vaults, vault)
		}
	}

	respVaults := make([]openapi.YggdrasilVault, len(vaults))
	for i, vault := range vaults {
		totalValue := cosmos.ZeroUint()

		// find the bond of this node account
		na, err := mgr.Keeper().GetNodeAccountByPubKey(ctx, vault.PubKey)
		if err != nil {
			ctx.Logger().Error("fail to get node account by pubkey", "error", err)
			continue
		}

		// calculate the total value of this yggdrasil vault
		for _, coin := range vault.Coins {
			if coin.Asset.IsRune() {
				totalValue = totalValue.Add(coin.Amount)
			} else {
				pool, err := mgr.Keeper().GetPool(ctx, coin.Asset)
				if err != nil {
					ctx.Logger().Error("fail to get pool", "error", err)
					continue
				}
				totalValue = totalValue.Add(pool.AssetValueInRune(coin.Amount))
			}
		}

		respVaults[i] = openapi.YggdrasilVault{
			BlockHeight:           wrapInt64(vault.BlockHeight),
			PubKey:                wrapString(vault.PubKey.String()),
			Coins:                 castCoins(vault.Coins...),
			Type:                  wrapString(vault.Type.String()),
			StatusSince:           wrapInt64(vault.StatusSince),
			Membership:            vault.Membership,
			Chains:                vault.Chains,
			InboundTxCount:        wrapInt64(vault.InboundTxCount),
			OutboundTxCount:       wrapInt64(vault.OutboundTxCount),
			PendingTxBlockHeights: vault.PendingTxBlockHeights,
			Routers:               castVaultRouters(vault.Routers),
			Status:                na.Status.String(),
			Bond:                  na.Bond.String(),
			TotalValue:            totalValue.String(),
			Addresses:             getVaultChainAddresses(ctx, vault),
		}
	}

	return jsonify(ctx, respVaults)
}

func queryVaultsPubkeys(ctx cosmos.Context, mgr *Mgrs) ([]byte, error) {
	var resp openapi.VaultPubkeysResponse
	resp.Asgard = make([]openapi.VaultInfo, 0)
	resp.Yggdrasil = make([]openapi.VaultInfo, 0) // TODO remove on hard fork
	resp.Inactive = make([]openapi.VaultInfo, 0)
	iter := mgr.Keeper().GetVaultIterator(ctx)

	active, err := mgr.Keeper().ListActiveValidators(ctx)
	if err != nil {
		return nil, err
	}
	cutOffAge := ctx.BlockHeight() - config.GetThornode().VaultPubkeysCutoffBlocks
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var vault Vault
		if err := mgr.Keeper().Cdc().Unmarshal(iter.Value(), &vault); err != nil {
			ctx.Logger().Error("fail to unmarshal vault", "error", err)
			return nil, fmt.Errorf("fail to unmarshal vault: %w", err)
		}
		if vault.IsYggdrasil() { // TODO remove ygg on hard fork
			na, err := mgr.Keeper().GetNodeAccountByPubKey(ctx, vault.PubKey)
			if err != nil {
				ctx.Logger().Error("fail to unmarshal vault", "error", err)
				return nil, fmt.Errorf("fail to unmarshal vault: %w", err)
			}
			if !na.Bond.IsZero() {
				resp.Yggdrasil = append(resp.Yggdrasil, openapi.VaultInfo{
					PubKey:  vault.PubKey.String(),
					Routers: castVaultRouters(vault.Routers),
				})
			}
		} else if vault.IsAsgard() {
			switch vault.Status {
			case ActiveVault, RetiringVault:
				resp.Asgard = append(resp.Asgard, openapi.VaultInfo{
					PubKey:  vault.PubKey.String(),
					Routers: castVaultRouters(vault.Routers),
				})
			case InactiveVault:
				// skip inactive vaults that have never received an inbound
				if vault.InboundTxCount == 0 {
					continue
				}

				// skip inactive vaults older than the cutoff age
				if vault.BlockHeight < cutOffAge {
					continue
				}

				activeMembers, err := vault.GetMembers(active.GetNodeAddresses())
				if err != nil {
					ctx.Logger().Error("fail to get active members of vault", "error", err)
					continue
				}
				allMembers := vault.Membership
				if HasSuperMajority(len(activeMembers), len(allMembers)) {
					resp.Inactive = append(resp.Inactive, openapi.VaultInfo{
						PubKey:  vault.PubKey.String(),
						Routers: castVaultRouters(vault.Routers),
					})
				}
			}
		}
	}
	return jsonify(ctx, resp)
}

func queryNetwork(ctx cosmos.Context, mgr *Mgrs) ([]byte, error) {
	data, err := mgr.Keeper().GetNetwork(ctx)
	if err != nil {
		ctx.Logger().Error("fail to get network", "error", err)
		return nil, fmt.Errorf("fail to get network: %w", err)
	}

	vaults, err := mgr.Keeper().GetAsgardVaultsByStatus(ctx, RetiringVault)
	if err != nil {
		return nil, fmt.Errorf("fail to get retiring vaults: %w", err)
	}
	vaultsMigrating := (len(vaults) != 0)

	nodeAccounts, err := mgr.Keeper().ListActiveValidators(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get active validators: %w", err)
	}

	effectiveSecurityBond := getEffectiveSecurityBond(nodeAccounts)

	targetOutboundFeeSurplus := mgr.Keeper().GetConfigInt64(ctx, constants.TargetOutboundFeeSurplusRune)
	maxMultiplierBasisPoints := mgr.Keeper().GetConfigInt64(ctx, constants.MaxOutboundFeeMultiplierBasisPoints)
	minMultiplierBasisPoints := mgr.Keeper().GetConfigInt64(ctx, constants.MinOutboundFeeMultiplierBasisPoints)
	outboundFeeMultiplier := mgr.gasMgr.CalcOutboundFeeMultiplier(ctx, cosmos.NewUint(uint64(targetOutboundFeeSurplus)), cosmos.NewUint(data.OutboundGasSpentRune), cosmos.NewUint(data.OutboundGasWithheldRune), cosmos.NewUint(uint64(maxMultiplierBasisPoints)), cosmos.NewUint(uint64(minMultiplierBasisPoints)))

	result := openapi.NetworkResponse{
		// Due to using openapi. this will be displayed in alphabetical order,
		// so its schema (and order here) should also be in alphabetical order.
		BondRewardRune:        data.BondRewardRune.String(),
		BurnedBep2Rune:        data.BurnedBep2Rune.String(),  // TODO remove on hard fork
		BurnedErc20Rune:       data.BurnedErc20Rune.String(), // TODO remove on hard fork
		TotalBondUnits:        data.TotalBondUnits.String(),
		EffectiveSecurityBond: effectiveSecurityBond.String(),
		TotalReserve:          mgr.Keeper().GetRuneBalanceOfModule(ctx, ReserveName).String(),
		VaultsMigrating:       vaultsMigrating,
		GasSpentRune:          cosmos.NewUint(data.OutboundGasSpentRune).String(),
		GasWithheldRune:       cosmos.NewUint(data.OutboundGasWithheldRune).String(),
		OutboundFeeMultiplier: wrapString(outboundFeeMultiplier.String()),
		NativeTxFeeRune:       mgr.Keeper().GetNativeTxFee(ctx).String(),
		NativeOutboundFeeRune: mgr.Keeper().GetOutboundTxFee(ctx).String(),
		TnsRegisterFeeRune:    mgr.Keeper().GetTHORNameRegisterFee(ctx).String(),
		TnsFeePerBlockRune:    mgr.Keeper().GetTHORNamePerBlockFee(ctx).String(),
		RunePriceInTor:        mgr.Keeper().DollarsPerRune(ctx).String(),
		TorPriceInRune:        mgr.Keeper().RunePerDollar(ctx).String(),
	}

	return jsonify(ctx, result)
}

func queryPOL(ctx cosmos.Context, mgr *Mgrs) ([]byte, error) {
	data, err := mgr.Keeper().GetPOL(ctx)
	if err != nil {
		ctx.Logger().Error("fail to get POL", "error", err)
		return nil, fmt.Errorf("fail to get POL: %w", err)
	}
	polValue, err := polPoolValue(ctx, mgr)
	if err != nil {
		ctx.Logger().Error("fail to fetch POL value", "error", err)
		return nil, fmt.Errorf("fail to fetch POL value: %w", err)
	}
	pnl := data.PnL(polValue)
	result := openapi.POLResponse{
		RuneDeposited:  data.RuneDeposited.String(),
		RuneWithdrawn:  data.RuneWithdrawn.String(),
		Value:          polValue.String(),
		Pnl:            pnl.String(),
		CurrentDeposit: data.CurrentDeposit().String(),
	}

	return jsonify(ctx, result)
}

func queryInboundAddresses(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	active, err := mgr.Keeper().GetAsgardVaultsByStatus(ctx, ActiveVault)
	if err != nil {
		ctx.Logger().Error("fail to get active vaults", "error", err)
		return nil, fmt.Errorf("fail to get active vaults: %w", err)
	}

	var resp []openapi.InboundAddress
	constAccessor := mgr.GetConstants()
	signingTransactionPeriod := constAccessor.GetInt64Value(constants.SigningTransactionPeriod)
	if mgr.Keeper() == nil {
		ctx.Logger().Error("keeper is nil, can't fulfill query")
		return nil, errors.New("keeper is nil, can't fulfill query")
	}
	// select vault that is most secure
	vault := mgr.Keeper().GetMostSecure(ctx, active, signingTransactionPeriod)

	chains := vault.GetChains()

	if len(chains) == 0 {
		chains = common.Chains{common.RuneAsset().Chain}
	}

	isGlobalTradingPaused := mgr.Keeper().IsGlobalTradingHalted(ctx)

	for _, chain := range chains {
		// tx send to thorchain doesn't need an address , thus here skip it
		if chain == common.THORChain {
			continue
		}

		isChainTradingPaused := mgr.Keeper().IsChainTradingHalted(ctx, chain)
		isChainLpPaused := mgr.Keeper().IsLPPaused(ctx, chain)

		vaultAddress, err := vault.PubKey.GetAddress(chain)
		if err != nil {
			ctx.Logger().Error("fail to get address for chain", "error", err)
			return nil, fmt.Errorf("fail to get address for chain: %w", err)
		}
		cc := vault.GetContract(chain)
		gasRate := mgr.GasMgr().GetGasRate(ctx, chain)
		networkFeeInfo, err := mgr.GasMgr().GetNetworkFee(ctx, chain)
		if err != nil {
			ctx.Logger().Error("fail to get network fee info", "error", err)
			return nil, fmt.Errorf("fail to get network fee info: %w", err)
		}

		// because THORNode is using 1e8, while GWei in ETH is in 1e9, thus the minimum THORNode can represent is 10Gwei
		// here convert the gas rate to Gwei , so api user don't need to convert it , make it easier for people to understand
		if chain.IsEVM() {
			gasRate = gasRate.MulUint64(10)
		}

		// Retrieve the outbound fee for the chain's gas asset - fee will be zero if no network fee has been posted/the pool doesn't exist
		outboundFee, _ := mgr.GasMgr().GetAssetOutboundFee(ctx, chain.GetGasAsset(), false)

		addr := openapi.InboundAddress{
			Chain:                wrapString(chain.String()),
			PubKey:               wrapString(vault.PubKey.String()),
			Address:              wrapString(vaultAddress.String()),
			Router:               wrapString(cc.Router.String()),
			Halted:               isGlobalTradingPaused || isChainTradingPaused,
			GlobalTradingPaused:  &isGlobalTradingPaused,
			ChainTradingPaused:   &isChainTradingPaused,
			ChainLpActionsPaused: &isChainLpPaused,
			GasRate:              wrapString(gasRate.String()),
			GasRateUnits:         wrapString(chain.GetGasUnits()),
			OutboundTxSize:       wrapString(cosmos.NewUint(networkFeeInfo.TransactionSize).String()),
			OutboundFee:          wrapString(outboundFee.String()),
			DustThreshold:        wrapString(chain.DustThreshold().String()),
		}

		resp = append(resp, addr)
	}

	return jsonify(ctx, resp)
}

// queryNode return the Node information related to the request node address
// /thorchain/node/{nodeaddress}
func queryNode(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.New("node address not provided")
	}
	nodeAddress := path[0]
	addr, err := cosmos.AccAddressFromBech32(nodeAddress)
	if err != nil {
		return nil, cosmos.ErrUnknownRequest("invalid account address")
	}

	nodeAcc, err := mgr.Keeper().GetNodeAccount(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("fail to get node accounts: %w", err)
	}

	slashPts, err := mgr.Keeper().GetNodeAccountSlashPoints(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("fail to get node slash points: %w", err)
	}
	jail, err := mgr.Keeper().GetNodeAccountJail(ctx, nodeAcc.NodeAddress)
	if err != nil {
		return nil, fmt.Errorf("fail to get node jail: %w", err)
	}

	bp, err := mgr.Keeper().GetBondProviders(ctx, nodeAcc.NodeAddress)
	if err != nil {
		return nil, fmt.Errorf("fail to get bond providers: %w", err)
	}
	bp.Adjust(mgr.GetVersion(), nodeAcc.Bond)

	active, err := mgr.Keeper().ListActiveValidators(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get all active node account: %w", err)
	}

	result := openapi.Node{
		NodeAddress: nodeAcc.NodeAddress.String(),
		Status:      nodeAcc.Status.String(),
		PubKeySet: openapi.NodePubKeySet{
			Secp256k1: wrapString(nodeAcc.PubKeySet.Secp256k1.String()),
			Ed25519:   wrapString(nodeAcc.PubKeySet.Ed25519.String()),
		},
		ValidatorConsPubKey: nodeAcc.ValidatorConsPubKey,
		ActiveBlockHeight:   nodeAcc.ActiveBlockHeight,
		StatusSince:         nodeAcc.StatusSince,
		NodeOperatorAddress: nodeAcc.BondAddress.String(),
		TotalBond:           nodeAcc.Bond.String(),
		SignerMembership:    nodeAcc.GetSignerMembership().Strings(),
		RequestedToLeave:    nodeAcc.RequestedToLeave,
		ForcedToLeave:       nodeAcc.ForcedToLeave,
		LeaveHeight:         int64(nodeAcc.LeaveScore), // OpenAPI can only represent uint64 as int64
		IpAddress:           nodeAcc.IPAddress,
		Version:             nodeAcc.GetVersion().String(),
		CurrentAward:        cosmos.ZeroUint().String(), // Default display for if not overwritten.
	}
	result.PeerId = getPeerIDFromPubKey(nodeAcc.PubKeySet.Secp256k1)
	result.SlashPoints = slashPts

	result.Jail = openapi.NodeJail{
		// Since redundant, leave out the node address
		ReleaseHeight: wrapInt64(jail.ReleaseHeight),
		Reason:        wrapString(jail.Reason),
	}

	var providers []openapi.NodeBondProvider
	// Leave this nil (null rather than []) if the source is nil.
	if bp.Providers != nil {
		providers = make([]openapi.NodeBondProvider, len(bp.Providers))
		for i := range bp.Providers {
			providers[i].BondAddress = wrapString(bp.Providers[i].BondAddress.String())
			providers[i].Bond = wrapString(bp.Providers[i].Bond.String())
		}
	}

	result.BondProviders = openapi.NodeBondProviders{
		// Since redundant, leave out the node address
		NodeOperatorFee: bp.NodeOperatorFee.String(),
		Providers:       providers,
	}

	// CurrentAward is an estimation of reward for node in active status
	// Node in other status should not have current reward
	if nodeAcc.Status == NodeActive && !nodeAcc.Bond.IsZero() {
		network, err := mgr.Keeper().GetNetwork(ctx)
		if err != nil {
			return nil, fmt.Errorf("fail to get network: %w", err)
		}
		vaults, err := mgr.Keeper().GetAsgardVaultsByStatus(ctx, ActiveVault)
		if err != nil {
			return nil, fmt.Errorf("fail to get active vaults: %w", err)
		}
		if len(vaults) == 0 {
			return nil, fmt.Errorf("no active vaults")
		}

		totalEffectiveBond, bondHardCap := getTotalEffectiveBond(active)

		// Note that unlike actual BondRewardRune distribution in manager_validator_current.go ,
		// this estimate treats lastChurnHeight as the block_height of the first (oldest) Asgard vault,
		// rather than the active_block_height of the youngest active node.
		// As an example, note from the below URLs that these are 5293728 and 5293733 respectively in block 5336942.
		// https://thornode.ninerealms.com/thorchain/vaults/asgard?height=5336942
		// https://thornode.ninerealms.com/thorchain/nodes?height=5336942
		// (Nodes .cxmy and .uy3a .)
		lastChurnHeight := vaults[0].BlockHeight

		reward, err := getNodeCurrentRewards(ctx, mgr, nodeAcc, lastChurnHeight, network.BondRewardRune, totalEffectiveBond, bondHardCap)
		if err != nil {
			return nil, fmt.Errorf("fail to get current node rewards: %w", err)
		}

		result.CurrentAward = reward.String()
	}

	// TODO: Represent this map as the field directly, instead of making an array?
	// It would then always be represented in alphabetical order.
	chainHeights, err := mgr.Keeper().GetLastObserveHeight(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("fail to get last observe chain height: %w", err)
	}
	// analyze-ignore(map-iteration)
	for c, h := range chainHeights {
		result.ObserveChains = append(result.ObserveChains, openapi.ChainHeight{
			Chain:  c.String(),
			Height: h,
		})
	}

	preflightCheckResult, err := getNodePreflightResult(ctx, mgr, nodeAcc)
	if err != nil {
		ctx.Logger().Error("fail to get node preflight result", "error", err)
	} else {
		result.PreflightStatus = preflightCheckResult
	}
	return jsonify(ctx, result)
}

func getNodePreflightResult(ctx cosmos.Context, mgr *Mgrs, nodeAcc NodeAccount) (openapi.NodePreflightStatus, error) {
	constAccessor := mgr.GetConstants()
	preflightResult := openapi.NodePreflightStatus{}
	status, err := mgr.ValidatorMgr().NodeAccountPreflightCheck(ctx, nodeAcc, constAccessor)
	preflightResult.Status = status.String()
	if err != nil {
		preflightResult.Reason = err.Error()
		preflightResult.Code = 1
	} else {
		preflightResult.Reason = "OK"
		preflightResult.Code = 0
	}
	return preflightResult, nil
}

// Estimates current rewards for the NodeAccount taking into account bond-weighted rewards and slash points
func getNodeCurrentRewards(ctx cosmos.Context, mgr *Mgrs, nodeAcc NodeAccount, lastChurnHeight int64, totalBondReward, totalEffectiveBond, bondHardCap cosmos.Uint) (cosmos.Uint, error) {
	slashPts, err := mgr.Keeper().GetNodeAccountSlashPoints(ctx, nodeAcc.NodeAddress)
	if err != nil {
		return cosmos.ZeroUint(), fmt.Errorf("fail to get node slash points: %w", err)
	}

	// Find number of blocks since the last churn (the last bond reward payout)
	totalActiveBlocks := ctx.BlockHeight() - lastChurnHeight

	// find number of blocks they were well behaved (ie active - slash points)
	earnedBlocks := totalActiveBlocks - slashPts
	if earnedBlocks < 0 {
		earnedBlocks = 0
	}

	naEffectiveBond := nodeAcc.Bond
	if naEffectiveBond.GT(bondHardCap) {
		naEffectiveBond = bondHardCap
	}

	// reward = totalBondReward * (naEffectiveBond / totalEffectiveBond) * (unslashed blocks since last churn / blocks since last churn)
	reward := common.GetUncappedShare(naEffectiveBond, totalEffectiveBond, totalBondReward)
	reward = common.GetUncappedShare(cosmos.NewUint(uint64(earnedBlocks)), cosmos.NewUint(uint64(totalActiveBlocks)), reward)
	return reward, nil
}

// queryNodes return all the nodes that has bond
// /thorchain/nodes
func queryNodes(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	nodeAccounts, err := mgr.Keeper().ListValidatorsWithBond(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get node accounts: %w", err)
	}

	active, err := mgr.Keeper().ListActiveValidators(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get all active node account: %w", err)
	}

	network, err := mgr.Keeper().GetNetwork(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get network: %w", err)
	}

	vaults, err := mgr.Keeper().GetAsgardVaultsByStatus(ctx, ActiveVault)
	if err != nil {
		return nil, fmt.Errorf("fail to get active vaults: %w", err)
	}
	if len(vaults) == 0 {
		return nil, fmt.Errorf("no active vaults")
	}

	totalEffectiveBond, bondHardCap := getTotalEffectiveBond(active)

	lastChurnHeight := vaults[0].BlockHeight
	version := mgr.GetVersion()
	result := make([]openapi.Node, len(nodeAccounts))
	for i, na := range nodeAccounts {
		if na.RequestedToLeave && na.Bond.LTE(cosmos.NewUint(common.One)) {
			// ignore the node , it left and also has very little bond
			// Set the default display for fields which would otherwise be "".
			result[i] = openapi.Node{
				Status:          types.NodeStatus_Unknown.String(),
				TotalBond:       cosmos.ZeroUint().String(),
				BondProviders:   openapi.NodeBondProviders{NodeOperatorFee: cosmos.ZeroUint().String()},
				Version:         semver.MustParse("0.0.0").String(),
				CurrentAward:    cosmos.ZeroUint().String(),
				PreflightStatus: openapi.NodePreflightStatus{Status: types.NodeStatus_Unknown.String()},
			}
			continue
		}

		slashPts, err := mgr.Keeper().GetNodeAccountSlashPoints(ctx, na.NodeAddress)
		if err != nil {
			return nil, fmt.Errorf("fail to get node slash points: %w", err)
		}

		result[i] = openapi.Node{
			NodeAddress: na.NodeAddress.String(),
			Status:      na.Status.String(),
			PubKeySet: openapi.NodePubKeySet{
				Secp256k1: wrapString(na.PubKeySet.Secp256k1.String()),
				Ed25519:   wrapString(na.PubKeySet.Ed25519.String()),
			},
			ValidatorConsPubKey: na.ValidatorConsPubKey,
			ActiveBlockHeight:   na.ActiveBlockHeight,
			StatusSince:         na.StatusSince,
			NodeOperatorAddress: na.BondAddress.String(),
			TotalBond:           na.Bond.String(),
			SignerMembership:    na.GetSignerMembership().Strings(),
			RequestedToLeave:    na.RequestedToLeave,
			ForcedToLeave:       na.ForcedToLeave,
			LeaveHeight:         int64(na.LeaveScore), // OpenAPI can only represent uint64 as int64
			IpAddress:           na.IPAddress,
			Version:             na.GetVersion().String(),
			CurrentAward:        cosmos.ZeroUint().String(), // Default display for if not overwritten.
		}
		result[i].PeerId = getPeerIDFromPubKey(na.PubKeySet.Secp256k1)
		result[i].SlashPoints = slashPts
		if na.Status == NodeActive {
			reward, err := getNodeCurrentRewards(ctx, mgr, na, lastChurnHeight, network.BondRewardRune, totalEffectiveBond, bondHardCap)
			if err != nil {
				return nil, fmt.Errorf("fail to get current node rewards: %w", err)
			}

			result[i].CurrentAward = reward.String()
		}

		jail, err := mgr.Keeper().GetNodeAccountJail(ctx, na.NodeAddress)
		if err != nil {
			return nil, fmt.Errorf("fail to get node jail: %w", err)
		}
		result[i].Jail = openapi.NodeJail{
			// Since redundant, leave out the node address
			ReleaseHeight: wrapInt64(jail.ReleaseHeight),
			Reason:        wrapString(jail.Reason),
		}

		// TODO: Represent this map as the field directly, instead of making an array?
		// It would then always be represented in alphabetical order.
		chainHeights, err := mgr.Keeper().GetLastObserveHeight(ctx, na.NodeAddress)
		if err != nil {
			return nil, fmt.Errorf("fail to get last observe chain height: %w", err)
		}
		// analyze-ignore(map-iteration)
		for c, h := range chainHeights {
			result[i].ObserveChains = append(result[i].ObserveChains, openapi.ChainHeight{
				Chain:  c.String(),
				Height: h,
			})
		}

		preflightCheckResult, err := getNodePreflightResult(ctx, mgr, na)
		if err != nil {
			ctx.Logger().Error("fail to get node preflight result", "error", err)
		} else {
			result[i].PreflightStatus = preflightCheckResult
		}

		bp, err := mgr.Keeper().GetBondProviders(ctx, na.NodeAddress)
		if err != nil {
			ctx.Logger().Error("fail to get bond providers", "error", err)
		}
		bp.Adjust(version, na.Bond)

		var providers []openapi.NodeBondProvider
		// Leave this nil (null rather than []) if the source is nil.
		if bp.Providers != nil {
			providers = make([]openapi.NodeBondProvider, len(bp.Providers))
			for i := range bp.Providers {
				providers[i].BondAddress = wrapString(bp.Providers[i].BondAddress.String())
				providers[i].Bond = wrapString(bp.Providers[i].Bond.String())
			}
		}

		result[i].BondProviders = openapi.NodeBondProviders{
			// Since redundant, leave out the node address
			NodeOperatorFee: bp.NodeOperatorFee.String(),
			Providers:       providers,
		}
	}

	return jsonify(ctx, result)
}

// queryBorrowers
func queryBorrowers(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.New("asset not provided")
	}
	asset, err := common.NewAsset(path[0])
	if err != nil {
		ctx.Logger().Error("fail to get parse asset", "error", err)
		return nil, fmt.Errorf("fail to parse asset: %w", err)
	}

	var loans Loans
	iterator := mgr.Keeper().GetLoanIterator(ctx, asset)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var loan Loan
		mgr.Keeper().Cdc().MustUnmarshal(iterator.Value(), &loan)
		if loan.CollateralDeposited.Equal(loan.CollateralWithdrawn) && loan.DebtIssued.Equal(loan.DebtRepaid) {
			continue
		}
		loans = append(loans, loan)
	}

	borrowers := make([]openapi.Borrower, len(loans))
	for i, loan := range loans {
		borrower := openapi.NewBorrower(
			loan.Owner.String(),
			loan.Asset.String(),
			loan.DebtIssued.String(),
			loan.DebtRepaid.String(),
			loan.Debt().String(),
			loan.CollateralDeposited.String(),
			loan.CollateralWithdrawn.String(),
			loan.Collateral().String(),
			loan.LastOpenHeight,
			loan.LastRepayHeight,
		)
		borrowers[i] = *borrower
	}

	return jsonify(ctx, borrowers)
}

// queryBorrower
func queryBorrower(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	if len(path) < 2 {
		return nil, errors.New("asset/loan not provided")
	}
	asset, err := common.NewAsset(path[0])
	if err != nil {
		ctx.Logger().Error("fail to get parse asset", "error", err)
		return nil, fmt.Errorf("fail to parse asset: %w", err)
	}

	addr, err := common.NewAddress(path[1])
	if err != nil {
		ctx.Logger().Error("fail to get parse address", "error", err)
		return nil, fmt.Errorf("fail to parse address: %w", err)
	}

	loan, err := mgr.Keeper().GetLoan(ctx, asset, addr)
	if err != nil {
		ctx.Logger().Error("fail to get borrower", "error", err)
		return nil, fmt.Errorf("fail to borrower: %w", err)
	}

	borrower := openapi.NewBorrower(
		loan.Owner.String(),
		loan.Asset.String(),
		loan.DebtIssued.String(),
		loan.DebtRepaid.String(),
		loan.Debt().String(),
		loan.CollateralDeposited.String(),
		loan.CollateralWithdrawn.String(),
		loan.Collateral().String(),
		loan.LastOpenHeight,
		loan.LastRepayHeight,
	)

	return jsonify(ctx, borrower)
}

func newSaver(lp LiquidityProvider, pool Pool) openapi.Saver {
	assetRedeemableValue := lp.GetSaversAssetRedeemValue(pool)

	gp := cosmos.NewDec(0)
	if !lp.AssetDepositValue.IsZero() {
		adv := cosmos.NewDec(lp.AssetDepositValue.BigInt().Int64())
		arv := cosmos.NewDec(assetRedeemableValue.BigInt().Int64())
		gp = arv.Sub(adv)
		gp = gp.Quo(adv)
	}

	return openapi.Saver{
		Asset:              lp.Asset.GetLayer1Asset().String(),
		AssetAddress:       lp.AssetAddress.String(),
		LastAddHeight:      wrapInt64(lp.LastAddHeight),
		LastWithdrawHeight: wrapInt64(lp.LastWithdrawHeight),
		Units:              lp.Units.String(),
		AssetDepositValue:  lp.AssetDepositValue.String(),
		AssetRedeemValue:   assetRedeemableValue.String(),
		GrowthPct:          gp.String(),
	}
}

// queryLiquidityProviders
// isSavers is true if request is for the savers of a Savers Pool, if false the request is for an L1 pool
func queryLiquidityProviders(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs, isSavers bool) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.New("asset not provided")
	}
	if isSavers {
		path[0] = strings.Replace(path[0], ".", "/", 1)
	}
	asset, err := common.NewAsset(path[0])
	if err != nil {
		ctx.Logger().Error("fail to get parse asset", "error", err)
		return nil, fmt.Errorf("fail to parse asset: %w", err)
	}
	if asset.IsDerivedAsset() {
		return nil, fmt.Errorf("must not be a derived asset")
	}
	if isSavers && !asset.IsVaultAsset() {
		return nil, fmt.Errorf("invalid request: requested pool is not a SaversPool")
	} else if !isSavers && asset.IsVaultAsset() {
		return nil, fmt.Errorf("invalid request: requested pool is a SaversPool")
	}

	poolAsset := asset
	if isSavers {
		poolAsset = asset.GetSyntheticAsset()
	}

	pool, err := mgr.Keeper().GetPool(ctx, poolAsset)
	if err != nil {
		ctx.Logger().Error("fail to get pool", "error", err)
		return nil, fmt.Errorf("fail to get pool: %w", err)
	}

	var lps []openapi.LiquidityProvider
	var savers []openapi.Saver
	iterator := mgr.Keeper().GetLiquidityProviderIterator(ctx, asset)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var lp LiquidityProvider
		mgr.Keeper().Cdc().MustUnmarshal(iterator.Value(), &lp)
		if !isSavers {
			lps = append(lps, openapi.LiquidityProvider{
				// No redeem or LUVI calculations for the array response.
				Asset:              lp.Asset.GetLayer1Asset().String(),
				RuneAddress:        wrapString(lp.RuneAddress.String()),
				AssetAddress:       wrapString(lp.AssetAddress.String()),
				LastAddHeight:      wrapInt64(lp.LastAddHeight),
				LastWithdrawHeight: wrapInt64(lp.LastWithdrawHeight),
				Units:              lp.Units.String(),
				PendingRune:        lp.PendingRune.String(),
				PendingAsset:       lp.PendingAsset.String(),
				PendingTxId:        wrapString(lp.PendingTxID.String()),
				RuneDepositValue:   lp.RuneDepositValue.String(),
				AssetDepositValue:  lp.AssetDepositValue.String(),
			})
		} else {
			savers = append(savers, newSaver(lp, pool))
		}
	}
	if !isSavers {
		return jsonify(ctx, lps)
	} else {
		return jsonify(ctx, savers)
	}
}

// queryLiquidityProvider
// isSavers is true if request is for the savers of a Savers Pool, if false the request is for an L1 pool
func queryLiquidityProvider(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs, isSavers bool) ([]byte, error) {
	if len(path) < 2 {
		return nil, errors.New("asset/lp not provided")
	}
	if isSavers {
		path[0] = strings.Replace(path[0], ".", "/", 1)
	}
	asset, err := common.NewAsset(path[0])
	if err != nil {
		ctx.Logger().Error("fail to get parse asset", "error", err)
		return nil, fmt.Errorf("fail to parse asset: %w", err)
	}

	if asset.IsDerivedAsset() {
		return nil, fmt.Errorf("must not be a derived asset")
	}

	if isSavers && !asset.IsVaultAsset() {
		return nil, fmt.Errorf("invalid request: requested pool is not a SaversPool")
	} else if !isSavers && asset.IsVaultAsset() {
		return nil, fmt.Errorf("invalid request: requested pool is a SaversPool")
	}

	addr, err := common.NewAddress(path[1])
	if err != nil {
		ctx.Logger().Error("fail to get parse address", "error", err)
		return nil, fmt.Errorf("fail to parse address: %w", err)
	}
	lp, err := mgr.Keeper().GetLiquidityProvider(ctx, asset, addr)
	if err != nil {
		ctx.Logger().Error("fail to get liquidity provider", "error", err)
		return nil, fmt.Errorf("fail to liquidity provider: %w", err)
	}

	poolAsset := asset
	if isSavers {
		poolAsset = asset.GetSyntheticAsset()
	}

	pool, err := mgr.Keeper().GetPool(ctx, poolAsset)
	if err != nil {
		ctx.Logger().Error("fail to get pool", "error", err)
		return nil, fmt.Errorf("fail to get pool: %w", err)
	}

	if !isSavers {
		synthSupply := mgr.Keeper().GetTotalSupply(ctx, poolAsset.GetSyntheticAsset())
		_, runeRedeemValue := lp.GetRuneRedeemValue(mgr.GetVersion(), pool, synthSupply)
		_, assetRedeemValue := lp.GetAssetRedeemValue(mgr.GetVersion(), pool, synthSupply)
		_, luviDepositValue := lp.GetLuviDepositValue(pool)
		_, luviRedeemValue := lp.GetLuviRedeemValue(runeRedeemValue, assetRedeemValue)

		lgp := cosmos.NewDec(0)
		if !luviDepositValue.IsZero() {
			ldv := cosmos.NewDec(luviDepositValue.BigInt().Int64())
			lrv := cosmos.NewDec(luviRedeemValue.BigInt().Int64())
			lgp = lrv.Sub(ldv)
			lgp = lgp.Quo(ldv)
		}

		liqp := openapi.LiquidityProvider{
			Asset:              lp.Asset.GetLayer1Asset().String(),
			RuneAddress:        wrapString(lp.RuneAddress.String()),
			AssetAddress:       wrapString(lp.AssetAddress.String()),
			LastAddHeight:      wrapInt64(lp.LastAddHeight),
			LastWithdrawHeight: wrapInt64(lp.LastWithdrawHeight),
			Units:              lp.Units.String(),
			PendingRune:        lp.PendingRune.String(),
			PendingAsset:       lp.PendingAsset.String(),
			PendingTxId:        wrapString(lp.PendingTxID.String()),
			RuneDepositValue:   lp.RuneDepositValue.String(),
			AssetDepositValue:  lp.AssetDepositValue.String(),
			RuneRedeemValue:    wrapString(runeRedeemValue.String()),
			AssetRedeemValue:   wrapString(assetRedeemValue.String()),
			LuviDepositValue:   wrapString(luviDepositValue.String()),
			LuviRedeemValue:    wrapString(luviRedeemValue.String()),
			LuviGrowthPct:      wrapString(lgp.String()),
		}
		return jsonify(ctx, liqp)
	} else {
		saver := newSaver(lp, pool)
		return jsonify(ctx, saver)
	}
}

func newStreamingSwap(streamingSwap StreamingSwap, msgSwap MsgSwap) openapi.StreamingSwap {
	var sourceAsset common.Asset
	// Leave the source_asset field empty if there is more than a single input Coin.
	if len(msgSwap.Tx.Coins) == 1 {
		sourceAsset = msgSwap.Tx.Coins[0].Asset
	}

	var failedSwaps []int64
	// Leave this nil (null rather than []) if the source is nil.
	if streamingSwap.FailedSwaps != nil {
		failedSwaps = make([]int64, len(streamingSwap.FailedSwaps))
		for i := range streamingSwap.FailedSwaps {
			failedSwaps[i] = int64(streamingSwap.FailedSwaps[i])
		}
	}

	return openapi.StreamingSwap{
		TxId:              wrapString(streamingSwap.TxID.String()),
		Interval:          wrapInt64(int64(streamingSwap.Interval)),
		Quantity:          wrapInt64(int64(streamingSwap.Quantity)),
		Count:             wrapInt64(int64(streamingSwap.Count)),
		LastHeight:        wrapInt64(streamingSwap.LastHeight),
		TradeTarget:       streamingSwap.TradeTarget.String(),
		SourceAsset:       wrapString(sourceAsset.String()),
		TargetAsset:       wrapString(msgSwap.TargetAsset.String()),
		Destination:       wrapString(msgSwap.Destination.String()),
		Deposit:           streamingSwap.Deposit.String(),
		In:                streamingSwap.In.String(),
		Out:               streamingSwap.Out.String(),
		FailedSwaps:       failedSwaps,
		FailedSwapReasons: streamingSwap.FailedSwapReasons,
	}
}

func queryStreamingSwaps(ctx cosmos.Context, mgr *Mgrs) ([]byte, error) {
	var streams []openapi.StreamingSwap
	iter := mgr.Keeper().GetStreamingSwapIterator(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var stream StreamingSwap
		mgr.Keeper().Cdc().MustUnmarshal(iter.Value(), &stream)

		var msgSwap MsgSwap
		// Check up to the first two indices (0 through 1) for the MsgSwap; if not found, leave the fields blank.
		for i := 0; i <= 1; i++ {
			swapQueueItem, err := mgr.Keeper().GetSwapQueueItem(ctx, stream.TxID, i)
			if err != nil {
				// GetSwapQueueItem returns an error if there is no MsgSwap set for that index, a normal occurrence here.
				continue
			}
			if !swapQueueItem.IsStreaming() {
				continue
			}
			// In case there are multiple streaming swaps with the same TxID, check the input amount.
			if len(swapQueueItem.Tx.Coins) == 0 || !swapQueueItem.Tx.Coins[0].Amount.Equal(stream.Deposit) {
				continue
			}
			msgSwap = swapQueueItem
			break
		}

		streams = append(streams, newStreamingSwap(stream, msgSwap))
	}
	return jsonify(ctx, streams)
}

func querySwapperClout(ctx cosmos.Context, path []string, mgr *Mgrs) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.New("address not provided")
	}
	addr, err := common.NewAddress(path[0])
	if err != nil {
		ctx.Logger().Error("fail to parse address", "error", err)
		return nil, fmt.Errorf("could not parse address: %w", err)
	}

	clout, err := mgr.Keeper().GetSwapperClout(ctx, addr)
	if err != nil {
		ctx.Logger().Error("fail to get swapper clout", "error", err)
		return nil, fmt.Errorf("could not get swapper clout: %w", err)
	}

	return jsonify(ctx, clout)
}

func queryStreamingSwap(ctx cosmos.Context, path []string, mgr *Mgrs) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.New("tx id not provided")
	}
	txid, err := common.NewTxID(path[0])
	if err != nil {
		ctx.Logger().Error("fail to parse txid", "error", err)
		return nil, fmt.Errorf("could not parse txid: %w", err)
	}

	streamingSwap, err := mgr.Keeper().GetStreamingSwap(ctx, txid)
	if err != nil {
		ctx.Logger().Error("fail to get streaming swap", "error", err)
		return nil, fmt.Errorf("could not get streaming swap: %w", err)
	}

	var msgSwap MsgSwap
	// Check up to the first two indices (0 through 1) for the MsgSwap; if not found, leave the fields blank.
	for i := 0; i <= 1; i++ {
		swapQueueItem, err := mgr.Keeper().GetSwapQueueItem(ctx, txid, i)
		if err != nil {
			// GetSwapQueueItem returns an error if there is no MsgSwap set for that index, a normal occurrence here.
			continue
		}
		if !swapQueueItem.IsStreaming() {
			continue
		}
		// In case there are multiple streaming swaps with the same TxID, check the input amount.
		if len(swapQueueItem.Tx.Coins) == 0 || !swapQueueItem.Tx.Coins[0].Amount.Equal(streamingSwap.Deposit) {
			continue
		}
		msgSwap = swapQueueItem
		break
	}

	result := newStreamingSwap(streamingSwap, msgSwap)

	return jsonify(ctx, result)
}

func queryPool(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.New("asset not provided")
	}
	asset, err := common.NewAsset(path[0])
	if err != nil {
		ctx.Logger().Error("fail to parse asset", "error", err)
		return nil, fmt.Errorf("could not parse asset: %w", err)
	}

	if asset.IsDerivedAsset() {
		return nil, fmt.Errorf("asset: %s is a derived asset", path[0])
	}

	pool, err := mgr.Keeper().GetPool(ctx, asset)
	if err != nil {
		ctx.Logger().Error("fail to get pool", "error", err)
		return nil, fmt.Errorf("could not get pool: %w", err)
	}
	if pool.IsEmpty() {
		return nil, fmt.Errorf("pool: %s doesn't exist", path[0])
	}

	// Get Savers Vault for this L1 pool if it's a gas asset
	saversAsset := pool.Asset.GetSyntheticAsset()
	saversPool, err := mgr.Keeper().GetPool(ctx, saversAsset)
	if err != nil {
		return nil, fmt.Errorf("fail to unmarshal savers vault: %w", err)
	}

	saversDepth := saversPool.BalanceAsset
	saversUnits := saversPool.LPUnits

	synthSupply := mgr.Keeper().GetTotalSupply(ctx, pool.Asset.GetSyntheticAsset())
	pool.CalcUnits(mgr.GetVersion(), synthSupply)

	synthMintPausedErr := isSynthMintPaused(ctx, mgr, saversAsset, cosmos.ZeroUint())
	synthSupplyRemaining, _ := getSynthSupplyRemainingV102(ctx, mgr, saversAsset)

	totalCollateral, err := mgr.Keeper().GetTotalCollateral(ctx, pool.Asset)
	if err != nil {
		return nil, fmt.Errorf("fail to fetch total loan collateral: %w", err)
	}

	loanHandler := NewLoanOpenHandler(mgr)
	// getPoolCR and GetLoanCollateralRemainingForPool
	// are expected to error for block heights earlier than 12241034
	// from negative MaxRuneSupply, so dropping the errors for both
	// and instead displaying both as zero when unretrievable.
	cr, _ := loanHandler.getPoolCR(ctx, pool, cosmos.OneUint())
	loanCollateralRemaining, _ := loanHandler.GetLoanCollateralRemainingForPool(ctx, pool)

	runeDepth, _, _ := mgr.NetworkMgr().CalcAnchor(ctx, mgr, asset)
	dpool, _ := mgr.Keeper().GetPool(ctx, asset.GetDerivedAsset())
	dbps := common.GetUncappedShare(dpool.BalanceRune, runeDepth, cosmos.NewUint(constants.MaxBasisPts))
	if dpool.Status != PoolAvailable {
		dbps = cosmos.ZeroUint()
	}

	p := openapi.Pool{
		Asset:               pool.Asset.String(),
		ShortCode:           wrapString(pool.Asset.ShortCode()),
		Status:              pool.Status.String(),
		Decimals:            wrapInt64(pool.Decimals),
		PendingInboundAsset: pool.PendingInboundAsset.String(),
		PendingInboundRune:  pool.PendingInboundRune.String(),
		BalanceAsset:        pool.BalanceAsset.String(),
		BalanceRune:         pool.BalanceRune.String(),
		PoolUnits:           pool.GetPoolUnits().String(),
		LPUnits:             pool.LPUnits.String(),
		SynthUnits:          pool.SynthUnits.String(),
	}
	p.SynthSupply = synthSupply.String()
	p.SaversDepth = saversDepth.String()
	p.SaversUnits = saversUnits.String()
	p.SynthMintPaused = (synthMintPausedErr != nil)
	p.SynthSupplyRemaining = synthSupplyRemaining.String()
	p.LoanCollateral = totalCollateral.String()
	p.LoanCollateralRemaining = loanCollateralRemaining.String()
	p.DerivedDepthBps = dbps.String()
	p.LoanCr = cr.String()

	if !pool.BalanceAsset.IsZero() && !pool.BalanceRune.IsZero() {
		dollarsPerRune := mgr.Keeper().DollarsPerRune(ctx)
		p.AssetTorPrice = dollarsPerRune.Mul(pool.BalanceRune).Quo(pool.BalanceAsset).String()
	}

	return jsonify(ctx, p)
}

func queryPools(ctx cosmos.Context, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	dollarsPerRune := mgr.Keeper().DollarsPerRune(ctx)
	pools := make([]openapi.Pool, 0)
	iterator := mgr.Keeper().GetPoolIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		var pool Pool
		if err := mgr.Keeper().Cdc().Unmarshal(iterator.Value(), &pool); err != nil {
			return nil, fmt.Errorf("fail to unmarshal pool: %w", err)
		}
		// ignore pool if no liquidity provider units
		if pool.LPUnits.IsZero() {
			continue
		}

		// Ignore synth asset pool (savers). Info will be on the L1 pool
		if pool.Asset.IsVaultAsset() {
			continue
		}

		// Ignore derived assets (except TOR)
		if pool.Asset.IsDerivedAsset() {
			continue
		}

		// Get Savers Vault
		saversAsset := pool.Asset.GetSyntheticAsset()
		saversPool, err := mgr.Keeper().GetPool(ctx, saversAsset)
		if err != nil {
			return nil, fmt.Errorf("fail to unmarshal savers vault: %w", err)
		}

		saversDepth := saversPool.BalanceAsset
		saversUnits := saversPool.LPUnits

		synthSupply := mgr.Keeper().GetTotalSupply(ctx, pool.Asset.GetSyntheticAsset())
		pool.CalcUnits(mgr.GetVersion(), synthSupply)

		synthMintPausedErr := isSynthMintPaused(ctx, mgr, pool.Asset, cosmos.ZeroUint())
		synthSupplyRemaining, _ := getSynthSupplyRemainingV102(ctx, mgr, pool.Asset)

		totalCollateral, err := mgr.Keeper().GetTotalCollateral(ctx, pool.Asset)
		if err != nil {
			return nil, fmt.Errorf("fail to fetch total loan collateral: %w", err)
		}

		loanHandler := NewLoanOpenHandler(mgr)
		// getPoolCR and GetLoanCollateralRemainingForPool
		// are expected to error for block heights earlier than 12241034
		// from negative MaxRuneSupply, so dropping the errors for both
		// and instead displaying both as zero when unretrievable.
		cr, _ := loanHandler.getPoolCR(ctx, pool, cosmos.OneUint())
		loanCollateralRemaining, _ := loanHandler.GetLoanCollateralRemainingForPool(ctx, pool)

		runeDepth, _, _ := mgr.NetworkMgr().CalcAnchor(ctx, mgr, pool.Asset)
		dpool, _ := mgr.Keeper().GetPool(ctx, pool.Asset.GetDerivedAsset())
		dbps := common.GetUncappedShare(dpool.BalanceRune, runeDepth, cosmos.NewUint(constants.MaxBasisPts))
		if dpool.Status != PoolAvailable {
			dbps = cosmos.ZeroUint()
		}

		p := openapi.Pool{
			Asset:               pool.Asset.String(),
			ShortCode:           wrapString(pool.Asset.ShortCode()),
			Status:              pool.Status.String(),
			Decimals:            wrapInt64(pool.Decimals),
			PendingInboundAsset: pool.PendingInboundAsset.String(),
			PendingInboundRune:  pool.PendingInboundRune.String(),
			BalanceAsset:        pool.BalanceAsset.String(),
			BalanceRune:         pool.BalanceRune.String(),
			PoolUnits:           pool.GetPoolUnits().String(),
			LPUnits:             pool.LPUnits.String(),
			SynthUnits:          pool.SynthUnits.String(),
		}

		p.SynthSupply = synthSupply.String()
		p.SaversDepth = saversDepth.String()
		p.SaversUnits = saversUnits.String()
		p.SynthMintPaused = (synthMintPausedErr != nil)
		p.SynthSupplyRemaining = synthSupplyRemaining.String()
		p.LoanCollateral = totalCollateral.String()
		p.LoanCollateralRemaining = loanCollateralRemaining.String()
		p.DerivedDepthBps = dbps.String()
		p.LoanCr = cr.String()

		if !pool.BalanceAsset.IsZero() && !pool.BalanceRune.IsZero() {
			p.AssetTorPrice = dollarsPerRune.Mul(pool.BalanceRune).Quo(pool.BalanceAsset).String()
		}

		pools = append(pools, p)
	}
	return jsonify(ctx, pools)
}

func queryPoolSlips(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	var assets []common.Asset
	if len(path) > 0 && len(path[0]) > 0 {
		// If an Asset has been specified, return information for just that Asset
		// (even if for instance a Derived Asset to check whether it has values set).
		asset, err := common.NewAsset(path[0])
		if err != nil {
			ctx.Logger().Error("fail to parse asset", "error", err, "asset", path[0])
			return nil, fmt.Errorf("fail to parse asset (%s): %w", path[0], err)
		}
		assets = []common.Asset{asset}
	} else {
		iterator := mgr.Keeper().GetPoolIterator(ctx)
		for ; iterator.Valid(); iterator.Next() {
			var pool Pool
			if err := mgr.Keeper().Cdc().Unmarshal(iterator.Value(), &pool); err != nil {
				return nil, fmt.Errorf("fail to unmarshal pool: %w", err)
			}

			// Display the swap slips of Available-pool Layer 1 assets.
			if pool.Status != PoolAvailable || pool.Asset.IsNative() {
				continue
			}
			assets = append(assets, pool.Asset)
		}
	}

	result := make([]openapi.PoolSlipResponseInner, len(assets))
	for i := range assets {
		result[i].Asset = assets[i].String()

		poolSlip, err := mgr.Keeper().GetPoolSwapSlip(ctx, ctx.BlockHeight(), assets[i])
		if err != nil {
			return nil, fmt.Errorf("fail to get swap slip for asset (%s) height (%d), err:%w", assets[i], ctx.BlockHeight(), err)
		}
		result[i].PoolSlip = poolSlip.Int64()

		rollupCount, err := mgr.Keeper().GetRollupCount(ctx, assets[i])
		if err != nil {
			return nil, fmt.Errorf("fail to get rollup count for asset (%s) height (%d), err:%w", assets[i], ctx.BlockHeight(), err)
		}
		result[i].RollupCount = rollupCount

		longRollup, err := mgr.Keeper().GetLongRollup(ctx, assets[i])
		if err != nil {
			return nil, fmt.Errorf("fail to get long rollup for asset (%s) height (%d), err:%w", assets[i], ctx.BlockHeight(), err)
		}
		result[i].LongRollup = longRollup

		rollup, err := mgr.Keeper().GetCurrentRollup(ctx, assets[i])
		if err != nil {
			return nil, fmt.Errorf("fail to get rollup count for asset (%s) height (%d), err:%w", assets[i], ctx.BlockHeight(), err)
		}
		result[i].Rollup = rollup
	}

	// For performance, only sum the rollup swap slip for comparison
	// when a single asset has been specified.
	if len(assets) == 1 {
		maxAnchorBlocks := mgr.Keeper().GetConfigInt64(ctx, constants.MaxAnchorBlocks)
		var summedRollup int64
		for i := ctx.BlockHeight() - maxAnchorBlocks; i < ctx.BlockHeight(); i++ {
			poolSlip, err := mgr.Keeper().GetPoolSwapSlip(ctx, i, assets[0])
			if err != nil {
				// Log the error, zero the sum, and exit the loop.
				ctx.Logger().Error("fail to get swap slip", "error", err, "asset", assets[0], "height", i)
				summedRollup = 0
				break
			}
			summedRollup += poolSlip.Int64()
		}
		result[0].SummedRollup = &summedRollup
	}

	return jsonify(ctx, result)
}

func queryDerivedPool(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.New("asset not provided")
	}
	asset, err := common.NewAsset(path[0])
	if err != nil {
		ctx.Logger().Error("fail to parse asset", "error", err)
		return nil, fmt.Errorf("could not parse asset: %w", err)
	}

	if !asset.IsDerivedAsset() {
		return nil, fmt.Errorf("asset is not a derived asset: %s", asset)
	}

	// call begin block so the derived depth matches the next block execution state
	_ = mgr.NetworkMgr().BeginBlock(ctx.WithBlockHeight(ctx.BlockHeight()+1), mgr)

	// sum rune depth of anchor pools
	runeDepth := sdk.ZeroUint()
	for _, anchor := range mgr.Keeper().GetAnchors(ctx, asset) {
		aPool, _ := mgr.Keeper().GetPool(ctx, anchor)
		runeDepth = runeDepth.Add(aPool.BalanceRune)
	}

	dpool, _ := mgr.Keeper().GetPool(ctx, asset.GetDerivedAsset())
	dbps := cosmos.ZeroUint()
	if dpool.Status == PoolAvailable {
		dbps = common.GetUncappedShare(dpool.BalanceRune, runeDepth, cosmos.NewUint(constants.MaxBasisPts))
	}

	p := openapi.DerivedPool{
		Asset:        dpool.Asset.String(),
		Status:       dpool.Status.String(),
		Decimals:     wrapInt64(dpool.Decimals),
		BalanceAsset: dpool.BalanceAsset.String(),
		BalanceRune:  dpool.BalanceRune.String(),
	}
	p.DerivedDepthBps = dbps.String()

	return jsonify(ctx, p)
}

func queryDerivedPools(ctx cosmos.Context, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	pools := make([]openapi.DerivedPool, 0)
	iterator := mgr.Keeper().GetPoolIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		var pool Pool
		if err := mgr.Keeper().Cdc().Unmarshal(iterator.Value(), &pool); err != nil {
			return nil, fmt.Errorf("fail to unmarshal pool: %w", err)
		}
		// Ignore derived assets (except TOR)
		if !pool.Asset.IsDerivedAsset() {
			continue
		}

		runeDepth, _, _ := mgr.NetworkMgr().CalcAnchor(ctx, mgr, pool.Asset)
		dpool, _ := mgr.Keeper().GetPool(ctx, pool.Asset.GetDerivedAsset())
		dbps := cosmos.ZeroUint()
		if dpool.Status == PoolAvailable {
			dbps = common.GetUncappedShare(dpool.BalanceRune, runeDepth, cosmos.NewUint(constants.MaxBasisPts))
		}

		p := openapi.DerivedPool{
			Asset:        dpool.Asset.String(),
			Status:       dpool.Status.String(),
			Decimals:     wrapInt64(dpool.Decimals),
			BalanceAsset: dpool.BalanceAsset.String(),
			BalanceRune:  dpool.BalanceRune.String(),
		}
		p.DerivedDepthBps = dbps.String()

		pools = append(pools, p)
	}
	return jsonify(ctx, pools)
}

func queryTradeUnit(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.New("asset not provided")
	}
	asset, err := common.NewAsset(path[0])
	if err != nil {
		ctx.Logger().Error("fail to parse asset", "error", err)
		return nil, fmt.Errorf("could not parse asset: %w", err)
	}

	tu, err := mgr.Keeper().GetTradeUnit(ctx, asset)
	if err != nil {
		ctx.Logger().Error("fail to get trade unit", "error", err)
		return nil, fmt.Errorf("could not get trade unit: %w", err)
	}
	tuResp := openapi.TradeUnitResponse{
		Asset: tu.Asset.String(),
		Units: tu.Units.String(),
		Depth: tu.Depth.String(),
	}
	return jsonify(ctx, tuResp)
}

func queryTradeUnits(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	pools, err := mgr.Keeper().GetPools(ctx)
	if err != nil {
		return nil, errors.New("failed to get pools")
	}
	units := make([]openapi.TradeUnitResponse, 0)
	for _, pool := range pools {
		// skip non-layer1 pools
		if pool.Asset.GetChain().IsTHORChain() {
			continue
		}
		asset := pool.Asset.GetTradeAsset()
		tu, err := mgr.Keeper().GetTradeUnit(ctx, asset)
		if err != nil {
			ctx.Logger().Error("fail to get trade unit", "error", err)
			return nil, fmt.Errorf("could not get trade unit: %w", err)
		}
		tuResp := openapi.TradeUnitResponse{
			Asset: tu.Asset.String(),
			Units: tu.Units.String(),
			Depth: tu.Depth.String(),
		}
		units = append(units, tuResp)
	}

	return jsonify(ctx, units)
}

func queryTradeAccounts(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.New("asset not provided")
	}
	asset, err := common.NewAsset(path[0])
	if err != nil {
		ctx.Logger().Error("fail to parse address", "error", err)
		return nil, fmt.Errorf("could not parse address: %w", err)
	}

	accounts := make([]openapi.TradeAccountResponse, 0)
	iter := mgr.Keeper().GetTradeAccountIterator(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var ta TradeAccount
		if err := mgr.Keeper().Cdc().Unmarshal(iter.Value(), &ta); err != nil {
			continue
		}
		if !ta.Asset.Equals(asset) {
			continue
		}
		if ta.Units.IsZero() {
			continue
		}
		taResp := openapi.TradeAccountResponse{
			Asset:              ta.Asset.String(),
			Units:              ta.Units.String(),
			Owner:              ta.Owner.String(),
			LastAddHeight:      &ta.LastAddHeight,
			LastWithdrawHeight: &ta.LastWithdrawHeight,
		}
		accounts = append(accounts, taResp)
	}

	return jsonify(ctx, accounts)
}

func queryTradeAccount(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.New("address not provided")
	}
	addr, err := cosmos.AccAddressFromBech32(path[0])
	if err != nil {
		ctx.Logger().Error("fail to parse address", "error", err)
		return nil, fmt.Errorf("could not parse address: %w", err)
	}

	accounts := make([]openapi.TradeAccountResponse, 0)
	iter := mgr.Keeper().GetTradeAccountIteratorWithAddress(ctx, addr)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var ta TradeAccount
		if err := mgr.Keeper().Cdc().Unmarshal(iter.Value(), &ta); err != nil {
			continue
		}
		if ta.Units.IsZero() {
			continue
		}

		taResp := openapi.TradeAccountResponse{
			Asset:              ta.Asset.String(),
			Units:              ta.Units.String(),
			Owner:              ta.Owner.String(),
			LastAddHeight:      &ta.LastAddHeight,
			LastWithdrawHeight: &ta.LastWithdrawHeight,
		}
		accounts = append(accounts, taResp)
	}

	return jsonify(ctx, accounts)
}

func extractVoter(ctx cosmos.Context, path []string, mgr *Mgrs) (common.TxID, ObservedTxVoter, error) {
	if len(path) == 0 {
		return "", ObservedTxVoter{}, errors.New("tx id not provided")
	}
	hash, err := common.NewTxID(path[0])
	if err != nil {
		ctx.Logger().Error("fail to parse tx id", "error", err)
		return "", ObservedTxVoter{}, fmt.Errorf("fail to parse tx id: %w", err)
	}
	voter, err := mgr.Keeper().GetObservedTxInVoter(ctx, hash)
	if err != nil {
		ctx.Logger().Error("fail to get observed tx voter", "error", err)
		return "", ObservedTxVoter{}, fmt.Errorf("fail to get observed tx voter: %w", err)
	}
	return hash, voter, nil
}

func queryTxVoters(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	hash, voter, err := extractVoter(ctx, path, mgr)
	if err != nil {
		return nil, err
	}
	// when tx in voter doesn't exist , double check tx out voter
	if len(voter.Txs) == 0 {
		voter, err = mgr.Keeper().GetObservedTxOutVoter(ctx, hash)
		if err != nil {
			return nil, fmt.Errorf("fail to get observed tx out voter: %w", err)
		}
		if len(voter.Txs) == 0 {
			return nil, fmt.Errorf("tx: %s doesn't exist", hash)
		}
	}

	var txs []openapi.ObservedTx
	// Leave this nil (null rather than []) if the source is nil.
	if voter.Txs != nil {
		txs = make([]openapi.ObservedTx, len(voter.Txs))
		for i := range voter.Txs {
			txs[i] = castObservedTx(voter.Txs[i])
		}
	}

	var actions []openapi.TxOutItem
	// Leave this nil (null rather than []) if the source is nil.
	if voter.Actions != nil {
		actions = make([]openapi.TxOutItem, len(voter.Actions))
		for i := range voter.Actions {
			actions[i] = castTxOutItem(voter.Actions[i], 0) // Omitted Height field
		}
	}

	var outTxs []openapi.Tx
	// Leave this nil (null rather than []) if the source is nil.
	if voter.OutTxs != nil {
		outTxs = make([]openapi.Tx, len(voter.OutTxs))
		for i := range voter.OutTxs {
			outTxs[i] = castTx(voter.OutTxs[i])
		}
	}

	result := openapi.TxDetailsResponse{
		TxId:            wrapString(voter.TxID.String()),
		Tx:              castObservedTx(voter.Tx),
		Txs:             txs,
		Actions:         actions,
		OutTxs:          outTxs,
		ConsensusHeight: wrapInt64(voter.Height),
		FinalisedHeight: wrapInt64(voter.FinalisedHeight),
		UpdatedVault:    wrapBool(voter.UpdatedVault),
		Reverted:        wrapBool(voter.Reverted),
		OutboundHeight:  wrapInt64(voter.OutboundHeight),
	}

	return jsonify(ctx, result)
}

// TODO: Remove isSwap and isPending code when SwapFinalised field deprecated.
func checkPending(ctx cosmos.Context, keeper keeper.Keeper, voter ObservedTxVoter) (isSwap, isPending, pending bool, streamingSwap StreamingSwap) {
	// If there's no (confirmation-counting-complete) consensus transaction yet, don't spend time checking the swap status.
	if voter.Tx.IsEmpty() || !voter.Tx.IsFinal() {
		return
	}

	pending = keeper.HasSwapQueueItem(ctx, voter.TxID, 0) || keeper.HasOrderBookItem(ctx, voter.TxID)

	// Only look for streaming information when a swap is pending.
	if pending {
		var err error
		streamingSwap, err = keeper.GetStreamingSwap(ctx, voter.TxID)
		if err != nil {
			// Log the error, but continue without streaming information.
			ctx.Logger().Error("fail to get streaming swap", "error", err)
		}
	}

	memo, err := ParseMemoWithTHORNames(ctx, keeper, voter.Tx.Tx.Memo)
	if err != nil {
		// If unable to parse, assume not a (valid) swap or limit order memo.
		return
	}

	memoType := memo.GetType()
	// If the memo asset is a synth, as with Savers add liquidity or withdraw, a swap is assumed to be involved.
	if memoType == TxSwap || memoType == TxLimitOrder || memo.GetAsset().IsVaultAsset() {
		isSwap = true
		// Only check the KVStore when the inbound transaction has already been finalised
		// and when there haven't been any Actions planned.
		// This will also check the KVStore when an inbound transaction has no output,
		// such as the output being not enough to cover a fee.
		if voter.FinalisedHeight != 0 && len(voter.Actions) == 0 {
			// Use of Swap Queue or Order Book depends on Mimir key EnableOrderBooks rather than memo type, so check both.
			isPending = pending
		}
	}

	return
}

// Get the largest number of signers for a not-final (pre-confirmation-counting) and final Txs respectively.
func countSigners(voter ObservedTxVoter) (*int64, int64) {
	var notFinalCount, finalCount int64
	for i, refTx := range voter.Txs {
		signersMap := make(map[string]bool)
		final := refTx.IsFinal()
		for f, tx := range voter.Txs {
			// Earlier Txs already checked against all, so no need to check,
			// but do include the signers of the current Txs.
			if f < i {
				continue
			}
			// Count larger number of signers for not-final and final observations separately.
			if tx.IsFinal() != final {
				continue
			}
			if !refTx.Tx.EqualsEx(tx.Tx) {
				continue
			}

			for _, signer := range tx.GetSigners() {
				signersMap[signer.String()] = true
			}
		}
		if final && int64(len(signersMap)) > finalCount {
			finalCount = int64(len(signersMap))
		} else if int64(len(signersMap)) > notFinalCount {
			notFinalCount = int64(len(signersMap))
		}
	}
	return wrapInt64(notFinalCount), finalCount
}

// Call newTxStagesResponse from both queryTxStatus (which includes the stages) and queryTxStages.
// TODO: Remove isSwap and isPending arguments when SwapFinalised deprecated in favour of SwapStatus.
// TODO: Deprecate InboundObserved.Started field in favour of the observation counting.
func newTxStagesResponse(ctx cosmos.Context, voter ObservedTxVoter, isSwap, isPending, pending bool, streamingSwap StreamingSwap) (result openapi.TxStagesResponse) {
	result.InboundObserved.PreConfirmationCount, result.InboundObserved.FinalCount = countSigners(voter)
	result.InboundObserved.Completed = !voter.Tx.IsEmpty()

	// If not Completed, fill in Started and do not proceed.
	if !result.InboundObserved.Completed {
		obStart := (len(voter.Txs) != 0)
		result.InboundObserved.Started = &obStart
		return result
	}

	// Current block height is relevant in the confirmation counting and outbound stages.
	currentHeight := ctx.BlockHeight()

	// Only fill in InboundConfirmationCounted when confirmation counting took place.
	if voter.Height != 0 {
		var confCount openapi.InboundConfirmationCountedStage

		// Set the Completed state first.
		extObsHeight := voter.Tx.BlockHeight
		extConfDelayHeight := voter.Tx.FinaliseHeight
		confCount.Completed = !(extConfDelayHeight > extObsHeight)

		// Only fill in other fields if not Completed.
		if !confCount.Completed {
			countStartHeight := voter.Height
			confCount.CountingStartHeight = wrapInt64(countStartHeight)
			confCount.Chain = wrapString(voter.Tx.Tx.Chain.String())
			confCount.ExternalObservedHeight = wrapInt64(extObsHeight)
			confCount.ExternalConfirmationDelayHeight = wrapInt64(extConfDelayHeight)

			estConfMs := voter.Tx.Tx.Chain.ApproximateBlockMilliseconds() * (extConfDelayHeight - extObsHeight)
			if currentHeight > countStartHeight {
				estConfMs -= (currentHeight - countStartHeight) * common.THORChain.ApproximateBlockMilliseconds()
			}
			estConfSec := estConfMs / 1000
			// Floor at 0.
			if estConfSec < 0 {
				estConfSec = 0
			}
			confCount.RemainingConfirmationSeconds = &estConfSec
		}

		result.InboundConfirmationCounted = &confCount
	}

	var inboundFinalised openapi.InboundFinalisedStage
	inboundFinalised.Completed = (voter.FinalisedHeight != 0)
	result.InboundFinalised = &inboundFinalised

	var swapStatus openapi.SwapStatus
	swapStatus.Pending = pending
	// Only display the SwapStatus stage's Streaming field when there's streaming information available.
	if streamingSwap.Valid() == nil {
		streaming := openapi.StreamingStatus{
			Interval: int64(streamingSwap.Interval),
			Quantity: int64(streamingSwap.Quantity),
			Count:    int64(streamingSwap.Count),
		}
		swapStatus.Streaming = &streaming
	}
	result.SwapStatus = &swapStatus

	// Whether there's an external outbound or not, show the SwapFinalised stage from the start.
	if isSwap {
		var swapFinalisedState openapi.SwapFinalisedStage

		swapFinalisedState.Completed = false
		if !isPending && result.InboundFinalised.Completed {
			// Record as completed only when not pending after the inbound has already been finalised.
			swapFinalisedState.Completed = true
		}

		result.SwapFinalised = &swapFinalisedState
	}

	// Only fill ExternalOutboundDelay and ExternalOutboundKeysign for inbound transactions with an external outbound;
	// namely, transactions with an outbound_height .
	if voter.OutboundHeight == 0 {
		return result
	}

	// Only display the OutboundDelay stage when there's a delay.
	if voter.OutboundHeight > voter.FinalisedHeight {
		var outDelay openapi.OutboundDelayStage

		// Set the Completed state first.
		outDelay.Completed = (currentHeight >= voter.OutboundHeight)

		// Only fill in other fields if not Completed.
		if !outDelay.Completed {
			remainBlocks := voter.OutboundHeight - currentHeight
			outDelay.RemainingDelayBlocks = &remainBlocks

			remainSec := remainBlocks * common.THORChain.ApproximateBlockMilliseconds() / 1000
			outDelay.RemainingDelaySeconds = &remainSec
		}

		result.OutboundDelay = &outDelay
	}

	var outSigned openapi.OutboundSignedStage

	// Set the Completed state first.
	outSigned.Completed = (voter.Tx.Status != types.Status_incomplete)

	// Only fill in other fields if not Completed.
	if !outSigned.Completed {
		scheduledHeight := voter.OutboundHeight
		outSigned.ScheduledOutboundHeight = &scheduledHeight

		// Only fill in BlocksSinceScheduled if the outbound delay is complete.
		if currentHeight >= scheduledHeight {
			sinceScheduled := currentHeight - scheduledHeight
			outSigned.BlocksSinceScheduled = &sinceScheduled
		}
	}

	result.OutboundSigned = &outSigned

	return result
}

func queryTxStages(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	// First, get the ObservedTxVoter of interest.
	_, voter, err := extractVoter(ctx, path, mgr)
	if err != nil {
		return nil, err
	}
	// when no TxIn voter don't check TxOut voter, as TxOut THORChain observation or not matters little to the user once signed and broadcast
	// Rather than a "tx: %s doesn't exist" result, allow a response to an existing-but-unobserved hash with Observation.Started 'false'.

	isSwap, isPending, pending, streamingSwap := checkPending(ctx, mgr.Keeper(), voter)

	result := newTxStagesResponse(ctx, voter, isSwap, isPending, pending, streamingSwap)

	return jsonify(ctx, result)
}

func queryTxStatus(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	// First, get the ObservedTxVoter of interest.
	_, voter, err := extractVoter(ctx, path, mgr)
	if err != nil {
		return nil, err
	}
	// when no TxIn voter don't check TxOut voter, as TxOut THORChain observation or not matters little to the user once signed and broadcast
	// Rather than a "tx: %s doesn't exist" result, allow a response to an existing-but-unobserved hash with Stages.Observation.Started 'false'.

	// TODO: Remove isSwap and isPending arguments when SwapFinalised deprecated.
	isSwap, isPending, pending, streamingSwap := checkPending(ctx, mgr.Keeper(), voter)

	var result openapi.TxStatusResponse

	// If there's a consensus Tx, display that.
	// If not, but there's at least one observation, display the first observation's Tx.
	// If there are no observations yet, don't display a Tx (only showing the 'Observation' stage with 'Started' false).
	if !voter.Tx.Tx.IsEmpty() {
		tx := castTx(voter.Tx.Tx)
		result.Tx = &tx
	} else if len(voter.Txs) > 0 {
		tx := castTx(voter.Txs[0].Tx)
		result.Tx = &tx
	}

	// Leave this nil (null rather than []) if the source is nil.
	if voter.Actions != nil {
		result.PlannedOutTxs = make([]openapi.PlannedOutTx, len(voter.Actions))
		for i := range voter.Actions {
			result.PlannedOutTxs[i] = openapi.PlannedOutTx{
				Chain:     voter.Actions[i].Chain.String(),
				ToAddress: voter.Actions[i].ToAddress.String(),
				Coin:      castCoin(voter.Actions[i].Coin),
				Refund:    strings.HasPrefix(voter.Actions[i].Memo, "REFUND"),
			}
		}
	}

	// Leave this nil (null rather than []) if the source is nil.
	if voter.OutTxs != nil {
		result.OutTxs = make([]openapi.Tx, len(voter.OutTxs))
		for i := range voter.OutTxs {
			result.OutTxs[i] = castTx(voter.OutTxs[i])
		}
	}

	result.Stages = newTxStagesResponse(ctx, voter, isSwap, isPending, pending, streamingSwap)

	return jsonify(ctx, result)
}

func queryTx(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	hash, voter, err := extractVoter(ctx, path, mgr)
	if err != nil {
		return nil, err
	}
	if len(voter.Txs) == 0 {
		voter, err = mgr.Keeper().GetObservedTxOutVoter(ctx, hash)
		if err != nil {
			return nil, fmt.Errorf("fail to get observed tx out voter: %w", err)
		}
		if len(voter.Txs) == 0 {
			return nil, fmt.Errorf("tx: %s doesn't exist", hash)
		}
	}

	nodeAccounts, err := mgr.Keeper().ListActiveValidators(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get node accounts: %w", err)
	}
	keysignMetric, err := mgr.Keeper().GetTssKeysignMetric(ctx, hash)
	if err != nil {
		ctx.Logger().Error("fail to get keysign metrics", "error", err)
	}
	result := struct {
		ObservedTx      openapi.ObservedTx     `json:"observed_tx"`
		ConsensusHeight int64                  `json:"consensus_height,omitempty"`
		FinalisedHeight int64                  `json:"finalised_height,omitempty"`
		OutboundHeight  int64                  `json:"outbound_height,omitempty"`
		KeysignMetrics  types.TssKeysignMetric `json:"keysign_metric"`
	}{
		ObservedTx:      castObservedTx(voter.GetTx(nodeAccounts)),
		ConsensusHeight: voter.Height,
		FinalisedHeight: voter.FinalisedHeight,
		OutboundHeight:  voter.OutboundHeight,
		KeysignMetrics:  *keysignMetric,
	}
	return jsonify(ctx, result)
}

func extractBlockHeight(ctx cosmos.Context, path []string) (int64, error) {
	if len(path) == 0 {
		return -1, errors.New("block height not provided")
	}
	height, err := strconv.ParseInt(path[0], 0, 64)
	if err != nil {
		ctx.Logger().Error("fail to parse block height", "error", err)
		return -1, fmt.Errorf("fail to parse block height: %w", err)
	}
	if height > ctx.BlockHeight() {
		return -1, fmt.Errorf("block height not available yet")
	}
	return height, nil
}

func queryKeygen(ctx cosmos.Context, kbs cosmos.KeybaseStore, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	height, err := extractBlockHeight(ctx, path)
	if err != nil {
		return nil, err
	}

	keygenBlock, err := mgr.Keeper().GetKeygenBlock(ctx, height)
	if err != nil {
		ctx.Logger().Error("fail to get keygen block", "error", err)
		return nil, fmt.Errorf("fail to get keygen block: %w", err)
	}

	if len(path) > 1 {
		pk, err := common.NewPubKey(path[1])
		if err != nil {
			ctx.Logger().Error("fail to parse pubkey", "error", err)
			return nil, fmt.Errorf("fail to parse pubkey: %w", err)
		}
		// only return those keygen contains the request pub key
		newKeygenBlock := NewKeygenBlock(keygenBlock.Height)
		for _, keygen := range keygenBlock.Keygens {
			if keygen.GetMembers().Contains(pk) {
				newKeygenBlock.Keygens = append(newKeygenBlock.Keygens, keygen)
			}
		}
		keygenBlock = newKeygenBlock
	}

	buf, err := json.Marshal(keygenBlock)
	if err != nil {
		ctx.Logger().Error("fail to marshal keygen block to json", "error", err)
		return nil, fmt.Errorf("fail to marshal keygen block to json: %w", err)
	}
	sig, _, err := kbs.Keybase.Sign("thorchain", buf)
	if err != nil {
		ctx.Logger().Error("fail to sign keygen", "error", err)
		return nil, fmt.Errorf("fail to sign keygen: %w", err)
	}

	var keygens []openapi.Keygen
	// Leave this nil (null rather than []) if the source is nil.
	if keygenBlock.Keygens != nil {
		keygens = make([]openapi.Keygen, len(keygenBlock.Keygens))
		for i := range keygenBlock.Keygens {
			keygens[i] = openapi.Keygen{
				Id:      wrapString(keygenBlock.Keygens[i].ID.String()),
				Type:    wrapString(keygenBlock.Keygens[i].Type.String()),
				Members: keygenBlock.Keygens[i].Members,
			}
		}
	}

	query := openapi.KeygenResponse{
		KeygenBlock: openapi.KeygenBlock{
			Height:  wrapInt64(keygenBlock.Height),
			Keygens: keygens,
		},
		Signature: base64.StdEncoding.EncodeToString(sig),
	}

	return jsonify(ctx, query)
}

func queryKeysign(ctx cosmos.Context, kbs cosmos.KeybaseStore, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	height, err := extractBlockHeight(ctx, path)
	if err != nil {
		return nil, err
	}

	pk := common.EmptyPubKey
	if len(path) > 1 {
		pk, err = common.NewPubKey(path[1])
		if err != nil {
			ctx.Logger().Error("fail to parse pubkey", "error", err)
			return nil, fmt.Errorf("fail to parse pubkey: %w", err)
		}
	}

	txs, err := mgr.Keeper().GetTxOut(ctx, height)
	if err != nil {
		ctx.Logger().Error("fail to get tx out array from key value store", "error", err)
		return nil, fmt.Errorf("fail to get tx out array from key value store: %w", err)
	}

	if !pk.IsEmpty() {
		newTxs := &TxOut{
			Height: txs.Height,
		}
		for _, tx := range txs.TxArray {
			if pk.Equals(tx.VaultPubKey) {
				zero := cosmos.ZeroUint()
				if tx.CloutSpent == nil {
					tx.CloutSpent = &zero
				}
				newTxs.TxArray = append(newTxs.TxArray, tx)
			}
		}
		txs = newTxs
	}

	buf, err := json.Marshal(txs)
	if err != nil {
		ctx.Logger().Error("fail to marshal keysign block to json", "error", err)
		return nil, fmt.Errorf("fail to marshal keysign block to json: %w", err)
	}
	sig, _, err := kbs.Keybase.Sign("thorchain", buf)
	if err != nil {
		ctx.Logger().Error("fail to sign keysign", "error", err)
		return nil, fmt.Errorf("fail to sign keysign: %w", err)
	}

	// TODO: use openapi type after Bifrost uses the same so signatures match.
	type QueryKeysign struct {
		Keysign   TxOut  `json:"keysign"`
		Signature string `json:"signature"`
	}

	query := QueryKeysign{
		Keysign:   *txs,
		Signature: base64.StdEncoding.EncodeToString(sig),
	}

	return jsonify(ctx, query)
}

// queryOutQueue - iterates over txout, counting how many transactions are waiting to be sent
func queryQueue(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	constAccessor := mgr.GetConstants()
	signingTransactionPeriod := constAccessor.GetInt64Value(constants.SigningTransactionPeriod)
	startHeight := ctx.BlockHeight() - signingTransactionPeriod
	var query openapi.QueueResponse
	scheduledOutboundValue := cosmos.ZeroUint()
	scheduledOutboundClout := cosmos.ZeroUint()

	iterator := mgr.Keeper().GetSwapQueueIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var msg MsgSwap
		if err := mgr.Keeper().Cdc().Unmarshal(iterator.Value(), &msg); err != nil {
			continue
		}
		query.Swap++
	}

	iter2 := mgr.Keeper().GetOrderBookItemIterator(ctx)
	defer iter2.Close()
	for ; iter2.Valid(); iter2.Next() {
		var msg MsgSwap
		if err := mgr.Keeper().Cdc().Unmarshal(iterator.Value(), &msg); err != nil {
			ctx.Logger().Error("failed to load MsgSwap", "error", err)
			continue
		}
		query.Swap++
	}

	for height := startHeight; height <= ctx.BlockHeight(); height++ {
		txs, err := mgr.Keeper().GetTxOut(ctx, height)
		if err != nil {
			ctx.Logger().Error("fail to get tx out array from key value store", "error", err)
			return nil, fmt.Errorf("fail to get tx out array from key value store: %w", err)
		}
		for _, tx := range txs.TxArray {
			if tx.OutHash.IsEmpty() {
				memo, _ := ParseMemoWithTHORNames(ctx, mgr.Keeper(), tx.Memo)
				if memo.IsInternal() {
					query.Internal++
				} else if memo.IsOutbound() {
					query.Outbound++
				}
			}
		}
	}

	// sum outbound value
	maxTxOutOffset, err := mgr.Keeper().GetMimir(ctx, constants.MaxTxOutOffset.String())
	if maxTxOutOffset < 0 || err != nil {
		maxTxOutOffset = constAccessor.GetInt64Value(constants.MaxTxOutOffset)
	}
	txOutDelayMax, err := mgr.Keeper().GetMimir(ctx, constants.TxOutDelayMax.String())
	if txOutDelayMax <= 0 || err != nil {
		txOutDelayMax = constAccessor.GetInt64Value(constants.TxOutDelayMax)
	}

	for height := ctx.BlockHeight() + 1; height <= ctx.BlockHeight()+txOutDelayMax; height++ {
		value, clout, err := mgr.Keeper().GetTxOutValue(ctx, height)
		if err != nil {
			ctx.Logger().Error("fail to get tx out array from key value store", "error", err)
			continue
		}
		if height > ctx.BlockHeight()+maxTxOutOffset && value.IsZero() {
			// we've hit our max offset, and an empty block, we can assume the
			// rest will be empty as well
			break
		}
		scheduledOutboundValue = scheduledOutboundValue.Add(value)
		scheduledOutboundClout = scheduledOutboundClout.Add(clout)
	}

	query.ScheduledOutboundValue = scheduledOutboundValue.String()
	query.ScheduledOutboundClout = scheduledOutboundClout.String()

	return jsonify(ctx, query)
}

func queryLastBlockHeights(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	var chains common.Chains
	if len(path) > 0 && len(path[0]) > 0 {
		var err error
		chain, err := common.NewChain(path[0])
		if err != nil {
			ctx.Logger().Error("fail to parse chain", "error", err, "chain", path[0])
			return nil, fmt.Errorf("fail to retrieve chain: %w", err)
		}
		chains = append(chains, chain)
	} else {
		asgards, err := mgr.Keeper().GetAsgardVaultsByStatus(ctx, ActiveVault)
		if err != nil {
			return nil, fmt.Errorf("fail to get active asgard: %w", err)
		}
		for _, vault := range asgards {
			chains = vault.GetChains().Distinct()
			break
		}
	}
	var result []openapi.LastBlock
	for _, c := range chains {
		if c == common.THORChain {
			continue
		}
		chainHeight, err := mgr.Keeper().GetLastChainHeight(ctx, c)
		if err != nil {
			return nil, fmt.Errorf("fail to get last chain height: %w", err)
		}

		signed, err := mgr.Keeper().GetLastSignedHeight(ctx)
		if err != nil {
			return nil, fmt.Errorf("fail to get last sign height: %w", err)
		}
		result = append(result, openapi.LastBlock{
			Chain:          c.String(),
			LastObservedIn: chainHeight,
			LastSignedOut:  signed,
			Thorchain:      ctx.BlockHeight(),
		})
	}

	return jsonify(ctx, result)
}

func queryConstantValues(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	constAccessor := mgr.GetConstants()
	return jsonify(ctx, constAccessor)
}

func queryVersion(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	v, hasV := mgr.Keeper().GetVersionWithCtx(ctx)
	if !hasV {
		// re-compute version if not stored
		v = mgr.Keeper().GetLowestActiveVersion(ctx)
	}

	minJoinLast, minJoinLastChangedHeight := mgr.Keeper().GetMinJoinLast(ctx)

	ver := openapi.VersionResponse{
		Current:         v.String(),
		Next:            minJoinLast.String(),
		NextSinceHeight: wrapInt64(minJoinLastChangedHeight), // omitted if 0
		Querier:         constants.SWVersion.String(),
	}
	return jsonify(ctx, ver)
}

func queryMimirWithKey(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	if len(path) == 0 && len(path[0]) == 0 {
		return nil, fmt.Errorf("no mimir key")
	}

	v, err := mgr.Keeper().GetMimir(ctx, path[0])
	if err != nil {
		return nil, fmt.Errorf("fail to get mimir with key:%s, err : %w", path[0], err)
	}
	return jsonify(ctx, v)
}

func queryMimirValues(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	values := map[string]int64{}

	// collect all keys with set values, not displaying those with votes but no set value
	keeper := mgr.Keeper()
	iter := keeper.GetMimirIterator(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		key := strings.TrimPrefix(string(iter.Key()), "mimir//")
		value, err := keeper.GetMimir(ctx, key)
		if err != nil {
			ctx.Logger().Error("fail to get mimir value", "error", err)
			continue
		}
		if value < 0 {
			ctx.Logger().Error("negative mimir value set", "key", key, "value", value)
			continue
		}
		values[key] = value
	}

	return jsonify(ctx, values)
}

func queryMimirV2IDs(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	mimirsMap := make(map[string]openapi.MimirV2IDsResponse, 0)
	iter := mgr.Keeper().GetNodeMimirIteratorV2(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		nm := NodeMimirs{}
		if err := mgr.Keeper().Cdc().Unmarshal(iter.Value(), &nm); err != nil {
			ctx.Logger().Error("fail to unmarshal node mimir value", "error", err)
			return nil, fmt.Errorf("fail to unmarshal node mimir value: %w", err)
		}

		for _, m := range nm.GetMimirs() {
			parts := strings.Split(m.Key, "-")
			id, err := strconv.Atoi(parts[0])
			if err != nil {
				continue
			}
			ref := parts[len(parts)-1]

			mv2, exists := mimir.GetMimir(mimir.Id(id), ref)
			if !exists {
				continue
			}

			if _, exists := mimirsMap[m.Key]; !exists {
				mimirsMap[m.Key] = openapi.MimirV2IDsResponse{
					Id:        id,
					Name:      mv2.Name(),
					Type:      mv2.Type().String(),
					VoteKey:   fmt.Sprintf("%d-%s", id, ref),
					LegacyKey: mv2.LegacyKey(ref),
					Votes:     make(map[string]int64), // OpenAPI only supports string keys.
				}
			}
			mimirsMap[m.Key].Votes[strconv.FormatInt(m.Value, 10)] += 1
		}
	}
	// jsonify's json.Marshal sorts the map keys alphabetically.
	return jsonify(ctx, mimirsMap)
}

func queryMimirV2Values(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	values := make(map[string]int64)

	active, err := mgr.Keeper().ListActiveValidators(ctx)
	if err != nil {
		ctx.Logger().Error("failed to get active validator set", "error", err)
	}

	iterNode := mgr.Keeper().GetNodeMimirIteratorV2(ctx)
	defer iterNode.Close()
	for ; iterNode.Valid(); iterNode.Next() {
		key := strings.TrimPrefix(string(iterNode.Key()), "nodemimirV2//")
		mimirs, err := mgr.Keeper().GetNodeMimirsV2(ctx, key)
		if err != nil {
			continue
		}
		parts := strings.Split(key, "-")
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		ref := parts[len(parts)-1]

		m, _ := mimir.GetMimir(mimir.Id(id), ref)
		value := int64(-1)
		switch m.Type() {
		case mimir.EconomicMimir:
			value, _ = mgr.Keeper().GetMimirV2(ctx, key)
			if value < 0 {
				value = mimirs.ValueOfEconomic(key, active.GetNodeAddresses())
			}
		case mimir.OperationalMimir:
			value = mimirs.ValueOfOperational(key, constants.MinMimirV2Vote, active.GetNodeAddresses())
		}
		if value >= 0 {
			values[m.Name()] = value
		}
	}

	return jsonify(ctx, values)
}

func queryMimirV2NodesAll(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	mimirs := NodeMimirs{}
	iter := mgr.Keeper().GetNodeMimirIteratorV2(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		m := NodeMimirs{}
		if err := mgr.Keeper().Cdc().Unmarshal(iter.Value(), &m); err != nil {
			ctx.Logger().Error("fail to unmarshal node mimir value", "error", err)
			return nil, fmt.Errorf("fail to unmarshal node mimir value: %w", err)
		}
		mimirs.Mimirs = append(mimirs.Mimirs, m.Mimirs...)
	}

	return jsonify(ctx, mimirs)
}

func queryMimirAdminValues(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	values := make(map[string]int64)
	iter := mgr.Keeper().GetMimirIterator(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		value := types.ProtoInt64{}
		if err := mgr.Keeper().Cdc().Unmarshal(iter.Value(), &value); err != nil {
			ctx.Logger().Error("fail to unmarshal mimir value", "error", err)
			return nil, fmt.Errorf("fail to unmarshal mimir value: %w", err)
		}
		k := strings.TrimPrefix(string(iter.Key()), "mimir//")
		values[k] = value.GetValue()
	}
	return jsonify(ctx, values)
}

func queryMimirNodesAllValues(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	mimirs := NodeMimirs{}
	iter := mgr.Keeper().GetNodeMimirIterator(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		m := NodeMimirs{}
		if err := mgr.Keeper().Cdc().Unmarshal(iter.Value(), &m); err != nil {
			ctx.Logger().Error("fail to unmarshal node mimir value", "error", err)
			return nil, fmt.Errorf("fail to unmarshal node mimir value: %w", err)
		}
		mimirs.Mimirs = append(mimirs.Mimirs, m.Mimirs...)
	}

	return jsonify(ctx, mimirs)
}

func queryMimirNodesValues(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	activeNodes, err := mgr.Keeper().ListActiveValidators(ctx)
	if err != nil {
		ctx.Logger().Error("fail to fetch active node accounts", "error", err)
		return nil, fmt.Errorf("fail to fetch active node accounts: %w", err)
	}
	active := activeNodes.GetNodeAddresses()

	values := make(map[string]int64)
	iter := mgr.Keeper().GetNodeMimirIterator(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		mimirs := NodeMimirs{}
		if err := mgr.Keeper().Cdc().Unmarshal(iter.Value(), &mimirs); err != nil {
			ctx.Logger().Error("fail to unmarshal node mimir value", "error", err)
			return nil, fmt.Errorf("fail to unmarshal node mimir value: %w", err)
		}
		k := strings.TrimPrefix(string(iter.Key()), "nodemimir//")
		if v, ok := mimirs.HasSuperMajority(k, active); ok {
			values[k] = v
		}
	}

	return jsonify(ctx, values)
}

func queryMimirNodeValues(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	acc, err := cosmos.AccAddressFromBech32(path[0])
	if err != nil {
		ctx.Logger().Error("fail to parse thor address", "error", err)
		return nil, fmt.Errorf("fail to parse thor address: %w", err)
	}

	values := make(map[string]int64)
	iter := mgr.Keeper().GetNodeMimirIterator(ctx)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		mimirs := NodeMimirs{}
		if err := mgr.Keeper().Cdc().Unmarshal(iter.Value(), &mimirs); err != nil {
			ctx.Logger().Error("fail to unmarshal node mimir v2 value", "error", err)
			return nil, fmt.Errorf("fail to unmarshal node mimir value: %w", err)
		}

		k := strings.TrimPrefix(string(iter.Key()), "nodemimir//")
		if v, ok := mimirs.Get(k, acc); ok {
			values[k] = v
		}
	}

	return jsonify(ctx, values)
}

func queryOutboundFees(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	var assets []common.Asset

	if len(path) > 0 && len(path[0]) > 0 {
		// If an Asset has been specified, return information for just that Asset
		// (even if for instance a Derived Asset to show its THORChain outbound fee).
		asset, err := common.NewAsset(path[0])
		if err != nil {
			ctx.Logger().Error("fail to parse asset", "error", err, "asset", path[0])
			return nil, fmt.Errorf("fail to parse asset (%s): %w", path[0], err)
		}
		assets = []common.Asset{asset}
	} else {
		// By default display the outbound fees of RUNE and all external-chain Layer 1 assets.
		// Even Staged pool Assets can incur outbound fees (from withdraw outbounds).
		assets = []common.Asset{common.RuneAsset()}
		iterator := mgr.Keeper().GetPoolIterator(ctx)
		for ; iterator.Valid(); iterator.Next() {
			var pool Pool
			if err := mgr.Keeper().Cdc().Unmarshal(iterator.Value(), &pool); err != nil {
				return nil, fmt.Errorf("fail to unmarshal pool: %w", err)
			}

			if pool.Asset.IsNative() {
				// To avoid clutter do not by default display the outbound fees
				// of THORChain Assets other than RUNE.
				continue
			}
			if pool.BalanceAsset.IsZero() || pool.BalanceRune.IsZero() {
				// A Layer 1 Asset's pool must have both depths be non-zero
				// for any outbound fee withholding or gas reimbursement to take place.
				// (This can take place even if the PoolUnits are zero and all liquidity is synths.)
				continue
			}

			assets = append(assets, pool.Asset)
		}
	}

	// Obtain the unchanging CalcOutboundFeeMultiplier arguments before the loop which calls it.
	targetSurplusRune := cosmos.NewUint(uint64(mgr.Keeper().GetConfigInt64(ctx, constants.TargetOutboundFeeSurplusRune)))
	maxMultiplier := cosmos.NewUint(uint64(mgr.Keeper().GetConfigInt64(ctx, constants.MaxOutboundFeeMultiplierBasisPoints)))
	minMultiplier := cosmos.NewUint(uint64(mgr.Keeper().GetConfigInt64(ctx, constants.MinOutboundFeeMultiplierBasisPoints)))

	// Due to the nature of pool iteration by key, this is expected to have RUNE at the top and then be in alphabetical order.
	result := make([]openapi.OutboundFee, 0, len(assets))
	for i := range assets {
		// Display the Asset's fee as the amount of that Asset deducted.
		outboundFee, err := mgr.GasMgr().GetAssetOutboundFee(ctx, assets[i], false)
		if err != nil {
			ctx.Logger().Error("fail to get asset outbound fee", "asset", assets[i], "error", err)
		}

		// Only display fields other than asset and outbound_fee when the Asset is external,
		// as a non-zero dynamic multiplier could be misleading otherwise.
		var outboundFeeWithheldRuneString, outboundFeeSpentRuneString, surplusRuneString, dynamicMultiplierBasisPointsString string
		if !assets[i].IsNative() {
			outboundFeeWithheldRune, err := mgr.Keeper().GetOutboundFeeWithheldRune(ctx, assets[i])
			if err != nil {
				ctx.Logger().Error("fail to get outbound fee withheld rune", "outbound asset", assets[i], "error", err)
				return nil, fmt.Errorf("fail to get outbound fee withheld rune for asset (%s): %w", assets[i], err)
			}
			outboundFeeWithheldRuneString = outboundFeeWithheldRune.String()

			outboundFeeSpentRune, err := mgr.Keeper().GetOutboundFeeSpentRune(ctx, assets[i])
			if err != nil {
				ctx.Logger().Error("fail to get outbound fee spent rune", "outbound asset", assets[i], "error", err)
				return nil, fmt.Errorf("fail to get outbound fee spent rune for asset (%s): %w", assets[i], err)
			}
			outboundFeeSpentRuneString = outboundFeeSpentRune.String()

			surplusRuneString = common.SafeSub(outboundFeeWithheldRune, outboundFeeSpentRune).String()

			dynamicMultiplierBasisPointsString = mgr.GasMgr().CalcOutboundFeeMultiplier(ctx, targetSurplusRune, outboundFeeSpentRune, outboundFeeWithheldRune, maxMultiplier, minMultiplier).String()
		}

		// As the entire endpoint is for outbounds, the term 'Outbound' is omitted from the field names.
		result = append(result, openapi.OutboundFee{
			Asset:                        assets[i].String(),
			OutboundFee:                  outboundFee.String(),
			FeeWithheldRune:              wrapString(outboundFeeWithheldRuneString),
			FeeSpentRune:                 wrapString(outboundFeeSpentRuneString),
			SurplusRune:                  wrapString(surplusRuneString),
			DynamicMultiplierBasisPoints: wrapString(dynamicMultiplierBasisPointsString),
		})

	}

	return jsonify(ctx, result)
}

func queryBan(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	if len(path) == 0 {
		return nil, errors.New("node address not available")
	}
	addr, err := cosmos.AccAddressFromBech32(path[0])
	if err != nil {
		ctx.Logger().Error("invalid node address", "error", err)
		return nil, fmt.Errorf("invalid node address: %w", err)
	}

	ban, err := mgr.Keeper().GetBanVoter(ctx, addr)
	if err != nil {
		ctx.Logger().Error("fail to get ban voter", "error", err)
		return nil, fmt.Errorf("fail to get ban voter: %w", err)
	}

	return jsonify(ctx, ban)
}

func queryScheduledOutbound(ctx cosmos.Context, mgr *Mgrs) ([]byte, error) {
	result := make([]openapi.TxOutItem, 0)
	constAccessor := mgr.GetConstants()
	maxTxOutOffset, err := mgr.Keeper().GetMimir(ctx, constants.MaxTxOutOffset.String())
	if maxTxOutOffset < 0 || err != nil {
		maxTxOutOffset = constAccessor.GetInt64Value(constants.MaxTxOutOffset)
	}
	for height := ctx.BlockHeight() + 1; height <= ctx.BlockHeight()+17280; height++ {
		txOut, err := mgr.Keeper().GetTxOut(ctx, height)
		if err != nil {
			ctx.Logger().Error("fail to get tx out array from key value store", "error", err)
			continue
		}
		if height > ctx.BlockHeight()+maxTxOutOffset && len(txOut.TxArray) == 0 {
			// we've hit our max offset, and an empty block, we can assume the
			// rest will be empty as well
			break
		}
		for _, toi := range txOut.TxArray {
			result = append(result, castTxOutItem(toi, height))
		}
	}

	return jsonify(ctx, result)
}

func queryPendingOutbound(ctx cosmos.Context, mgr *Mgrs) ([]byte, error) {
	constAccessor := mgr.GetConstants()
	signingTransactionPeriod := constAccessor.GetInt64Value(constants.SigningTransactionPeriod)
	rescheduleCoalesceBlocks := mgr.Keeper().GetConfigInt64(ctx, constants.RescheduleCoalesceBlocks)
	startHeight := ctx.BlockHeight() - signingTransactionPeriod
	if startHeight < 1 {
		startHeight = 1
	}

	// outbounds can be rescheduled to a future height which is the rounded-up nearest multiple of reschedule coalesce blocks
	lastOutboundHeight := ctx.BlockHeight()
	if rescheduleCoalesceBlocks > 1 {
		overBlocks := lastOutboundHeight % rescheduleCoalesceBlocks
		if overBlocks != 0 {
			lastOutboundHeight += rescheduleCoalesceBlocks - overBlocks
		}
	}

	result := make([]openapi.TxOutItem, 0)
	for height := startHeight; height <= lastOutboundHeight; height++ {
		txs, err := mgr.Keeper().GetTxOut(ctx, height)
		if err != nil {
			ctx.Logger().Error("fail to get tx out array from key value store", "error", err)
			return nil, fmt.Errorf("fail to get tx out array from key value store: %w", err)
		}
		for _, tx := range txs.TxArray {
			if tx.OutHash.IsEmpty() {
				result = append(result, castTxOutItem(tx, height))
			}
		}
	}

	return jsonify(ctx, result)
}

func querySwapQueue(ctx cosmos.Context, mgr *Mgrs) ([]byte, error) {
	result := make([]openapi.MsgSwap, 0)

	iterator := mgr.Keeper().GetSwapQueueIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var msg MsgSwap
		if err := mgr.Keeper().Cdc().Unmarshal(iterator.Value(), &msg); err != nil {
			continue
		}
		result = append(result, castMsgSwap(msg))
	}

	return jsonify(ctx, result)
}

func queryTssKeygenMetric(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	var pubKeys common.PubKeys
	if len(path) > 0 {
		pkey, err := common.NewPubKey(path[0])
		if err != nil {
			return nil, fmt.Errorf("fail to parse pubkey(%s) err:%w", path[0], err)
		}
		pubKeys = append(pubKeys, pkey)
	}
	var result []*types.TssKeygenMetric
	for _, pkey := range pubKeys {
		m, err := mgr.Keeper().GetTssKeygenMetric(ctx, pkey)
		if err != nil {
			return nil, fmt.Errorf("fail to get tss keygen metric for pubkey(%s):%w", pkey, err)
		}
		result = append(result, m)
	}
	return jsonify(ctx, result)
}

func queryTssMetric(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	var pubKeys common.PubKeys
	// get all active asgard
	vaults, err := mgr.Keeper().GetAsgardVaultsByStatus(ctx, ActiveVault)
	if err != nil {
		return nil, fmt.Errorf("fail to get active asgards:%w", err)
	}
	for _, v := range vaults {
		pubKeys = append(pubKeys, v.PubKey)
	}
	var keygenMetrics []*types.TssKeygenMetric
	for _, pkey := range pubKeys {
		m, err := mgr.Keeper().GetTssKeygenMetric(ctx, pkey)
		if err != nil {
			return nil, fmt.Errorf("fail to get tss keygen metric for pubkey(%s):%w", pkey, err)
		}
		if len(m.NodeTssTimes) == 0 {
			continue
		}
		keygenMetrics = append(keygenMetrics, m)
	}
	keysignMetric, err := mgr.Keeper().GetLatestTssKeysignMetric(ctx)
	if err != nil {
		return nil, fmt.Errorf("fail to get keysign metric:%w", err)
	}
	m := struct {
		KeygenMetrics []*types.TssKeygenMetric `json:"keygen"`
		KeysignMetric *types.TssKeysignMetric  `json:"keysign"`
	}{
		KeygenMetrics: keygenMetrics,
		KeysignMetric: keysignMetric,
	}
	return jsonify(ctx, m)
}

func queryInvariants(ctx cosmos.Context, mgr *Mgrs) ([]byte, error) {
	result := openapi.InvariantsResponse{}
	for _, route := range mgr.Keeper().InvariantRoutes() {
		result.Invariants = append(result.Invariants, route.Route)
	}
	return jsonify(ctx, result)
}

func queryInvariant(ctx cosmos.Context, path []string, mgr *Mgrs) ([]byte, error) {
	if len(path) < 1 {
		return nil, fmt.Errorf("invalid path: %v", path)
	}
	for _, route := range mgr.Keeper().InvariantRoutes() {
		if strings.EqualFold(route.Route, path[0]) {
			msg, broken := route.Invariant(ctx)
			result := openapi.InvariantResponse{
				Invariant: route.Route,
				Broken:    broken,
				Msg:       msg,
			}
			return jsonify(ctx, result)
		}
	}
	return nil, fmt.Errorf("invariant not registered: %s", path[0])
}

func queryBlock(ctx cosmos.Context, mgr *Mgrs) ([]byte, error) {
	initTendermintOnce.Do(initTendermint)
	height := ctx.BlockHeight()

	// get the block and results from tendermint rpc
	block, err := tendermintClient.Block(ctx.Context(), &height)
	if err != nil {
		return nil, fmt.Errorf("fail to get block from tendermint rpc: %w", err)
	}
	results, err := tendermintClient.BlockResults(ctx.Context(), &height)
	if err != nil {
		return nil, fmt.Errorf("fail to get block results from tendermint rpc: %w", err)
	}

	res := types.QueryBlockResponse{
		BlockResponse: openapi.BlockResponse{
			Id: openapi.BlockResponseId{
				Hash: block.BlockID.Hash.String(),
				Parts: openapi.BlockResponseIdParts{
					Total: int64(block.BlockID.PartSetHeader.Total),
					Hash:  block.BlockID.PartSetHeader.Hash.String(),
				},
			},
			Header: openapi.BlockResponseHeader{
				Version: openapi.BlockResponseHeaderVersion{
					Block: strconv.FormatUint(block.Block.Header.Version.Block, 10),
					App:   strconv.FormatUint(block.Block.Header.Version.App, 10),
				},
				ChainId: block.Block.Header.ChainID,
				Height:  block.Block.Header.Height,
				Time:    block.Block.Header.Time.Format(time.RFC3339Nano),
				LastBlockId: openapi.BlockResponseId{
					Hash: block.Block.Header.LastBlockID.Hash.String(),
					Parts: openapi.BlockResponseIdParts{
						Total: int64(block.Block.Header.LastBlockID.PartSetHeader.Total),
						Hash:  block.Block.Header.LastBlockID.PartSetHeader.Hash.String(),
					},
				},
				LastCommitHash:     block.Block.Header.LastCommitHash.String(),
				DataHash:           block.Block.Header.DataHash.String(),
				ValidatorsHash:     block.Block.Header.ValidatorsHash.String(),
				NextValidatorsHash: block.Block.Header.NextValidatorsHash.String(),
				ConsensusHash:      block.Block.Header.ConsensusHash.String(),
				AppHash:            block.Block.Header.AppHash.String(),
				LastResultsHash:    block.Block.Header.LastResultsHash.String(),
				EvidenceHash:       block.Block.Header.EvidenceHash.String(),
				ProposerAddress:    block.Block.Header.ProposerAddress.String(),
			},
			BeginBlockEvents: []map[string]string{},
			EndBlockEvents:   []map[string]string{},
		},
		Txs: make([]types.QueryBlockTx, len(block.Block.Txs)),
	}

	// parse the events
	for _, event := range results.BeginBlockEvents {
		res.BeginBlockEvents = append(res.BeginBlockEvents, eventMap(sdk.Event(event)))
	}
	for _, event := range results.EndBlockEvents {
		res.EndBlockEvents = append(res.EndBlockEvents, eventMap(sdk.Event(event)))
	}

	for i, tx := range block.Block.Txs {
		res.Txs[i].Hash = strings.ToUpper(hex.EncodeToString(tx.Hash()))

		// decode the protobuf and encode to json
		dtx, err := authtx.DefaultTxDecoder(mgr.cdc.(*codec.ProtoCodec))(tx)
		if err != nil {
			return nil, fmt.Errorf("fail to decode tx: %w", err)
		}
		res.Txs[i].Tx, err = authtx.DefaultJSONTxEncoder(mgr.cdc.(*codec.ProtoCodec))(dtx)
		if err != nil {
			return nil, fmt.Errorf("fail to encode tx: %w", err)
		}

		// parse the tx events
		code := int64(results.TxsResults[i].Code)
		res.Txs[i].Result.Code = &code
		res.Txs[i].Result.Data = wrapString(string(results.TxsResults[i].Data))
		res.Txs[i].Result.Log = wrapString(results.TxsResults[i].Log)
		res.Txs[i].Result.Info = wrapString(results.TxsResults[i].Info)
		res.Txs[i].Result.GasWanted = wrapString(strconv.FormatInt(results.TxsResults[i].GasWanted, 10))
		res.Txs[i].Result.GasUsed = wrapString(strconv.FormatInt(results.TxsResults[i].GasUsed, 10))
		res.Txs[i].Result.Events = []map[string]string{}
		for _, event := range results.TxsResults[i].Events {
			res.Txs[i].Result.Events = append(res.Txs[i].Result.Events, eventMap(sdk.Event(event)))
		}
	}

	return jsonify(ctx, res)
}

// -------------------------------------------------------------------------------------
// Generic Helpers
// -------------------------------------------------------------------------------------

func wrapBool(b bool) *bool {
	if !b {
		return nil
	}
	return &b
}

func wrapString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func wrapInt64(d int64) *int64 {
	if d == 0 {
		return nil
	}
	return &d
}

func wrapUintPtr(uintPtr *cosmos.Uint) *string {
	if uintPtr == nil {
		return nil
	}
	return wrapString(uintPtr.String())
}

func castCoin(sourceCoin common.Coin) openapi.Coin {
	return openapi.Coin{
		Asset:    sourceCoin.Asset.String(),
		Amount:   sourceCoin.Amount.String(),
		Decimals: wrapInt64(sourceCoin.Decimals),
	}
}

func castCoins(sourceCoins ...common.Coin) []openapi.Coin {
	// Leave this nil (null rather than []) if the source is nil.
	if sourceCoins == nil {
		return nil
	}

	coins := make([]openapi.Coin, len(sourceCoins))
	for i := range sourceCoins {
		coins[i] = castCoin(sourceCoins[i])
	}
	return coins
}

func castTxOutItem(toi TxOutItem, height int64) openapi.TxOutItem {
	return openapi.TxOutItem{
		Chain:       toi.Chain.String(),
		ToAddress:   toi.ToAddress.String(),
		VaultPubKey: wrapString(toi.VaultPubKey.String()),
		Coin:        castCoin(toi.Coin),
		Memo:        wrapString(toi.Memo),
		MaxGas:      castCoins(toi.MaxGas...),
		GasRate:     wrapInt64(toi.GasRate),
		InHash:      wrapString(toi.InHash.String()),
		OutHash:     wrapString(toi.OutHash.String()),
		Height:      wrapInt64(height), // Omitted if 0, for use in openapi.TxDetailsResponse
		CloutSpent:  wrapUintPtr(toi.CloutSpent),
	}
}

func castTx(tx common.Tx) openapi.Tx {
	return openapi.Tx{
		Id:          wrapString(tx.ID.String()),
		Chain:       wrapString(tx.Chain.String()),
		FromAddress: wrapString(tx.FromAddress.String()),
		ToAddress:   wrapString(tx.ToAddress.String()),
		Coins:       castCoins(tx.Coins...),
		Gas:         castCoins(tx.Gas...),
		Memo:        wrapString(tx.Memo),
	}
}

func castObservedTx(observedTx ObservedTx) openapi.ObservedTx {
	// Only display the Status if it is "done", not if "incomplete".
	var status *string
	if observedTx.Status != types.Status_incomplete {
		status = wrapString(observedTx.Status.String())
	}

	return openapi.ObservedTx{
		Tx:                              castTx(observedTx.Tx),
		ObservedPubKey:                  wrapString(observedTx.ObservedPubKey.String()),
		ExternalObservedHeight:          wrapInt64(observedTx.BlockHeight),
		ExternalConfirmationDelayHeight: wrapInt64(observedTx.FinaliseHeight),
		Aggregator:                      wrapString(observedTx.Aggregator),
		AggregatorTarget:                wrapString(observedTx.AggregatorTarget),
		AggregatorTargetLimit:           wrapUintPtr(observedTx.AggregatorTargetLimit),
		Signers:                         observedTx.Signers,
		KeysignMs:                       wrapInt64(observedTx.KeysignMs),
		OutHashes:                       observedTx.OutHashes,
		Status:                          status,
	}
}

func castMsgSwap(msg MsgSwap) openapi.MsgSwap {
	// Only display the OrderType if it is "limit", not if "market".
	var orderType *string
	if msg.OrderType != types.OrderType_market {
		orderType = wrapString(msg.OrderType.String())
	}
	// TODO: After order books implementation,
	// always display the OrderType?

	return openapi.MsgSwap{
		Tx:                      castTx(msg.Tx),
		TargetAsset:             msg.TargetAsset.String(),
		Destination:             wrapString(msg.Destination.String()),
		TradeTarget:             msg.TradeTarget.String(),
		AffiliateAddress:        wrapString(msg.AffiliateAddress.String()),
		AffiliateBasisPoints:    msg.AffiliateBasisPoints.String(),
		Signer:                  wrapString(msg.Signer.String()),
		Aggregator:              wrapString(msg.Aggregator),
		AggregatorTargetAddress: wrapString(msg.AggregatorTargetAddress),
		AggregatorTargetLimit:   wrapUintPtr(msg.AggregatorTargetLimit),
		OrderType:               orderType,
		StreamQuantity:          wrapInt64(int64(msg.StreamQuantity)),
		StreamInterval:          wrapInt64(int64(msg.StreamInterval)),
	}
}

func castVaultRouters(chainContracts []ChainContract) []openapi.VaultRouter {
	// Leave this nil (null rather than []) if the source is nil.
	if chainContracts == nil {
		return nil
	}

	routers := make([]openapi.VaultRouter, len(chainContracts))
	for i := range chainContracts {
		routers[i] = openapi.VaultRouter{
			Chain:  wrapString(chainContracts[i].Chain.String()),
			Router: wrapString(chainContracts[i].Router.String()),
		}
	}
	return routers
}

// TODO: Migrate callers to use simulate instead.
func simulateInternal(ctx cosmos.Context, mgr *Mgrs, msg sdk.Msg) (sdk.Events, error) {
	// validate
	err := msg.ValidateBasic()
	if err != nil {
		return nil, fmt.Errorf("failed validate: %w", err)
	}

	// intercept events and avoid modifying state
	cms := ctx.MultiStore().CacheMultiStore() // never call cms.Write()
	em := cosmos.NewEventManager()
	ctx = ctx.WithMultiStore(cms).WithEventManager(em)

	// disable logging
	ctx = ctx.WithLogger(nullLogger)

	// simulate the message handler
	_, err = NewInternalHandler(mgr)(ctx, msg)
	return em.Events(), err
}

func eventMap(e sdk.Event) map[string]string {
	m := map[string]string{}
	m["type"] = e.Type
	for _, a := range e.Attributes {
		m[string(a.Key)] = string(a.Value)
	}
	return m
}

func simulate(ctx cosmos.Context, mgr Manager, msg sdk.Msg) (sdk.Events, error) {
	// use the first active node account as the signer
	nodeAccounts, err := mgr.Keeper().ListActiveValidators(ctx)
	if err != nil {
		return nil, fmt.Errorf("no active node accounts: %w", err)
	}

	// set the signer
	switch m := msg.(type) {
	case *MsgLoanOpen:
		m.Signer = nodeAccounts[0].NodeAddress
	case *MsgLoanRepayment:
		m.Signer = nodeAccounts[0].NodeAddress
	}

	// set random txid
	txid := common.TxID(common.RandHexString(64))
	ctx = ctx.WithValue(constants.CtxLoanTxID, txid)

	// validate
	err = msg.ValidateBasic()
	if err != nil {
		return nil, fmt.Errorf("failed to validate message: %w", err)
	}

	// intercept events and avoid modifying state
	cms := ctx.MultiStore().CacheMultiStore() // never call cms.Write()
	em := cosmos.NewEventManager()
	ctx = ctx.WithMultiStore(cms).WithEventManager(em)

	// disable logging
	// ctx = ctx.WithLogger(nullLogger)

	// reset the swap queue
	iter := mgr.Keeper().GetSwapQueueIterator(ctx)
	for ; iter.Valid(); iter.Next() {
		mgr.Keeper().DeleteKey(ctx, string(iter.Key()))
	}
	iter.Close()

	// save pool state
	pools, err := mgr.Keeper().GetPools(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pools: %w", err)
	}

	// simulate the handler
	_, err = NewInternalHandler(mgr)(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to simulate handler: %w", err)
	}

	// simulate end block, loop it until the swap queue is empty
	var count int64
	for count < 1000 {
		err = mgr.SwapQ().EndBlock(ctx.WithBlockHeight(ctx.BlockHeight()+count), mgr)
		if err != nil {
			return nil, fmt.Errorf("failed to simulate end block: %w", err)
		}

		for _, pool := range pools {
			_ = mgr.Keeper().SetPool(ctx, pool)
		}

		count += 1
		queueEmpty := true
		iter = mgr.Keeper().GetSwapQueueIterator(ctx)
		for ; iter.Valid(); iter.Next() {
			queueEmpty = false
			break
		}
		iter.Close()
		if queueEmpty {
			break
		}
	}

	return em.Events(), nil
}
