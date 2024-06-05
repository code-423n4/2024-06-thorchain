package thorchain

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

// SlasherV125 is VCUR implementation of slasher
type SlasherV125 struct {
	keeper   keeper.Keeper
	eventMgr EventManager
}

// newSlasherV125 create a new instance of Slasher
func newSlasherV125(keeper keeper.Keeper, eventMgr EventManager) *SlasherV125 {
	return &SlasherV125{keeper: keeper, eventMgr: eventMgr}
}

// BeginBlock called when a new block get proposed to detect whether there are duplicate vote
func (s *SlasherV125) BeginBlock(ctx cosmos.Context, req abci.RequestBeginBlock, constAccessor constants.ConstantValues) {
	var doubleSignEvidence []abci.Evidence
	// Iterate through any newly discovered evidence of infraction
	// Slash any validators (and since-unbonded liquidity within the unbonding period)
	// who contributed to valid infractions
	for _, evidence := range req.ByzantineValidators {
		// TODO: Remove on next hard fork.
		// The consensus failure occurred at block 7971846 and we give a few block buffer.
		if evidence.Height > 7971840 && evidence.Height < 7971850 {
			continue
		}
		switch evidence.Type {
		case abci.EvidenceType_DUPLICATE_VOTE:
			doubleSignEvidence = append(doubleSignEvidence, evidence)
		default:
			ctx.Logger().Error("ignored unknown evidence type", "type", evidence.Type)
		}
	}

	// Identify validators which didn't sign the previous block
	var missingSignAddresses []crypto.Address
	for _, voteInfo := range req.LastCommitInfo.Votes {
		if voteInfo.SignedLastBlock {
			continue
		}
		missingSignAddresses = append(missingSignAddresses, voteInfo.Validator.Address)
	}

	// Do not continue if there is no action to take.
	if len(doubleSignEvidence)+len(missingSignAddresses) == 0 {
		return
	}

	// Derive Active node validator addresses once.
	nas, err := s.keeper.ListActiveValidators(ctx)
	if err != nil {
		ctx.Logger().Error("fail to list active validators", "error", err)
		return
	}
	var validatorAddresses []nodeAddressValidatorAddressPair
	for _, na := range nas {
		pk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeConsPub, na.ValidatorConsPubKey)
		if err != nil {
			ctx.Logger().Error("fail to derive validator address", "error", err)
			continue
		}
		var pair nodeAddressValidatorAddressPair
		pair.nodeAddress = na.NodeAddress
		pair.validatorAddress = pk.Address()
		validatorAddresses = append(validatorAddresses, pair)
	}

	// Act on double signs.
	for _, evidence := range doubleSignEvidence {
		if err := s.HandleDoubleSign(ctx, evidence.Validator.Address, evidence.Height, constAccessor, validatorAddresses); err != nil {
			ctx.Logger().Error("fail to slash for double signing a block", "error", err)
		}
	}

	// Act on missing signs.
	for _, missingSignAddress := range missingSignAddresses {
		if err := s.HandleMissingSign(ctx, missingSignAddress, constAccessor, validatorAddresses); err != nil {
			ctx.Logger().Error("fail to slash for missing signing a block", "error", err)
		}
	}
}

