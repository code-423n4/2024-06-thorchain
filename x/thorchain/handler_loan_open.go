package thorchain

import (
	"fmt"

	"github.com/blang/semver"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
	"gitlab.com/thorchain/thornode/mimir"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

// LoanOpenHandler a handler to process bond
type LoanOpenHandler struct {
	mgr Manager
}

// NewLoanOpenHandler create new LoanOpenHandler
func NewLoanOpenHandler(mgr Manager) LoanOpenHandler {
	return LoanOpenHandler{
		mgr: mgr,
	}
}

// Run execute the handler
func (h LoanOpenHandler) Run(ctx cosmos.Context, m cosmos.Msg) (*cosmos.Result, error) {
	msg, ok := m.(*MsgLoanOpen)
	if !ok {
		return nil, errInvalidMessage
	}
	ctx.Logger().Info("receive MsgLoanOpen",
		"owner", msg.Owner,
		"col_asset", msg.CollateralAsset,
		"col_amount", msg.CollateralAmount,
		"target_address", msg.TargetAddress,
		"target_asset", msg.TargetAsset,
		"affiliate", msg.AffiliateAddress,
		"affiliate_basis_points", msg.AffiliateBasisPoints,
	)

	if err := h.validate(ctx, *msg); err != nil {
		ctx.Logger().Error("msg loan fail validation", "error", err)
		return nil, err
	}

	err := h.handle(ctx, *msg)
	if err != nil {
		ctx.Logger().Error("fail to process msg loan", "error", err)
		return nil, err
	}

	return &cosmos.Result{}, nil
}

func (h LoanOpenHandler) validate(ctx cosmos.Context, msg MsgLoanOpen) error {
	version := h.mgr.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.128.0")):
		return h.validateV128(ctx, msg)
	case version.GTE(semver.MustParse("1.121.0")):
		return h.validateV121(ctx, msg)
	case version.GTE(semver.MustParse("1.111.0")):
		return h.validateV111(ctx, msg)
	case version.GTE(semver.MustParse("1.108.0")):
		return h.validateV108(ctx, msg)
	case version.GTE(semver.MustParse("1.107.0")):
		return h.validateV107(ctx, msg)
	default:
		return errBadVersion
	}
}

