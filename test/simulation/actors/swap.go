package actors

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ecommon "github.com/ethereum/go-ethereum/common"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/test/simulation/pkg/evm"
	"gitlab.com/thorchain/thornode/test/simulation/pkg/thornode"
	. "gitlab.com/thorchain/thornode/test/simulation/pkg/types"
)

////////////////////////////////////////////////////////////////////////////////////////
// SwapActor
////////////////////////////////////////////////////////////////////////////////////////

type SwapActor struct {
	Actor

	account *User

	// starting balances
	from        common.Asset
	fromAddress common.Address
	to          common.Asset
	toAddress   common.Address
	toBalance   sdk.Uint

	swapAmount sdk.Uint
	swapTxID   string

	// expected range for received amount, including outbound fee
	minExpected sdk.Uint
	maxExpected sdk.Uint
}

func NewSwapActor(from, to common.Asset) *Actor {
	a := &SwapActor{
		Actor: *NewActor(fmt.Sprintf("Swap %s => %s", from, to)),
		from:  from,
		to:    to,
	}

	// lock an account with from balance
	a.Ops = append(a.Ops, a.acquireUser)

	// generate swap quote
	a.Ops = append(a.Ops, a.getQuote)

	// send swap inbound
	if from.Chain.IsEVM() && !from.IsGasAsset() {
		a.Ops = append(a.Ops, a.sendTokenSwap)
	} else {
		a.Ops = append(a.Ops, a.sendSwap)
	}

	// verify the swap within expected range
	a.Ops = append(a.Ops, a.verifyOutbound)

	return &a.Actor
}

////////////////////////////////////////////////////////////////////////////////////////
// Ops
////////////////////////////////////////////////////////////////////////////////////////

func (a *SwapActor) acquireUser(config *OpConfig) OpResult {
	// swap 0.5% of from pool depth
	pool, err := thornode.GetPool(a.from)
	if err != nil {
		return OpResult{
			Continue: false,
		}
	}
	a.swapAmount = sdk.NewUintFromString(pool.BalanceAsset).QuoUint64(200)

	for _, user := range config.Users {
		// skip users already being used
		if !user.Acquire() {
			continue
		}

		cl := a.Log().With().
			Str("user", user.Name()).
			Stringer("from", a.from).
			Stringer("to", a.to).
			Logger()
		a.SetLogger(cl)

		// skip users that don't have from asset balance
		fromAcct, err := user.ChainClients[a.from.Chain].GetAccount(nil)
		if err != nil {
			a.Log().Error().Err(err).Msg("failed to get from account")
			user.Release()
			continue
		}
		if fromAcct.Coins.GetCoin(a.from).Amount.LT(a.swapAmount) {
			a.Log().Error().Msg("user has insufficient from balance")
			user.Release()
			continue
		}

		// get l1 address to store in state context
		a.fromAddress, err = user.PubKey().GetAddress(a.from.Chain)
		if err != nil {
			a.Log().Fatal().Err(err).Msg("failed to get L1 address")
		}
		a.toAddress, err = user.PubKey().GetAddress(a.to.Chain)
		if err != nil {
			a.Log().Fatal().Err(err).Msg("failed to get L1 address")
		}

		// get to asset balance for tracking
		toAcct, err := user.ChainClients[a.to.Chain].GetAccount(nil)
		if err != nil {
			a.Log().Warn().Err(err).Msg("failed to get to account")
		} else {
			a.toBalance = toAcct.Coins.GetCoin(a.to).Amount
		}

		// set acquired account and amounts in state context
		a.account = user

		break
	}

	// continue if we acquired a user
	if a.account != nil {
		a.Log().Info().Msg("acquired user")
		return OpResult{
			Continue: true,
		}
	}

	// remain pending if no user is available
	a.Log().Info().Msg("waiting for user with sufficient balance")
	return OpResult{
		Continue: false,
	}
}

