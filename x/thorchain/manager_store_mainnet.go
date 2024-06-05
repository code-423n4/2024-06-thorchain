//go:build !stagenet && !mocknet && !regtest
// +build !stagenet,!mocknet,!regtest

package thorchain

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cometbft/cometbft/crypto/tmhash"
	ctypes "github.com/cosmos/cosmos-sdk/types"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
)

func migrateStoreV86(ctx cosmos.Context, mgr *Mgrs) {}

func importPreRegistrationTHORNames(ctx cosmos.Context, mgr Manager) error {
	oneYear := mgr.Keeper().GetConfigInt64(ctx, constants.BlocksPerYear)
	names, err := getPreRegisterTHORNames(ctx, ctx.BlockHeight()+oneYear)
	if err != nil {
		return err
	}

	for _, name := range names {
		mgr.Keeper().SetTHORName(ctx, name)
	}
	return nil
}

func migrateStoreV88(ctx cosmos.Context, mgr Manager) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v88", "error", err)
		}
	}()

	err := importPreRegistrationTHORNames(ctx, mgr)
	if err != nil {
		ctx.Logger().Error("fail to migrate store to v88", "error", err)
	}
}

// no op
func migrateStoreV102(ctx cosmos.Context, mgr Manager) {}

func migrateStoreV103(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v102", "error", err)
		}
	}()

	// MAINNET REFUND
	// A user sent two 4,500 RUNE swap out txs (to USDT), but the external asset matching had a conflict and the outbounds were dropped. Txs:

	// https://viewblock.io/thorchain/tx/B07A6B1B40ADBA2E404D9BCE1BEF6EDE6F70AD135E83806E4F4B6863CF637D0B
	// https://viewblock.io/thorchain/tx/4795A3C036322493A9692B5D44E7D4FF29C3E2C1E848637184E98FE8B05FD06E

	// The below methodology was tested on Stagenet, results are documented here: https://gitlab.com/thorchain/thornode/-/merge_requests/2596#note_1216814315

	// The RUNE was swapped to ETH, but the outbound swap out was dropped by Bifrost. This means RUNE was added, ETH was removed from
	// the pool. This must be reversed and the RUNE sent back to the user.
	// So:
	// 1. Credit the total ETH amount back the pool, this ETH is already in the pool since the outbound was dropped.
	// 2. Deduct the RUNE balance from the ETH pool, this will be sent back to the user.
	// 3. Send the user RUNE from Asgard.
	//
	// Note: the Asgard vault does not need to be credited the ETH since the outbound was never sent, thus never observed (which
	// is where Vault funds are subtracted)

	firstSwapOut := DroppedSwapOutTx{
		inboundHash: "B07A6B1B40ADBA2E404D9BCE1BEF6EDE6F70AD135E83806E4F4B6863CF637D0B",
		gasAsset:    common.ETHAsset,
	}

	err := refundDroppedSwapOutFromRUNE(ctx, mgr, firstSwapOut)
	if err != nil {
		ctx.Logger().Error("fail to migrate store to v103 refund failed", "error", err, "tx hash", "B07A6B1B40ADBA2E404D9BCE1BEF6EDE6F70AD135E83806E4F4B6863CF637D0B")
	}

	secondSwapOut := DroppedSwapOutTx{
		inboundHash: "4795A3C036322493A9692B5D44E7D4FF29C3E2C1E848637184E98FE8B05FD06E",
		gasAsset:    common.ETHAsset,
	}

	err = refundDroppedSwapOutFromRUNE(ctx, mgr, secondSwapOut)
	if err != nil {
		ctx.Logger().Error("fail to migrate store to v103 refund failed", "error", err, "tx hash", "4795A3C036322493A9692B5D44E7D4FF29C3E2C1E848637184E98FE8B05FD06E")
	}
}

func migrateStoreV106(ctx cosmos.Context, mgr *Mgrs) {
	// refund tx stuck in pending state: https://thorchain.net/tx/BC12B3B715546053A2D5615ADB4B3C2C648613D44AA9E942F2BDE40AB09EAA86
	// pool module still contains 4884 synth eth/thor: https://thornode.ninerealms.com/cosmos/bank/v1beta1/balances/thor1g98cy3n9mmjrpn0sxmn63lztelera37n8n67c0?height=9221024
	// deduct 4884 from pool module, create 4884 to user address: thor1vlzlsjfx2l3anh6wsh293fv2e8yh9rwpg7u723
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v106", "error", err)
		}
	}()

	recipient, err := cosmos.AccAddressFromBech32("thor1vlzlsjfx2l3anh6wsh293fv2e8yh9rwpg7u723")
	if err != nil {
		ctx.Logger().Error("fail to create acc address from bech32", err)
		return
	}

	coins := cosmos.NewCoins(cosmos.NewCoin(
		"eth/thor-0xa5f2211b9b8170f694421f2046281775e8468044",
		cosmos.NewInt(488432852150),
	))
	if err := mgr.coinKeeper.SendCoinsFromModuleToAccount(ctx, AsgardName, recipient, coins); err != nil {
		ctx.Logger().Error("fail to SendCoinsFromModuleToAccount", err)
	}
}

func migrateStoreV108(ctx cosmos.Context, mgr *Mgrs) {
	// Requeue four BCH.BCH txout (dangling actions) items swallowed in a chain halt.
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v108", "error", err)
		}
	}()

	danglingInboundTxIDs := []common.TxID{
		"5840920B63CDB9A02028ABB844B28F0305C2B37ADA4F936B69EBEFA04E2F826B",
		"BFACE691A12E85083DD2E4E4ADFBE702299DA6F08A98E6B6F7CF95A9D1D71632",
		"395EBDADA6D0975CF4D3F2E2BD7E246037C672C9CAB97DBFB744CC0F2FFABE95",
		"5881692D0522D0D5221A61FC0704B018ED51A6C43475063ADF6AC912D748208D",
	}

	requeueDanglingActionsV108(ctx, mgr, danglingInboundTxIDs)
}

func migrateStoreV109(ctx cosmos.Context, mgr *Mgrs) {
	// Requeue ETH-chain dangling actions swallowed in a chain halt.
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v109", "error", err)
		}
	}()

	danglingInboundTxIDs := []common.TxID{
		"91C72EFCCF18AE043D036E2A207CC03A063E60024899E050AA7070EF15956BD7",
		"8D17D78A9E3168B88EFDBC30C5ADB3B09459C981B784D8F63C931988295DFE3B",
		"AD88EC612C188E62352F6157B26B97D76BD981744CE4C5AAC672F6338737F011",
		"88FD1BE075C55F18E73DD176E82A870F93B0E4692D514C36C8BF23692B139DED",
		"037254E2534D979FA196EC7B42C62A121B7A46D6854F9EC6FBE33C24B237EF0C",
	}

	requeueDanglingActionsV108(ctx, mgr, danglingInboundTxIDs)
	createFakeTxInsAndMakeObservations(ctx, mgr)
}

// TXs
// - 1771d234f38e13fd9e4672fe469342fd598b6a2931f311d01b12dd4f35e9ce5c - 0.1 BTC - asg-9lf
// - 5c4ad18723fe385946288574760b2d563f52a8917cdaf850d66958cd472db07a - 0.1 BTC - asg-9lf
// - 96eca0eb4be36ac43fa2b2488fd3468aa2079ae02ae361ef5c08a4ace5070ed1 - 0.2 BTC - asg-9lf
// - 5338aa51f6a7ce8e7f7bc4c98ac06b47c50a3cf335d61e69cf06c0e11b945ea5 - 0.2 BTC - asg-9lf
// - 63d92b111b9dc1b09e030d5a853120917e6205ed43d536a25a335ae96930469d - 0.2 BTC - asg-9lf
// - 6a747fdf782fa87693183b865b261f39b32790df4b6959482c4c8d16c54c1273 - 0.2 BTC - asg-9lf
// - 4209f36cb89ff216fcf6b02f6badf22d64f1596a876c9805a9d6978c4e09190a - 0.2 BTC - asg-9lf
// - f09faaec7d3f84e89ef184bcf568e44b39296b2ad55d464743dd2a656720e6c1 - 0.2 BTC - asg-qev
// - ec7e201eda9313a434313376881cb61676b8407960df2d7cc9d17e65cbc8ba82 - 0.2 BTC - asg-qev

// Asgards
// - 9lf: 1.2 BTC (bc1q8my83gyjy95dya9e0j8vzsjz5hz786zll9w9lf) pubkey (thorpub1addwnpepqdlyqz7renj8u8hqsvynxwgwnfufcwmh7ttsx5n259cva8nctwre5qx29zu)
// - qev 0.4 BTC (bc1qe65v2vmxnplwfg8z0fwsps79sly2wrfn5tlqev) pubkey (thorpub1addwnpepqw6ckwjel98vpsfyd2cq6cvwdeqh6jfcshnsgdlpzng6uhg3f69ssawhg99)
func createFakeTxInsAndMakeObservations(ctx cosmos.Context, mgr *Mgrs) {
	userAddr, err := common.NewAddress("bc1qqfmzftwe7xtfjq5ukwar59yk9ts40u42mnznwr")
	if err != nil {
		ctx.Logger().Error("fail to create addr", "addr", userAddr.String(), "error", err)
		return
	}
	asg9lf, err := common.NewAddress("bc1q8my83gyjy95dya9e0j8vzsjz5hz786zll9w9lf")
	if err != nil {
		ctx.Logger().Error("fail to create addr", "addr", asg9lf.String(), "error", err)
		return
	}
	asg9lfPubKey, err := common.NewPubKey("thorpub1addwnpepqdlyqz7renj8u8hqsvynxwgwnfufcwmh7ttsx5n259cva8nctwre5qx29zu")
	if err != nil {
		ctx.Logger().Error("fail to create pubkey for vault", "addr", asg9lf.String(), "error", err)
		return
	}
	asgQev, err := common.NewAddress("bc1qe65v2vmxnplwfg8z0fwsps79sly2wrfn5tlqev")
	if err != nil {
		ctx.Logger().Error("fail to create addr", "addr", asgQev.String(), "error", err)
		return
	}
	asgQevPubKey, err := common.NewPubKey("thorpub1addwnpepqw6ckwjel98vpsfyd2cq6cvwdeqh6jfcshnsgdlpzng6uhg3f69ssawhg99")
	if err != nil {
		ctx.Logger().Error("fail to create pubkey for vault", "addr", asg9lf.String(), "error", err)
		return
	}

	// include savers add memo
	memo := "+:BTC/BTC"
	blockHeight := ctx.BlockHeight()

	unobservedTxs := ObservedTxs{
		NewObservedTx(common.Tx{
			ID:          "1771d234f38e13fd9e4672fe469342fd598b6a2931f311d01b12dd4f35e9ce5c",
			Chain:       common.BTCChain,
			FromAddress: userAddr,
			ToAddress:   asg9lf,
			Coins: common.NewCoins(common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(0.1 * common.One),
			}),
			Gas: common.Gas{common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1),
			}},
			Memo: memo,
		}, blockHeight, asg9lfPubKey, blockHeight),
		NewObservedTx(common.Tx{
			ID:          "5c4ad18723fe385946288574760b2d563f52a8917cdaf850d66958cd472db07a",
			Chain:       common.BTCChain,
			FromAddress: userAddr,
			ToAddress:   asg9lf,
			Coins: common.NewCoins(common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(0.1 * common.One),
			}),
			Gas: common.Gas{common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1),
			}},
			Memo: memo,
		}, blockHeight, asg9lfPubKey, blockHeight),
		NewObservedTx(common.Tx{
			ID:          "96eca0eb4be36ac43fa2b2488fd3468aa2079ae02ae361ef5c08a4ace5070ed1",
			Chain:       common.BTCChain,
			FromAddress: userAddr,
			ToAddress:   asg9lf,
			Coins: common.NewCoins(common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(0.2 * common.One),
			}),
			Gas: common.Gas{common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1),
			}},
			Memo: memo,
		}, blockHeight, asg9lfPubKey, blockHeight),
		NewObservedTx(common.Tx{
			ID:          "5338aa51f6a7ce8e7f7bc4c98ac06b47c50a3cf335d61e69cf06c0e11b945ea5",
			Chain:       common.BTCChain,
			FromAddress: userAddr,
			ToAddress:   asg9lf,
			Coins: common.NewCoins(common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(0.2 * common.One),
			}),
			Gas: common.Gas{common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1),
			}},
			Memo: memo,
		}, blockHeight, asg9lfPubKey, blockHeight),
		NewObservedTx(common.Tx{
			ID:          "63d92b111b9dc1b09e030d5a853120917e6205ed43d536a25a335ae96930469d",
			Chain:       common.BTCChain,
			FromAddress: userAddr,
			ToAddress:   asg9lf,
			Coins: common.NewCoins(common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(0.2 * common.One),
			}),
			Gas: common.Gas{common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1),
			}},
			Memo: memo,
		}, blockHeight, asg9lfPubKey, blockHeight),
		NewObservedTx(common.Tx{
			ID:          "6a747fdf782fa87693183b865b261f39b32790df4b6959482c4c8d16c54c1273",
			Chain:       common.BTCChain,
			FromAddress: userAddr,
			ToAddress:   asg9lf,
			Coins: common.NewCoins(common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(0.2 * common.One),
			}),
			Gas: common.Gas{common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1),
			}},
			Memo: memo,
		}, blockHeight, asg9lfPubKey, blockHeight),
		NewObservedTx(common.Tx{
			ID:          "4209f36cb89ff216fcf6b02f6badf22d64f1596a876c9805a9d6978c4e09190a",
			Chain:       common.BTCChain,
			FromAddress: userAddr,
			ToAddress:   asg9lf,
			Coins: common.NewCoins(common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(0.2 * common.One),
			}),
			Gas: common.Gas{common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1),
			}},
			Memo: memo,
		}, blockHeight, asg9lfPubKey, blockHeight),
		NewObservedTx(common.Tx{
			ID:          "f09faaec7d3f84e89ef184bcf568e44b39296b2ad55d464743dd2a656720e6c1",
			Chain:       common.BTCChain,
			FromAddress: userAddr,
			ToAddress:   asgQev,
			Coins: common.NewCoins(common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(0.2 * common.One),
			}),
			Gas: common.Gas{common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1),
			}},
			Memo: memo,
		}, blockHeight, asgQevPubKey, blockHeight),
		NewObservedTx(common.Tx{
			ID:          "ec7e201eda9313a434313376881cb61676b8407960df2d7cc9d17e65cbc8ba82",
			Chain:       common.BTCChain,
			FromAddress: userAddr,
			ToAddress:   asgQev,
			Coins: common.NewCoins(common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(0.2 * common.One),
			}),
			Gas: common.Gas{common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1),
			}},
			Memo: memo,
		}, blockHeight, asgQevPubKey, blockHeight),
	}

	err = makeFakeTxInObservation(ctx, mgr, unobservedTxs)
	if err != nil {
		ctx.Logger().Error("failed to migrate v109", "error", err)
	}
}

