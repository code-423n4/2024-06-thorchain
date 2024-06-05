# Math

## ​Math Library

The following libraries implement the math below.

{{#embed <https://gitlab.com/thorchain/asgardex-common/asgardex-util/-/tree/master/src/calc> }}

{{#embed <https://github.com/xchainjs/xchainjs-lib/blob/master/packages/xchain-thorchain-query/src/utils/swap.ts> }}

### Example Data

All the examples below use the following snapshotted BTC and BUSD Pool data

​[https://thornode.ninerealms.com/thorchain/pool/BTC.BTC](https://thornode.ninerealms.com/thorchain/pool/BTC.BTC)​ ​[https://thornode.ninerealms.com/thorchain/pool/BNB.BUSD-BD1](https://thornode.ninerealms.com/thorchain/pool/BNB.BUSD-BD1)​

```json
{
 "LP_units": "117582615428135",
 "asset": "BNB.BUSD-BD1",
 "balance_asset": "952382623537567",
 "balance_rune": "508868258770825",
 "pending_inbound_asset": "310872270739",
 "pending_inbound_rune": "1701596418307",
 "pool_units": "134664599295503",
 "status": "Available",
 "synth_supply": "241616351972821",
 "synth_units": "17081983867368"
},
{
 "LP_units": "476785169622350",
 "asset": "BTC.BTC",
 "balance_asset": "81439552768",
 "balance_rune": "863897777396922",
 "pending_inbound_asset": "386699833",
 "pending_inbound_rune": "216117023429",
 "pool_units": "492710913491074",
 "status": "Available",
 "synth_supply": "5264691415",
 "synth_units": "15925743868724"
}
```

## Prices

$$
price = \frac{quoteBalance}{baseBalance}=\frac{USD}{RUNE} = $RUNE
$$

Prices of all assets on THORChain are in ratios of each other, based on the depths of the pools. The quote asset is the "pricing" asset, and the base asset is the asset to be quoted. Ie, for the $ value of RUNE, the quote asset is USD and the base asset is RUNE.

### Example

Let's take the BTC and BUSD Pool data

[https://thornode.ninerealms.com/thorchain/pool/BTC.BTC](https://thornode.ninerealms.com/thorchain/pool/BTC.BTC)

[https://thornode.ninerealms.com/thorchain/pool/BNB.BUSD-BD1](https://thornode.ninerealms.com/thorchain/pool/BNB.BUSD-BD1)

The $BTC Price of RUNE is `BTC/RUNE = 81439552768/863897777396922 = 0.000094 BTC`

The $BUSD price of BTC is `(BUSD/RUNE) * (RUNE/BTC) = (952382623537567/508868258770825) * (863897777396922/81439552768) = 19,854 BUSD`

```javascript
export const getValueOfAssetInRune = (inputAsset: BaseAmount, pool: PoolData): BaseAmount => {
  // formula: ((a * R) / A) => R per A (Runeper$)
  const t = inputAsset.amount()
  const R = pool.runeBalance.amount()
  const A = pool.assetBalance.amount()
  const result = t.times(R).div(A)
  return baseAmount(result)
}

export const getValueOfRuneInAsset = (inputRune: BaseAmount, pool: PoolData): BaseAmount => {
  // formula: ((r * A) / R) => A per R ($perRune)
  const r = inputRune.amount()
  const R = pool.runeBalance.amount()
  const A = pool.assetBalance.amount()
  const result = r.times(A).div(R)
  return baseAmount(result)
}
export const getValueOfAsset1InAsset2 = (inputAsset: BaseAmount, pool1: PoolData, pool2: PoolData): BaseAmount => {
  // formula: (A2 / R) * (R / A1) => A2/A1 => A2 per A1 ($ per Asset)
  const oneAsset = assetToBase(assetAmount(1))
  // Note: All calculation needs to be done in `AssetAmount` (not `BaseAmount`)
  const A2perR = baseToAsset(getValueOfRuneInAsset(oneAsset, pool2))
  const RperA1 = baseToAsset(getValueOfAssetInRune(inputAsset, pool1))
  const result = A2perR.amount().times(RperA1.amount())
  // transform result back from `AssetAmount` into `BaseAmount`
  return assetToBase(assetAmount(result))
}
```

## Slippage

$$
slippage = \frac{inputAmount}{(inputAmount + inputBalance)}
$$

Slippage is simply the transaction divided by its corresponding depth.

Ie, swapping 10 BTC to RUNE = `1000000000 / (1000000000 + 81439552768) = 0.012 = 1.2%`

Since THORChain has all pools in RUNE, a cross-asset (double) swap would involve two swaps in two pools, thus the slip needs to be doubled.

Here's a reference implementation of calculating slip for a double swap:

```javascript
// Calculate swap output with slippage
function calcSwapOutput(inputAmount, pool, toRune) {
  // formula: (inputAmount * inputBalance * outputBalance) / (inputAmount + inputBalance) ^ 2
  const inputBalance = toRune ? pool.assetBalance : pool.runeBalance; // input is asset if toRune
  const outputBalance = toRune ? pool.runeBalance : pool.assetBalance; // output is rune if toRune
  const numerator = inputAmount * inputBalance * outputBalance;
  const denominator = Math.pow(inputAmount + inputBalance, 2);
  const result = numerator / denominator;
  return result;
}

// Calculate swap slippage
function calcSwapSlip(inputAmount, pool, toRune) {
  // formula: (inputAmount) / (inputAmount + inputBalance)
  const inputBalance = toRune ? pool.assetBalance : pool.runeBalance; // input is asset if toRune
  const result = inputAmount / (inputAmount + inputBalance);
  return result;
}

// Calculate swap slippage for double swap
function calcDoubleSwapSlip(inputAmount, pool1, pool2) {
  // formula: calcSwapSlip1(input1) + calcSwapSlip2(calcSwapOutput1 => input2)
  const swapSlip1 = calcSwapSlip(inputAmount, pool1, true);
  const r = calcSwapOutput(inputAmount, pool1, true);
  const swapSlip2 = calcSwapSlip(r, pool2, false);
  const result = swapSlip1 + swapSlip2;
  return result;
}
```

_Source:_ [_https://gitlab.com/thorchain/asgardex-common/asgardex-util/-/blob/master/src/calc/swap.ts_](https://gitlab.com/thorchain/asgardex-common/asgardex-util/-/blob/master/src/calc/swap.ts)

## Swap Output

$$
output = \frac{(inputAmount * outputBalance * inputBalance)}{(inputAmount + inputBalance)^{2}}
$$

​The output in a swap is the CLP formula.

Ie, output after swapping 10 BTC: `(1000000000 * 81439552768 * 863897777396922)/ (1000000000 + 81439552768)^2 = 10352052898302 = 103520 RUNE`

```javascript
export const getSwapOutput = (inputAmount: BaseAmount, pool: PoolData, toRune: boolean): BaseAmount => {
  // formula: (x * X * Y) / (x + X) ^ 2
  const x = inputAmount.amount()
  const X = toRune ? pool.assetBalance.amount() : pool.runeBalance.amount() // input is asset if toRune
  const Y = toRune ? pool.runeBalance.amount() : pool.assetBalance.amount() // output is rune if toRune
  const numerator = x.times(X).times(Y)
  const denominator = x.plus(X).pow(2)
  const result = numerator.div(denominator)
  return baseAmount(result)
}
export const getDoubleSwapOutput = (inputAmount: BaseAmount, pool1: PoolData, pool2: PoolData): BaseAmount => {
  // formula: getSwapOutput(pool1) => getSwapOutput(pool2)
  const r = getSwapOutput(inputAmount, pool1, true)
  const output = getSwapOutput(r, pool2, false)
  return output
}

```

## Swap Input

$$
\text{input} = \frac{X(-\sqrt{Y (Y - 4y) - 2y + Y})}{2y}
$$

`X = inputBalance`

`Y = outputBalance`, `y = outputAmount`

The swap formula can be reversed to specify what needs to be deposited to get a certain output.

```javascript
export const getSwapInput = (toRune: boolean, pool: PoolData, outputAmount: BaseAmount): BaseAmount => {
  // formula: (((X*Y)/y - 2*X) - sqrt(((X*Y)/y - 2*X)^2 - 4*X^2))/2
  // (part1 - sqrt(part1 - part2))/2
  const X = toRune ? pool.assetBalance.amount() : pool.runeBalance.amount() // input is asset if toRune
  const Y = toRune ? pool.runeBalance.amount() : pool.assetBalance.amount() // output is rune if toRune
  const y = outputAmount.amount()
  const part1 = X.times(Y).div(y).minus(X.times(2))
  const part2 = X.pow(2).times(4)
  const result = part1.minus(part1.pow(2).minus(part2).sqrt()).div(2)
  return baseAmount(result)
}

```

### LP Units Add

$$
\text{units} = \frac{P(R_a + rA)}{2RA}
$$

- \(P\): Existing Pool Units
- \(R\): runeBalance, \(A\): assetBalance
- \(r\): runeAdded, \(a\): assetAdded

The units to give an LP depend on the existing units, as well as the assets they are adding, and the depths of the pool they are adding to.

## LP Units Withdrawn

$$
\text{output} = \frac{\text{basisPoints} \cdot L \cdot X}{10000 \cdot P}
$$

- \(L\): Liquidity units owned
- \(P\): Pool Units
- \(X\): depth of side

THORChain allows LPs to redeem a Basis Points amount of their position (out of 10000). To find out how much the user will get, multiply this by each side.

```javascript
export const getPoolShare = (unitData: UnitData, pool: PoolData): StakeData => {
  // formula: (rune * part) / total; (asset * part) / total
  const units = unitData.stakeUnits.amount()
  const total = unitData.totalUnits.amount()
  const R = pool.runeBalance.amount()
  const T = pool.assetBalance.amount()
  const asset = T.times(units).div(total)
  const rune = R.times(units).div(total)
  const stakeData = {
    asset: baseAmount(asset),
    rune: baseAmount(rune)
  }
  return stakeData
}

```
