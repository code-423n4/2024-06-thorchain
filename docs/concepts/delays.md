# Delays

## Overview

There are four phases of a transaction sent to THORChain.

1. [Inbound Confirmation](delays.md#inbound-confirmation)
2. [Observation Counting](delays.md#observation-counting)
3. [Confirmation Counting](delays.md#confirmation-counting)
4. [Outbound Delay](delays.md#outbound-delay)
5. [Outbound Confirmation](delays.md#txout-delay)

Wait times can be between a few seconds to several hours. The assets being swapped, the size of the swap and the current network traffic within THORChain will determine the wait time.

### Inbound Confirmation

This depends purely on the host chain and is out of the control of THORChain.

- Bitcoin/BitcoinCash: \~10 minutes
- Litecoin: \~2.5 minutes
- Dogecoin: \~60 seconds
- ETH: \~15 seconds
- Cosmos: \~6 seconds

### Observation Counting

THORNodes have to witness to THORChain when they see a transaction. It could seconds to minutes depending on how fast nodes can scan their blockchains to find transactions. Once 67% of THORNodes see a tx, then it can be confirmed. You can count the number of nodes that have seen a tx by counting the signatures in the `signers` parameter or look at the `status` field on the `/tx` endpoint.

Example: [https://thornode.ninerealms.com/thorchain/tx/0AAA205438B6409CBA11DED8C8F63794D719CF4E3818B85117259311E3ADEA0E](https://thornode.ninerealms.com/thorchain/tx/0AAA205438B6409CBA11DED8C8F63794D719CF4E3818B85117259311E3ADEA0E)

### Confirmation Counting

THORChain has to defend against 51% attacks, which it does by counting to economic finality for each block (the value of the block relative to the value of the block reward). It tracks both, then computes the number of blocks to wait. It then populates this on the `/tx` endpoint.

Example: [https://thornode.ninerealms.com/thorchain/tx/0AAA205438B6409CBA11DED8C8F63794D719CF4E3818B85117259311E3ADEA0E](https://thornode.ninerealms.com/thorchain/tx/0AAA205438B6409CBA11DED8C8F63794D719CF4E3818B85117259311E3ADEA0E)

`block_height` is the external height it first saw it.

`finalise_height` is the external height it needs to see before it will confirm it.

```admonish warning
An event is not sent until the external block height crosses `finalise_height` so Midgard will NOT see the tx until confirmation-counted.
```

Examples:

- 10 BTC: 2 blocks
- 50 ETH: 16 blocks
- 100 LTC: 9 blocks

### Outbound Delay

THORChain throttles all outputs to prevent fund loss attacks. The maximum delay is 720 blocks which is approx 1 hour. Outbound delay worked out by computing the value of the outbound transaction in RUNE then applying an artificial delay. If the tx is in "scheduled", it will be delayed by a number of blocks. Once it is "outbound" it is being processed. See more information [here](https://docs.thorchain.org/how-it-works/security#b905-1).

```admonish info
Arbs and Traders who have trade history can have a reduced wait time to do [Swapper Clout.](./swapper-clout.md)
```

**Queue:**

[https://thornode.ninerealms.com/thorchain/queue](https://thornode.ninerealms.com/thorchain/queue)

**Delayed txOuts:**

[https://thornode.ninerealms.com/thorchain/queue/scheduled](https://thornode.ninerealms.com/thorchain/queue/scheduled)

**Finalised txOuts:**

[https://thornode.ninerealms.com/thorchain/queue/outbound](https://thornode.ninerealms.com/thorchain/queue/outbound)

**Swap Clout:**

[Clout](./swapper-clout.md) is the cumulative total fees paid (in RUNE) for a given address. Swappers with a high clout have proven themselves to be highly aligned to the project and therefore can reap the rewards by getting faster trade execution. Clout is deducted from the RUNE value of an outbound before a delay relative to the outbound value is applied. Clout is removed from an address when an outbound using it is scheduled, and the clout is reclaimed when the outbound achieves observation consensus.

### Outbound Confirmation

This depends purely on the host chain and is out of the control of THORChain.

- Bitcoin/BitcoinCash: \~10 minutes
- Litecoin: \~2.5 minutes
- Dogecoin: \~60 seconds
- ETH: \~15 seconds
- Cosmos: \~6 seconds

## How to Handle Delays

Follow these guidelines

1. Use the Quote endpoint to get the estimated fee.
2. Don't leave the user with a swap screen spinner, instead, move the swap to a "pending state" with a 10minute countdown. Let the user exit the app, perhaps even send them a notification after.
3. Every minute, poll Midgard and see if the swap is processed.
4. Once processed, you can inform the user, perhaps surprise them if the swap is done faster
