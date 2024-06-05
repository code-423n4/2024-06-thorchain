package evm

import (
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	ctypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"gitlab.com/thorchain/thornode/bifrost/metrics"
	"gitlab.com/thorchain/thornode/bifrost/pkg/chainclients/shared/evm/types"
	"gitlab.com/thorchain/thornode/bifrost/pubkeymanager"
	"gitlab.com/thorchain/thornode/bifrost/thorclient"
	ttypes "gitlab.com/thorchain/thornode/bifrost/thorclient/types"
	"gitlab.com/thorchain/thornode/cmd"
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/config"
	openapi "gitlab.com/thorchain/thornode/openapi/gen"
	types2 "gitlab.com/thorchain/thornode/x/thorchain/types"
	. "gopkg.in/check.v1"
)

type UnstuckTestSuite struct {
	thorKeys *thorclient.Keys
	bridge   thorclient.ThorchainBridge
	m        *metrics.Metrics
	server   *httptest.Server
}

var _ = Suite(&UnstuckTestSuite{})

func (s *UnstuckTestSuite) SetUpTest(c *C) {
	s.m = GetMetricForTest(c)
	c.Assert(s.m, NotNil)
	types2.SetupConfigForTest()
	c.Assert(os.Setenv("NET", "mocknet"), IsNil)

	cfg := config.BifrostClientConfiguration{
		ChainID:      "thorchain",
		SignerName:   "bob",
		SignerPasswd: "password",
	}

	kb := cKeys.NewInMemory()
	_, _, err := kb.NewMnemonic(cfg.SignerName, cKeys.English, cmd.THORChainHDPath, cfg.SignerPasswd, hd.Secp256k1)
	c.Assert(err, IsNil)
	s.thorKeys = thorclient.NewKeysWithKeybase(kb, cfg.SignerName, cfg.SignerPasswd)

	// get public key
	priv, err := s.thorKeys.GetPrivateKey()
	c.Assert(err, IsNil)
	temp, err := codec.ToTmPubKeyInterface(priv.PubKey())
	c.Assert(err, IsNil)
	pubkey, err := common.NewPubKeyFromCrypto(temp)
	c.Assert(err, IsNil)

	var vaultAddr common.Address
	vaultAddr, err = pubkey.GetAddress(common.ETHChain)
	c.Assert(err, IsNil)

	lastBroadcastTx := ""
	var lastBroadcastTxJSON []byte

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.RequestURI {
		case thorclient.ThorchainConstants:
			httpTestHandler(c, rw, "../../../../test/fixtures/endpoints/constants/constants.json")
		case thorclient.PubKeysEndpoint:
			var content []byte
			content, err = os.ReadFile("../../../../test/fixtures/endpoints/vaults/pubKeys.json")
			c.Assert(err, IsNil)
			var pubKeysVault openapi.VaultPubkeysResponse
			c.Assert(json.Unmarshal(content, &pubKeysVault), IsNil)
			var buf []byte
			buf, err = json.MarshalIndent(pubKeysVault, "", "	")
			c.Assert(err, IsNil)
			_, err = rw.Write(buf)
			c.Assert(err, IsNil)
		case thorclient.LastBlockEndpoint:
			httpTestHandler(c, rw, "../../../../test/fixtures/eth/last_block_height.json")
		case thorclient.InboundAddressesEndpoint:
			httpTestHandler(c, rw, "../../../../test/fixtures/endpoints/inbound_addresses/inbound_addresses.json")
		case thorclient.AsgardVault:
			httpTestHandler(c, rw, "../../../../test/fixtures/endpoints/vaults/asgard.json")
		case thorclient.NodeAccountEndpoint:
			httpTestHandler(c, rw, "../../../../test/fixtures/endpoints/nodeaccount/template.json")
		default:
			var body []byte
			body, err = io.ReadAll(req.Body)
			c.Assert(err, IsNil)

			type RPCRequest struct {
				JSONRPC string          `json:"jsonrpc"`
				ID      interface{}     `json:"id"`
				Method  string          `json:"method"`
				Params  json.RawMessage `json:"params"`
			}

			// the only thing in this test that uses batching is fetch transaction receipts
			if body[0] == '[' {
				var rpcRequests []RPCRequest
				err = json.Unmarshal(body, &rpcRequests)
				c.Assert(err, IsNil)
				c.Assert(len(rpcRequests), Equals, 1)
				rpcRequest := rpcRequests[0]
				c.Assert(rpcRequest.Method, Equals, "eth_getTransactionReceipt")
				var params []string
				err = json.Unmarshal(rpcRequest.Params, &params)
				c.Assert(err, IsNil)
				c.Assert(params[0], Equals, lastBroadcastTx)
				_, err = rw.Write([]byte(`[{
						"jsonrpc": "2.0",
						"id": 5,
						"result": {
							"blockHash": "0x96395fbdb39e33293999dc1a0a3b87c8a9e51185e177760d1482c2155bb35b87",
							"blockNumber": "0x1",
							"contractAddress": null,
							"cumulativeGasUsed": "0x1",
							"from": "` + vaultAddr.String() + `",
							"gasUsed": "0x1",
							"logs": [],
							"logsBloom": "0x` + strings.Repeat("00", 256) + `",
							"status": "0x1",
							"to": "` + vaultAddr.String() + `",
							"transactionHash": "` + lastBroadcastTx + `",
							"transactionIndex": "0x1"
						}
					}]`))
				c.Assert(err, IsNil)
				return
			}

			var rpcRequest RPCRequest
			err = json.Unmarshal(body, &rpcRequest)
			c.Assert(err, IsNil)

			switch rpcRequest.Method {
			case "eth_chainId":
				_, err = rw.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x539"}`))
				c.Assert(err, IsNil)
				return
			case "eth_getTransactionByHash":
				var hashes []string
				c.Assert(json.Unmarshal(rpcRequest.Params, &hashes), IsNil)
				switch hashes[0] {
				case "0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b":
					_, err = rw.Write([]byte(`{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "nonce": "0x2",
        "gasPrice": "0x1",
        "gas": "0x13990",
        "to": "0xe65e9d372f8cacc7b6dfcd4af6507851ed31bb44",
        "value": "0x22b1c8c1227a00000",
        "input": "0x1fece7b4000000000000000000000000f6da288748ec4c77642f6c5543717539b3ae001b00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000045345454400000000000000000000000000000000000000000000000000000000",
        "v": "0xa96",
        "r": "0x4fed375d064158c79dd0ee1e35cbfbe6e19ed7c0005763ca1edc10121124d1fd",
        "s": "0x56d194669c9188176ed87e96b4bd2e2b3869cdb959c1153a557e3d8d8d48c12c",
        "hash": "0x81604fe8c8df8b5e32daafa00acd06ec97281ed3056ab368cf57e2dcacd7e2d1"
    }
}`))
					c.Assert(err, IsNil)
				case "0x96395fbdb39e33293999dc1a0a3b87c8a9e51185e177760d1482c2155bb35b87":
					_, err = rw.Write([]byte(`{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "blockHash": "0x96395fbdb39e33293999dc1a0a3b87c8a9e51185e177760d1482c2155bb35b87",
        "blockNumber": "0x32",
        "from": "0xfabb9cc6ec839b1214bb11c53377a56a6ed81762",
        "gas": "0x26fca",
        "gasPrice": "0x1",
        "hash": "0xc416a0332b4346f8090818981d1b2bf491d67b22cfec44ed8ec9a897b3631db2",
        "input": "0x1fece7b40000000000000000000000008d8f3199e684c76f25eeb9c0ce922d15bf72dfa200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000384144443a4554482e4554483a7474686f7231777a3738716d726b706c726468793337747730746e766e30746b6d35707164367a64703235370000000000000000",
        "nonce": "0x0",
        "to": "0xe65e9d372f8cacc7b6dfcd4af6507851ed31bb44",
        "transactionIndex": "0x0",
        "value": "0x58d15e176280000",
        "v": "0xa96",
        "r": "0x3c5e5945cadf1429bdb3e45a74003d86d62dd4091d9664f0c8163a951fe6f10e",
        "s": "0x1b86c85a84b0a76ea5d75ff08ea667f7f1e1883245b70deb4d4c347de3585ece"
    }
}`))
					c.Assert(err, IsNil)
				}
				return
			case "eth_gasPrice":
				_, err = rw.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x1"}`))
				c.Assert(err, IsNil)
				return
			case "eth_sendRawTransaction":
				// read tx data from the request
				var params []string
				err = json.Unmarshal(rpcRequest.Params, &params)
				c.Assert(err, IsNil)
				tx := params[0]

				// decode tx
				var data []byte
				data, err = hex.DecodeString(tx[2:])
				c.Assert(err, IsNil)
				var txData ctypes.Transaction
				err = rlp.DecodeBytes(data, &txData)
				c.Assert(err, IsNil)
				lastBroadcastTxJSON, err = json.MarshalIndent(&txData, "", "  ")
				c.Assert(err, IsNil)
				lastBroadcastTx = txData.Hash().Hex()

				_, err = rw.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"}`))
				c.Assert(err, IsNil)
				return
			case "eth_call":
				if string(rpcRequest.Params) == `[{"data":"0x03b6a6730000000000000000000000009f4aab49a9cd8fc54dcb3701846f608a6f2c44da0000000000000000000000003b7fa4dd21c6f9ba3ca375217ead7cab9d6bf483","from":"0x9f4aab49a9cd8fc54dcb3701846f608a6f2c44da","to":"0xe65e9d372f8cacc7b6dfcd4af6507851ed31bb44"},"latest"]` {
					_, err = rw.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x0000000000000000000000000000000000000000000000000000000000000012"}`))
					c.Assert(err, IsNil)
				} else if string(rpcRequest.Params) == `[{"data":"0x95d89b41","from":"0x0000000000000000000000000000000000000000","to":"0x3b7fa4dd21c6f9ba3ca375217ead7cab9d6bf483"},"latest"]` {
					_, err = rw.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":"0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003544b4e0000000000000000000000000000000000000000000000000000000000"}`))
					c.Assert(err, IsNil)
				}
			case "eth_getBlockByNumber":
				_, err = rw.Write([]byte(`{
						"jsonrpc": "2.0",
						"id": 1,
						"result": {
							"difficulty": "0x2",
							"extraData": "0xd88301091a846765746888676f312e31352e36856c696e757800000000000000e86d9af8b427b780cd1e6f7cabd2f9231ccac25d313ed475351ed64ac19f21491461ed1fae732d3bbf73a5866112aec23b0ca436185685b9baee4f477a950f9400",
							"gasLimit": "0x9e0f54",
							"gasUsed": "0xabd3",
							"hash": "0xb273789207ce61a1ec0314fdb88efe6c6b554a9505a97ff3dff05aa691e220ac",
							"logsBloom": "0x00010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000040000000000000000010000200020000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000040000020000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000010000000000000000000000000000000000000000000000020000000000000",
							"miner": "0x0000000000000000000000000000000000000000",
							"mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
							"nonce": "0x0000000000000000",
							"number": "0x6b",
							"parentHash": "0xf18470c54efec284fb5ad57c0ee4afe2774d61393bd5224ac5484b39a0a07556",
							"receiptsRoot": "0x794a74d56ec50769a1400f7ae0887061b0ec3ea6702589a0b45b9102df2c9954",
							"sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
							"size": "0x30a",
							"stateRoot": "0x1c84090d7f5dc8137d6762e3d4babe10b30bf61fa827618346ae1ba8600a9629",
							"timestamp": "0x6008f03a",
							"totalDifficulty": "0xd7",
							"transactions": [` + string(lastBroadcastTxJSON) + `],
							"transactionsRoot": "0x4247bb112edbe20ee8cf406864b335f4a3aa215f65ea686c9820f056c637aca6",
							"uncles": []
						}
					}
					`))
				c.Assert(err, IsNil)
			}
		}
	}))
	s.server = server

	cfg.ChainHost = server.Listener.Addr().String()
	s.bridge, err = thorclient.NewThorchainBridge(cfg, s.m, s.thorKeys)
	c.Assert(err, IsNil)
}