func migrateStoreV110(ctx cosmos.Context, mgr *Mgrs) {
	resetObservationHeights(ctx, mgr, 110, common.BTCChain, 788640)
}

func migrateStoreV111(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v111", "error", err)
		}
	}()

	// these were the node addresses missed in the last migration
	bech32Addrs := []string{
		"thor10rgvc7c44mq5vpcq07dx5fg942eykagm9p6gxh",
		"thor12espg8k5fxqmclx9vyte7cducmmvrtxll40q7z",
		"thor169fahg7x70vkv909h06c2mspphrzqgy7g6prr4",
		"thor1gukvqaag4vk2l3uq3kjme5x9xy8556pgv5rw4k",
		"thor1h6h54d7jutljwt46qzt2w7nnyuswwv045kmshl",
		"thor1raylctzthcvjc0a5pv5ckzjr3rgxk5qcwu7af2",
		"thor1s76zxv0kpr78za293kvj0eep4tfqljacknsjzc",
		"thor1w8mntay3xuk3c77j8fgvyyt0nfvl2sk398a3ww",
	}

	for _, addr := range bech32Addrs {
		// convert to cosmos address
		na, err := ctypes.AccAddressFromBech32(addr)
		if err != nil {
			ctx.Logger().Error("failed to convert bech32 address", "address", addr, "error", err)
			continue
		}

		// set observation height back
		mgr.Keeper().ForceSetLastObserveHeight(ctx, common.BTCChain, na, 788640)
	}
}

func migrateStoreV113(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v113", "error", err)
		}
	}()

	// block: 11227005, tx: 5AC64AC48219456C8701E67CB4E6ACA13495F8A8042EBC0E5B4E9DA9CF963A9B

	poolSlashRune := cosmos.NewUint(8101892874988)
	poolSlashBTC := cosmos.NewUint(297035619)

	// send coins from pool to bond module
	if err := mgr.Keeper().SendFromModuleToModule(ctx, AsgardName, BondName, common.Coins{common.NewCoin(common.RuneNative, poolSlashRune)}); err != nil {
		ctx.Logger().Error("fail to transfer coin from reserve to bond module", "error", err)
		return
	}

	// send coins from reserve to bond module
	if err := mgr.Keeper().SendFromModuleToModule(ctx, ReserveName, BondName, common.Coins{common.NewCoin(common.RuneNative, poolSlashRune)}); err != nil {
		ctx.Logger().Error("fail to transfer coin from reserve to bond module", "error", err)
		return
	}

	// revert pool slash
	pool, err := mgr.Keeper().GetPool(ctx, common.BTCAsset)
	if err != nil {
		ctx.Logger().Error("fail to get pool", "error", err)
		return
	}
	pool.BalanceAsset = pool.BalanceAsset.Add(poolSlashBTC)
	pool.BalanceRune = common.SafeSub(pool.BalanceRune, poolSlashRune)

	// store updated pool
	if err := mgr.Keeper().SetPool(ctx, pool); err != nil {
		ctx.Logger().Error("fail to set pool", "error", err)
		return
	}

	// emit inverted slash event for midgard
	poolSlashAmt := []PoolAmt{
		{
			Asset:  common.BTCAsset,
			Amount: int64(poolSlashBTC.Uint64()),
		},
		{
			Asset:  common.RuneAsset(),
			Amount: 0 - int64(poolSlashRune.Uint64()),
		},
	}
	eventSlash := NewEventSlash(common.BTCAsset, poolSlashAmt)
	if err := mgr.EventMgr().EmitEvent(ctx, eventSlash); err != nil {
		ctx.Logger().Error("fail to emit slash event", "error", err)
	}

	// credits from node slashes (sum to 2x the RUNE amount from pool slash)
	credits := []struct {
		address string
		amount  cosmos.Uint
	}{
		{address: "thor10rgvc7c44mq5vpcq07dx5fg942eykagm9p6gxh", amount: cosmos.NewUint(956154881499)},
		{address: "thor1pt8zkvkccj4397kemxeq8sjcyl7y6vacaedpvx", amount: cosmos.NewUint(761044063699)},
		{address: "thor1nlsfq25y74u8qt2hqmuzh5wd9t4uv28ghc258g", amount: cosmos.NewUint(973107929821)},
		{address: "thor1u5pfv07xtxz6aj59pnejaxh2dy7ew5s79ds8cw", amount: cosmos.NewUint(1063814699290)},
		{address: "thor1ypjwdplx07vf42qdfkex39dp8zxqnaects270v", amount: cosmos.NewUint(917937526969)},
		{address: "thor1vt207wgvefjgk88mtfjuurcl3vw6z4d2gu5psw", amount: cosmos.NewUint(1000265002165)},
		{address: "thor1vp29289yyvfar0ektscjk08r0tufvl24tn6xf9", amount: cosmos.NewUint(1021124834581)},
		{address: "thor1u9dnzza6hpesrwq4p8j2f29v6jsyeq4le66j3c", amount: cosmos.NewUint(978832200788)},
		{address: "thor1xk362wwunmr0gzew05j3euvdkjcvfmfyhmzd82", amount: cosmos.NewUint(1010886872701)},
		{address: "thor183fwfzgdfxzf5338acw32kplscgltf28j9s68j", amount: cosmos.NewUint(966449181925)},
		{address: "thor170xscqs5d469chdt83fxatjntc79zucrygsfxj", amount: cosmos.NewUint(1083603612921)},
		{address: "thor12espg8k5fxqmclx9vyte7cducmmvrtxll40q7z", amount: cosmos.NewUint(996350776100)},
		{address: "thor18nlluv0zw5g8930sx3r5xn7tqpsvwd7axxfynv", amount: cosmos.NewUint(1027540783824)},
		{address: "thor1faa0c6sqryr0am6ls9u8y6zs22ju2y7yw8ju9g", amount: cosmos.NewUint(603270429244)},
		{address: "thor1dqlmsm67h363nuxpd68esg54kt2t7xw2xewqml", amount: cosmos.NewUint(973292135373)},
		{address: "thor1gukvqaag4vk2l3uq3kjme5x9xy8556pgv5rw4k", amount: cosmos.NewUint(986737992322)},
		{address: "thor10f40m6nv7ulc0fvhmt07szn3n7ajd7e8xhghc3", amount: cosmos.NewUint(883372826754)},
	}

	for _, credit := range credits {
		ctx.Logger().Info("credit", "node", credit.address, "amount", credit.amount)

		// get addresses
		addr, err := cosmos.AccAddressFromBech32(credit.address)
		if err != nil {
			ctx.Logger().Error("fail to parse node address", "error", err)
			return
		}

		// get node account
		na, err := mgr.Keeper().GetNodeAccount(ctx, addr)
		if err != nil {
			ctx.Logger().Error("fail to get node account", "error", err)
			return
		}

		// update node bond
		na.Bond = na.Bond.Add(credit.amount)

		// store updated records
		if err := mgr.Keeper().SetNodeAccount(ctx, na); err != nil {
			ctx.Logger().Error("fail to save node account", "error", err)
			return
		}
	}
}

func migrateStoreV114(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV116(ctx cosmos.Context, mgr *Mgrs) {
	// query /thorchain/invariant/[bond,asgard,thorchain]
	bondRuneOver := cosmos.NewUint(6936532592883)
	asgardRuneUnder := cosmos.NewUint(5082320319988)
	thorchainRuneOver := cosmos.NewUint(100000000)

	// non-gas synth assets
	avaxUSDC, _ := common.NewAsset("avax/usdc-0xb97ef9ef8734c71904d8002f8b6bc66dd9c48a6e")
	ethFOX, _ := common.NewAsset("eth/fox-0xc770eefad204b5180df6a14ee197d99d808ee52d")
	ethUSDC, _ := common.NewAsset("eth/usdc-0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")
	ethUSDT, _ := common.NewAsset("eth/usdt-0xdac17f958d2ee523a2206206994597c13d831ec7")
	terraUST, _ := common.NewAsset("terra/ust")

	actions := []ModuleBalanceAction{
		// send rune from bond oversolvency to fix asgard insolvency
		{
			ModuleName:     BondName,
			RuneRecipient:  AsgardName,
			RuneToTransfer: asgardRuneUnder,
			SynthsToBurn:   common.Coins{},
		},

		// send remaining bond rune oversolvency to reserve
		{
			ModuleName:     BondName,
			RuneRecipient:  ReserveName,
			RuneToTransfer: common.SafeSub(bondRuneOver, asgardRuneUnder),
			SynthsToBurn:   common.Coins{},
		},

		// burn synths from asgard to fix oversolvencies
		{
			ModuleName:     AsgardName,
			RuneRecipient:  AsgardName, // noop
			RuneToTransfer: cosmos.ZeroUint(),
			SynthsToBurn: common.Coins{
				{
					Asset:  common.AVAXAsset.GetSyntheticAsset(),
					Amount: cosmos.NewUint(1000001),
				},
				{
					Asset:  avaxUSDC,
					Amount: cosmos.NewUint(4581),
				},
				{
					Asset:  common.BCHAsset.GetSyntheticAsset(),
					Amount: cosmos.NewUint(20529),
				},
				{
					Asset:  common.BNBAsset.GetSyntheticAsset(),
					Amount: cosmos.NewUint(99999),
				},
				{
					Asset:  common.BTCAsset.GetSyntheticAsset(),
					Amount: cosmos.NewUint(942067),
				},
				{
					Asset:  common.DOGEAsset.GetSyntheticAsset(),
					Amount: cosmos.NewUint(4400752724),
				},
				{
					Asset:  common.ETHAsset.GetSyntheticAsset(),
					Amount: cosmos.NewUint(4455527),
				},
				{
					Asset:  ethFOX,
					Amount: cosmos.NewUint(215666666666),
				},
				{
					Asset:  ethUSDC,
					Amount: cosmos.NewUint(21884753549),
				},
				{
					Asset:  ethUSDT,
					Amount: cosmos.NewUint(281542),
				},
				{
					Asset:  common.ATOMAsset.GetSyntheticAsset(),
					Amount: cosmos.NewUint(1039626),
				},
				{
					Asset:  common.LUNAAsset.GetSyntheticAsset(),
					Amount: cosmos.NewUint(2527),
				},
				{
					Asset:  terraUST,
					Amount: cosmos.NewUint(29955102645),
				},
			},
		},

		// transfer rune from thorchain to reserve to clear thorchain balances
		{
			ModuleName:     ModuleName,
			RuneRecipient:  ReserveName,
			RuneToTransfer: thorchainRuneOver,
			SynthsToBurn:   common.Coins{},
		},
	}

	processModuleBalanceActions(ctx, mgr.Keeper(), actions)
}

func migrateStoreV117(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v117", "error", err)
		}
	}()

	subBalance := func(ctx cosmos.Context, mgr *Mgrs, amountUint64 uint64, assetString, pkString string) {
		amount := cosmos.NewUint(amountUint64)

		asset, err := common.NewAsset(assetString)
		if err != nil {
			ctx.Logger().Error("fail to make asset", "error", err)
			return
		}

		pubkey, err := common.NewPubKey(pkString)
		if err != nil {
			ctx.Logger().Error("fail to make pubkey", "error", err)
			return
		}

		vault, err := mgr.Keeper().GetVault(ctx, pubkey)
		if err != nil {
			ctx.Logger().Error("fail to get vault", "error", err)
			return
		}

		coins := common.NewCoins(common.NewCoin(asset, amount))
		vault.SubFunds(coins)

		if err := mgr.Keeper().SetVault(ctx, vault); err != nil {
			ctx.Logger().Error("fail to set vault", "error", err)
			return
		}
	}

	subBalance(ctx, mgr, 636462, "ETH.ETH", "thorpub1addwnpepqf654umpm7vzgmegae0k4yq0xe69kpvvp3w437hvy7rpyk3svxgtszw2tu9")
	subBalance(ctx, mgr, 8640, "ETH.ETH", "thorpub1addwnpepqfx4cxtsthazf8609lhfcxxlu60er2t90utta66u2xz2xtdhpts9slmkc93")
}

