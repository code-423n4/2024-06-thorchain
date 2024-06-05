# Memos

## Swap Memo (from [here](https://gitlab.com/thorchain/thornode/-/merge_requests/2218))

In order to support SwapOut DEX Aggregation feature, a few more fields added into the swap memo.

**`SWAP:ASSET:DESTADDR:LIM:AFFILIATE:FEE:DEXAggregatorAddr:FinalTokenAddr:MinAmountOut|`**

| Parameter            | Nodes                                                     | Conditions                                                                                                                                                                                                                                                                                                                                                       |
| -------------------- | --------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `:DEXAggregatorAddr` | The whitelisted aggregator contract.                      | Can use the last x characters of the address to fuzz match it.                                                                                                                                                                                                                                                                                                   |
| `:FinalTokenAddr`    | The final token (must be on 1INCH Whitelist)              | Can be shortened                                                                                                                                                                                                                                                                                                                                                 |
| `:minAmountOut`      | The parameter to pass into AmountOutMin in AMM contracts. | Handled by the aggregator, so:<br>1. Can be 0 (no protection). <br>2. Can be in any decimals<br>3. Can be in % or BasisPoints, then converted to a price at the time of swap by the aggregator contract.<br><br> Thornode accepts integers and scientific notation. Both `100000000` and `1e8` would forward `uint256 100000000` to the aggregator contract.</p> |

```admonish success
If you include a vertical pipe (|) at the end of the memo, any data following it will be sent as an outbound memo to the specified outbound address. This feature enables developers to send generic data to contracts cross-chain.
```

### Additional ObserveTxIn field

In order to support SwapOut Dex Aggregation feature safely , a few more fields have been added into tx out item

```json
[
  {
    "chain": "ETH",
    "to_address": "0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473",
    "vault_pub_key": "tthorpub1addwnpepq02wq8hwgmwge6p9yzwscyfp0kjv823kres7l7tcv89nn2zfu3jguu5s4qa",
    "coin": {
      "asset": "ETH.ETH",
      "amount": "19053206"
    },
    "memo": "OUT:EA7D80B3EB709319A6577AF6CF4DEFF67975D4F5A93CD8817E7FF04A048D1C5C",
    "max_gas": [
      {
        "asset": "ETH.ETH",
        "amount": "240000",
        "decimals": 8
      }
    ],
    "gas_rate": 3,
    "in_hash": "EA7D80B3EB709319A6577AF6CF4DEFF67975D4F5A93CD8817E7FF04A048D1C5C",
    "aggregator": "0x69800327b38A4CeF30367Dec3f64c2f2386f3848",    <-------------------- NEW
    "aggregator_target_asset": "0x0a44986b70527154e9F4290eC14e5f0D1C861822", <-------------------- NEW
    "aggregator_target_limit": "1000" <-------------------- NEW , but optional
  }
]
```

Also the same fields have been added to `ObservedTx` so THORNode can verify that bifrost did send out the transaction per instruction, use the aggregator per instructed , and pass target asset and target limit to the aggregator correctly

### How to swap out with dex aggregator?

If i want to swap RUNE to random ERC20 asset that is not list on THORChain , but is list on SushiSwap for example

```text
thornode tx thorchain deposit 200000000000 RUNE '=:ETH.ETH:0x3fd2d4ce97b082d4bce3f9fee2a3d60668d2f473::::2386f3848:0x0a44986b70527154e9F4290eC14e5f0D1C861822' --chain-id thorchain --node tcp://$THORNODE_IP:26657 --from {from user} --keyring-backend=file --yes --gas 20000000
```

**Note:**

1. Swap asset is `ETH.ETH`
2. `2386f3848` is the last nine characters of the aggregator contract address
3. `0x0a44986b70527154e9F4290eC14e5f0D1C861822` is the final asset address
4. Keep in mind SwapOut is best effort, when aggregator contract failed to perform the requested swap , then user will get ETH.ETH instead of the final asset it request
