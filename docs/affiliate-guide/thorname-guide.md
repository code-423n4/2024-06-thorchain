# THORName Guide

## Summary

[THORNames](https://docs.thorchain.org/how-it-works/thorchain-name-service) are THORChain's vanity address system that allows affiliates to collect fees and track their user's transactions. THORNames exist on the THORChain L1, so you will need a THORChain address and $RUNE to create and manage a THORName.;

THORNames have the following properties:

- **Name:** The THORName's string. Between 1-30 hexadecimal characters and `-_+` special characters.;
- **Owner**: This is the THORChain address that owns the THORName
- **Aliases**: THORNames can have an alias address for any external chain supported by THORChain, and can have an alias for the THORChain L1 that is different than the owner.
- **Expiry:** THORChain Block-height at which the THORName expires.
- **Preferred Asset:** The asset to pay out affiliate fees in. This can be any asset supported by THORChain.;

## Create a THORName

THORNames are created by posting a `MsgDeposit` to the THORChain network with the appropriate [memo](../concepts/memos.md) and enough $RUNE to cover the registration fee and to pay for the amount of blocks the THORName should be registered for.;

- **Registration fee**: `tns_register_fee_rune` on the [Network endpoint](https://thornode.ninerealms.com/thorchain/network). This value is in 1e8, so `100000000 = 1 $RUNE`
- **Per block fee**: `tns_fee_per_block_rune` on the same endpoint, also in 1e8.;

For example, for a new THORName to be registered for 10 years the amount paid would be:

`amt = tns_register_fee_rune + tns_fee_per_block_rune * 10 * 5256000`

`5256000 = avg # of blocks per year`

The expiration of the THORName will automatically be set to the number of blocks in the future that was paid for minus the registration fee.

**Memo Format:**

Memo template is: `~:name:chain:address:?owner:?preferredAsset:?expiry`

- **name**: Your THORName. Must be unique, between 1-30 characters, hexadecimal and `-_+` special characters.;
- **chain:** The chain of the alias to set.;
- **address**: The alias address. Must be an address of chain.
- **owner**: THORChain address of owner (optional).
- **preferredAsset:** Asset to receive fees in. Must be supported be an active pool on THORChain. Value should be `asset` property from the [Pools endpoint](https://thornode.ninerealms.com/thorchain/pools).;

```admonish info
Example: `~:ODIN:BTC:bc1Address:thorAddress:BTC.BTC`
```

This will register a new THORName called `ODIN` with a Bitcoin alias of `bc1Address` owner of `thorAddress` and preferred asset of BTC.BTC.

```admonish info
You can use [Asgardex](https://github.com/thorchain/asgardex-electron) to post a MsgDeposit with a custom memo. Load your wallet, then open your THORChain wallet page > Deposit > Custom.;
```

```admonish info
View your THORName's configuration at the THORName endpoint:

e.g. [https://thornode.ninerealms.com/thorchain/thorname/](https://thornode.ninerealms.com/thorchain/thorname/ac-test){name}
```

## Renewing your THORName

All THORName's have a expiration represented by a THORChain block-height. Once the expiration block-height has passed, another THORChain address can claim the THORName and any associated balance in the Affiliate Fee Collector Module (Read [#preferred-asset-for-affiliate-fees](thorname-guide.md#preferred-asset-for-affiliate-fees "mention")), so it's important to monitor this and renew your THORName as needed.;

To keep your THORName registered you can extend the registration period (move back the expiration block height), by posting a `MsgDeposit` with the correct THORName memo and $RUNE amount.;

**Memo:**

`~:ODIN:THOR:<thor-alias-address>`

_(Chain and alias address are required, so just use current values to keep alias unchanged)._

**$RUNE Amount:**

`rune_amt = num_blocks_to_extend * tns_fee_per_block`

_(Remember this value will be in 1e8, so adjust accordingly for your transaction)._

## Preferred Asset for Affiliate Fees

Starting in THORNode V116, affiliates can collect their fees in the asset of their choice (choosing from the assets that have a pool on THORChain). In order to collect fees in a preferred asset, affiliates must use a [THORName](https://docs.thorchain.org/how-it-works/thorchain-name-service) in their swap memos.;

### How it Works

If an affiliate's THORName has the proper preferred asset configuration set, the network will begin collecting their affiliate fees in $RUNE in the [AffiliateCollector module](https://thornode.ninerealms.com/thorchain/balance/module/affiliate_collector). Once the accrued RUNE in the module is greater than [`PreferredAssetOutboundFeeMultiplier`](https://gitlab.com/thorchain/thornode/-/blob/develop/constants/constants_v1.go#L107)`* outbound_fee` of the preferred asset's chain, the network initiates a swap from $RUNE -> Preferred Asset on behalf of the affiliate. At the time of writing, `PreferredAssetOutboundFeeMultiplier` is set to `100`, so the preferred asset swap happens when the outbound fee is 1% of the accrued $RUNE.;

**Configuring a Preferred Asset for a THORName:**

1. Register your THORName following instructions above.
2. Set your preferred asset's chain alias (the address you'll be paid out to), and your preferred asset. _Note: your preferred asset must be currently supported by THORChain._

For example, if you wanted to be paid out in USDC you would:

1. Grab the full USDC name from the [Pools](https://thornode.ninerealms.com/thorchain/pools) endpoint: `ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48`
2. Post a `MsgDeposit` to the THORChain network with the appropriate memo to register your THORName, set your preferred asset as USDC, and set your Ethereum network address alias. Assuming the following info:

   1. THORChain address: `thor1dl7un46w7l7f3ewrnrm6nq58nerjtp0dradjtd`
   2. THORName: `ac-test`
   3. ETH payout address: `0x6621d872f17109d6601c49edba526ebcfd332d5d`;

   The full memo would look like:

   > `~:ac-test:ETH:0x6621d872f17109d6601c49edba526ebcfd332d5d:thor1dl7un46w7l7f3ewrnrm6nq58nerjtp0dradjtd:ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48`

```admonish info
You will also need a THOR alias set to collect affiliate fees. Use another MsgDeposit with memo: `~:<thorname>:THOR:<thorchain-address>` to set your THOR alias. Your THOR alias address can be the same as your owner address, but won't be used for anything if a preferred asset is set.;
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
