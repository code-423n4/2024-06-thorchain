package utxo

import (
	"bytes"
	"fmt"
	"sort"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcjson"
	btcchaincfg "github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	btcwire "github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	btctxscript "gitlab.com/thorchain/bifrost/txscript"

	dogeec "github.com/eager7/dogd/btcec"
	dogechaincfg "github.com/eager7/dogd/chaincfg"
	dogewire "github.com/eager7/dogd/wire"
	"github.com/eager7/dogutil"
	dogetxscript "gitlab.com/thorchain/bifrost/dogd-txscript"

	"github.com/gcash/bchd/bchec"
	bchchaincfg "github.com/gcash/bchd/chaincfg"
	bchwire "github.com/gcash/bchd/wire"
	"github.com/gcash/bchutil"
	bchtxscript "gitlab.com/thorchain/bifrost/bchd-txscript"

	ltcec "github.com/ltcsuite/ltcd/btcec"
	ltcchaincfg "github.com/ltcsuite/ltcd/chaincfg"
	ltcwire "github.com/ltcsuite/ltcd/wire"
	"github.com/ltcsuite/ltcutil"
	ltctxscript "gitlab.com/thorchain/bifrost/ltcd-txscript"

	"gitlab.com/thorchain/thornode/bifrost/pkg/chainclients/utxo/rpc"
	"gitlab.com/thorchain/thornode/bifrost/thorclient"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"

	. "gitlab.com/thorchain/thornode/test/simulation/pkg/types"
)

////////////////////////////////////////////////////////////////////////////////////////
// Client
////////////////////////////////////////////////////////////////////////////////////////

//go:generate go run utxo_generate.go

type Client struct {
	sync.Mutex
	chain common.Chain
	rpc   *rpc.Client

	keys    *thorclient.Keys
	privKey *btcec.PrivateKey
	pubKey  common.PubKey
	address common.Address
}

var _ LiteChainClient = &Client{}

func NewConstructor(host string) LiteChainClientConstructor {
	return func(chain common.Chain, keys *thorclient.Keys) (LiteChainClient, error) {
		return NewClient(chain, host, keys)
	}
}