func (s *UnstuckTestSuite) TearDownTest(c *C) {
	c.Assert(os.Unsetenv("NET"), IsNil)
}

func (s *UnstuckTestSuite) TestUnstuckProcess(c *C) {
	config.Init()
	pubkeyMgr, err := pubkeymanager.NewPubKeyManager(s.bridge, s.m)
	c.Assert(err, IsNil)
	poolMgr := thorclient.NewPoolMgr(s.bridge)
	e, err := NewEVMClient(s.thorKeys, config.BifrostChainConfiguration{
		ChainID:        common.AVAXChain,
		RPCHost:        "http://" + s.server.Listener.Addr().String(),
		SolvencyBlocks: 10,
		BlockScanner: config.BifrostBlockScannerConfiguration{
			StartBlockHeight:   1, // avoids querying thorchain for block height
			HTTPRequestTimeout: time.Second * 10,
			GasCacheBlocks:     40,
			Concurrency:        1,
		},
	}, nil, s.bridge, s.m, pubkeyMgr, poolMgr)
	c.Assert(err, IsNil)
	c.Assert(e, NotNil)
	c.Assert(pubkeyMgr.Start(), IsNil)
	defer func() { c.Assert(pubkeyMgr.Stop(), IsNil) }()
	pubkey := e.kw.GetPubKey().String()

	txID1 := types2.GetRandomTxHash().String()
	txID2 := types2.GetRandomTxHash().String()
	// add some thing here
	c.Assert(e.evmScanner.blockMetaAccessor.AddSignedTxItem(types.SignedTxItem{
		Hash:        txID1,
		Height:      1022,
		VaultPubKey: pubkey,
	}), IsNil)
	c.Assert(e.evmScanner.blockMetaAccessor.AddSignedTxItem(types.SignedTxItem{
		Hash:        txID2,
		Height:      1024,
		VaultPubKey: pubkey,
	}), IsNil)
	// this should not do anything , because because all the tx has not been
	e.unstuckAction()
	items, err := e.evmScanner.blockMetaAccessor.GetSignedTxItems()
	c.Assert(err, IsNil)
	c.Assert(items, HasLen, 2)
	c.Assert(e.evmScanner.blockMetaAccessor.RemoveSignedTxItem(txID1), IsNil)
	c.Assert(e.evmScanner.blockMetaAccessor.RemoveSignedTxItem(txID2), IsNil)

	inHash := types2.GetRandomTxHash()
	pKey := types2.GetRandomPubKey()
	toAddr, _ := pKey.GetAddress(e.GetChain())
	toi := ttypes.TxOutItem{
		Chain:       e.GetChain(),
		ToAddress:   toAddr,
		VaultPubKey: e.kw.GetPubKey(),
		Coins:       common.Coins{common.NewCoin(e.GetChain().GetGasAsset(), cosmos.NewUint(100000000))},
		Memo:        "OUT:" + inHash.String(),
		MaxGas:      common.Gas(common.NewCoins(common.NewCoin(e.GetChain().GetGasAsset(), cosmos.NewUint(100000000)))),
		GasRate:     5,
		InHash:      inHash,
		Height:      800,
	}
	stuckTxID := "0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b"
	err = e.signerCacheManager.SetSigned(toi.CacheHash(), toi.CacheVault(e.GetChain()), stuckTxID)
	c.Assert(err, IsNil)

	c.Assert(e.evmScanner.blockMetaAccessor.AddSignedTxItem(types.SignedTxItem{
		Hash:        stuckTxID,
		Height:      800,
		VaultPubKey: pubkey,
		TxOutItem:   &toi,
	}), IsNil)
	c.Assert(e.evmScanner.blockMetaAccessor.AddSignedTxItem(types.SignedTxItem{
		Hash:        "0x96395fbdb39e33293999dc1a0a3b87c8a9e51185e177760d1482c2155bb35b87",
		Height:      800,
		VaultPubKey: pubkey,
	}), IsNil)
	// this should try to check 0x88df016429689c079f3b2f6ad39fa052532c56795b733da78a91ebe6a713944b
	e.unstuckAction()
	items, err = e.evmScanner.blockMetaAccessor.GetSignedTxItems()
	c.Assert(err, IsNil)

	// the stuck tx should be removed
	c.Assert(items, HasLen, 0)

	// signer cache item should not have been removed yet
	c.Assert(e.signerCacheManager.HasSigned(toi.CacheHash()), Equals, true)

	// removal should occur after the block containing the cancel transaction is scanned
	txIn, err := e.evmScanner.FetchTxs(int64(1), int64(1))
	c.Assert(err, IsNil)
	c.Check(len(txIn.TxArray), Equals, 1)

	// after one block the goroutine should remove from signer cache
	c.Assert(e.signerCacheManager.HasSigned(toi.CacheHash()), Equals, false)
}
