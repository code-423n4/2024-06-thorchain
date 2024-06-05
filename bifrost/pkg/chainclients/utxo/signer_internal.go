package utxo

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/blang/semver"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"

	"github.com/eager7/dogutil"
	dogetxscript "gitlab.com/thorchain/bifrost/dogd-txscript"

	"github.com/gcash/bchutil"
	bchtxscript "gitlab.com/thorchain/bifrost/bchd-txscript"

	"github.com/ltcsuite/ltcutil"
	ltctxscript "gitlab.com/thorchain/bifrost/ltcd-txscript"

	"github.com/btcsuite/btcutil"
	btctxscript "gitlab.com/thorchain/bifrost/txscript"

	stypes "gitlab.com/thorchain/thornode/bifrost/thorclient/types"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	mem "gitlab.com/thorchain/thornode/x/thorchain/memo"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

////////////////////////////////////////////////////////////////////////////////////////
// UTXO Selection
////////////////////////////////////////////////////////////////////////////////////////

func (c *Client) getMaximumUtxosToSpend() int64 {
	const mimirMaxUTXOsToSpend = `MaxUTXOsToSpend`
	utxosToSpend, err := c.bridge.GetMimir(mimirMaxUTXOsToSpend)
	if err != nil {
		c.log.Err(err).Msg("fail to get MaxUTXOsToSpend")
	}
	if utxosToSpend <= 0 {
		utxosToSpend = c.cfg.UTXO.MaxUTXOsToSpend
	}
	return utxosToSpend
}

// getAllUtxos will iterate unspend utxos for the given address and return the oldest
// set of utxos that can cover the amount.
func (c *Client) getUtxoToSpend(pubkey common.PubKey, total float64) ([]btcjson.ListUnspentResult, error) {
	// get all unspent utxos
	addr, err := pubkey.GetAddress(c.cfg.ChainID)
	if err != nil {
		return nil, fmt.Errorf("fail to get address from pubkey(%s): %w", pubkey, err)
	}
	utxos, err := c.rpc.ListUnspent(addr.String())
	if err != nil {
		return nil, fmt.Errorf("fail to get UTXOs: %w", err)
	}

	// spend UTXO older to younger
	sort.SliceStable(utxos, func(i, j int) bool {
		if utxos[i].Confirmations > utxos[j].Confirmations {
			return true
		} else if utxos[i].Confirmations < utxos[j].Confirmations {
			return false
		}
		return utxos[i].TxID < utxos[j].TxID
	})

	var result []btcjson.ListUnspentResult
	var toSpend float64
	minUTXOAmt := btcutil.Amount(c.cfg.ChainID.DustThreshold().Uint64()).ToBTC()
	utxosToSpend := c.getMaximumUtxosToSpend() // can be set by mimir

	for _, item := range utxos {
		if !c.isValidUTXO(item.ScriptPubKey) {
			c.log.Warn().Str("script", item.ScriptPubKey).Msgf("invalid utxo, unable to spend")
			continue
		}

		if item.Confirmations < c.cfg.UTXO.MinUTXOConfirmations || item.Amount < minUTXOAmt {
			// use all UTXOs sent from asgard, regardless of confirmations or dust threshold
			isSelfTx := c.isSelfTransaction(item.TxID)

			// confirm sender of the UTXO is not asgard in case of lost block meta
			if !isSelfTx {
				isSelfTx = c.isFromAsgard(item.TxID)
			}
			if !isSelfTx {
				continue
			}
		}

		result = append(result, item)
		toSpend += item.Amount

		// in the scenario that there are too many unspent utxos available, make sure it
		// doesn't spend too much as too much UTXO will cause huge pressure on TSS, also
		// make sure it will spend at least maxUTXOsToSpend so the UTXOs will be
		// consolidated
		if int64(len(result)) >= utxosToSpend && toSpend >= total {
			break
		}
	}
	return result, nil
}