func (h LoanOpenHandler) validateV128(ctx cosmos.Context, msg MsgLoanOpen) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	pauseLoans := fetchConfigInt64(ctx, h.mgr, constants.PauseLoans)
	if pauseLoans > 0 {
		return fmt.Errorf("loans are currently paused")
	}

	if msg.TargetAsset.IsTradeAsset() || msg.CollateralAsset.IsTradeAsset() {
		return fmt.Errorf("trade assets may not be used for loans")
	}

	// ensure that while derived assets are disabled, borrower cannot receive a
	// derived asset as their debt
	enableDerived := fetchConfigInt64(ctx, h.mgr, constants.EnableDerivedAssets)
	if enableDerived == 0 && msg.TargetAsset.IsDerivedAsset() {
		return fmt.Errorf("cannot receive derived asset")
	}

	// Do not allow a network module as the target address.
	targetAccAddr, err := msg.TargetAddress.AccAddress()
	// A network module address would be resolvable,
	// so if not resolvable it should not be a network module address.
	if err == nil && IsModuleAccAddress(h.mgr.Keeper(), targetAccAddr) {
		return fmt.Errorf("a network module cannot be the target address of a loan open memo")
	}

	// Circuit Breaker: check if we're hit the max supply
	supply := h.mgr.Keeper().GetTotalSupply(ctx, common.RuneAsset())
	maxAmt := fetchConfigInt64(ctx, h.mgr, constants.MaxRuneSupply)
	if maxAmt <= 0 {
		return fmt.Errorf("no max supply set")
	}
	if supply.GTE(cosmos.NewUint(uint64(maxAmt))) {
		return fmt.Errorf("loans are currently paused, due to rune supply cap (%d/%d)", supply.Uint64(), maxAmt)
	}

	// ensure collateral pool exists
	if !h.mgr.Keeper().PoolExist(ctx, msg.CollateralAsset) {
		return fmt.Errorf("collateral asset does not have a pool")
	}

	// The lending key for the ETH.ETH pool would be LENDING-THOR-ETH .
	key := "LENDING-" + msg.CollateralAsset.GetDerivedAsset().MimirString()
	val, err := h.mgr.Keeper().GetMimir(ctx, key)
	if err != nil {
		ctx.Logger().Error("fail to fetch LENDING key", "pool", msg.CollateralAsset.GetDerivedAsset().String(), "error", err)
		return err
	}
	if val <= 0 {
		return fmt.Errorf("Lending is not available for this collateral asset")
	}

	// convert collateral asset back to layer1 asset
	// NOTE: if the symbol of a derived asset isn't the chain, this won't work
	// (ie TERRA.LUNA)
	msg.CollateralAsset.Chain, err = common.NewChain(msg.CollateralAsset.Symbol.String())
	if err != nil {
		return err
	}

	totalCollateral, err := h.mgr.Keeper().GetTotalCollateral(ctx, msg.CollateralAsset)
	if err != nil {
		return err
	}
	totalRune, err := h.getTotalLiquidityRUNELoanPools(ctx)
	if err != nil {
		return err
	}
	if totalRune.IsZero() {
		return fmt.Errorf("no liquidity, lending unavailable")
	}
	lever := fetchConfigInt64(ctx, h.mgr, constants.LendingLever)
	runeBurnt := common.SafeSub(cosmos.NewUint(uint64(maxAmt)), supply)
	totalAvailableRuneForProtocol := common.GetSafeShare(cosmos.NewUint(uint64(lever)), cosmos.NewUint(10_000), runeBurnt) // calculate how much of that rune is available for loans
	if totalAvailableRuneForProtocol.IsZero() {
		return fmt.Errorf("no availability (0), lending unavailable")
	}
	pool, err := h.mgr.Keeper().GetPool(ctx, msg.CollateralAsset)
	if err != nil {
		ctx.Logger().Error("fail to get pool", "error", err)
		return err
	}
	totalAvailableRuneForPool := common.GetSafeShare(pool.BalanceRune, totalRune, totalAvailableRuneForProtocol)
	totalAvailableAssetForPool := pool.RuneValueInAsset(totalAvailableRuneForPool)
	if totalCollateral.Add(msg.CollateralAmount).GT(totalAvailableAssetForPool) {
		return fmt.Errorf("no availability (%d/%d), lending unavailable", totalCollateral.Add(msg.CollateralAmount).Uint64(), totalAvailableAssetForPool.Uint64())
	}

	return nil
}

func (h LoanOpenHandler) handle(ctx cosmos.Context, msg MsgLoanOpen) error {
	version := h.mgr.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.113.0")):
		return h.handleV113(ctx, msg)
	case version.GTE(semver.MustParse("1.111.0")):
		return h.handleV111(ctx, msg)
	case version.GTE(semver.MustParse("1.108.0")):
		return h.handleV108(ctx, msg)
	case version.GTE(semver.MustParse("1.107.0")):
		return h.handleV107(ctx, msg)
	default:
		return errBadVersion
	}
}

func (h LoanOpenHandler) openLoan(ctx cosmos.Context, msg MsgLoanOpen) error {
	version := h.mgr.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.121.0")):
		return h.openLoanV121(ctx, msg)
	case version.GTE(semver.MustParse("1.113.0")):
		return h.openLoanV113(ctx, msg)
	case version.GTE(semver.MustParse("1.112.0")):
		return h.openLoanV112(ctx, msg)
	case version.GTE(semver.MustParse("1.111.0")):
		return h.openLoanV111(ctx, msg)
	case version.GTE(semver.MustParse("1.108.0")):
		return h.openLoanV108(ctx, msg)
	case version.GTE(semver.MustParse("1.107.0")):
		return h.openLoanV107(ctx, msg)
	default:
		return errBadVersion
	}
}

func (h LoanOpenHandler) swap(ctx cosmos.Context, msg MsgLoanOpen) error {
	version := h.mgr.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.132.0")):
		return h.swapV132(ctx, msg)
	case version.GTE(semver.MustParse("1.121.0")):
		return h.swapV121(ctx, msg)
	case version.GTE(semver.MustParse("1.113.0")):
		return h.swapV113(ctx, msg)
	case version.GTE(semver.MustParse("1.108.0")):
		return h.swapV108(ctx, msg)
	case version.GTE(semver.MustParse("1.107.0")):
		return h.swapV107(ctx, msg)
	default:
		return errBadVersion
	}
}

