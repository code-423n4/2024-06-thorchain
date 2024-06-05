package thorchain

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/blang/semver"
	sdkerrs "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/legacy/legacytx"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/common/tokenlist"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
)

// MsgHandler is an interface expect all handler to implement
type MsgHandler interface {
	Run(ctx cosmos.Context, msg cosmos.Msg) (*cosmos.Result, error)
}

// NewExternalHandler returns a handler for "thorchain" type messages.
func NewExternalHandler(mgr Manager) cosmos.Handler {
	return func(ctx cosmos.Context, msg cosmos.Msg) (_ *cosmos.Result, err error) {
		// TODO: remove the outer if-check on hard fork, always add deferred recover
		if mgr.GetVersion().GTE(semver.MustParse("1.106.0")) {
			defer func() {
				if r := recover(); r != nil {
					// print stack
					stack := make([]byte, 1024)
					length := runtime.Stack(stack, true)
					ctx.Logger().Error("panic", "msg", msg)
					fmt.Println(string(stack[:length]))
					err = fmt.Errorf("panic: %v", r)
				}
			}()
		}

		ctx = ctx.WithEventManager(cosmos.NewEventManager())
		if mgr.GetVersion().LT(semver.MustParse("1.90.0")) {
			_ = mgr.Keeper().GetLowestActiveVersion(ctx) // TODO: remove me on hard fork
		}
		handlerMap := getHandlerMapping(mgr)
		legacyMsg, ok := msg.(legacytx.LegacyMsg)
		if !ok {
			return nil, cosmos.ErrUnknownRequest("unknown message type")
		}
		h, ok := handlerMap[legacyMsg.Type()]
		if !ok {
			errMsg := fmt.Sprintf("Unrecognized thorchain Msg type: %v", legacyMsg.Type())
			return nil, cosmos.ErrUnknownRequest(errMsg)
		}
		result, err := h.Run(ctx, msg)
		if err != nil {
			// TODO: remove version condition on hard fork
			if mgr.GetVersion().GTE(semver.MustParse("1.132.0")) {
				if _, code, _ := sdkerrs.ABCIInfo(err, false); code == 1 {
					// This would be redacted, so wrap it.
					err = sdkerrs.Wrap(errInternal, err.Error())
				}
			}
			return nil, err
		}
		if result == nil {
			result = &cosmos.Result{}
		}
		if len(ctx.EventManager().Events()) > 0 {
			result.Events = ctx.EventManager().ABCIEvents()
		}
		return result, nil
	}
}

func getHandlerMapping(mgr Manager) map[string]MsgHandler {
	return getHandlerMappingV65(mgr)
}

func getHandlerMappingV65(mgr Manager) map[string]MsgHandler {
	// New arch handlers
	m := make(map[string]MsgHandler)

	// Consensus handlers - can only be sent by addresses in
	//   the active validator set.
	m[MsgTssPool{}.Type()] = NewTssHandler(mgr)
	m[MsgObservedTxIn{}.Type()] = NewObservedTxInHandler(mgr)
	m[MsgObservedTxOut{}.Type()] = NewObservedTxOutHandler(mgr)
	m[MsgTssKeysignFail{}.Type()] = NewTssKeysignHandler(mgr)
	m[MsgErrataTx{}.Type()] = NewErrataTxHandler(mgr)
	m[MsgBan{}.Type()] = NewBanHandler(mgr)
	m[MsgNetworkFee{}.Type()] = NewNetworkFeeHandler(mgr)
	m[MsgSolvency{}.Type()] = NewSolvencyHandler(mgr)

	// cli handlers (non-consensus)
	m[MsgMimir{}.Type()] = NewMimirHandler(mgr)
	m[MsgSetNodeKeys{}.Type()] = NewSetNodeKeysHandler(mgr)
	m[MsgSetVersion{}.Type()] = NewVersionHandler(mgr)
	m[MsgSetIPAddress{}.Type()] = NewIPAddressHandler(mgr)
	m[MsgNodePauseChain{}.Type()] = NewNodePauseChainHandler(mgr)

	// native handlers (non-consensus)
	m[MsgSend{}.Type()] = NewSendHandler(mgr)
	m[MsgDeposit{}.Type()] = NewDepositHandler(mgr)
	return m
}

