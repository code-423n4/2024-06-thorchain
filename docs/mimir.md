# Constants and Mimir

## Overview

The network launched with a set number of constants, which has not changed. Constants can be overridden via Mimir and nodes have the ability to [vote on](../thornodes/overview.md#node-voting) and change Mimir values. \
See [Halt Management](./concepts/network-halts.md) for halt and pause specific settings.

Mimir setting can be created and changed without a corresponding Constant.

### Values

- Constant Values:[https://midgard.ninerealms.com/v2/thorchain/constants](https://thornode.ninerealms.com/thorchain/constants)
- Mimir Values: [https://thornode.ninerealms.com/thorchain/mimir](https://thornode.ninerealms.com/thorchain/mimir)

### Key

- No Star or Hash - Constant only, no Mimir override.
- Star (\*) indicates a Mimir override of a Constant
- Hash (#) indicates Mimir with no Constant.

## Outbound Transactions

- `OutboundTransactionFee`: Amount of rune to withhold on all outbound transactions (1e8 notation)
- `RescheduleCoalesceBlocks`\*: The number of blocks to coalesce rescheduled outbounds
- `MaxOutboundAttempts`: The maximum retries to reschedule a transaction

### Scheduled Outbound

- `MaxTxOutOffset`: Maximum number of blocks a scheduled outbound transaction can be delayed
- `MinTxOutVolumeThreshold`: Quantity of outbound value (in 1e8 rune) in a block before it's considered "full" and additional value is pushed into the next block
- `TxOutDelayMax`: Maximum number of blocks a scheduled transaction can be delayed
- `TxOutDelayRate`\*: Rate of which scheduled transactions are delayed

## Swapping

- `HaltTrading`#: Pause all trading
- `Halt<chain>Trading`#: Pause trading on a specific chain
- `MaxSwapsPerBlock`: Artificial limit on the number of swaps that a single block with process
- `MinSwapsPerBlock`: Process all swaps if the queue is equal to or smaller than this number
- `EnableDerivedAssets`: Enable/disable derived asset swapping (excludes lending)
- `StreamingSwapMinBPFee`\*: Minimum swap fee (in basis points) for a streaming swap trade
- `StreamingSwapMaxLength`: Maximum number of blocks a streaming swap can trade for
- `StreamingSwapMaxLengthNative`: Maximum number of blocks native streaming swaps can trade over
- `TradeAccountsEnabled`: Enable/disable trade account

## Synths

- `MaxSynthPerAssetDepth`: The amount of synths allowed per pool relative to the pool depth
- `MaxSynthPerPoolDepth`\*: The percentage (in basis points) of how many synths are allowed relative to pool depth of the related pool
- `BurnSynths`#: Enable/Disable burning synths
- `MintSynths`#: Enable/Disable minting synths
- `VirtualMultSynths`: The amount of increase the pool depths for calculating swap fees of synths

## Savers

- `MaxSynthsForSaversYield`\*: The percentage (in basis points) synth per pool where synth yield reaches 0%
- `SaversStreamingSwapsInterval`\*: For Savers deposits and withdraws, the streaming swaps interval to use for the Native <> Synth swap

### POL Management

- `POL-<Asset-Asset>`#: Enabled POL for that pool. E.g. `POL-BTC-BTC" = 1` enabled POL for the BTC pool.
- `POLBuffer`\*: the buffer around the POL synth utilization (basis points added to/subtracted from POLTargetSynthPerPoolDepth basis points)
- `POLMaxNetworkDeposit`\*: Maximum amount of rune deposited into the pools
- `POLMaxPoolMovement`\*: Maximum amount of rune to enter/exit a pool per iteration - 1 equals one hundredth of a basis point of pool rune depth
- `POLTargetSynthPerPoolDepth`\*: The target synth per pool depth for POL (basis points)

## Lending

- `LendingLever`: Controls (in basis points) how much lending is allowed relative to rune supply
- `LoanRepaymentMaturity`: Number of blocks before loan has reached maturity and can be repaid
- `MinCR`\*: Minimum collateralization ratio (basis pts)
- `MaxCR`\*: Maximum collateralization ratio (basis pts)
- `Lending-THOR-<Asset>`#: Lending key for an asset, allows that Asset to be used as colloteral. The lending key for the `ETH.ETH` pool would be `LENDING-THOR-ETH` and enabled the `THOR-ETH` virtual pool.
- `LoanStreamingSwapsInterval`\*: The block interval between each streaming swap of opening or closing a loan

## Derived Assets

- `DerivedDepthBasisPts`: Allows mimir to increase or decrease the default derived asset pool depth relative to the anchor pools. 10k == 1x, 20k == 2x, 5k == 0.5x
- `DerivedMinDepth`: Sets the minimum derived asset depth in basis points, or pool depth floor.
- `MaxAnchorSlip`\*: Percentage (in basis points) of how much price slip in the anchor pools will cause the derived asset pool depths to decrease to
- `DerivedMinDepth`. For example, 8k basis pts will mean that when there has been 80% price slip in the last `MaxAnchorBlocks`, the derived asset pool depth will be `DerivedMinDepth`. So this controls the "reactiveness" of the derived asset pool to the layer1 trade volume.
- `MaxAnchorBlocks`: Number of blocks that are summed to get total pool slip. This is the number used to be applied to `MaxAnchorSlip`
- `TORAnchor-<Asset>`#: Enables an asset to be used in the TOR price calculation

## LP Management

- `PauseLP`#: Pauses the ability for LPs to add/remove liquidity
- `PauseLP<chain>`#: Pauses the ability for LPs to add/remove liquidity, per chain
- `MaximumLiquidityRune`#: Maximum RUNE capped on the pools known as the ‘soft cap’
- `LiquidityLockUpBlocks`: The number of blocks LP can withdraw after their liquidity
- `PendingLiquidityAgeLimit`: The number of blocks the network waits before initiating pending liquidity cleanup. Cleanup of all pools lasts for the same duration.

## Chain Management

- `HaltChainGlobal`#: Pause observations on all chains (chain clients)
- `HaltTrading`: Stops swaps and additions, if done, will result in refunds. Observations still occur.
- `Halt<chain>Chain`#: Pause a specific blockchain
- `NodePauseChainGlobal`#: Individual node controlled means to pause all chains
- `NodePauseChainBlocks`: Number of block a node operator can pause/resume the chains for
- `BlocksPerYear`: Blocks in a year
- `MaxUTXOsToSpend`#: Max UTXOs to be spent in one block
- `MinimumNodesForBFT`: Minimum node count to keep the network running. Below this, Ragnarök is performed

### Fee Management

- `NativeTransactionFee`: RUNE fee on all on chain txs
- `TNSRegisterFee`: Registration fee for new THORName, in RUNE
- `TNSFeeOnSale`: fee for TNS sale in basis points
- `TNSFeePerBlock`: per block cost for TNS, in RUNE
- `PreferredAssetOutboundFeeMultiplier`\*: The multiplier of the current preferred asset outbound fee, if RUNE balance > multiplier \* outbound_fee, a preferred asset swap is triggered
- `MinOutboundFeeMultiplierBasisPoints`\*: Minimum multiplier applied to base outbound fee charged to user, in basis points
- `MaxOutboundFeeMultiplierBasisPoints`: Maximum multiplier applied to base outbound fee charged to user, in basis points

### Solvency Checker

- `StopSolvencyCheck`#: Enable/Disable Solvency Checker
- `StopSolvencyCheck<chain>`#: Enable/Disable Solvency Checker, per chain
- `PermittedSolvencyGap`: The amount of funds permitted to be "insolvent". This gives the network a little bit of "wiggle room" for margin of error

## Node Management

- `MinimumBondInRune`\*: Sets a lower bound on bond for a node to be considered to be churned in
- `ValidatorMaxRewardRatio`\*: the ratio to MinimumBondInRune at which validators stop receiving rewards proportional to their bond
- `MaxBondProviders`\*: Maximum number of bond providers per mode
- `NodeOperatorFee`\*: Minimum node operator fee
- `SignerConcurrency`\*: Number of concurrent signers for active and retiring vaults

### Yggdrasil Management

```admonish note
**Yggdrasil** Vaults are deprecated since [ADR-002](./architecture/adr-002-removeyggvaults.md).
```

- `YggFundLimit`: Funding limit for yggdrasil vaults (percentage)
- `YggFundRetry`\*: Number of blocks to wait before attempting to fund a yggdrasil again
- `StopFundYggdrasil`#: Enable/Disable yggdrasil funding
- `ObservationDelayFlexibility`\*: Number of blocks of flexibility for a validator to get their slash points taken off for making an observation
- `PoolDepthForYggFundingMin`\*: The minimum pool depth in RUNE required for ygg funding
- `MinimumNodesForYggdrasil`: No yggdrasil vaults if THORNode have less than 6 active nodes

### Slashing Management

- `LackOfObservationPenalty`: Add two slash points for each block where a node does not observe
- `SigningTransactionPeriod`: How many blocks before a request to sign a tx by yggdrasil pool, is counted as delinquent.
- `DoubleSignMaxAge`: Number of blocks to limit double signing a block
- `FailKeygenSlashPoints`: Slash for 720 blocks, which equals 1 hour
- `FailKeysignSlashPoints`: Slash for 2 blocks
- `ObserveSlashPoints`: The number of slashpoints for making an observation (redeems later if observation reaches consensus)
- `ObservationDelayFlexibility`: Number of blocks of flexibility for a validator to get their slash points taken off for making an observation
- `JailTimeKeygen`: Blocks a node account is jailed for failing to keygen. DO NOT drop below TSS timeout
- `JailTimeKeysign`: Blocks a node account is jailed for failing to keysign. DO NOT drop below TSS timeout

### Churning

- `AsgardSize`\*: Defines the number of members to an Asgard vault
- `MinSlashPointsForBadValidator`: Minimum quantity of slash points needed to be considered "bad" and be marked for churn out
- `BondLockupPeriod`: Lockout period that a node must wait before being allowed to unbond
- `ChurnInterval`\*: Number of blocks between each churn
- `HaltChurning`: Pause churning
- `DesiredValidatorSet`: Maximum number of validators
- `FundMigrationInterval`\*: Number of blocks between attempts to migrate funds between asgard vaults during a migration
- `NumberOfNewNodesPerChurn`#: Number of targeted additional nodes added to the validator set each churn
- `BadValidatorRedline`\*: Redline multiplier to find a multitude of bad actors
- `LowBondValidatorRate`: Rate to mark a validator to be rotated out for low bond
- `MaxNodeToChurnOutForLowVersion`\*: Maximum number of validators to churn out for low version each churn

## Economics

- `EmissionCurve`\*: How quickly rune is emitted from the reserve in block rewards
- `IncentiveCurve`\*: The split between nodes and LPs while the balance is optimal
- `MaxAvailablePools`: Maximum number of pools allowed on the network. Gas pools (native pools) are excluded from this.
- `MinRunePoolDepth`\*: Minimum number of RUNE to be considered to become active
- `PoolCycle`\*: Number of blocks the network will churn the pools (add/remove new available pools)
- `StagedPoolCost`: Number of RUNE (1e8 notation) that a stage pool is deducted on each pool cycle.
- `KillSwitchStart`\*: Block height to start to kill BEP2 and ERC20 RUNE
- `KillSwitchDuration`: Duration (in blocks) until switching is deprecated
- `MinimumPoolLiquidityFee`: Minimum liquidity fee an active pool should accumulate to avoid being demoted, set to 0 to disable demote pool based on liquidity fee
- `MaxRuneSupply`\*: Maximum supply of RUNE

## Miscellaneous

- `DollarsPerRune`: Manual override of number of dollars per one RUNE. Used for metrics data collection and RUNE calculation from MinimumL1OutboundFeeUSD
- `THORNames`: Enable/Disable THORNames
- `TNSRegisterFee`: TNS registration fee of new names
- `TNSFeePerBlock`: TNS cost per block to retain ownership of a name
- `ArtificialRagnarokBlockHeight`: Triggers a chain shutdown and ragnarok
- `NativeTransactionFee`: The RUNE fee for a native transaction (gas cost in 1e8 notation)
- `HALTSIGNING<chain>`#: Halt signing in a specific chain
- `HALTSIGNING#`: Halt signing globally
- `Ragnarok-<Asset>`#: Ragnaroks a specific pool

### Router Upgrading (DO NOT TOUCH!)

#### Old keys (pre 1.94.0)

- `MimirRecallFund`: Recalls Chain funds, typically used for router upgrades only
- `MimirUpgradeContract`: Upgrades contract, typically used for router upgrades only

#### New keys (1.94.0 and on)

- `MimirRecallFund<CHAIN>`: Recalls Chain funds, typically used for router upgrades only
- `MimirUpgradeContract<CHAIN>`: Upgrades contract, typically used for router upgrades only
