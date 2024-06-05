# AMM Package

While the Query package allows quick and powerful information retrieval from THORChain such as a swap estimate., this package performs the actions (sends the transactions), such as a swap, add and remove.

As a general rule, this package should be used in conjunction with the query package to first check if an action is going to be possible be performing the action.

Example: call estimateSwap first to see if the swap is going to be successful before calling doSwap, as doSwap will not check.

## Code examples in Replit

Currently implemented functions are listed below with code examples. Press the Run button to run the code and see the output. Press,`Show Files`, and select `index.ts` to see the code. Select `package.json` to see all the package dependencies. [Github link](https://github.com/xchainjs/xchainjs-lib/tree/master/packages/xchain-thorchain-amm) and [install instructions](overview.md#install-procedures).

### DoSwap

Executes a swap from a given wallet. This will result in the inbound transaction into THORChain.

DoSwap runs [EstimateSwap](query-package.md#estimate-swap) first then if successful [sends a transaction](../../concepts/sending-transactions.md) with a constructed [transaction memo](../../concepts/memos.md#overview) using a [chain client](client-packages.md). Do swap can be used with an existing xchain client implementation or a custom wallet and will return the transaction hash of the inbound transaction.

A seed is provided in the example but it has no funds so it will error.

<iframe width="100%" height="600" src="https://replit.com/@thorchain/doSwap-Single?embed=true" ></iframe>

### Savers

Adds and removed liquidity from Savers. Requires a seed with funds.

<iframe width="100%" height="600" src="https://replit.com/@thorchain/saversTs?embed=true" ></iframe>

### Add Liquidity

Adds liquidity to a pool. Provide both assets for the pool. lp type is determined from the amount of the asset. The example is a single-sided rune deposit. A seed is provided in the example but it has no funds so it will error.

<iframe width="100%" height="600" src="https://replit.com/@thorchain/addLiquidity?embed=true" ></iframe>

### Remove Liquidity

Removes Liquidity from a pool. The opposite of adding liquidity.

<iframe width="100%" height="600" src="https://replit.com/@thorchain/removeLiquidity?embed=true" ></iframe>