// NewInternalHandler returns a handler for "thorchain" internal type messages.
func NewInternalHandler(mgr Manager) cosmos.Handler {
	return func(ctx cosmos.Context, msg cosmos.Msg) (*cosmos.Result, error) {
		version := mgr.GetVersion()
		if version.LT(semver.MustParse("1.90.0")) {
			version = mgr.Keeper().GetLowestActiveVersion(ctx) // TODO remove me on hard fork
		}
		handlerMap := getInternalHandlerMapping(mgr)
		legacyMsg, ok := msg.(legacytx.LegacyMsg)
		if !ok {
			return nil, cosmos.ErrUnknownRequest("invalid message type")
		}
		h, ok := handlerMap[legacyMsg.Type()]
		if !ok {
			errMsg := fmt.Sprintf("Unrecognized thorchain Msg type: %v", legacyMsg.Type())
			return nil, cosmos.ErrUnknownRequest(errMsg)
		}
		// TODO: remove if-check on hardfork. Always use CacheContext.
		if version.GTE(semver.MustParse("1.88.1")) {
			// CacheContext() returns a context which caches all changes and only forwards
			// to the underlying context when commit() is called. Call commit() only when
			// the handler succeeds, otherwise return error and the changes will be discarded.
			// On commit, cached events also have to be explicitly emitted.
			cacheCtx, commit := ctx.CacheContext()
			res, err := h.Run(cacheCtx, msg)
			if err == nil {
				// Success, commit the cached changes and events
				commit()
				ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
			}
			return res, err
		}
		return h.Run(ctx, msg)
	}
}

func getInternalHandlerMapping(mgr Manager) map[string]MsgHandler {
	version := mgr.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.128.0")):
		return getInternalHandlerMappingV128(mgr)
	case version.GTE(semver.MustParse("1.124.0")):
		return getInternalHandlerMappingV124(mgr)
	case version.GTE(semver.MustParse("1.117.0")):
		return getInternalHandlerMappingV117(mgr)
	default:
		return getInternalHandlerMappingV116(mgr)
	}
}

func getInternalHandlerMappingV128(mgr Manager) map[string]MsgHandler {
	// New arch handlers
	m := make(map[string]MsgHandler)
	m[MsgOutboundTx{}.Type()] = NewOutboundTxHandler(mgr)
	m[MsgSwap{}.Type()] = NewSwapHandler(mgr)
	m[MsgReserveContributor{}.Type()] = NewReserveContributorHandler(mgr)
	m[MsgBond{}.Type()] = NewBondHandler(mgr)
	m[MsgUnBond{}.Type()] = NewUnBondHandler(mgr)
	m[MsgLeave{}.Type()] = NewLeaveHandler(mgr)
	m[MsgDonate{}.Type()] = NewDonateHandler(mgr)
	m[MsgWithdrawLiquidity{}.Type()] = NewWithdrawLiquidityHandler(mgr)
	m[MsgAddLiquidity{}.Type()] = NewAddLiquidityHandler(mgr)
	m[MsgRefundTx{}.Type()] = NewRefundHandler(mgr)
	m[MsgMigrate{}.Type()] = NewMigrateHandler(mgr)
	m[MsgRagnarok{}.Type()] = NewRagnarokHandler(mgr)
	m[MsgNoOp{}.Type()] = NewNoOpHandler(mgr)
	m[MsgConsolidate{}.Type()] = NewConsolidateHandler(mgr)
	m[MsgManageTHORName{}.Type()] = NewManageTHORNameHandler(mgr)
	m[MsgLoanOpen{}.Type()] = NewLoanOpenHandler(mgr)
	m[MsgLoanRepayment{}.Type()] = NewLoanRepaymentHandler(mgr)
	m[MsgTradeAccountDeposit{}.Type()] = NewTradeAccountDepositHandler(mgr)
	m[MsgTradeAccountWithdrawal{}.Type()] = NewTradeAccountWithdrawalHandler(mgr)
	return m
}