// HandleDoubleSign - slashes a validator for signing two blocks at the same
// block height
// https://blog.cosmos.network/consensus-compare-casper-vs-tendermint-6df154ad56ae
func (s *SlasherV125) HandleDoubleSign(ctx cosmos.Context, addr crypto.Address, infractionHeight int64, constAccessor constants.ConstantValues, validatorAddresses []nodeAddressValidatorAddressPair) error {
	// check if we're recent enough to slash for this behavior
	maxAge := constAccessor.GetInt64Value(constants.DoubleSignMaxAge)
	if (ctx.BlockHeight() - infractionHeight) > maxAge {
		ctx.Logger().Info("double sign detected but too old to be slashed", "infraction height", fmt.Sprintf("%d", infractionHeight), "address", addr.String())
		return nil
	}

	for _, pair := range validatorAddresses {
		if addr.String() != pair.validatorAddress.String() {
			continue
		}

		na, err := s.keeper.GetNodeAccount(ctx, pair.nodeAddress)
		if err != nil {
			return err
		}

		if na.Bond.IsZero() {
			return fmt.Errorf("found account to slash for double signing, but did not have any bond to slash: %s", addr)
		}
		// take 5% of the minimum bond, and put it into the reserve
		minBond, err := s.keeper.GetMimir(ctx, constants.MinimumBondInRune.String())
		if minBond < 0 || err != nil {
			minBond = constAccessor.GetInt64Value(constants.MinimumBondInRune)
		}
		slashAmount := cosmos.NewUint(uint64(minBond)).MulUint64(5).QuoUint64(100)
		if slashAmount.GT(na.Bond) {
			slashAmount = na.Bond
		}

		slashFloat, _ := new(big.Float).SetInt(slashAmount.BigInt()).Float32()
		telemetry.IncrCounterWithLabels(
			[]string{"thornode", "bond_slash"},
			slashFloat,
			[]metrics.Label{
				telemetry.NewLabel("address", addr.String()),
				telemetry.NewLabel("reason", "double_block_sign"),
			},
		)

		na.Bond = common.SafeSub(na.Bond, slashAmount)
		coin := common.NewCoin(common.RuneNative, slashAmount)
		if err := s.keeper.SendFromModuleToModule(ctx, BondName, ReserveName, common.NewCoins(coin)); err != nil {
			ctx.Logger().Error("fail to transfer funds from bond to reserve", "error", err)
			return fmt.Errorf("fail to transfer funds from bond to reserve: %w", err)
		}

		return s.keeper.SetNodeAccount(ctx, na)
	}

	return fmt.Errorf("could not find active node account with validator address: %s", addr)
}

// HandleMissingSign - slashes a validator for not signing a block
func (s *SlasherV125) HandleMissingSign(ctx cosmos.Context, addr crypto.Address, constAccessor constants.ConstantValues, validatorAddresses []nodeAddressValidatorAddressPair) error {
	missBlockSignSlashPoints := s.keeper.GetConfigInt64(ctx, constants.MissBlockSignSlashPoints)

	for _, pair := range validatorAddresses {
		if addr.String() != pair.validatorAddress.String() {
			continue
		}

		na, err := s.keeper.GetNodeAccount(ctx, pair.nodeAddress)
		if err != nil {
			return err
		}

		slashCtx := ctx.WithContext(context.WithValue(ctx.Context(), constants.CtxMetricLabels, []metrics.Label{
			telemetry.NewLabel("address", na.NodeAddress.String()),
			telemetry.NewLabel("reason", "miss_block_sign"),
		}))
		if err := s.keeper.IncNodeAccountSlashPoints(slashCtx, na.NodeAddress, missBlockSignSlashPoints); err != nil {
			ctx.Logger().Error("fail to increase node account slash points", "error", err, "address", na.NodeAddress.String())
		}

		return s.keeper.SetNodeAccount(ctx, na)
	}

	return fmt.Errorf("could not find active node account with validator address: %s", addr)
}

// LackObserving Slash node accounts that didn't observe a single inbound txn
func (s *SlasherV125) LackObserving(ctx cosmos.Context, constAccessor constants.ConstantValues) error {
	signingTransPeriod := constAccessor.GetInt64Value(constants.SigningTransactionPeriod)
	height := ctx.BlockHeight()
	if height < signingTransPeriod {
		return nil
	}
	heightToCheck := height - signingTransPeriod
	tx, err := s.keeper.GetTxOut(ctx, heightToCheck)
	if err != nil {
		return fmt.Errorf("fail to get txout for block height(%d): %w", heightToCheck, err)
	}
	// no txout , return
	if tx == nil || tx.IsEmpty() {
		return nil
	}
	for _, item := range tx.TxArray {
		if item.InHash.IsEmpty() {
			continue
		}
		if item.InHash.Equals(common.BlankTxID) {
			continue
		}
		if err := s.slashNotObserving(ctx, item.InHash, constAccessor); err != nil {
			ctx.Logger().Error("fail to slash not observing", "error", err)
		}
	}

	return nil
}

