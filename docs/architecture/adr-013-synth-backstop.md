# ADR 013: Synth Backstop

## Changelog

- 2024-04-07: Proposed

## Context

### Synthetic Backstop

The synth code has been in the codebase since before MCCN, but took a year to flip on. Its initial purpose was to give arb bots a much more efficient way to arb the pools without going through exogenous blockchains. According to an analysis from Delphi Digital, these assets had 16x the efficiency of other pool assets.

Later, the community opted to utilize synths for "savers". This feature has effectively doubled the depth of the pools, doubled the volume, and created significant buy pressure for the $RUNE asset. Savers created a more significant demand center for the synthetic asset than originally considered. This increase in synths was beneficial for the protocol, but also increased the risk to the protocol. The community increased the value of synths and its risk once stable savers were added.

It has come time for the community to implement a backstop to synths to curtail any large-scale risk to the protocol. Nodes are asked to vote over the the next 1 week over two proposed solutions. The solution with the most votes will win.

#### Solution 1 - Forced Ejection

In this design, once the synth utilization grows above the cap (currently 60%), the protocol will force savers to exit. Likely on a "last in, first out" system (in interest to reward our most loyal community members). Savers would be ejected until the synth utilization drops below the cap.

This approach is the most "user friendly", but creates $RUNE sell pressure in a down market while also reducing the pool depths (and trade volume).

#### Solution 2 - Savers Haircut

In this design, once the synth utilization grows above the cap, any withdrawal of savers would experience a haircut relative to how far over the synth cap the pool is. This would effectively be an "exit fee" that users opt into when they choose to withdraw. This caps the risk that LPs can experience from a down $RUNE price and moves it towards the savers.

This approach is most beneficial to the protocol as it doesn't sell $RUNE and keeps its pools deep, but may cause a PR problem around the savers experience etc.