func (h LoanOpenHandler) handleV113(ctx cosmos.Context, msg MsgLoanOpen) error {
	// if the inbound asset is TOR, then lets repay the loan. If not, lets
	// swap first and try again later
	if msg.CollateralAsset.IsDerivedAsset() {
		return h.openLoan(ctx, msg)
	} else {
		return h.swap(ctx, msg)
	}
}

func (h LoanOpenHandler) openLoanV121(ctx cosmos.Context, msg MsgLoanOpen) error {
	var err error
	zero := cosmos.ZeroUint()

	// convert collateral asset back to layer1 asset
	// NOTE: if the symbol of a derived asset isn't the chain, this won't work
	// (ie TERRA.LUNA)
	msg.CollateralAsset.Chain, err = common.NewChain(msg.CollateralAsset.Symbol.String())
	if err != nil {
		return err
	}

	pool, err := h.mgr.Keeper().GetPool(ctx, msg.CollateralAsset)
	if err != nil {
		ctx.Logger().Error("fail to get pool", "error", err)
		return err
	}
	loan, err := h.mgr.Keeper().GetLoan(ctx, msg.CollateralAsset, msg.Owner)
	if err != nil {
		ctx.Logger().Error("fail to get loan", "error", err)
		return err
	}
	totalCollateral, err := h.mgr.Keeper().GetTotalCollateral(ctx, msg.CollateralAsset)
	if err != nil {
		return err
	}

	// move derived asset collateral into lending module
	// TODO: on hard fork, change lending module to an actual module (created as account)
	lendingAcc := h.mgr.Keeper().GetModuleAccAddress(LendingName)
	collateral := common.NewCoin(msg.CollateralAsset.GetDerivedAsset(), msg.CollateralAmount)
	if err := h.mgr.Keeper().SendFromModuleToAccount(ctx, AsgardName, lendingAcc, common.NewCoins(collateral)); err != nil {
		return fmt.Errorf("fail to send collateral funds: %w", err)
	}

	// get configs
	enableDerived := fetchConfigInt64(ctx, h.mgr, constants.EnableDerivedAssets)

	// calculate CR
	cr, err := h.getPoolCR(ctx, pool, msg.CollateralAmount)
	if err != nil {
		return err
	}

	price := h.mgr.Keeper().DollarsPerRune(ctx)
	if price.IsZero() {
		return fmt.Errorf("TOR price cannot be zero")
	}

	collateralValueInRune := pool.AssetValueInRune(msg.CollateralAmount)
	collateralValueInTOR := collateralValueInRune.Mul(price).QuoUint64(1e8)
	debt := collateralValueInTOR.Quo(cr).MulUint64(10_000)
	ctx.Logger().Info("Loan Details", "collateral", common.NewCoin(msg.CollateralAsset, msg.CollateralAmount), "debt", debt.Uint64(), "rune price", price.Uint64(), "colRune", collateralValueInRune.Uint64(), "colTOR", collateralValueInTOR.Uint64())

	// sanity checks
	if debt.IsZero() {
		return fmt.Errorf("debt cannot be zero")
	}

	// if the user has over-repayed the loan, credit the difference on the next open
	cumulativeDebt := debt
	if loan.DebtRepaid.GT(loan.DebtIssued) {
		cumulativeDebt = cumulativeDebt.Add(loan.DebtRepaid.Sub(loan.DebtIssued))
	}

	// update Loan record
	loan.DebtIssued = loan.DebtIssued.Add(cumulativeDebt)
	loan.CollateralDeposited = loan.CollateralDeposited.Add(msg.CollateralAmount)
	loan.LastOpenHeight = ctx.BlockHeight()

	if msg.TargetAsset.Equals(common.TOR) && enableDerived > 0 {
		toi := TxOutItem{
			Chain:      msg.TargetAsset.GetChain(),
			ToAddress:  msg.TargetAddress,
			Coin:       common.NewCoin(common.TOR, cumulativeDebt),
			ModuleName: ModuleName,
		}
		ok, err := h.mgr.TxOutStore().TryAddTxOutItem(ctx, h.mgr, toi, zero)
		if err != nil {
			return err
		}
		if !ok {
			return errFailAddOutboundTx
		}
	} else {
		txID, ok := ctx.Value(constants.CtxLoanTxID).(common.TxID)
		if !ok {
			return fmt.Errorf("fail to get txid")
		}

		torCoin := common.NewCoin(common.TOR, cumulativeDebt)

		if err := h.mgr.Keeper().MintToModule(ctx, ModuleName, torCoin); err != nil {
			return fmt.Errorf("fail to mint loan tor debt: %w", err)
		}
		mintEvt := NewEventMintBurn(MintSupplyType, torCoin.Asset.Native(), torCoin.Amount, "swap")
		if err := h.mgr.EventMgr().EmitEvent(ctx, mintEvt); err != nil {
			ctx.Logger().Error("fail to emit mint event", "error", err)
		}

		if err := h.mgr.Keeper().SendFromModuleToModule(ctx, ModuleName, AsgardName, common.NewCoins(torCoin)); err != nil {
			return fmt.Errorf("fail to send TOR vault funds: %w", err)
		}

		// Get streaming swaps interval to use for loan swap
		ssInterval := h.mgr.Keeper().GetConfigInt64(ctx, constants.LoanStreamingSwapsInterval)
		if ssInterval <= 0 || !msg.MinOut.IsZero() {
			ssInterval = 0
		}

		// As this is to be a swap from TOR which has been sent to AsgardName, the ToAddress should be AsgardName's address.
		tx := common.NewTx(txID, common.NoopAddress, common.NoopAddress, common.NewCoins(torCoin), nil, "noop")
		// we do NOT pass affiliate info here as it was already taken out on the swap of the collateral to derived asset
		swapMsg := NewMsgSwap(tx, msg.TargetAsset, msg.TargetAddress, msg.MinOut, common.NoAddress, zero, msg.Aggregator, msg.AggregatorTargetAddress, &msg.AggregatorTargetLimit, 0, 0, uint64(ssInterval), msg.Signer)
		if ssInterval == 0 {
			handler := NewSwapHandler(h.mgr)
			if _, err := handler.Run(ctx, swapMsg); err != nil {
				ctx.Logger().Error("fail to make second swap when opening a loan", "error", err)
				return err
			}
		} else {
			if err := h.mgr.Keeper().SetSwapQueueItem(ctx, *swapMsg, 1); err != nil {
				ctx.Logger().Error("fail to add swap to queue", "error", err)
				return err
			}
		}
	}

	// update kvstore
	h.mgr.Keeper().SetLoan(ctx, loan)
	h.mgr.Keeper().SetTotalCollateral(ctx, msg.CollateralAsset, totalCollateral.Add(msg.CollateralAmount))

	// emit events and metrics
	evt := NewEventLoanOpen(msg.CollateralAmount, cr, debt, msg.CollateralAsset, msg.TargetAsset, msg.Owner, msg.TxID)
	if err := h.mgr.EventMgr().EmitEvent(ctx, evt); nil != err {
		ctx.Logger().Error("fail to emit loan open event", "error", err)
	}

	return nil
}

