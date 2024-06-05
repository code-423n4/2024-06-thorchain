# Memo Length Reduction

## Reducing Memo Size

Given the complexity of memos, they can become very long, beyond the limits of chains like Bitcoin. Various methods have been developed to significantly shorten memo length. \
\
**Example**:

1. `SWAP:ETH.ETH:0xe6a30f4f3bad978910e2cbb4d97581f5b5a0ade0:1612345678:thor1el4ufmhll3yw7zxzszvfakrk66j7fx0tvcslym:10`

   can be reduced to `=:e:dx:161e6:t:10`

2. `SWAP:ETH.USDT-0xdac17f958d2ee523a2206206994597c13d831ec7:0xe6a30f4f3bad978910e2cbb4d97581f5b5a0ade0:10012345678:thor1el4ufmhll3yw7zxzszvfakrk66j7fx0tvcslym:10`

   can be reduced to `=:ETH.USDT:dx:100e7:t:10`

The examples below use the following features to reduce memo length:

1. [Shortened Asset Names](memo-length-reduction.md#shortened-asset-names)
2. [THORNames](memo-length-reduction.md#mechanism-for-transaction-intent-1)
3. [Shortened Function](memo-length-reduction.md#mechanism-for-transaction-intent-2)
4. [Asset Abbreviations](memo-length-reduction.md#asset-abbreviations)
5. [Scientific Notation](memo-length-reduction.md#scientific-notation)

### **Shortened Asset Names**

Native asset names can be shortened to reduce the length of the memo. The exact list is [here](https://gitlab.com/thorchain/thornode/-/blob/develop/common/asset.go#L231).

| Shorten Asset | Asset Notation |
| ------------- | -------------- |
| r             | THOR.RUNE      |
| a             | AVAX.AVAX      |
| b             | BTC.BTC        |
| c             | BCH.BCH        |
| e             | ETH.ETH        |
| g             | GAIA.ATOM      |
| n             | BNB.BNB        |
| s             | BSC.BNB        |
| d             | DOGE.DOGE      |
| e             | ETH.ETH        |
| l             | LTC.LTC        |

**Example Swaps**:

- `=:ETH.ETH:0x388C818CA8B9251b393131C08a736A67ccB19297` is reduced to `=:e:0x388C818CA8B9251b393131C08a736A67ccB19297,` Swap for Ether.
- `=:r:thor1el4ufmhll3yw7zxzszvfakrk66j7fx0tvcslym` - Swap to RUNE.

### Asset Abbreviations

Assets can be abbreviated using fuzzy logic. The following will all be matched appropriately. If there are conflicts, then the deepest pool is matched to prevent attacks.

| Notation                                            |
| --------------------------------------------------- |
| ETH.USDT                                            |
| ETH.USDT-ec7                                        |
| ETH.USDT-6994597c13d831ec7                          |
| ETH.USDT-0xdac17f958d2ee523a2206206994597c13d831ec7 |

### THORNames

THORNames allows a custom name to be assigned to an address, like an alias, so the address does not need to be specified.

Example:

- `thor1nt2d4kmj0xd6xxm3m82tac3d20y05dm0vv7ur3` can be specified as `tr`.

See the [THORName Creation Guide](../affiliate-guide/thorname-guide.md) to create your own. This is used greatly to specify the affiliate address.

### Shortened Functions

Memos contain functions such as Swap or Add, which describe the user's intent and are sent along with specific parameters. Functions can be reduced in the following way:

| Function      | Abbreviated | Recommended |
| ------------- | ----------- | ----------- |
| Swap          | s           | =           |
| Add / Deposit |             | +           |
| Withdraw      | wd          | -           |
| Loan Open     | Loan+       | $+          |
| Loan Close    | Loan-       | $-          |
| THORName      | name, n     | \~          |
| Limit Order   | Limitto     | lo          |

**Example**:

`SWAP:e:0x388C818CA8B9251b393131C08a736A67ccB19297` is reduced to `=:e:0x388C818CA8B9251b393131C08a736A67ccB19297`

### Scientific Notation

In THORChain memos and the state machine, asset amounts are expressed as Base in 1e8 format requiring many digits to express an amount. E.g. 0.01 BTC is expressed as `1000000` and 5 Ether is expressed as `500000000`.

To help save space in memos, scientific notation can shorten memos by specifying both significant digits and the amount of trailing zeros. Note that using scientific notation in memos always leads to a loss of precision, so ensure enough significant digits are used to express the amount properly. For example, using `161e6` to represent `1612345678` results in a loss of precision.

**Examples:**

- In memo: `1e8` -> THORChain reads: `100000000`
- In memo: `51e7` -> THORChain reads: `510000000`

**Full Memo Example:**

- `SWAP:ETH.ETH:0x388C818CA8B9251b393131C08a736A67ccB19297:1612341234:thor19emplkuphjk2y9gkkv06m8vcstc0ufn4pevv5u:10`

  is reduced to `=:e:0x388C818CA8B9251b393131C08a736A67ccB19297:161e6:t:10`