func migrateStoreV121(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v121", "error", err)
		}
	}()

	// For any in-progress streaming swaps to non-RUNE Native coins,
	// mint the current Out amount to the Pool Module.
	var coinsToMint common.Coins

	iterator := mgr.Keeper().GetSwapQueueIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var msg MsgSwap
		if err := mgr.Keeper().Cdc().Unmarshal(iterator.Value(), &msg); err != nil {
			ctx.Logger().Error("fail to fetch swap msg from queue", "error", err)
			continue
		}

		if !msg.IsStreaming() || !msg.TargetAsset.IsNative() || msg.TargetAsset.IsRune() {
			continue
		}

		swp, err := mgr.Keeper().GetStreamingSwap(ctx, msg.Tx.ID)
		if err != nil {
			ctx.Logger().Error("fail to fetch streaming swap", "error", err)
			continue
		}

		if !swp.Out.IsZero() {
			mintCoin := common.NewCoin(msg.TargetAsset, swp.Out)
			coinsToMint = coinsToMint.Add_deprecated(mintCoin)
		}
	}

	// The minted coins are for in-progress swaps, so keeping the "swap" in the event field and logs.
	var coinsToTransfer common.Coins
	for _, mintCoin := range coinsToMint {
		if err := mgr.Keeper().MintToModule(ctx, ModuleName, mintCoin); err != nil {
			ctx.Logger().Error("fail to mint coins during swap", "error", err)
		} else {
			mintEvt := NewEventMintBurn(MintSupplyType, mintCoin.Asset.Native(), mintCoin.Amount, "swap")
			if err := mgr.EventMgr().EmitEvent(ctx, mintEvt); err != nil {
				ctx.Logger().Error("fail to emit mint event", "error", err)
			}
			coinsToTransfer = coinsToTransfer.Add_deprecated(mintCoin)
		}
	}

	if len(coinsToTransfer) > 0 {
		if err := mgr.Keeper().SendFromModuleToModule(ctx, ModuleName, AsgardName, coinsToTransfer); err != nil {
			ctx.Logger().Error("fail to move coins during swap", "error", err)
		}
	}
}

func mainnetUnobservedTxsV122(ctx cosmos.Context, mgr *Mgrs) {
	unobservedTxs := ObservedTxs{}

	// Manually refunded by treasury in:
	// https://blockstream.info/tx/50442a094f14d937056c17697f9e1909fc6fdfea980a48e3caa991bde983d4f6
	// Fake observation will use the treasury address as the sender for refund.
	treasuryAddr, err := common.NewAddress("bc1qq2z2f4gs4nd7t0a9jjp90y9l9zzjtegu4nczha")
	if err != nil {
		ctx.Logger().Error("fail to create addr", "addr", treasuryAddr.String(), "error", err)
		return
	}
	asgjh5h, err := common.NewAddress("bc1qj9trqrwtp4c8c5hpt63pyavqkgk90av2zsc0vj")
	if err != nil {
		ctx.Logger().Error("fail to create addr", "addr", asgjh5h.String(), "error", err)
		return
	}
	asgjh5hPubKey, err := common.NewPubKey("thorpub1addwnpepqgq37yke5ya53rkwx57z65zv0k8e80paxncfpgfzr83pfz9atywdgvgjh5h")
	if err != nil {
		ctx.Logger().Error("fail to create pubkey for vault", "addr", asgjh5h.String(), "error", err)
		return
	}
	unobservedTxs = append(
		unobservedTxs,
		NewObservedTx(common.Tx{
			ID:          "919a5bd8a69426d61f141aa93423e33a77c608986c92e124cc2444afe38ef6aa",
			Chain:       common.BTCChain,
			FromAddress: treasuryAddr,
			ToAddress:   asgjh5h,
			Coins: common.NewCoins(common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(39000000),
			}),
			Gas: common.Gas{common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1),
			}},
			Memo: "",
		}, 807061, asgjh5hPubKey, 807061),
	)

	// Refund bRUTE (Discord).
	userAddr, err := common.NewAddress("bc1q7cem2mjtl7uzk67xa9n9pvawamwhcttqfryntf")
	if err != nil {
		ctx.Logger().Error("fail to create addr", "addr", treasuryAddr.String(), "error", err)
		return
	}
	asg5mza, err := common.NewAddress("bc1qjx8uw8l3vkf59l493nz23vq6mhtf8k5lu2yrqq")
	if err != nil {
		ctx.Logger().Error("fail to create addr", "addr", asg5mza.String(), "error", err)
		return
	}
	asg5mzaPubKey, err := common.NewPubKey("thorpub1addwnpepq0j3d5j45kkfnf0arkqxx0g40zmalhp5uxmg2fnfzgt7tz8mp5f9chd5mza")
	if err != nil {
		ctx.Logger().Error("fail to create pubkey for vault", "addr", asg5mza.String(), "error", err)
		return
	}
	unobservedTxs = append(
		unobservedTxs,
		NewObservedTx(common.Tx{
			ID:          "3e9bf3b1e92732d75a58d42d43a07f9f45f3549daa18f57f021741a3a56a3414",
			Chain:       common.BTCChain,
			FromAddress: userAddr,
			ToAddress:   asg5mza,
			Coins: common.NewCoins(common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(4000000),
			}),
			Gas: common.Gas{common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1),
			}},
			Memo: "",
		}, 808099, asg5mzaPubKey, 808099),
	)

	err = makeFakeTxInObservation(ctx, mgr, unobservedTxs)
	if err != nil {
		ctx.Logger().Error("failed to migrate v122", "error", err)
	}
}

// The following refunds are for the instances of Ethereum double spends on 2023/10/1.
// The reserve will refund the full bond slash and leave the slash that went to pools.
// Amounts were taken from:
// tci nodes slash-diff --height 12687242 --blocks 1
// tci nodes slash-diff --height 12836219 --blocks 1
// tci nodes slash-diff --height 12837936 --blocks 1
//
// Totals for the above verified to be slash event changes with the following (amounts
// taken from slash-diff instead since the events contain the operator address and not
// the node address):
//
//	http https://thornode-v1.ninerealms.com/thorchain/block height==12687242 | \
//		jq -c '[.txs[]|.result.events[]|select((.type=="bond") and (.bond_type=="\u0003"))|.amount|tonumber]|add'
//
//	http https://thornode-v1.ninerealms.com/thorchain/block height==12836219 | \
//		jq -c '[.txs[]|.result.events[]|select((.type=="bond") and (.bond_type=="\u0003"))|.amount|tonumber]|add'
//
//	http https://thornode-v1.ninerealms.com/thorchain/block height==12837936 | \
//		jq -c '[.txs[]|.result.events[]|select((.type=="bond") and (.bond_type=="\u0003"))|.amount|tonumber]|add'
func mainnetBondRefundsV122(ctx cosmos.Context, mgr *Mgrs) {
	credits := []struct {
		address string
		amount  cosmos.Uint
	}{
		// height: 12687242
		{address: "thor1dwt6szf098rd4vlnjn83w249zky3penf76cuxy", amount: cosmos.NewUint(84168028924)},
		{address: "thor1muc7w8s4k2v94lz9mhda5dav9nyyf9c9959g89", amount: cosmos.NewUint(90072158640)},
		{address: "thor1nlxtkz6wjrsz3wcez0vz577kl6xx7m5mdmysvy", amount: cosmos.NewUint(82448629506)},
		{address: "thor183fwfzgdfxzf5338acw32kplscgltf28j9s68j", amount: cosmos.NewUint(83016415361)},
		{address: "thor1z3dmy779shx8x9903ldnyqnt3a3g6vjqx68hkt", amount: cosmos.NewUint(91153484540)},
		{address: "thor1nw2jdqn5u8xsx4j0n4e8cmndapxqj47z8zhcs3", amount: cosmos.NewUint(95306185316)},
		{address: "thor1raylctzthcvjc0a5pv5ckzjr3rgxk5qcwu7af2", amount: cosmos.NewUint(85327552827)},
		{address: "thor1u9dnzza6hpesrwq4p8j2f29v6jsyeq4le66j3c", amount: cosmos.NewUint(88209401326)},
		{address: "thor1xd4j3gk9frpxh8r22runntnqy34lwzrdkazldh", amount: cosmos.NewUint(89740664937)},
		{address: "thor1zga95gkv87356lmjj0mvw3geylfuv3ph7wa9t0", amount: cosmos.NewUint(90932733476)},
		{address: "thor18fqat7ta4mdxlzq8xdhuel23ng7plm00qrzdre", amount: cosmos.NewUint(81748865505)},
		{address: "thor186k9w7hw4zdmd0kyqsrfzhzgvpc8v9shd9qe7u", amount: cosmos.NewUint(90544695864)},
		{address: "thor1aulde7ynkh8jd9qpuxw5srafew00vpauw5np52", amount: cosmos.NewUint(76386938018)},
		{address: "thor1faa0c6sqryr0am6ls9u8y6zs22ju2y7yw8ju9g", amount: cosmos.NewUint(71628288242)},
		{address: "thor1gqtwzazgdncthm2cuu947d0mvk3w5fkahm40qp", amount: cosmos.NewUint(86683624579)},
		{address: "thor1lgms9fnlgz8den685z0fs5f2vm60jauvkjf6pm", amount: cosmos.NewUint(91474396878)},
		{address: "thor1lxm4ahz43va3s2mwyed63l5k0mua0ecr9qhmmm", amount: cosmos.NewUint(85358066170)},
		// height: 12836219
		{address: "thor1pcylx2quurhr44fg35jgvlrypvag70aszgd2t3", amount: cosmos.NewUint(815835260165)},
		{address: "thor12jrhy6mqxtff6utq4kkavtvmqz4qxtztxxnk4j", amount: cosmos.NewUint(786637113773)},
		{address: "thor1errw9wx5pv8rhevexfxa950jx6tux0qywrlwlp", amount: cosmos.NewUint(739209960394)},
		{address: "thor1z0zph6u9e00y407d7vg7dh2c5knz5an0wzv8v4", amount: cosmos.NewUint(829270201151)},
		{address: "thor13xa9eseegag6lcg4qa3eaj7uhljc9kmtxld84v", amount: cosmos.NewUint(775486690428)},
		{address: "thor1xd4j3gk9frpxh8r22runntnqy34lwzrdkazldh", amount: cosmos.NewUint(808331250813)},
		{address: "thor1nlsfq25y74u8qt2hqmuzh5wd9t4uv28ghc258g", amount: cosmos.NewUint(752359380806)},
		{address: "thor12z69uvtwxlj2j9c5cqrnnfqy7s2twrqmvqvj20", amount: cosmos.NewUint(801803186341)},
		{address: "thor1qp8288u08r2da9sj9pkzv3fkh0ugfutkl9gqdj", amount: cosmos.NewUint(765278325822)},
		{address: "thor1jnmj9jszmwjctxfd5leczrvl7kdcml3uyq43yn", amount: cosmos.NewUint(808893248809)},
		{address: "thor1haadhysqf9z5hq92eya78e89qehx0wkpm3jkgu", amount: cosmos.NewUint(875592259674)},
		{address: "thor1hsga5e2ul8jsy4t6tnxuqakulhmxs32ln7zy2g", amount: cosmos.NewUint(747523029022)},
		{address: "thor1hue0dwzd3lsxyq3qgecyzzmxrhq96qytwdvwj0", amount: cosmos.NewUint(789507173516)},
		{address: "thor1krcz33mejvc5f6grj2c5w5x3kuj7mnjhgqltj8", amount: cosmos.NewUint(746423338852)},
		{address: "thor1gqtwzazgdncthm2cuu947d0mvk3w5fkahm40qp", amount: cosmos.NewUint(780797985626)},
		{address: "thor1raylctzthcvjc0a5pv5ckzjr3rgxk5qcwu7af2", amount: cosmos.NewUint(768569889338)},
		{address: "thor1vwqz5hhh5un28qlz6x5f8zczj39jqwel38q2kc", amount: cosmos.NewUint(875216599105)},
		{address: "thor1sqf8fjuj050wq3m2p83af8l93g7s6ucn42eqa0", amount: cosmos.NewUint(819428936375)},
		// height: 12837936
		{address: "thor1yak0z56elhcfqw7xn7wjmp43ndnnxgfcmwkwex", amount: cosmos.NewUint(61353663980)},
		{address: "thor18fqat7ta4mdxlzq8xdhuel23ng7plm00qrzdre", amount: cosmos.NewUint(52153321557)},
		{address: "thor1nw2jdqn5u8xsx4j0n4e8cmndapxqj47z8zhcs3", amount: cosmos.NewUint(59357869146)},
		{address: "thor1asnulx9f4hr8e8fsa40wg3yxsyrdesj38vwndn", amount: cosmos.NewUint(54543153298)},
		{address: "thor1n5ylq3kyylr7jrq6zdksy2jrtyffxssra22tm3", amount: cosmos.NewUint(54698705550)},
		{address: "thor12g0es965kj3nql8k244unkznqx37r23ytns75x", amount: cosmos.NewUint(52590255978)},
		{address: "thor1agftrgu74z84hef6dt6ykhe7cmjf3f8dcpkfun", amount: cosmos.NewUint(55325194037)},
		{address: "thor1h3pvd8x44v63qj488lku6pzzcq3g5p8tc2nd6c", amount: cosmos.NewUint(52068842540)},
		{address: "thor1dwt6szf098rd4vlnjn83w249zky3penf76cuxy", amount: cosmos.NewUint(52338666870)},
		{address: "thor1ffz7rvtjvckuj3l05n4xp55v4zsqpxavej9dtr", amount: cosmos.NewUint(55964770042)},
		{address: "thor1aulde7ynkh8jd9qpuxw5srafew00vpauw5np52", amount: cosmos.NewUint(47599879644)},
		{address: "thor13r9p8upgtpff05nxy2kagy70qe0ljumhxe6qyf", amount: cosmos.NewUint(55957461500)},
		{address: "thor1dqlmsm67h363nuxpd68esg54kt2t7xw2xewqml", amount: cosmos.NewUint(53875467917)},
		{address: "thor1muc7w8s4k2v94lz9mhda5dav9nyyf9c9959g89", amount: cosmos.NewUint(56127352503)},
		{address: "thor10f40m6nv7ulc0fvhmt07szn3n7ajd7e8xhghc3", amount: cosmos.NewUint(49100123844)},
		{address: "thor1ucwcatnqwjyucfrf7vv2xnmfzfaplvrkqzl337", amount: cosmos.NewUint(56999681069)},
		{address: "thor1vt207wgvefjgk88mtfjuurcl3vw6z4d2gu5psw", amount: cosmos.NewUint(55682338872)},
		{address: "thor1kj56aupxnkhhy0rpdcp2gjncm4y78nnhjv496v", amount: cosmos.NewUint(56017963303)},
	}

	// sum amounts to get the total we will refund to nodes from the reserve
	total := cosmos.ZeroUint()
	for _, credit := range credits {
		total = total.Add(credit.amount)
	}

	// assertion for sanity check (~167k RUNE)
	if !total.Equal(cosmos.NewUint(16732118671769)) {
		ctx.Logger().Error("total refund amount is not correct", "total", total)
		return
	}

	// send coins from reserve to bond module
	if err := mgr.Keeper().SendFromModuleToModule(ctx, ReserveName, BondName, common.Coins{common.NewCoin(common.RuneNative, total)}); err != nil {
		ctx.Logger().Error("fail to transfer coin from reserve to bond module", "error", err)
		return
	}

	for _, credit := range credits {
		ctx.Logger().Info("credit", "node", credit.address, "amount", credit.amount)

		// get addresses
		addr, err := cosmos.AccAddressFromBech32(credit.address)
		if err != nil {
			ctx.Logger().Error("fail to parse node address", "error", err)
			return
		}

		// get node account
		na, err := mgr.Keeper().GetNodeAccount(ctx, addr)
		if err != nil {
			ctx.Logger().Error("fail to get node account", "error", err)
			return
		}

		// update node bond
		na.Bond = na.Bond.Add(credit.amount)

		// store updated records
		if err := mgr.Keeper().SetNodeAccount(ctx, na); err != nil {
			ctx.Logger().Error("fail to save node account", "error", err)
			return
		}
	}
}