func (s *SlasherV125) slashNotObserving(ctx cosmos.Context, txHash common.TxID, constAccessor constants.ConstantValues) error {
	voter, err := s.keeper.GetObservedTxInVoter(ctx, txHash)
	if err != nil {
		return fmt.Errorf("fail to get observe txin voter (%s): %w", txHash.String(), err)
	}

	if len(voter.Txs) == 0 {
		return nil
	}

	nodes, err := s.keeper.ListActiveValidators(ctx)
	if err != nil {
		return fmt.Errorf("unable to get list of active accounts: %w", err)
	}
	if len(voter.Txs) > 0 {
		tx := voter.Tx
		if !tx.IsEmpty() && len(tx.Signers) > 0 {
			height := voter.Height
			if tx.IsFinal() {
				height = voter.FinalisedHeight
			}
			// as long as the node has voted one of the tx , regardless finalised or not , it should not be slashed
			var allSigners []cosmos.AccAddress
			for _, item := range voter.Txs {
				allSigners = append(allSigners, item.GetSigners()...)
			}
			s.checkSignerAndSlash(ctx, nodes, height, allSigners, constAccessor)
		}
	}
	return nil
}

func (s *SlasherV125) checkSignerAndSlash(ctx cosmos.Context, nodes NodeAccounts, blockHeight int64, signers []cosmos.AccAddress, constAccessor constants.ConstantValues) {
	for _, na := range nodes {
		// the node is active after the tx finalised
		if na.ActiveBlockHeight > blockHeight {
			continue
		}
		found := false
		for _, addr := range signers {
			if na.NodeAddress.Equals(addr) {
				found = true
				break
			}
		}
		// this na is not found, therefore it should be slashed
		if !found {
			lackOfObservationPenalty := constAccessor.GetInt64Value(constants.LackOfObservationPenalty)
			slashCtx := ctx.WithContext(context.WithValue(ctx.Context(), constants.CtxMetricLabels, []metrics.Label{
				telemetry.NewLabel("reason", "not_observing"),
			}))
			if err := s.keeper.IncNodeAccountSlashPoints(slashCtx, na.NodeAddress, lackOfObservationPenalty); err != nil {
				ctx.Logger().Error("fail to inc slash points", "error", err)
			}
		}
	}
}