// vinsUnspent will return true if all the vins are unspent.
func (c *Client) vinsUnspent(tx stypes.TxOutItem, vins []*wire.TxIn) (bool, error) {
	// get all unspent utxos
	addr, err := tx.VaultPubKey.GetAddress(c.cfg.ChainID)
	if err != nil {
		return false, fmt.Errorf("fail to get address from pubkey(%s): %w", tx.VaultPubKey, err)
	}
	utxos, err := c.rpc.ListUnspent(addr.String())
	if err != nil {
		return false, fmt.Errorf("fail to get UTXOs: %w", err)
	}
	unspent := make(map[string]bool, len(utxos))
	for _, utxo := range utxos {
		unspent[utxo.TxID] = true
	}

	// return false if any vin is spent
	allUnspent := true
	for _, vin := range vins {
		if !unspent[vin.PreviousOutPoint.Hash.String()] {
			c.log.Warn().
				Stringer("in_hash", tx.InHash).
				Stringer("vin", vin.PreviousOutPoint).
				Msg("vin is spent")
			allUnspent = false
		}
	}

	return allUnspent, nil
}

// isSelfTransaction check the block meta to see whether the transactions is broadcast
// by ourselves if the transaction is broadcast by ourselves, then we should be able to
// spend the UTXO even it is still in mempool as such we could daisy chain the outbound
// transaction
func (c *Client) isSelfTransaction(txID string) bool {
	bms, err := c.temporalStorage.GetBlockMetas()
	if err != nil {
		c.log.Err(err).Msg("fail to get block metas")
		return false
	}
	for _, item := range bms {
		for _, tx := range item.SelfTransactions {
			if strings.EqualFold(tx, txID) {
				c.log.Debug().Msgf("%s is self transaction", txID)
				return true
			}
		}
	}
	return false
}

func (c *Client) getPaymentAmount(tx stypes.TxOutItem) float64 {
	amtToPay1e8 := tx.Coins.GetCoin(c.cfg.ChainID.GetGasAsset()).Amount.Uint64()
	amtToPay := btcutil.Amount(int64(amtToPay1e8)).ToBTC()
	if !tx.MaxGas.IsEmpty() {
		gasAmt := tx.MaxGas.ToCoins().GetCoin(c.cfg.ChainID.GetGasAsset()).Amount
		amtToPay += btcutil.Amount(int64(gasAmt.Uint64())).ToBTC()
	}
	return amtToPay
}

// getSourceScript retrieve pay to addr script from tx source
func (c *Client) getSourceScript(tx stypes.TxOutItem) ([]byte, error) {
	sourceAddr, err := tx.VaultPubKey.GetAddress(c.cfg.ChainID)
	if err != nil {
		return nil, fmt.Errorf("fail to get source address: %w", err)
	}

	switch c.cfg.ChainID {
	case common.DOGEChain:
		var addr dogutil.Address
		addr, err = dogutil.DecodeAddress(sourceAddr.String(), c.getChainCfgDOGE())
		if err != nil {
			return nil, fmt.Errorf("fail to decode source address(%s): %w", sourceAddr.String(), err)
		}
		return dogetxscript.PayToAddrScript(addr)
	case common.BCHChain:
		var addr bchutil.Address
		addr, err = bchutil.DecodeAddress(sourceAddr.String(), c.getChainCfgBCH())
		if err != nil {
			return nil, fmt.Errorf("fail to decode source address(%s): %w", sourceAddr.String(), err)
		}
		return bchtxscript.PayToAddrScript(addr)
	case common.LTCChain:
		var addr ltcutil.Address
		addr, err = ltcutil.DecodeAddress(sourceAddr.String(), c.getChainCfgLTC())
		if err != nil {
			return nil, fmt.Errorf("fail to decode source address(%s): %w", sourceAddr.String(), err)
		}
		return ltctxscript.PayToAddrScript(addr)
	case common.BTCChain:
		var addr btcutil.Address
		addr, err = btcutil.DecodeAddress(sourceAddr.String(), c.getChainCfgBTC())
		if err != nil {
			return nil, fmt.Errorf("fail to decode source address(%s): %w", sourceAddr.String(), err)
		}
		return btctxscript.PayToAddrScript(addr)
	default:
		c.log.Fatal().Msg("unsupported chain")
		return nil, nil
	}
}

