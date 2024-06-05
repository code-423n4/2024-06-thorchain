# Aggregator Overview

## Overview

THORChain will only support a set number of assets and is not designed to support long tailed assets. If a user wants to swap from a long tail ERC20 asset to Bitcoin, they have to use an Ethereum AMM like Sushi Swap to swap the ERC20 asset to ETH then they can swap the ETH to BTC.;

The same process applies for long tail tokens on other chains such as Avalanche and Cosmos.;

**Aggregator** is the ability for a user swap long tail assets via leveraging a supported on-chain AMMs and THORChain in one transaction.;

To support cross-chain aggregation, THORChain whitelists [aggregator contracts](https://gitlab.com/thorchain/thornode/-/blob/develop/x/thorchain/aggregators/dex_mainnet.go) that can call into THORChain (**Swap In)**, or receive calls (**Swap Out**). Chains that do not have on-chain AMMs (like Bitcoin) cannot support **SwapIn**, but they can support **SwapOut**, since they can pass a memo to THORChain.;

ETH swap contracts such as Sushi Swap to convert to/from THORChain support L1 tokens such as BTC. Example, in one transaction:

1. User swaps long-tail ERC20 to ETH in SushiSwap, then swaps that ETH to BTC.
2. User swaps BTC into ETH, then swaps that ETH into long-tail ERC20

There can be multiple `aggregators`. The first `thorchain aggregator` will use Sushiswap only and use ETH as the base asset. Aggregators need to follow a spec for compatibility with THORChain. Any THORChain ecosystem project can launch their own aggregator and get it whitelisted into THORChain. They can add custom/exotic routing logic if they wish.

```admonish warning
Destination addresses should only be user-controlled addresses, not smart contract addresses.;
```

### SwapIn

The SwapIn is called by the User, which then passes a memo to THORChain to do the final swap.;

`User -> Call Into Aggregator -> Swap Via AMM -> Deposit into THORChain -> Swap to Base Asset`

Eg: Swap long tail ERC20 via Sushiswap into BTC on THORChain.;

[Transaction Example](https://etherscan.io/tx/0x7905c41daaa214fbb3bad79ef63bb69aafcb15147f53cd9cf621d4049c2cea4d) using [UniSwap](https://etherscan.io/address/0x86904eb2b3c743400d03f929f2246efa80b91215) to swap ETH.ENJ to BNB.BNB.

### SwapOut

The SwapOut is called by the User invoking the aggregator memo on THORChain.;

The User needs to pass the aggregator contract address in the memo. THORChain will perform the swap to the preferred Base Asset for that chain. The rest of the parameters, being `to, asset, limit` are what is passed by THORChain in the SwapOut call for further execution.

`User -> Deposit into THORChain -> Swap to Base Asset -> Call into Aggregator -> Swap Via AMM`

Eg: Swap from BTC on THORChain to long tail ERC20 via Sushiswap. See [Memos](memos.md).;

### Combined

A user can combine the two. Swapping In first, then passing an Aggregator Memo to THORChain. This will cause THORChain to perform a **SwapOut**.

`User -> Swap In -> THORChain -> Swap Out`

Eg: Swap long tail ERC20 via Sushiswap into ETH on THORChain to LUNA then long tail CW20 via TerraSwap.

### [EVM Implementation](evm-implementation.md)

### CosmWasm Implementation

For **SwapIn** The caller must first execute a `MsgExecuteContract`, then call a `MsgSend` into THORChain vaults with the correct memo.;

For **SwapOut** THORChain will execute a `MsgExecuteContract` which then sends the final asset to the user. If failed, THORChain will execute the fallback and send the member the base asset instead.

### Deploying An Aggregator

If you would like to deploy your own aggregator with your own custom logic, deploy it with the principles above, then submit a PR for it to get whitelisted on THORChain.

Example: [https://gitlab.com/thorchain/thornode/-/merge_requests/2132](https://gitlab.com/thorchain/thornode/-/merge_requests/2132)