// LackSigning slash account that fail to sign tx
func (s *SlasherV125) LackSigning(ctx cosmos.Context, mgr Manager) error {
	var resultErr error
	signingTransPeriod := mgr.GetConstants().GetInt64Value(constants.SigningTransactionPeriod)
	if ctx.BlockHeight() < signingTransPeriod {
		return nil
	}
	height := ctx.BlockHeight() - signingTransPeriod
	txs, err := s.keeper.GetTxOut(ctx, height)
	if err != nil {
		return fmt.Errorf("fail to get txout from block height(%d): %w", height, err)
	}
	for i, tx := range txs.TxArray {
		if !common.CurrentChainNetwork.SoftEquals(tx.ToAddress.GetNetwork(mgr.GetVersion(), tx.Chain)) {
			continue // skip this transaction
		}
		if tx.OutHash.IsEmpty() {
			// Slash node account for not sending funds
			vault, err := s.keeper.GetVault(ctx, tx.VaultPubKey)
			if err != nil {
				// in some edge cases the vault may no longer exists, in which
				// case log and continue with rescheduling
				ctx.Logger().Error("Unable to get vault", "error", err, "vault pub key", tx.VaultPubKey.String())
			}

			// don't reschedule transactions on frozen vaults. This will cause
			// txns to be trapped in a specific asgard forever, which is the
			// expected result. This is here to protect the network from a
			// round7 attack
			if len(vault.Frozen) > 0 {
				chains, err := common.NewChains(vault.Frozen)
				if err != nil {
					ctx.Logger().Error("failed to convert chains", "error", err)
				}
				if chains.Has(tx.Coin.Asset.GetChain()) {
					etx := common.Tx{
						ID:        tx.InHash,
						Chain:     tx.Chain,
						ToAddress: tx.ToAddress,
						Coins:     []common.Coin{tx.Coin},
						Gas:       tx.MaxGas,
						Memo:      tx.Memo,
					}
					eve := NewEventSecurity(etx, "skipping reschedule on frozen vault")
					if err := mgr.EventMgr().EmitEvent(ctx, eve); err != nil {
						ctx.Logger().Error("fail to emit security event", "error", err)
					}
					continue // skip this transaction
				}
			}

			memo, _ := ParseMemoWithTHORNames(ctx, s.keeper, tx.Memo) // ignore err
			if memo.IsInternal() {
				// there is a different mechanism for rescheduling outbound
				// transactions for migration transactions
				continue
			}
			var voter ObservedTxVoter
			if !memo.IsType(TxRagnarok) {
				voter, err = s.keeper.GetObservedTxInVoter(ctx, tx.InHash)
				if err != nil {
					ctx.Logger().Error("fail to get observed tx voter", "error", err)
					resultErr = fmt.Errorf("failed to get observed tx voter: %w", err)
					continue
				}
			}

			maxOutboundAttempts := mgr.Keeper().GetConfigInt64(ctx, constants.MaxOutboundAttempts)
			if maxOutboundAttempts > 0 {
				age := ctx.BlockHeight() - voter.FinalisedHeight
				attempts := age / signingTransPeriod
				if attempts >= maxOutboundAttempts {
					ctx.Logger().Info("txn dropped, too many attempts", "hash", tx.InHash)
					continue
				}
			}

			// if vault is inactive, do not reassign the outbound txn to
			// another vault
			if vault.Status == InactiveVault {
				ctx.Logger().Info("cannot reassign tx from inactive vault", "hash", tx.InHash)
				continue
			}

			if s.needsNewVault(ctx, mgr, vault, signingTransPeriod, voter.FinalisedHeight, tx.InHash, tx.VaultPubKey) {
				active, err := s.keeper.GetAsgardVaultsByStatus(ctx, ActiveVault)
				if err != nil {
					return fmt.Errorf("fail to get active asgard vaults: %w", err)
				}
				available := active.Has(tx.Coin).SortBy(tx.Coin.Asset)
				if len(available) == 0 {
					// we need to give it somewhere to send from, even if that
					// asgard doesn't have enough funds. This is because if we
					// don't the transaction will just be dropped on the floor,
					// which is bad. Instead it may try to send from an asgard that
					// doesn't have enough funds, fail, and then get rescheduled
					// again later. Maybe by then the network will have enough
					// funds to satisfy.
					// TODO add split logic to send it out from multiple asgards in
					// this edge case.
					ctx.Logger().Error("unable to determine asgard vault to send funds, trying first asgard")
					if len(active) > 0 {
						vault = active[0]
					}
				} else {
					rep := int(tx.InHash.Int64() + ctx.BlockHeight())
					if vault.PubKey.Equals(available[rep%len(available)].PubKey) {
						// looks like the new vault is going to be the same as the
						// old vault, increment rep to ensure a differ asgard is
						// chosen (if there is more than one option)
						rep++
					}
					vault = available[rep%len(available)]
				}
				if !memo.IsType(TxRagnarok) {
					// update original tx action in observed tx
					// check observedTx has done status. Skip if it does already.
					voterTx := voter.GetTx(NodeAccounts{})
					if voterTx.IsDone(len(voter.Actions)) {
						if len(voterTx.OutHashes) > 0 && len(voterTx.GetOutHashes()) > 0 {
							txs.TxArray[i].OutHash = voterTx.GetOutHashes()[0]
						}
						continue
					}

					// update the actions in the voter with the new vault pubkey
					for i, action := range voter.Actions {
						if action.Equals(tx) {
							voter.Actions[i].VaultPubKey = vault.PubKey

							if tx.Aggregator != "" || tx.AggregatorTargetAsset != "" || tx.AggregatorTargetLimit != nil {
								ctx.Logger().Info("clearing aggregator fields on outbound reassignment", "hash", tx.InHash)

								// Here, simultaneously clear the Aggregator information for a reassigned TxOutItem and its Actions item
								// so that a SwapOut will send the THORChain output asset instead of cycling and swallowingif
								// (and maybe failing with slashes) if something goes wrong.
								tx.Aggregator = ""
								tx.AggregatorTargetAsset = ""
								tx.AggregatorTargetLimit = nil
								voter.Actions[i].Aggregator = ""
								voter.Actions[i].AggregatorTargetAsset = ""
								voter.Actions[i].AggregatorTargetLimit = nil
							}
						}
					}
					s.keeper.SetObservedTxInVoter(ctx, voter)

				}
				// Save the tx to as a new tx, select Asgard to send it this time.
				tx.VaultPubKey = vault.PubKey

				// update max gas
				maxGas, err := mgr.GasMgr().GetMaxGas(ctx, tx.Chain)
				if err != nil {
					ctx.Logger().Error("fail to get max gas", "error", err)
				} else {
					tx.MaxGas = common.Gas{maxGas}
					// Update MaxGas in ObservedTxVoter action as well
					if err := updateTxOutGas(ctx, s.keeper, tx, common.Gas{maxGas}); err != nil {
						ctx.Logger().Error("Failed to update MaxGas of action in ObservedTxVoter", "hash", tx.InHash, "error", err)
					}
				}
				// Equals checks GasRate so update actions GasRate too (before updating in the queue item)
				// for future updates of MaxGas, which must match for matchActionItem in AddOutTx.
				gasRate := int64(mgr.GasMgr().GetGasRate(ctx, tx.Chain).Uint64())
				if err := updateTxOutGasRate(ctx, s.keeper, tx, gasRate); err != nil {
					ctx.Logger().Error("Failed to update GasRate of action in ObservedTxVoter", "hash", tx.InHash, "error", err)
				}
				tx.GasRate = gasRate
			}

			// if a pool with the asset name doesn't exist, skip rescheduling
			if !tx.Coin.Asset.IsRune() && !s.keeper.PoolExist(ctx, tx.Coin.Asset) {
				ctx.Logger().Error("fail to add outbound tx", "error", "coin is not rune and does not have an associated pool")
				continue
			}

			err = mgr.TxOutStore().UnSafeAddTxOutItem(ctx, mgr, tx, ctx.BlockHeight())
			if err != nil {
				ctx.Logger().Error("fail to add outbound tx", "error", err)
				resultErr = fmt.Errorf("failed to add outbound tx: %w", err)
				continue
			}
			// because the txout item has been rescheduled, thus mark the replaced tx out item as already send out, even it is not
			// in this way bifrost will not send it out again cause node to be slashed
			txs.TxArray[i].OutHash = common.BlankTxID
		}
	}
	if !txs.IsEmpty() {
		if err := s.keeper.SetTxOut(ctx, txs); err != nil {
			return fmt.Errorf("fail to save tx out : %w", err)
		}
	}

	return resultErr
}

