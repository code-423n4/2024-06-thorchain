# UTXO

## Chain Client

Example for Bitcoin.

### Observer

{{#embed https://gitlab.com/thorchain/thornode/-/blob/develop/bifrost/pkg/chainclients/bitcoin/client.go }}

#### Signer

{{#embed https://gitlab.com/thorchain/thornode/-/blob/develop/bifrost/pkg/chainclients/bitcoin/signer.go }}

### Scanning Blocks

The block scanner monitors the Asgard Addresses and looks for incoming UTXOs spending to those addresses. When it sees one performs validation on it and witnesses to THORChain. For Bitcoin, it looks that at least 1 output is spent to Asgard, and searches for another output to have an `OP_RETURN`. These two outputs form the `amount` and `memo` witness to THORChain.

[https://gitlab.com/thorchain/thornode/-/blob/develop/bifrost/blockscanner/blockscanner.go](https://gitlab.com/thorchain/thornode/-/blob/develop/bifrost/blockscanner/blockscanner.go)

### Confirmation Counting

The `txValue` is the sum of all transactions received in a block to Asgard vaults. The `blockValue` is the coinbase value, which includes fees and subsidy. If a miner forgets to add a coinbase value (it has happened) a default of 6.25 is used. (This should be updated every 4 years, or use logic to auto-update).

[https://gitlab.com/thorchain/thornode/-/blob/4bcaf4f80787d0aaee711388578ed453959ef673/bifrost/pkg/chainclients/bitcoin/client.go#L1084](https://gitlab.com/thorchain/thornode/-/blob/4bcaf4f80787d0aaee711388578ed453959ef673/bifrost/pkg/chainclients/bitcoin/client.go#L1084)

### Re-orgs

Bifrost tracks the `BlockCacheSize = 144` blocks of transactions reported in a local KV store. Every time it detects a new block at a previous height it has seen, it checks for the presence of every transaction it has reported. If the transaction is missing then it has been re-orged out. The missing txID is reported to THORChain as an `ErrataTx`

[https://gitlab.com/thorchain/thornode/-/blob/4bcaf4f80787d0aaee711388578ed453959ef673/bifrost/pkg/chainclients/bitcoin/client.go#L337](https://gitlab.com/thorchain/thornode/-/blob/4bcaf4f80787d0aaee711388578ed453959ef673/bifrost/pkg/chainclients/bitcoin/client.go#L337)

### Network Fees

Reported as `sats/byte` where the fee rate is computed over the last block. Reports the highest seen in the last 20 blocks. [https://gitlab.com/thorchain/thornode/-/blob/4bcaf4f80787d0aaee711388578ed453959ef673/bifrost/pkg/chainclients/bitcoin/client.go#L668](https://gitlab.com/thorchain/thornode/-/blob/4bcaf4f80787d0aaee711388578ed453959ef673/bifrost/pkg/chainclients/bitcoin/client.go#L668)

### Handling Gas

The gas amount for a transaction is just the difference between outputs and inputs.\
[https://gitlab.com/thorchain/thornode/-/blob/4bcaf4f80787d0aaee711388578ed453959ef673/bifrost/pkg/chainclients/bitcoin/client.go#L988](https://gitlab.com/thorchain/thornode/-/blob/4bcaf4f80787d0aaee711388578ed453959ef673/bifrost/pkg/chainclients/bitcoin/client.go#L988)

### Other Considerations

#### UTXO consolidation

UTXOs consume inputs, and these inputs need to be signed independently. Thus consuming 15 inputs requires 15 times the TSS bandwidth than a single input. To prevent runaway liabilities the client will automatically enter a TSS signing ceremony for Asgard every 15 inputs to consolidate them back to one. This transaction uses the `consolidate` memo and can be seen regularly on THORChain vaults.

[https://gitlab.com/thorchain/thornode/-/blob/4bcaf4f80787d0aaee711388578ed453959ef673/bifrost/pkg/chainclients/bitcoin/client.go#L915](https://gitlab.com/thorchain/thornode/-/blob/4bcaf4f80787d0aaee711388578ed453959ef673/bifrost/pkg/chainclients/bitcoin/client.go#L915)

#### ChildPaysForParent

Asgard cannot consume a pending transaction spent to it, since THORChain requires at least 1 confirmation. However, Ygg Vaults will consume pending transactions, since they continually spend back to themselves and are only funded by Asgard. To do this, outbound transactions from Ygg Vaults are actually witnessed when in the mempool, instead of being confirmed. This allows Ygg vaults to have high swap throughput, even if the swaps are still pending in the mempool.

```admonish warning
Ygg vaults have historically been subject to dust attacks which spend large-size transactions with low fees, causing vaults to lock up. To prevent this, Ygg vaults only consume pending transactions spent to itself.
```

#### ReplaceByFee

RBF transactions allow the spender to double-spend with a higher fee. Users can use RBF transactions to spend to Asgard, but RBF does not need to be used in the THORChain vaults. Ygg vaults have CPFP instead.

## Wallet Client

UTXO clients implemented in XChainJS have the following nuances:

### Fees

The wallet client should spend with a fee rate at least equal to what is reported on `inbound_addresses` - if not it risks not being confirmed by the time the vault migrates.

### Pending UTXOS

Do not consume pending transactions when spending to Asgard (with a memo) since it may consume a low-fee tx and get stuck.

### MEMO

The memo is inserted as an OP_RETURN in an output. It can be any output. The MEMO is limited to 80bytes, so it should be trimmed and use abbreviated memos or Asset identifiers where possible.
