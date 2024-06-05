# Swapper Clout

Swapper Clout allows traders to have immediate and faster swaps while maintaining the security features of [delayed outbounds](./delays.md#outbound-delay).

Clout is established by fees paid (in RUNE) for swaps/trades. The more fees paid, the higher the clout. Swappers with a high clout have proven themselves to be highly aligned to the project and therefore can reap the rewards by getting faster trade execution. The higher the clout, the less their outbound transactions are delayed or removed entirely if their score is high enough.

By reducing traders's outbound delay time, Swapper clout reduces the outbound queue allowing for a better UX for normal users.

## Implementation Detail

For each swap, the fees paid (in RUNE) are divided by 2, and associated with the sender and recipient addresses. This number only increases. When/if the feature is implemented, devs will do a historical chain analysis and create every address's initial clout score, based on historical trades.

When the outbound delay is calculated, the clout score is subtracted from the outbound value (in RUNE), causing the delay amount to be reduced (maybe even eliminated). If there is already a scheduled outbound with the same address, the value of the clout applied reduces (or removes) the clout applied to this outbound txn (increasing delay). This is to ensure that clout is collectively applied to all current outbounds, not on a one-by-one basis. This is to ensure that an individual cannot have a clout score of 100 RUNE, and make infinite zero-delay swaps of 100 RUNE value simultaneously.

To calculate delay:

- ov: Outbound transaction value in RUNE
- c: Current address clout score
- cu: Current Clout Utilisation

$$
\text{delay} = delaycalc(ov - (c - cu))
$$

This feature rewards power users and arbitrage bots to facilitate a better trading experience, potentially enabling near-instant trade execution. Specifically, it aids arbitrage bots in operating more efficiently and managing pool prices effectively.
Since the majority of trade volume originates from arbitrage bots and power users, their outbound transactions do not contribute to the delay experienced by infrequent organic traders. This means that even a one-time trade from a known address will experience significantly reduced delays. However, delays will persist for individuals conducting large swaps from an unknown address.

More information in the original [MR](https://gitlab.com/thorchain/thornode/-/issues/1723)
