# ADR 012: Scale Lending

## Changelog

- 2021-01-30: Proposed

## Context

### Lending

Lending was launched in Q3 2023. As of writing:

1. 1.3k loans opened with $24m in collateral from ~700 borrowers, to issue $7m in debt
2. 3.68m RUNE burnt from collateral, with 2m RUNE minted to create debt
3. 88% percent "full" based on the lending lever
4. System Risk extremely low, with Full Closure Scenario resulting in 700k RUNE burnt forever (safe)
5. No observed attacks on the Lending Lever and/or TOR Anchor pool prices

Lending has been shown to be safe, controlled and has almost reached caps. It is time to Scale Lending.

https://flipsidecrypto.xyz/banbannard/⚡-thor-chain-lending-thorchain-lending-fOAKej

https://dashboards.ninerealms.com/#lending

### Standby Reserve

The Standby Reserve has 60,000,000 RUNE held in idle, waiting to be deployed. This was allocated since Genesis to be used by the protocol when it was ready for it.

https://runescan.io/address/thor1lj62pg6ryxv2htekqx04nv7wd3g98qf9gfvamy

## Current Params

```text
"MAXCR": 50000, //500%
"MINCR": 20000, //200%
```

## Proposal

To scale Lending, reduce risk and increase Lending Lever safely the following is proposed:

1. Set `"MAXCR": 20000` (200%)
2. Burn the 60,000,000 Standby Reserve

### MAXCR

Currently the CR for BTC is at maximum (500%) whilst the CR for ETH is at ~300%.

By reducing the Max Collaterisation Ratio down to 200% the following is achieved

1. More loans will be opened since terms are favourable
2. Full Closure Scenario is less likely to net-mint RUNE because the debt-collateral ratio is 2x, not 5x
3. More loans are unlikely to be closed in bear markets because $collateral is likely to be less than $debt (since only has to fall 2x instead of 5x)

### Burn Reserve

By burning the 60m Standby Reserve, the Lending Lever is:

1. Scaled by 5x (15m to 75m)
2. Safer for current borrowers since there is more buffer to fill before the RUNE supply circuit breaker is hit at 500m supply
3. Safer for the protocol by taking the 60m RESERVE out of possible circulation (if it was added to MAIN Reserve)

## Other Considerations

The Standby Reserve has always been a "Plan B" for the protocol, a joker up its sleeve to play. It's time to play it. By burning it, it makes all outstanding RUNE notionally more valuable.

Burning it also frees up space for the protocol to Scale Lending (which is a successful feature), as well as space to launch future features:

1. Stablecoin (will require space under the cap to mint)
2. Perpetuals (will use derived assets, so needs space)

## References

- [ADR-011](https://gitlab.com/thorchain/thornode/-/blob/develop/docs/architecture/adr-011-lending.md)
