# Quickstart Guide

## Introduction

THORChain allows native L1 Swaps. On-chain [Memos](../concepts/memos.md) are used instruct THORChain how to swap, with the option to add [price limits](quickstart-guide.md#price-limits) and [affiliate fees](quickstart-guide.md#affiliate-fees). THORChain nodes observe the inbound transactions and when the majority have observed the transactions, the transaction is processed by threshold-signature transactions from THORChain vaults.

Let's demonstrate decentralized, non-custodial cross-chain swaps. In this example, we will build a transaction that instructs THORChain to swap native Bitcoin to native Ethereum in one transaction.

```admonish info
The following examples use a free, hosted API provided by [Nine Realms](https://twitter.com/ninerealms_cap). If you want to run your own full node, please see [connecting-to-thorchain.md](../concepts/connecting-to-thorchain.md "mention").
```

### 1. Determine the correct asset name

THORChain uses a specific [asset notation](../concepts/asset-notation.md#layer-1-assets). Available assets are at: [Pools Endpoint.](https://thornode.ninerealms.com/thorchain/pools)

BTC => `BTC.BTC`\
ETH => `ETH.ETH`

```admonish info
Only available pools can be used. (`where 'status' == Available)`
```

### 2. Query for a swap quote

```admonish info
All amounts are 1e8. Multiply native asset amounts by 100000000 when dealing with amounts in THORChain. 1 BTC = 100,000,000.
```

**Request**: _Swap 1 BTC to ETH and send the ETH to_ `0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430`.

[https://thornode.ninerealms.com/thorchain/quote/swap?from_asset=BTC.BTC\&to_asset=ETH.ETH\&amount=100000000\&destination=0x86d526d6624AbC0178cF7296cD538Ecc080A95F1](https://thornode.ninerealms.com/thorchain/quote/swap?from_asset=BTC.BTC&to_asset=ETH.ETH&amount=100000000&destination=0x86d526d6624AbC0178cF7296cD538Ecc080A95F1)

**Response**:

```json
{
  "dust_threshold": "10000",
  "expected_amount_out": "1619355520",
  "expiry": 1689143119,
  "fees": {
    "affiliate": "0",
    "asset": "ETH.ETH",
    "outbound": "240000"
  },
  "inbound_address": "bc1qpzs9rm82m08u48842ka59hyxu36wsgzqlt6e3t",
  "inbound_confirmation_blocks": 1,
  "inbound_confirmation_seconds": 600,
  "max_streaming_quantity": 0,
  "memo": "=:ETH.ETH:0x86d526d6624AbC0178cF7296cD538Ecc080A95F1",
  "notes": "First output should be to inbound_address, second output should be change back to self, third output should be OP_RETURN, limited to 80 bytes. Do not send below the dust threshold. Do not use exotic spend scripts, locks or address formats (P2WSH with Bech32 address format preferred).",
  "outbound_delay_blocks": 305,
  "outbound_delay_seconds": 1830,
  "recommended_min_amount_in": "60000",
  "slippage_bps": 49,
  "streaming_swap_blocks": 0,
  "total_swap_seconds": 2430,
  "warning": "Do not cache this response. Do not send funds after the expiry."
}
```

_If you send 1 BTC to `bc1qlccxv985m20qvd8g5yp6g9lc0wlc70v6zlalz8` with the memo `=:ETH.ETH:0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430`, you can expect to receive `13.4493552` ETH._

_For security reasons, your inbound transaction will be delayed by 600 seconds (1 BTC Block) and 2040 seconds (or 136 native THORChain blocks) for the outbound transaction,_ 2640 seconds all up*. You will pay an outbound gas fee of 0.0048 ETH and will incur 41 basis points (0.41%) of slippage.*

```admonish info
Full quote swap endpoint specification can be found here: [https://thornode.ninerealms.com/thorchain/doc/](https://thornode.ninerealms.com/thorchain/doc/).

See an example implementation [here.](https://replit.com/@thorchain/quoteSwap#index.js)
```

If you'd prefer to calculate the swap yourself, see the [Fees](fees-and-wait-times.md) section to understand what fees need to be accounted for in the output amount. Also, review the [Transaction Memos](../concepts/memos.md) section to understand how to create the swap memos.

### 3. Sign and send transactions on the from_asset chain

Construct, sign and broadcast a transaction on the BTC network with the following parameters:

Amount => `1.0`

Recipient => `bc1qlccxv985m20qvd8g5yp6g9lc0wlc70v6zlalz8`

Memo => `=:ETH.ETH:0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430`

```admonish warning
Never cache inbound addresses! Quotes should only be considered valid for 10 minutes. Sending funds to an old inbound address will result in loss of funds.
```

```admonish info
Learn more about how to construct inbound transactions for each chain type here: [Sending Transactions](../concepts/sending-transactions.md)
```

### 4. Receive tokens

Once a majority of nodes have observed your inbound BTC transaction, they will sign the Ethereum funds out of the network and send them to the address specified in your transaction. You have just completed a non-custodial, cross-chain swap by simply sending a native L1 transaction.

## Additional Considerations

```admonish warning
There is a rate limit of 1 request per second per IP address on /quote endpoints. It is advised to put a timeout on frontend components input fields, so that a request for quote only fires at most once per second. If not implemented correctly, you will receive 503 errors.
```

```admonish success
For best results, request a new quote right before the user submits a transaction. This will tell you whether the _expected_amount_out_ has changed or if the _inbound_address_ has changed. Ensuring that the _expected_amount_out_ is still valid will lead to better user experience and less frequent failed transactions.
```

### Price Limits

Specify _tolerance_bps_ to give users control over the maximum slip they are willing to experience before canceling the trade. If not specified, users will pay an unbounded amount of slip.

[https://thornode.ninerealms.com/thorchain/quote/swap?amount=100000000\&from_asset=BTC.BTC\&to_asset=ETH.ETH\&destination=0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430\&tolerance_bps=100](https://thornode.ninerealms.com/thorchain/quote/swap?amount=100000000&from_asset=BTC.BTC&to_asset=ETH.ETH&destination=0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430&tolerance_bps=100)

`https://thornode.ninerealms.com/thorchain/quote/swap?amount=100000000&from_asset=BTC.BTC&to_asset=ETH.ETH&destination=0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430&tolerance_bps=100`

Notice how a minimum amount (1342846539 / \~13.42 ETH) has been appended to the end of the memo. This tells THORChain to revert the transaction if the transacted amount is more than 100 basis points less than what the _expected_amount_out_ returns.

### Affiliate Fees

Specify `affiliate` and `affiliate_bps` to skim a percentage of the swap as an affiliate fee. When a valid affiliate address and affiliate basis points are present in the memo, the protocol will skim affiliate_bps from the inbound swap amount and swap this to $RUNE with the affiliate address as the destination address.

Params:

- **affiliate**: Can be a THORName or valid THORChain address
- **affiliate_bps**: 0-1000 basis points

Memo format:
`=:BTC.BTC:<destination_addr>:<limit>:<affiliate>:<affiliate_bps>`

Quote example:

[https://thornode.ninerealms.com/thorchain/quote/swap?amount=100000000\&from_asset=BTC.BTC\&to_asset=ETH.ETH\&destination=0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430\&affiliate=thorname\&affiliate_bps=10](https://thornode.ninerealms.com/thorchain/quote/swap?amount=100000000&from_asset=BTC.BTC&to_asset=ETH.ETH&destination=0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430&affiliate=thorname&affiliate_bps=10)

```json
{
  "dust_threshold": "10000",
  "expected_amount_out": "1603383828",
  "expiry": 1688973775,
  "fees": {
    "affiliate": "1605229",
    "asset": "ETH.ETH",
    "outbound": "240000"
  },
  "inbound_address": "bc1qhkutxeluztncm5pq0ckpm75hztrv7m7nhhh94d",
  "inbound_confirmation_blocks": 1,
  "inbound_confirmation_seconds": 600,
  "max_streaming_quantity": 0,
  "memo": "=:ETH.ETH:0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430::thorname:10",
  "notes": "First output should be to inbound_address, second output should be change back to self, third output should be OP_RETURN, limited to 80 bytes. Do not send below the dust threshold. Do not use exotic spend scripts, locks or address formats (P2WSH with Bech32 address format preferred).",
  "outbound_delay_blocks": 303,
  "outbound_delay_seconds": 1818,
  "recommended_min_amount_in": "72000",
  "slippage_bps": 49,
  "streaming_swap_blocks": 0,
  "total_swap_seconds": 2418,
  "warning": "Do not cache this response. Do not send funds after the expiry."
}
```

Notice how `thorname:10` has been appended to the end of the memo. This instructs THORChain to skim 10 basis points from the swap. The user should still expect to receive the _expected_amount_out,_ meaning the affiliate fee has already been subtracted from this number.

For more information on affiliate fees: [fees.md](../concepts/fees.md "mention").

### Streaming Swaps

[_Streaming Swaps_](streaming-swaps.md) _can be used to break up the trade to reduce slip fees._

Params:

- **streaming_interval**: # of THORChain blocks between each subswap. Larger # of blocks gives arb bots more time to rebalance pools. For deeper/more active pools a value of `1` is most likely okay. For shallower/less active pools a larger value should be considered.
- **streaming_quantity**: # of subswaps to execute. If this value is omitted or set to `0` the protocol will calculate the # of subswaps such that each subswap has a slippage of 5 bps.

Memo format:
`=:BTC.BTC:<destination_addr>:<limit>/<streaming_interval>/<streaming_quantity>`

Quote example:

[_https://stagenet-thornode.ninerealms.com/thorchain/quote/swap?amount=100000000\&from_asset=BTC.BTC\&to_asset=ETH.ETH\&destination=0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430\&streaming_interval=10_](https://stagenet-thornode.ninerealms.com/thorchain/quote/swap?amount=100000000&from_asset=BTC.BTC&to_asset=ETH.ETH&destination=0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430&streaming_interval=10)

```json
{
  "approx_streaming_savings": 0.99930555,
  "dust_threshold": "10000",
  "expected_amount_out": "145448080",
  "expiry": 1689117597,
  "fees": {
    "affiliate": "0",
    "asset": "ETH.ETH",
    "outbound": "480000"
  },
  "inbound_address": "bc1qk2z8luw2afwuugndynegn72dkv45av5hyjrtm8",
  "inbound_confirmation_blocks": 1,
  "inbound_confirmation_seconds": 600,
  "max_streaming_quantity": 1440,
  "memo": "=:ETH.ETH:0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430:0/10/1440",
  "notes": "First output should be to inbound_address, second output should be change back to self, third output should be OP_RETURN, limited to 80 bytes. Do not send below the dust threshold. Do not use exotic spend scripts, locks or address formats (P2WSH with Bech32 address format preferred).",
  "outbound_delay_blocks": 76,
  "outbound_delay_seconds": 456,
  "recommended_min_amount_in": "158404",
  "slippage_bps": 8176,
  "streaming_swap_blocks": 14400,
  "streaming_swap_seconds": 86400,
  "total_swap_seconds": 87456,
  "warning": "Do not cache this response. Do not send funds after the expiry."
}
```

Notice how `approx_streaming_savings` shows the savings by using streaming swaps. `total_swap_seconds` also shows the amount of time the swap will take.

### Custom Refund Address

By default, in the case of a refund the protocol will return the inbound swap to the original sender. However, in the case of protocol <> protocol interactions, many times the original sender is a smart contract, and not the user's EOA. In these cases, a custom refund address can be defined in the memo, which will ensure the user will receive the refund and not the smart contract.

Params:

- **refund_address**: User's refund address. Needs to be a valid address for the inbound asset, otherwise refunds will be returned to the sender

Memo format:
`=:BTC.BTC:<destination>/<refund_address>`

Quote example:
[https://thornode.ninerealms.com/thorchain/quote/swap?amount=100000000&from_asset=ETH.ETH&to_asset=BTC.BTC&destination=bc1qyl7wjm2ldfezgnjk2c78adqlk7dvtm8sd7gn0q&refund_address=0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430](https://thornode.ninerealms.com/thorchain/quote/swap?amount=100000000&from_asset=ETH.ETH&to_asset=BTC.BTC&destination=bc1qyl7wjm2ldfezgnjk2c78adqlk7dvtm8sd7gn0q&refund_address=0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430)

```json
{
  ...
  "memo": "=:BTC.BTC:bc1qyl7wjm2ldfezgnjk2c78adqlk7dvtm8sd7gn0q/0x3021c479f7f8c9f1d5c7d8523ba5e22c0bcb5430",
  ...
}
```

### Error Handling

The quote swap endpoint simulates all of the logic of an actual swap transaction. It ships with comprehensive error handling.

![Price Tolerance Error](../.gitbook/assets/image (6).png)
_This error means the swap cannot be completed given your price tolerance._

![Destination Address Error](../.gitbook/assets/image (1).png)
_This error ensures the destination address is for the chain specified by `to_asset`._

![Affiliate Address Length Error](../.gitbook/assets/image (4).png)
_This error is due to the fact the affiliate address is too long given the source chain's memo length requirements. Try registering a THORName to shorten the memo._

![Asset Not Found Error](../.gitbook/assets/image (2).png)
_This error means the requested asset does not exist._

![Bound Checks Error](../.gitbook/assets/image (3).png)
_Bound checks are made on both `affiliate_bps` and `tolerance_bps`._

### Support

Developers experiencing issues with these APIs can go to the [Developer Discord](https://discord.gg/2Vw3RsQ7) for assistance. Interface developers should subscribe to the #interface-alerts channel for information pertinent to the endpoints and functionality discussed here.