func (a *SwapActor) getQuote(config *OpConfig) OpResult {
	quote, err := thornode.GetSwapQuote(a.from, a.to, a.swapAmount)
	if err != nil {
		a.Log().Error().Err(err).Str("amount", a.swapAmount.String()).Msg("failed to get swap quote")
		return OpResult{
			Continue: false,
		}
	}

	// store expected range to fail if received amount is outside 5% tolerance
	quoteOut := sdk.NewUintFromString(quote.ExpectedAmountOut)
	tolerance := quoteOut.QuoUint64(20)
	if quote.Fees.Outbound != nil {
		outboundFee := sdk.NewUintFromString(*quote.Fees.Outbound)
		quoteOut = quoteOut.Add(outboundFee)

		// handle 2x gas rate fluctuation (add 1x outbound fee to tolerance)
		tolerance = tolerance.Add(outboundFee)
	}
	a.minExpected = quoteOut.Sub(tolerance)
	a.maxExpected = quoteOut.Add(tolerance)

	return OpResult{
		Continue: true,
	}
}

func (a *SwapActor) sendSwap(config *OpConfig) OpResult {
	// get inbound address
	inboundAddr, _, err := thornode.GetInboundAddress(a.from.Chain)
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to get inbound address")
		return OpResult{
			Continue: false,
		}
	}

	// if on a utxo chain, shorten the to asset to fuzzy match
	to := a.to.String()
	if a.from.Chain.IsUTXO() && !a.to.IsGasAsset() {
		to = strings.Split(to, "-")[0]
	}

	// create tx out
	memo := fmt.Sprintf("=:%s:%s", to, a.toAddress)
	tx := SimTx{
		Chain:     a.from.Chain,
		ToAddress: inboundAddr,
		Coin:      common.NewCoin(a.from, a.swapAmount),
		Memo:      memo,
	}

	client := a.account.ChainClients[a.from.Chain]

	// sign transaction
	signed, err := client.SignTx(tx)
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to sign tx")
		return OpResult{
			Continue: false,
		}
	}

	// broadcast transaction
	txid, err := client.BroadcastTx(signed)
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to broadcast tx")
		return OpResult{
			Continue: false,
		}
	}
	a.swapTxID = txid

	a.Log().Info().Str("txid", txid).Msg("broadcasted swap tx")
	return OpResult{
		Continue: true,
	}
}

func (a *SwapActor) sendTokenSwap(config *OpConfig) OpResult {
	// get router address
	inboundAddr, routerAddr, err := thornode.GetInboundAddress(a.from.Chain)
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to get inbound address")
		return OpResult{
			Continue: false,
		}
	}
	if routerAddr == nil {
		a.Log().Error().Msg("failed to get router address")
		return OpResult{
			Continue: false,
		}
	}

	token := evm.Tokens(a.from.Chain)[a.from]

	// convert amount to token decimals
	factor := big.NewInt(1).Exp(big.NewInt(10), big.NewInt(int64(token.Decimals)), nil)
	tokenAmount := a.swapAmount.Mul(cosmos.NewUintFromBigInt(factor))
	tokenAmount = tokenAmount.QuoUint64(common.One)

	// approve the router
	eRouterAddr := ecommon.HexToAddress(routerAddr.String())
	tx := SimContractTx{
		Chain:    a.from.Chain,
		Contract: common.Address(token.Address),
		ABI:      evm.ERC20ABI(),
		Method:   "approve",
		Args:     []interface{}{eRouterAddr, tokenAmount.BigInt()},
	}

	iClient := a.account.ChainClients[a.from.Chain]
	client, ok := iClient.(*evm.Client)
	if !ok {
		a.Log().Fatal().Msg("failed to get evm client")
	}

	// sign approve transaction
	signed, err := client.SignContractTx(tx)
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to sign tx")
		return OpResult{
			Continue: false,
		}
	}

	// broadcast approve transaction
	txid, err := client.BroadcastTx(signed)
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to broadcast tx")
		return OpResult{
			Continue: false,
		}
	}
	a.Log().Info().Str("txid", txid).Msg("broadcasted router approve tx")

	// call depositWithExpiry
	memo := fmt.Sprintf("=:%s:%s", a.to, a.toAddress)
	expiry := time.Now().Add(time.Hour).Unix()
	eInboundAddr := ecommon.HexToAddress(inboundAddr.String())
	eTokenAddr := ecommon.HexToAddress(token.Address)
	tx = SimContractTx{
		Chain:    a.from.Chain,
		Contract: *routerAddr,
		ABI:      evm.RouterABI(),
		Method:   "depositWithExpiry",
		Args: []interface{}{
			eInboundAddr,
			eTokenAddr,
			tokenAmount.BigInt(),
			memo,
			big.NewInt(expiry),
		},
	}

	// sign deposit transaction
	signed, err = client.SignContractTx(tx)
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to sign tx")
		return OpResult{
			Continue: false,
		}
	}

	// broadcast deposit transaction
	txid, err = client.BroadcastTx(signed)
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to broadcast tx")
		return OpResult{
			Continue: false,
		}
	}

	a.swapTxID = txid

	a.Log().Info().Str("txid", txid).Msg("broadcasted swap tx")
	return OpResult{
		Continue: true,
	}
}