func getMsgSwapFromMemo(memo SwapMemo, tx ObservedTx, signer cosmos.AccAddress) (cosmos.Msg, error) {
	if memo.Destination.IsEmpty() {
		memo.Destination = tx.Tx.FromAddress
	}
	return NewMsgSwap(tx.Tx, memo.GetAsset(), memo.Destination, memo.SlipLimit, memo.AffiliateAddress, memo.AffiliateBasisPoints, memo.GetDexAggregator(), memo.GetDexTargetAddress(), memo.GetDexTargetLimit(), memo.GetOrderType(), memo.GetStreamQuantity(), memo.GetStreamInterval(), signer), nil
}

func getMsgWithdrawFromMemo(memo WithdrawLiquidityMemo, tx ObservedTx, signer cosmos.AccAddress) (cosmos.Msg, error) {
	withdrawAmount := cosmos.NewUint(MaxWithdrawBasisPoints)
	if !memo.GetAmount().IsZero() {
		withdrawAmount = memo.GetAmount()
	}
	return NewMsgWithdrawLiquidity(tx.Tx, tx.Tx.FromAddress, withdrawAmount, memo.GetAsset(), memo.GetWithdrawalAsset(), signer), nil
}

func getMsgAddLiquidityFromMemo(ctx cosmos.Context, memo AddLiquidityMemo, tx ObservedTx, signer cosmos.AccAddress) (cosmos.Msg, error) {
	// Extract the Rune amount and the asset amount from the transaction. At least one of them must be
	// nonzero. If THORNode saw two types of coins, one of them must be the asset coin.
	runeCoin := tx.Tx.Coins.GetCoin(common.RuneAsset())
	assetCoin := tx.Tx.Coins.GetCoin(memo.GetAsset())

	var runeAddr common.Address
	var assetAddr common.Address
	if tx.Tx.Chain.Equals(common.THORChain) {
		runeAddr = tx.Tx.FromAddress
		assetAddr = memo.GetDestination()
	} else {
		runeAddr = memo.GetDestination()
		assetAddr = tx.Tx.FromAddress
	}
	// in case we are providing native rune and another native asset
	if memo.GetAsset().Chain.Equals(common.THORChain) {
		assetAddr = runeAddr
	}

	return NewMsgAddLiquidity(tx.Tx, memo.GetAsset(), runeCoin.Amount, assetCoin.Amount, runeAddr, assetAddr, memo.AffiliateAddress, memo.AffiliateBasisPoints, signer), nil
}

func getMsgDonateFromMemo(memo DonateMemo, tx ObservedTx, signer cosmos.AccAddress) (cosmos.Msg, error) {
	runeCoin := tx.Tx.Coins.GetCoin(common.RuneAsset())
	assetCoin := tx.Tx.Coins.GetCoin(memo.GetAsset())
	return NewMsgDonate(tx.Tx, memo.GetAsset(), runeCoin.Amount, assetCoin.Amount, signer), nil
}

func getMsgRefundFromMemo(memo RefundMemo, tx ObservedTx, signer cosmos.AccAddress) (cosmos.Msg, error) {
	return NewMsgRefundTx(tx, memo.GetTxID(), signer), nil
}

func getMsgOutboundFromMemo(memo OutboundMemo, tx ObservedTx, signer cosmos.AccAddress) (cosmos.Msg, error) {
	return NewMsgOutboundTx(tx, memo.GetTxID(), signer), nil
}

func getMsgMigrateFromMemo(memo MigrateMemo, tx ObservedTx, signer cosmos.AccAddress) (cosmos.Msg, error) {
	return NewMsgMigrate(tx, memo.GetBlockHeight(), signer), nil
}

func getMsgRagnarokFromMemo(memo RagnarokMemo, tx ObservedTx, signer cosmos.AccAddress) (cosmos.Msg, error) {
	return NewMsgRagnarok(tx, memo.GetBlockHeight(), signer), nil
}

func getMsgLeaveFromMemo(memo LeaveMemo, tx ObservedTx, signer cosmos.AccAddress) (cosmos.Msg, error) {
	return NewMsgLeave(tx.Tx, memo.GetAccAddress(), signer), nil
}

