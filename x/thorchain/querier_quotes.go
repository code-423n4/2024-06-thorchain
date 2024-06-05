package thorchain

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/url"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"
	abci "github.com/tendermint/tendermint/abci/types"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
	"gitlab.com/thorchain/thornode/log"
	openapi "gitlab.com/thorchain/thornode/openapi/gen"
	mem "gitlab.com/thorchain/thornode/x/thorchain/memo"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

// -------------------------------------------------------------------------------------
// Config
// -------------------------------------------------------------------------------------

const (
	heightParam               = "height"
	fromAssetParam            = "from_asset"
	toAssetParam              = "to_asset"
	assetParam                = "asset"
	addressParam              = "address"
	loanOwnerParam            = "loan_owner"
	withdrawBasisPointsParam  = "withdraw_bps"
	amountParam               = "amount"
	repayBpsParam             = "repay_bps"
	destinationParam          = "destination"
	toleranceBasisPointsParam = "tolerance_bps"
	affiliateParam            = "affiliate"
	affiliateBpsParam         = "affiliate_bps"
	minOutParam               = "min_out"
	intervalParam             = "streaming_interval"
	quantityParam             = "streaming_quantity"
	refundAddressParam        = "refund_address"

	quoteWarning         = "Do not cache this response. Do not send funds after the expiry."
	quoteExpiration      = 15 * time.Minute
	ethBlockRewardAndFee = 3 * 1e18
)

var nullLogger = &log.TendermintLogWrapper{Logger: zerolog.New(io.Discard)}

// -------------------------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------------------------

func quoteErrorResponse(err error) ([]byte, error) {
	return json.Marshal(map[string]string{"error": err.Error()})
}

func quoteParseParams(data []byte) (params url.Values, err error) {
	// parse the query parameters
	u, err := url.ParseRequestURI(string(data))
	if err != nil {
		return nil, fmt.Errorf("bad params: %w", err)
	}

	// error if parameters were not provided
	if len(u.Query()) == 0 {
		return nil, fmt.Errorf("no parameters provided")
	}

	return u.Query(), nil
}

func quoteParseAddress(ctx cosmos.Context, mgr *Mgrs, addrString string, chain common.Chain) (common.Address, error) {
	if addrString == "" {
		return common.NoAddress, nil
	}

	// attempt to parse a raw address
	addr, err := common.NewAddress(addrString)
	if err == nil {
		return addr, nil
	}

	// attempt to lookup a thorname address
	name, err := mgr.Keeper().GetTHORName(ctx, addrString)
	if err != nil {
		return common.NoAddress, fmt.Errorf("unable to parse address: %w", err)
	}

	// find the address for the correct chain
	for _, alias := range name.Aliases {
		if alias.Chain.Equals(chain) {
			return alias.Address, nil
		}
	}

	return common.NoAddress, fmt.Errorf("no thorname alias for chain %s", chain)
}

func quoteHandleAffiliate(ctx cosmos.Context, mgr *Mgrs, params url.Values, amount sdk.Uint) (affiliate common.Address, memo string, bps, newAmount, affiliateAmt sdk.Uint, err error) {
	// parse affiliate
	affAmt := cosmos.ZeroUint()
	memo = "" // do not resolve thorname for the memo
	if len(params[affiliateParam]) > 0 {
		affiliate, err = quoteParseAddress(ctx, mgr, params[affiliateParam][0], common.THORChain)
		if err != nil {
			err = fmt.Errorf("bad affiliate address: %w", err)
			return
		}
		memo = params[affiliateParam][0]
	}

	// parse affiliate fee
	bps = sdk.NewUint(0)
	if len(params[affiliateBpsParam]) > 0 {
		bps, err = sdk.ParseUint(params[affiliateBpsParam][0])
		if err != nil {
			err = fmt.Errorf("bad affiliate fee: %w", err)
			return
		}
	}

	// verify affiliate fee
	if bps.GT(sdk.NewUint(10000)) {
		err = fmt.Errorf("affiliate fee must be less than 10000 bps")
		return
	}

	// compute the new swap amount if an affiliate fee will be taken first
	if affiliate != common.NoAddress && !bps.IsZero() {
		// calculate the affiliate amount
		affAmt = common.GetSafeShare(
			bps,
			cosmos.NewUint(10000),
			amount,
		)

		// affiliate fee modifies amount at observation before the swap
		amount = amount.Sub(affAmt)
	}

	return affiliate, memo, bps, amount, affAmt, nil
}

func hasSuffixMatch(suffix string, values []string) bool {
	for _, value := range values {
		if strings.HasSuffix(value, suffix) {
			return true
		}
	}
	return false
}

// quoteConvertAsset - converts amount to target asset using THORChain pools
func quoteConvertAsset(ctx cosmos.Context, mgr *Mgrs, fromAsset common.Asset, amount sdk.Uint, toAsset common.Asset) (sdk.Uint, error) {
	// no conversion necessary
	if fromAsset.Equals(toAsset) {
		return amount, nil
	}

	// convert to rune
	if !fromAsset.IsRune() {
		// get the fromPool for the from asset
		fromPool, err := mgr.Keeper().GetPool(ctx, fromAsset.GetLayer1Asset())
		if err != nil {
			return sdk.ZeroUint(), fmt.Errorf("failed to get pool: %w", err)
		}

		// ensure pool exists
		if fromPool.IsEmpty() {
			return sdk.ZeroUint(), fmt.Errorf("pool does not exist")
		}

		amount = fromPool.AssetValueInRune(amount)
	}

	// convert to target asset
	if !toAsset.IsRune() {

		toPool, err := mgr.Keeper().GetPool(ctx, toAsset.GetLayer1Asset())
		if err != nil {
			return sdk.ZeroUint(), fmt.Errorf("failed to get pool: %w", err)
		}

		// ensure pool exists
		if toPool.IsEmpty() {
			return sdk.ZeroUint(), fmt.Errorf("pool does not exist")
		}

		amount = toPool.RuneValueInAsset(amount)
	}

	return amount, nil
}

func quoteReverseFuzzyAsset(ctx cosmos.Context, mgr *Mgrs, asset common.Asset) (common.Asset, error) {
	// get all pools
	pools, err := mgr.Keeper().GetPools(ctx)
	if err != nil {
		return asset, fmt.Errorf("failed to get pools: %w", err)
	}

	// return the asset if no symbol to shorten
	aSplit := strings.Split(asset.Symbol.String(), "-")
	if len(aSplit) == 1 {
		return asset, nil
	}

	// find all other assets that match the chain and ticker
	// (without exactly matching the symbol)
	addressMatches := []string{}
	for _, p := range pools {
		if p.IsAvailable() && !p.IsEmpty() && !p.Asset.IsVaultAsset() &&
			!p.Asset.Symbol.Equals(asset.Symbol) &&
			p.Asset.Chain.Equals(asset.Chain) && p.Asset.Ticker.Equals(asset.Ticker) {
			pSplit := strings.Split(p.Asset.Symbol.String(), "-")
			if len(pSplit) != 2 {
				return asset, fmt.Errorf("ambiguous match: %s", p.Asset.Symbol)
			}
			addressMatches = append(addressMatches, pSplit[1])
		}
	}

	if len(addressMatches) == 0 { // if only one match, drop the address
		asset.Symbol = common.Symbol(asset.Ticker)
	} else { // find the shortest unique suffix of the asset symbol
		address := aSplit[1]

		for i := len(address) - 1; i > 0; i-- {
			if !hasSuffixMatch(address[i:], addressMatches) {
				asset.Symbol = common.Symbol(
					fmt.Sprintf("%s-%s", asset.Ticker, address[i:]),
				)
				break
			}
		}
	}

	return asset, nil
}

