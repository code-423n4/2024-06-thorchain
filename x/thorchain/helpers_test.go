package thorchain

import (
	"fmt"
	"strings"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

type HelperSuite struct{}

var _ = Suite(&HelperSuite{})

type TestRefundBondKeeper struct {
	keeper.KVStoreDummy
	pool    Pool
	na      NodeAccount
	vaults  Vaults
	modules map[string]int64
	consts  constants.ConstantValues
}

func (k *TestRefundBondKeeper) GetConfigInt64(ctx cosmos.Context, key constants.ConstantName) int64 {
	return k.consts.GetInt64Value(key)
}

func (k *TestRefundBondKeeper) GetAsgardVaultsByStatus(_ cosmos.Context, _ VaultStatus) (Vaults, error) {
	return k.vaults, nil
}

func (k *TestRefundBondKeeper) VaultExists(_ cosmos.Context, pk common.PubKey) bool {
	return true
}

func (k *TestRefundBondKeeper) GetLeastSecure(ctx cosmos.Context, vaults Vaults, signingTransPeriod int64) Vault {
	return vaults[0]
}

func (k *TestRefundBondKeeper) GetPool(_ cosmos.Context, asset common.Asset) (Pool, error) {
	if k.pool.Asset.Equals(asset) {
		return k.pool, nil
	}
	return NewPool(), errKaboom
}

func (k *TestRefundBondKeeper) SetNodeAccount(_ cosmos.Context, na NodeAccount) error {
	k.na = na
	return nil
}

func (k *TestRefundBondKeeper) SetPool(_ cosmos.Context, p Pool) error {
	if k.pool.Asset.Equals(p.Asset) {
		k.pool = p
		return nil
	}
	return errKaboom
}

func (k *TestRefundBondKeeper) SetBondProviders(ctx cosmos.Context, _ BondProviders) error {
	return nil
}

func (k *TestRefundBondKeeper) GetBondProviders(ctx cosmos.Context, add cosmos.AccAddress) (BondProviders, error) {
	return BondProviders{}, nil
}

func (k *TestRefundBondKeeper) SendFromModuleToModule(_ cosmos.Context, from, to string, coins common.Coins) error {
	k.modules[from] -= int64(coins[0].Amount.Uint64())
	k.modules[to] += int64(coins[0].Amount.Uint64())
	return nil
}

func (s *HelperSuite) TestRefundBondHappyPath(c *C) {
	ctx, _ := setupKeeperForTest(c)
	na := GetRandomValidatorNode(NodeActive)
	na.Bond = cosmos.NewUint(12098 * common.One)
	pk := GetRandomPubKey()
	na.PubKeySet.Secp256k1 = pk
	keeper := &TestRefundBondKeeper{
		modules: make(map[string]int64),
		consts:  constants.GetConstantValues(GetCurrentVersion()),
	}
	na.Status = NodeStandby
	mgr := NewDummyMgrWithKeeper(keeper)
	tx := GetRandomTx()
	tx.FromAddress, _ = common.NewAddress(na.BondAddress.String())
	err := refundBond(ctx, tx, na.NodeAddress, cosmos.ZeroUint(), &na, mgr)
	c.Assert(err, IsNil)
	items, err := mgr.TxOutStore().GetOutboundItems(ctx)
	c.Assert(err, IsNil)
	c.Assert(items, HasLen, 0)
}

func (s *HelperSuite) TestRefundBondDisableRequestToLeaveNode(c *C) {
	ctx, _ := setupKeeperForTest(c)
	na := GetRandomValidatorNode(NodeActive)
	na.Bond = cosmos.NewUint(12098 * common.One)
	pk := GetRandomPubKey()
	na.PubKeySet.Secp256k1 = pk
	keeper := &TestRefundBondKeeper{
		modules: make(map[string]int64),
		consts:  constants.GetConstantValues(GetCurrentVersion()),
	}
	na.Status = NodeStandby
	na.RequestedToLeave = true
	mgr := NewDummyMgrWithKeeper(keeper)
	tx := GetRandomTx()
	err := refundBond(ctx, tx, na.NodeAddress, cosmos.ZeroUint(), &na, mgr)
	c.Assert(err, IsNil)
	items, err := mgr.TxOutStore().GetOutboundItems(ctx)
	c.Assert(err, IsNil)
	c.Assert(items, HasLen, 0)
	c.Assert(err, IsNil)
	c.Assert(keeper.na.Status == NodeDisabled, Equals, true)
}

func (s *HelperSuite) TestDollarsPerRune(c *C) {
	ctx, k := setupKeeperForTest(c)
	mgr := NewDummyMgrWithKeeper(k)
	mgr.Keeper().SetMimir(ctx, "TorAnchor-BNB-BUSD-BD1", 1) // enable BUSD pool as a TOR anchor
	busd, err := common.NewAsset("BNB.BUSD-BD1")
	c.Assert(err, IsNil)
	pool := NewPool()
	pool.Asset = busd
	pool.Status = PoolAvailable
	pool.BalanceRune = cosmos.NewUint(85515078103667)
	pool.BalanceAsset = cosmos.NewUint(709802235538353)
	pool.Decimals = 8
	c.Assert(k.SetPool(ctx, pool), IsNil)

	runeUSDPrice := telem(mgr.Keeper().DollarsPerRune(ctx))
	c.Assert(runeUSDPrice, Equals, float32(8.300317))

	// Now try with a second pool, identical depths.
	mgr.Keeper().SetMimir(ctx, "TorAnchor-ETH-USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48", 1) // enable USDC pool as a TOR anchor
	usdc, err := common.NewAsset("ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48")
	c.Assert(err, IsNil)
	pool = NewPool()
	pool.Asset = usdc
	pool.Status = PoolAvailable
	pool.BalanceRune = cosmos.NewUint(85515078103667)
	pool.BalanceAsset = cosmos.NewUint(709802235538353)
	pool.Decimals = 8
	c.Assert(k.SetPool(ctx, pool), IsNil)

	runeUSDPrice = telem(mgr.Keeper().DollarsPerRune(ctx))
	c.Assert(runeUSDPrice, Equals, float32(8.300317))
}

func (s *HelperSuite) TestTelem(c *C) {
	value := cosmos.NewUint(12047733)
	c.Assert(value.Uint64(), Equals, uint64(12047733))
	c.Assert(telem(value), Equals, float32(0.12047733))
}

type addGasFeesKeeperHelper struct {
	keeper.Keeper
	errGetNetwork bool
	errSetNetwork bool
	errGetPool    bool
	errSetPool    bool
}

func newAddGasFeesKeeperHelper(keeper keeper.Keeper) *addGasFeesKeeperHelper {
	return &addGasFeesKeeperHelper{
		Keeper: keeper,
	}
}

func (h *addGasFeesKeeperHelper) GetNetwork(ctx cosmos.Context) (Network, error) {
	if h.errGetNetwork {
		return Network{}, errKaboom
	}
	return h.Keeper.GetNetwork(ctx)
}

func (h *addGasFeesKeeperHelper) SetNetwork(ctx cosmos.Context, data Network) error {
	if h.errSetNetwork {
		return errKaboom
	}
	return h.Keeper.SetNetwork(ctx, data)
}

func (h *addGasFeesKeeperHelper) SetPool(ctx cosmos.Context, pool Pool) error {
	if h.errSetPool {
		return errKaboom
	}
	return h.Keeper.SetPool(ctx, pool)
}

func (h *addGasFeesKeeperHelper) GetPool(ctx cosmos.Context, asset common.Asset) (Pool, error) {
	if h.errGetPool {
		return Pool{}, errKaboom
	}
	return h.Keeper.GetPool(ctx, asset)
}

type addGasFeeTestHelper struct {
	ctx cosmos.Context
	na  NodeAccount
	mgr Manager
}

func newAddGasFeeTestHelper(c *C) addGasFeeTestHelper {
	ctx, mgr := setupManagerForTest(c)
	keeper := newAddGasFeesKeeperHelper(mgr.Keeper())
	mgr.K = keeper
	pool := NewPool()
	pool.Asset = common.BNBAsset
	pool.BalanceAsset = cosmos.NewUint(100 * common.One)
	pool.BalanceRune = cosmos.NewUint(100 * common.One)
	pool.Status = PoolAvailable
	c.Assert(mgr.Keeper().SetPool(ctx, pool), IsNil)

	poolBTC := NewPool()
	poolBTC.Asset = common.BTCAsset
	poolBTC.BalanceAsset = cosmos.NewUint(100 * common.One)
	poolBTC.BalanceRune = cosmos.NewUint(100 * common.One)
	poolBTC.Status = PoolAvailable
	c.Assert(mgr.Keeper().SetPool(ctx, poolBTC), IsNil)

	na := GetRandomValidatorNode(NodeActive)
	c.Assert(mgr.Keeper().SetNodeAccount(ctx, na), IsNil)
	vault := NewVault(ctx.BlockHeight(), ActiveVault, AsgardVault, na.PubKeySet.Secp256k1, common.Chains{common.BNBChain}.Strings(), []ChainContract{})
	// TODO:  Perhaps make this vault entirely unrelated to the NodeAccount pubkey, such as with an addGasFeeTestHelper 'vault' field.
	c.Assert(mgr.Keeper().SetVault(ctx, vault), IsNil)
	version := GetCurrentVersion()
	constAccessor := constants.GetConstantValues(version)
	mgr.gasMgr = newGasMgrVCUR(constAccessor, keeper)
	return addGasFeeTestHelper{
		ctx: ctx,
		mgr: mgr,
		na:  na,
	}
}

func (s *HelperSuite) TestAddGasFees(c *C) {
	testCases := []struct {
		name        string
		txCreator   func(helper addGasFeeTestHelper) ObservedTx
		runner      func(helper addGasFeeTestHelper, tx ObservedTx) error
		expectError bool
		validator   func(helper addGasFeeTestHelper, c *C)
	}{
		{
			name: "empty Gas should just return nil",
			txCreator: func(helper addGasFeeTestHelper) ObservedTx {
				return GetRandomObservedTx()
			},

			expectError: false,
		},
		{
			name: "normal BNB gas",
			txCreator: func(helper addGasFeeTestHelper) ObservedTx {
				tx := ObservedTx{
					Tx: common.Tx{
						ID:          GetRandomTxHash(),
						Chain:       common.BNBChain,
						FromAddress: GetRandomBNBAddress(),
						ToAddress:   GetRandomBNBAddress(),
						Coins: common.Coins{
							common.NewCoin(common.BNBAsset, cosmos.NewUint(5*common.One)),
							common.NewCoin(common.RuneAsset(), cosmos.NewUint(8*common.One)),
						},
						Gas: common.Gas{
							common.NewCoin(common.BNBAsset, BNBGasFeeSingleton[0].Amount),
						},
						Memo: "",
					},
					Status:         types.Status_done,
					OutHashes:      nil,
					BlockHeight:    helper.ctx.BlockHeight(),
					Signers:        []string{helper.na.NodeAddress.String()},
					ObservedPubKey: helper.na.PubKeySet.Secp256k1,
				}
				return tx
			},
			runner: func(helper addGasFeeTestHelper, tx ObservedTx) error {
				return addGasFees(helper.ctx, helper.mgr, tx)
			},
			expectError: false,
			validator: func(helper addGasFeeTestHelper, c *C) {
				expected := common.NewCoin(common.BNBAsset, BNBGasFeeSingleton[0].Amount)
				c.Assert(helper.mgr.GasMgr().GetGas(), HasLen, 1)
				c.Assert(helper.mgr.GasMgr().GetGas()[0].Equals(expected), Equals, true)
			},
		},
		{
			name: "normal BTC gas",
			txCreator: func(helper addGasFeeTestHelper) ObservedTx {
				tx := ObservedTx{
					Tx: common.Tx{
						ID:          GetRandomTxHash(),
						Chain:       common.BTCChain,
						FromAddress: GetRandomBTCAddress(),
						ToAddress:   GetRandomBTCAddress(),
						Coins: common.Coins{
							common.NewCoin(common.BTCAsset, cosmos.NewUint(5*common.One)),
						},
						Gas: common.Gas{
							common.NewCoin(common.BTCAsset, cosmos.NewUint(2000)),
						},
						Memo: "",
					},
					Status:         types.Status_done,
					OutHashes:      nil,
					BlockHeight:    helper.ctx.BlockHeight(),
					Signers:        []string{helper.na.NodeAddress.String()},
					ObservedPubKey: helper.na.PubKeySet.Secp256k1,
				}
				return tx
			},
			runner: func(helper addGasFeeTestHelper, tx ObservedTx) error {
				return addGasFees(helper.ctx, helper.mgr, tx)
			},
			expectError: false,
			validator: func(helper addGasFeeTestHelper, c *C) {
				expected := common.NewCoin(common.BTCAsset, cosmos.NewUint(2000))
				c.Assert(helper.mgr.GasMgr().GetGas(), HasLen, 1)
				c.Assert(helper.mgr.GasMgr().GetGas()[0].Equals(expected), Equals, true)
			},
		},
	}
	for _, tc := range testCases {
		helper := newAddGasFeeTestHelper(c)
		tx := tc.txCreator(helper)
		var err error
		if tc.runner == nil {
			err = addGasFees(helper.ctx, helper.mgr, tx)
		} else {
			err = tc.runner(helper, tx)
		}

		if err != nil && !tc.expectError {
			c.Errorf("test case: %s,didn't expect error however it got : %s", tc.name, err)
			c.FailNow()
		}
		if err == nil && tc.expectError {
			c.Errorf("test case: %s, expect error however it didn't", tc.name)
			c.FailNow()
		}
		if !tc.expectError && tc.validator != nil {
			tc.validator(helper, c)
			continue
		}
	}
}

func (s *HelperSuite) TestEmitPoolStageCostEvent(c *C) {
	ctx, mgr := setupManagerForTest(c)
	emitPoolBalanceChangedEvent(ctx,
		NewPoolMod(common.BTCAsset, cosmos.NewUint(1000), false, cosmos.ZeroUint(), false), "test", mgr)
	found := false
	for _, e := range ctx.EventManager().Events() {
		if strings.EqualFold(e.Type, types.PoolBalanceChangeEventType) {
			found = true
			break
		}
	}
	c.Assert(found, Equals, true)
}

func (s *HelperSuite) TestIsSynthMintPause(c *C) {
	ctx, mgr := setupManagerForTest(c)

	mgr.Keeper().SetMimir(ctx, constants.MaxSynthPerPoolDepth.String(), 1500)

	pool := types.Pool{
		Asset:        common.BTCAsset,
		BalanceAsset: cosmos.NewUint(100 * common.One),
		BalanceRune:  cosmos.NewUint(100 * common.One),
	}
	c.Assert(mgr.Keeper().SetPool(ctx, pool), IsNil)

	coins := cosmos.NewCoins(cosmos.NewCoin("btc/btc", cosmos.NewInt(29*common.One))) // 29% utilization
	c.Assert(mgr.coinKeeper.MintCoins(ctx, ModuleName, coins), IsNil)

	c.Assert(isSynthMintPaused(ctx, mgr, common.BTCAsset, cosmos.ZeroUint()), IsNil)

	// A swap that outputs 0.5 synth BTC would not surpass the synth utilization cap (29% -> 29.5%)
	c.Assert(isSynthMintPaused(ctx, mgr, common.BTCAsset, cosmos.NewUint(0.5*common.One)), IsNil)
	// A swap that outputs 1 synth BTC would not surpass the synth utilization cap (29% -> 30%)
	c.Assert(isSynthMintPaused(ctx, mgr, common.BTCAsset, cosmos.NewUint(1*common.One)), IsNil)
	// A swap that outputs 1.1 synth BTC would surpass the synth utilization cap (29% -> 30.1%)
	c.Assert(isSynthMintPaused(ctx, mgr, common.BTCAsset, cosmos.NewUint(1.1*common.One)), NotNil)

	coins = cosmos.NewCoins(cosmos.NewCoin("btc/btc", cosmos.NewInt(1*common.One))) // 30% utilization
	c.Assert(mgr.coinKeeper.MintCoins(ctx, ModuleName, coins), IsNil)

	c.Assert(isSynthMintPaused(ctx, mgr, common.BTCAsset, cosmos.ZeroUint()), IsNil)

	coins = cosmos.NewCoins(cosmos.NewCoin("btc/btc", cosmos.NewInt(1*common.One))) // 31% utilization
	c.Assert(mgr.coinKeeper.MintCoins(ctx, ModuleName, coins), IsNil)

	c.Assert(isSynthMintPaused(ctx, mgr, common.BTCAsset, cosmos.ZeroUint()), NotNil)
}

func (s *HelperSuite) TestUpdateTxOutGas(c *C) {
	ctx, mgr := setupManagerForTest(c)

	// Create ObservedVoter and add a TxOut
	txVoter := GetRandomObservedTxVoter()
	txOut := GetRandomTxOutItem()
	txVoter.Actions = append(txVoter.Actions, txOut)
	mgr.Keeper().SetObservedTxInVoter(ctx, txVoter)

	// Try to set new gas, should return error as TxOut InHash doesn't match
	newGas := common.Gas{common.NewCoin(common.LUNAAsset, cosmos.NewUint(2000000))}
	err := updateTxOutGas(ctx, mgr.K, txOut, newGas)
	c.Assert(err.Error(), Equals, fmt.Sprintf("fail to find tx out in ObservedTxVoter %s", txOut.InHash))

	// Update TxOut InHash to match, should update gas
	txOut.InHash = txVoter.TxID
	txVoter.Actions[1] = txOut
	mgr.Keeper().SetObservedTxInVoter(ctx, txVoter)

	// Err should be Nil
	err = updateTxOutGas(ctx, mgr.K, txOut, newGas)
	c.Assert(err, IsNil)

	// Keeper should have updated gas of TxOut in Actions
	txVoter, err = mgr.Keeper().GetObservedTxInVoter(ctx, txVoter.TxID)
	c.Assert(err, IsNil)

	didUpdate := false
	for _, item := range txVoter.Actions {
		if item.Equals(txOut) && item.MaxGas.Equals(newGas) {
			didUpdate = true
			break
		}
	}

	c.Assert(didUpdate, Equals, true)
}

func (s *HelperSuite) TestUpdateTxOutGasRate(c *C) {
	ctx, mgr := setupManagerForTest(c)

	// Create ObservedVoter and add a TxOut
	txVoter := GetRandomObservedTxVoter()
	txOut := GetRandomTxOutItem()
	txVoter.Actions = append(txVoter.Actions, txOut)
	mgr.Keeper().SetObservedTxInVoter(ctx, txVoter)

	// Try to set new gas rate, should return error as TxOut InHash doesn't match
	newGasRate := int64(25)
	err := updateTxOutGasRate(ctx, mgr.K, txOut, newGasRate)
	c.Assert(err.Error(), Equals, fmt.Sprintf("fail to find tx out in ObservedTxVoter %s", txOut.InHash))

	// Update TxOut InHash to match, should update gas
	txOut.InHash = txVoter.TxID
	txVoter.Actions[1] = txOut
	mgr.Keeper().SetObservedTxInVoter(ctx, txVoter)

	// Err should be Nil
	err = updateTxOutGasRate(ctx, mgr.K, txOut, newGasRate)
	c.Assert(err, IsNil)

	// Now that the actions have been updated (dependent on Equals which checks GasRate),
	// update the GasRate in the outbound queue item.
	txOut.GasRate = newGasRate

	// Keeper should have updated gas of TxOut in Actions
	txVoter, err = mgr.Keeper().GetObservedTxInVoter(ctx, txVoter.TxID)
	c.Assert(err, IsNil)

	didUpdate := false
	for _, item := range txVoter.Actions {
		if item.Equals(txOut) && item.GasRate == newGasRate {
			didUpdate = true
			break
		}
	}

	c.Assert(didUpdate, Equals, true)
}

func (s *HelperSuite) TestPOLPoolValue(c *C) {
	ctx, mgr := setupManagerForTest(c)

	polAddress, err := mgr.Keeper().GetModuleAddress(ReserveName)
	c.Assert(err, IsNil)

	btcPool := NewPool()
	btcPool.Asset = common.BTCAsset
	btcPool.BalanceRune = cosmos.NewUint(2000 * common.One)
	btcPool.BalanceAsset = cosmos.NewUint(20 * common.One)
	btcPool.LPUnits = cosmos.NewUint(1600)
	c.Assert(mgr.Keeper().SetPool(ctx, btcPool), IsNil)

	coin := common.NewCoin(common.BTCAsset.GetSyntheticAsset(), cosmos.NewUint(10*common.One))
	c.Assert(mgr.Keeper().MintToModule(ctx, ModuleName, coin), IsNil)

	lps := LiquidityProviders{
		{
			Asset:             btcPool.Asset,
			RuneAddress:       GetRandomBNBAddress(),
			AssetAddress:      GetRandomBTCAddress(),
			LastAddHeight:     5,
			Units:             btcPool.LPUnits.QuoUint64(2),
			PendingRune:       cosmos.ZeroUint(),
			PendingAsset:      cosmos.ZeroUint(),
			AssetDepositValue: cosmos.ZeroUint(),
			RuneDepositValue:  cosmos.ZeroUint(),
		},
		{
			Asset:             btcPool.Asset,
			RuneAddress:       polAddress,
			AssetAddress:      common.NoAddress,
			LastAddHeight:     10,
			Units:             btcPool.LPUnits.QuoUint64(2),
			PendingRune:       cosmos.ZeroUint(),
			PendingAsset:      cosmos.ZeroUint(),
			AssetDepositValue: cosmos.ZeroUint(),
			RuneDepositValue:  cosmos.ZeroUint(),
		},
	}
	for _, lp := range lps {
		mgr.Keeper().SetLiquidityProvider(ctx, lp)
	}

	value, err := polPoolValue(ctx, mgr)
	c.Assert(err, IsNil)
	c.Check(value.Uint64(), Equals, uint64(150023441162), Commentf("%d", value.Uint64()))
}

// This including the test of getTotalEffectiveBond.
func (s *HelperSuite) TestSecurityBond(c *C) {
	nas := make(NodeAccounts, 0)
	c.Assert(getEffectiveSecurityBond(nas).Uint64(), Equals, uint64(0), Commentf("%d", getEffectiveSecurityBond(nas).Uint64()))
	totalEffectiveBond, _ := getTotalEffectiveBond(nas)
	c.Assert(totalEffectiveBond.Uint64(), Equals, uint64(0), Commentf("%d", totalEffectiveBond.Uint64()))

	nas = NodeAccounts{
		NodeAccount{Bond: cosmos.NewUint(10)},
	}
	c.Assert(getEffectiveSecurityBond(nas).Uint64(), Equals, uint64(10), Commentf("%d", getEffectiveSecurityBond(nas).Uint64()))
	totalEffectiveBond, _ = getTotalEffectiveBond(nas)
	c.Assert(totalEffectiveBond.Uint64(), Equals, uint64(10), Commentf("%d", totalEffectiveBond.Uint64()))

	nas = NodeAccounts{
		NodeAccount{Bond: cosmos.NewUint(10)},
		NodeAccount{Bond: cosmos.NewUint(20)},
	}
	c.Assert(getEffectiveSecurityBond(nas).Uint64(), Equals, uint64(30), Commentf("%d", getEffectiveSecurityBond(nas).Uint64()))
	totalEffectiveBond, _ = getTotalEffectiveBond(nas)
	c.Assert(totalEffectiveBond.Uint64(), Equals, uint64(30), Commentf("%d", totalEffectiveBond.Uint64()))

	nas = NodeAccounts{
		NodeAccount{Bond: cosmos.NewUint(10)},
		NodeAccount{Bond: cosmos.NewUint(20)},
		NodeAccount{Bond: cosmos.NewUint(30)},
	}
	c.Assert(getEffectiveSecurityBond(nas).Uint64(), Equals, uint64(30), Commentf("%d", getEffectiveSecurityBond(nas).Uint64()))
	totalEffectiveBond, _ = getTotalEffectiveBond(nas)
	c.Assert(totalEffectiveBond.Uint64(), Equals, uint64(50), Commentf("%d", totalEffectiveBond.Uint64()))
	// Only 20 of the top-bond's node is effective.

	nas = NodeAccounts{
		NodeAccount{Bond: cosmos.NewUint(10)},
		NodeAccount{Bond: cosmos.NewUint(20)},
		NodeAccount{Bond: cosmos.NewUint(30)},
		NodeAccount{Bond: cosmos.NewUint(40)},
	}
	c.Assert(getEffectiveSecurityBond(nas).Uint64(), Equals, uint64(60), Commentf("%d", getEffectiveSecurityBond(nas).Uint64()))
	totalEffectiveBond, _ = getTotalEffectiveBond(nas)
	c.Assert(totalEffectiveBond.Uint64(), Equals, uint64(90), Commentf("%d", totalEffectiveBond.Uint64()))
	// Only 30 of the top-bond's node is effective.

	nas = NodeAccounts{
		NodeAccount{Bond: cosmos.NewUint(10)},
		NodeAccount{Bond: cosmos.NewUint(20)},
		NodeAccount{Bond: cosmos.NewUint(30)},
		NodeAccount{Bond: cosmos.NewUint(40)},
		NodeAccount{Bond: cosmos.NewUint(50)},
	}
	c.Assert(getEffectiveSecurityBond(nas).Uint64(), Equals, uint64(100), Commentf("%d", getEffectiveSecurityBond(nas).Uint64()))
	totalEffectiveBond, _ = getTotalEffectiveBond(nas)
	c.Assert(totalEffectiveBond.Uint64(), Equals, uint64(140), Commentf("%d", totalEffectiveBond.Uint64()))
	// Only 40 of the top-bond's node is effective.

	nas = NodeAccounts{
		NodeAccount{Bond: cosmos.NewUint(10)},
		NodeAccount{Bond: cosmos.NewUint(20)},
		NodeAccount{Bond: cosmos.NewUint(30)},
		NodeAccount{Bond: cosmos.NewUint(40)},
		NodeAccount{Bond: cosmos.NewUint(50)},
		NodeAccount{Bond: cosmos.NewUint(60)},
	}
	c.Assert(getEffectiveSecurityBond(nas).Uint64(), Equals, uint64(100), Commentf("%d", getEffectiveSecurityBond(nas).Uint64()))
	totalEffectiveBond, _ = getTotalEffectiveBond(nas)
	c.Assert(totalEffectiveBond.Uint64(), Equals, uint64(180), Commentf("%d", totalEffectiveBond.Uint64()))
	// Only 40 each of the top-bonds two nodes is effective.
}

func (s *HelperSuite) TestGetHardBondCap(c *C) {
	nas := make(NodeAccounts, 0)
	c.Assert(getHardBondCap(nas).Uint64(), Equals, uint64(0), Commentf("%d", getHardBondCap(nas).Uint64()))

	nas = NodeAccounts{
		NodeAccount{Bond: cosmos.NewUint(10)},
	}
	c.Assert(getHardBondCap(nas).Uint64(), Equals, uint64(10), Commentf("%d", getHardBondCap(nas).Uint64()))

	nas = NodeAccounts{
		NodeAccount{Bond: cosmos.NewUint(10)},
		NodeAccount{Bond: cosmos.NewUint(20)},
	}
	c.Assert(getHardBondCap(nas).Uint64(), Equals, uint64(20), Commentf("%d", getHardBondCap(nas).Uint64()))

	nas = NodeAccounts{
		NodeAccount{Bond: cosmos.NewUint(10)},
		NodeAccount{Bond: cosmos.NewUint(20)},
		NodeAccount{Bond: cosmos.NewUint(30)},
	}
	c.Assert(getHardBondCap(nas).Uint64(), Equals, uint64(20), Commentf("%d", getHardBondCap(nas).Uint64()))

	nas = NodeAccounts{
		NodeAccount{Bond: cosmos.NewUint(10)},
		NodeAccount{Bond: cosmos.NewUint(20)},
		NodeAccount{Bond: cosmos.NewUint(30)},
		NodeAccount{Bond: cosmos.NewUint(40)},
	}
	c.Assert(getHardBondCap(nas).Uint64(), Equals, uint64(30), Commentf("%d", getHardBondCap(nas).Uint64()))

	nas = NodeAccounts{
		NodeAccount{Bond: cosmos.NewUint(10)},
		NodeAccount{Bond: cosmos.NewUint(20)},
		NodeAccount{Bond: cosmos.NewUint(30)},
		NodeAccount{Bond: cosmos.NewUint(40)},
		NodeAccount{Bond: cosmos.NewUint(50)},
	}
	c.Assert(getHardBondCap(nas).Uint64(), Equals, uint64(40), Commentf("%d", getHardBondCap(nas).Uint64()))

	nas = NodeAccounts{
		NodeAccount{Bond: cosmos.NewUint(10)},
		NodeAccount{Bond: cosmos.NewUint(20)},
		NodeAccount{Bond: cosmos.NewUint(30)},
		NodeAccount{Bond: cosmos.NewUint(40)},
		NodeAccount{Bond: cosmos.NewUint(50)},
		NodeAccount{Bond: cosmos.NewUint(60)},
	}
	c.Assert(getHardBondCap(nas).Uint64(), Equals, uint64(40), Commentf("%d", getHardBondCap(nas).Uint64()))
}

func (HandlerSuite) TestIsSignedByActiveNodeAccounts(c *C) {
	ctx, mgr := setupManagerForTest(c)

	r := isSignedByActiveNodeAccounts(ctx, mgr.Keeper(), []cosmos.AccAddress{})
	c.Check(r, Equals, false,
		Commentf("empty signers should return false"))

	nodeAddr := GetRandomBech32Addr()
	r = isSignedByActiveNodeAccounts(ctx, mgr.Keeper(), []cosmos.AccAddress{nodeAddr})
	c.Check(r, Equals, false,
		Commentf("empty node account should return false"))

	nodeAccount1 := GetRandomValidatorNode(NodeWhiteListed)
	c.Assert(mgr.Keeper().SetNodeAccount(ctx, nodeAccount1), IsNil)
	r = isSignedByActiveNodeAccounts(ctx, mgr.Keeper(), []cosmos.AccAddress{nodeAccount1.NodeAddress})
	c.Check(r, Equals, false,
		Commentf("non-active node account should return false"))

	nodeAccount1.Status = NodeActive
	c.Assert(mgr.Keeper().SetNodeAccount(ctx, nodeAccount1), IsNil)
	r = isSignedByActiveNodeAccounts(ctx, mgr.Keeper(), []cosmos.AccAddress{nodeAccount1.NodeAddress})
	c.Check(r, Equals, true,
		Commentf("active node account should return true"))

	r = isSignedByActiveNodeAccounts(ctx, mgr.Keeper(), []cosmos.AccAddress{nodeAccount1.NodeAddress, nodeAddr})
	c.Check(r, Equals, false,
		Commentf("should return false if any signer is not an active validator"))

	nodeAccount1.Type = NodeTypeVault
	c.Assert(mgr.Keeper().SetNodeAccount(ctx, nodeAccount1), IsNil)
	r = isSignedByActiveNodeAccounts(ctx, mgr.Keeper(), []cosmos.AccAddress{nodeAccount1.NodeAddress})
	c.Check(r, Equals, false,
		Commentf("non-validator node should return false"))

	asgardAddr := mgr.Keeper().GetModuleAccAddress(AsgardName)
	r = isSignedByActiveNodeAccounts(ctx, mgr.Keeper(), []cosmos.AccAddress{asgardAddr})
	c.Check(r, Equals, true,
		Commentf("asgard module address should return true"))
}

func (HandlerSuite) TestWillSwapSucceed(c *C) {
	ctx, mgr := setupManagerForTest(c)

	// Set up some pools
	pool := NewPool()
	pool.Asset = common.BTCAsset
	pool.Status = PoolAvailable
	pool.BalanceRune = cosmos.NewUint(100_000 * common.One)
	pool.BalanceAsset = cosmos.NewUint(100 * common.One)
	pool.Decimals = 8
	c.Assert(mgr.Keeper().SetPool(ctx, pool), IsNil)

	pool2 := NewPool()
	pool2.Asset = common.ETHAsset
	pool2.Status = PoolAvailable
	pool2.BalanceRune = cosmos.NewUint(100_000 * common.One)
	pool2.BalanceAsset = cosmos.NewUint(1000 * common.One)
	pool2.Decimals = 8
	c.Assert(mgr.Keeper().SetPool(ctx, pool2), IsNil)

	// Set Network fees
	networkFee := NewNetworkFee(common.ETHChain, 1, 1000)
	c.Assert(mgr.Keeper().SaveNetworkFee(ctx, common.ETHChain, networkFee), IsNil)

	networkFee = NewNetworkFee(common.BTCChain, 1000, 10)
	c.Assert(mgr.Keeper().SaveNetworkFee(ctx, common.BTCChain, networkFee), IsNil)

	tx := common.NewTx(
		GetRandomTxHash(),
		GetRandomBTCAddress(),
		GetRandomBTCAddress(),
		common.Coins{common.NewCoin(common.BTCAsset, cosmos.NewUint(common.One))},
		common.Gas{
			{Asset: common.BTCAsset, Amount: cosmos.NewUint(37500)},
		},
		"",
	)

	// swap from BTC to ETH

	// no limit, should succeed
	msg := NewMsgSwap(tx, common.ETHAsset, GetRandomBTCAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(), "", "", nil, MarketOrder, 0, 0, GetRandomBech32Addr())
	c.Assert(willSwapOutputExceedLimitAndFees(ctx, mgr, *msg), Equals, true)

	// no limit, but small swap, should fail
	tx.Coins = common.Coins{common.NewCoin(common.BTCAsset, cosmos.NewUint(1))}
	msg = NewMsgSwap(tx, common.ETHAsset, GetRandomBTCAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(), "", "", nil, MarketOrder, 0, 0, GetRandomBech32Addr())
	c.Assert(willSwapOutputExceedLimitAndFees(ctx, mgr, *msg), Equals, false)

	// limit too high, should fail
	tx.Coins = common.Coins{common.NewCoin(common.BTCAsset, cosmos.NewUint(common.One))}
	msg = NewMsgSwap(tx, common.ETHAsset, GetRandomBTCAddress(), cosmos.NewUint(100*common.One), common.NoAddress, cosmos.ZeroUint(), "", "", nil, MarketOrder, 0, 0, GetRandomBech32Addr())
	c.Assert(willSwapOutputExceedLimitAndFees(ctx, mgr, *msg), Equals, false)

	// limit not too high, should succeed
	tx.Coins = common.Coins{common.NewCoin(common.BTCAsset, cosmos.NewUint(common.One))}
	msg = NewMsgSwap(tx, common.ETHAsset, GetRandomBTCAddress(), cosmos.NewUint(1*common.One), common.NoAddress, cosmos.ZeroUint(), "", "", nil, MarketOrder, 0, 0, GetRandomBech32Addr())
	c.Assert(willSwapOutputExceedLimitAndFees(ctx, mgr, *msg), Equals, true)

	runeTx := common.NewTx(
		GetRandomTxHash(),
		GetRandomTHORAddress(),
		GetRandomTHORAddress(),
		common.Coins{common.NewCoin(common.RuneNative, cosmos.NewUint(common.One*50))},
		common.Gas{
			{Asset: common.RuneNative, Amount: cosmos.NewUint(20000)},
		},
		"",
	)

	// swaps from RUNE

	// swap from RUNE no limit, should succeed
	msg = NewMsgSwap(runeTx, common.BTCAsset, GetRandomBTCAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(), "", "", nil, MarketOrder, 0, 0, GetRandomBech32Addr())
	c.Assert(willSwapOutputExceedLimitAndFees(ctx, mgr, *msg), Equals, true)

	// swap from RUNE, no limit, but small swap, should fail
	runeTx.Coins = common.Coins{common.NewCoin(common.RuneNative, cosmos.NewUint(1))}
	msg = NewMsgSwap(runeTx, common.BTCAsset, GetRandomBTCAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(), "", "", nil, MarketOrder, 0, 0, GetRandomBech32Addr())
	c.Assert(willSwapOutputExceedLimitAndFees(ctx, mgr, *msg), Equals, false)

	// swap from RUNE, limit too high, should fail
	runeTx.Coins = common.Coins{common.NewCoin(common.RuneNative, cosmos.NewUint(common.One*50))}
	msg = NewMsgSwap(runeTx, common.BTCAsset, GetRandomBTCAddress(), cosmos.NewUint(100*common.One), common.NoAddress, cosmos.ZeroUint(), "", "", nil, MarketOrder, 0, 0, GetRandomBech32Addr())
	c.Assert(willSwapOutputExceedLimitAndFees(ctx, mgr, *msg), Equals, false)

	// swap from RUNE, limit not too high, should succeed
	msg = NewMsgSwap(runeTx, common.BTCAsset, GetRandomBTCAddress(), cosmos.NewUint(0.01*common.One), common.NoAddress, cosmos.ZeroUint(), "", "", nil, MarketOrder, 0, 0, GetRandomBech32Addr())
	c.Assert(willSwapOutputExceedLimitAndFees(ctx, mgr, *msg), Equals, true)

	// swaps to RUNE

	// swap to RUNE, no limit, should succeed
	msg = NewMsgSwap(tx, common.RuneNative, GetRandomTHORAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(), "", "", nil, MarketOrder, 0, 0, GetRandomBech32Addr())
	c.Assert(willSwapOutputExceedLimitAndFees(ctx, mgr, *msg), Equals, true)

	// swap to RUNE, no limit, but small swap, should fail
	tx.Coins = common.Coins{common.NewCoin(common.BTCAsset, cosmos.NewUint(1))}
	msg = NewMsgSwap(tx, common.RuneNative, GetRandomTHORAddress(), cosmos.ZeroUint(), common.NoAddress, cosmos.ZeroUint(), "", "", nil, MarketOrder, 0, 0, GetRandomBech32Addr())
	c.Assert(willSwapOutputExceedLimitAndFees(ctx, mgr, *msg), Equals, false)

	// swap to RUNE, limit too high, should fail
	tx.Coins = common.Coins{common.NewCoin(common.BTCAsset, cosmos.NewUint(common.One))}
	msg = NewMsgSwap(tx, common.RuneNative, GetRandomTHORAddress(), cosmos.NewUint(100_000*common.One), common.NoAddress, cosmos.ZeroUint(), "", "", nil, MarketOrder, 0, 0, GetRandomBech32Addr())
	c.Assert(willSwapOutputExceedLimitAndFees(ctx, mgr, *msg), Equals, false)

	// swap to RUNE, limit not too high, should succeed
	tx.Coins = common.Coins{common.NewCoin(common.BTCAsset, cosmos.NewUint(common.One))}
	msg = NewMsgSwap(tx, common.RuneNative, GetRandomTHORAddress(), cosmos.NewUint(1*common.One), common.NoAddress, cosmos.ZeroUint(), "", "", nil, MarketOrder, 0, 0, GetRandomBech32Addr())
	c.Assert(willSwapOutputExceedLimitAndFees(ctx, mgr, *msg), Equals, true)
}

func (HandlerSuite) TestNewSwapMemo(c *C) {
	ctx, mgr := setupManagerForTest(c)
	addr := GetRandomBTCAddress()
	memo := NewSwapMemo(ctx, mgr, common.BTCAsset, addr, cosmos.ZeroUint(), "test", cosmos.ZeroUint())
	c.Assert(memo, Equals, fmt.Sprintf("=:BTC.BTC:%s:0:test:0", addr.String()))

	memo = NewSwapMemo(ctx, mgr, common.BTCAsset, addr, cosmos.NewUint(100), "test", cosmos.NewUint(50))
	c.Assert(memo, Equals, fmt.Sprintf("=:BTC.BTC:%s:100:test:50", addr.String()))

	addr = GetRandomTHORAddress()
	memo = NewSwapMemo(ctx, mgr, common.RuneNative, addr, cosmos.NewUint(0), "", cosmos.NewUint(0))
	c.Assert(memo, Equals, fmt.Sprintf("=:THOR.RUNE:%s:0::0", addr.String()))
}
