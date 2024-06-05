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
- Technically any ERC20 can interact with the THORChain Router, but there is a whitelist implemented in Bifrost that limits which tokens a tx will be processed for. The whitelist is defined here: https://gitlab.com/thorchain/thornode/-/blob/develop/common/tokenlist/ethtokens/eth_mainnet_latest.json?ref_type=heads

If an un-whitelisted token interacts with the Router (for example through the depositWithExpiry function), the tx will effectively be dropped on the floor, as the whitelist is checked in the scp.assetResolver function. The asset resolver + whitelist functionality can be seen here: https://gitlab.com/thorchain/thornode/-/blob/develop/bifrost/pkg/chainclients/ethereum/ethereum_block_scanner.go#L762

- Nothing related to ERC20 RUNE (https://etherscan.io/address/0x3155BA85D5F96b2d030a4966AF206230e46849cb) is in scope. This asset is deprecated and cannot interact with the network anymore. 

✅ SCOUTS: Please format the response above 👆 so its not a wall of text and its readable.

# Overview

[ ⭐️ SPONSORS: add info here ]

## Links

- **Previous audits:**  https://github.com/thorchain/Resources/blob/master/Audits/Halborn-StateMachine-Router-Bifrost-Audit-Sep2021.pdf

https://github.com/thorchain/Resources/blob/master/Audits/THORChain-TrailOfBits-FullAudit-Aug2021.pdf
  - ✅ SCOUTS: If there are multiple report links, please format them in a list.
- **Documentation:** https://docs.thorchain.org/
- **Website:** https://thorchain.org/
- **X/Twitter:** https://x.com/THORChain

---

# Scope

[ ✅ SCOUTS: add scoping and technical details here ]

### Files in scope
- ✅ This should be completed using the `metrics.md` file
- ✅ Last row of the table should be Total: SLOC
- ✅ SCOUTS: Have the sponsor review and and confirm in text the details in the section titled "Scoping Q amp; A"

*For sponsors that don't use the scoping tool: list all files in scope in the table below (along with hyperlinks) -- and feel free to add notes to emphasize areas of focus.*

This contest is focused on **THORChain Removal of Whitelisting on Router**

- Only whitelisted contracts can call into and receive calls from TC Router
https://gitlab.com/thorchain/thornode/-/blob/develop/bifrost/pkg/chainclients/ethereum/whitelist_smartcontract.go?ref_type=headshttps://gitlab.com/thorchain/thornode/-/blob/develop/bifrost/pkg/chainclients/ethereum/whitelist_smartcontract_aggregators.go?ref_type=heads
- This is because the July 2021 hacks were all from attack contracts into the router, which faked deposits or tricked the bifrost into refunds
https://rekt.news/thorchain-rekt/https://rekt.news/thorchain-rekt2/
- TC wants to remove the whitelisting; but wants to make sure there is no attack paths possible on the router
https://gitlab.com/thorchain/thornode/-/blob/develop/chain/ethereum/contracts/THORChain_Router.sol
- Focus on how the Bifrost scans ETH events
https://gitlab.com/thorchain/thornode/-/blob/develop/bifrost/pkg/chainclients/ethereum/ethereum_block_scanner.go#L818

| Contract | SLOC | Purpose | Libraries used |  
| ----------- | ----------- | ----------- | ----------- |
| [contracts/folder/sample.sol](https://github.com/code-423n4/repo-name/blob/contracts/folder/sample.sol) | 123 | This contract does XYZ | [`@openzeppelin/*`](https://openzeppelin.com/contracts/) |

### Files out of scope
✅ SCOUTS: List files/directories out of scope

## Scoping Q &amp; A

### General questions
### Are there any ERC20's in scope?: Yes

✅ SCOUTS: If the answer above 👆 is "Yes", please add the tokens below 👇 to the table. Otherwise, update the column with "None".

Any (all possible ERC20s)


### Are there any ERC777's in scope?: No

✅ SCOUTS: If the answer above 👆 is "Yes", please add the tokens below 👇 to the table. Otherwise, update the column with "None".



### Are there any ERC721's in scope?: No

✅ SCOUTS: If the answer above 👆 is "Yes", please add the tokens below 👇 to the table. Otherwise, update the column with "None".



### Are there any ERC1155's in scope?: No

✅ SCOUTS: If the answer above 👆 is "Yes", please add the tokens below 👇 to the table. Otherwise, update the column with "None".



✅ SCOUTS: Once done populating the table below, please remove all the Q/A data above.

| Question                                | Answer                       |
| --------------------------------------- | ---------------------------- |
| ERC20 used by the protocol              |       🖊️             |
| Test coverage                           | ✅ SCOUTS: Please populate this after running the test coverage command                          |
| ERC721 used  by the protocol            |            🖊️              |
| ERC777 used by the protocol             |           🖊️                |
| ERC1155 used by the protocol            |              🖊️            |
| Chains the protocol will be deployed on | Ethereum,Avax,BSC |

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
N/A

✅ SCOUTS: Please format the response above 👆 using the template below👇

| Question                                | Answer                       |
| --------------------------------------- | ---------------------------- |
| src/Token.sol                           | ERC20, ERC721                |
| src/NFT.sol                             | ERC721                       |


# Additional context

## Main invariants

- a transaction calling depositWithExpiry should be rejected if it is confirmed after the expiration parameter
- transferOut should only update the allowance mapping if the ERC20 transfer was successful 
- transferOut should only allow ERC20s to be transferred from the Router if msg.sender has the appropriate allowance for the asset stored in the _vaultAllowances map
- transferAllowance should only update the _vaultAllowances map if msg.sender already has the appropriate allowance for the asset
- deposit should forward ETH directly to the provided vault address
- deposit should keep ERC20s on the router contract and update the vault allowance
- only valid events emitted from the Router contract itself should result in the txInItem parameter being populated in the GetTxInItem function of the smartcontract_log_parser https://gitlab.com/thorchain/thornode/-/blob/develop/bifrost/pkg/chainclients/shared/evm/smartcontract_log_parser.go#L166
- a tx with more logs than `max_contract_tx_logs` should be ignored by the GetTxInItem function of the smartcontract_log_parser

✅ SCOUTS: Please format the response above 👆 so its not a wall of text and its readable.

## Attack ideas (where to focus for bugs)
- The primary concern is a contract interacting with the Router and tricking the smartcontract_log_parser, and therefore the network, into thinking a DepositEvent has been emitted by the THORChain router when in fact it was emitted by a different contract. If this is possible, then an attacker could send in a fake swap or other transaction and extract value from THORChain's liquidity pools. 
- A malicious smart contract that interacts with the Router should also not be able to trick the smartcontract_log_parser that the DepositEvent has a different amount of ETH or ERC20s that were actually sent in by the transaction. 
- The contract stores a map of each vault's allowance for each ERC20 token stored. There should be no way for a malicious contract or attack to use a vault's allowance (stored in _vaultAllowance) to transfer out ERC20 tokens from the Router using transferOut, transferOutAndCall, transferOutV5, transferOutAndCallV5, batchTransferOut, or batchTransferOutAndCallV5
- There should be no way for an attacker to abscond a vault's allowance using the transferAllowance function or otherwise

✅ SCOUTS: Please format the response above 👆 so its not a wall of text and its readable.

## All trusted roles in the protocol

N/A

✅ SCOUTS: Please format the response above 👆 using the template below👇

| Role                                | Description                       |
| --------------------------------------- | ---------------------------- |
| Owner                          | Has superpowers                |
| Administrator                             | Can change fees                       |

## Describe any novel or unique curve logic or mathematical models implemented in the contracts:

N/A

✅ SCOUTS: Please format the response above 👆 so its not a wall of text and its readable.

## Running tests

Build, compile, and test instructions found in contract repo README: https://gitlab.com/thorchain/thornode/-/tree/develop/chain/ethereum

✅ SCOUTS: Please format the response above 👆 using the template below👇

```bash
git clone https://github.com/code-423n4/2023-08-arbitrum
git submodule update --init --recursive
cd governance
foundryup
make install
make build
make sc-election-test
```
To run code coverage
```bash
make coverage
```
To run gas benchmarks
```bash
make gas
```

✅ SCOUTS: Add a screenshot of your terminal showing the gas report
✅ SCOUTS: Add a screenshot of your terminal showing the test coverage

## Miscellaneous
Employees of Thorchain and employees' family members are ineligible to participate in this audit.