func getMsgLoanOpenFromMemo(ctx cosmos.Context, keeper keeper.Keeper, memo LoanOpenMemo, tx ObservedTx, signer cosmos.AccAddress, txid common.TxID) (cosmos.Msg, error) {
	version := keeper.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.120.0")):
		return getMsgLoanOpenFromMemoV120(ctx, keeper, memo, tx, signer, txid)
	default:
		return getMsgLoanOpenFromMemoV1(memo, tx, signer, txid)
	}
}

func getMsgLoanOpenFromMemoV120(ctx cosmos.Context, keeper keeper.Keeper, memo LoanOpenMemo, tx ObservedTx, signer cosmos.AccAddress, txid common.TxID) (cosmos.Msg, error) {
	memo.TargetAsset = fuzzyAssetMatch(ctx, keeper, memo.TargetAsset)
	return NewMsgLoanOpen(tx.Tx.FromAddress, tx.Tx.Coins[0].Asset, tx.Tx.Coins[0].Amount, memo.TargetAddress, memo.TargetAsset, memo.GetMinOut(), memo.GetAffiliateAddress(), memo.GetAffiliateBasisPoints(), memo.GetDexAggregator(), memo.GetDexTargetAddress(), memo.DexTargetLimit, signer, txid), nil
}

func getMsgLoanRepaymentFromMemo(memo LoanRepaymentMemo, from common.Address, coin common.Coin, signer cosmos.AccAddress, txid common.TxID) (cosmos.Msg, error) {
	return NewMsgLoanRepayment(memo.Owner, memo.Asset, memo.MinOut, from, coin, signer, txid), nil
}

func getMsgBondFromMemo(memo BondMemo, tx ObservedTx, signer cosmos.AccAddress) (cosmos.Msg, error) {
	coin := tx.Tx.Coins.GetCoin(common.RuneAsset())
	return NewMsgBond(tx.Tx, memo.GetAccAddress(), coin.Amount, tx.Tx.FromAddress, memo.BondProviderAddress, signer, memo.NodeOperatorFee), nil
}

func getMsgUnbondFromMemo(memo UnbondMemo, tx ObservedTx, signer cosmos.AccAddress) (cosmos.Msg, error) {
	return NewMsgUnBond(tx.Tx, memo.GetAccAddress(), memo.GetAmount(), tx.Tx.FromAddress, memo.BondProviderAddress, signer), nil
}

func getMsgManageTHORNameFromMemo(memo ManageTHORNameMemo, tx ObservedTx, signer cosmos.AccAddress) (cosmos.Msg, error) {
	if len(tx.Tx.Coins) == 0 {
		return nil, fmt.Errorf("transaction must have rune in it")
	}
	return NewMsgManageTHORName(memo.Name, memo.Chain, memo.Address, tx.Tx.Coins[0], memo.Expire, memo.PreferredAsset, memo.Owner, signer), nil
}

func processOneTxIn(ctx cosmos.Context, version semver.Version, keeper keeper.Keeper, tx ObservedTx, signer cosmos.AccAddress) (cosmos.Msg, error) {
	switch {
	case version.GTE(semver.MustParse("1.128.0")):
		return processOneTxInV128(ctx, keeper, tx, signer)
	case version.GTE(semver.MustParse("1.124.0")):
		return processOneTxInV124(ctx, keeper, tx, signer)
	case version.GTE(semver.MustParse("1.120.0")):
		return processOneTxInV120(ctx, keeper, tx, signer)
	case version.GTE(semver.MustParse("1.117.0")):
		return processOneTxInV117(ctx, keeper, tx, signer)
	case version.GTE(semver.MustParse("1.107.0")):
		return processOneTxInV107(ctx, keeper, tx, signer)
	case version.GTE(semver.MustParse("0.63.0")):
		return processOneTxInV63(ctx, keeper, tx, signer)
	}
	return nil, errBadVersion
}

