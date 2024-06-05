# Eth Vault Contract

Vault Contract for Bifrost.

- Incoming deposits register a `memo` and either 1) forward funds to chosen vault 2) give vault allowance to spend for specified asset
- Vaults can transfer allowance to spend to other vaults (asgard, yggdrasil)
- Vaults can transferOut up to their allowance to spend

This mirrors how vaults behave on other chains.

## THORChain Integration

### Bifrost - Observer

**ETH & ERC-20 Deposits**
Users query THORChain for the correct asgard, then call deposit with a memo:

- ETH: `await VAULT.depositETH(asgard, amount, memo)`
- ERC20: `await VAULT.deposit(asgard, asset, amount, memo)`

Bifrost should parse contract events to read the `memo` for asset transfers.

_Note: ETHER is `0x0000000000000000000000000000000000000000`_

### Bifrost - Signer

Asgard/Ygg vaults should compose the transactions:

- ETH: `await web3.eth.sendTransaction({ from:0xself, value:value, to:to, data:memo })`
- ERC20: `await VAULT.transferOut(to, asset, value, memo)`

### Funding Ygg (only ERC20)

Ygg can be funded by transferring an allowance to spend:
`await VAULT.transferAllance(vault, asset, amount, memo)`

Ygg can return assets by batch transferring back:
`await VAULT.batchTransferAllance(asgard[], assets[], amounts[], memos[])`

### Churning (only ERC20)

Asgard can transfer allowance to spend across 5 iterations:
`await VAULT.transferAllance(newAsgard, asset, amount, memo)`

## Contract Design

### Public Getters

Tracks vault allowances

```solidity
mapping(address => mapping(address => uint)) public vaultAllowance;
```

### Events

```solidity
event Deposit(address asset, uint value, string memo);
event TransferOut(address to, address asset, uint value, string memo);
event TransferAllowance(address vault, address asset, uint value, string memo);
```

### Testing

```bash
npx hardhat clean
npx hardhat compile
npx hardhat test
```

## Deployment

Get ABI and BYTECODES from `/artifacts`