func (h LoanOpenHandler) getPoolCR(ctx cosmos.Context, pool Pool, collateralAmount cosmos.Uint) (cosmos.Uint, error) {
	minCR := fetchConfigInt64(ctx, h.mgr, constants.MinCR)
	maxCR := fetchConfigInt64(ctx, h.mgr, constants.MaxCR)
	lever := fetchConfigInt64(ctx, h.mgr, constants.LendingLever)

	currentRuneSupply := h.mgr.Keeper().GetTotalSupply(ctx, common.RuneAsset())
	maxRuneSupply := fetchConfigInt64(ctx, h.mgr, constants.MaxRuneSupply)
	if maxRuneSupply <= 0 {
		return cosmos.ZeroUint(), fmt.Errorf("no max supply set")
	}
	runeBurnt := common.SafeSub(cosmos.NewUint(uint64(maxRuneSupply)), currentRuneSupply)
	totalAvailableRuneForProtocol := common.GetSafeShare(cosmos.NewUint(uint64(lever)), cosmos.NewUint(10_000), runeBurnt) // calculate how much of that rune is available for loans
	if totalAvailableRuneForProtocol.IsZero() {
		return cosmos.ZeroUint(), fmt.Errorf("no availability (0), lending unavailable")
	}

	totalCollateral, err := h.mgr.Keeper().GetTotalCollateral(ctx, pool.Asset)
	if err != nil {
		return cosmos.ZeroUint(), err
	}

	totalRune, err := h.getTotalLiquidityRUNELoanPools(ctx)
	if err != nil {
		return cosmos.ZeroUint(), err
	}
	if totalRune.IsZero() {
		return cosmos.ZeroUint(), fmt.Errorf("no liquidity, lending unavailable")
	}

	totalAvailableRuneForPool := common.GetSafeShare(pool.BalanceRune, totalRune, totalAvailableRuneForProtocol)
	totalAvailableAssetForPool := pool.RuneValueInAsset(totalAvailableRuneForPool)
	if totalCollateral.Add(collateralAmount).GT(totalAvailableAssetForPool) {
		return cosmos.ZeroUint(), fmt.Errorf("no availability (%d/%d), lending unavailable", totalCollateral.Add(collateralAmount).Uint64(), totalAvailableAssetForPool.Uint64())
	}
	cr := h.calcCR(totalCollateral.Add(collateralAmount), totalAvailableAssetForPool, minCR, maxCR)

	return cr, nil
}

