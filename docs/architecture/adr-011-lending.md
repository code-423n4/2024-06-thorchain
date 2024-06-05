# ADR 011: THORFi Lending Feature

## Overview

## Changelog

- 2023-07-14: Initial Outline
- 2023-08-01: Summary of Risk
- 2023-08-07: Update

## Status

- Proposed

## Decision

> This section records the decision that was made.

## Context

Lending will bring fresh exogenous capital from new users to scale liquidity (both TVL and security), and drive up the capital efficiency of the pools which increases system income and real yield.

Lending allows users to deposit collateral, then create debt at some collateralization ratio CR (collateralization ratio). The debt is always denominated in TOR (despite what asset the user receives). TOR is an algo-usd-stable used internally for lending and cannot be acquired or transferred. These loans have 0% interest, no liquidations, and no expiration. Risk is contained by caps on launch, slip-based fees when opening and closing loans, dynamic CR, and a circuit breaker on RUNE supply.

### ELI5 (explain like I'm 5)

"THORChain onboards layer1 collateral to issue TORdollar-based debt. The collateral is stored as equity (RUNE). The more collateral is taken onboard relative to the pool depths, the higher the collateralization ratio for new loans, so users want to get in first. The higher the collateralization ratio, the safer the system becomes. Since there are no liquidations and no interest rate, users are less likely to pay back their loans, which is beneficial to the protocol as the equity-value increases. Also, since the loans remove RUNE from the pools, the pools open back up for more TVL, causing THORChain to scale liquidity and security. The protocol takes on the liability of repaying collateral when the loan closes, so net-benefit is realized when RUNE appreciates faster than collateral assets."

## Detailed Design

