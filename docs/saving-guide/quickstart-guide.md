# Quickstart Guide

## Introduction

THORChain allows users to deposit Layer1 assets into its network to earn asset-denominated yield without RUNE asset exposure, or being aware of THORChain’s network.

There is no permission, authentication or prior steps, so developers can get started and allow their users to earn asset-denominated yield simply by sending layer1 transactions to THORChain vaults.

Under the hood, THORChain deposits the user’s Layer1 asset into a liquidity pool which earns yield. This yield is tracked and paid to the user’s deposit value. Users can withdraw their Layer1 asset, including the yield earned. There is no slashing, penalties, timelocks, or account minimum/maximums. The only fees paid are the Layer1 fees to make a deposit and withdraw transaction (as necessitated), and a slip-based fee on entry and exit to stop price manipulation attacks. Both of these are transparent and within the user’s control.

```admonish indo
[Streaming swaps](../swap-guide/streaming-swaps.md) is enabled for savers.
```

## Quote for a Savers Quote

Savers Quote endpoints have been created to simplify the implementation process.

**Add 1 BTC to Savers.**

**Request:** _Add 1_ _BTC to Savers_

[https://thornode.ninerealms.com/thorchain/quote/saver/deposit?asset=BTC.BTC\&amount=100000000](https://thornode.ninerealms.com/thorchain/quote/saver/deposit?asset=BTC.BTC&amount=100000000)

**Response:**

```json
{
  "dust_threshold": "10000",
  "expected_amount_deposit": "99932291",
  "expected_amount_out": "99932291",
  "expiry": 1700263119,
  "fees": {
    "affiliate": "0",
    "asset": "BTC/BTC",
    "liquidity": "67672",
    "outbound": "355",
    "slippage_bps": 6,
    "total": "68027",
    "total_bps": 6
  },
  "inbound_address": "bc1qe7lfmet2l5j7ypsd6ln300jt8mg3dt2q3darj8",
  "inbound_confirmation_blocks": 1,
  "inbound_confirmation_seconds": 600,
  "memo": "+:BTC/BTC",
  "notes": "First output should be to inbound_address, second output should be change back to self, third output should be OP_RETURN, limited to 80 bytes. Do not send below the dust threshold. Do not use exotic spend scripts, locks or address formats (P2WSH with Bech32 address format preferred).",
  "recommended_min_amount_in": "10000",
  "slippage_bps": 13,
  "warning": "Do not cache this response. Do not send funds after the expiry."
}
```

_If you send 1 BTC to_ bc1quuf5sr444km2zlgrg654mjdfgkuzayfs7nqrfm*with the memo* `+:BTC/BTC`_, you can expect `0.99932` BTC will and will incur 13 basis points (0.13%) of slippage._

```admonish danger
The `Inbound_Address` changes regularly, do not cache!
```

```admonish danger
Inbound transactions should not be delayed for any reason else there is risk funds will be sent to an unreachable address. Use standard transactions, check the `inbound address` before sending and use the recommended [`gas rate`](../concepts/querying-thorchain.md#getting-the-asgard-vault) to ensure transactions are confirmed in the next block to the latest `Inbound_Address`.
```

_For security reasons, your inbound transaction will be delayed by 1 BTC Block._

```admonish info
Full quote saving endpoint specification can be found here: [https://thornode.ninerealms.com/thorchain/doc/](https://thornode.ninerealms.com/thorchain/doc/).

See an example implementation [here](https://replit.com/@thorchain/quoteSavers#index.js).
```

**User withdrawing all of their BTC Saver's position.**

**Request:** _Withdraw 100% of BTC Savers for_ `bc1qy9rjlz5w3tqn7m3reh3y48n8del4y8z42sswx5`

[https://thornode.ninerealms.com/thorchain/quote/saver/withdraw?asset=BTC.BTC\&address=bc1qy9rjlz5w3tqn7m3reh3y48n8del4y8z42sswx5\&withdraw_bps=10000](https://thornode.ninerealms.com/thorchain/quote/saver/withdraw?asset=BTC.BTC&address=bc1qy9rjlz5w3tqn7m3reh3y48n8del4y8z42sswx5&withdraw_bps=10000)

**Response:**

```json
{
  "dust_amount": "20000",
  "dust_threshold": "10000",
  "expected_amount_out": "297234276",
  "expiry": 1698901306,
  "fees": {
    "affiliate": "0",
    "asset": "BTC.BTC",
    "liquidity": "150576",
    "outbound": "39000",
    "slippage_bps": 5,
    "total": "189576",
    "total_bps": 6
  },
  "inbound_address": "bc1qmed4v5am2hcg8furkeff2pczdnt0qu4flke420",
  "memo": "-:BTC/BTC:10000",
  "notes": "First output should be to inbound_address, second output should be change back to self, third output should be OP_RETURN, limited to 80 bytes. Do not send below the dust threshold. Do not use exotic spend scripts, locks or address formats (P2WSH with Bech32 address format preferred).",
  "outbound_delay_blocks": 548,
  "outbound_delay_seconds": 3288,
  "slippage_bps": 60,
  "warning": "Do not cache this response. Do not send funds after the expiry."
}
```

```admonish warning
Deposit and withdraw interfaces will return `inbound_address` and `memo` fields that can be used to construct the transaction. Do not cache the`inbound_address` field!
```

### Basic Mechanics

Users can add assets to a vault by sending assets directly to the chain’s vault `address` found on the `/thorchain/inbound_addresses` endpoint. Quote endpoints will also return this.

#### 1. Find the L1 vault address

[`https://thornode.ninerealms.com/thorchain/inbound_addresses`](https://thornode.ninerealms.com/thorchain/inbound_addresses)

Example:

```text
curl -SL https://thornode.ninerealms.com/thorchain/inbound_addresses | jq '.[] | select(.chain == "BTC") | .address'
=> “bc1q556ljv5y4rkdt4p46usx86esljs3xqjxyntlyd”
```

#### 2. Determine if there is capacity available to mint new synths

There is a cap on how many synths can be minted as a function of liquidity depth. To do this, find `synth_mint_paused = false` on the `/pool` endpoint

```text
curl -SL https://thornode.ninerealms.com/thorchain/pools | jq '.[] | select(.asset == "BTC.BTC") | .synth_mint_paused'
```

#### 3. Send memoless savers transactions

Both Saver **Deposit** and **Withdraw** transactions can be done without memos _(optional memos can be included if a wallet wishes, see_ [_`Transaction Memos`_](../concepts/memos.md)_, since there is a marginal transaction cost savings to including memos)._

To **deposit**, users should send any amount of asset they wish (avoiding dust amounts). The network will read the deposit and user address, then add them into the Saver Vault automatically.

To **withdraw**, the user should send a specific dust amount of asset (avoiding the dust threshold), from an amount 0 units above the dust threshold, to an amount 10,000 units above the threshold. \
10000 units is read as “withdraw 10000 basis points”, which is 100%.

```admonish info
The dust threshold is the point at which the network will ignore the amount sent to stop dust attacks (widely seen on UTXO chains).
```

Specific rules for each chain and action are as follows:

- Each chain has a defined `dust_threshold` in base units
- For asset amounts in the range: `[ dust_threshold + 1 : dust_threshold + 10,000]`, the network will withdraw `dust_threshold - 10,000` basis points from the user’s Savers position
- For asset amounts greater than `dust_threshold + 10,000`, the network will add to the user’s Savers position

The `dust_threshold` for each chain are defined as:

- BTC: 10,000 sats
- BCH: 10,000 sats
- LTC: 10,000 sats
- DOGE: 100,000,000 sats
- ETH,AVAX: 0 wei
- ATOM: 0 uatom
- BNB: 0 nbnb

```admonish info
Transactions with asset amounts equal to or below the `dust_threshold` for the chain will be ignored to prevent dust attacks. Ensure you are converting the “human readable” amount (1 BTC) to the correct gas units (100,000,000 sats)
```

**Examples:**

- User wants to deposit 100,000 sats (0.001 BTC): Wallet signs an inbound tx to THORChain’s BTC `/inbound_addresses` vault address from the user with 100,000 sats. This will be added to the user’s Savers position.
- User wants to withdraw 50% of their BTC Savers position: Wallet signs an inbound with 15,000 sats `50% = 5,000 basis points + 10,000[BTC dust_threshold` to THORChain’s BTC vault
- User wants to withdraw 10% of their ETH Savers position: Wallet signs an inbound with 1,000 wei `(10% = 1,000 basis points + 0 [ETH dust_threshold])` to THORChain’s ETH vault
- User wants to deposit 10,000 sats to their DOGE Savers position: Not possible transactions below the `dust_threshold` for each chain are ignored to prevent dust attacks.
- User wants to deposit 20,000 sats to their BTC Savers position: Not possible with memoless, the user’s deposit will be interpreted as a `withdraw:100%`. Instead the user should use a memo.

_translates to: “withdraw 10,000 basis points, or 100% of address’ savings._

### Historical Data & Performance

An important consideration for UIs when implementing this feature is how to display:

- an address’ present performance (targeted at retaining current savers)
- past performance of savings vaults (targeted at attracting potential savers)

#### Present Performance

A user is likely to want to know the following things:

- What is the redeemable value of my share in the Savings Vault?
- What is the absolute amount and % yield I have earned to date on my stake?

The latter can be derived from the former.

`yield_percent = (1 - (depositValue / redeemableValue)) * 100`

```text
saver’s address: bc1qcxssye4j6730h7ehgega3gyykkuwgdgmmpu62n
myUnits => curl -SL https://thornode.ninerealms.com/thorchain/pool/BTC.BTC/savers | jq '.[] | select(.asset_address == "bc1qcxssye4j6730h7ehgega3gyykkuwgdgmmpu62n") | .units'
saverUnits => curl -SL https://thornode.ninerealms.com/thorchain/pools | jq '.[] | select(.asset == "BTC.BTC") | .savers_units'
saverDepth => curl -SL https://thornode.ninerealms.com/thorchain/pools | jq '.[] | select(.asset == "BTC.BTC") | .savers_depth'
```

#### Past Performance

The easy way to determine lifetime performance of the savers vault is to look back 7 days, find the saver value, then compare it with the current saver value.

Example code:\

<iframe width="100%" height="600" src="https://replit.com/@thorchain/THORChain-Savers-Tracker?embed=true" ></iframe>

```admonish info
[https://thornode.ninerealms.com/thorchain/pool/BTC.BTC/savers](https://thornode.ninerealms.com/thorchain/pool/BTC.BTC/savers) will show all BTC Savers
```

### Support

Developers experiencing issues with these APIs can go to the [Developer Discord](https://discord.gg/2Vw3RsQ7) for assistance. Interface developers should subscribe to the #interface-alerts channel for information pertinent to the endpoints and functionality discussed here.