func migrateStoreV122(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v122", "error", err)
		}
	}()

	mainnetUnobservedTxsV122(ctx, mgr)
	mainnetBondRefundsV122(ctx, mgr)
}

func migrateStoreV123(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v122", "error", err)
		}
	}()

	// Reque streaming swap outbound that got swallowed by preferred asset swap
	danglingInboundTxIDs := []common.TxID{
		"15D60EDF026A662E5D8BF3A36A50EBF9C0BD7B169340F669C4816376FCF5605E",
	}
	requeueDanglingActionsV123(ctx, mgr, danglingInboundTxIDs)

	// Requeue dropped attempt to rescue funds sent to old vault
	// Original tx: F9FA0745290D5EDB287F8641B390171B45BD84C7628A1A45DADB876F9359B4F8
	// Sent to bc1qkrcd6cfhmur80lsc0dxj2h3cge6lytaxdy5rl9
	// thorpub1addwnpepqdu0wrnvx63eqf6gf5qyfz2k95gj9c96hsgecdcfwawckvsgqy3ezh0rxwt
	originalTxID := "F9FA0745290D5EDB287F8641B390171B45BD84C7628A1A45DADB876F9359B4F8"
	droppedRescue := TxOutItem{
		Chain:       common.BTCChain,
		ToAddress:   common.Address("bc1q0vu0a7zpmgfrjuke7jeg5nlttfknsn2lee2qrx"),
		VaultPubKey: common.PubKey("thorpub1addwnpepqdu0wrnvx63eqf6gf5qyfz2k95gj9c96hsgecdcfwawckvsgqy3ezh0rxwt"),
		Coin:        common.NewCoin(common.BTCAsset, cosmos.NewUint(149896000)),
		Memo:        fmt.Sprintf("REFUND:%s", originalTxID),
		InHash:      common.TxID(originalTxID),
		GasRate:     94,
		MaxGas:      common.Gas{common.NewCoin(common.BTCAsset, cosmos.NewUint(94500))},
	}

	err := mgr.txOutStore.UnSafeAddTxOutItem(ctx, mgr, droppedRescue, ctx.BlockHeight())
	if err != nil {
		ctx.Logger().Error("fail to requeue BTC rescue tx", "error", err)
		return
	}
}

func migrateStoreV124(ctx cosmos.Context, mgr *Mgrs) {
	// Second attempt. Previously tried in V123 but amount was not fully spendable
	// Requeue dropped attempt to rescue funds sent to old vault
	// Original tx: F9FA0745290D5EDB287F8641B390171B45BD84C7628A1A45DADB876F9359B4F8
	// Sent to bc1qkrcd6cfhmur80lsc0dxj2h3cge6lytaxdy5rl9
	// thorpub1addwnpepqdu0wrnvx63eqf6gf5qyfz2k95gj9c96hsgecdcfwawckvsgqy3ezh0rxwt
	originalTxID := "F9FA0745290D5EDB287F8641B390171B45BD84C7628A1A45DADB876F9359B4F8"
	droppedRescue := TxOutItem{
		Chain:       common.BTCChain,
		ToAddress:   common.Address("bc1q0vu0a7zpmgfrjuke7jeg5nlttfknsn2lee2qrx"),
		VaultPubKey: common.PubKey("thorpub1addwnpepqdu0wrnvx63eqf6gf5qyfz2k95gj9c96hsgecdcfwawckvsgqy3ezh0rxwt"),
		Coin:        common.NewCoin(common.BTCAsset, cosmos.NewUint(147000000)),
		Memo:        fmt.Sprintf("REFUND:%s", originalTxID),
		InHash:      common.TxID(originalTxID),
		GasRate:     94,
		MaxGas:      common.Gas{common.NewCoin(common.BTCAsset, cosmos.NewUint(94500))},
	}

	err := mgr.txOutStore.UnSafeAddTxOutItem(ctx, mgr, droppedRescue, ctx.BlockHeight())
	if err != nil {
		ctx.Logger().Error("fail to requeue BTC rescue tx", "error", err)
		return
	}
}

