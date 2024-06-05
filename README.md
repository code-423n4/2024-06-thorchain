# Thorchain audit details
- Total Prize Pool: $36,500 in USDC
  - HM awards: $28,800 in USDC
  - QA awards: $1,200 in USDC
  - Judge awards: $3,000 in USDC
  - Validator awards: $3,000 in USDC
  - Scout awards: $500 in USDC
- Join [C4 Discord](https://discord.gg/code4rena) to register
- Submit findings [using the C4 form](https://code4rena.com/contests/2024-06-thorchain/submit)
- [Read our guidelines for more details](https://docs.code4rena.com/roles/wardens)
- Starts June 5, 2024 20:00 UTC
- Ends June 12, 2024 20:00 UTC

## Automated Findings / Publicly Known Issues

The 4naly3er report can be found [here](https://github.com/code-423n4/2024-06-thorchain/blob/main/4naly3er-report.md).



_Note for C4 wardens: Anything included in this `Automated Findings / Publicly Known Issues` section is considered a publicly known issue and is ineligible for awards._
- Technically any ERC20 can interact with the THORChain Router, but there is a whitelist implemented in Bifrost that limits which tokens a tx will be processed for. The whitelist is defined [here](https://gitlab.com/thorchain/thornode/-/blob/develop/common/tokenlist/ethtokens/eth_mainnet_latest.json?ref_type=heads)

If an un-whitelisted token interacts with the Router (for example through the depositWithExpiry function), the tx will effectively be dropped on the floor, as the whitelist is checked in the scp.assetResolver function. The asset resolver + whitelist functionality can be seen [here](https://gitlab.com/thorchain/thornode/-/blob/develop/bifrost/pkg/chainclients/ethereum/ethereum_block_scanner.go#L762)

- Nothing related to [ERC20 RUNE](https://etherscan.io/address/0x3155BA85D5F96b2d030a4966AF206230e46849cb) is in scope. This asset is deprecated and cannot interact with the network anymore. 


# Overview

## THORChain Overview

THORChain is a decentralized, cross-chain liquidity network that replicates the functionality of centralized exchanges on-chain. Its features include native asset swapping, time-weighted average price (TWAP) trades, savings, and lending. The network is maintained by over 100 anonymous and independent validators, each operating a sophisticated stack to secure Layer 1 assets, monitor and sign transactions on external blockchains, and execute business logic on THORChain’s proprietary Cosmos SDK-based blockchain.

## THORChain Architecture Overview

THORChain’s Architecture consists of 4 major components. The two components relevant for this contest are the Router Smart Contract and the Bifrost observation & signing interface.

**Router Smart Contract (EVM Only) (In scope)**
https://gitlab.com/thorchain/thornode/-/tree/develop/chain/ethereum

The THORChain router serves as both the entry and exit point for EVM-based gas assets and ERC20 tokens. Users deposit assets using the router's depositWithExpiry function, providing the necessary parameters. THORChain's validators monitor the emitted events from this function to determine the appropriate actions, such as swaps, savings, or loan initiation. For EVM outbound transactions, THORChain’s validators sign and broadcast a transferOut or transferOutAndCall function call on the router. The Router contract holds the network’s ERC20 tokens and manages allowances for each active vault, while gas assets are forwarded to and from the contract and the active vaults. The Router is also used in vault "churn". Retiring vaults will call the Router's transferAllowance function to move ERC20 allowance to newly active vaults. 

**Bifrost Observation and Signing Interface (Partially in scope)**
https://gitlab.com/thorchain/thornode/-/tree/develop/bifrost

Bifrost is the observation and signing interface between THORChain and each external blockchain. For each connected blockchain, Bifrost scans each block and monitors for inbound transactions to the network. When a valid inbound transaction is observed, Bifrost posts an observation transaction to THORChain’s Layer 1. Bifrost also handles signing outbound transactions to external blockchains. Once an outbound transaction is assigned to a vault, the Bifrost daemons for each vault invoke THORChain’s TSS library to sign the outbound with its keyshare. After the threshold signature is complete, Bifrost broadcasts the transaction to the external chain.

**TSS Library (Not in scope)**
https://gitlab.com/thorchain/tss/tss-lib
https://gitlab.com/thorchain/tss/go-tss

THORChain’s Threshold Signature Scheme (TSS) library is a fork of Binance’s Go-based TSS implementation of the G20 TSS algorithm. The TSS library performs KeyGen and KeySign ceremonies to create, sign outbounds for, and rotate THORChain’s Asgard vaults, which protect the network's external Layer 1 assets.

**THORNode - Cosmos SDK Layer 1 (Not in scope)**

THORChain’s Layer 1 is a Cosmos chain that executes all business logic for the network, including feature functionality, validator management, vault rotation, rewards distribution, and more.

## Router Overview
### Public functions
***Note***: although all of the Router's functions are public, only `depositWithExpiry` is "meant" to be called by external parties. Of course, since the functions are public anyone can call them, but an external party should not be able to impersonate a vault in order to move the network's ERC20s or other assets. This is achieved through the private `_vaultAllowance` map, which keeps track of each vault's spendable ERC20s. As users call `depositWithExpiry` the target vault's allowance is incremented, and a `transferOut`, `transferOutAndCall`, or `transferAllowance` call will decrement a vault's allowance. 

**depositWithExpiry** - all EVM-based actions on THORChain are initiated through this function. This should increase the target vault's allowance (if an ERC20 is deposited), or increase the vault's ETH balance if ETH is sent in.
**Params**:
- vault (address): The THORChain vault to deposit the asset to
- asset (address): The asset being deposited (null address for ETH)
- amount (uint256): Amount of asset being deposited in asset’s decimals. Not required for ETH, transaction value used instead.
- memo (string): The THORChain memo indicating transaction intent. More information about memos here. 
- expiration (uint256): The expiration of this deposit in seconds. Transactions confirmed after this expiry should be reverted. 

**Events**:
- DepositEvent: if emitted, Bifrost will observe and post a MsgObservedTxIn transaction to THORChain to initiate the action. 

**transferOut** - intended to only be called by a THORChain vault. transferOut sends gas assets or ERC20s to fulfill an EVM outbound. 
**Params**:
- to (address): User’s destination address for the outbound
- asset (address): The asset being transferred out
- amount (uint256): The amount being transferred out
- memo (string): Outbound transaction memo. Format (OUT:<in-hash> where in-hash is the tx hash of the inbound that triggered this outbound)

**Events**:
- TransferOut: if emitted, Bifrost will observe and post a MsgObservedTxOut to THORChain to alert the network of an outbound from a vault. 

**transferOutAndCall** - (V4 version) intended to only be called by a THORChain vault. transferOutAndCall calls a whitelisted “aggregator” contract’s `swapOut` function that in turn initiates a transaction on a 3rd party protocol. This is used to daisy chain swaps with external liquidity, enabling swaps such as $DOGE -> $SHIB (DOGE → ETH on THORChain, and ETH → SHIB on Uniswap). 
**Params**:
- aggregator (address) - The address of the whitelisted aggregator
- finalToken (address) - The final token of the aggregator call (e.g. SHIB)
- to (address) - User’s destination address
- amountOutMin (uint256) - Limit for the final swap 
- memo (string) - THORChain’s outbound memo 

**transferOutAndCallV5** - intended to only be called by a THORChain vault. transferOutAndCall calls a whitelisted “aggregator” contract’s `swapOut` function that in turn initiates a transaction on a 3rd party protocol. This is used to daisy chain swaps with external liquidity, enabling swaps such as $DOGE -> $SHIB (DOGE → ETH on THORChain, and ETH → SHIB on Uniswap). 
**Params**:
- aggregator (address) - The address of the whitelisted aggregator
- finalToken (address) - The final token of the aggregator call (e.g. SHIB)
- to (address) - User’s destination address
- amountOutMin (uint256) - Limit for the final swap 
- memo (string) - THORChain’s outbound memo
- payload (bytes) - arbitrary bytes passed to aggregator contract
- originAddress (string)(optional) - source address of the swap 

**Events**:
- TransferOutAndCall - if emitted, Bifrost will observe and post a MsgObservedTxOut to THORChain to alert the network of an outbound from a vault.

**transferAllowance** - intended to only be called by a THORChain vault. transferAllowance is used during “churns” (i.e. vault rotation). Retiring vaults transfer their allowance to spend ERC20 tokens on the router to the new vaults. 
Params:
- router (address) - The address of the router where ERC20s are stored
- newVault (address) - The address of the new vault
- asset (address) - The asset whose allowance is being transferred
- amount (uint256) - Amount of allowance being transferred
- memo (string) - THORChain’s migrate memo (MIGRATE:<thorchain-block-height>)
Events:
- TransferAllowance - if emitted, Bifrost will observe abd post a MsgMigrate transaction to THORChain to alert the network of a vault migration transaction. 

**_transferOutV5** - same functionality as transferOut, but without reentrancy protection so it can be used in the batchTransferOut function

**transferOutV5** - non-reentrant wrapper around _transferOutV5

**batchTransferOutV5** - intended to only be called by a THORChain vault. Batch an array of transferOutV5 calls. 



## Bifrost Overview
The most important part of Bifrost for this audit contest is the smartcontract_log_parser, and specifically the GetTxInItem function. The GetTxInItem is run for each obvserve transaction on EVM chains. The function iterates through each log emitted by the transaction, and determines if any valid logs were emitted by the THORChain router. If it determines that a log was emitted by the THORChain router, it parses that log and forwards the appropriate details to the THORChain Layer 1 for processing. This is an incredibly sensitive and crucial piece of the codebase, as if it can be tricked, this means malicious contracts can trick Bifrost into informing THORChain of something that didn’t actually happen, leading to loss of funds. In fact, in July of 2021, THORChain’s old router + Bifrost were hacked in this exact way. Details of those hacks:

https://rekt.news/thorchain-rekt/ 
https://rekt.news/thorchain-rekt2/

## Whitelist Removal
After the July 2021 hacks, THORChain implemented a whitelist for smart contracts that are able to interact with the Router - this was implemented in the Bifrost layer (txs to contracts not whitelisted would be ignored). THORChain is now removing this whitelist, so an important focus of this audit contest is if removing the whitelist will open any vulnerabilities. Wardens will need to understand both the Router functionality and the Bifrost transaction parsing to properly delve into this problem space.


## Links

- **Previous audits:**  
1. [Last](https://github.com/thorchain/Resources/blob/master/Audits/Halborn-StateMachine-Router-Bifrost-Audit-Sep2021.pdf)
2. [Second Last](https://github.com/thorchain/Resources/blob/master/Audits/THORChain-TrailOfBits-FullAudit-Aug2021.pdf)
- **Documentation:** https://docs.thorchain.org/
- **Website:** https://thorchain.org/
- **X/Twitter:** https://x.com/THORChain

---

# Scope



### Files in scope


This contest is focused on **THORChain Removal of Whitelisting on Router and Router V5**

- Only whitelisted contracts can call into and receive calls from TC Router.
[See here](https://gitlab.com/thorchain/thornode/-/blob/develop/bifrost/pkg/chainclients/ethereum/whitelist_smartcontract.go?ref_type=headshttps://gitlab.com/thorchain/thornode/-/blob/develop/bifrost/pkg/chainclients/ethereum/whitelist_smartcontract_aggregators.go?ref_type=heads)
- This is because the July 2021 hacks were all from attack contracts into the router, which faked deposits or tricked the bifrost into refunds.
[See here](https://rekt.news/thorchain-rekt/https://rekt.news/thorchain-rekt2/)
- TC wants to remove the whitelisting; but wants to make sure there is no attack paths possible on the router.
[See here](https://gitlab.com/thorchain/thornode/-/blob/develop/chain/ethereum/contracts/THORChain_Router.sol)
- Focus on how the Bifrost scans ETH events.
[See here](https://gitlab.com/thorchain/thornode/-/blob/develop/bifrost/pkg/chainclients/ethereum/ethereum_block_scanner.go#L818)


| Contract                                                                                                                                                                                     | SLOC | Purpose | Libraries used |
|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---- | ------- |:-------------- |
| [chain/ethereum/contracts/THORChain_Router.sol](https://github.com/code-423n4/2024-06-thorchain/blob/main/chain/ethereum/contracts/THORChain_Router.sol)                                     | 378  |         |                |
| [bifrost/pkg/chainclients/shared/evm/smartcontract_log_parser.go](https://github.com/code-423n4/2024-06-thorchain/blob/main/bifrost/pkg/chainclients/shared/evm/smartcontract_log_parser.go) | 315   |         |                |
| [bifrost/pkg/chainclients/ethereum/ethereum_block_scanner.go](https://github.com/code-423n4/2024-06-thorchain/blob/main/bifrost/pkg/chainclients/ethereum/ethereum_block_scanner.go)         | 824  |         |                |
| TOTAL                                                                                                                                                                                        | 1517 |         |                |

### Files out of scope  
All files not listed above are Out Of Scope.
This function of Router V5 is out of scope: `returnVaultAssets`

## Scoping Q &amp; A

### General questions

| Question                                | Answer                       |
| --------------------------------------- | ---------------------------- |
| ERC20 used by the protocol              |       Any (all possible ERC20s)             |
| Test coverage                           |  Ethereum:  Functions - 46.91% , Lines - 49.33% , Avalanche: Functions - 100%, Lines - 97.44%                   |
| ERC721 used  by the protocol            |            None              |
| ERC777 used by the protocol             |           None                |
| ERC1155 used by the protocol            |              None            |
| Chains the protocol will be deployed on | Ethereum, Avax, BSC |

### ERC20 token behaviors in scope

| Question                                                                                                                                                   | Answer |
| ---------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| [Missing return values](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#missing-return-values)                                                      |   Yes  |
| [Fee on transfer](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#fee-on-transfer)                                                                  |  Yes  |
| [Balance changes outside of transfers](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#balance-modifications-outside-of-transfers-rebasingairdrops) | Yes    |
| [Upgradeability](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#upgradable-tokens)                                                                 |   Yes  |
| [Flash minting](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#flash-mintable-tokens)                                                              | Yes    |
| [Pausability](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#pausable-tokens)                                                                      | Yes    |
| [Approval race protections](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#approval-race-protections)                                              | Yes    |
| [Revert on approval to zero address](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#revert-on-approval-to-zero-address)                            | Yes    |
| [Revert on zero value approvals](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#revert-on-zero-value-approvals)                                    | Yes    |
| [Revert on zero value transfers](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#revert-on-zero-value-transfers)                                    | Yes    |
| [Revert on transfer to the zero address](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#revert-on-transfer-to-the-zero-address)                    | Yes    |
| [Revert on large approvals and/or transfers](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#revert-on-large-approvals--transfers)                  | Yes    |
| [Doesn't revert on failure](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#no-revert-on-failure)                                                   |  Yes   |
| [Multiple token addresses](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#revert-on-zero-value-transfers)                                          | Yes    |
| [Low decimals ( < 6)](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#low-decimals)                                                                 |   Yes  |
| [High decimals ( > 18)](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#high-decimals)                                                              | Yes    |
| [Blocklists](https://github.com/d-xo/weird-erc20?tab=readme-ov-file#tokens-with-blocklists)                                                                | Yes    |

### External integrations (e.g., Uniswap) behavior in scope:


| Question                                                  | Answer |
| --------------------------------------------------------- | ------ |
| Enabling/disabling fees (e.g. Blur disables/enables fees) | No   |
| Pausability (e.g. Uniswap pool gets paused)               |  No   |
| Upgradeability (e.g. Uniswap gets upgraded)               |   No  |


### EIP compliance checklist
None


# Additional context

## Main invariants

- A transaction calling `depositWithExpiry` should be rejected if it is confirmed after the expiration parameter.
- `transferOut` should only update the allowance mapping if the ERC20 transfer was successful.
- `transferOut` should only allow ERC20s to be transferred from the Router if msg.sender has the appropriate allowance for the asset stored in the `_vaultAllowances` map.
- `transferAllowance` should only update the `_vaultAllowances` map if `msg.sender` already has the appropriate allowance for the asset.
- `deposit` should forward ETH directly to the provided vault address.
- `deposit` should keep ERC20s on the router contract and update the vault allowance.
- Only valid events emitted from the Router contract itself should result in the txInItem parameter being populated in the `GetTxInItem` function of the [smartcontract_log_parser](https://gitlab.com/thorchain/thornode/-/blob/develop/bifrost/pkg/chainclients/shared/evm/smartcontract_log_parser.go#L166)
- A tx with more logs than `max_contract_tx_logs` should be ignored by the `GetTxInItem` function of the `smartcontract_log_parser`


## Attack ideas (where to focus for bugs)
- The primary concern is a contract interacting with the Router and tricking the `smartcontract_log_parser`, and therefore the network, into thinking a DepositEvent has been emitted by the THORChain router when in fact it was emitted by a different contract. If this is possible, then an attacker could send in a fake swap or other transaction and extract value from THORChain's liquidity pools. 
- A malicious smart contract that interacts with the Router should also not be able to trick the `smartcontract_log_parser` that the DepositEvent has a different amount of ETH or ERC20s that were actually sent in by the transaction. 
- The contract stores a map of each vault's allowance for each ERC20 token stored. There should be no way for a malicious contract or attack to use a vault's allowance (stored in `_vaultAllowance`) to transfer out ERC20 tokens from the Router using `transferOut`, `transferOutAndCall`, `transferOutV5`, `transferOutAndCallV5`, `batchTransferOut`, or `batchTransferOutAndCallV5`
- There should be no way for an attacker to abscond a vault's allowance using the `transferAllowance` function or otherwise



## All trusted roles in the protocol

None


## Describe any novel or unique curve logic or mathematical models implemented in the contracts:

None


## Running tests

> [!NOTE]
> More detailed build, compile, and test instructions found in contract repo README: https://gitlab.com/thorchain/thornode/-/tree/develop/chain/ethereum


```bash
git clone https://github.com/code-423n4/2024-06-thorchain.git
git submodule update --init --recursive
cd ethereum
npx hardhat clean
npx hardhat compile
npx hardhat test
cd avalanche
npx hardhat clean
npx hardhat compile
npx hardhat test
```
To run code coverage
```bash
npx hardhat coverage
```
![Screenshot from 2024-06-05 20-38-56](https://github.com/code-423n4/2024-06-thorchain/assets/65364747/70bf0d43-d92e-409e-a2ec-2a6a05fa4cf3)
![Screenshot from 2024-06-05 20-28-34](https://github.com/code-423n4/2024-06-thorchain/assets/65364747/1041d87d-776d-47e0-a59d-d9e53e2a9db3)


## Miscellaneous
Employees of Thorchain and employees' family members are ineligible to participate in this audit.