// SlashVault thorchain keep monitoring the outbound tx from asgard pool
// usually the txout is triggered by thorchain itself by
// adding an item into the txout array, refer to TxOutItem for the detail, the
// TxOutItem contains a specific coin and amount.  if somehow thorchain
// discover signer send out fund more than the amount specified in TxOutItem,
// it will slash the node account who does that by taking 1.5 * extra fund from
// node account's bond and subsidise the pool that actually lost it.
func (s *SlasherV125) SlashVault(ctx cosmos.Context, vaultPK common.PubKey, coins common.Coins, mgr Manager) error {
	if coins.IsEmpty() {
		return nil
	}

	vault, err := s.keeper.GetVault(ctx, vaultPK)
	if err != nil {
		return fmt.Errorf("fail to get slash vault (pubkey %s), %w", vaultPK, err)
	}
	membership := vault.GetMembership()

	// sum the total bond of membership of the vault
	totalBond := cosmos.ZeroUint()
	for _, member := range membership {
		na, err := s.keeper.GetNodeAccountByPubKey(ctx, member)
		if err != nil {
			ctx.Logger().Error("fail to get node account bond", "pk", member, "error", err)
			continue
		}
		totalBond = totalBond.Add(na.Bond)
	}

	for _, coin := range coins {
		if coin.IsEmpty() {
			continue
		}
		pool, err := s.keeper.GetPool(ctx, coin.Asset)
		if err != nil {
			ctx.Logger().Error("fail to get pool for slash", "asset", coin.Asset, "error", err)
			continue
		}
		// THORChain doesn't even have a pool for the asset
		if pool.IsEmpty() {
			ctx.Logger().Error("cannot slash for an empty pool", "asset", coin.Asset)
			continue
		}

		stolenAssetValue := coin.Amount
		vaultAmount := vault.GetCoin(coin.Asset).Amount
		if stolenAssetValue.GT(vaultAmount) {
			stolenAssetValue = vaultAmount
		}
		if stolenAssetValue.GT(pool.BalanceAsset) {
			stolenAssetValue = pool.BalanceAsset
		}

		// stolenRuneValue is the value in RUNE of the missing funds
		stolenRuneValue := pool.AssetValueInRune(stolenAssetValue)

		if stolenRuneValue.IsZero() {
			continue
		}

		penaltyPts := mgr.Keeper().GetConfigInt64(ctx, constants.SlashPenalty)
		// total slash amount is penaltyPts the RUNE value of the missing funds
		totalRuneToSlash := common.GetUncappedShare(cosmos.NewUint(uint64(penaltyPts)), cosmos.NewUint(10_000), stolenRuneValue)
		totalRuneSlashed := cosmos.ZeroUint()
		pauseOnSlashThreshold := mgr.Keeper().GetConfigInt64(ctx, constants.PauseOnSlashThreshold)
		if pauseOnSlashThreshold > 0 && totalRuneToSlash.GTE(cosmos.NewUint(uint64(pauseOnSlashThreshold))) {
			// set mimirs to pause the chain
			key := fmt.Sprintf("Halt%sChain", coin.Asset.Chain)
			s.keeper.SetMimir(ctx, key, ctx.BlockHeight())
			mimirEvent := NewEventSetMimir(strings.ToUpper(key), strconv.FormatInt(ctx.BlockHeight(), 10))
			if err := mgr.EventMgr().EmitEvent(ctx, mimirEvent); err != nil {
				ctx.Logger().Error("fail to emit set_mimir event", "error", err)
			}
		}
		for _, member := range membership {
			na, err := s.keeper.GetNodeAccountByPubKey(ctx, member)
			if err != nil {
				ctx.Logger().Error("fail to get node account for slash", "pk", member, "error", err)
				continue
			}
			if na.Bond.IsZero() {
				ctx.Logger().Info("validator's bond is zero, can't be slashed", "node address", na.NodeAddress.String())
				continue
			}
			runeSlashed := s.slashAndUpdateNodeAccount(ctx, na, coin, vault, totalBond, totalRuneToSlash)
			totalRuneSlashed = totalRuneSlashed.Add(runeSlashed)
		}

		//  2/3 of the total slashed RUNE value to asgard
		//  1/3 of the total slashed RUNE value to reserve
		runeToAsgard := stolenRuneValue

		// stolenRuneValue is the total value in RUNE of stolen coins, but totalRuneSlashed is
		// the total amount able to be slashed from Nodes, credit the pool only totalRuneSlashed
		if totalRuneSlashed.LT(stolenRuneValue) {
			// this should theoretically never happen
			ctx.Logger().Info("total slashed bond amount is less than RUNE value", "slashed_bond", totalRuneSlashed.String(), "rune_value", stolenRuneValue.String())
			runeToAsgard = totalRuneSlashed
		}
		runeToReserve := common.SafeSub(totalRuneSlashed, runeToAsgard)

		if !runeToReserve.IsZero() {
			if err := s.keeper.SendFromModuleToModule(ctx, BondName, ReserveName, common.NewCoins(common.NewCoin(common.RuneAsset(), runeToReserve))); err != nil {
				ctx.Logger().Error("fail to send slash funds to reserve module", "pk", vaultPK, "error", err)
			}
		}
		if !runeToAsgard.IsZero() {
			if err := s.keeper.SendFromModuleToModule(ctx, BondName, AsgardName, common.NewCoins(common.NewCoin(common.RuneAsset(), runeToAsgard))); err != nil {
				ctx.Logger().Error("fail to send slash fund to asgard module", "pk", vaultPK, "error", err)
			}
			s.updatePoolFromSlash(ctx, pool, common.NewCoin(coin.Asset, stolenAssetValue), runeToAsgard, mgr)
		}
	}

	return nil
}