func (a *SwapActor) verifyOutbound(config *OpConfig) OpResult {
	// get swap stages
	stages, err := thornode.GetTxStages(a.swapTxID)
	if err != nil {
		a.Log().Warn().Err(err).Msg("failed to get tx stages")
		return OpResult{
			Continue: false,
		}
	}

	// wait for outbound to be marked complete
	if stages.OutboundSigned == nil || !stages.OutboundSigned.Completed {
		return OpResult{
			Continue: false,
		}
	}

	// get tx details
	details, err := thornode.GetTxDetails(a.swapTxID)
	if err != nil {
		a.Log().Warn().Err(err).Msg("failed to get tx details")
		return OpResult{
			Continue: false,
		}
	}

	// verify exactly one out transaction
	if len(details.OutTxs) != 1 {
		return OpResult{
			Error:  fmt.Errorf("expected exactly one out transaction"),
			Finish: true,
		}
	}

	// verify exactly one action
	if len(details.Actions) != 1 {
		return OpResult{
			Error:  fmt.Errorf("expected exactly one action"),
			Finish: true,
		}
	}

	// verify outbound amount + max gas within expected range
	action := details.Actions[0]
	out := details.OutTxs[0]
	outAmountPlusMaxGas := cosmos.NewUintFromString(out.Coins[0].Amount)
	maxGas := action.MaxGas[0]
	if maxGas.Asset == a.to.String() {
		outAmountPlusMaxGas = outAmountPlusMaxGas.Add(cosmos.NewUintFromString(maxGas.Amount))
	} else {
		var maxGasAssetValue sdk.Uint
		maxGasAssetValue, err = thornode.ConvertAssetAmount(maxGas, a.to.String())
		if err != nil {
			a.Log().Warn().Err(err).Msg("failed to convert asset")
			return OpResult{
				Continue: false,
			}
		}
		outAmountPlusMaxGas = outAmountPlusMaxGas.Add(maxGasAssetValue)
	}

	// retrieve L1 balance
	toAcct, err := a.account.ChainClients[a.to.Chain].GetAccount(nil)
	if err != nil {
		a.Log().Warn().Err(err).Msg("failed to get to account")
		return OpResult{
			Continue: false,
		}
	}

	// check received amount
	received := toAcct.Coins.GetCoin(a.to).Amount.Sub(a.toBalance)
	a.Log().Info().
		Stringer("received", received).
		Stringer("outAmountPlusMaxGas", outAmountPlusMaxGas).
		Stringer("minExpected", a.minExpected).
		Stringer("maxExpected", a.maxExpected).
		Msg("swap complete")

	// fail if received amount is outside expected range
	if outAmountPlusMaxGas.LT(a.minExpected) || outAmountPlusMaxGas.GT(a.maxExpected) {
		err = fmt.Errorf("out amount plus gas outside tolerance")
	}

	// release user
	a.account.Release()

	return OpResult{
		Finish: true,
		Error:  err,
	}
}