// NOTE: streamingQuantity > 0 is a precondition.
func quoteSimulateSwap(ctx cosmos.Context, mgr *Mgrs, amount sdk.Uint, msg *MsgSwap, streamingQuantity uint64) (
	res *openapi.QuoteSwapResponse, emitAmount, outboundFeeAmount sdk.Uint, err error,
) {
	// should be unreachable
	if streamingQuantity == 0 {
		return nil, sdk.ZeroUint(), sdk.ZeroUint(), fmt.Errorf("streaming quantity must be greater than zero")
	}

	msg.Tx.Coins[0].Amount = msg.Tx.Coins[0].Amount.QuoUint64(streamingQuantity)

	// if the generated memo is too long for the source chain send error
	maxMemoLength := msg.Tx.Coins[0].Asset.Chain.MaxMemoLength()
	if !msg.Tx.Coins[0].Asset.Synth && len(msg.Tx.Memo) > maxMemoLength {
		return nil, sdk.ZeroUint(), sdk.ZeroUint(), fmt.Errorf("generated memo too long for source chain")
	}

	// use the first active node account as the signer
	nodeAccounts, err := mgr.Keeper().ListActiveValidators(ctx)
	if err != nil {
		return nil, sdk.ZeroUint(), sdk.ZeroUint(), fmt.Errorf("no active node accounts: %w", err)
	}
	msg.Signer = nodeAccounts[0].NodeAddress

	// simulate the swap
	events, err := simulateInternal(ctx, mgr, msg)
	if err != nil {
		return nil, sdk.ZeroUint(), sdk.ZeroUint(), err
	}

	// extract events
	var swaps []map[string]string
	var fee map[string]string
	for _, e := range events {
		switch e.Type {
		case "swap":
			swaps = append(swaps, eventMap(e))
		case "fee":
			fee = eventMap(e)
		}
	}
	finalSwap := swaps[len(swaps)-1]

	// parse outbound fee from event (except on trade assets with no outbound fee)
	outboundFeeAmount = sdk.ZeroUint()
	if !msg.TargetAsset.IsTradeAsset() {
		outboundFeeCoin, err := common.ParseCoin(fee["coins"])
		if err != nil {
			return nil, sdk.ZeroUint(), sdk.ZeroUint(), fmt.Errorf("unable to parse outbound fee coin: %w", err)
		}
		outboundFeeAmount = outboundFeeCoin.Amount
	}

	// parse outbound amount from event
	emitCoin, err := common.ParseCoin(finalSwap["emit_asset"])
	if err != nil {
		return nil, sdk.ZeroUint(), sdk.ZeroUint(), fmt.Errorf("unable to parse emit coin: %w", err)
	}
	emitAmount = emitCoin.Amount.MulUint64(streamingQuantity)

	// sum the liquidity fees and convert to target asset
	liquidityFee := sdk.ZeroUint()
	for _, s := range swaps {
		liquidityFee = liquidityFee.Add(sdk.NewUintFromString(s["liquidity_fee_in_rune"]))
	}
	var targetPool types.Pool
	if !msg.TargetAsset.IsNativeRune() {
		targetPool, err = mgr.Keeper().GetPool(ctx, msg.TargetAsset.GetLayer1Asset())
		if err != nil {
			return nil, sdk.ZeroUint(), sdk.ZeroUint(), fmt.Errorf("unable to get pool: %w", err)
		}
		liquidityFee = targetPool.RuneValueInAsset(liquidityFee)
	}
	liquidityFee = liquidityFee.MulUint64(streamingQuantity)

	// approximate the affiliate fee in the target asset
	affiliateFee := sdk.ZeroUint()
	if msg.AffiliateAddress != common.NoAddress && !msg.AffiliateBasisPoints.IsZero() {
		inAsset := msg.Tx.Coins[0].Asset.GetLayer1Asset()
		if !inAsset.IsNativeRune() {
			pool, err := mgr.Keeper().GetPool(ctx, msg.Tx.Coins[0].Asset.GetLayer1Asset())
			if err != nil {
				return nil, sdk.ZeroUint(), sdk.ZeroUint(), fmt.Errorf("unable to get pool: %w", err)
			}
			amount = pool.AssetValueInRune(amount)
		}
		affiliateFee = common.GetUncappedShare(msg.AffiliateBasisPoints, cosmos.NewUint(10_000), amount)
		if !msg.TargetAsset.IsNativeRune() {
			affiliateFee = targetPool.RuneValueInAsset(affiliateFee)
		}
	}

	// compute slip based on emit amount instead of slip in event to handle double swap
	slippageBps := liquidityFee.MulUint64(10000).Quo(emitAmount.Add(liquidityFee))

	// build fees
	totalFees := affiliateFee.Add(liquidityFee).Add(outboundFeeAmount)
	fees := openapi.QuoteFees{
		Asset:       msg.TargetAsset.String(),
		Affiliate:   wrapString(affiliateFee.String()),
		Liquidity:   liquidityFee.String(),
		Outbound:    wrapString(outboundFeeAmount.String()),
		Total:       totalFees.String(),
		SlippageBps: slippageBps.BigInt().Int64(),
		TotalBps:    totalFees.MulUint64(10000).Quo(emitAmount.Add(totalFees)).BigInt().Int64(),
	}

	// build response from simulation result events
	return &openapi.QuoteSwapResponse{
		ExpectedAmountOut: emitAmount.String(),
		Fees:              fees,
		// TODO: notify clients to migrate to fees object and deprecate
		SlippageBps: slippageBps.BigInt().Int64(),
	}, emitAmount, outboundFeeAmount, nil
}

func convertThorchainAmountToWei(amt *big.Int) *big.Int {
	return big.NewInt(0).Mul(amt, big.NewInt(common.One*100))
}

func quoteInboundInfo(ctx cosmos.Context, mgr *Mgrs, amount sdk.Uint, chain common.Chain, asset common.Asset) (address, router common.Address, confirmations int64, err error) {
	// If inbound chain is THORChain there is no inbound address
	if chain.IsTHORChain() {
		address = common.NoAddress
		router = common.NoAddress
	} else {
		// get the most secure vault for inbound
		active, err := mgr.Keeper().GetAsgardVaultsByStatus(ctx, ActiveVault)
		if err != nil {
			return common.NoAddress, common.NoAddress, 0, err
		}
		constAccessor := mgr.GetConstants()
		signingTransactionPeriod := constAccessor.GetInt64Value(constants.SigningTransactionPeriod)
		vault := mgr.Keeper().GetMostSecure(ctx, active, signingTransactionPeriod)
		address, err = vault.PubKey.GetAddress(chain)
		if err != nil {
			return common.NoAddress, common.NoAddress, 0, err
		}

		router = common.NoAddress
		if chain.IsEVM() {
			router = vault.GetContract(chain).Router
		}
	}

	// estimate the inbound confirmation count blocks: ceil(amount/coinbase * conf adjustment)
	confMul, err := mgr.Keeper().GetMimirWithRef(ctx, constants.MimirTemplateConfMultiplierBasisPoints, chain.String())
	if confMul < 0 || err != nil {
		confMul = int64(constants.MaxBasisPts)
	}
	if chain.DefaultCoinbase() > 0 {
		confValue := common.GetUncappedShare(cosmos.NewUint(uint64(confMul)), cosmos.NewUint(constants.MaxBasisPts), cosmos.NewUint(uint64(chain.DefaultCoinbase())*common.One))
		confirmations = amount.Quo(confValue).BigInt().Int64()
		if !amount.Mod(confValue).IsZero() {
			confirmations++
		}
	} else if chain.Equals(common.ETHChain) {
		// copying logic from getBlockRequiredConfirmation of ethereum.go
		// convert amount to ETH
		gasAssetAmount, err := quoteConvertAsset(ctx, mgr, asset, amount, chain.GetGasAsset())
		if err != nil {
			return common.NoAddress, common.NoAddress, 0, fmt.Errorf("unable to convert asset: %w", err)
		}
		gasAssetAmountWei := convertThorchainAmountToWei(gasAssetAmount.BigInt())
		confValue := common.GetUncappedShare(cosmos.NewUint(uint64(confMul)), cosmos.NewUint(constants.MaxBasisPts), cosmos.NewUintFromBigInt(big.NewInt(ethBlockRewardAndFee)))
		confirmations = int64(cosmos.NewUintFromBigInt(gasAssetAmountWei).MulUint64(2).Quo(confValue).Uint64())
	}

	// max confirmation adjustment for btc and eth
	if chain.Equals(common.BTCChain) || chain.Equals(common.ETHChain) {
		maxConfirmations, err := mgr.Keeper().GetMimirWithRef(ctx, constants.MimirTemplateMaxConfirmations, chain.String())
		if maxConfirmations < 0 || err != nil {
			maxConfirmations = 0
		}
		if maxConfirmations > 0 && confirmations > maxConfirmations {
			confirmations = maxConfirmations
		}
	}

	// min confirmation adjustment
	confFloor := map[common.Chain]int64{
		common.ETHChain:  2,
		common.DOGEChain: 2,
	}
	if floor := confFloor[chain]; confirmations < floor {
		confirmations = floor
	}

	return address, router, confirmations, nil
}

func quoteOutboundInfo(ctx cosmos.Context, mgr *Mgrs, coin common.Coin) (int64, error) {
	toi := TxOutItem{
		Memo: "OUT:-",
		Coin: coin,
	}
	outboundHeight, _, err := mgr.txOutStore.CalcTxOutHeight(ctx, mgr.GetVersion(), toi)
	if err != nil {
		return 0, err
	}
	return outboundHeight - ctx.BlockHeight(), nil
}

// -------------------------------------------------------------------------------------
// Swap
// -------------------------------------------------------------------------------------