// Bond refunds for AVAX double spend slash. See PR for details on determining amounts.
func migrateStoreV125(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v125", "error", err)
		}
	}()

	credits := []struct {
		address string
		amount  cosmos.Uint
	}{
		{address: "thor104un8h7jslr28xq7mxlrcljhnxtqcl0v4zl5pq", amount: cosmos.NewUint(319321625569)},
		{address: "thor10czf2s89h79fsjmqqck85cdqeq536hw5ngz4lt", amount: cosmos.NewUint(127832591722)},
		{address: "thor10fgzvdajq2f0gc2a5pmfh9up4qqajuk9je8lnc", amount: cosmos.NewUint(118728742428)},
		{address: "thor10rgvc7c44mq5vpcq07dx5fg942eykagm9p6gxh", amount: cosmos.NewUint(252712045018)},
		{address: "thor10zrxnnd75u2kwmdygehuurgapk36wtz9atsmkh", amount: cosmos.NewUint(89892458720)},
		{address: "thor1283lplant40dy3aq6k22rxamxuypqx9sk4qv2h", amount: cosmos.NewUint(145226714721)},
		{address: "thor12963n028s8gj5h048x7x5jjx8zjd6ev2uqysgn", amount: cosmos.NewUint(421019486106)},
		{address: "thor12espg8k5fxqmclx9vyte7cducmmvrtxll40q7z", amount: cosmos.NewUint(72210471615)},
		{address: "thor12g0es965kj3nql8k244unkznqx37r23ytns75x", amount: cosmos.NewUint(434885263871)},
		{address: "thor12jrhy6mqxtff6utq4kkavtvmqz4qxtztxxnk4j", amount: cosmos.NewUint(254177439761)},
		{address: "thor12nfq8smgr93mk845zlqpdel8r8mjk477g99426", amount: cosmos.NewUint(432212517832)},
		{address: "thor12z69uvtwxlj2j9c5cqrnnfqy7s2twrqmvqvj20", amount: cosmos.NewUint(430168608993)},
		{address: "thor13r9p8upgtpff05nxy2kagy70qe0ljumhxe6qyf", amount: cosmos.NewUint(148157042392)},
		{address: "thor13xa9eseegag6lcg4qa3eaj7uhljc9kmtxld84v", amount: cosmos.NewUint(250726994016)},
		{address: "thor140wms8h9pm5dmj832lwnhw45qvt25v0ps888sf", amount: cosmos.NewUint(73308903138)},
		{address: "thor14289vlld8lq7qcp67j6w66gfx4l76havr697m2", amount: cosmos.NewUint(117339493525)},
		{address: "thor14zyn9xdkv4au3frpj5ze6fkf4fwt80s7uy2394", amount: cosmos.NewUint(116666471790)},
		{address: "thor15hg7xk9k2rh0yyhn2atwj96h9srh4c6ys7vxn5", amount: cosmos.NewUint(422621892057)},
		{address: "thor15vzju96yvcpuhqk9u2mevdsud96k7cvjxhwk0e", amount: cosmos.NewUint(120985416555)},
		{address: "thor16ery22gma35h2fduxr0swdfvz4s6yvy6yhskf6", amount: cosmos.NewUint(72457288757)},
		{address: "thor16ta9xjecs0ju6w4u8f7udh4405lxpvje77kynk", amount: cosmos.NewUint(227093608519)},
		{address: "thor175jn909du3pwkwradsufuypzgjg4q5d8e26gh0", amount: cosmos.NewUint(515421323080)},
		{address: "thor17s8u4s635kee8g2u34htqxxg6jalvgvjpdpxsu", amount: cosmos.NewUint(120839256151)},
		{address: "thor186k9w7hw4zdmd0kyqsrfzhzgvpc8v9shd9qe7u", amount: cosmos.NewUint(149060273172)},
		{address: "thor18fqat7ta4mdxlzq8xdhuel23ng7plm00qrzdre", amount: cosmos.NewUint(247012541721)},
		{address: "thor18nlluv0zw5g8930sx3r5xn7tqpsvwd7axxfynv", amount: cosmos.NewUint(121984253317)},
		{address: "thor19dxkzzp09egc06wv4jcg5q2l6zy2yg3yyxplcw", amount: cosmos.NewUint(16674098338)},
		{address: "thor19m4kqulyqvya339jfja84h6qp8tkjgxuxa4n4a", amount: cosmos.NewUint(451694082161)},
		{address: "thor19xxxeetxrjvn2qchx00l4xxm0dlcwjkx0h2r82", amount: cosmos.NewUint(464165882575)},
		{address: "thor1afpy3lt4dh4x3jwvm0ncss7anm3zrd4wfwsk8n", amount: cosmos.NewUint(434511468328)},
		{address: "thor1agftrgu74z84hef6dt6ykhe7cmjf3f8dcpkfun", amount: cosmos.NewUint(116036966383)},
		{address: "thor1asnulx9f4hr8e8fsa40wg3yxsyrdesj38vwndn", amount: cosmos.NewUint(116659084404)},
		{address: "thor1aulde7ynkh8jd9qpuxw5srafew00vpauw5np52", amount: cosmos.NewUint(311365547994)},
		{address: "thor1c0grac9sstxey0jvhmzan9a9dplerres8rdg4n", amount: cosmos.NewUint(121015785523)},
		{address: "thor1dqlmsm67h363nuxpd68esg54kt2t7xw2xewqml", amount: cosmos.NewUint(84858891179)},
		{address: "thor1dwt6szf098rd4vlnjn83w249zky3penf76cuxy", amount: cosmos.NewUint(66724389440)},
		{address: "thor1e5v06j4w5u683n6yypg6apcuvft264qdqm5ayd", amount: cosmos.NewUint(391827212146)},
		{address: "thor1ee5ec0mnvhqvu84lgtxpgcc4yrdhalylfty246", amount: cosmos.NewUint(119818920029)},
		{address: "thor1ejtuux8f5pzg4h3q54l9hgw6j4969huh7xpamz", amount: cosmos.NewUint(60422728898)},
		{address: "thor1en5cc7sahcy6phqrvcetyyxsu05h8d05yeg88x", amount: cosmos.NewUint(259902590738)},
		{address: "thor1errw9wx5pv8rhevexfxa950jx6tux0qywrlwlp", amount: cosmos.NewUint(71521162233)},
		{address: "thor1ffz7rvtjvckuj3l05n4xp55v4zsqpxavej9dtr", amount: cosmos.NewUint(112764711483)},
		{address: "thor1fsk42s4trwrrc7k9elwhr0njlt5jfqlw8269mz", amount: cosmos.NewUint(166051691388)},
		{address: "thor1fy4njncghzmuce87c63mrtwpdpyzusdlfr54k6", amount: cosmos.NewUint(18053756311)},
		{address: "thor1gqtwzazgdncthm2cuu947d0mvk3w5fkahm40qp", amount: cosmos.NewUint(238082090162)},
		{address: "thor1gukvqaag4vk2l3uq3kjme5x9xy8556pgv5rw4k", amount: cosmos.NewUint(260514693954)},
		{address: "thor1h3pvd8x44v63qj488lku6pzzcq3g5p8tc2nd6c", amount: cosmos.NewUint(393026086148)},
		{address: "thor1h6h54d7jutljwt46qzt2w7nnyuswwv045kmshl", amount: cosmos.NewUint(480572829252)},
		{address: "thor1haadhysqf9z5hq92eya78e89qehx0wkpm3jkgu", amount: cosmos.NewUint(468358544129)},
		{address: "thor1hpt4l30qgr3pugg2wdp4hmrv2g0qlg2z7z9m7h", amount: cosmos.NewUint(160670428442)},
		{address: "thor1hue0dwzd3lsxyq3qgecyzzmxrhq96qytwdvwj0", amount: cosmos.NewUint(17622544431)},
		{address: "thor1hx3gayvwx0j92nf0ev6c837km4yg4kdk0w6fjl", amount: cosmos.NewUint(516735010159)},
		{address: "thor1jayc3hvwmgyex2ftmg3hr2mk4ujjhrv9eua74x", amount: cosmos.NewUint(421026220242)},
		{address: "thor1jk93saw08vwmqdj684f4w5y6v4dkgdfy2jkrdp", amount: cosmos.NewUint(122635027611)},
		{address: "thor1k2e50ws3d9lce9ycr7ppaazx3ygaa7lxj8kkny", amount: cosmos.NewUint(72959319676)},
		{address: "thor1k42q8hvzk3r3uy0w2udaxy743d7u996s60jwv8", amount: cosmos.NewUint(90504693795)},
		{address: "thor1kchgh8t790zlfatdun975mu04xvumq3qjms65a", amount: cosmos.NewUint(130912256935)},
		{address: "thor1kj56aupxnkhhy0rpdcp2gjncm4y78nnhjv496v", amount: cosmos.NewUint(131295769218)},
		{address: "thor1krcz33mejvc5f6grj2c5w5x3kuj7mnjhgqltj8", amount: cosmos.NewUint(136535375373)},
		{address: "thor1lgms9fnlgz8den685z0fs5f2vm60jauvkjf6pm", amount: cosmos.NewUint(242603659660)},
		{address: "thor1luvyfs2r4cedmmv0fqr2dqdcr4j4fd764hquu3", amount: cosmos.NewUint(109433819474)},
		{address: "thor1lxm4ahz43va3s2mwyed63l5k0mua0ecr9qhmmm", amount: cosmos.NewUint(487542300305)},
		{address: "thor1lzzvchm4ldm66rem8u85n8ytj2nzhcnprmtwqr", amount: cosmos.NewUint(219635557647)},
		{address: "thor1m45tc3uw4egzw9v2j39x47ds926ynfducvt9fx", amount: cosmos.NewUint(17880282436)},
		{address: "thor1muc7w8s4k2v94lz9mhda5dav9nyyf9c9959g89", amount: cosmos.NewUint(429271927236)},
		{address: "thor1n5ylq3kyylr7jrq6zdksy2jrtyffxssra22tm3", amount: cosmos.NewUint(542248910653)},
		{address: "thor1nd4n9s9shgdp4lnn859zq6snx8pnp2u2zc2mqc", amount: cosmos.NewUint(453931314209)},
		{address: "thor1nl0hc33pllze4athmvyaj9ky35le30vvhx3a25", amount: cosmos.NewUint(125190793050)},
		{address: "thor1nlsfq25y74u8qt2hqmuzh5wd9t4uv28ghc258g", amount: cosmos.NewUint(110223575724)},
		{address: "thor1nlxtkz6wjrsz3wcez0vz577kl6xx7m5mdmysvy", amount: cosmos.NewUint(267941321436)},
		{address: "thor1nw2jdqn5u8xsx4j0n4e8cmndapxqj47z8zhcs3", amount: cosmos.NewUint(139279915982)},
		{address: "thor1p4mykaudddvnkzpfvtutn3vfjhyy43wktgxxf9", amount: cosmos.NewUint(353051378515)},
		{address: "thor1pcylx2quurhr44fg35jgvlrypvag70aszgd2t3", amount: cosmos.NewUint(562992832814)},
		{address: "thor1pszqlupqmp90w8w3368auraw9nnczxysr9em0l", amount: cosmos.NewUint(152131729277)},
		{address: "thor1pt8zkvkccj4397kemxeq8sjcyl7y6vacaedpvx", amount: cosmos.NewUint(519029350311)},
		{address: "thor1qp8288u08r2da9sj9pkzv3fkh0ugfutkl9gqdj", amount: cosmos.NewUint(314125641973)},
		{address: "thor1r0jtp4y6kr627fsgdj5mmqfumkk488sn4k2c3x", amount: cosmos.NewUint(443279651214)},
		{address: "thor1r6fmvdx85mq55qn59qgun6kyhlgvcy8nm46dz0", amount: cosmos.NewUint(528997248894)},
		{address: "thor1r9027rfu48kyvs0curxjur7wzywk36ukrk7mc4", amount: cosmos.NewUint(122415521033)},
		{address: "thor1raylctzthcvjc0a5pv5ckzjr3rgxk5qcwu7af2", amount: cosmos.NewUint(112240778405)},
		{address: "thor1s7lu6rfxgw0c9xypmrmjxzv0glva02cmn40rde", amount: cosmos.NewUint(144720033501)},
		{address: "thor1sn88hq7n85a5ju4x9pjghzgqu070h2epnyj53w", amount: cosmos.NewUint(130709199876)},
		{address: "thor1sngd0zz6pwdx2e20sml27354vzkrwa4fnjxvnc", amount: cosmos.NewUint(73083108474)},
		{address: "thor1sqf8fjuj050wq3m2p83af8l93g7s6ucn42eqa0", amount: cosmos.NewUint(119538313490)},
		{address: "thor1szv77gjy2ruvtuqnhd09nckx9kec3x5y96e50x", amount: cosmos.NewUint(218555070403)},
		{address: "thor1u5pfv07xtxz6aj59pnejaxh2dy7ew5s79ds8cw", amount: cosmos.NewUint(251196444191)},
		{address: "thor1u9dnzza6hpesrwq4p8j2f29v6jsyeq4le66j3c", amount: cosmos.NewUint(549696983774)},
		{address: "thor1ucwcatnqwjyucfrf7vv2xnmfzfaplvrkqzl337", amount: cosmos.NewUint(133626712498)},
		{address: "thor1v882x92avxcmkaucm2h6cp6j2qc4lhll0vfsex", amount: cosmos.NewUint(459913606302)},
		{address: "thor1v8shd72ns62j9g6za6yupmw2jlvf5ezczcc7vg", amount: cosmos.NewUint(118068825494)},
		{address: "thor1vp29289yyvfar0ektscjk08r0tufvl24tn6xf9", amount: cosmos.NewUint(432573328768)},
		{address: "thor1vt207wgvefjgk88mtfjuurcl3vw6z4d2gu5psw", amount: cosmos.NewUint(380057621427)},
		{address: "thor1vtcpemkcgr72jl8ael0jvt07dfc4pmqc2jj7vf", amount: cosmos.NewUint(398051284534)},
		{address: "thor1vwqz5hhh5un28qlz6x5f8zczj39jqwel38q2kc", amount: cosmos.NewUint(258025724038)},
		{address: "thor1w8mntay3xuk3c77j8fgvyyt0nfvl2sk398a3ww", amount: cosmos.NewUint(59098407417)},
		{address: "thor1wed8wsu98kvphxx39hgjnqtzw0s69u85ckuh40", amount: cosmos.NewUint(554560176371)},
		{address: "thor1wymexemzfhexhp0kars6fdl9mmteg20tcngm5m", amount: cosmos.NewUint(237970977507)},
		{address: "thor1x2whgc2nt665y0kc44uywhynazvp0l8tp0vtu6", amount: cosmos.NewUint(465158711923)},
		{address: "thor1xczxhtnu4vmtmkazny5vkdplfm7gtdfwcerdss", amount: cosmos.NewUint(117963931597)},
		{address: "thor1xd4j3gk9frpxh8r22runntnqy34lwzrdkazldh", amount: cosmos.NewUint(119946139794)},
		{address: "thor1xjm3wnxp0fl3a7vrhkagp4slnez0vspv4hfwle", amount: cosmos.NewUint(112248676846)},
		{address: "thor1xn9whf9nnhcr4mvtq0acjsu43n0x47nna3fyzq", amount: cosmos.NewUint(146444976424)},
		{address: "thor1yak0z56elhcfqw7xn7wjmp43ndnnxgfcmwkwex", amount: cosmos.NewUint(73476881807)},
		{address: "thor1ypjwdplx07vf42qdfkex39dp8zxqnaects270v", amount: cosmos.NewUint(82610353774)},
		{address: "thor1ytvzjwmf9pwuq95mdya4y9gale3864jz2ryu3r", amount: cosmos.NewUint(118287628643)},
		{address: "thor1yzwjdkujv956lx5j2r7a4fk4ajjjm5773v7v3h", amount: cosmos.NewUint(424889623867)},
		{address: "thor1z3dmy779shx8x9903ldnyqnt3a3g6vjqx68hkt", amount: cosmos.NewUint(18325429283)},
		{address: "thor1zfy2dm8urvwzc6shcmfpewdxamf8v35zq593ev", amount: cosmos.NewUint(92205991554)},
		{address: "thor1zga95gkv87356lmjj0mvw3geylfuv3ph7wa9t0", amount: cosmos.NewUint(119876405200)},
		{address: "thor1zkt4hzha6he8d0nlkyeer6pfudvkjckr6hj0p6", amount: cosmos.NewUint(88257910154)},
	}

	// sum amounts to get the total we will refund to nodes from the reserve
	total := cosmos.ZeroUint()
	for _, credit := range credits {
		total = total.Add(credit.amount)
	}

	// assertion for sanity check (~255k RUNE)
	if !total.Equal(cosmos.NewUint(25580168572803)) {
		ctx.Logger().Error("total refund amount is not correct", "total", total)
		return
	}

	// send coins from reserve to bond module
	if err := mgr.Keeper().SendFromModuleToModule(ctx, ReserveName, BondName, common.Coins{common.NewCoin(common.RuneNative, total)}); err != nil {
		ctx.Logger().Error("fail to transfer coin from reserve to bond module", "error", err)
		return
	}

	for _, credit := range credits {
		ctx.Logger().Info("credit", "node", credit.address, "amount", credit.amount)

		// get addresses
		addr, err := cosmos.AccAddressFromBech32(credit.address)
		if err != nil {
			ctx.Logger().Error("fail to parse node address", "error", err)
			return
		}

		// get node account
		na, err := mgr.Keeper().GetNodeAccount(ctx, addr)
		if err != nil {
			ctx.Logger().Error("fail to get node account", "error", err)
			return
		}

		// update node bond
		na.Bond = na.Bond.Add(credit.amount)

		// store updated records
		if err := mgr.Keeper().SetNodeAccount(ctx, na); err != nil {
			ctx.Logger().Error("fail to save node account", "error", err)
			return
		}
	}

	for _, item := range getInitClout() {
		addr, err := common.NewAddress(item.address)
		if err != nil {
			ctx.Logger().Error("failed to parse address during init clout", "address", item.address, "error", err)
			continue
		}
		c := NewSwapperClout(addr)
		c.Score = item.amount
		if err := mgr.Keeper().SetSwapperClout(ctx, c); err != nil {
			ctx.Logger().Error("failed to set swapper clout", "address", item.address, "score", item.amount, "error", err)
		}
	}
}

func migrateStoreV126(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v126", "error", err)
		}
	}()

	for _, item := range getInitCloutTHORNames() {
		addr, err := common.NewAddress(item.address)
		if err != nil {
			ctx.Logger().Error("failed to parse address during init clout", "address", item.address, "error", err)
			continue
		}
		c := NewSwapperClout(addr)
		c.Score = item.amount
		if err := mgr.Keeper().SetSwapperClout(ctx, c); err != nil {
			ctx.Logger().Error("failed to set swapper clout", "address", item.address, "score", item.amount, "error", err)
		}
	}
}

