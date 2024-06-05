# Query Package

This package is designed to obtain any desired information from THORChain that a developer could want. While it uses Midgard and Thornode packages it is much more powerful as it knows where to get the best information and how to package the information for easy consumption by developers or higher functions.

It exposes several simple functions that implement all of THORChain's complexity to allow easy information retrieval. The Query package does not perform any actions on THORChain or send transactions, that is the job of the [Thorchain-AMM package](amm-package.md).

## Code examples in Replit

Currently implemented functions are listed below with code examples. Press the Run button to run the code and see the output. Press Show files, and select index.ts to see the code. Select package.json to see all the package dependencies. [Repo link](https://github.com/xchainjs/xchainjs-lib/tree/master/packages/xchain-thorchain-query) and [install instructions](overview.md#install-procedures).

### Estimate Swap

Provides estimated swap information for a single or double swap for any given two assets within THORChain. Designed to be used by interfaces, see more info [here](coding-guide.md#query). EstimateSwap will do the following:

- Validate swap inputs
- Checks for [network or chain halts](../../concepts/network-halts.md)
- Get the latest pool data from [Midgard](../../concepts/connecting-to-thorchain.md#midgard)
- Work out the swap [slip](../../concepts/math.md#slippage), swap [fee](../../concepts/fees.md#fee-ordering-for-swaps) and [output](../../concepts/math.md#swap-output)
- Deducts all [fees](../../concepts/fees.md#overview) from the input amount (inbound, swap, outbound and any affiliate) in the correct order to produce `netOutput` and detail fees in `totalFees`
- Ensures `totalFees` is not greater than `input`.
- Work out the expected [wait time](../../concepts/delays.md#overview) including confirmation counting and outbound delay.
- Get the current [Asgard Vault address](../../concepts/querying-thorchain.md#getting-the-asgard-vault) for the inbound asset
- Advise if a swap is possible and provide a reason if it is not.

Note: This will be the best estimate using the information available, however exact values will be different depending on pool depth changes and network congestion.

<iframe width="100%" height="600" src="https://replit.com/@thorchain/estimateSwap?embed=true"></iframe>

### Savers

Shows use of the savers quote endpoints.

<iframe width="100%" height="600" src="https://replit.com/@thorchain/quoteSaversTS?embed=true"></iframe>

### Check Balance

Checks the liquidity position of a given address in a given pool. Retrieves information such as current value, ILP coverage, pool share and last added.

<iframe width="100%" height="600" src="https://replit.com/@thorchain/checkLiquidity?embed=true"></iframe>

### Check Transaction

Provide the status of a transaction given an input hash (e.g. the output of doSwap). Looks at the different stages a transaction can be in and report.

In development

### Estimate Add Liquidity

Provides an estimate for given add liquidity parameters such as slip fee, transaction fees, and expected wait time. Supports symmetrical, asymmetrical and uneven additions.

<iframe width="100%" height="600" src="https://replit.com/@thorchain/estimateWithdrawLiquidity?embed=true" ></iframe>

### Estimate Remove Liquidity

Provides information for a planned withdrawal for a given liquidity position and withdrawal percentage. Information such as fees, wait time and ILP coverage

<iframe width="100%" height="600" src="https://replit.com/@thorchain/estimateWithdrawLiquidity?embed=true" ></iframe>

### List Pools

Lists all the pool detail within THORChain.

<iframe width="100%" height="600" src="https://replit.com/@thorchain/listPools?embed=true" ></iframe>

### Network Values

List current network values from constants and mimir. If mimir override exists, it is displayed instead.

<iframe width="100%" height="600" src="https://replit.com/@thorchain/networkSetting?embed=true" ></iframe>

If there is a function you want to be added, reach out in Telegram or the dev discord server.