// calculateMinSwapAmount returns the recommended minimum swap amount The recommended
// min swap amount is: - MAX(
//
//	  outbound_fee(src_chain) * 4,
//	  outbound_fee(dest_chain) * 4,
//	  (native_tx_fee_rune * 2) * 10,000 / affiliateBps
//	)
//
// The reason the base value is the MAX of the outbound fees of each chain is because if
// the swap is refunded the input amount will need to cover the outbound fee of the
// source chain. A 4x buffer is applied because outbound fees can spike quickly, meaning
// the original input amount could be less than the new outbound fee. If this happens
// and the swap is refunded, the refund will fail, and the user will lose the entire
// input amount. The min amount could also be determined by the affiliate bps of the
// swap. The affiliate bps of the input amount needs to be enough to cover the native tx fee for the
// affiliate swap to RUNE. In this case, we give a 2x buffer on the native_tx_fee so the
// affiliate receives some amount after the fee is deducted.
func calculateMinSwapAmount(ctx cosmos.Context, mgr *Mgrs, fromAsset, toAsset common.Asset, affiliateBps cosmos.Uint) (cosmos.Uint, error) {
	srcOutboundFee, err := mgr.GasMgr().GetAssetOutboundFee(ctx, fromAsset, false)
	if err != nil {
		return cosmos.ZeroUint(), fmt.Errorf("fail to get outbound fee for source chain gas asset %s: %w", fromAsset, err)
	}
	destOutboundFee, err := mgr.GasMgr().GetAssetOutboundFee(ctx, toAsset, false)
	if err != nil {
		return cosmos.ZeroUint(), fmt.Errorf("fail to get outbound fee for destination chain gas asset %s: %w", toAsset, err)
	}

	if fromAsset.GetChain().IsTHORChain() && toAsset.GetChain().IsTHORChain() {
		// If this is a purely THORChain swap, no need to give a 4x buffer since outbound fees do not change
		// 2x buffer should suffice
		return srcOutboundFee.Mul(cosmos.NewUint(2)), nil
	}

	destInSrcAsset, err := quoteConvertAsset(ctx, mgr, toAsset, destOutboundFee, fromAsset)
	if err != nil {
		return cosmos.ZeroUint(), fmt.Errorf("fail to convert dest fee to src asset %w", err)
	}

	minSwapAmount := srcOutboundFee
	if destInSrcAsset.GT(srcOutboundFee) {
		minSwapAmount = destInSrcAsset
	}

	minSwapAmount = minSwapAmount.Mul(cosmos.NewUint(4))

	if affiliateBps.GT(cosmos.ZeroUint()) {
		nativeTxFeeRune, err := mgr.GasMgr().GetAssetOutboundFee(ctx, common.RuneNative, true)
		if err != nil {
			return cosmos.ZeroUint(), fmt.Errorf("fail to get native tx fee for rune: %w", err)
		}
		affSwapAmountRune := nativeTxFeeRune.Mul(cosmos.NewUint(2))
		mainSwapAmountRune := affSwapAmountRune.Mul(cosmos.NewUint(10_000)).Quo(affiliateBps)

		mainSwapAmount, err := quoteConvertAsset(ctx, mgr, common.RuneAsset(), mainSwapAmountRune, fromAsset)
		if err != nil {
			return cosmos.ZeroUint(), fmt.Errorf("fail to convert main swap amount to src asset %w", err)
		}

		if mainSwapAmount.GT(minSwapAmount) {
			minSwapAmount = mainSwapAmount
		}
	}

	return minSwapAmount, nil
}

