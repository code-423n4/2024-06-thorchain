# Fees

## Overview

There are 4 different fees the user should know about.

1. Inbound Fee (sourceChain: gasRate \* txSize)
2. Affiliate Fee (affiliateFee \* swapAmount)
3. Liquidity Fee (swapSlip \* swapAmount)
4. Outbound Fee (destinationChain: gasRate \* txSize)

### **Terms**

- **SourceChain**: the chain the user is swapping from
- **DestinationChain**: the chain the user is swapping to txSize: the size of the transaction in bytes (or units)
- **gasRate**: the current gas rate of the external network
- **swapAmount**: the amount the user is swapping swapSlip: the slip created by the
- **swapAmount**, as a function of poolDepth
- **affiliateFee**: optional fee set by interface in basis points

## Fees Detail

### Inbound Fee

This is the fee the user pays to make a transaction on the source chain, which the user pays directly themselves. The gas rate recommended to use is `fast` where the tx is guaranteed to be committed in the next block. Any longer and the user will be waiting a long time for their swap and their price will be invalid (thus they may get an unnecessary refund).

$$
inboundFee = txSize * gasRate
$$

```admonish success
THORChain calculates and posts fee rates at [`https://thornode.ninerealms.com/thorchain/inbound_addresses`](https://thornode.ninerealms.com/thorchain/inbound_addresses)
```

```admonish warning
Always use a "fast" or "fastest" fee, if the transaction is not confirmed in time, it could be abandoned by the network or failed due to old prices. You should allow your users to cancel or re-try with higher fees.
```

### Liquidity Fee

This is simply the slip created by the transaction multiplied by its amount. It is priced and deducted from the destination amount automatically by the protocol.

$$
slip = \frac{swapAmount}{swapAmount + poolDepth}
$$

$$
fee =slip * swapAmount
$$

### Affiliate Fee

In the swap transaction you build for your users you can include an affiliate fee for your exchange (accepted in $RUNE or a synthetic asset, so you will need a $RUNE address).

- The affiliate fee is in basis points (0-10,000) and will be deducted from the inbound swap amount from the user.
- If the inbound swap asset is a native THORChain asset ($RUNE or synth) the affiliate fee amount will be deducted directly from the transaction amount.
- If the inbound swap asset is on any other chain the network will submit a swap to $RUNE with the destination address as your affiliate fee address.
- If the affiliate is added to an ADDLP tx, then the affiliate is included in the network as an LP.

`SWAP:CHAIN.ASSET:DESTINATION:LIMIT:AFFILIATE:FEE`

Read [https://medium.com/thorchain/affiliate-fees-on-thorchain-17cbc176a11b](https://medium.com/thorchain/affiliate-fees-on-thorchain-17cbc176a11b) for more information.

$$
affliateFee = \frac{feeInBasisPoints * swapAmount}{10000}
$$

### Preferred Asset for Affiliate Fees

Affiliates can collect their fees in the asset of their choice (choosing from the assets that have a pool on THORChain). In order to collect fees in a preferred asset, affiliates must use a [THORName](../affiliate-guide/thorname-guide.md) in their swap [memos](memos.md#swap).

### How it Works

If an affiliate's THORName has the proper preferred asset configuration set, the network will begin collecting their affiliate fees in $RUNE in the [AffiliateCollector module](https://thornode.ninerealms.com/thorchain/balance/module/affiliate_collector). Once the accrued RUNE in the module is greater than [`PreferredAssetOutboundFeeMultiplier`](https://gitlab.com/thorchain/thornode/-/blob/develop/constants/constants_v1.go#L107)`* outbound_fee` of the preferred asset's chain, the network initiates a swap from $RUNE -> Preferred Asset on behalf of the affiliate. At the time of writing, `PreferredAssetOutboundFeeMultiplier` is set to `100`, so the preferred asset swap happens when the outbound fee is 1% of the accrued $RUNE.

**Configuring a Preferred Asset for a THORName.**

1. [**Register a THORName**](../affiliate-guide/thorname-guide.md) if not done already. This is done with a `MsgDeposit` posted to the THORChain network.
2. Set your preferred asset's chain alias (the address you'll be paid out to), and your preferred asset. _Note: your preferred asset must be currently supported by THORChain._

For example, if you wanted to be paid out in USDC you would:

1. Grab the full USDC name from the [Pools](https://thornode.ninerealms.com/thorchain/pools) endpoint: `ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48`
2. Post a `MsgDeposit` to the THORChain network with the appropriate memo to register your THORName, set your preferred asset as USDC, and set your Ethereum network address alias. Assuming the following info:

   1. THORChain address: `thor1dl7un46w7l7f3ewrnrm6nq58nerjtp0dradjtd`
   2. THORName: `ac-test`
   3. ETH payout address: `0x6621d872f17109d6601c49edba526ebcfd332d5d`

   The full memo would look like:

   > `~:ac-test:ETH:0x6621d872f17109d6601c49edba526ebcfd332d5d:thor1dl7un46w7l7f3ewrnrm6nq58nerjtp0dradjtd:ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48`

```admonish info
You can use [Asgardex](https://github.com/thorchain/asgardex-electron) to post a MsgDeposit with a custom memo. Load your wallet, then open your THORChain wallet page > Deposit > Custom.
```

```admonish info
You will also need a THOR alias set to collect affiliate fees. Use another MsgDeposit with memo: `~:<thorname>:THOR:<thorchain-address>` to set your THOR alias. Your THOR alias address can be the same as your owner address, but won't be used for anything if a preferred asset is set.
```

Once you successfully post your MsgDeposit you can verify that your THORName is configured properly. View your THORName info from THORNode at the following endpoint:\
[https://thornode.ninerealms.com/thorchain/thorname/ac-test](https://thornode.ninerealms.com/thorchain/thorname/ac-test)

The response should look like:

```json
{
  "affiliate_collector_rune": "0",
  "aliases": [
    {
      "address": "0x6621d872f17109d6601c49edba526ebcfd332d5d",
      "chain": "ETH"
    },
    {
      "address": "thor1dl7un46w7l7f3ewrnrm6nq58nerjtp0dradjtd",
      "chain": "THOR"
    }
  ],
  "expire_block_height": 22061405,
  "name": "ac-test",
  "owner": "thor1dl7un46w7l7f3ewrnrm6nq58nerjtp0dradjtd",
  "preferred_asset": "ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48"
}
```

Your THORName is now properly configured and any affiliate fees will begin accruing in the AffiliateCollector module. You can verify that fees are being collected by checking the `affiliate_collector_rune` value of the above endpoint.

### Outbound Fee

This is the fee the Network pays on behalf of the user to send the outbound transaction. To adequately pay for network resources (TSS, compute, state storage) the fee is marked up from what nodes actually pay on-chain by an "Outbound Fee Multiplier" (OFM).

The OFM moves between a `MaxOutboundFeeMultiplier` and a `MinOutboundFeeMultiplier`(defined as [Network Constants](https://gitlab.com/thorchain/thornode/-/blob/develop/constants/constants_v1.go) or as [Mimir Values](https://thornode.ninerealms.com/thorchain/mimir)), based on the network's current outbound fee "surplus" in relation to a "target surplus". The outbound fee "surplus" is the cumulative difference (in $RUNE) between what the users are charged for outbound fees and what the nodes actually pay. As the network books a "surplus" the OFM slowly decreases from the Max to the Min. Current values for the OFM can be found on the [Network Endpoint](https://thornode.ninerealms.com/thorchain/network).

$$
outboundFee = txSize * gasRate * OFM
$$

The minimum Outbound Layer1 Fee the network will charge is on `/thorchain/mimir` and is priced in USD (based on THORChain's USD pool prices). This means really cheap chains still pay their fair share. It is currently set to `100000000` = $1.00

See [Outbound Fee](https://docs.thorchain.org/how-it-works/fees#outbound-fee) for more information.

## Fee Ordering for Swaps

Fees are taken in the following order when conducting a swap.

1. Inbound Fee (user wallet controlled, not THORChain controlled)
2. Affiliate Fee (if any) - skimmed from the input.
3. Swap Fee (denoted in output asset)
4. Outbound Fee (taken from the swap output)

To work out the total fees, fees should be converted to a common asset (e.g. RUNE or USD) then added up. Total fees should be less than the input else it is likely to result in a refund.

### Refunds and Minimum Swappable Amount

If a transaction fails, it is refunded, thus it will pay the `outboundFee` for the **SourceChain** not the DestinationChain. Thus devs should always swap an amount that is a maximum of the following, multiplier by a buffer of at least 4x to allow for sudden gas spikes:

1. The Destination Chain outbound_fee
2. The Source Chain outbound_fee
3. $1.00 (the minimum)

The outbound_fee for each chain is returned on the [Inbound Addresses](https://thornode.ninerealms.com/thorchain/inbound_addresses) endpoint, priced in the gas asset.

It is strongly recommended to use the `recommended_min_amount_in` value that is included on the [Swap Quote](broken-reference) endpoint, which is the calculation described above. This value is priced in the inbound asset of the quote request (in 1e8). This should be the minimum-allowed swap amount for the requested quote.

_Remember, if the swap limit is not met or the swap is otherwise refunded the outbound_fee of the Source Chain will be deducted from the input amount, so give your users enough room._

### Understanding gas_rate

THORNode keeps track of current gas prices. Access these at the `/inbound_addresses` endpoint of the [THORNode API](https://dev.thorchain.org/thorchain-dev/wallets/connecting-to-thorchain#thornode). The response is an array of objects like this:

```json
{
    "chain": "ETH",
    "pub_key": "thorpub1addwnpepqdlx0avvuax3x9skwcpvmvsvhdtnw6hr5a0398vkcvn9nk2ytpdx5cpp70n",
    "address": "0x74ce1c3556a6d864de82575b36c3d1fb9c303a80",
    "router": "0x3624525075b88B24ecc29CE226b0CEc1fFcB6976",
    "halted": false,
    "gas_rate": "10"
    "gas_rate_units": "satsperbyte",
    "outbound_fee": "30000",
    "outbound_tx_size": "1000",
}
```

The `gas_rate` property can be used to estimate network fees for each chain the swap interacts with. For example, if the swap is `BTC -> ETH` the swap will incur fees on the bitcoin network and Ethereum network. The `gas_rate` property works differently on each chain "type" (e.g. EVM, UTXO, BFT).

The `gas_rate_units` explain what the rate is for chain, as a prompt to the developer.

The `outbound_tx_size` is what THORChain internally budgets as a typical transaction size for each chain.

The `outbound_fee` is `gas_rate * outbound_tx_size * OFM` and developers can use this to budget for the fee to be charged to the user. The current Outbound Fee Multiplier (OFM) can be found on the [Network Endpoint](https://thornode.ninerealms.com/thorchain/network).

Keep in mind the `outbound_fee` is priced in the gas asset of each chain. For chains with tokens, be sure to convert the `outbound_fee` to the outbound token to determine how much will be taken from the outbound amount. To do this, use the `getValueOfAsset1InAsset2` formula described in the [`Math`](https://dev.thorchain.org/thorchain-dev/interface-guide/math#example-1) section.

## Fee Calculation by Chain

### **THORChain (Native Rune)**

The THORChain blockchain has a set 0.02 RUNE fee. This is set within the THORChain [Constants](https://thornode.ninerealms.com/thorchain/constants) by `NativeTransactionFee`. As THORChain is 1e8, `2000000 TOR = 0.02 RUNE`

### Binance Chain

THORChain uses the gas_rate as the flat Binance Chain transaction fee.

E.g. If the `gas_rate` = 11250 then fee is 0.0011250 BNB.

### UTXO Chains like Bitcoin

For UXTO chains link Bitcoin, `gas_rate`is denoted in Satoshis. The `gas_rate` is calculated by looking at the average previous block fee seen by the THORNodes.

All THORChain transactions use BECH32 so a standard tx size of 250 bytes can be used. The standard UTXO fee is then `gas_rate`\* 250.

### EVM Chains like Ethereum

For EVM chains like Ethereum, `gas_rate`is denoted in GWEI. The `gas_rate` is calculated by looking at the average previous block fee seen by the THORNodes

An Ether Tx fee is: `gasRate * 10^9 (GWEI) * 21000 (units).`

An ERC20 Tx is larger: `gasRate * 10^9 (GWEI) * 70000 (units)`

```admonish success
THORChain calculates and posts gas fee rates at [`https://thornode.ninerealms.com/thorchain/inbound_addresses`](https://thornode.ninerealms.com/thorchain/inbound_addresses)
```
