package utxo

import (
	"github.com/btcsuite/btcd/btcec"
	btcjson "github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"gitlab.com/thorchain/thornode/common"
)

func bech32AccountPubKey(key *btcec.PrivateKey) (common.PubKey, error) {
	buf := key.PubKey().SerializeCompressed()
	pk := secp256k1.PubKey(buf)
	return common.NewPubKeyFromCrypto(pk)
}

func sumVoutSats(tx *btcjson.TxRawResult) (uint64, error) {
	var sumVout uint64 = 0
	for _, vout := range tx.Vout {
		amount, err := btcutil.NewAmount(vout.Value)
		if err != nil {
			return 0, err
		}
		sumVout += uint64(amount.ToUnit(btcutil.AmountSatoshi))
	}
	return sumVout, nil
}