func queryQuoteSwap(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	// extract parameters
	params, err := quoteParseParams(req.Data)
	if err != nil {
		return quoteErrorResponse(err)
	}

	// validate required parameters
	for _, p := range []string{fromAssetParam, toAssetParam, amountParam} {
		if len(params[p]) == 0 {
			return quoteErrorResponse(fmt.Errorf("missing required parameter %s", p))
		}
	}

	// parse assets
	fromAsset, err := common.NewAssetWithShortCodes(mgr.GetVersion(), params[fromAssetParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad from asset: %w", err))
	}
	fromAsset = fuzzyAssetMatch(ctx, mgr.Keeper(), fromAsset)
	toAsset, err := common.NewAssetWithShortCodes(mgr.GetVersion(), params[toAssetParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad to asset: %w", err))
	}
	toAsset = fuzzyAssetMatch(ctx, mgr.Keeper(), toAsset)

	// parse amount
	amount, err := cosmos.ParseUint(params[amountParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad amount: %w", err))
	}

	// parse streaming interval
	streamingInterval := uint64(0) // default value
	if len(params[intervalParam]) > 0 {
		streamingInterval, err = strconv.ParseUint(params[intervalParam][0], 10, 64)
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("bad streaming interval amount: %w", err))
		}
	}
	streamingQuantity := uint64(0) // default value
	if len(params[quantityParam]) > 0 {
		streamingQuantity, err = strconv.ParseUint(params[quantityParam][0], 10, 64)
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("bad streaming quantity amount: %w", err))
		}
	}
	swp := StreamingSwap{
		Interval: streamingInterval,
		Deposit:  amount,
	}
	maxSwapQuantity, err := getMaxSwapQuantity(ctx, mgr, fromAsset, toAsset, swp)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("failed to calculate max streaming swap quantity: %w", err))
	}

	// cap the streaming quantity to the max swap quantity
	if streamingQuantity > maxSwapQuantity {
		streamingQuantity = maxSwapQuantity
	}

	// if from asset is a synth, transfer asset to asgard module
	if fromAsset.IsSyntheticAsset() {
		// mint required coins to asgard so swap can be simulated
		err = mgr.Keeper().MintToModule(ctx, ModuleName, common.NewCoin(fromAsset, amount))
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("failed to mint coins to module: %w", err))
		}

		err = mgr.Keeper().SendFromModuleToModule(ctx, ModuleName, AsgardName, common.NewCoins(common.NewCoin(fromAsset, amount)))
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("failed to send coins to asgard: %w", err))
		}
	}

	// parse affiliate
	affiliate, affiliateMemo, affiliateBps, swapAmount, affAmt, err := quoteHandleAffiliate(ctx, mgr, params, amount)
	if err != nil {
		return quoteErrorResponse(err)
	}

	// simulate/validate the affiliate swap
	if affAmt.GT(sdk.ZeroUint()) {
		if fromAsset.IsNativeRune() {
			fee := mgr.Keeper().GetNativeTxFee(ctx)
			if affAmt.LTE(fee) {
				return quoteErrorResponse(fmt.Errorf("affiliate amount must be greater than native fee %s", fee))
			}
		} else {
			// validate affiliate address
			affiliateSwapMsg := &types.MsgSwap{
				Tx: common.Tx{
					ID:          common.BlankTxID,
					Chain:       fromAsset.Chain,
					FromAddress: common.NoopAddress,
					ToAddress:   common.NoopAddress,
					Coins: []common.Coin{
						{
							Asset:  fromAsset,
							Amount: affAmt,
						},
					},
					Gas: []common.Coin{{
						Asset:  common.RuneAsset(),
						Amount: sdk.NewUint(1),
					}},
					Memo: "",
				},
				TargetAsset:          common.RuneAsset(),
				TradeTarget:          cosmos.ZeroUint(),
				Destination:          affiliate,
				AffiliateAddress:     common.NoAddress,
				AffiliateBasisPoints: cosmos.ZeroUint(),
			}

			nodeAccounts, err := mgr.Keeper().ListActiveValidators(ctx)
			if err != nil {
				return nil, fmt.Errorf("no active node accounts: %w", err)
			}
			affiliateSwapMsg.Signer = nodeAccounts[0].NodeAddress

			// simulate the swap
			_, err = simulateInternal(ctx, mgr, affiliateSwapMsg)
			if err != nil {
				return quoteErrorResponse(fmt.Errorf("affiliate swap failed: %w", err))
			}
		}
	}

	// parse destination address or generate a random one
	sendMemo := true
	var destination common.Address
	if len(params[destinationParam]) > 0 {
		destination, err = quoteParseAddress(ctx, mgr, params[destinationParam][0], toAsset.Chain)
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("bad destination address: %w", err))
		}

	} else {
		chain := common.THORChain
		if !toAsset.IsSyntheticAsset() {
			chain = toAsset.Chain
		}
		destination, err = types.GetRandomPubKey().GetAddress(chain)
		if err != nil {
			return nil, fmt.Errorf("failed to generate address: %w", err)
		}
		sendMemo = false // do not send memo if destination was random
	}

	// parse tolerance basis points
	limit := sdk.ZeroUint()
	if len(params[toleranceBasisPointsParam]) > 0 {
		// validate tolerance basis points
		toleranceBasisPoints, err := sdk.ParseUint(params[toleranceBasisPointsParam][0])
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("bad tolerance basis points: %w", err))
		}
		if toleranceBasisPoints.GT(sdk.NewUint(10000)) {
			return quoteErrorResponse(fmt.Errorf("tolerance basis points must be less than 10000"))
		}

		// convert to a limit of target asset amount assuming zero fees and slip
		feelessEmit, err := quoteConvertAsset(ctx, mgr, fromAsset, swapAmount, toAsset)
		if err != nil {
			return quoteErrorResponse(err)
		}

		limit = feelessEmit.MulUint64(10000 - toleranceBasisPoints.Uint64()).QuoUint64(10000)
	}

	// custom refund addr
	refundAddress := common.NoAddress
	if len(params[refundAddressParam]) > 0 {
		refundAddress, err = quoteParseAddress(ctx, mgr, params[refundAddressParam][0], fromAsset.Chain)
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("bad refund address: %w", err))
		}
	}

	// create the memo
	memo := &SwapMemo{
		MemoBase: mem.MemoBase{
			TxType: TxSwap,
			Asset:  toAsset,
		},
		Destination:          destination,
		SlipLimit:            limit,
		AffiliateAddress:     common.Address(affiliateMemo),
		AffiliateBasisPoints: affiliateBps,
		StreamInterval:       streamingInterval,
		StreamQuantity:       streamingQuantity,
		RefundAddress:        refundAddress,
	}

	// if from asset chain has memo length restrictions use a prefix
	memoString := memo.String()
	if !fromAsset.Synth && len(memoString) > fromAsset.Chain.MaxMemoLength() {
		if len(memo.ShortString()) < len(memoString) { // use short codes if available
			memoString = memo.ShortString()
		} else { // otherwise attempt to shorten
			fuzzyAsset, err := quoteReverseFuzzyAsset(ctx, mgr, toAsset)
			if err == nil {
				memo.Asset = fuzzyAsset
				memoString = memo.String()
			}
		}

		// this is the shortest we can make it
		if len(memoString) > fromAsset.Chain.MaxMemoLength() {
			return quoteErrorResponse(fmt.Errorf("generated memo too long for source chain"))
		}
	}

	// trade assets must have from address on the source tx
	fromChain := fromAsset.Chain
	if fromAsset.IsSyntheticAsset() || fromAsset.IsDerivedAsset() || fromAsset.IsTradeAsset() {
		fromChain = common.THORChain
	}
	fromPubkey := types.GetRandomPubKey()
	fromAddress, err := fromPubkey.GetAddress(fromChain)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad from address: %w", err))
	}

	// if from asset is a trade asset, create fake balance
	if fromAsset.IsTradeAsset() {
		thorAddr, err := fromPubkey.GetThorAddress()
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("failed to get thor address: %w", err))
		}
		_, err = mgr.TradeAccountManager().Deposit(ctx, fromAsset, amount, thorAddr, common.NoAddress, common.BlankTxID)
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("failed to deposit trade asset: %w", err))
		}
	}

	// create the swap message
	msg := &types.MsgSwap{
		Tx: common.Tx{
			ID:          common.BlankTxID,
			Chain:       fromAsset.Chain,
			FromAddress: fromAddress,
			ToAddress:   common.NoopAddress,
			Coins: []common.Coin{
				{
					Asset:  fromAsset,
					Amount: swapAmount,
				},
			},
			Gas: []common.Coin{{
				Asset:  common.RuneAsset(),
				Amount: sdk.NewUint(1),
			}},
			Memo: memoString,
		},
		TargetAsset:          toAsset,
		TradeTarget:          limit,
		Destination:          destination,
		AffiliateAddress:     affiliate,
		AffiliateBasisPoints: affiliateBps,
	}

	// simulate the swap
	res, emitAmount, outboundFeeAmount, err := quoteSimulateSwap(ctx, mgr, amount, msg, 1)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("failed to simulate swap: %w", err))
	}

	// if we're using a streaming swap, calculate emit amount by a sub-swap amount instead
	// of the full amount, then multiply the result by the swap count
	if streamingInterval > 0 && streamingQuantity == 0 {
		streamingQuantity = maxSwapQuantity
	}
	res.StreamingSlippageBps = res.SlippageBps
	if streamingInterval > 0 && streamingQuantity > 0 {
		msg.TradeTarget = msg.TradeTarget.QuoUint64(streamingQuantity)
		// simulate the swap
		var streamRes *openapi.QuoteSwapResponse
		streamRes, emitAmount, _, err = quoteSimulateSwap(ctx, mgr, amount, msg, streamingQuantity)
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("failed to simulate swap: %w", err))
		}
		res.StreamingSlippageBps = streamRes.SlippageBps
		res.Fees = streamRes.Fees
	}

	// TODO: After UIs have transitioned everything below the message definition above
	// should reduce to the following:
	//
	// if streamingInterval > 0 && streamingQuantity == 0 {
	//   streamingQuantity = maxSwapQuantity
	// }
	// if streamingInterval > 0 && streamingQuantity > 0 {
	//   msg.TradeTarget = msg.TradeTarget.QuoUint64(streamingQuantity)
	// }
	// res, emitAmount, outboundFeeAmount, err := quoteSimulateSwap(ctx, mgr, amount, msg, streamingQuantity)
	// if err != nil {
	//   return quoteErrorResponse(fmt.Errorf("failed to simulate swap: %w", err))
	// }

	// check invariant
	if emitAmount.LT(outboundFeeAmount) {
		return quoteErrorResponse(fmt.Errorf("invariant broken: emit %s less than outbound fee %s", emitAmount, outboundFeeAmount))
	}

	// the amount out will deduct the outbound fee
	res.ExpectedAmountOut = emitAmount.Sub(outboundFeeAmount).String()

	// TODO: temporary for transition to only use expected amount out.
	if streamingInterval > 0 && streamingQuantity > 0 {
		res.ExpectedAmountOutStreaming = res.ExpectedAmountOut
	}

	maxQ := int64(maxSwapQuantity)
	res.MaxStreamingQuantity = &maxQ
	var streamSwapBlocks int64
	if streamingQuantity > 0 {
		streamSwapBlocks = int64(streamingInterval) * int64(streamingQuantity-1)
	}
	res.StreamingSwapBlocks = &streamSwapBlocks
	res.StreamingSwapSeconds = wrapInt64(streamSwapBlocks * common.THORChain.ApproximateBlockMilliseconds() / 1000)

	// estimate the inbound info
	inboundAddress, routerAddress, inboundConfirmations, err := quoteInboundInfo(ctx, mgr, amount, fromAsset.GetChain(), fromAsset)
	if err != nil {
		return quoteErrorResponse(err)
	}
	res.InboundAddress = wrapString(inboundAddress.String())
	if inboundConfirmations > 0 {
		res.InboundConfirmationBlocks = wrapInt64(inboundConfirmations)
		res.InboundConfirmationSeconds = wrapInt64(inboundConfirmations * msg.Tx.Chain.ApproximateBlockMilliseconds() / 1000)
	}

	res.OutboundDelayBlocks = 0
	res.OutboundDelaySeconds = 0
	if !toAsset.Chain.IsTHORChain() {
		// estimate the outbound info
		outboundDelay, err := quoteOutboundInfo(ctx, mgr, common.Coin{Asset: toAsset, Amount: emitAmount})
		if err != nil {
			return quoteErrorResponse(err)
		}
		res.OutboundDelayBlocks = outboundDelay
		res.OutboundDelaySeconds = outboundDelay * common.THORChain.ApproximateBlockMilliseconds() / 1000
	}

	totalSeconds := res.OutboundDelaySeconds
	if res.StreamingSwapSeconds != nil && res.OutboundDelaySeconds < *res.StreamingSwapSeconds {
		totalSeconds = *res.StreamingSwapSeconds
	}
	if inboundConfirmations > 0 {
		totalSeconds += *res.InboundConfirmationSeconds
	}
	res.TotalSwapSeconds = wrapInt64(totalSeconds)

	// send memo if the destination was provided
	if sendMemo {
		res.Memo = wrapString(memoString)
	}

	// set info fields
	if fromAsset.Chain.IsEVM() {
		res.Router = wrapString(routerAddress.String())
	}
	if !fromAsset.Chain.DustThreshold().IsZero() {
		res.DustThreshold = wrapString(fromAsset.Chain.DustThreshold().String())
	}

	res.Notes = fromAsset.GetChain().InboundNotes()
	res.Warning = quoteWarning
	res.Expiry = time.Now().Add(quoteExpiration).Unix()
	minSwapAmount, err := calculateMinSwapAmount(ctx, mgr, fromAsset, toAsset, affiliateBps)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("Failed to calculate min amount in: %s", err.Error()))
	}
	res.RecommendedMinAmountIn = wrapString(minSwapAmount.String())

	// set inbound recommended gas for non-native swaps
	if !fromAsset.Chain.IsTHORChain() {
		inboundGas := mgr.GasMgr().GetGasRate(ctx, fromAsset.Chain)
		res.RecommendedGasRate = wrapString(inboundGas.String())
		res.GasRateUnits = wrapString(fromAsset.Chain.GetGasUnits())
	}

	return json.MarshalIndent(res, "", "  ")
}

// -------------------------------------------------------------------------------------
// Saver Deposit
// -------------------------------------------------------------------------------------

