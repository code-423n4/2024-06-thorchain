// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract CallerContract {
    // call depositWithExpiry function
    function callDeposit(address routerAddr, address payable vault, address asset, uint amount, string memory memo) external payable {
        THORChainRouter router = THORChainRouter(routerAddr);
        router.depositWithExpiry{value: msg.value}(vault, asset, amount, memo, block.timestamp + 3600);
    }

    // call depositWithExpiry function after emitting logs
    function callDepositWithLogs(address routerAddr, address payable vault, address asset, uint amount, string memory memo) external payable {
        THORChainRouter router = THORChainRouter(routerAddr);
        for (uint i = 0; i < 10; i++) {
            emit BasicLog(i);
        }
        router.depositWithExpiry{value: msg.value}(vault, asset, amount, memo, block.timestamp + 3600);
    }

    // BasicLog
    event BasicLog(uint index);
}

// Interface for the other contract
interface THORChainRouter {
    function depositWithExpiry(address payable vault, address asset, uint amount, string memory memo, uint expiration) external payable;
}