func processOneTxInV128(ctx cosmos.Context, keeper keeper.Keeper, tx ObservedTx, signer cosmos.AccAddress) (cosmos.Msg, error) {
	if len(tx.Tx.Coins) != 1 {
		return nil, cosmos.ErrInvalidCoins("only send 1 coins per message")
	}

	memo, err := ParseMemoWithTHORNames(ctx, keeper, tx.Tx.Memo)
	if err != nil {
		ctx.Logger().Error("fail to parse memo", "error", err)
		return nil, err
	}

	// THORNode should not have one tx across chain, if it is cross chain it should be separate tx
	var newMsg cosmos.Msg
	// interpret the memo and initialize a corresponding msg event
	switch m := memo.(type) {
	case AddLiquidityMemo:
		m.Asset = fuzzyAssetMatch(ctx, keeper, m.Asset)
		newMsg, err = getMsgAddLiquidityFromMemo(ctx, m, tx, signer)
	case WithdrawLiquidityMemo:
		m.Asset = fuzzyAssetMatch(ctx, keeper, m.Asset)
		newMsg, err = getMsgWithdrawFromMemo(m, tx, signer)
	case SwapMemo:
		m.Asset = fuzzyAssetMatch(ctx, keeper, m.Asset)
		m.DexTargetAddress = externalAssetMatch(keeper.GetVersion(), m.Asset.GetChain(), m.DexTargetAddress)
		newMsg, err = getMsgSwapFromMemo(m, tx, signer)
	case DonateMemo:
		m.Asset = fuzzyAssetMatch(ctx, keeper, m.Asset)
		newMsg, err = getMsgDonateFromMemo(m, tx, signer)
	case RefundMemo:
		newMsg, err = getMsgRefundFromMemo(m, tx, signer)
	case OutboundMemo:
		newMsg, err = getMsgOutboundFromMemo(m, tx, signer)
	case MigrateMemo:
		newMsg, err = getMsgMigrateFromMemo(m, tx, signer)
	case BondMemo:
		newMsg, err = getMsgBondFromMemo(m, tx, signer)
	case UnbondMemo:
		newMsg, err = getMsgUnbondFromMemo(m, tx, signer)
	case RagnarokMemo:
		newMsg, err = getMsgRagnarokFromMemo(m, tx, signer)
	case LeaveMemo:
		newMsg, err = getMsgLeaveFromMemo(m, tx, signer)
	case ReserveMemo:
		res := NewReserveContributor(tx.Tx.FromAddress, tx.Tx.Coins.GetCoin(common.RuneAsset()).Amount)
		newMsg = NewMsgReserveContributor(tx.Tx, res, signer)
	case NoOpMemo:
		newMsg = NewMsgNoOp(tx, signer, m.Action)
	case ConsolidateMemo:
		newMsg = NewMsgConsolidate(tx, signer)
	case ManageTHORNameMemo:
		newMsg, err = getMsgManageTHORNameFromMemo(m, tx, signer)
	case LoanOpenMemo:
		newMsg, err = getMsgLoanOpenFromMemo(ctx, keeper, m, tx, signer, tx.Tx.ID)
	case LoanRepaymentMemo:
		m.Asset = fuzzyAssetMatch(ctx, keeper, m.Asset)
		from := common.NoAddress
		if keeper.GetVersion().GTE(semver.MustParse("1.110.0")) {
			from = tx.Tx.FromAddress
		}
		newMsg, err = getMsgLoanRepaymentFromMemo(m, from, tx.Tx.Coins[0], signer, tx.Tx.ID)
	case TradeAccountDepositMemo:
		coin := tx.Tx.Coins[0]
		newMsg = NewMsgTradeAccountDeposit(coin.Asset, coin.Amount, m.GetAccAddress(), signer, tx.Tx)
	case TradeAccountWithdrawalMemo:
		coin := tx.Tx.Coins[0]
		newMsg = NewMsgTradeAccountWithdrawal(coin.Asset, coin.Amount, m.GetAddress(), signer, tx.Tx)
	default:
		return nil, errInvalidMemo
	}

	if err != nil {
		return newMsg, err
	}
	// MsgAddLiquidity & MsgSwap has a new version of validateBasic
	switch m := newMsg.(type) {
	case *MsgAddLiquidity:
		switch {
		case keeper.GetVersion().GTE(semver.MustParse("1.98.0")):
			return newMsg, m.ValidateBasicV98()
		case keeper.GetVersion().GTE(semver.MustParse("1.93.0")):
			return newMsg, m.ValidateBasicV93()
		default:
			return newMsg, m.ValidateBasicV63()
		}
	case *MsgSwap:
		return newMsg, m.ValidateBasicV63()
	}
	return newMsg, newMsg.ValidateBasic()
}

