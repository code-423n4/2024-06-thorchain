# Coding Guide

A coding overview to xchainjs.

## **General**

The foundation of xchainjs is defined in the [xchain-util](https://github.com/xchainjs/xchainjs-lib/tree/master/packages/xchain-util) package

- `Address`: a crypto address as a string.
- `BaseAmount`: a bigNumber in 1e8 format. E.g. 1 BTC = 100,000,000 in BaseAmount
- `AssetAmount`: a BaseAmount\*10^8. E.g. 1 BTC = 1 in Asset Amount.
- `Asset`: Asset details {Chain, symbol, ticker, isSynth}

```admonish info
All `Assets` must conform to the [Asset Notation](../../concepts/memos.md#asset-notation)

`assetFromString()` is used to quickly create assets and will assign chain and synth.
```

- `CryptoAmount:` is a class that has:

```javascript
 baseAmount: BaseAmount
 readonly asset: Asset
```

All crypto should use the `CryptoAmount` object with the understanding they are in BaseAmount format. An example to switch between them:

```typescript
// Define 1 BTC as CryptoAmount
oneBtc = new CryptoAmount(
  assetToBase(assetAmount(1)),
  assetFromStringEx(`BTC.BTC`),
);
// Print 1 BTC in BaseAmount
console.log(oneBtc.amount().toNumber().toFixed()); // 100000000
// Print 1 BTC in Asset Amount
console.log(oneBtc.AssetAmount().amount().toNumber().toFixed()); // 1
```

## Query

Major data types for the thorchain-query package.

- [Package description](query-package.md)
- [Github source code](https://github.com/xchainjs/xchainjs-lib/tree/master/packages/xchain-thorchain-query)
- [Code examples](query-package.md#code-examples-in-replit)
- [Install procedures](overview.md#install-procedures)

### **EstimateSwapParams**

### **SwapInput**

The input Type for `estimateSwap`. This is designed to be created by interfaces and passed into EstimateSwap. Also see [Swap Memo](../../concepts/memos.md#swap) for more information.

| Variable              | Data Type    | Comments                        |
| --------------------- | ------------ | ------------------------------- |
| `input`               | CryptoAmount | Inbound asset and amount        |
| `destinationAsset`    | Asset        | Outbound asset                  |
| `destinationAddress`  | String       | Outbound asset address          |
| `slipLimit`           | BigNumber    | Optional: Used to set LIM       |
| `affiliateFeePercent` | number       | Optional: 0-0.1 allowed         |
| `affiliateAddress`    | Address      | Optional: THOR address          |
| `interfaceID`         | string       | Optional: Assigned interface ID |

### **SwapEstimate**

The internal return type is used within `estimateSwap` after the calculation is done.

| Variable          | Data Type    | Comments                            |
| ----------------- | ------------ | ----------------------------------- |
| `totalFees`       | TotalFees    | All fees for swap                   |
| `slipPercentage`  | BigNumber    | Actual slip of the swap             |
| `netOutput`       | CryptoAmount | Input - totalFees                   |
| `waitTimeSeconds` | number       | Estimated time for the swap         |
| `canSwap`         | boolean      | False if there is an issue          |
| `errors`          | string array | Contains info if `canSwap` is false |

### **TxDetails**

Return type of `estimateSwap`. This is designed to be used by interfaces to give them all the information they need to display to the user.

| Variable     | Data Type    | Comments                                                    |
| ------------ | ------------ | ----------------------------------------------------------- |
| `txEstimate` | SwapEstimate | Swap details                                                |
| `memo`       | string       | Constructed memo THORChain will understand                  |
| `expiry`     | DateTime     | When the `SwapEstimate` information will no longer be valid |
| `toAddress`  | string       | Current Asgard Vault address from `inbound_address`         |

```admonish danger
Do not use `toAddress` after `expiry` as the Asgard vault rotates
```

## AMM

Major data types for the thorchain-query package.

- [Package description](amm-package.md)
- [Github source code](https://github.com/xchainjs/xchainjs-lib/tree/master/packages/xchain-thorchain-amm)
- [Code examples](amm-package.md#code-examples-in-replit)
- [Install procedures](overview.md#install-procedures)

### **ExecuteSwap**

Input Type for doSwap where a swap will be actually conducted. Based on EstimateSwapParams.

### **TxSubmitted**

| Variable        |        |                             |
| --------------- | ------ | --------------------------- |
| hash            | string | inbound Tx Hash             |
| url             | string | Block exploer url           |
| waitTimeSeconds | number | Estimated time for the swap |