func queryQuoteSaverDeposit(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	// extract parameters
	params, err := quoteParseParams(req.Data)
	if err != nil {
		return quoteErrorResponse(err)
	}

	// validate required parameters
	for _, p := range []string{assetParam, amountParam} {
		if len(params[p]) == 0 {
			return quoteErrorResponse(fmt.Errorf("missing required parameter %s", p))
		}
	}

	// parse asset
	asset, err := common.NewAssetWithShortCodes(mgr.GetVersion(), params[assetParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad asset: %w", err))
	}
	asset = fuzzyAssetMatch(ctx, mgr.Keeper(), asset)

	// parse amount
	amount, err := cosmos.ParseUint(params[amountParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad amount: %w", err))
	}

	// parse affiliate
	affiliate, affiliateMemo, affiliateBps, depositAmount, _, err := quoteHandleAffiliate(ctx, mgr, params, amount)
	if err != nil {
		return quoteErrorResponse(err)
	}

	// generate deposit memo
	depositMemoComponents := []string{
		"+",
		asset.GetSyntheticAsset().String(),
		"",
		affiliateMemo,
		affiliateBps.String(),
	}
	depositMemo := strings.Join(depositMemoComponents[:2], ":")
	if affiliate != common.NoAddress && !affiliateBps.IsZero() {
		depositMemo = strings.Join(depositMemoComponents, ":")
	}

	q := url.Values{}
	q.Add("from_asset", asset.String())
	q.Add("to_asset", asset.GetSyntheticAsset().String())
	q.Add("amount", depositAmount.String())
	q.Add("destination", string(GetRandomTHORAddress())) // required param, not actually used, spoof it

	ssInterval := mgr.Keeper().GetConfigInt64(ctx, constants.SaversStreamingSwapsInterval)
	if ssInterval > 0 {
		q.Add("streaming_interval", fmt.Sprintf("%d", ssInterval))
		q.Add("streaming_quantity", fmt.Sprintf("%d", 0))
	}

	swapReq := abci.RequestQuery{Data: []byte("/thorchain/quote/swap?" + q.Encode())}
	swapResRaw, err := queryQuoteSwap(ctx, nil, swapReq, mgr)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("unable to queryQuoteSwap: %w", err))
	}

	var swapRes *openapi.QuoteSwapResponse
	err = json.Unmarshal(swapResRaw, &swapRes)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("unable to unmarshal swapRes: %w", err))
	}

	expectedAmountOut, _ := sdk.ParseUint(swapRes.ExpectedAmountOut)
	outboundFee, _ := sdk.ParseUint(*swapRes.Fees.Outbound)
	depositAmount = expectedAmountOut.Add(outboundFee)

	// use the swap result info to generate the deposit quote
	res := &openapi.QuoteSaverDepositResponse{
		// TODO: deprecate ExpectedAmountOut in future version
		ExpectedAmountOut:          wrapString(depositAmount.String()),
		ExpectedAmountDeposit:      depositAmount.String(),
		Fees:                       swapRes.Fees,
		SlippageBps:                swapRes.SlippageBps,
		InboundConfirmationBlocks:  swapRes.InboundConfirmationBlocks,
		InboundConfirmationSeconds: swapRes.InboundConfirmationSeconds,
		Memo:                       depositMemo,
	}

	// estimate the inbound info
	inboundAddress, _, inboundConfirmations, err := quoteInboundInfo(ctx, mgr, amount, asset.GetLayer1Asset().Chain, asset)
	if err != nil {
		return quoteErrorResponse(err)
	}
	res.InboundAddress = inboundAddress.String()
	res.InboundConfirmationBlocks = wrapInt64(inboundConfirmations)

	// set info fields
	chain := asset.GetLayer1Asset().Chain
	if !chain.DustThreshold().IsZero() {
		res.DustThreshold = wrapString(chain.DustThreshold().String())
		res.RecommendedMinAmountIn = res.DustThreshold
	}
	res.Notes = chain.InboundNotes()
	res.Warning = quoteWarning
	res.Expiry = time.Now().Add(quoteExpiration).Unix()

	// set inbound recommended gas
	inboundGas := mgr.GasMgr().GetGasRate(ctx, chain)
	res.RecommendedGasRate = inboundGas.String()
	res.GasRateUnits = chain.GetGasUnits()

	return json.MarshalIndent(res, "", "  ")
}

// -------------------------------------------------------------------------------------
// Saver Withdraw
// -------------------------------------------------------------------------------------

func queryQuoteSaverWithdraw(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	// extract parameters
	params, err := quoteParseParams(req.Data)
	if err != nil {
		return quoteErrorResponse(err)
	}

	// validate required parameters
	for _, p := range []string{assetParam, addressParam, withdrawBasisPointsParam} {
		if len(params[p]) == 0 {
			return quoteErrorResponse(fmt.Errorf("missing required parameter %s", p))
		}
	}

	// parse asset
	asset, err := common.NewAssetWithShortCodes(mgr.GetVersion(), params[assetParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad asset: %w", err))
	}
	asset = fuzzyAssetMatch(ctx, mgr.Keeper(), asset)
	asset = asset.GetSyntheticAsset() // always use the vault asset

	// parse address
	address, err := common.NewAddress(params[addressParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad address: %w", err))
	}

	// parse basis points
	basisPoints, err := cosmos.ParseUint(params[withdrawBasisPointsParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad basis points: %w", err))
	}

	// validate basis points
	if basisPoints.GT(sdk.NewUint(10_000)) {
		return quoteErrorResponse(fmt.Errorf("basis points must be less than 10000"))
	}

	// get liquidity provider
	lp, err := mgr.Keeper().GetLiquidityProvider(ctx, asset, address)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("failed to get liquidity provider: %w", err))
	}

	// get the pool
	pool, err := mgr.Keeper().GetPool(ctx, asset)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("failed to get pool: %w", err))
	}

	// get the liquidity provider share of the pool
	lpShare := lp.GetSaversAssetRedeemValue(pool)

	// calculate the withdraw amount
	amount := common.GetSafeShare(basisPoints, sdk.NewUint(10_000), lpShare)

	q := url.Values{}
	q.Add("from_asset", asset.String())
	q.Add("to_asset", asset.GetLayer1Asset().String())
	q.Add("amount", amount.String())
	q.Add("destination", address.String()) // required param, not actually used, spoof it

	ssInterval := mgr.Keeper().GetConfigInt64(ctx, constants.SaversStreamingSwapsInterval)
	if ssInterval > 0 {
		q.Add("streaming_interval", fmt.Sprintf("%d", ssInterval))
		q.Add("streaming_quantity", fmt.Sprintf("%d", 0))
	}

	swapReq := abci.RequestQuery{Data: []byte("/thorchain/quote/swap?" + q.Encode())}
	swapResRaw, err := queryQuoteSwap(ctx, nil, swapReq, mgr)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("unable to queryQuoteSwap: %w", err))
	}

	var swapRes *openapi.QuoteSwapResponse
	err = json.Unmarshal(swapResRaw, &swapRes)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("unable to unmarshal swapRes: %w", err))
	}

	// use the swap result info to generate the withdraw quote
	res := &openapi.QuoteSaverWithdrawResponse{
		ExpectedAmountOut: swapRes.ExpectedAmountOut,
		Fees:              swapRes.Fees,
		SlippageBps:       swapRes.SlippageBps,
		Memo:              fmt.Sprintf("-:%s:%s", asset.String(), basisPoints.String()),
		DustAmount:        asset.GetLayer1Asset().Chain.DustThreshold().Add(basisPoints).String(),
	}

	// estimate the inbound info
	inboundAddress, _, _, err := quoteInboundInfo(ctx, mgr, amount, asset.GetLayer1Asset().Chain, asset)
	if err != nil {
		return quoteErrorResponse(err)
	}
	res.InboundAddress = inboundAddress.String()

	// estimate the outbound info
	expectedAmountOut, _ := sdk.ParseUint(swapRes.ExpectedAmountOut)
	outboundCoin := common.Coin{Asset: asset.GetLayer1Asset(), Amount: expectedAmountOut}
	outboundDelay, err := quoteOutboundInfo(ctx, mgr, outboundCoin)
	if err != nil {
		return quoteErrorResponse(err)
	}
	res.OutboundDelayBlocks = outboundDelay
	res.OutboundDelaySeconds = outboundDelay * common.THORChain.ApproximateBlockMilliseconds() / 1000

	// set info fields
	chain := asset.GetLayer1Asset().Chain
	if !chain.DustThreshold().IsZero() {
		res.DustThreshold = wrapString(chain.DustThreshold().String())
	}
	res.Notes = chain.InboundNotes()
	res.Warning = quoteWarning
	res.Expiry = time.Now().Add(quoteExpiration).Unix()

	// set inbound recommended gas
	inboundGas := mgr.GasMgr().GetGasRate(ctx, chain)
	res.RecommendedGasRate = inboundGas.String()
	res.GasRateUnits = chain.GetGasUnits()

	return json.MarshalIndent(res, "", "  ")
}

// -------------------------------------------------------------------------------------
// Loan Open
// -------------------------------------------------------------------------------------