func migrateStoreV128(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v128", "error", err)
		}
	}()

	// burn 60m RUNE from standby reserve
	toBurn := common.NewCoin(common.RuneNative, cosmos.NewUint(60000000*1e8))
	standbyReserveAddress, err := cosmos.AccAddressFromBech32("thor1lj62pg6ryxv2htekqx04nv7wd3g98qf9gfvamy")
	if err != nil {
		ctx.Logger().Error("unable to AccAddressFromBech32 in v128 store migration", "error", err)
		return
	}

	// cannot burn directly from account. send to thorchain module, then burn
	err = mgr.Keeper().SendFromAccountToModule(ctx, standbyReserveAddress, ModuleName, common.Coins{toBurn})
	if err != nil {
		ctx.Logger().Error("unable to SendFromAccountToModule in v128 store migration", "error", err)
		return
	}

	// burn coins
	err = mgr.Keeper().BurnFromModule(ctx, ModuleName, toBurn)
	if err != nil {
		ctx.Logger().Error("unable to BurnFromModule in v128 store migration", "error", err)
		return
	}
	burnEvt := NewEventMintBurn(BurnSupplyType, toBurn.Asset.Native(), toBurn.Amount, "adr")
	if err := mgr.EventMgr().EmitEvent(ctx, burnEvt); err != nil {
		ctx.Logger().Error("fail to emit burn event in v128 store migration", "error", err)
	}
	ctx.Logger().Info("Burned 60m RUNE")
}

// The https://lends.so/ UI prompted user to sign an unsupported taproot tx to a
// thorchain vault. It is clearly documented here not to do that
// https://dev.thorchain.org/concepts/sending-transactions.html#utxo-chains. Despite the
// fact, Lends denies any blame whatsoever. The user inbound amount will be refunded to
// the treasury, which will keep 20% as a bounty and return the remainder to the user.
//
// Inbound: https://blockstream.info/tx/a1354d5d8ac67ae540cc61accbae64dbe2ed0c25a5334925ab28a02416581f44
func migrateStoreV129(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v129", "error", err)
		}
	}()

	treasuryAddr, err := common.NewAddress("bc1q5s9rxyu94n8twggw25agldy8xl4v55e76l02vn")
	if err != nil {
		ctx.Logger().Error("fail to create treasury addr", "error", err)
		return
	}
	vaultPubkey, err := common.NewPubKey("thorpub1addwnpepq0y6xyun0469ngddnsufuglvem6rkh0lwnm5dq3p6uhu690g953kgku4uhx")
	if err != nil {
		ctx.Logger().Error("fail to create pubkey for vault", "error", err)
		return
	}
	vaultAddr, err := common.NewAddress("bc1qrca9ta0x2znmqahnwwsfpkddjlg73sn79pcq02")
	if err != nil {
		ctx.Logger().Error("fail to create vault addr", "error", err)
		return
	}
	txID, err := common.NewTxID("A1354D5D8AC67AE540CC61ACCBAE64DBE2ED0C25A5334925AB28A02416581F44")
	if err != nil {
		ctx.Logger().Error("fail to create tx id", "error", err)
		return
	}

	badMemo := "BOUNTY" // trigger refund for bad inbound memo
	externalHeight := int64(834329)

	unobservedTxs := ObservedTxs{
		NewObservedTx(common.Tx{
			ID:          txID,
			Chain:       common.BTCChain,
			FromAddress: treasuryAddr, // this is faked to refund to treasury
			ToAddress:   vaultAddr,
			Coins: common.NewCoins(common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1_9980_0000),
			}),
			Gas: common.Gas{common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1),
			}},
			Memo: badMemo,
		}, externalHeight, vaultPubkey, externalHeight),
	}

	err = makeFakeTxInObservation(ctx, mgr, unobservedTxs)
	if err != nil {
		ctx.Logger().Error("failed to migrate v129", "error", err)
	}
}

