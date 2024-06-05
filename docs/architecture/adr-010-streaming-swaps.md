# ADR 010: Introduction of Streaming Swaps

## Changelog

- Initial commit: July 23, 2023

## Status

Launching Feature

## Context

The current model of network swap fees and price execution on THORChain is
directly proportional to the depth of the pool. Larger trades lead to higher
fees and consequently, less favorable price execution. This has resulted in a
trend where approximately 99% of swaps on THORChain are under $10k in value,
as users executing larger trades often find better price execution on other
exchanges, typically centralized ones (CEXs). If one does market analysis, one
would see that whales control the majority of the spot market which is largely
unavailable to THORChain.

To capture a larger market share of trading, the network needs to offer more
competitive price execution, particularly for larger (whale) trades.

## Proposed Change

This ADR introduces an enhancement to swaps, enabling users to optionally
divide larger trades into several autonomous smaller trades. This division
allows arbitrage bots to adjust the price multiple times during the swap
process.

It is worth make a note that while this will reduce swap fees, it will not
reduce other fees such as gas fees, outbound fees, and affiliate fees.

For a comprehensive description of this feature, please refer to this [GitLab
issue](https://gitlab.com/thorchain/thornode/-/issues/1514).

As of the date of this writing, the feature has been deployed on our stagenet
for several weeks and undergone extensive testing by our developers. The
testing document is accessible
[here](https://docs.google.com/document/d/1QMHtYi-pH0Ie4i3QKCBDaekweC8IR9_zPpYX1vhxS4E/edit#heading=h.bcxnyhl6rm1c).

### Advantages

This proposal offers substantial benefits for the network and its users:

1. **Improved Price Execution:** This feature allows the network to determine
   any swap fee (in basis points) for trades of all sizes, between any two
   supported assets, whilst maintaining the benefits and safeguards that the
   slip-based fee model provides. This flexibility allows the community to choose
   its level of competitiveness against other exchanges (CEXs or DEXs).

2. **Increased Trade Volume and User Base:** The improved price execution
   should lead to a significant increase in trade volume, a rise in unique
   swappers, and an expansion of our market share.

3. **Enhanced Capital Efficiency:** Due to swaps being spread out over time,
   each swap will affect the pool price less, causing the AMM to become
   significantly more capital efficient.

4. **Support for New Trading Strategies:** THORChain will be able to support
   new trading strategies, such as time-weighted average price (TWAP) and
   dollar-cost averaging (DCA), further extending its user base and community.

5. **Enhanced Value Proposition:** Other THORChain features can also leverage
   streaming swaps to enhance their value propositions. For example, cheaper
   entry and exit for savers, and order books with partial fulfillment and better
   price execution.

### Potential Drawbacks

This feature may lead to the network collecting fewer fees per swap, depending
on the trade size. Although there may be a decrease in system income, it is
anticipated that the increase in trade volume and number of swappers will
compensate for this. This proposal represents a shift in THORChain's
priorities from profitability to increased adoption and growth (in the short
to medium term).

Another potential issue is with the increase in trade volume into the network
(but outside of the pools) might result in an unbounded amount of liquidity
that the network is securing. While this has always been true for the network,
trades have traditionally been "instant", allowing funds to leave the network
as quickly as they entered. Now, the network will permit trade value to remain
in the network for up to 24 hours, which is adjustable via mimir.