func (h LoanOpenHandler) calcCR(a, b cosmos.Uint, minCR, maxCR int64) cosmos.Uint {
	// (maxCR - minCR) / (b / a) + minCR
	// NOTE: a should include the collateral currently being deposited
	crCalc := cosmos.NewUint(uint64(maxCR - minCR))
	cr := common.GetUncappedShare(a, b, crCalc)
	return cr.AddUint64(uint64(minCR))
}

func (h LoanOpenHandler) swapV132(ctx cosmos.Context, msg MsgLoanOpen) error {
	txID, ok := ctx.Value(constants.CtxLoanTxID).(common.TxID)
	if !ok {
		return fmt.Errorf("fail to get txid")
	}
	// ensure TxID does NOT have a collision with another swap, this could
	// happen if the user submits two identical loan requests in the same
	// block
	if ok := h.mgr.Keeper().HasSwapQueueItem(ctx, txID, 0); ok {
		return fmt.Errorf("txn hash conflict")
	}

	toAddress, ok := ctx.Value(constants.CtxLoanToAddress).(common.Address)
	// An empty ToAddress fails Tx validation,
	// and a querier quote or unit test has no provided ToAddress.
	// As this only affects emitted swap event contents, do not return an error.
	if !ok || toAddress.IsEmpty() {
		toAddress = "no to address available"
	}

	// Get streaming swaps interval to use for loan swap
	ssInterval := h.mgr.Keeper().GetConfigInt64(ctx, constants.LoanStreamingSwapsInterval)
	if ssInterval <= 0 || !msg.MinOut.IsZero() {
		ssInterval = 0
	}

	collateral := common.NewCoin(msg.CollateralAsset, msg.CollateralAmount)
	maxAffPoints := mimir.NewAffiliateFeeBasisPointsMax().FetchValue(ctx, h.mgr.Keeper())

	// only take affiliate fee if parameters are set and it's the original swap (not the derived asset swap)
	if !msg.AffiliateBasisPoints.IsZero() && msg.AffiliateBasisPoints.LTE(cosmos.NewUint(uint64(maxAffPoints))) && !msg.AffiliateAddress.IsEmpty() && !msg.CollateralAsset.IsNative() {
		newAmt, err := h.handleAffiliateSwap(ctx, msg, collateral)
		if err != nil {
			ctx.Logger().Error("fail to handle affiliate swap", "error", err)
		} else {
			collateral.Amount = newAmt
		}
	}

	memo := fmt.Sprintf("loan+:%s:%s:%d:%s:%d:%s:%s:%d", msg.TargetAsset, msg.TargetAddress, msg.MinOut.Uint64(), msg.AffiliateAddress, msg.AffiliateBasisPoints.Uint64(), msg.Aggregator, msg.AggregatorTargetAddress, msg.AggregatorTargetLimit.Uint64())
	fakeGas := common.NewCoin(msg.CollateralAsset.GetChain().GetGasAsset(), cosmos.OneUint())
	tx := common.NewTx(txID, msg.Owner, toAddress, common.NewCoins(collateral), common.Gas{fakeGas}, memo)
	swapMsg := NewMsgSwap(tx, msg.CollateralAsset.GetDerivedAsset(), common.NoopAddress, cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(), "", "", nil, 0, 0, uint64(ssInterval), msg.Signer)
	if err := h.mgr.Keeper().SetSwapQueueItem(ctx, *swapMsg, 0); err != nil {
		ctx.Logger().Error("fail to add swap to queue", "error", err)
		return err
	}

	return nil
}