func migrateStoreV131(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v131", "error", err)
		}
	}()

	// Helper function
	subBalance := func(ctx cosmos.Context, mgr *Mgrs, amountUint64 uint64, assetString, pkString string) {
		amount := cosmos.NewUint(amountUint64)

		asset, err := common.NewAsset(assetString)
		if err != nil {
			ctx.Logger().Error("fail to make asset", "assetString", assetString, "error", err)
			return
		}

		pubkey, err := common.NewPubKey(pkString)
		if err != nil {
			ctx.Logger().Error("fail to make pubkey", "pkString", pkString, "error", err)
			return
		}

		vault, err := mgr.Keeper().GetVault(ctx, pubkey)
		if err != nil {
			ctx.Logger().Error("fail to get vault", "pubkey", pubkey, "error", err)
			return
		}

		coins := common.NewCoins(common.NewCoin(asset, amount))
		vault.SubFunds(coins)

		if err := mgr.Keeper().SetVault(ctx, vault); err != nil {
			ctx.Logger().Error("fail to set vault", "vault", vault, "error", err)
			return
		}
	}

	// This migration address 6 separate issues.

	// Issue 1: https://gitlab.com/thorchain/thornode/-/issues/1898
	// Attempt to rescue 1 BTC from old vault, original tx unobserved by majority of nodes
	originalTxID := common.TxID("9F465AE9A43655619E0ADCC0EC52187B8856BF7F025E9142A424C272394CFA50")
	vaultPubKey := common.PubKey("thorpub1addwnpepqw68vqcyfyerpqvmfn7r39myxgn9q5hwd7crk82radl7ggl5dvmm2uqxzwk")
	externalHeight := int64(837277)
	unobservedTxs := ObservedTxs{
		NewObservedTx(common.Tx{
			ID:          originalTxID,
			Chain:       common.BTCChain,
			FromAddress: "bc1qneepjkjy2p3rzk9aryvygm39gv9kzk0tvzt7g7",
			ToAddress:   "bc1qgydyyxt7t3nq73vk98t0850t4kxhx8tmfw78js",
			Coins: common.NewCoins(common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(100000000),
			}),
			Gas: common.Gas{common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1),
			}},
			Memo: "badMemo",
		}, externalHeight, vaultPubKey, externalHeight),
	}
	if err := makeFakeTxInObservation(ctx, mgr, unobservedTxs); err != nil {
		ctx.Logger().Error("fail to handle issue 1", "err", err)
	} else {
		ctx.Logger().Info("successfully handled issue 1")
	}

	// Issue 2: https://gitlab.com/thorchain/thornode/-/issues/1880
	// The target asset (partial fill) sent OK, but THORChain did not reschedule
	// the source asset refund (USDT) and it dropped off.
	// The vault has since migrated and has zero allowance for USDT on the router contract.
	// The sum of vault asset amount ~= the sum of vault allowances on the router contract.
	// the pool asset amount is -300k than the vault amount. so THORChain knows it has 300k more USDT than the pool.
	// Therefore, for this transaction, balance was already deducted from pool, but remains in vault.
	// Methodology: create the outbound without VaultPubKey specified, allowing THORChain to choose a vault.
	// NOTE: per the issue above, the UI advanced the refund to the customer, so refund the wallet that sent that tx.
	// https://etherscan.io/tx/0x705b941fc7b5d619836b03b711bb5cbe8b1a1251cff8fe3a94a357a27ad52905
	originalTxID = "8C7A9E17BDDE0E1F90EB2251E65ED4633B9F44E8EF3D0BA04308DB2C4D5CA92D"
	usdt, err := common.NewAsset("ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7")
	if err != nil {
		ctx.Logger().Error("unable to create NewAsset in ETH.USDT rescue tx", "error", err)
	} else {
		toi := TxOutItem{
			Chain:     common.ETHChain,
			ToAddress: common.Address("0xc85fef7a1b039a9e080aadf80ff6f1536dada088"),
			Coin:      common.NewCoin(usdt, cosmos.NewUint(11225335474500)),
			Memo:      fmt.Sprintf("REFUND:%s", originalTxID),
			InHash:    originalTxID,
		}
		_, err := mgr.txOutStore.TryAddTxOutItem(ctx, mgr, toi, cosmos.NewUint(0))
		if err != nil {
			ctx.Logger().Error("fail to attempt ETH.USDT rescue tx", "error", err)
		} else {
			ctx.Logger().Info("successfully handled issue 2")
		}
	}

	// Issue 3: https://gitlab.com/thorchain/thornode/-/issues/1926
	// V129 fix for an edge case in vault accounting for UTXO chains.
	// Pool accounting is correct, but vaults think they have more asset (DOGE) than they do.
	// Remove those improperly credited balances, so that vault accounting more closely matches pool accounting.
	subBalance(ctx, mgr, 1160074_8155_8121, "DOGE.DOGE", "thorpub1addwnpepq0y6xyun0469ngddnsufuglvem6rkh0lwnm5dq3p6uhu690g953kgku4uhx")
	subBalance(ctx, mgr, 1004979_3428_7850, "DOGE.DOGE", "thorpub1addwnpepq0xtamtm6l35efh3f5wlt5fql7cx9j94fywtumz83vvnzagx46h76yk8sa3")
	subBalance(ctx, mgr, 414685_0000_0000, "DOGE.DOGE", "thorpub1addwnpepq23r8srfathem5jgu8szm9mrjylx9y9atjeawwz4kajqpg3y9vvy7xghtst")
	subBalance(ctx, mgr, 156000_0000_0000, "DOGE.DOGE", "thorpub1addwnpepqf3h9xa9qantpfjmvv37cuvzm7u0zy48qrh25tafctyzvyn2l8fr7e66mqq")
	subBalance(ctx, mgr, 124812_2354_6484, "DOGE.DOGE", "thorpub1addwnpepqdaln9pasrj33vzcupezwp60hdgk8v797shp75jxp23xfvtu5gyd546uwcy")
	subBalance(ctx, mgr, 500_0000_0000, "DOGE.DOGE", "thorpub1addwnpepq2v44c4392cwa80mt9enycql7m90ghukuw2x2u72t60e9knlqer7gdcprrr")
	ctx.Logger().Info("finished handling issue 3")

	// Issue 4: https://gitlab.com/thorchain/thornode/-/issues/1881
	// Savers add with affiliates had a bug where not only would the add liquidity fail,
	// but the user would only be partially refunded..
	// Return the customer the difference between what they got and what they were owed.

	// First refund (LTC)
	originalTxID = "8EC6D7B459136D708BF03FCB55C3BBA04F0F32E924354E7BBB3F40637F0E56AD"
	toi := TxOutItem{
		Chain:     common.LTCChain,
		ToAddress: common.Address("ltc1qe2jleyfezmahfcmlgls2flqht9mlu8aaal7xhm"),
		Coin:      common.NewCoin(common.LTCAsset, cosmos.NewUint(5978633839)),
		Memo:      fmt.Sprintf("REFUND:%s", originalTxID),
		InHash:    originalTxID,
	}
	_, err = mgr.txOutStore.TryAddTxOutItem(ctx, mgr, toi, cosmos.NewUint(0))
	if err != nil {
		ctx.Logger().Error("fail to attempt LTC rescue tx", "error", err)
	} else {
		ctx.Logger().Info("successfully handled issue 4.1")
	}

	// Second refund (ETH.USDC)
	originalTxID = "B42EBA71BCF8D7ACEB14BE738A25670327A9441FCC29D70754AAF1D68A793529"
	usdc, err := common.NewAsset("ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48")
	if err != nil {
		ctx.Logger().Error("unable to create NewAsset in ETH.USDC rescue tx", "error", err)
	} else {
		toi := TxOutItem{
			Chain:     common.ETHChain,
			ToAddress: common.Address("0x58fa7c9c34e1e54e9b66d586203ec074ce412501"),
			Coin:      common.NewCoin(usdc, cosmos.NewUint(4751145371600)),
			Memo:      fmt.Sprintf("REFUND:%s", originalTxID),
			InHash:    originalTxID,
		}
		_, err := mgr.txOutStore.TryAddTxOutItem(ctx, mgr, toi, cosmos.NewUint(0))
		if err != nil {
			ctx.Logger().Error("fail to attempt ETH.USDC rescue tx", "error", err)
		} else {
			ctx.Logger().Info("successfully handled issue 4.2")
		}
	}

	// Third refund (BNB)
	originalTxID = "374E024BB5162CF62F3A8F1C0DA7B0A52F2B4B395E154EECDC4F3EDB0C5E3D43"
	bnb, err := common.NewAsset("BSC.BNB")
	if err != nil {
		ctx.Logger().Error("fail to create NewAsset in BNB rescue tx", "error", err)
	} else {
		toi := TxOutItem{
			Chain:     common.BSCChain,
			ToAddress: common.Address("0x04c5998ded94f89263370444ce64a99b7dbc9f46"),
			Coin:      common.NewCoin(bnb, cosmos.NewUint(2691840026)),
			Memo:      fmt.Sprintf("REFUND:%s", originalTxID),
			InHash:    originalTxID,
		}
		_, err := mgr.txOutStore.TryAddTxOutItem(ctx, mgr, toi, cosmos.NewUint(0))
		if err != nil {
			ctx.Logger().Error("fail to attempt BNB rescue tx", "error", err)
		} else {
			ctx.Logger().Info("successfully handled issue 4.3")
		}
	}

	// Issue 5: https://gitlab.com/thorchain/thornode/-/issues/1935
	// BNB Beacon ragnarok resulted in ~8229 RUNE sent to RESERVE that should have been
	// left in the Pool Module. Bug fixed in separate PR.
	ragnarokPoolModuleInsolvency := cosmos.NewUint(822893701389)
	// send coins from reserve to pool module
	if err := mgr.Keeper().SendFromModuleToModule(ctx, ReserveName, AsgardName, common.Coins{common.NewCoin(common.RuneNative, ragnarokPoolModuleInsolvency)}); err != nil {
		ctx.Logger().Error("fail to transfer coin from reserve to pool module", "error", err)
	} else {
		ctx.Logger().Info("successfully handled issue 5")
	}

	// Issue 6: https://gitlab.com/thorchain/thornode/-/issues/1936
	// Following ragnarok, attempt to donate remaining tokens in BNB vaults to treasury
	// A subsequent store migration will drop any remaining vault balances from the KV store
	type TreasuryRecovery struct {
		Asset       string
		Amount      cosmos.Uint
		VaultPubKey string
	}
	recoveries := []TreasuryRecovery{}

	// bnb1rca9ta0x2znmqahnwwsfpkddjlg73sn73le82f
	// manually decremented for 0.2 BNB to ensure there is enough for gas
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BNB", Amount: cosmos.NewUint(35858610344), VaultPubKey: "thorpub1addwnpepq0y6xyun0469ngddnsufuglvem6rkh0lwnm5dq3p6uhu690g953kgku4uhx"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.AVA-645", Amount: cosmos.NewUint(16463772289), VaultPubKey: "thorpub1addwnpepq0y6xyun0469ngddnsufuglvem6rkh0lwnm5dq3p6uhu690g953kgku4uhx"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.TWT-8C2", Amount: cosmos.NewUint(669391046606), VaultPubKey: "thorpub1addwnpepq0y6xyun0469ngddnsufuglvem6rkh0lwnm5dq3p6uhu690g953kgku4uhx"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.ETH-1C9", Amount: cosmos.NewUint(415113), VaultPubKey: "thorpub1addwnpepq0y6xyun0469ngddnsufuglvem6rkh0lwnm5dq3p6uhu690g953kgku4uhx"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BUSD-BD1", Amount: cosmos.NewUint(23596032), VaultPubKey: "thorpub1addwnpepq0y6xyun0469ngddnsufuglvem6rkh0lwnm5dq3p6uhu690g953kgku4uhx"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BTCB-1DE", Amount: cosmos.NewUint(248757), VaultPubKey: "thorpub1addwnpepq0y6xyun0469ngddnsufuglvem6rkh0lwnm5dq3p6uhu690g953kgku4uhx"})

	// bnb1pg8zcjzrjjnknzkqh2eenyeaxj7qlax6rl9a0z
	// this vault needs a donation of BNB to pay for gas
	// recoveries = append(recoveries, TreasuryRecovery{VaultAddress: "bnb1pg8zcjzrjjnknzkqh2eenyeaxj7qlax6rl9a0z", Asset: "BNB.BNB", Amount: cosmos.NewUint(106412), VaultPubKey: "thorpub1addwnpepq23r8srfathem5jgu8szm9mrjylx9y9atjeawwz4kajqpg3y9vvy7xghtst"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BUSD-BD1", Amount: cosmos.NewUint(14317902698143), VaultPubKey: "thorpub1addwnpepq23r8srfathem5jgu8szm9mrjylx9y9atjeawwz4kajqpg3y9vvy7xghtst"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BTCB-1DE", Amount: cosmos.NewUint(21192302), VaultPubKey: "thorpub1addwnpepq23r8srfathem5jgu8szm9mrjylx9y9atjeawwz4kajqpg3y9vvy7xghtst"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.ETH-1C9", Amount: cosmos.NewUint(3453816695), VaultPubKey: "thorpub1addwnpepq23r8srfathem5jgu8szm9mrjylx9y9atjeawwz4kajqpg3y9vvy7xghtst"})

	// bnb1wy9467ppep0z8sdmzd5nc3xtt0ygu092l9a93u
	// this vault needs a donation of BNB to pay for gas
	// recoveries = append(recoveries, TreasuryRecovery{VaultAddress: "bnb1wy9467ppep0z8sdmzd5nc3xtt0ygu092l9a93u", Asset: "BNB.BNB", Amount: cosmos.NewUint(27689), VaultPubKey: "thorpub1addwnpepqf3h9xa9qantpfjmvv37cuvzm7u0zy48qrh25tafctyzvyn2l8fr7e66mqq"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BTCB-1DE", Amount: cosmos.NewUint(2507), VaultPubKey: "thorpub1addwnpepqf3h9xa9qantpfjmvv37cuvzm7u0zy48qrh25tafctyzvyn2l8fr7e66mqq"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.AVA-645", Amount: cosmos.NewUint(25227036220), VaultPubKey: "thorpub1addwnpepqf3h9xa9qantpfjmvv37cuvzm7u0zy48qrh25tafctyzvyn2l8fr7e66mqq"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.TWT-8C2", Amount: cosmos.NewUint(1867406976001), VaultPubKey: "thorpub1addwnpepqf3h9xa9qantpfjmvv37cuvzm7u0zy48qrh25tafctyzvyn2l8fr7e66mqq"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.ETH-1C9", Amount: cosmos.NewUint(119945), VaultPubKey: "thorpub1addwnpepqf3h9xa9qantpfjmvv37cuvzm7u0zy48qrh25tafctyzvyn2l8fr7e66mqq"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BUSD-BD1", Amount: cosmos.NewUint(1896139686), VaultPubKey: "thorpub1addwnpepqf3h9xa9qantpfjmvv37cuvzm7u0zy48qrh25tafctyzvyn2l8fr7e66mqq"})

	// bnb1hal0a2ywtd5qt97zfc2mmjwelxwlshgjn27xae
	// this vault needs a donation of BNB to pay for gas
	// recoveries = append(recoveries, TreasuryRecovery{VaultAddress: "bnb1hal0a2ywtd5qt97zfc2mmjwelxwlshgjn27xae", Asset: "BNB.BNB", Amount: cosmos.NewUint(150552), VaultPubKey: "thorpub1addwnpepqdaln9pasrj33vzcupezwp60hdgk8v797shp75jxp23xfvtu5gyd546uwcy"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.ETH-1C9", Amount: cosmos.NewUint(2301791), VaultPubKey: "thorpub1addwnpepqdaln9pasrj33vzcupezwp60hdgk8v797shp75jxp23xfvtu5gyd546uwcy"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.TWT-8C2", Amount: cosmos.NewUint(387443465935), VaultPubKey: "thorpub1addwnpepqdaln9pasrj33vzcupezwp60hdgk8v797shp75jxp23xfvtu5gyd546uwcy"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BTCB-1DE", Amount: cosmos.NewUint(329860), VaultPubKey: "thorpub1addwnpepqdaln9pasrj33vzcupezwp60hdgk8v797shp75jxp23xfvtu5gyd546uwcy"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BUSD-BD1", Amount: cosmos.NewUint(7542272440680), VaultPubKey: "thorpub1addwnpepqdaln9pasrj33vzcupezwp60hdgk8v797shp75jxp23xfvtu5gyd546uwcy"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.AVA-645", Amount: cosmos.NewUint(46237488), VaultPubKey: "thorpub1addwnpepqdaln9pasrj33vzcupezwp60hdgk8v797shp75jxp23xfvtu5gyd546uwcy"})

	// bnb1faprhnsvxmv586zwuyv3dugtegsa8yqz8nnmnh
	// manually decremented for 0.2 BNB to ensure there is enough for gas
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BNB", Amount: cosmos.NewUint(44312404328), VaultPubKey: "thorpub1addwnpepq0xtamtm6l35efh3f5wlt5fql7cx9j94fywtumz83vvnzagx46h76yk8sa3"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BUSD-BD1", Amount: cosmos.NewUint(23758202063840), VaultPubKey: "thorpub1addwnpepq0xtamtm6l35efh3f5wlt5fql7cx9j94fywtumz83vvnzagx46h76yk8sa3"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.ETH-1C9", Amount: cosmos.NewUint(612881), VaultPubKey: "thorpub1addwnpepq0xtamtm6l35efh3f5wlt5fql7cx9j94fywtumz83vvnzagx46h76yk8sa3"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BTCB-1DE", Amount: cosmos.NewUint(235740384), VaultPubKey: "thorpub1addwnpepq0xtamtm6l35efh3f5wlt5fql7cx9j94fywtumz83vvnzagx46h76yk8sa3"})

	// bnb1tdrfuq4my89fk39lq8u9zp3x8zh9qnk7s0vgqu
	// manually decremented for 0.2 BNB to ensure there is enough for gas
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BNB", Amount: cosmos.NewUint(166821351), VaultPubKey: "thorpub1addwnpepq2v44c4392cwa80mt9enycql7m90ghukuw2x2u72t60e9knlqer7gdcprrr"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.ETH-1C9", Amount: cosmos.NewUint(26936939), VaultPubKey: "thorpub1addwnpepq2v44c4392cwa80mt9enycql7m90ghukuw2x2u72t60e9knlqer7gdcprrr"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BUSD-BD1", Amount: cosmos.NewUint(13784540734936), VaultPubKey: "thorpub1addwnpepq2v44c4392cwa80mt9enycql7m90ghukuw2x2u72t60e9knlqer7gdcprrr"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.BTCB-1DE", Amount: cosmos.NewUint(720964), VaultPubKey: "thorpub1addwnpepq2v44c4392cwa80mt9enycql7m90ghukuw2x2u72t60e9knlqer7gdcprrr"})
	recoveries = append(recoveries, TreasuryRecovery{Asset: "BNB.TWT-8C2", Amount: cosmos.NewUint(1859478374319), VaultPubKey: "thorpub1addwnpepq2v44c4392cwa80mt9enycql7m90ghukuw2x2u72t60e9knlqer7gdcprrr"})

	gasRate := mgr.GasMgr().GetGasRate(ctx, common.BNBChain)
	maxGas, err := mgr.GasMgr().GetMaxGas(ctx, common.BNBChain)
	if err != nil {
		ctx.Logger().Error("fail to GetMaxGas for BNB recovery", "error", err)
		return
	}
	requeueTxBnb := func(ctx cosmos.Context, mgr *Mgrs, recovery TreasuryRecovery) {
		asset, err := common.NewAsset(recovery.Asset)
		if err != nil {
			ctx.Logger().Error("fail to create asset", "error", err, "recovery", recovery)
			return
		}
		recoveryTx := TxOutItem{
			Chain:       common.BNBChain,
			VaultPubKey: common.PubKey(recovery.VaultPubKey),
			ToAddress:   common.Address("bnb1pa6hpjs7qv0vkd5ks5tqa2xtt2gk5n08yw7v7f"),
			Coin:        common.NewCoin(asset, recovery.Amount),
			Memo:        "",
			InHash:      "",
			GasRate:     int64(gasRate.Uint64()),
			MaxGas:      common.Gas{maxGas},
		}

		txBytes, err := json.Marshal(recoveryTx)
		if err != nil {
			ctx.Logger().Error("fail to Marshal recoveryTx", "error", err, "recovery", recovery)
			return
		}

		hash := strings.ToUpper(hex.EncodeToString(tmhash.Sum(txBytes)))
		recoveryTx.Memo = fmt.Sprintf("REFUND:%s", hash)
		recoveryTx.InHash = common.TxID(hash)

		err = mgr.txOutStore.UnSafeAddTxOutItem(ctx, mgr, recoveryTx, ctx.BlockHeight())
		if err != nil {
			ctx.Logger().Error("fail to add BNB recovery tx", "error", err, "tx", recoveryTx)
		}
	}

	for _, recovery := range recoveries {
		requeueTxBnb(ctx, mgr, recovery)
	}
}

