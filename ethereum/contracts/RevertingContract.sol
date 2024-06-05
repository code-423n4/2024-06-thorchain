// SPDX-License-Identifier: MIT
// -------------------
// RevertingContract v1.0
// -------------------

pragma solidity 0.8.22;

contract RevertingContract {
  receive() external payable {
    revert();
  }
}