func NewClient(chain common.Chain, host string, keys *thorclient.Keys) (LiteChainClient, error) {
	// create rpc client
	retries := 5
	rpc, err := rpc.NewClient(host, "thorchain", "password", retries, zerolog.Nop())
	if err != nil {
		return nil, fmt.Errorf("fail to create rpc client: %w", err)
	}

	// extract the private key
	privateKey, err := keys.GetPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("fail to get private key: %w", err)
	}
	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKey.Bytes())

	// derive the public key
	buf := privKey.PubKey().SerializeCompressed()
	pk := secp256k1.PubKey(buf)
	pubkey, err := common.NewPubKeyFromCrypto(pk)
	if err != nil {
		return nil, fmt.Errorf("fail to create pubkey: %w", err)
	}

	// get pubkey address for the chain
	address, err := pubkey.GetAddress(chain)
	if err != nil {
		return nil, fmt.Errorf("fail to get address from pubkey(%s): %w", pk, err)
	}

	// import address to wallet, rescan if master account
	if keys.GetSignerInfo().GetName() == "master" {
		err = rpc.ImportAddressRescan(address.String())
	} else {
		err = rpc.ImportAddress(address.String())
	}
	if err != nil {
		return nil, fmt.Errorf("fail to import address(%s): %w", address, err)
	}

	return &Client{
		chain:   chain,
		rpc:     rpc,
		keys:    keys,
		privKey: privKey,
		pubKey:  pubkey,
		address: address,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// GetAccount
////////////////////////////////////////////////////////////////////////////////////////

func (c *Client) GetAccount(pk *common.PubKey) (*common.Account, error) {
	// default to the client key address
	var err error
	addr := c.address
	if pk != nil {
		addr, err = pk.GetAddress(c.chain)
		if err != nil {
			return nil, fmt.Errorf("fail to get address from pubkey(%s): %w", pk, err)
		}
	}

	// get unspent utxos
	utxos, err := c.rpc.ListUnspent(addr.String())
	if err != nil {
		return nil, fmt.Errorf("fail to get UTXOs: %w", err)
	}

	// sum balance
	total := 0.0
	for _, item := range utxos {
		total += item.Amount
	}
	totalAmt, err := btcutil.NewAmount(total)
	if err != nil {
		return nil, fmt.Errorf("fail to convert total amount: %w", err)
	}

	// create account
	coin := common.NewCoin(c.chain.GetGasAsset(), cosmos.NewUint(uint64(totalAmt)))
	a := common.NewAccount(0, 0, common.NewCoins(coin), false)

	return &a, nil
}

////////////////////////////////////////////////////////////////////////////////////////
// SignTx
////////////////////////////////////////////////////////////////////////////////////////

func (c *Client) SignContractTx(SimContractTx) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

func (c *Client) SignTx(tx SimTx) ([]byte, error) {
	sourceScript, err := c.getSourceScript(tx)
	if err != nil {
		return nil, fmt.Errorf("fail to get source pay to address script: %w", err)
	}

	// build transaction
	redeemTx, amounts, err := c.buildTx(tx, sourceScript)
	if err != nil {
		return nil, err
	}

	// create the list of signing requests
	signings := []struct{ idx, amount int64 }{}
	totalAmount := int64(0)
	for idx, txIn := range redeemTx.TxIn {
		key := fmt.Sprintf("%s-%d", txIn.PreviousOutPoint.Hash, txIn.PreviousOutPoint.Index)
		outputAmount := amounts[key]
		totalAmount += outputAmount
		signings = append(signings, struct{ idx, amount int64 }{int64(idx), outputAmount})
	}

	// convert the wire tx to the chain specific tx for signing
	var stx interface{}
	switch c.chain {
	case common.DOGEChain:
		stx = wireToDOGE(redeemTx)
	case common.BCHChain:
		stx = wireToBCH(redeemTx)
	case common.LTCChain:
		stx = wireToLTC(redeemTx)
	case common.BTCChain:
		stx = wireToBTC(redeemTx)
	default:
		log.Fatal().Msg("unsupported chain")
	}

	// sign the tx
	wg := &sync.WaitGroup{}
	wg.Add(len(signings))
	mu := &sync.Mutex{}
	var utxoErr error
	for _, signing := range signings {
		go func(i int, amount int64) {
			defer wg.Done()

			// trunk-ignore(golangci-lint/govet): shadow
			var err error

			// trunk-ignore-all(golangci-lint/forcetypeassert)
			switch c.chain {
			case common.DOGEChain:
				err = c.signUTXODOGE(stx.(*dogewire.MsgTx), amount, sourceScript, i)
			case common.BCHChain:
				err = c.signUTXOBCH(stx.(*bchwire.MsgTx), amount, sourceScript, i)
			case common.LTCChain:
				err = c.signUTXOLTC(stx.(*ltcwire.MsgTx), amount, sourceScript, i)
			case common.BTCChain:
				err = c.signUTXOBTC(stx.(*btcwire.MsgTx), amount, sourceScript, i)
			default:
				log.Fatal().Msg("unsupported chain")
			}

			if err != nil {
				mu.Lock()
				utxoErr = multierror.Append(utxoErr, err)
				mu.Unlock()
			}
		}(int(signing.idx), signing.amount)
	}
	wg.Wait()
	if utxoErr != nil {
		return nil, fmt.Errorf("fail to sign the message: %w", err)
	}

	// convert back to wire tx
	switch c.chain {
	case common.DOGEChain:
		redeemTx = dogeToWire(stx.(*dogewire.MsgTx))
	case common.BCHChain:
		redeemTx = bchToWire(stx.(*bchwire.MsgTx))
	case common.LTCChain:
		redeemTx = ltcToWire(stx.(*ltcwire.MsgTx))
	case common.BTCChain:
		redeemTx = btcToWire(stx.(*btcwire.MsgTx))
	default:
		log.Fatal().Msg("unsupported chain")
	}

	// calculate the final transaction size
	var signedTx bytes.Buffer
	if err = redeemTx.Serialize(&signedTx); err != nil {
		return nil, fmt.Errorf("fail to serialize tx to bytes: %w", err)
	}

	return signedTx.Bytes(), nil
}

////////////////////////////////////////////////////////////////////////////////////////
// BroadcastTx
////////////////////////////////////////////////////////////////////////////////////////

func (c *Client) BroadcastTx(payload []byte) (string, error) {
	redeemTx := btcwire.NewMsgTx(btcwire.TxVersion)
	buf := bytes.NewBuffer(payload)
	if err := redeemTx.Deserialize(buf); err != nil {
		return "", fmt.Errorf("fail to deserialize payload: %w", err)
	}

	var maxFee any
	switch c.chain {
	case common.DOGEChain, common.BCHChain:
		maxFee = true // "allowHighFees"
	case common.LTCChain, common.BTCChain:
		maxFee = 10_000_000
	}

	// broadcast tx
	return c.rpc.SendRawTransaction(redeemTx, maxFee)
}

////////////////////////////////////////////////////////////////////////////////////////
// Internal
////////////////////////////////////////////////////////////////////////////////////////

func (c *Client) getSourceScript(tx SimTx) ([]byte, error) {
	switch c.chain {
	case common.DOGEChain:
		var addr dogutil.Address
		addr, err := dogutil.DecodeAddress(c.address.String(), c.getChainCfgDOGE())
		if err != nil {
			return nil, fmt.Errorf("fail to decode source address(%s): %w", c.address, err)
		}
		return dogetxscript.PayToAddrScript(addr)
	case common.BCHChain:
		var addr bchutil.Address
		addr, err := bchutil.DecodeAddress(c.address.String(), c.getChainCfgBCH())
		if err != nil {
			return nil, fmt.Errorf("fail to decode source address(%s): %w", c.address, err)
		}
		return bchtxscript.PayToAddrScript(addr)
	case common.LTCChain:
		var addr ltcutil.Address
		addr, err := ltcutil.DecodeAddress(c.address.String(), c.getChainCfgLTC())
		if err != nil {
			return nil, fmt.Errorf("fail to decode source address(%s): %w", c.address, err)
		}
		return ltctxscript.PayToAddrScript(addr)
	case common.BTCChain:
		var addr btcutil.Address
		addr, err := btcutil.DecodeAddress(c.address.String(), c.getChainCfgBTC())
		if err != nil {
			return nil, fmt.Errorf("fail to decode source address(%s): %w", c.address, err)
		}
		return btctxscript.PayToAddrScript(addr)
	default:
		log.Fatal().Msg("unsupported chain")
		return nil, nil
	}
}

func (c *Client) getUtxoToSpend(total float64) ([]btcjson.ListUnspentResult, error) {
	utxos, err := c.rpc.ListUnspent(c.address.String())
	if err != nil {
		return nil, fmt.Errorf("fail to get UTXOs: %w", err)
	}

	// spend UTXOs older to younger
	sort.SliceStable(utxos, func(i, j int) bool {
		if utxos[i].Confirmations > utxos[j].Confirmations {
			return true
		} else if utxos[i].Confirmations < utxos[j].Confirmations {
			return false
		}
		return utxos[i].TxID < utxos[j].TxID
	})

	// collect UTXOs to spend
	var result []btcjson.ListUnspentResult
	var toSpend float64
	for _, item := range utxos {
		result = append(result, item)
		toSpend += item.Amount
		if toSpend >= total {
			break
		}
	}

	// error if there is insufficient balance
	if toSpend < total {
		return nil, fmt.Errorf("insufficient balance: %f < %f", toSpend, total)
	}

	return result, nil
}

func (c *Client) getGasSats() float64 {
	switch c.chain {
	case common.DOGEChain:
		return 0.1
	case common.BCHChain:
		return 0.0001
	case common.LTCChain:
		return 0.0001
	case common.BTCChain:
		return 0.0001
	default:
		log.Fatal().Msg("unsupported chain")
		return 0.0
	}
}

func (c *Client) getPaymentAmount(tx SimTx) float64 {
	amtToPay1e8 := tx.Coin.Amount.Uint64()
	amtToPay := btcutil.Amount(int64(amtToPay1e8)).ToBTC()
	amtToPay += c.getGasSats()
	return amtToPay
}

func (c *Client) buildTx(tx SimTx, sourceScript []byte) (
	*btcwire.MsgTx, map[string]int64, error,
) {
	txes, err := c.getUtxoToSpend(c.getPaymentAmount(tx))
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
	switch c.chain {
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
		log.Fatal().Msg("unsupported chain")
	}

	// pay to customer
	redeemTxOut := wire.NewTxOut(int64(tx.Coin.Amount.Uint64()), buf)
	redeemTx.AddTxOut(redeemTxOut)

	// add output to pay the balance back ourselves
	balance := totalAmt - redeemTxOut.Value - int64(c.getGasSats()*common.One)
	if balance > 0 {
		redeemTx.AddTxOut(wire.NewTxOut(balance, sourceScript))
	}

	// memo
	var nullDataScript []byte
	switch c.chain {
	case common.DOGEChain:
		nullDataScript, err = dogetxscript.NullDataScript([]byte(tx.Memo))
	case common.BCHChain:
		nullDataScript, err = bchtxscript.NullDataScript([]byte(tx.Memo))
	case common.LTCChain:
		nullDataScript, err = ltctxscript.NullDataScript([]byte(tx.Memo))
	case common.BTCChain:
		nullDataScript, err = btctxscript.NullDataScript([]byte(tx.Memo))
	default:
		log.Fatal().Msg("unsupported chain")
	}
	if err != nil {
		return nil, nil, fmt.Errorf("fail to generate null data script: %w", err)
	}
	redeemTx.AddTxOut(wire.NewTxOut(0, nullDataScript))

	return redeemTx, individualAmounts, nil
}

// ------------------------------ chain config ------------------------------

func (c *Client) getChainCfgBTC() *btcchaincfg.Params {
	switch common.CurrentChainNetwork {
	case common.MockNet:
		return &btcchaincfg.RegressionNetParams
	case common.TestNet:
		return &btcchaincfg.TestNet3Params
	case common.MainNet:
		return &btcchaincfg.MainNetParams
	case common.StageNet:
		return &btcchaincfg.MainNetParams
	default:
		log.Fatal().Msg("unsupported network")
		return nil
	}
}

func (c *Client) getChainCfgLTC() *ltcchaincfg.Params {
	cn := common.CurrentChainNetwork
	switch cn {
	case common.MockNet:
		return &ltcchaincfg.RegressionNetParams
	case common.MainNet:
		return &ltcchaincfg.MainNetParams
	case common.StageNet:
		return &ltcchaincfg.MainNetParams
	}
	return nil
}

func (c *Client) getChainCfgBCH() *bchchaincfg.Params {
	switch common.CurrentChainNetwork {
	case common.MockNet:
		return &bchchaincfg.RegressionNetParams
	case common.MainNet:
		return &bchchaincfg.MainNetParams
	case common.StageNet:
		return &bchchaincfg.MainNetParams
	default:
		log.Fatal().Msg("unsupported network")
		return nil
	}
}

func (c *Client) getChainCfgDOGE() *dogechaincfg.Params {
	switch common.CurrentChainNetwork {
	case common.MockNet:
		return &dogechaincfg.RegressionNetParams
	case common.MainNet:
		return &dogechaincfg.MainNetParams
	case common.StageNet:
		return &dogechaincfg.MainNetParams
	default:
		log.Fatal().Msg("unsupported network")
		return nil
	}
}

// ------------------------------ signing ------------------------------

func (c *Client) signUTXODOGE(redeemTx *dogewire.MsgTx, amount int64, sourceScript []byte, idx int) error {
	signable := dogetxscript.NewPrivateKeySignable((*dogeec.PrivateKey)(c.privKey))
	sig, err := dogetxscript.RawTxInSignature(redeemTx, idx, sourceScript, dogetxscript.SigHashAll, signable)
	if err != nil {
		return fmt.Errorf("fail to get witness: %w", err)
	}

	pkData := signable.GetPubKey().SerializeCompressed()
	sigscript, err := dogetxscript.NewScriptBuilder().AddData(sig).AddData(pkData).Script()
	if err != nil {
		return fmt.Errorf("fail to build signature script: %w", err)
	}
	redeemTx.TxIn[idx].SignatureScript = sigscript
	flag := dogetxscript.StandardVerifyFlags
	engine, err := dogetxscript.NewEngine(sourceScript, redeemTx, idx, flag, nil, nil, amount)
	if err != nil {
		return fmt.Errorf("fail to create engine: %w", err)
	}
	if err = engine.Execute(); err != nil {
		return fmt.Errorf("fail to execute the script: %w", err)
	}
	return nil
}

func (c *Client) signUTXOBCH(redeemTx *bchwire.MsgTx, amount int64, sourceScript []byte, idx int) error {
	signable := bchtxscript.NewPrivateKeySignable((*bchec.PrivateKey)(c.privKey))
	sig, err := bchtxscript.RawTxInECDSASignature(redeemTx, idx, sourceScript, bchtxscript.SigHashAll, signable, amount)
	if err != nil {
		return fmt.Errorf("fail to get witness: %w", err)
	}

	pkData := signable.GetPubKey().SerializeCompressed()
	sigscript, err := bchtxscript.NewScriptBuilder().AddData(sig).AddData(pkData).Script()
	if err != nil {
		return fmt.Errorf("fail to build signature script: %w", err)
	}
	redeemTx.TxIn[idx].SignatureScript = sigscript
	flag := bchtxscript.StandardVerifyFlags
	engine, err := bchtxscript.NewEngine(sourceScript, redeemTx, idx, flag, nil, nil, amount)
	if err != nil {
		return fmt.Errorf("fail to create engine: %w", err)
	}
	if err = engine.Execute(); err != nil {
		return fmt.Errorf("fail to execute the script: %w", err)
	}
	return nil
}

func (c *Client) signUTXOLTC(redeemTx *ltcwire.MsgTx, amount int64, sourceScript []byte, idx int) error {
	sigHashes := ltctxscript.NewTxSigHashes(redeemTx)
	signable := ltctxscript.NewPrivateKeySignable((*ltcec.PrivateKey)(c.privKey))
	witness, err := ltctxscript.WitnessSignature(redeemTx, sigHashes, idx, amount, sourceScript, ltctxscript.SigHashAll, signable, true)
	if err != nil {
		return fmt.Errorf("fail to get witness: %w", err)
	}

	redeemTx.TxIn[idx].Witness = witness
	flag := ltctxscript.StandardVerifyFlags
	engine, err := ltctxscript.NewEngine(sourceScript, redeemTx, idx, flag, nil, nil, amount)
	if err != nil {
		return fmt.Errorf("fail to create engine: %w", err)
	}
	if err = engine.Execute(); err != nil {
		return fmt.Errorf("fail to execute the script: %w", err)
	}
	return nil
}

func (c *Client) signUTXOBTC(redeemTx *btcwire.MsgTx, amount int64, sourceScript []byte, idx int) error {
	sigHashes := btctxscript.NewTxSigHashes(redeemTx)
	signable := btctxscript.NewPrivateKeySignable(c.privKey)
	witness, err := btctxscript.WitnessSignature(redeemTx, sigHashes, idx, amount, sourceScript, btctxscript.SigHashAll, signable, true)
	if err != nil {
		return fmt.Errorf("fail to get witness: %w", err)
	}

	redeemTx.TxIn[idx].Witness = witness
	flag := btctxscript.StandardVerifyFlags
	engine, err := btctxscript.NewEngine(sourceScript, redeemTx, idx, flag, nil, nil, amount)
	if err != nil {
		return fmt.Errorf("fail to create engine: %w", err)
	}
	if err = engine.Execute(); err != nil {
		return fmt.Errorf("fail to execute the script: %w", err)
	}
	return nil
}