////////////////////////////////////////////////////////////////////////////////////////
// Build Transaction
////////////////////////////////////////////////////////////////////////////////////////

// TODO: Cleanup magic numbers and/or improve comment specificity.
// estimateTxSize will create a temporary MsgTx, and use it to estimate the final tx size
// the value in the temporary MsgTx is not real
// https://bitcoinops.org/en/tools/calc-size/
func (c *Client) estimateTxSize(memo string, txes []btcjson.ListUnspentResult) int64 {
	switch c.cfg.ChainID {
	case common.DOGEChain, common.BCHChain:
		// overhead - 10
		// Per input - 148
		// Per output - 34 , we might have 1 / 2 output , depends on the circumstances , here we only count 1  output , would rather underestimate
		// so we won't hit absurd high fee issue
		// overhead for NULL DATA - 9 , len(memo) is the size of memo
		return int64(10 + 148*len(txes) + 34 + 9 + len([]byte(memo)))
	case common.LTCChain, common.BTCChain:
		// overhead - 10.75
		// Per Input - 67.75
		// Per output - 31 , we sometimes have 2 output , and sometimes only have 1 , it depends ,here we only count 1
		// it is better to underestimate rather than over estimate
		// 10.5 overhead for null data
		// len(memo) is the size of memo put in null data
		// these get us very close to the final vbytes.
		// multiple by 100 , and then add, so don't need to deal with float
		return int64((1075+6775*len(txes)+1050)/100) + int64(31+len([]byte(memo)))
	default:
		c.log.Fatal().Msg("unsupported chain")
		return 0
	}
}

func (c *Client) getGasCoin(tx stypes.TxOutItem, vSize int64) common.Coin {
	gasRate := tx.GasRate

	// if the gas rate is zero, try to get from last transaction fee
	if gasRate == 0 {
		fee, vBytes, err := c.temporalStorage.GetTransactionFee()
		if err != nil {
			c.log.Error().Err(err).Msg("fail to get previous transaction fee from local storage")
			return common.NewCoin(c.cfg.ChainID.GetGasAsset(), cosmos.NewUint(uint64(vSize*gasRate)))
		}
		if fee != 0.0 && vSize != 0 {
			var amt btcutil.Amount
			amt, err = btcutil.NewAmount(fee)
			if err != nil {
				c.log.Err(err).Msg("fail to convert amount from float64 to int64")
			} else {
				gasRate = int64(amt) / int64(vBytes) // sats per vbyte
			}
		}
	}

	// default to configured value
	if gasRate == 0 {
		gasRate = c.cfg.UTXO.DefaultSatsPerVByte
	}

	return common.NewCoin(c.cfg.ChainID.GetGasAsset(), cosmos.NewUint(uint64(gasRate*vSize)))
}