func fuzzyAssetMatch(ctx cosmos.Context, keeper keeper.Keeper, asset common.Asset) common.Asset {
	version := keeper.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.126.0")):
		return fuzzyAssetMatchV126(ctx, keeper, asset)
	case version.GTE(semver.MustParse("1.103.0")):
		return fuzzyAssetMatchV103(ctx, keeper, asset)
	case version.GTE(semver.MustParse("1.83.0")):
		return fuzzyAssetMatchV83(ctx, keeper, asset)
	default:
		return fuzzyAssetMatchV1(ctx, keeper, asset)
	}
}

func fuzzyAssetMatchV126(ctx cosmos.Context, keeper keeper.Keeper, origAsset common.Asset) common.Asset {
	asset := origAsset.GetLayer1Asset()
	// if it's already an exact match with successfully-added liquidity, return it immediately
	pool, err := keeper.GetPool(ctx, asset)
	if err != nil {
		return origAsset
	}
	// Only check BalanceRune after checking the error so that no panic if there were an error.
	if !pool.BalanceRune.IsZero() {
		return origAsset
	}

	parts := strings.Split(asset.Symbol.String(), "-")
	hasNoSymbol := len(parts) < 2 || len(parts[1]) == 0
	var symbol string
	if !hasNoSymbol {
		symbol = strings.ToLower(parts[1])
	}
	winner := NewPool()
	// if no asset found, return original asset
	winner.Asset = origAsset
	iterator := keeper.GetPoolIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		if err = keeper.Cdc().Unmarshal(iterator.Value(), &pool); err != nil {
			ctx.Logger().Error("fail to fetch pool", "asset", asset, "err", err)
			continue
		}

		// check chain match
		if !asset.Chain.Equals(pool.Asset.Chain) {
			continue
		}

		// check ticker match
		if !asset.Ticker.Equals(pool.Asset.Ticker) {
			continue
		}

		// check if no symbol given (ie "USDT" or "USDT-")
		if hasNoSymbol {
			// Use LTE rather than LT so this function can only return origAsset or a match
			if winner.BalanceRune.LTE(pool.BalanceRune) {
				winner = pool
			}
			continue
		}

		if strings.HasSuffix(strings.ToLower(pool.Asset.Symbol.String()), symbol) {
			// Use LTE rather than LT so this function can only return origAsset or a match
			if winner.BalanceRune.LTE(pool.BalanceRune) {
				winner = pool
			}
			continue
		}
	}
	winner.Asset.Synth = origAsset.Synth
	return winner.Asset
}

func externalAssetMatch(version semver.Version, chain common.Chain, hint string) string {
	switch {
	case version.GTE(semver.MustParse("1.126.0")):
		return externalAssetMatchV126(version, chain, hint)
	case version.GTE(semver.MustParse("1.95.0")):
		return externalAssetMatchV95(version, chain, hint)
	case version.GTE(semver.MustParse("1.93.0")):
		return externalAssetMatchV93(version, chain, hint)
	default:
		return hint
	}
}

func externalAssetMatchV126(version semver.Version, chain common.Chain, hint string) string {
	if len(hint) == 0 {
		return hint
	}
	if chain.IsEVM() {
		// find all potential matches
		firstMatch := ""
		addrHint := strings.ToLower(hint)
		for _, token := range tokenlist.GetEVMTokenList(chain, version).Tokens {
			if strings.HasSuffix(strings.ToLower(token.Address), addrHint) {
				// store first found address
				if firstMatch == "" {
					firstMatch = token.Address
				} else {
					return hint
				}
			}
		}
		if firstMatch != "" {
			return firstMatch
		}
	}
	return hint
}