[Detailed Design on GitLab](https://gitlab.com/thorchain/thornode/-/issues/1412)

## Economic Reasoning

### Scaling Liquidity

Due to PoL and Savers, the network is likely to max out pooled RUNE and send all yield to nodes. The protocol stops scaling until it can recruit new nodes, but this takes time. A faster mechanism to scale security (and TVL) is to attract new capital, buy out the RUNE from the pools and open up caps again. This is the lending design.
Whilst the loan was opened, it created a net reduction in RUNE in the pools, allowing more TVL to enter. It also bids on RUNE, allowing security to increase. Additionally, +ve RUNE Price action boosts the value of yields (block rewards).
This is why lending scales THORChain.

### Contracting Safely

If borrowers start paying back their loans, this reverses the above sequence of events and the system contracts. RUNE is minted and sold, causing sell pressure on the RUNE asset and an increase of the RUNE supply.
However, the CR will start dropping for new loans, until a point is hit and new loans are opened under favorable CR terms, encouraging new loans to be taken.

### Circuit Breaker

If the RUNE price drops drastically against the majority of its collateral assets (BTC, ETH), then net inflation of RUNE will occur if users start paying back their loans, and it exceeds the margin of RUNE already historically burnt. This inflation could hit the breaker limit. At this point the system pauses new loans to be created and sunsets (turns off) lending (note, all other features of TC still function). At this point, no further inflation of RUNE can occur and the supply arrested. The RESERVE will cover the remaining collateral payouts.

### Unwind Path

If the Circuit Breaker is hit, then the protocol will draw from RESERVES, until RESERVES are zero’d. After this, the protocol will attempt to mint over 500m and fire another trigger which auto-halts the entire protocol (swaps, savers, lending, sends). At this point there is still 60m RUNE in the Standby RESERVE that can be used to:

1. Re-stock the Protocol RESERVE, then
2. Place remainder in a Claims Module, where collateral holders can pro-rata claim based on their collateral at the time of the halt, paid out in RUNE.
3. After some time, turn off the Claims module, and redraw unclaimed RUNE back to Standby.

Important to note, the above Unwind Path is not coded, and will require governance to approve and devs to implement. It is just a suggestion for the Community to know there is an Unwind Path possible that does not result in destruction of the entire network.

### Proactive Dynamic Burn

Currently, all derived asset slip fees are burnt and this presents a continuous burn on RUNE supply.

A suggestion by devs to prevent the Circuit Breaker ever being hit is a dynamic burn of System Income, removing the 500m auto-halt:
At 490m Supply, start burning, linearly from 0% to 20% at 500m Supply
Thus there is an ability for protocol to “burn down” its liability over a period of time.
Since 500m can never be exceeded, the dynamic burn will allow small borrowers to exit as soon as space for them to mint is opened back up.
Large borrowers will be the last to leave.
The protocol will keep burning down to 490m, and at some point, market sentiments will change and Borrowers will come back to THORChain.

An additional suggestion is to build more features (such as Perps) which will add more burning into the protocol (via more derived asset pool fees and liquidations).

### Loan Caps

For safe scaling, lending is capped based on the amount of outstanding RUNE supply. The monetary policy of the asset is a cap of 500m. Currently, the network is missing ~15m RUNE due to individuals not upgrading original rune assets of BEP2 or ERC20 RUNE. Since we have this 15m gap, this is used to help "fund" lending to start.
To help protect the network from hitting the circuit breaker, only 1/3rd of the outstanding RUNE will be used to guide scaling lending (~5m RUNE) in terms of loan collateral value. 1/3rd is selected due to historical analysis of RUNE's price movements, where its largest micro-economic price event (several hacks in early beta days, 2 years ago), we saw a 3x price change downward (relative to BTC). This means that however amount of loans that are taken out, we could see a 3x downward price change AND 100% of loans could close, and we still wouldn't hit the 500m cap (circuit breaker).
As more rune is burnt, and the gap is increased from 15m --> 20m, this opens up more space for more loans to be opened. In addition, the market buy and burning of 5m+ RUNE should cause rune's price to naturally out perform BTC. As rune's price out performance BTC, and the ratio in the BTC pool shifts, and more loans can be opened, without changing our risk profile.

## Parameter Recommendations

The following recommendation were based off of ranges from Block Science Simulation: https://hackmd.io/R_ksPSG0T6mtjsEE7U9q-w#Parameter-Recommendations

| Mimir                 | Description                                                                                                                                                | Recommendation |
| --------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------- |
| LendingLever          | Relative to the amount of missing rune from max supply. Throttles loan CR, higher is more risky, lower is less risky, but retards lending growth/adoption. | 3,333          |
| MaxCR                 | Max CR ratio when LendingLever consumed (pro-rata to pool)                                                                                                 | 50,000         |
| MinCR                 | First loan gets this CR                                                                                                                                    | 20,000         |
| LoanRepaymentMaturity | Minimum number of blocks a loan must remain open                                                                                                           | 432000         |
| MaxAnchorBlocks       | Number of blocks to sum the slip for Anchor                                                                                                                | 300            |
| MaxAnchorSlip         | Max slip value for anchor sensitivity                                                                                                                      | 6,000          |
| DerivedDepthBasisPts  | Depth of derives pools out the gate compared to L1 pool                                                                                                    | 10,000         |
| DerivedMinDepth       | Min pool depth to shrink to in periods of high volatility                                                                                                  | 1000           |

## Summary of Largest Risks

Block Science reviewed the lending mechanisms exhaustively. The Output of their research includes:

- [Risk Report](https://hackmd.io/@blockscience/H1Q-erh_n)
- [Simulation Summary](https://hackmd.io/R_ksPSG0T6mtjsEE7U9q-w)
- [CadCad simulation framework](https://gitlab.com/thorchain/misc/cadcad-thorchain/-/tree/main/documentation)

### Below are a summary of the largest risks identified

Block Science reviewed the lending mechanisms exhaustively. The Output of their research includes:

Risk Report: https://hackmd.io/@blockscience/H1Q-erh_n
Simulation Summary: https://hackmd.io/R_ksPSG0T6mtjsEE7U9q-w
CadCad simulation framework: https://gitlab.com/thorchain/misc/cadcad-thorchain/-/tree/main/documentation

## Testing Guide

Lending has been tested for success and failure modes on Stagenet. This can be used by devs to review for feature readiness.  
[Testing Document](https://docs.google.com/document/d/1-kOHRk-P1ooJHRzkh37o63LftsXu83rAOjEysLKJVeo/edit)

[Regression Tests](https://gitlab.com/thorchain/thornode/-/tree/develop/test/regression/suites/lending)

## References

- [Issue on GitLab](https://gitlab.com/thorchain/thornode/-/issues/1412)
- [PR on GitLab](https://gitlab.com/thorchain/thornode/-/merge_requests/2713)
- [Block Science Risk Report](https://hackmd.io/@blockscience/H1Q-erh_n)
- [Block Science Simulation Summary](https://hackmd.io/R_ksPSG0T6mtjsEE7U9)