// handleAffiliateSwap handles the affiliate swap for the loan open and returns updated
// collateral amount for the loan with affiliate amount deducted
func (h LoanOpenHandler) handleAffiliateSwap(ctx cosmos.Context, msg MsgLoanOpen, collateral common.Coin) (cosmos.Uint, error) {
	version := h.mgr.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.132.0")):
		return h.handleAffiliateSwapV132(ctx, msg, collateral)
	default:
		return cosmos.ZeroUint(), errBadVersion
	}
}

func (h LoanOpenHandler) handleAffiliateSwapV132(ctx cosmos.Context, msg MsgLoanOpen, collateral common.Coin) (cosmos.Uint, error) {
	// Setup affiliate swap
	affAmt := common.GetSafeShare(
		msg.AffiliateBasisPoints,
		cosmos.NewUint(constants.MaxBasisPts),
		msg.CollateralAmount,
	)

	affCoin := common.NewCoin(msg.CollateralAsset, affAmt)
	gasCoin := common.NewCoin(msg.CollateralAsset.GetChain().GetGasAsset(), cosmos.OneUint())
	fakeTx := common.NewTx(msg.TxID, common.NoopAddress, common.NoopAddress, common.NewCoins(affCoin), common.Gas{gasCoin}, "noop")
	affiliateSwap := NewMsgSwap(fakeTx, common.RuneAsset(), msg.AffiliateAddress, cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(), "", "", nil, 0, 0, 0, msg.Signer)

	var affThorname *types.THORName
	voter, err := h.mgr.Keeper().GetObservedTxInVoter(ctx, msg.TxID)
	if err == nil {
		memo, err := ParseMemoWithTHORNames(ctx, h.mgr.Keeper(), voter.Tx.Tx.Memo)
		if err != nil {
			ctx.Logger().Error("fail to parse memo", "error", err)
		}
		affThorname = memo.GetAffiliateTHORName()
	}

	// PreferredAsset set, swap to the AffiliateCollector Module + check if the
	// preferred asset swap should be triggered
	if affThorname != nil && !affThorname.PreferredAsset.IsEmpty() {
		affcol, err := h.mgr.Keeper().GetAffiliateCollector(ctx, affThorname.Owner)
		if err != nil {
			return collateral.Amount, err
		}
		affColAddress, err := h.mgr.Keeper().GetModuleAddress(AffiliateCollectorName)
		if err != nil {
			return collateral.Amount, err
		}
		// Set AffiliateCollector Module as destination and populate the AffiliateAddress
		// so that the swap handler can increment the emitted RUNE for the affiliate in
		// the AffiliateCollector KVStore.
		affiliateSwap.Destination = affColAddress
		affiliateSwap.AffiliateAddress = msg.AffiliateAddress
		// Need to set the memo as a normal swap, so that the swap handler doesn't run the
		// open loan handler for the affiliate swaps
		affiliateSwap.Tx.Memo = NewSwapMemo(ctx, h.mgr, common.RuneNative, msg.AffiliateAddress, cosmos.ZeroUint(), affThorname.Name, cosmos.ZeroUint())

		// Check if accrued RUNE is 100x current outbound fee of preferred asset, if
		// so trigger the preferred asset swap
		ofRune, err := h.mgr.GasMgr().GetAssetOutboundFee(ctx, affThorname.PreferredAsset, true)
		if err != nil {
			ctx.Logger().Error("fail to get outbound fee", "err", err)
		}
		multiplier := h.mgr.Keeper().GetConfigInt64(ctx, constants.PreferredAssetOutboundFeeMultiplier)
		threshold := ofRune.Mul(cosmos.NewUint(uint64(multiplier)))
		if err == nil && affcol.RuneAmount.GT(threshold) {
			if err = triggerPreferredAssetSwap(ctx, h.mgr, msg.AffiliateAddress, msg.TxID, *affThorname, affcol, 3); err != nil {
				ctx.Logger().Error("fail to queue preferred asset swap", "thorname", affThorname.Name, "err", err)
			}
		}
	}

	// If the affiliate swap would exceed the native tx fee, add it to the queue
	if willSwapOutputExceedLimitAndFees(ctx, h.mgr, *affiliateSwap) {
		if err := h.mgr.Keeper().SetSwapQueueItem(ctx, *affiliateSwap, 2); err != nil {
			return collateral.Amount, fmt.Errorf("fail to add affiliate swap to queue: %w", err)
		}
		collateral.Amount = common.SafeSub(collateral.Amount, affAmt)
	}

	return collateral.Amount, nil
}

