package actors

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
	openapi "gitlab.com/thorchain/thornode/openapi/gen"
	"gitlab.com/thorchain/thornode/test/simulation/pkg/thornode"
	. "gitlab.com/thorchain/thornode/test/simulation/pkg/types"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

////////////////////////////////////////////////////////////////////////////////////////
// ArbActor
////////////////////////////////////////////////////////////////////////////////////////

type ArbActor struct {
	Actor

	account     *User
	thorAddress cosmos.AccAddress

	// originalPools maps the asset to the first seen available pool state (arb target)
	originalPools map[string]types.Pool
}

func NewArbActor() *Actor {
	a := &ArbActor{
		Actor:         *NewActor("Arbitrage"),
		originalPools: make(map[string]types.Pool),
	}
	a.Timeout = time.Hour
	a.Interval = 5 * time.Second // roughly once per block

	// lock an account to use for arb
	a.Ops = append(a.Ops, a.acquireUser)

	// enable trade assets
	a.Ops = append(a.Ops, a.enableTradeAssets)

	// convert all assets to trade assets
	a.Ops = append(a.Ops, a.bootstrapTradeAssets)

	// arb until pools are drained
	a.Ops = append(a.Ops, a.arb)

	return &a.Actor
}

////////////////////////////////////////////////////////////////////////////////////////
// Ops
////////////////////////////////////////////////////////////////////////////////////////

func (a *ArbActor) acquireUser(config *OpConfig) OpResult {
	for _, user := range config.Users {
		// skip users already being used
		if !user.Acquire() {
			continue
		}

		cl := a.Log().With().Str("user", user.Name()).Logger()
		a.SetLogger(cl)

		// set acquired account and amounts in state context
		a.account = user

		// set thorchain address for later use
		thorAddress, err := user.PubKey().GetThorAddress()
		if err != nil {
			a.Log().Error().Err(err).Msg("failed to get thor address")
			user.Release()
			continue
		}
		a.thorAddress = thorAddress

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

func (a *ArbActor) enableTradeAssets(config *OpConfig) OpResult {
	// wait to acquire the admin user
	if !config.AdminUser.Acquire() {
		return OpResult{
			Continue: false,
		}
	}

	// enable trade assets
	accAddr, err := config.AdminUser.PubKey().GetThorAddress()
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to get thor address")
		return OpResult{
			Continue: false,
		}
	}
	mimirMsg := types.NewMsgMimir("TradeAccountsEnabled", 1, accAddr)
	txid, err := config.AdminUser.Thorchain.Broadcast(mimirMsg)
	if err != nil {
		a.Log().Fatal().Err(err).Msg("failed to broadcast tx")
	}

	a.Log().Info().
		Stringer("txid", txid).
		Msg("broadcasted admin mimir tx to enable trade assets")

	return OpResult{
		Continue: true,
	}
}

func (a *ArbActor) bootstrapTradeAssets(config *OpConfig) OpResult {
	// get all pools
	pools, err := thornode.GetPools()
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to get pools")
		return OpResult{
			Continue: false,
			Error:    err,
		}
	}

	// deposit trade assets for all pools
	for _, pool := range pools {
		asset, err := common.NewAsset(pool.Asset)
		if err != nil {
			a.Log().Fatal().Err(err).Str("asset", pool.Asset).Msg("failed to create asset")
		}

		// get deposit parameters for 90% of asset balance
		client := a.account.ChainClients[asset.Chain]
		memo := fmt.Sprintf("trade+:%s", a.thorAddress)
		l1Acct, err := a.account.ChainClients[asset.Chain].GetAccount(nil)
		if err != nil {
			a.Log().Fatal().Err(err).Msg("failed to get L1 account")
		}
		depositAmount := l1Acct.Coins.GetCoin(asset).Amount.QuoUint64(10).MulUint64(9)

		// make deposit
		var txid string
		if asset.Chain.IsEVM() && !asset.IsGasAsset() {
			txid, err = depositL1Token(a.Log(), client, asset, memo, depositAmount)
		} else {
			txid, err = depositL1(a.Log(), client, asset, memo, depositAmount)
		}
		if err != nil {
			a.Log().Fatal().
				Err(err).
				Str("asset", asset.String()).
				Msg("failed to deposit trade asset")
		}
		a.Log().Info().
			Stringer("asset", asset).
			Str("txid", txid).
			Msg("deposited trade asset")
	}

	// mark actor as backgrounded
	a.Log().Info().Msg("moving arbitrage actor to background")
	a.Background()

	return OpResult{
		Continue: true,
	}
}