func queryQuoteLoanOpen(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	// extract parameters
	params, err := quoteParseParams(req.Data)
	if err != nil {
		return quoteErrorResponse(err)
	}

	// validate required parameters
	for _, p := range []string{fromAssetParam, amountParam, toAssetParam} {
		if len(params[p]) == 0 {
			return quoteErrorResponse(fmt.Errorf("missing required parameter %s", p))
		}
	}

	// invalidate unexpected parameters
	allowed := map[string]bool{
		heightParam:       true,
		fromAssetParam:    true,
		amountParam:       true,
		minOutParam:       true,
		toAssetParam:      true,
		destinationParam:  true,
		affiliateParam:    true,
		affiliateBpsParam: true,
	}
	for p := range params {
		if !allowed[p] {
			return quoteErrorResponse(fmt.Errorf("unexpected parameter %s", p))
		}
	}

	// parse asset
	asset, err := common.NewAssetWithShortCodes(mgr.GetVersion(), params[fromAssetParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad asset: %w", err))
	}
	asset = fuzzyAssetMatch(ctx, mgr.Keeper(), asset)

	// parse amount
	amount, err := cosmos.ParseUint(params[amountParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad amount: %w", err))
	}

	// parse min out
	minOut := sdk.ZeroUint()
	if len(params[minOutParam]) > 0 {
		minOut, err = cosmos.ParseUint(params[minOutParam][0])
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("bad min out: %w", err))
		}
	}

	// Affiliate fee in RUNE
	affiliateRuneAmt := sdk.ZeroUint()

	// parse affiliate
	affiliate, affiliateMemo, affiliateBps, amt, affiliateAmt, err := quoteHandleAffiliate(ctx, mgr, params, amount)
	if err != nil {
		return quoteErrorResponse(err)
	}

	if affiliate != common.NoAddress && !affiliateBps.IsZero() {
		affCoin := common.NewCoin(asset, affiliateAmt)
		gasCoin := common.NewCoin(asset.GetChain().GetGasAsset(), cosmos.OneUint())
		fakeTx := common.NewTx(common.BlankTxID, common.NoopAddress, common.NoopAddress, common.NewCoins(affCoin), common.Gas{gasCoin}, "noop")
		affiliateSwap := NewMsgSwap(fakeTx, common.RuneAsset(), affiliate, cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(), "", "", nil, 0, 0, 0, nil)

		_, affiliateRuneAmt, _, err = quoteSimulateSwap(ctx, mgr, affiliateAmt, affiliateSwap, 1)
		if err == nil {
			// skim fee off collateral amount
			amount = amt
		} else {
			affiliateRuneAmt = sdk.ZeroUint()
		}
	}

	// parse target asset
	targetAsset, err := common.NewAssetWithShortCodes(mgr.GetVersion(), params[toAssetParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad target asset: %w", err))
	}
	targetAsset = fuzzyAssetMatch(ctx, mgr.Keeper(), targetAsset)

	// parse destination address or generate a random one
	sendMemo := true
	var destination common.Address
	if len(params[destinationParam]) > 0 {
		destination, err = quoteParseAddress(ctx, mgr, params[destinationParam][0], targetAsset.Chain)
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("bad destination address: %w", err))
		}

	} else {
		destination, err = types.GetRandomPubKey().GetAddress(targetAsset.Chain)
		if err != nil {
			return nil, fmt.Errorf("failed to generate address: %w", err)
		}
		sendMemo = false // do not send memo if destination was random
	}

	// check that destination and affiliate are not the same
	if destination.Equals(affiliate) {
		return quoteErrorResponse(fmt.Errorf("destination and affiliate should not be the same"))
	}

	// generate random address for collateral owner
	collateralOwner, err := types.GetRandomPubKey().GetAddress(asset.Chain)
	if err != nil {
		return nil, fmt.Errorf("failed to generate address: %w", err)
	}

	// create message for simulation
	msg := &types.MsgLoanOpen{
		Owner:            collateralOwner,
		CollateralAsset:  asset,
		CollateralAmount: amount,
		TargetAddress:    destination,
		TargetAsset:      targetAsset,
		MinOut:           minOut,

		// We calculate the affiliate fee manually as handler_open_loan expects a TxVoter to
		// get the affiliate params from the memo
		AffiliateBasisPoints: cosmos.ZeroUint(),

		// TODO: support aggregator
		Aggregator:              "",
		AggregatorTargetAddress: "",
		AggregatorTargetLimit:   sdk.ZeroUint(),
	}

	// simulate message handling
	events, err := simulate(ctx, mgr, msg)
	if err != nil {
		return quoteErrorResponse(err)
	}

	// create response
	res := &openapi.QuoteLoanOpenResponse{
		Fees: openapi.QuoteFees{
			Asset: targetAsset.String(),
		},
		Expiry:  time.Now().Add(quoteExpiration).Unix(),
		Warning: quoteWarning,
		Notes:   asset.Chain.InboundNotes(),
	}

	// estimate the inbound info
	inboundAddress, routerAddress, inboundConfirmations, err := quoteInboundInfo(ctx, mgr, amount, asset.Chain, asset)
	if err != nil {
		return quoteErrorResponse(err)
	}
	res.InboundAddress = wrapString(inboundAddress.String())
	if inboundConfirmations > 0 {
		res.InboundConfirmationBlocks = wrapInt64(inboundConfirmations)
		res.InboundConfirmationSeconds = wrapInt64(inboundConfirmations * asset.Chain.ApproximateBlockMilliseconds() / 1000)
	}

	// set info fields
	if asset.Chain.IsEVM() {
		res.Router = wrapString(routerAddress.String())
	}
	if !asset.Chain.DustThreshold().IsZero() {
		res.DustThreshold = wrapString(asset.Chain.DustThreshold().String())
	}

	// sum liquidity fees in rune from all swap events
	outboundFee := sdk.ZeroUint()
	liquidityFee := sdk.ZeroUint()
	affiliateFee := affiliateRuneAmt
	expectedAmountOut := sdk.ZeroUint()
	finalEmitAmount := sdk.ZeroUint() // used to calculate slippage
	streamingSwapBlocks := int64(0)
	streamingSwapSeconds := int64(0)

	// iterate events in reverse order
	for i := len(events) - 1; i >= 0; i-- {
		e := events[i]
		em := eventMap(e)

		switch e.Type {

		// use final outbound event as expected amount - scheduled_outbound (L1) or outbound (native)
		case "scheduled_outbound":
			if res.ExpectedAmountOut == "" { // if not empty we already saw the last outbound event
				res.ExpectedAmountOut = em["coin_amount"]
				expectedAmountOut = sdk.NewUintFromString(em["coin_amount"])
				if em["coin_asset"] != targetAsset.String() { // should be unreachable
					return quoteErrorResponse(fmt.Errorf("unexpected outbound asset: %s", em["coin_asset"]))
				}

				// estimate the outbound info
				outboundDelay, err := quoteOutboundInfo(ctx, mgr, common.NewCoin(targetAsset, sdk.NewUintFromString(res.ExpectedAmountOut)))
				if err != nil {
					return quoteErrorResponse(err)
				}
				res.OutboundDelayBlocks = outboundDelay
				res.OutboundDelaySeconds = outboundDelay * common.THORChain.ApproximateBlockMilliseconds() / 1000
			}
		case "outbound":
			coin, err := common.ParseCoin(em["coin"])
			if err != nil {
				return quoteErrorResponse(fmt.Errorf("failed to parse coin: %w", err))
			}
			toAddress, _ := common.NewAddress(em["to"])

			// check for the outbound event
			if toAddress.Equals(destination) {
				res.ExpectedAmountOut = coin.Amount.String()
				expectedAmountOut = coin.Amount

				if !coin.Asset.Equals(targetAsset) { // should be unreachable
					return quoteErrorResponse(fmt.Errorf("unexpected outbound asset: %s", coin.Asset))
				}
			}

		// sum liquidity fee in rune for all swap events
		case "swap":
			liquidityFee = liquidityFee.Add(sdk.NewUintFromString(em["liquidity_fee_in_rune"]))
			coin, err := common.ParseCoin(em["emit_asset"])
			if err != nil {
				return quoteErrorResponse(fmt.Errorf("failed to parse coin: %w", err))
			}
			if coin.Asset.Equals(targetAsset) {
				finalEmitAmount = coin.Amount
			}
			swapQuantity, err := cosmos.ParseUint(em["streaming_swap_quantity"])
			if err != nil {
				return quoteErrorResponse(fmt.Errorf("bad quantity: %w", err))
			}
			streamingSwapBlocks += swapQuantity.BigInt().Int64()

		// extract loan data from loan open event
		case "loan_open":
			res.ExpectedCollateralizationRatio = em["collateralization_ratio"]
			res.ExpectedCollateralDeposited = em["collateral_deposited"]
			res.ExpectedDebtIssued = em["debt_issued"]

		// catch refund if there was an issue
		case "refund":
			if em["reason"] != "" {
				return quoteErrorResponse(fmt.Errorf("failed to simulate swap: %s", em["reason"]))
			}

		// set outbound fee from fee event
		case "fee":
			coin, err := common.ParseCoin(em["coins"])
			if err != nil {
				return quoteErrorResponse(fmt.Errorf("failed to parse coin: %w", err))
			}
			res.Fees.Outbound = wrapString(coin.Amount.String()) // already in target asset
			res.Fees.Asset = coin.Asset.String()
			outboundFee = coin.Amount

			if !coin.Asset.Equals(targetAsset) { // should be unreachable
				return quoteErrorResponse(fmt.Errorf("unexpected fee asset: %s", coin.Asset))
			}
		}
	}

	// convert fees to target asset if it is not rune
	if !targetAsset.Equals(common.RuneNative) {
		targetPool, err := mgr.Keeper().GetPool(ctx, targetAsset)
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("failed to get pool: %w", err))
		}
		affiliateFee = targetPool.RuneValueInAsset(affiliateRuneAmt)
		liquidityFee = targetPool.RuneValueInAsset(liquidityFee)
	}
	slippageBps := liquidityFee.MulUint64(10000).Quo(finalEmitAmount.Add(liquidityFee))

	// set fee info
	res.Fees.Liquidity = liquidityFee.String()
	totalFees := liquidityFee.Add(outboundFee).Add(affiliateFee)
	res.Fees.Total = totalFees.String()
	res.Fees.SlippageBps = slippageBps.BigInt().Int64()
	res.Fees.TotalBps = totalFees.MulUint64(10000).Quo(expectedAmountOut.Add(totalFees)).BigInt().Int64()
	if !affiliateFee.IsZero() {
		res.Fees.Affiliate = wrapString(affiliateFee.String())
	}

	// generate memo
	if sendMemo {
		memo := &mem.LoanOpenMemo{
			MemoBase: mem.MemoBase{
				TxType: TxLoanOpen,
			},
			TargetAsset:          targetAsset,
			TargetAddress:        destination,
			MinOut:               minOut,
			AffiliateAddress:     common.Address(affiliateMemo),
			AffiliateBasisPoints: affiliateBps,
			DexTargetLimit:       sdk.ZeroUint(),
		}

		// if from asset chain has memo length restrictions use a prefix
		memoString := memo.String()
		if len(memoString) > asset.Chain.MaxMemoLength() {
			if len(memo.ShortString()) < len(memoString) { // use short codes if available
				memoString = memo.ShortString()
			} else { // otherwise attempt to shorten
				fuzzyAsset, err := quoteReverseFuzzyAsset(ctx, mgr, targetAsset)
				if err == nil {
					memo.TargetAsset = fuzzyAsset
					memoString = memo.String()
				}
			}

			// this is the shortest we can make it
			if len(memoString) > asset.Chain.MaxMemoLength() {
				return quoteErrorResponse(fmt.Errorf("generated memo too long for source chain"))
			}
		}

		res.Memo = wrapString(memoString)
	}

	minLoanOpenAmount, err := calculateMinSwapAmount(ctx, mgr, asset, targetAsset, cosmos.ZeroUint())
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("Failed to calculate min amount in: %s", err.Error()))
	}
	res.RecommendedMinAmountIn = wrapString(minLoanOpenAmount.String())

	streamingSwapSeconds += streamingSwapBlocks * common.THORChain.ApproximateBlockMilliseconds() / 1000

	if res.InboundConfirmationSeconds != nil {
		value := *res.InboundConfirmationSeconds
		res.TotalOpenLoanSeconds = streamingSwapSeconds + res.OutboundDelaySeconds + value
	} else {
		res.TotalOpenLoanSeconds = streamingSwapSeconds + res.OutboundDelaySeconds
	}

	res.StreamingSwapBlocks = streamingSwapBlocks
	res.StreamingSwapSeconds = streamingSwapSeconds

	// set inbound recommended gas
	inboundGas := mgr.GasMgr().GetGasRate(ctx, asset.Chain)
	res.RecommendedGasRate = inboundGas.String()
	res.GasRateUnits = asset.Chain.GetGasUnits()

	return json.MarshalIndent(res, "", "  ")
}

