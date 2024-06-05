//go:build regtest
// +build regtest

package thorchain

import (
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

func migrateStoreV86(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV88(ctx cosmos.Context, mgr Manager) {}

func migrateStoreV102(ctx cosmos.Context, mgr Manager) {}

func migrateStoreV103(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV106(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV108(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV109(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV110(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV111(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV113(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV114(ctx cosmos.Context, mgr *Mgrs) {}

// migrateStoreV116 subset of mainnet migration
func migrateStoreV116(ctx cosmos.Context, mgr *Mgrs) {
	bondRuneOver := cosmos.NewUint(6936522592883)
	asgardRuneUnder := cosmos.NewUint(5082320319988)
	thorchainRuneOver := cosmos.NewUint(100000000)

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

		// transfer rune from thorchain to reserve to clear thorchain balances
		{
			ModuleName:     ModuleName,
			RuneRecipient:  ReserveName,
			RuneToTransfer: thorchainRuneOver,
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
			},
		},
	}

	processModuleBalanceActions(ctx, mgr.Keeper(), actions)
}

func migrateStoreV117(ctx cosmos.Context, mgr *Mgrs) {}

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
			coinsToMint = coinsToMint.Add(mintCoin)
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
			coinsToTransfer = coinsToTransfer.Add(mintCoin)
		}
	}

	if len(coinsToTransfer) > 0 {
		if err := mgr.Keeper().SendFromModuleToModule(ctx, ModuleName, AsgardName, coinsToTransfer); err != nil {
			ctx.Logger().Error("fail to move coins during swap", "error", err)
		}
	}
}

func migrateStoreV122(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV123(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV124(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV125(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV126(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV128(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV129(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV131(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV132(ctx cosmos.Context, mgr *Mgrs) {}

func migrateStoreV133(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v133", "error", err)
		}
	}()

	treasuryAddr, err := mgr.Keeper().GetModuleAddress(TreasuryName)
	if err != nil {
		ctx.Logger().Error("fail to get treasury module address", "error", err)
		return
	}

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

	changeLPOwnership(ctx, mgr, common.Address("tthor1uuds8pd92qnnq0udw0rpg0szpgcslc9p8lluej"), treasuryAddr)
}