// slashAndUpdateNodeAccount slashes a NodeAccount a portion of the value of coin based on their
// portion of the total bond of the offending Vault's membership. Return the amount of RUNE slashed
func (s SlasherV125) slashAndUpdateNodeAccount(ctx cosmos.Context, na types.NodeAccount, coin common.Coin, vault types.Vault, totalBond, totalSlashAmountInRune cosmos.Uint) cosmos.Uint {
	slashAmountRune := common.GetSafeShare(na.Bond, totalBond, totalSlashAmountInRune)
	if slashAmountRune.GT(na.Bond) {
		ctx.Logger().Info("slash amount is larger than bond", "slash amount", slashAmountRune, "bond", na.Bond)
		slashAmountRune = na.Bond
	}

	ctx.Logger().Info("slash node account", "node address", na.NodeAddress.String(), "amount", slashAmountRune.String(), "total slash amount", totalSlashAmountInRune)
	na.Bond = common.SafeSub(na.Bond, slashAmountRune)

	tx := common.Tx{}
	tx.ID = common.BlankTxID
	tx.FromAddress = na.BondAddress
	bondEvent := NewEventBond(slashAmountRune, BondCost, tx)
	if err := s.eventMgr.EmitEvent(ctx, bondEvent); err != nil {
		ctx.Logger().Error("fail to emit bond event", "error", err)
	}

	metricLabels, _ := ctx.Context().Value(constants.CtxMetricLabels).([]metrics.Label)
	slashAmountRuneFloat, _ := new(big.Float).SetInt(slashAmountRune.BigInt()).Float32()
	telemetry.IncrCounterWithLabels(
		[]string{"thornode", "bond_slash"},
		slashAmountRuneFloat,
		append(
			metricLabels,
			telemetry.NewLabel("address", na.NodeAddress.String()),
			telemetry.NewLabel("coin_symbol", coin.Asset.Symbol.String()),
			telemetry.NewLabel("coin_chain", string(coin.Asset.Chain)),
			telemetry.NewLabel("vault_type", vault.Type.String()),
			telemetry.NewLabel("vault_status", vault.Status.String()),
		),
	)

	if err := s.keeper.SetNodeAccount(ctx, na); err != nil {
		ctx.Logger().Error("fail to save node account for slash", "error", err)
	}

	return slashAmountRune
}