// -------------------------------------------------------------------------------------
// Loan Close
// -------------------------------------------------------------------------------------

func quoteSimulateCloseLoan(ctx cosmos.Context, mgr *Mgrs, msg *MsgLoanRepayment) (
	res *openapi.QuoteLoanCloseResponse, data []byte, err error,
) {
	res = &openapi.QuoteLoanCloseResponse{
		Fees: openapi.QuoteFees{
			Asset: msg.CollateralAsset.String(),
		},
		Expiry:  time.Now().Add(quoteExpiration).Unix(),
		Warning: quoteWarning,
		Notes:   msg.Coin.Asset.Chain.InboundNotes(),
	}

	// simulate message handling
	events, err := simulate(ctx, mgr, msg)
	if err != nil {
		data, err = quoteErrorResponse(err)
		return nil, data, err
	}

	// estimate the inbound info
	inboundAddress, routerAddress, inboundConfirmations, err := quoteInboundInfo(ctx, mgr, msg.Coin.Amount, msg.Coin.Asset.GetChain(), msg.Coin.Asset)
	if err != nil {
		data, err = quoteErrorResponse(err)
		return nil, data, err
	}
	res.InboundAddress = wrapString(inboundAddress.String())
	if inboundConfirmations > 0 {
		res.InboundConfirmationBlocks = wrapInt64(inboundConfirmations)
		res.InboundConfirmationSeconds = wrapInt64(inboundConfirmations * msg.Coin.Asset.GetChain().ApproximateBlockMilliseconds() / 1000)
	}

	// set info fields
	if msg.Coin.Asset.Chain.IsEVM() {
		res.Router = wrapString(routerAddress.String())
	}
	if !msg.Coin.Asset.Chain.DustThreshold().IsZero() {
		res.DustThreshold = wrapString(msg.Coin.Asset.Chain.DustThreshold().String())
	}

	// sum liquidity fees in rune from all swap events
	outboundFee := sdk.ZeroUint()
	repaymentLiquidityFee := sdk.ZeroUint()
	outboundLiquidityFee := sdk.ZeroUint()
	affiliateFee := sdk.ZeroUint()
	expectedAmountOut := sdk.ZeroUint()
	streamingSwapBlocks := int64(0)
	streamingSwapSeconds := int64(0)
	var repaymentEmit, outboundEmit common.Coin

	// iterate events in reverse order
	for i := len(events) - 1; i >= 0; i-- {
		e := events[i]
		em := eventMap(e)

		switch e.Type {

		// use final outbound event as expected amount - scheduled_outbound (L1) or outbound (native)
		case "scheduled_outbound":
			if res.ExpectedAmountOut == "" { // if not empty we already saw the last outbound event
				res.ExpectedAmountOut = em["coin_amount"]
				expectedAmountOut = sdk.NewUintFromString(em["coin_amount"])
				if em["coin_asset"] != msg.CollateralAsset.String() { // should be unreachable
					data, err = quoteErrorResponse(fmt.Errorf("unexpected outbound asset: %s", em["coin_asset"]))
					return nil, data, err
				}

				// estimate the outbound info
				outboundDelay, err := quoteOutboundInfo(ctx, mgr, common.NewCoin(msg.CollateralAsset, sdk.NewUintFromString(res.ExpectedAmountOut)))
				if err != nil {
					data, err = quoteErrorResponse(err)
					return nil, data, err
				}
				res.OutboundDelayBlocks = outboundDelay
				res.OutboundDelaySeconds = outboundDelay * common.THORChain.ApproximateBlockMilliseconds() / 1000
			}
		case "outbound":
			// track coin and to address
			coin, err := common.ParseCoin(em["coin"])
			if err != nil {
				data, err = quoteErrorResponse(fmt.Errorf("failed to parse coin: %w", err))
				return nil, data, err
			}
			toAddress, _ := common.NewAddress(em["to"])

			// check for the outbound event
			if toAddress.Equals(msg.Owner) {
				res.ExpectedAmountOut = coin.Amount.String()
				expectedAmountOut = coin.Amount

				if !coin.Asset.Equals(msg.CollateralAsset) { // should be unreachable
					data, err = quoteErrorResponse(fmt.Errorf("unexpected outbound asset: %s", coin.Asset))
					return nil, data, err
				}
			}

		// sum liquidity fee in rune for all swap events
		case "swap":
			coin, err := common.ParseCoin(em["emit_asset"])
			if err != nil {
				data, err = quoteErrorResponse(fmt.Errorf("failed to parse coin: %w", err))
				return nil, data, err
			}
			swapQuantity, err := cosmos.ParseUint(em["streaming_swap_quantity"])
			if err != nil {
				data, err = quoteErrorResponse(fmt.Errorf("bad amount: %w", err))
				return nil, data, err
			}
			streamingSwapBlocks += swapQuantity.BigInt().Int64()
			switch {
			case coin.Asset.Equals(common.TOR):
				repaymentEmit = coin
				repaymentLiquidityFee = repaymentLiquidityFee.Add(sdk.NewUintFromString(em["liquidity_fee_in_rune"]))
			case !coin.Asset.IsNativeRune():
				outboundEmit = coin
				outboundLiquidityFee = outboundLiquidityFee.Add(sdk.NewUintFromString(em["liquidity_fee_in_rune"]))
			default:
				inCoin, err := common.ParseCoin(em["coin"])
				if err != nil {
					data, err = quoteErrorResponse(fmt.Errorf("failed to parse coin: %w", err))
					return nil, data, err
				}
				if inCoin.Asset.IsDerivedAsset() {
					outboundLiquidityFee = outboundLiquidityFee.Add(sdk.NewUintFromString(em["liquidity_fee_in_rune"]))
				} else {
					repaymentLiquidityFee = repaymentLiquidityFee.Add(sdk.NewUintFromString(em["liquidity_fee_in_rune"]))
				}
			}

		// extract loan data from loan close event
		case "loan_repayment":
			res.ExpectedCollateralWithdrawn = em["collateral_withdrawn"]
			res.ExpectedDebtRepaid = em["debt_repaid"]

		// catch refund if there was an issue
		case "refund":
			if em["reason"] != "" {
				data, err = quoteErrorResponse(fmt.Errorf("failed to simulate loan close: %s", em["reason"]))
				return nil, data, err
			}

		// set outbound fee from fee event
		case "fee":
			coin, err := common.ParseCoin(em["coins"])
			if err != nil {
				data, err = quoteErrorResponse(fmt.Errorf("failed to parse coin: %w", err))
				return nil, data, err
			}
			res.Fees.Outbound = wrapString(coin.Amount.String()) // already in collateral asset
			res.Fees.Asset = coin.Asset.String()
			outboundFee = coin.Amount

			if !coin.Asset.Equals(msg.CollateralAsset) { // should be unreachable
				data, err = quoteErrorResponse(fmt.Errorf("unexpected fee asset: %s", coin.Asset))
				return nil, data, err
			}

		}
	}

	// calculate emit values in rune
	torPool, err := mgr.Keeper().GetPool(ctx, common.TOR)
	if err != nil {
		data, err = quoteErrorResponse(fmt.Errorf("failed to get pool: %w", err))
		return nil, data, err
	}
	repaymentEmitRune := torPool.RuneValueInAsset(repaymentEmit.Amount)
	outPool, err := mgr.Keeper().GetPool(ctx, outboundEmit.Asset)
	if err != nil {
		data, err = quoteErrorResponse(fmt.Errorf("failed to get pool: %w", err))
		return nil, data, err
	}
	outboundEmitRune := outPool.RuneValueInAsset(outboundEmit.Amount)

	// slippage calculation is weighted to repayment and outbound amounts
	outboundSlip := sdk.ZeroUint()
	if !outboundEmitRune.IsZero() {
		outboundSlip = outboundLiquidityFee.MulUint64(10000).Quo(outboundEmitRune.Add(outboundLiquidityFee))
	}
	repaymentSlip := repaymentLiquidityFee.MulUint64(10000).Quo(repaymentEmitRune.Add(repaymentLiquidityFee))
	slippageBps := repaymentSlip.Mul(repaymentEmitRune).Add(outboundSlip.Mul(outboundEmitRune)).Quo(repaymentEmitRune.Add(outboundEmitRune))

	// convert fees to target asset if it is not rune
	liquidityFee := repaymentLiquidityFee.Add(outboundLiquidityFee)
	if !msg.CollateralAsset.Equals(common.RuneNative) {
		loanPool, err := mgr.Keeper().GetPool(ctx, msg.CollateralAsset)
		if err != nil {
			data, err = quoteErrorResponse(fmt.Errorf("failed to get pool: %w", err))
			return nil, data, err
		}
		affiliateFee = loanPool.RuneValueInAsset(affiliateFee)
		liquidityFee = loanPool.RuneValueInAsset(liquidityFee)
	}

	// set fee info
	res.Fees.Liquidity = liquidityFee.String()
	totalFees := liquidityFee.Add(outboundFee).Add(affiliateFee)
	res.Fees.Total = totalFees.String()
	res.Fees.SlippageBps = slippageBps.BigInt().Int64()
	if !expectedAmountOut.IsZero() {
		res.Fees.TotalBps = totalFees.MulUint64(10000).Quo(expectedAmountOut).BigInt().Int64()
	} else {
		res.Fees.TotalBps = res.Fees.SlippageBps
	}
	if !affiliateFee.IsZero() {
		res.Fees.Affiliate = wrapString(affiliateFee.String())
	}

	// generate memo
	memo := &mem.LoanRepaymentMemo{
		MemoBase: mem.MemoBase{
			TxType: TxLoanRepayment,
			Asset:  msg.CollateralAsset,
		},
		Owner:  msg.Owner,
		MinOut: msg.MinOut,
	}
	res.Memo = memo.String()

	minLoanCloseAmount, err := calculateMinSwapAmount(ctx, mgr, msg.Coin.Asset, msg.CollateralAsset, cosmos.ZeroUint())
	if err != nil {
		data, err = quoteErrorResponse(fmt.Errorf("Failed to calculate min amount in: %s", err.Error()))
		return nil, data, err
	}
	res.RecommendedMinAmountIn = wrapString(minLoanCloseAmount.String())

	streamingSwapSeconds += streamingSwapBlocks * common.THORChain.ApproximateBlockMilliseconds() / 1000

	if res.InboundConfirmationSeconds != nil {
		value := *res.InboundConfirmationSeconds
		res.TotalRepaySeconds = streamingSwapSeconds + res.OutboundDelaySeconds + value
	} else {
		res.TotalRepaySeconds = streamingSwapSeconds + res.OutboundDelaySeconds
	}

	res.StreamingSwapBlocks = streamingSwapBlocks
	res.StreamingSwapSeconds = streamingSwapSeconds
	res.ExpectedAmountIn = msg.Coin.Amount.String()

	return res, nil, nil
}

