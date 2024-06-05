package thorchain

import (
	"fmt"
	"strings"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
)

type SwapQueueV121Suite struct{}

var _ = Suite(&SwapQueueV121Suite{})

func (s SwapQueueV121Suite) TestGetTodoNum(c *C) {
	queue := newSwapQueueV121(keeper.KVStoreDummy{})

	c.Check(queue.getTodoNum(50, 10, 100), Equals, int64(25))     // halves it
	c.Check(queue.getTodoNum(11, 10, 100), Equals, int64(5))      // halves it
	c.Check(queue.getTodoNum(10, 10, 100), Equals, int64(10))     // does all of them
	c.Check(queue.getTodoNum(1, 10, 100), Equals, int64(1))       // does all of them
	c.Check(queue.getTodoNum(0, 10, 100), Equals, int64(0))       // does none
	c.Check(queue.getTodoNum(10000, 10, 100), Equals, int64(100)) // does max 100
	c.Check(queue.getTodoNum(200, 10, 100), Equals, int64(100))   // does max 100
}

func (s SwapQueueV121Suite) TestScoreMsgs(c *C) {
	ctx, k := setupKeeperForTest(c)

	pool := NewPool()
	pool.Asset = common.BNBAsset
	pool.BalanceRune = cosmos.NewUint(143166 * common.One)
	pool.BalanceAsset = cosmos.NewUint(1000 * common.One)
	c.Assert(k.SetPool(ctx, pool), IsNil)
	pool = NewPool()
	pool.Asset = common.BTCAsset
	pool.BalanceRune = cosmos.NewUint(73708333 * common.One)
	pool.BalanceAsset = cosmos.NewUint(1000 * common.One)
	c.Assert(k.SetPool(ctx, pool), IsNil)
	pool = NewPool()
	pool.Asset = common.ETHAsset
	pool.BalanceRune = cosmos.NewUint(1000 * common.One)
	pool.BalanceAsset = cosmos.NewUint(1000 * common.One)
	pool.Status = PoolStaged
	c.Assert(k.SetPool(ctx, pool), IsNil)

	queue := newSwapQueueV121(k)

	// check that we sort by liquidity ok
	msgs := []*MsgSwap{
		NewMsgSwap(common.Tx{
			ID:    common.TxID("5E1DF027321F1FE37CA19B9ECB11C2B4ABEC0D8322199D335D9CE4C39F85F115"),
			Coins: common.Coins{common.NewCoin(common.RuneAsset(), cosmos.NewUint(2*common.One))},
		}, common.BNBAsset, GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
		NewMsgSwap(common.Tx{
			ID:    common.TxID("53C1A22436B385133BDD9157BB365DB7AAC885910D2FA7C9DC3578A04FFD4ADC"),
			Coins: common.Coins{common.NewCoin(common.BNBAsset, cosmos.NewUint(50*common.One))},
		}, common.RuneAsset(), GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
		NewMsgSwap(common.Tx{
			ID:    common.TxID("6A470EB9AFE82981979A5EEEED3296E1E325597794BD5BFB3543A372CAF435E5"),
			Coins: common.Coins{common.NewCoin(common.RuneAsset(), cosmos.NewUint(1*common.One))},
		}, common.BNBAsset, GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
		NewMsgSwap(common.Tx{
			ID:    common.TxID("5EE9A7CCC55A3EBAFA0E542388CA1B909B1E3CE96929ED34427B96B7CCE9F8E8"),
			Coins: common.Coins{common.NewCoin(common.RuneAsset(), cosmos.NewUint(100*common.One))},
		}, common.BNBAsset, GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
		NewMsgSwap(common.Tx{
			ID:    common.TxID("0FF2A521FB11FFEA4DFE3B7AD4066FF0A33202E652D846F8397EFC447C97A91B"),
			Coins: common.Coins{common.NewCoin(common.RuneAsset(), cosmos.NewUint(10*common.One))},
		}, common.BNBAsset, GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),

		NewMsgSwap(common.Tx{
			ID:    GetRandomTxHash(),
			Coins: common.Coins{common.NewCoin(common.BNBAsset, cosmos.NewUint(150*common.One))},
		}, common.RuneAsset(), GetRandomTHORAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),

		NewMsgSwap(common.Tx{
			ID:    GetRandomTxHash(),
			Coins: common.Coins{common.NewCoin(common.BNBAsset, cosmos.NewUint(151*common.One))},
		}, common.RuneAsset(), GetRandomTHORAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),

		// synthetics can be redeemed on unavailable pools, should score
		NewMsgSwap(common.Tx{
			ID:    GetRandomTxHash(),
			Coins: common.Coins{common.NewCoin(common.ETHAsset.GetSyntheticAsset(), cosmos.NewUint(3*common.One))},
		}, common.RuneAsset(), GetRandomTHORAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
	}

	swaps := make(swapItems, len(msgs))
	for i, msg := range msgs {
		swaps[i] = swapItem{
			msg:  *msg,
			fee:  cosmos.ZeroUint(),
			slip: cosmos.ZeroUint(),
		}
	}
	swaps, err := queue.scoreMsgs(ctx, swaps, 10_000)
	c.Assert(err, IsNil)
	swaps = swaps.Sort()
	c.Check(swaps, HasLen, 8)
	c.Check(swaps[0].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(151*common.One)), Equals, true, Commentf("%d", swaps[0].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[1].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(150*common.One)), Equals, true, Commentf("%d", swaps[1].msg.Tx.Coins[0].Amount.Uint64()))
	// 50 BNB is worth more than 100 RUNE
	c.Check(swaps[2].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(50*common.One)), Equals, true, Commentf("%d", swaps[2].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[3].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(3*common.One)), Equals, true, Commentf("%d", swaps[3].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[4].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(100*common.One)), Equals, true, Commentf("%d", swaps[4].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[5].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(10*common.One)), Equals, true, Commentf("%d", swaps[5].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[6].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(2*common.One)), Equals, true, Commentf("%d", swaps[6].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[7].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(1*common.One)), Equals, true, Commentf("%d", swaps[7].msg.Tx.Coins[0].Amount.Uint64()))

	// check that slip is taken into account
	msgs = []*MsgSwap{
		NewMsgSwap(common.Tx{
			ID:    GetRandomTxHash(),
			Coins: common.Coins{common.NewCoin(common.BNBAsset, cosmos.NewUint(2*common.One))},
		}, common.RuneAsset(), GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
		NewMsgSwap(common.Tx{
			ID:    GetRandomTxHash(),
			Coins: common.Coins{common.NewCoin(common.BNBAsset, cosmos.NewUint(50*common.One))},
		}, common.RuneAsset(), GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
		NewMsgSwap(common.Tx{
			ID:    GetRandomTxHash(),
			Coins: common.Coins{common.NewCoin(common.BNBAsset, cosmos.NewUint(1*common.One))},
		}, common.RuneAsset(), GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
		NewMsgSwap(common.Tx{
			ID:    GetRandomTxHash(),
			Coins: common.Coins{common.NewCoin(common.BNBAsset, cosmos.NewUint(100*common.One))},
		}, common.RuneAsset(), GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
		NewMsgSwap(common.Tx{
			ID:    GetRandomTxHash(),
			Coins: common.Coins{common.NewCoin(common.BNBAsset, cosmos.NewUint(10*common.One))},
		}, common.RuneAsset(), GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
		NewMsgSwap(common.Tx{
			ID:    GetRandomTxHash(),
			Coins: common.Coins{common.NewCoin(common.BTCAsset, cosmos.NewUint(2*common.One))},
		}, common.RuneAsset(), GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
		NewMsgSwap(common.Tx{
			ID:    GetRandomTxHash(),
			Coins: common.Coins{common.NewCoin(common.BTCAsset, cosmos.NewUint(50*common.One))},
		}, common.RuneAsset(), GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
		NewMsgSwap(common.Tx{
			ID:    GetRandomTxHash(),
			Coins: common.Coins{common.NewCoin(common.BTCAsset, cosmos.NewUint(1*common.One))},
		}, common.RuneAsset(), GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
		NewMsgSwap(common.Tx{
			ID:    GetRandomTxHash(),
			Coins: common.Coins{common.NewCoin(common.BTCAsset, cosmos.NewUint(100*common.One))},
		}, common.RuneAsset(), GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
		NewMsgSwap(common.Tx{
			ID:    GetRandomTxHash(),
			Coins: common.Coins{common.NewCoin(common.BTCAsset, cosmos.NewUint(10*common.One))},
		}, common.RuneAsset(), GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),

		NewMsgSwap(common.Tx{
			ID:    GetRandomTxHash(),
			Coins: common.Coins{common.NewCoin(common.BTCAsset, cosmos.NewUint(10*common.One))},
		}, common.BNBAsset, GetRandomBNBAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(),
			"", "", nil,
			MarketOrder,
			0, 0, GetRandomBech32Addr()),
	}

	swaps = make(swapItems, len(msgs))
	for i, msg := range msgs {
		swaps[i] = swapItem{
			msg:  *msg,
			fee:  cosmos.ZeroUint(),
			slip: cosmos.ZeroUint(),
		}
	}
	swaps, err = queue.scoreMsgs(ctx, swaps, 10_000)
	c.Assert(err, IsNil)
	swaps = swaps.Sort()
	c.Assert(swaps, HasLen, 11)

	c.Check(swaps[0].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(10*common.One)), Equals, true, Commentf("%d", swaps[0].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[0].msg.Tx.Coins[0].Asset.Equals(common.BTCAsset), Equals, true)

	c.Check(swaps[1].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(100*common.One)), Equals, true, Commentf("%d", swaps[1].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[1].msg.Tx.Coins[0].Asset.Equals(common.BTCAsset), Equals, true)

	c.Check(swaps[2].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(100*common.One)), Equals, true, Commentf("%d", swaps[2].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[2].msg.Tx.Coins[0].Asset.Equals(common.BNBAsset), Equals, true)

	c.Check(swaps[3].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(50*common.One)), Equals, true, Commentf("%d", swaps[3].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[3].msg.Tx.Coins[0].Asset.Equals(common.BTCAsset), Equals, true)

	c.Check(swaps[4].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(50*common.One)), Equals, true, Commentf("%d", swaps[4].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[4].msg.Tx.Coins[0].Asset.Equals(common.BNBAsset), Equals, true)

	c.Check(swaps[5].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(10*common.One)), Equals, true, Commentf("%d", swaps[5].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[5].msg.Tx.Coins[0].Asset.Equals(common.BTCAsset), Equals, true)

	c.Check(swaps[6].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(10*common.One)), Equals, true, Commentf("%d", swaps[6].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[6].msg.Tx.Coins[0].Asset.Equals(common.BNBAsset), Equals, true)

	c.Check(swaps[7].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(2*common.One)), Equals, true, Commentf("%d", swaps[7].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[7].msg.Tx.Coins[0].Asset.Equals(common.BTCAsset), Equals, true)

	c.Check(swaps[8].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(2*common.One)), Equals, true, Commentf("%d", swaps[8].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[8].msg.Tx.Coins[0].Asset.Equals(common.BNBAsset), Equals, true)

	c.Check(swaps[9].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(1*common.One)), Equals, true, Commentf("%d", swaps[9].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[9].msg.Tx.Coins[0].Asset.Equals(common.BTCAsset), Equals, true)

	c.Check(swaps[10].msg.Tx.Coins[0].Amount.Equal(cosmos.NewUint(1*common.One)), Equals, true, Commentf("%d", swaps[10].msg.Tx.Coins[0].Amount.Uint64()))
	c.Check(swaps[10].msg.Tx.Coins[0].Asset.Equals(common.BNBAsset), Equals, true)
}

func (s SwapQueueV121Suite) TestStreamingSwapSelection(c *C) {
	ctx, k := setupKeeperForTest(c)
	queue := newSwapQueueV121(k)

	bnbAddr := GetRandomBNBAddress()
	txID := GetRandomTxHash()
	tx := common.NewTx(
		txID,
		bnbAddr,
		bnbAddr,
		common.NewCoins(common.NewCoin(common.RuneAsset(), cosmos.NewUint(common.One*100))),
		BNBGasFeeSingleton,
		"",
	)

	// happy path
	msg := NewMsgSwap(tx, common.BNBAsset.GetSyntheticAsset(), GetRandomTHORAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(), "", "", nil, MarketOrder, 10, 20, GetRandomBech32Addr())
	c.Assert(k.SetSwapQueueItem(ctx, *msg, 0), IsNil)

	// no saved streaming swap, should swap now
	items, err := queue.FetchQueue(ctx)
	c.Assert(err, IsNil)
	c.Check(items, HasLen, 1)

	// save streaming swap data, should have same result
	swp := msg.GetStreamingSwap()
	k.SetStreamingSwap(ctx, swp)
	items, err = queue.FetchQueue(ctx)
	c.Assert(err, IsNil)
	c.Check(items, HasLen, 1)

	// last height is this block, no result
	swp.LastHeight = ctx.BlockHeight()
	k.SetStreamingSwap(ctx, swp)
	items, err = queue.FetchQueue(ctx)
	c.Assert(err, IsNil)
	c.Check(items, HasLen, 0)

	// last height is halfway there
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + (int64(swp.Interval) / 2))
	items, err = queue.FetchQueue(ctx)
	c.Assert(err, IsNil)
	c.Check(items, HasLen, 0)

	// last height is interval blocks ago
	ctx = ctx.WithBlockHeight(swp.LastHeight + int64(swp.Interval))
	items, err = queue.FetchQueue(ctx)
	c.Assert(err, IsNil)
	c.Check(items, HasLen, 1)
}

func (s SwapQueueV121Suite) TestStreamingSwapOutbounds(c *C) {
	ctx, mgr := setupManagerForTest(c)
	mgr.txOutStore = NewTxStoreDummy()

	pool := NewPool()
	pool.Asset = common.BNBAsset
	pool.BalanceRune = cosmos.NewUint(143166 * common.One)
	pool.BalanceAsset = cosmos.NewUint(1000 * common.One)
	c.Assert(mgr.Keeper().SetPool(ctx, pool), IsNil)
	pool.Asset = common.BNBAsset
	c.Assert(mgr.Keeper().SetPool(ctx, pool), IsNil)

	queue := newSwapQueueV121(mgr.Keeper())

	badHandler := func(mgr Manager) cosmos.Handler {
		return func(ctx cosmos.Context, msg cosmos.Msg) (*cosmos.Result, error) {
			return nil, fmt.Errorf("failed handler")
		}
	}
	/*
		goodHandler := func(mgr Manager) cosmos.Handler {
			return func(ctx cosmos.Context, msg cosmos.Msg) (*cosmos.Result, error) {
				return nil, fmt.Errorf("failed handler")
			}
		}
	*/

	bnbAddr := GetRandomBNBAddress()
	btcAddr := GetRandomBTCAddress()
	txID := GetRandomTxHash()
	tx := common.NewTx(
		txID,
		bnbAddr,
		bnbAddr,
		common.NewCoins(common.NewCoin(common.BNBAsset, cosmos.NewUint(common.One*100))),
		BNBGasFeeSingleton,
		fmt.Sprintf("=:BTC.BTC:%s", btcAddr),
	)

	msg := NewMsgSwap(tx, common.BTCAsset, btcAddr, cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(), "", "", nil, MarketOrder, 10, 20, GetRandomBech32Addr())
	swp := msg.GetStreamingSwap()
	mgr.Keeper().SetStreamingSwap(ctx, swp)
	c.Assert(mgr.Keeper().SetSwapQueueItem(ctx, *msg, 0), IsNil)

	// test that the refund handler works
	queue.handler = badHandler
	c.Assert(queue.EndBlock(ctx, mgr), IsNil)
	items, err := mgr.TxOutStore().GetOutboundItems(ctx)
	c.Assert(err, IsNil)
	c.Assert(items, HasLen, 1)
	c.Check(strings.HasPrefix(items[0].Memo, "REFUND:"), Equals, true)
	// ensure swp has been deleted
	c.Check(mgr.Keeper().StreamingSwapExists(ctx, txID), Equals, false)
	// ensure swap queue item is gone
	_, err = mgr.Keeper().GetSwapQueueItem(ctx, txID, 0)
	c.Assert(err, NotNil)
	c.Check(err.Error(), Equals, "not found")
	mgr.TxOutStore().ClearOutboundItems(ctx)

	// test we DO NOT send outbound while streaming swap isn't done
	swp.In = swp.Deposit.QuoUint64(2)
	swp.Out = cosmos.NewUint(12345)
	swp.Count = 5
	mgr.Keeper().SetStreamingSwap(ctx, swp)
	c.Assert(mgr.Keeper().SetSwapQueueItem(ctx, *msg, 0), IsNil)
	c.Assert(queue.EndBlock(ctx, mgr), IsNil)
	items, err = mgr.TxOutStore().GetOutboundItems(ctx)
	c.Assert(err, IsNil)
	c.Assert(items, HasLen, 0)
	// make sure we haven't delete the streaming swap entity
	c.Check(mgr.Keeper().StreamingSwapExists(ctx, txID), Equals, true)
	// ensure swap queue item is NOT gone
	_, err = mgr.Keeper().GetSwapQueueItem(ctx, txID, 0)
	c.Assert(err, IsNil)
	mgr.TxOutStore().ClearOutboundItems(ctx)

	// test we DO send outbounds while streaming swap is done
	swp.In = swp.Deposit.QuoUint64(3)
	swp.Out = cosmos.NewUint(12345)
	swp.Count = swp.Quantity
	mgr.Keeper().SetStreamingSwap(ctx, swp)
	c.Assert(queue.EndBlock(ctx, mgr), IsNil)
	items, err = mgr.TxOutStore().GetOutboundItems(ctx)
	c.Assert(err, IsNil)
	c.Assert(items, HasLen, 2)
	c.Check(items[0].Memo, Equals, "") // ensure its not a refund tx
	c.Check(items[1].Memo, Equals, "") // ensure its not a refund tx
	c.Check(items[0].Coin.Equals(common.NewCoin(common.BTCAsset, cosmos.NewUint(12345))), Equals, true, Commentf("%s", items[0].Coin.String()))
	c.Check(items[1].Coin.Equals(common.NewCoin(common.BNBAsset, cosmos.NewUint(6666666667))), Equals, true, Commentf("%s", items[1].Coin.String()))
	// make sure we have deleted the streaming swap entity
	c.Check(mgr.Keeper().StreamingSwapExists(ctx, txID), Equals, false)
	// ensure swap queue item is gone
	_, err = mgr.Keeper().GetSwapQueueItem(ctx, txID, 0)
	c.Assert(err, NotNil)
	mgr.TxOutStore().ClearOutboundItems(ctx)

	// test we do send send the outbound (no refund needed)
	swp.In = swp.Deposit
	swp.Out = cosmos.NewUint(12345)
	mgr.Keeper().SetStreamingSwap(ctx, swp)
	c.Assert(mgr.Keeper().SetSwapQueueItem(ctx, *msg, 0), IsNil)
	c.Assert(queue.EndBlock(ctx, mgr), IsNil)
	items, err = mgr.TxOutStore().GetOutboundItems(ctx)
	c.Assert(err, IsNil)
	c.Assert(items, HasLen, 1)
	c.Check(items[0].Memo, Equals, "") // ensure its not a refund tx
	c.Check(items[0].Coin.Equals(common.NewCoin(common.BTCAsset, cosmos.NewUint(12345))), Equals, true, Commentf("%s", items[0].Coin.String()))
	// make sure we have deleted the streaming swap entity
	c.Check(mgr.Keeper().StreamingSwapExists(ctx, txID), Equals, false)
	// ensure swap queue item is gone
	_, err = mgr.Keeper().GetSwapQueueItem(ctx, txID, 0)
	c.Assert(err, NotNil)
	mgr.TxOutStore().ClearOutboundItems(ctx)
}
