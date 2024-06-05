# Streaming Swaps

Streaming Swaps is a means for a swapper to get better price execution if they are patient. This ensures Capital Efficiency while still keeping with the philosophy "impatient people pay more".

There are two important parts to streaming swaps:

1. The **interval** part of the stream allows arbs enough time to rebalance intra-swap - this means the capital demands of swaps are met throughout, instead of _after_.
2. The **quantity** part of the stream allows the swapper to reduce the size of their sub-swap so each is executed with less slip (so the total swap will be executed with less slip) _without_ losing capital to on-chain L1 fees.

If a swapper is willing to be patient, they can execute the swap with a better price, by allowing arbs to rebalance the pool between the streaming swaps.

Once all swaps are executed and the streaming swap is completed, the target token is sent to the user (minus outbound fees).

Streaming Swaps is similar to a Time Weighted Average Price (TWAP) trade however it is restricted to 24 hours (Mimir `STREAMINGSWAPMAXLENGTH = 14400` blocks).

## Using Streaming Swaps

To utilise a streaming swap, use the following within a [Memo](../concepts/memos.md#swap):

Trade Target or Limit / Swap Interval / Swap Quantity.

- **Limit** or Trade Target: Uses the trade limit to set the maximum asset ratio at which a mini-swap can occur; otherwise, a refund is issued.
- **Interval**: Block separation of each swap. For example, a value of 10 means a mini-swap is performed every 10 blocks.
- **Quantity**: The number of swaps to be conducted. If set to 0, the network will determine the appropriate quantity.

Using the values Limit/10/5 would conduct five mini-swaps with a block separation of 10. Only swaps that achieve the specified asset ratio (defined by Limit) will be performed, while others will result in a refund.

On each swap attempt, the network will track how much (in funds) failed to swap and how much was successful. After all swap attempts are made (specified by "swap quantity"), the network will send out all successfully swapped value, and the remaining source asset via refund (that failed to swap for some reason, most likely due to the trade target).

If the first swap attempt fails for some reason, the entire streaming swap is refunded and no further attempts will be made. If the `swap quantity` is set to zero, the network will determine the number of swaps on its own with a focus on the lowest fees and maximize the number of trades.

## Minimum Swap Size

A min swap size is placed on the network for streaming swaps (Mimir `StreamingSwapMinBPFee = 5` Basis Points). This is the minimum slip for each individual swap within a streaming swap allowed. This also puts a cap on the number of swaps in a streaming swap. This allows the network to be more friendly to large trades, while also keeping revenues up for small or medium-sized trades.

### Calculate Optimal Swap

The network works out the optimal streaming swap solution based on the Mimumn Swap Size and the swapAmount.

**Single Swap**: To calculate the minimum swap size for a single swap, you take 2.5 basis points (bps) of the depth of the pool. The formula is as follows:

$$
{MinimumSwapSize} = MinBPStreamingSwap * Rune Pool Depth
$$

Example using BTC Pool:

- BTC Rune Depth = 20,007,476 RUNE
- StreamingSwapMinBPFee = 5 bp

MinimumSwapSize = 0.0005 \* 20,007,476 = 10,003. RUNE

**Double Swap**: When dealing with two pools of arbitrary depths and aiming for a precise 5 bps swap fee (set by `StreamingSwapMinBPFee`), you need to create a virtual pool size called `runeDepth` using the following formula:

$$
virtualRuneDepth  =(2*r1*r2) / (r1+r2)
$$

`r1` represents the rune depth of pool1, and `r2` represents the rune depth of pool2.

The `runeDepth` is then used with 1.25 bps (half of 2.5 bps since there are two swaps), which gives you the minimum swap size that results in a 5 bps swap fee.

$$
{MinimumSwapSize} = (MinBPStreamingSwap / 2) * virtualRuneDepth
$$

```admonish success
The larger the difference between the pools, the more the virtual pool skews towards the smaller pool. This results in less rewards given to the larger pool, and more rewards given to the smaller pool.
```

Example using BTC and ETH Pool

- BTC Rune Depth = 20,007,476 RUNE
- ETH Rune Depth = 8,870,648 RUNE
- StreamingSwapMinBPFee = 5 bp

virtualRuneDepth = (2\*20,007,476\*8,870,648) / (20,007,476 + 8,870,648) = 12,291,607 RUNE

MinimumSwapSize = (0.0005/4) \* 12,291,607 = 1536.45 RUNE

### Swap Count

The number of swaps required is determined by dividing the `swap Amount` by the minimum swap size calculated in the previous step.

$$
swapCount = swapAmount / MinimumSwapSize
$$

The `swapAmount` represents the total amount to be swapped.

Example: swap 20,000 RUNE worth of BTC to ETH. (approx 0.653 BTC).

20,000 / 3,072.90 = 6.5 = 7 Swaps.

### Comparing Price Execution

The difference between streaming swaps and non-streaming swaps can be calculated using the swap count with the following formula:

$$
difference = (swapCount - 1) / swapCount
$$

The `difference`value represents the percentage of the swap fee saved compared to doing the same swap with a regular fee structure. There higher the swapCount, the bigger the difference.

Example:

- (7-1)/7 = 6/7 = 85% better price execution by being patient.