func (h LoanOpenHandler) getTotalLiquidityRUNELoanPools(ctx cosmos.Context) (cosmos.Uint, error) {
	version := h.mgr.GetVersion()
	switch {
	case version.GTE(semver.MustParse("1.108.0")):
		return h.getTotalLiquidityRUNELoanPoolsV108(ctx)
	case version.GTE(semver.MustParse("1.107.0")):
		return h.getTotalLiquidityRUNELoanPoolsV107(ctx)
	default:
		return cosmos.ZeroUint(), errBadVersion
	}
}

// getTotalLiquidityRUNE we have in all pools
func (h LoanOpenHandler) getTotalLiquidityRUNELoanPoolsV108(ctx cosmos.Context) (cosmos.Uint, error) {
	pools, err := h.mgr.Keeper().GetPools(ctx)
	if err != nil {
		return cosmos.ZeroUint(), fmt.Errorf("fail to get pools from data store: %w", err)
	}
	total := cosmos.ZeroUint()
	for _, p := range pools {
		// ignore suspended pools
		if p.Status == PoolSuspended {
			continue
		}
		if p.Asset.IsVaultAsset() {
			continue
		}
		if p.Asset.IsDerivedAsset() {
			continue
		}

		key := "LENDING-" + p.Asset.GetDerivedAsset().MimirString()
		val, err := h.mgr.Keeper().GetMimir(ctx, key)
		if err != nil {
			continue
		}
		if val <= 0 {
			continue
		}
		total = total.Add(p.BalanceRune)
	}
	return total, nil
}

func (h LoanOpenHandler) GetLoanCollateralRemainingForPool(ctx cosmos.Context, pool Pool) (cosmos.Uint, error) {
	lever := fetchConfigInt64(ctx, h.mgr, constants.LendingLever)

	currentRuneSupply := h.mgr.Keeper().GetTotalSupply(ctx, common.RuneAsset())
	maxRuneSupply := fetchConfigInt64(ctx, h.mgr, constants.MaxRuneSupply)
	if maxRuneSupply <= 0 {
		return cosmos.ZeroUint(), fmt.Errorf("no max supply set")
	}
	runeBurnt := common.SafeSub(cosmos.NewUint(uint64(maxRuneSupply)), currentRuneSupply)
	// calculate total rune available for loans
	totalAvailableRuneForProtocol := common.GetSafeShare(cosmos.NewUint(uint64(lever)), cosmos.NewUint(constants.MaxBasisPts), runeBurnt)
	totalCollateral, err := h.mgr.Keeper().GetTotalCollateral(ctx, pool.Asset)
	if err != nil {
		return cosmos.ZeroUint(), err
	}

	totalRune, err := h.getTotalLiquidityRUNELoanPools(ctx)
	if err != nil {
		return cosmos.ZeroUint(), err
	}

	totalAvailableRuneForPool := common.GetSafeShare(pool.BalanceRune, totalRune, totalAvailableRuneForProtocol)
	totalAvailableAssetForPool := pool.RuneValueInAsset(totalAvailableRuneForPool)

	loanCollateralRemainingForPool := common.SafeSub(totalAvailableAssetForPool, totalCollateral)

	return loanCollateralRemainingForPool, nil
}
