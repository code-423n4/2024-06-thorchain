# Fees and Wait Times

## **Fees**

Users pay two kinds of fees when entering or exiting Savings Vaults:

1. **Layer1 Network Fees** (gas): paid by the user when depositing or paid by the network when withdrawing and subtracted from the user's redemption value.
2. **Slip Fees**: protects the pool from being manipulated by large deposits/withdraws. Calculated as a function of transaction size and current pool depth.

The following are required to determine approximate deposit / withdrawal fees:

```json
outboundFee = curl -SL https://thornode.ninerealms.com/thorchain/inbound_addresses | jq '.[] | select(.chain == "BTC") | .outbound_fee'
=> 30000

poolDepth = curl -SL https://thornode.ninerealms.com/thorchain/pools | jq '.[] | select(.asset == "BTC.BTC") | .balance_asset'
=> 68352710830 => 683.5 BTC
```

```admonish info
The Quote endpoints will return fee estimates.
```

### Deposit Fees

_Example:_ user is depositing 1.0 BTC into the network, which has 1000 BTC in the pool, with 30k sats `outboundFee.`

The user will pay \~1/3rd of the THORChain's outbound fee to send assets to Savings Vault, using their typical wallet fee settings (note, this is an estimate only).

```go
totalFee = networkFee + liquidityFee

networkFee = 0.33 * outboundFee = 10,000 sats

liquidityFee = depositAmount / (depositAmount + poolDepth) * depositAmount
liquidityFee = 1.0 / (1.0+10000) * 1.0 = 99000 sats

total fee = 109,000 sats
```

### Withdrawal Fees

Example: user is withdrawing 1.1 BTC from the network, which has 1000 BTC in the pool, with 30k `outboundFee.`

```go
totalFee = networkFee + liquidityFee

networkFee = outboundFee = 30,000 sats

liquidityFee = withdrawAmount / (withdrawAmount + poolDepth) * withdrawAmount
liquidityFee = 1.1 / (1.1 + 1001.1) * 1.1 = 120,734 sats

total fee = 150,734 sats
```

```admonish info
Remember, the **liquidityFee** is entirely dependent on the size of the transaction the user is wishing to do. They may wish to do smaller transactions over a period of time to reduce fees.
```

## Wait Times

When **depositing**, there are three phases to the transaction.

1. **Layer1 Inbound Confirmation -** assuming the inbound Tx will be confirmed in the next block, it is the source blockchain block time.
2. **Observation Counting** - time for 67% THORChain Nodes to observe and agree on the inbound Tx.
3. **Confirmation Counting** - for non-instant finality blockchains, the amount of time THORChain will wait before processing to protect against double spends and re-org attacks.

When **withdrawing** using the dust threshold, there are three phases to the transaction

1. **Layer1 Inbound Confirmation -** assuming the inbound Tx will be confirmed in the next block, it is the source blockchain block time.
2. **Observation Counting** - time for 67% THORChain Nodes to observe and agree on the inbound Tx.
3. **Outbound Delay** - dependent on size and network traffic. Large outbounds will be delayed.
4. **Layer1 Outbound Confirmation** - Outbound blockchain block time.

Wait times can be between a few seconds up to an hour. The assets being swapped, the size of the swap and the current network traffic within THORChain will determine the wait time

```admonish info
The Quote endpoint will calculate wait times.
```

See the [delays.md](../concepts/delays.md "mention") section for full details.
