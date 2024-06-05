# EVM Implementation

THORChain Aggregator Example

{{#embed https://gitlab.com/thorchain/ethereum/eth-router/-/blob/master/contracts/THORChain_Aggregator.sol }}

Tokens must be on the [ETH Whitelist](https://gitlab.com/thorchain/thornode/-/blob/develop/docs/evm_whitelist_procedure.md#dex-token). The destination address should be a user control address, not a contract address.

## SwapIn

The aggregator contract needs a **swapIn** function similar to the one below. First, swap the token via an on-chain AMM, then call into THORChain and pass the correct memo to execute the next swap.

```go
function swapIn(
    address tcVault,
    address tcRouter,
    string calldata tcMemo,
    address token,
    uint amount,
    uint amountOutMin,
    uint256 deadline
    ) public nonReentrant {
    uint256 _safeAmount = safeTransferFrom(token, amount); // Transfer asset
    safeApprove(token, address(swapRouter), amount);
    address[] memory path = new address[](2);
    path[0] = token; path[1] = WETH;
    swapRouter.swapExactTokensForETH(_safeAmount, amountOutMin, path, address(this), deadline);
    _safeAmount = address(this).balance;
    iROUTER(tcRouter).depositWithExpiry{value:_safeAmount}(payable(tcVault), ETH, _safeAmount, tcMemo, deadline);
}
```

[Transaction Example](https://etherscan.io/tx/0x7905c41daaa214fbb3bad79ef63bb69aafcb15147f53cd9cf621d4049c2cea4d). Note the destination address is not a contract address.

## SwapOut

The THORChain router uses `transferOutAndCall()` to call the aggregator with a max GasLimit of 400k units.

It is a particular function that also handles a swap fail by sending the user the base asset directly (ie, breached AmountOutMin, or could not find the finaltoken). The user will need to do the swap manually.

The parameters for this function are passed to THORChain by the user's original memo.

```go
function transferOutAndCall(address payable target, address finalToken, address to, uint256 amountOutMin, string memory memo) public payable nonReentrant {
        uint256 _safeAmount = msg.value;
        (bool success, ) = target.call{value:_safeAmount}(abi.encodeWithSignature("swapOut(address,address,uint256)", finalToken, to, amountOutMin));
        if (!success) {
            payable(address(to)).transfer(_safeAmount); // If can't swap, just send the recipient the ETH
        }
        emit TransferOutAndCall(msg.sender, target, address(0), _safeAmount, finalToken, to, amountOutMin, memo);
    }
```

The **swapOut** function will only be passed three parameters from the THORChain Router and it must comply with the function signature (name, parameters). It can then call an on-chain AMM to execute the swap. It will only ever be given the base asset (eg ETH).

Here is an example to call UniV2 router:

```go
function swapOut(address token, address to, uint256 amountOutMin) public payable nonReentrant {
        address[] memory path = nelw address[](2);
        path[0] = WETH; path[1] = token;
        swapRouter.swapExactETHForTokens{value: msg.value}(amountOutMin, path, to, type(uint).max);
    }
```
