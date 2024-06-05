# Fees and Wait Times

## Fee Types

Users pay up to four kinds of fees when conducting a swap.

1. **Layer1 Network Fees** (gas): paid by the user when sending the asset to THORChain to be swapped. This is controlled by the user's wallet.
2. **Slip Fee**: protects the pool from being manipulated by large swaps. Calculated as a function of transaction size and current pool depth. The slip fee formula is explained [here](https://docs.thorchain.org/thorchain-finance/continuous-liquidity-pools#clp-derivation) and an example implementation is [here](https://gitlab.com/thorchain/asgardex-common/asgardex-util/-/blob/master/src/calc/swap.ts#L57).
3. **Affiliate Fee** - (optional) a percentage skimmed from the inbound amount that can be paid to exchanges or wallet providers. _Wallets can now accept fees in any THORChain-supported asset (USDC, BTC, etc). Check the "Preferred Asset for Affiliate Fees" section in_ [fees.md](../concepts/fees.md "mention") _for more details and setup information._
4. **Outbound Fee** - the fee the Network pays on behalf of the user to send the outbound transaction. See [Outbound Fee](../concepts/fees.md#outbound-fee).

```admonish info
The Swap Quote endpoint will calculate and show all fees.
```

See the [fees](../concepts/fees.md) section for full details.

### Refunds and Minimum Swap Amount

If a transaction fails, it is refunded, thus it will pay the `outboundFee` for the **SourceChain** not the DestinationChain. Thus devs should always swap an amount that is a maximum of the following, multiplied by at least a 4x buffer to allow for gas spikes:

1. The Destination Chain outboundFee, or
2. The Source Chain outboundFee, or
3. $1.00 (the minimum outboundFee).

For convenience, a `recommended_min_amount_in` is included on the [Swap Quote](broken-reference) endpoint, which is the value described above. This value is priced in the inbound asset of the quote request (in 1e8). This should be the minimum-allowed swap amount for the requested quote.

## Wait Times

There are four phases of a transaction sent to THORChain each taking time to complete.

1. **Layer1 Inbound Confirmation -** assuming the inboundTx will be confirmed in the next block, it is the source blockchain block time.
2. **Observation Counting** - time for 67% THORChain Nodes to observe and agree on the inboundTx.
3. **Confirmation Counting** - for non-instant finality blockchains, the amount of time THORChain will wait before processing to protect against double spends and re-org attacks.
4. **Outbound Delay** - dependent on size and network traffic. Large outbounds will be delayed.
5. **Layer1 Outbound Confirmation** - Outbound blockchain block time.

Wait times can be between a few seconds up to an hour. The assets being swapped, the size of the swap and the current network traffic within THORChain will determine the wait time.

```admonish info
The Swap Quote endpoint will calculate points 3 and 4.
```

See the [delays.md](../concepts/delays.md "mention") section for full details.
