# ADR 006: Enable POL

## Changelog

- February 17, 2022: Initial commit
- September 26, 2023: Add all Saver pools to PoL targets

## Status

Implemented

## Update Nov 23: Lower POL Exit Criteria

As of Nov 23, POL enters at 50% and exits at 40% (4500):

> - `POLTargetSynthPerPoolDepth` to `4500`
> - `POLBuffer` to `500`
>   PoL will enter at `4500 + 500 = 50%` but exit at `4500 - 500 = 40%`

The issue is that POL does not stay in the pools long enough to make enough yield to compensate for the Impermanent Loss experienced from the price change as Synth Utilisation drops from 50% to 40%:

```text
BlockHeight: 13,326,840
Overall RUNE deposited: 7,590,445.22 RUNE
Overall RUNE Withdrawn: 5,704,572.68 RUNE
Current RUNE PnL: -430,760.85 RUNE
```

To let PoL stay in the pools for much longer (but still exit if a pool is being removed from the network or utilisation drops off), mimir should refine PoL parameters:

- `POLTargetSynthPerPoolDepth` to `3000`
- `POLBuffer` to `2000`

PoL will enter at `3000 + 2000 = 50%` but exit at `3000 - 2000 = 10%`. This should give PoL enough time to make yield on deposits and is not losing to Impermanent Loss.

## Update Sep 23: Add All Saver Pools to PoL Targets

Recently additional pools (stablecoins) were enabled for Saver positions, but PoL was not activated on those pools. The original Pol ADR was explicit in which pools would receive PoL, but it is not wise to have Saver Pools without PoL protection. This ADR amendment sets out that all Saver pools should receive PoL treatment.

> PoL reduces dual-LP leverage and keeps Synth utilization away from Synth Caps. If synths exceed their caps, then the L1 pool has more synthetic counterparts than L1 assets, and becomes top-heavy. PoL adds L1 liquidity to prevent this.

Going forward, any pool activated for Savers should also enable PoL.
To sync the pools, the following should be set:

L1 Pools

```text
AVAX.AVAX 1
LTC.LTC 1
BCH.BCH 1
DOGE.DOGE 1
BSC.BNB 1
BNB.BNB 1
DOGE.DOGE 1
GAIA.ATOM 1
```

StableCoin (TOR Anchor) Pools

```text
POL-AVAX.USDC-0XB97EF9EF8734C71904D8002F8B6BC66DD9C48A6E 1
POL-BNB.BUSD-BD1 1
POL-ETH.USDC-0XA0B86991C6218B36C1D19D4A2E9EB0CE3606EB48 1
POL-ETH.USDT-0XDAC17F958D2EE523A2206206994597C13D831EC7 1
```

## Context

Protocol Owned Liquidity is a mechanism whereby the protocol utilizes the Protocol Reserve to deposit $RUNE asymmetrically into liquidity pools. In effect, it is taking the RUNE-side exposure in dual-sided LPs, reducing synth utilization, so that Savers Vaults can grow. Protocol Owned Liquidity may generate profit or losses to the Protocol Reserve, and care should be taken to determine the timing, assets and amount of POL that is deployed to the network.

A vote is currently underway to raise the `MAXSYNTHPERPOOLDEPTH` from `5000` to `6000`. Nodes have already been instructed that raising the vote to `6000` comes with an implicit understanding that Protocol Owned Liquidity (POL) will be activated as a result (https://discord.com/channels/838986635756044328/839001804812451873/1074682919886528542). This ADR serves to codify the exact parameters being proposed to enable POL.

## Proposed Change

- `POLTargetSynthPerPoolDepth` to `4500`: POL will continue adding RUNE to a pool until the synth depth of that pool is 45%.
- `POLBuffer` to `500`: Synth utilization must be >5% from the target synth per pool depth in order to add liquidity / remove liquidity. In this context, liquidity will be withdrawn below 40% synth utilization and deposited above 50% synth utilization.
- `POLMaxPoolMovement` to `1`: POL will move the pool price at most 0.01% in one block
- `POLMaxNetworkDeposit` to `1000000000000`: start at 10,000 RUNE, with authorization to add up to 10,000,000 RUNE on an incremental basis at developer's discretion. After 10m RUNE, a new vote must be called to further raise the `POLMaxNetworkDeposit`.
- `POL-BTC-BTC` to `1`: POL will start adding to the BTC pool immediately, as the pool has reached its synth cap at the time of publication.
- `POL-ETH-ETH` to `1`: POL will start adding to the ETH pool once it has reached the its synth cap.

The threshold for this ADR to pass are as follows, in chronological order:

- `MAXSYNTHPERPOOLDEPTH` to `6000` achieves 2/3rds node vote consensus
- If the author requests a Motion to Bypass and fewer than 16% of nodes dissent within 7 days (by setting `DISSENTPOL` to `1`)
- `ENABLEPOL` to `1` achieves 2/3rds node vote consensus

## Alternatives Considered

The pros/cons and alternatives to Protocol Owned Liquidity have been discussed on Discord ad neauseum. Check the [#economic-design](https://discord.com/channels/838986635756044328/839002361749438485) channel for discussion, as most topics have been covered there. The benefits and risks of POL are complex and cannot be summarized impartially by the author of this ADR. Get involved in the discussion and do your own research.

## References

- [GitLab Issue](https://gitlab.com/thorchain/thornode/-/issues/1342#protocol-owned-liquidity-pol)
- [GrassRoots Crypto](https://www.youtube.com/watch?v=Up2-arSzH5k)