// migrateStoreV132 retries several migrations from v131:
// 1) BNB treasury recovery - retry substantial remaining balances for BNB, BUSD, and TWT.
// 2) Issue 2 and 4 from previous version - both failed on TryTxOutAddItem due to voter.OutboundHeight being set at a lower height, when the partial refund occurred.
// Prefer to use UnSafeAddTxOutItem w/ discoverOutbounds helper functions to skip other handling logic and prepare multiple outbounds for an explicit TxOutItem
func migrateStoreV132(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v132", "error", err)
		}
	}()

	// list of vault assets to eject
	busd, err := common.NewAsset("BNB.BUSD-BD1")
	if err != nil {
		ctx.Logger().Error("fail to create busd asset", "error", err)
		return
	}
	twt, err := common.NewAsset("BNB.TWT-8C2")
	if err != nil {
		ctx.Logger().Error("fail to create twt asset", "error", err)
		return
	}
	assets := []common.Asset{busd, twt, common.BNBAsset}

	// treasury bnb address
	treasuryBnbAddress, err := common.NewAddress("bnb1pa6hpjs7qv0vkd5ks5tqa2xtt2gk5n08yw7v7f")
	if err != nil {
		ctx.Logger().Error("fail to create treasury bnb address", "error", err)
		return
	}

	// get active vaults
	activeVaults, err := mgr.Keeper().GetAsgardVaultsByStatus(ctx, ActiveVault)
	if err != nil {
		ctx.Logger().Error("fail to get active vaults", "error", err)
		return
	}

	// get gas information for tx outs
	gasRate := mgr.GasMgr().GetGasRate(ctx, common.BNBChain)
	maxGas, err := mgr.GasMgr().GetMaxGas(ctx, common.BNBChain)
	if err != nil {
		ctx.Logger().Error("fail to get max gas", "error", err)
		return
	}

	// iterate all active vaults and send the full balance of eject assets to treasury
	for _, vault := range activeVaults {
		for _, asset := range assets {
			coin := vault.GetCoin(asset)

			// if this is BNB then deduct 0.1 BNB to leave for gas
			if asset.IsBNB() {
				coin.Amount = common.SafeSub(coin.Amount, cosmos.NewUint(1000_0000))
			}

			// skip in the event the vault has no balance of the asset
			if coin.Amount.IsZero() {
				continue
			}

			// send the full balance of the asset to treasury
			toi := TxOutItem{
				Chain:       common.BNBChain,
				VaultPubKey: vault.PubKey,
				ToAddress:   treasuryBnbAddress,
				Coin:        common.NewCoin(asset, coin.Amount),
				Memo:        "",
				InHash:      "",
				GasRate:     int64(gasRate.Uint64()),
				MaxGas:      common.Gas{maxGas},
			}

			// fake hash for outbound
			txBytes, err := json.Marshal(toi)
			if err != nil {
				ctx.Logger().Error("fail to marshal outbound", "error", err)
				continue
			}
			hash := strings.ToUpper(hex.EncodeToString(tmhash.Sum(txBytes)))
			toi.Memo = fmt.Sprintf("REFUND:%s", hash)
			toi.InHash = common.TxID(hash)

			// add the tx out item
			err = mgr.txOutStore.UnSafeAddTxOutItem(ctx, mgr, toi, ctx.BlockHeight())
			if err != nil {
				ctx.Logger().Error("fail to add recovery outbound", "error", err, "tx", toi)
			}
		}
	}

	// Requeue failed attempt Issue 1 from previous version migration
	// See L1316 from this file (manager_store_mainnet.go)
	originalTxID := "9F465AE9A43655619E0ADCC0EC52187B8856BF7F025E9142A424C272394CFA50"
	maxGas, err = mgr.gasMgr.GetMaxGas(ctx, common.BTCChain)
	if err != nil {
		ctx.Logger().Error("unable to GetMaxGas while retrying issue 1", "err", err)
	} else {
		gasRate = mgr.gasMgr.GetGasRate(ctx, common.BTCChain)
		droppedRescue := TxOutItem{
			Chain:       common.BTCChain,
			ToAddress:   common.Address("bc1qneepjkjy2p3rzk9aryvygm39gv9kzk0tvzt7g7"),
			VaultPubKey: common.PubKey("thorpub1addwnpepqw68vqcyfyerpqvmfn7r39myxgn9q5hwd7crk82radl7ggl5dvmm2uqxzwk"), // vaults/asgard?height=14766979
			Coin:        common.NewCoin(common.BTCAsset, cosmos.NewUint(99000000)),
			Memo:        fmt.Sprintf("REFUND:%s", originalTxID),
			InHash:      common.TxID(originalTxID),
			GasRate:     int64(gasRate.Uint64()),
			MaxGas:      common.Gas{maxGas},
		}

		err := mgr.txOutStore.UnSafeAddTxOutItem(ctx, mgr, droppedRescue, ctx.BlockHeight())
		if err != nil {
			ctx.Logger().Error("fail to retry BTC rescue tx", "error", err)
		}
	}

	// Retry failed attempt Issue 2 from previous version migration
	// See L1344 from this file (manager_store_mainnet.go)
	originalTxID = "8C7A9E17BDDE0E1F90EB2251E65ED4633B9F44E8EF3D0BA04308DB2C4D5CA92D"
	usdt, err := common.NewAsset("ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7")
	if err != nil {
		ctx.Logger().Error("unable to create NewAsset in Issue 2 retry", "error", err)
	} else {
		toi := TxOutItem{
			Chain:     common.ETHChain,
			ToAddress: common.Address("0xc85fef7a1b039a9e080aadf80ff6f1536dada088"),
			Coin:      common.NewCoin(usdt, cosmos.NewUint(11225335474500)),
			Memo:      fmt.Sprintf("REFUND:%s", originalTxID),
			InHash:    common.TxID(originalTxID),
		}
		outbounds, err := discoverOutbounds(ctx, mgr, toi)
		if err != nil {
			ctx.Logger().Error("unable to discoverOutbounds for Issue 2 retry", "error", err)
		}
		for _, outbound := range outbounds {
			if err := mgr.txOutStore.UnSafeAddTxOutItem(ctx, mgr, outbound, ctx.BlockHeight()); err != nil {
				ctx.Logger().Error("unable to UnSafeAddTxOutItem Issue 2 retry", "error", err, "outbound", outbound)
			}
		}
	}

	// Retry failed attempt(s) Issue 4 from previous version migration
	// See 1386 from this file (manager_store_mainnet.go)
	originalTxID = "8EC6D7B459136D708BF03FCB55C3BBA04F0F32E924354E7BBB3F40637F0E56AD"
	toi := TxOutItem{
		Chain:     common.LTCChain,
		ToAddress: common.Address("ltc1qe2jleyfezmahfcmlgls2flqht9mlu8aaal7xhm"),
		Coin:      common.NewCoin(common.LTCAsset, cosmos.NewUint(5978633839)),
		Memo:      fmt.Sprintf("REFUND:%s", originalTxID),
		InHash:    common.TxID(originalTxID),
	}
	outbounds, err := discoverOutbounds(ctx, mgr, toi)
	if err != nil {
		ctx.Logger().Error("unable to discoverOutbounds for Issue 4 retry", "error", err)
	}
	for _, outbound := range outbounds {
		if err := mgr.txOutStore.UnSafeAddTxOutItem(ctx, mgr, outbound, ctx.BlockHeight()); err != nil {
			ctx.Logger().Error("unable to UnSafeAddTxOutItem Issue 4 retry", "error", err, "outbound", outbound)
		}
	}

	// Second refund (ETH.USDC)
	originalTxID = "B42EBA71BCF8D7ACEB14BE738A25670327A9441FCC29D70754AAF1D68A793529"
	usdc, err := common.NewAsset("ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48")
	if err != nil {
		ctx.Logger().Error("unable to NewAsset in Issue 4 retry", "error", err)
	} else {
		toi := TxOutItem{
			Chain:     common.ETHChain,
			ToAddress: common.Address("0x58fa7c9c34e1e54e9b66d586203ec074ce412501"),
			Coin:      common.NewCoin(usdc, cosmos.NewUint(4751145371600)),
			Memo:      fmt.Sprintf("REFUND:%s", originalTxID),
			InHash:    common.TxID(originalTxID),
		}
		outbounds, err := discoverOutbounds(ctx, mgr, toi)
		if err != nil {
			ctx.Logger().Error("unable to discoverOutbounds for Issue 4 retry", "error", err)
		}
		for _, outbound := range outbounds {
			if err := mgr.txOutStore.UnSafeAddTxOutItem(ctx, mgr, outbound, ctx.BlockHeight()); err != nil {
				ctx.Logger().Error("unable to UnSafeAddTxOutItem Issue 4 retry", "error", err, "outbound", outbound)
			}
		}
	}

	// Third refund (BNB)
	originalTxID = "374E024BB5162CF62F3A8F1C0DA7B0A52F2B4B395E154EECDC4F3EDB0C5E3D43"
	bnb, err := common.NewAsset("BSC.BNB")
	if err != nil {
		ctx.Logger().Error("unable to NewAsset in Issue 4 retry", "error", err)
	} else {
		toi := TxOutItem{
			Chain:     common.BSCChain,
			ToAddress: common.Address("0x04c5998ded94f89263370444ce64a99b7dbc9f46"),
			Coin:      common.NewCoin(bnb, cosmos.NewUint(2691840026)),
			Memo:      fmt.Sprintf("REFUND:%s", originalTxID),
			InHash:    common.TxID(originalTxID),
		}
		outbounds, err := discoverOutbounds(ctx, mgr, toi)
		if err != nil {
			ctx.Logger().Error("unable to discoverOutbounds for Issue 4 retry", "error", err)
		}
		for _, outbound := range outbounds {
			if err := mgr.txOutStore.UnSafeAddTxOutItem(ctx, mgr, outbound, ctx.BlockHeight()); err != nil {
				ctx.Logger().Error("unable to UnSafeAddTxOutItem Issue 4 retry", "error", err, "outbound", outbound)
			}
		}
	}
}

func rescueTaprootBTC(ctx cosmos.Context, mgr *Mgrs, refundAddr, vaultAddr, vaultPubkey, txId string, externalBlockHeight, btcAmount int64) error {
	refundAddress, err := common.NewAddress(refundAddr)
	if err != nil {
		return fmt.Errorf("fail to create treasury addr: %w", err)
	}
	vaultPubKey, err := common.NewPubKey(vaultPubkey)
	if err != nil {
		return fmt.Errorf("fail to create pubkey for vault: %w", err)
	}
	vaultAddress, err := common.NewAddress(vaultAddr)
	if err != nil {
		return fmt.Errorf("fail to create vault addr: %w", err)
	}
	txID, err := common.NewTxID(txId)
	if err != nil {
		return fmt.Errorf("fail to create tx id: %w", err)
	}

	badMemo := "BOUNTY" // trigger refund for bad inbound memo
	externalHeight := externalBlockHeight

	unobservedTxs := ObservedTxs{
		NewObservedTx(common.Tx{
			ID:          txID,
			Chain:       common.BTCChain,
			FromAddress: refundAddress,
			ToAddress:   vaultAddress,
			Coins: common.NewCoins(common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(uint64(btcAmount)),
			}),
			Gas: common.Gas{common.Coin{
				Asset:  common.BTCAsset,
				Amount: cosmos.NewUint(1),
			}},
			Memo: badMemo,
		}, externalHeight, vaultPubKey, externalHeight),
	}

	return makeFakeTxInObservation(ctx, mgr, unobservedTxs)
}

func migrateStoreV133(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v133", "error", err)
		}
	}()

	// Rescue 10 BTC Taproot tx
	// https://mempool.space/tx/bb7bab03caa29fd7b6ef1c169f53bb86e47fb75c016d7503540f7420441ad60b
	err := rescueTaprootBTC(ctx, mgr, "bc1qkw8n70m5f5h5v736fepm0knx8lfdyqft4mxxzs", "bc1qnv6w9tqtsl6whh2tcdam9plqrggx63vuvna3gv", "thorpub1addwnpepq2vp5ydfzmpxt32fm3clwk5562mkju67amtgj9h9fanfhvqzp4w324a9yw8", "BB7BAB03CAA29FD7B6EF1C169F53BB86E47FB75C016D7503540F7420441AD60B", 844286, 10_00000000)
	if err != nil {
		ctx.Logger().Error("failed to rescue BTC in migrate v133", "error", err)
	}

	vaults, err := mgr.Keeper().GetAsgardVaults(ctx)
	if err != nil {
		ctx.Logger().Error("fail to get asgard vaults", "error", err)
		return
	}

	// Zero all BNB Asset Amounts (following Ragnarok, in preparation for BEP2 sunset).
	for _, vault := range vaults {
		for i := range vault.Coins {
			if vault.Coins[i].Asset.Chain.IsBNB() {
				vault.Coins[i].Amount = cosmos.ZeroUint()
			}
		}
		if err := mgr.Keeper().SetVault(ctx, vault); err != nil {
			ctx.Logger().Error("fail to save vault", "error", err)
		}
	}

	// https://gitlab.com/thorchain/thornode/-/issues/1545
	// Zero out L1 address and set RUNE address (owner) to TreasuryModule for all treasury LPs

	// Mint and send smallest amount possible to initialize module account
	oneRune := common.NewCoin(common.RuneNative, cosmos.NewUint(1))
	if err := mgr.Keeper().MintToModule(ctx, ModuleName, oneRune); err != nil {
		ctx.Logger().Error("fail to MintToModule", "error", err)
		return
	}
	if err := mgr.Keeper().SendFromModuleToModule(ctx, ModuleName, TreasuryName, common.Coins{oneRune}); err != nil {
		ctx.Logger().Error("fail to SendFromModuleToModule", "error", err)
		return
	}

	// Change ownership of existing treasury LPs to new module account
	treasuryAddr, err := mgr.Keeper().GetModuleAddress(TreasuryName)
	if err != nil {
		ctx.Logger().Error("fail to get treasury module address", "error", err)
		return
	}
	changeLPOwnership(ctx, mgr, common.Address("thor1egxvam70a86jafa8gcg3kqfmfax3s0m2g3m754"), treasuryAddr)
	changeLPOwnership(ctx, mgr, common.Address("thor1wfe7hsuvup27lx04p5al4zlcnx6elsnyft7dzm"), treasuryAddr)
}