// IncSlashPoints will increase the given account's slash points
func (s *SlasherV125) IncSlashPoints(ctx cosmos.Context, point int64, addresses ...cosmos.AccAddress) {
	for _, addr := range addresses {
		if err := s.keeper.IncNodeAccountSlashPoints(ctx, addr, point); err != nil {
			ctx.Logger().Error("fail to increase node account slash point", "error", err, "address", addr.String())
		}
	}
}

// DecSlashPoints will decrease the given account's slash points
func (s *SlasherV125) DecSlashPoints(ctx cosmos.Context, point int64, addresses ...cosmos.AccAddress) {
	for _, addr := range addresses {
		if err := s.keeper.DecNodeAccountSlashPoints(ctx, addr, point); err != nil {
			ctx.Logger().Error("fail to decrease node account slash point", "error", err, "address", addr.String())
		}
	}
}

// updatePoolFromSlash updates a pool's depths and emits appropriate events after a slash
func (s *SlasherV125) updatePoolFromSlash(ctx cosmos.Context, pool types.Pool, stolenAsset common.Coin, runeCreditAmt cosmos.Uint, mgr Manager) {
	pool.BalanceAsset = common.SafeSub(pool.BalanceAsset, stolenAsset.Amount)
	pool.BalanceRune = pool.BalanceRune.Add(runeCreditAmt)
	if err := s.keeper.SetPool(ctx, pool); err != nil {
		ctx.Logger().Error("fail to save pool for slash", "asset", pool.Asset, "error", err)
	}
	poolSlashAmt := []PoolAmt{
		{
			Asset:  pool.Asset,
			Amount: 0 - int64(stolenAsset.Amount.Uint64()),
		},
		{
			Asset:  common.RuneAsset(),
			Amount: int64(runeCreditAmt.Uint64()),
		},
	}
	eventSlash := NewEventSlash(pool.Asset, poolSlashAmt)
	if err := mgr.EventMgr().EmitEvent(ctx, eventSlash); err != nil {
		ctx.Logger().Error("fail to emit slash event", "error", err)
	}
}