func (c *Client) buildTx(tx stypes.TxOutItem, sourceScript []byte) (*wire.MsgTx, map[string]int64, error) {
	txes, err := c.getUtxoToSpend(tx.VaultPubKey, c.getPaymentAmount(tx))
	if err != nil {
		return nil, nil, fmt.Errorf("fail to get unspent UTXO")
	}
	redeemTx := wire.NewMsgTx(wire.TxVersion)
	totalAmt := int64(0)
	individualAmounts := make(map[string]int64, len(txes))
	for _, item := range txes {
		var txID *chainhash.Hash
		txID, err = chainhash.NewHashFromStr(item.TxID)
		if err != nil {
			return nil, nil, fmt.Errorf("fail to parse txID(%s): %w", item.TxID, err)
		}
		// double check that the utxo is still valid
		outputPoint := wire.NewOutPoint(txID, item.Vout)
		sourceTxIn := wire.NewTxIn(outputPoint, nil, nil)
		redeemTx.AddTxIn(sourceTxIn)
		var amt btcutil.Amount
		amt, err = btcutil.NewAmount(item.Amount)
		if err != nil {
			return nil, nil, fmt.Errorf("fail to parse amount(%f): %w", item.Amount, err)
		}
		individualAmounts[fmt.Sprintf("%s-%d", txID, item.Vout)] = int64(amt)
		totalAmt += int64(amt)
	}

	var buf []byte
	switch c.cfg.ChainID {
	case common.DOGEChain:
		var outputAddr dogutil.Address
		outputAddr, err = dogutil.DecodeAddress(tx.ToAddress.String(), c.getChainCfgDOGE())
		if err != nil {
			return nil, nil, fmt.Errorf("fail to decode next address: %w", err)
		}
		buf, err = dogetxscript.PayToAddrScript(outputAddr)
		if err != nil {
			return nil, nil, fmt.Errorf("fail to get pay to address script: %w", err)
		}
	case common.BCHChain:
		var outputAddr bchutil.Address
		outputAddr, err = bchutil.DecodeAddress(tx.ToAddress.String(), c.getChainCfgBCH())
		if err != nil {
			return nil, nil, fmt.Errorf("fail to decode next address: %w", err)
		}
		buf, err = bchtxscript.PayToAddrScript(outputAddr)
		if err != nil {
			return nil, nil, fmt.Errorf("fail to get pay to address script: %w", err)
		}
	case common.LTCChain:
		var outputAddr ltcutil.Address
		outputAddr, err = ltcutil.DecodeAddress(tx.ToAddress.String(), c.getChainCfgLTC())
		if err != nil {
			return nil, nil, fmt.Errorf("fail to decode next address: %w", err)
		}
		buf, err = ltctxscript.PayToAddrScript(outputAddr)
		if err != nil {
			return nil, nil, fmt.Errorf("fail to get pay to address script: %w", err)
		}
	case common.BTCChain:
		var outputAddr btcutil.Address
		outputAddr, err = btcutil.DecodeAddress(tx.ToAddress.String(), c.getChainCfgBTC())
		if err != nil {
			return nil, nil, fmt.Errorf("fail to decode next address: %w", err)
		}
		buf, err = btctxscript.PayToAddrScript(outputAddr)
		if err != nil {
			return nil, nil, fmt.Errorf("fail to get pay to address script: %w", err)
		}
	default:
		c.log.Fatal().Msg("unsupported chain")
	}

	coinToCustomer := tx.Coins.GetCoin(c.cfg.ChainID.GetGasAsset())
	totalSize := c.estimateTxSize(tx.Memo, txes)

	// maxFee in sats
	maxFeeSats := totalSize * c.cfg.UTXO.MaxSatsPerVByte
	gasCoin := c.getGasCoin(tx, totalSize)
	gasAmtSats := gasCoin.Amount.Uint64()

	// make sure the transaction fee is not more than the max, otherwise it might reject the transaction
	if gasAmtSats > uint64(maxFeeSats) {
		diffSats := gasAmtSats - uint64(maxFeeSats) // in sats
		c.log.Info().Msgf("gas amount: %d is larger than maximum fee: %d, diff: %d", gasAmtSats, uint64(maxFeeSats), diffSats)
		gasAmtSats = uint64(maxFeeSats)
	} else if gasAmtSats < c.minRelayFeeSats {
		diffStats := c.minRelayFeeSats - gasAmtSats
		c.log.Info().Msgf("gas amount: %d is less than min relay fee: %d, diff remove from customer: %d", gasAmtSats, c.minRelayFeeSats, diffStats)
		gasAmtSats = c.minRelayFeeSats
	}

	// TODO: remove version check after v132
	version, err := c.bridge.GetThorchainVersion()
	if err != nil {
		c.log.Err(err).Msg("fail to get thorchain version")
	}
	// TODO: after v132 remove 'var memo mem.Memo' and change 'memo, err =' to 'memo, err :='
	var memo mem.Memo
	if err == nil && version.GTE(semver.MustParse("1.132.0")) {
		// Parse the memo to be able to identify Migrate or Consolidate outbounds.
		memo, err = mem.ParseMemo(common.LatestVersion, tx.Memo)
		if err != nil {
			return nil, nil, fmt.Errorf("fail to parse memo: %w", err)
		}

		// if the total gas spend is more than max gas , then we have to take away some from the amount pay to customer
		if !tx.MaxGas.IsEmpty() {
			maxGasCoin := tx.MaxGas.ToCoins().GetCoin(c.cfg.ChainID.GetGasAsset())
			if gasAmtSats > maxGasCoin.Amount.Uint64() {
				c.log.Info().Msgf("max gas: %s, however estimated gas need %d", tx.MaxGas, gasAmtSats)
				gasAmtSats = maxGasCoin.Amount.Uint64()
			} else if gasAmtSats < maxGasCoin.Amount.Uint64() && memo.GetType() == mem.TxMigrate {
				// if the tx spend less gas then the estimated MaxGas , then the extra can be added to the coinToCustomer
				gap := maxGasCoin.Amount.Uint64() - gasAmtSats
				c.log.Info().Msgf("max gas is: %s, however only: %d is required, gap: %d goes to the vault migrated to", tx.MaxGas, gasAmtSats, gap)
				coinToCustomer.Amount = coinToCustomer.Amount.Add(cosmos.NewUint(gap))
			}
		} else if memo.GetType() == mem.TxConsolidate {
			gap := gasAmtSats
			c.log.Info().Msgf("consolidate tx, need gas: %d", gap)
			coinToCustomer.Amount = common.SafeSub(coinToCustomer.Amount, cosmos.NewUint(gap))
		}
	} else {
		// if the total gas spend is more than max gas , then we have to take away some from the amount pay to customer
		if !tx.MaxGas.IsEmpty() {
			maxGasCoin := tx.MaxGas.ToCoins().GetCoin(c.cfg.ChainID.GetGasAsset())
			if gasAmtSats > maxGasCoin.Amount.Uint64() {
				c.log.Info().Msgf("max gas: %s, however estimated gas need %d", tx.MaxGas, gasAmtSats)
				gasAmtSats = maxGasCoin.Amount.Uint64()
			} else if gasAmtSats < maxGasCoin.Amount.Uint64() {
				// if the tx spend less gas then the estimated MaxGas , then the extra can be added to the coinToCustomer
				gap := maxGasCoin.Amount.Uint64() - gasAmtSats
				c.log.Info().Msgf("max gas is: %s, however only: %d is required, gap: %d goes to customer", tx.MaxGas, gasAmtSats, gap)
				coinToCustomer.Amount = coinToCustomer.Amount.Add(cosmos.NewUint(gap))
			}
		} else {
			var memo mem.Memo
			memo, err = mem.ParseMemo(common.LatestVersion, tx.Memo)
			if err != nil {
				return nil, nil, fmt.Errorf("fail to parse memo: %w", err)
			}
			if memo.GetType() == mem.TxConsolidate {
				gap := gasAmtSats
				c.log.Info().Msgf("consolidate tx, need gas: %d", gap)
				coinToCustomer.Amount = common.SafeSub(coinToCustomer.Amount, cosmos.NewUint(gap))
			}
		}
	}

	gasAmt := btcutil.Amount(gasAmtSats)
	if err = c.temporalStorage.UpsertTransactionFee(gasAmt.ToBTC(), int32(totalSize)); err != nil {
		c.log.Err(err).Msg("fail to save gas info to UTXO storage")
	}

	// pay to customer
	redeemTxOut := wire.NewTxOut(int64(coinToCustomer.Amount.Uint64()), buf)
	redeemTx.AddTxOut(redeemTxOut)

	// balance to ourselves
	// add output to pay the balance back ourselves
	balance := totalAmt - redeemTxOut.Value - int64(gasAmt)
	c.log.Info().Msgf("total: %d, to customer: %d, gas: %d", totalAmt, redeemTxOut.Value, int64(gasAmt))
	if balance < 0 {
		return nil, nil, fmt.Errorf("not enough balance to pay customer: %d", balance)
	}
	if balance > 0 {
		c.log.Info().Msgf("send %d back to self", balance)
		redeemTx.AddTxOut(wire.NewTxOut(balance, sourceScript))
	}

	// memo
	if len(tx.Memo) != 0 {
		var nullDataScript []byte
		switch c.cfg.ChainID {
		case common.DOGEChain:
			nullDataScript, err = dogetxscript.NullDataScript([]byte(tx.Memo))
		case common.BCHChain:
			nullDataScript, err = bchtxscript.NullDataScript([]byte(tx.Memo))
		case common.LTCChain:
			nullDataScript, err = ltctxscript.NullDataScript([]byte(tx.Memo))
		case common.BTCChain:
			nullDataScript, err = btctxscript.NullDataScript([]byte(tx.Memo))
		default:
			c.log.Fatal().Msg("unsupported chain")
		}
		if err != nil {
			return nil, nil, fmt.Errorf("fail to generate null data script: %w", err)
		}
		redeemTx.AddTxOut(wire.NewTxOut(0, nullDataScript))
	}

	return redeemTx, individualAmounts, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// UTXO Consolidation
////////////////////////////////////////////////////////////////////////////////////////

// consolidateUTXOs only required when there is a new block
func (c *Client) consolidateUTXOs() {
	defer func() {
		c.wg.Done()
		c.consolidateInProgress.Store(false)
	}()

	nodeStatus, err := c.bridge.FetchNodeStatus()
	if err != nil {
		c.log.Err(err).Msg("fail to get node status")
		return
	}
	if nodeStatus != types.NodeStatus_Active {
		c.log.Info().Msgf("node is not active , doesn't need to consolidate utxos")
		return
	}
	vaults, err := c.bridge.GetAsgards()
	if err != nil {
		c.log.Err(err).Msg("fail to get current asgards")
		return
	}
	utxosToSpend := c.getMaximumUtxosToSpend()
	for _, vault := range vaults {
		if !vault.Contains(c.nodePubKey) {
			// Not part of this vault , don't need to consolidate UTXOs for this Vault
			continue
		}
		// the amount used here doesn't matter , just to see whether there are more than 15 UTXO available or not
		var utxos []btcjson.ListUnspentResult
		utxos, err = c.getUtxoToSpend(vault.PubKey, 0.01)
		if err != nil {
			c.log.Err(err).Msg("fail to get utxos to spend")
			continue
		}
		// doesn't have enough UTXOs , don't need to consolidate
		if int64(len(utxos)) < utxosToSpend {
			continue
		}
		total := 0.0
		for _, item := range utxos {
			total += item.Amount
		}
		var addr common.Address
		addr, err = vault.PubKey.GetAddress(c.cfg.ChainID)
		if err != nil {
			c.log.Err(err).Msgf("fail to get address for pubkey: %s", vault.PubKey)
			continue
		}
		// THORChain usually pay 1.5 of the last observed fee rate
		feeRate := math.Ceil(float64(c.lastFeeRate) * 3 / 2)
		var amt btcutil.Amount
		amt, err = btcutil.NewAmount(total)
		if err != nil {
			c.log.Err(err).Msgf("fail to convert to amount: %f", total)
			continue
		}

		txOutItem := stypes.TxOutItem{
			Chain:       c.cfg.ChainID,
			ToAddress:   addr,
			VaultPubKey: vault.PubKey,
			Coins: common.Coins{
				common.NewCoin(c.cfg.ChainID.GetGasAsset(), cosmos.NewUint(uint64(amt))),
			},
			Memo:    mem.NewConsolidateMemo().String(),
			MaxGas:  nil,
			GasRate: int64(feeRate),
		}
		var height int64
		height, err = c.bridge.GetBlockHeight()
		if err != nil {
			c.log.Err(err).Msg("fail to get THORChain block height")
			continue
		}
		var rawTx []byte
		rawTx, _, _, err = c.SignTx(txOutItem, height)
		if err != nil {
			c.log.Err(err).Msg("fail to sign consolidate txout item")
			continue
		}
		var txID string
		txID, err = c.BroadcastTx(txOutItem, rawTx)
		if err != nil {
			c.log.Err(err).Str("signed", string(rawTx)).Msg("fail to broadcast consolidate tx")
			continue
		}
		c.log.Info().Msgf("broadcast consolidate tx successfully, hash:%s", txID)
	}
}