func (a *ArbActor) arb(config *OpConfig) OpResult {
	// get all pools
	pools, err := thornode.GetPools()
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to get pools")
		return OpResult{
			Continue: false,
			Error:    err,
		}
	}

	// if pools are drained then we are done
	if len(pools) == 0 {
		a.account.Release()
		a.Log().Info().Msg("pools are drained, nothing more to arb")
		return OpResult{
			Finish: true,
			Error:  nil,
		}
	}

	// gather pools we have seen
	arbPools := []openapi.Pool{}
	for _, pool := range pools {
		// skip unavailable pools and those with no liquidity
		if pool.BalanceRune == "0" || pool.BalanceAsset == "0" || pool.Status != types.PoolStatus_Available.String() {
			continue
		}

		// if this is the first time we see the pool, store it to use as the target price
		if _, ok := a.originalPools[pool.Asset]; !ok {
			a.originalPools[pool.Asset] = types.Pool{
				BalanceRune:  cosmos.NewUintFromString(pool.BalanceRune),
				BalanceAsset: cosmos.NewUintFromString(pool.BalanceAsset),
			}
			continue
		}

		arbPools = append(arbPools, pool)
	}

	// skip if there are not enough pools to arb
	if len(arbPools) < 2 {
		a.Log().Info().Msg("not enough pools to arb")
		return OpResult{
			Continue: false,
		}
	}

	// sort pools by price change
	priceChangeBps := func(pool openapi.Pool) int64 {
		originalPool := a.originalPools[pool.Asset]
		originalPrice := originalPool.BalanceRune.MulUint64(1e8).Quo(originalPool.BalanceAsset)
		currentPrice := cosmos.NewUintFromString(pool.BalanceRune).MulUint64(1e8).Quo(cosmos.NewUintFromString(pool.BalanceAsset))
		return int64(constants.MaxBasisPts) - int64(originalPrice.MulUint64(constants.MaxBasisPts).Quo(currentPrice).Uint64())
	}
	sort.Slice(arbPools, func(i, j int) bool {
		return priceChangeBps(arbPools[i]) > priceChangeBps(arbPools[j])
	})

	send := arbPools[0]
	receive := arbPools[len(arbPools)-1]

	// skip if none have diverged more than 10 basis points
	adjustmentBps := common.Min(common.Abs(priceChangeBps(send)), common.Abs(priceChangeBps(receive)))
	if adjustmentBps < 10 {
		a.Log().Info().
			Int64("maxShift", priceChangeBps(send)).
			Int64("minShift", priceChangeBps(receive)).
			Msg("pools have not diverged enough to arb")
		return OpResult{
			Continue: false,
		}
	}

	// build the swap
	memo := fmt.Sprintf("=:%s", strings.Replace(receive.Asset, ".", "~", 1))
	asset, err := common.NewAsset(strings.Replace(send.Asset, ".", "~", 1))
	if err != nil {
		a.Log().Fatal().Err(err).Str("asset", send.Asset).Msg("failed to create asset")
	}
	amount := cosmos.NewUint(uint64(adjustmentBps / 2)).Mul(cosmos.NewUintFromString(send.BalanceAsset)).QuoUint64(constants.MaxBasisPts)
	coin := common.NewCoin(asset, amount)

	// build the swap
	deposit := types.NewMsgDeposit(common.NewCoins(coin), memo, a.thorAddress)
	a.Log().Info().Interface("deposit", deposit).Msg("arbing most diverged pool")

	// broadcast the swap
	txid, err := a.account.Thorchain.Broadcast(deposit)
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to broadcast tx")
		return OpResult{
			Continue: false,
		}
	}

	a.Log().Info().Stringer("txid", txid).Str("memo", memo).Msg("broadcasted arb tx")

	return OpResult{
		Continue: false,
	}
}