func (s *SlasherV125) needsNewVault(ctx cosmos.Context, mgr Manager, vault Vault, signingTransPeriod, startHeight int64, inhash common.TxID, pk common.PubKey) bool {
	outhashes := mgr.Keeper().GetObservedLink(ctx, inhash)
	if len(outhashes) == 0 {
		return true
	}

	for _, hash := range outhashes {
		voter, err := mgr.Keeper().GetObservedTxOutVoter(ctx, hash)
		if err != nil {
			ctx.Logger().Error("fail to get txout voter", "hash", hash, "error", err)
			continue
		}
		// in the event there are multiple outbounds for a given inhash, we
		// focus on the matching pubkey
		signers := make(map[string]bool)
		for _, tx1 := range voter.Txs {
			if tx1.ObservedPubKey.Equals(pk) {
				for _, tx := range voter.Txs {
					if !tx.Tx.ID.Equals(hash) {
						continue
					}
					for _, signer := range tx.Signers {
						// Uniquely record each signer for this outbound hash.
						signers[signer] = true
					}
				}
			}
		}
		if len(signers) > 0 {
			var count int // count the number of signers from the assigned vault
			for _, member := range vault.Membership {
				pk, err := common.NewPubKey(member)
				if err != nil {
					continue
				}
				addr, err := pk.GetThorAddress()
				if err != nil {
					continue
				}
				if _, ok := signers[addr.String()]; ok {
					count++
				}
			}
			// if a super majority of vault members have observed the outbound,
			// then we should not reschedule. If a vault says it sent it, it
			// sent it and shouldn't get another vault to send it (potentially
			// a second time)
			if count > 0 && HasSuperMajority(count, len(vault.Membership)) {
				return false
			}
			maxHeight := startHeight + ((int64(len(signers)) + 1) * signingTransPeriod)
			return maxHeight < ctx.BlockHeight()
		}

	}

	return true
}