func queryQuoteLoanClose(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	// extract parameters
	params, err := quoteParseParams(req.Data)
	if err != nil {
		return quoteErrorResponse(err)
	}

	// validate required parameters
	for _, p := range []string{fromAssetParam, repayBpsParam, toAssetParam, loanOwnerParam} {
		if len(params[p]) == 0 {
			return quoteErrorResponse(fmt.Errorf("missing required parameter %s", p))
		}
	}

	// invalidate unexpected parameters
	allowed := map[string]bool{
		heightParam:    true,
		fromAssetParam: true,
		repayBpsParam:  true,
		toAssetParam:   true,
		loanOwnerParam: true,
		minOutParam:    true,
	}
	for p := range params {
		if !allowed[p] {
			return quoteErrorResponse(fmt.Errorf("unexpected parameter %s", p))
		}
	}

	// parse asset
	asset, err := common.NewAssetWithShortCodes(mgr.GetVersion(), params[fromAssetParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad asset: %w", err))
	}
	asset = fuzzyAssetMatch(ctx, mgr.Keeper(), asset)

	// parse repayment bps
	repayBps, err := cosmos.ParseUint(params[repayBpsParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad amount: %w", err))
	}

	// parse min out
	minOut := sdk.ZeroUint()
	if len(params[minOutParam]) > 0 {
		minOut, err = cosmos.ParseUint(params[minOutParam][0])
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("bad min out: %w", err))
		}
	}

	// parse loan asset
	loanAsset, err := common.NewAssetWithShortCodes(mgr.GetVersion(), params[toAssetParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad loan asset: %w", err))
	}
	loanAsset = fuzzyAssetMatch(ctx, mgr.Keeper(), loanAsset)

	// parse loan owner
	loanOwner, err := common.NewAddress(params[loanOwnerParam][0])
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad loan owner: %w", err))
	}

	// generate random from address
	fromAddress, err := types.GetRandomPubKey().GetAddress(asset.Chain)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("bad from address: %w", err))
	}

	// validate if it is valid collateral asset
	key := "LENDING-" + loanAsset.GetDerivedAsset().MimirString()
	val, err := mgr.Keeper().GetMimir(ctx, key)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("ail to fetch LENDING key: %w", err))
	}
	if val <= 0 {
		return quoteErrorResponse(fmt.Errorf("Lending is not available for this collateral asset"))
	}

	loan, err := mgr.Keeper().GetLoan(ctx, loanAsset, loanOwner)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("failed to get loan: %w", err))
	}

	poolRepayment, err := mgr.Keeper().GetPool(ctx, asset)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("failed to get pool: %w", err))
	}

	poolThorAsset, err := mgr.Keeper().GetPool(ctx, common.TOR)
	if err != nil {
		return quoteErrorResponse(fmt.Errorf("failed to get pool: %w", err))
	}

	pendingDebt := loan.DebtIssued.Sub(loan.DebtRepaid)
	totalPendingDebtInRune := poolThorAsset.AssetValueInRune(pendingDebt)
	totalPendingDebtInRepaymentAsset := totalPendingDebtInRune

	if !asset.IsRune() {
		totalPendingDebtInRepaymentAsset = poolRepayment.RuneValueInAsset(totalPendingDebtInRune)
	}

	minBP := mgr.Keeper().GetConfigInt64(ctx, constants.StreamingSwapMinBPFee)
	initialThresholdBasisPoints := sdk.NewUint(uint64(minBP)) // Initial threshold to start looking for the target amount
	amountInTorToRepay := pendingDebt.Mul(repayBps).Quo(sdk.NewUint(10_000))
	amountToRepay := totalPendingDebtInRepaymentAsset.Mul(repayBps).Quo(sdk.NewUint(10_000))
	incrementBasedOnThreshold := amountToRepay.Mul(initialThresholdBasisPoints).Quo(sdk.NewUint(10_000))
	amountPlusThresholdToRepay := amountToRepay.Add(incrementBasedOnThreshold)

	msg := &types.MsgLoanRepayment{
		Owner:           loanOwner,
		CollateralAsset: loanAsset,
		Coin:            common.NewCoin(asset, amountPlusThresholdToRepay),
		From:            fromAddress,
		MinOut:          minOut,
	}

	res, data, err := quoteSimulateCloseLoan(ctx, mgr, msg)
	if data != nil {
		return data, err
	}

	thresholdBasisPoint := initialThresholdBasisPoints

	for thresholdBasisPoint.LTE(sdk.NewUint(1500)) { // Arbitrary cap for the threshold of 1500 BPS to avoid harmful requests.

		exptectedDebtRepaid, err := cosmos.ParseUint(res.ExpectedDebtRepaid)
		if err != nil {
			return quoteErrorResponse(fmt.Errorf("bad exptectedDebtRepaid: %w", err))
		}

		if exptectedDebtRepaid.GTE(amountInTorToRepay) {
			break
		}

		// Arbitrarily increment by 10 BPS per iteration until the target is met. A higher amount results in less server load but also less accurate calculations
		thresholdBasisPoint = thresholdBasisPoint.Add(sdk.NewUint(10))

		// Resimulate with new threshold
		increment := amountToRepay.Mul(thresholdBasisPoint).Quo(sdk.NewUint(10_000))
		newAmount := amountToRepay.Add(increment)
		msg.Coin.Amount = newAmount
		res, data, err = quoteSimulateCloseLoan(ctx, mgr, msg)
		if err != nil {
			return data, err
		}
	}

	// set inbound recommended gas for non-native in asset
	if !asset.Chain.IsTHORChain() {
		inboundGas := mgr.GasMgr().GetGasRate(ctx, asset.Chain)
		res.RecommendedGasRate = wrapString(inboundGas.String())
		res.GasRateUnits = wrapString(asset.Chain.GetGasUnits())
	}

	return json.MarshalIndent(res, "", "  ")
}
