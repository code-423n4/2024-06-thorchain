# Eth Vault Contract

Vault Contract for Bifrost

- Incoming deposits register a `memo` and either
  - A. forward funds to chosen vault
  - B. give vault allowance to spend for specified asset
- Vaults can transfer allowance to spend to other vaults (asgard, yggdrasil)
- Vaults can transferOut up to their allowance to spend

This mirrors how vaults behave on other chains.

## THORChain Integration

### Bifrost - Observer

### ETH & ERC-20 Deposits

Users query THORChain for the correct asgard, then call deposit with a memo:

- ETH/ERC20: `await VAULT.depositWithExpiry(vault, asset, amount, memo, expiration)`

Bifrost should parse contract events to read the `memo` for asset transfers.

_Note 1: ETHER/Gas asset uses asset `0x0000000000000000000000000000000000000000`_
_Note 2: expiration is in seconds_
_Note 3: `deposit` method without expiry was deprecated in RouterV4.1 and removed in RouterV5_
_Note 4: V5 functions use structure to reduce local variables and prevent stack too deep errors_

### Bifrost - Signer

Asgard/Ygg vaults should compose the transactions:

- ETH: `await web3.eth.sendTransaction({ from:0xself, value:value, to:to, data:memo })`
- ERC20: `await VAULT.transferOut(to, asset, value, memo)`
- Aggregation `await VAULT.transferOutAndCall(target, fromAsset, fromAmount, toAsset, recipient, amountOutMin, memo, payload, originAddress)`

### Funding Ygg (only ERC20)

Ygg can be funded by transferring an allowance to spend:
`await VAULT.TransferAllowance(vault, asset, amount, memo)`

Ygg can return assets by batch transferring back:
`await VAULT.batchTransferAllowance(asgard[], assets[], amounts[], memos[])`

### Churning (only ERC20)

Asgard can transfer allowance to spend across 5 iterations:
`await VAULT.transferAllowance(newAsgard, asset, amount, memo)`

### Batch outbounds

RouterV5 re-introduces batch outbound transactions to optimize gas required for normal operations.

`batchTransferOut` supports outbounds for gas assets and tokens in the same batch. `TransferOut` logs are emitted normally for each outbound.

`batchTransferOutAndCall` works the same as `batchTransferOut` but with dex aggregation calldata and flow.

Attempts have been made to use a try/catch pattern to prevent 1 failing outbound to revert the whole batch. Try/catch is currently limited to external functions only. Using `_internal` and `external` functions combo to call as external however changes the transaction context and the `msg.sender` gets replaced from the original Asgard Vault to the router itself, breaking allowance logic.

## Contract Design

### Public Getters

Tracks vault allowances

```solidity
mapping(address => mapping(address => uint)) public vaultAllowance;
```

### Events

```solidity
event Deposit(address indexed to, address indexed asset, uint amount, string memo);
event TransferOut(address indexed vault, address indexed to, address asset, uint amount, string memo);
event TransferOutAndCall(address indexed vault, address target, uint amount, address finalAsset, address indexed to, uint256 amountOutMin, string memo, bytes payload, string originAddress);
event TransferAllowance(address indexed oldVault, address indexed newVault, address asset, uint amount, string memo);
event VaultTransfer(address indexed oldVault, address indexed newVault, Coin[] coins, string memo);
```

### Testing

```bash
npx hardhat clean
npx hardhat compile
npx hardhat test
```

## Deployment

Get ABI and BYTECODES from `/artifacts`
