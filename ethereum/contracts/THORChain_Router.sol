// SPDX-License-Identifier: MIT
// -------------------
// Router Version: 5.0
// -------------------
pragma solidity 0.8.22;

// ERC20 Interface
interface iERC20 {
  function balanceOf(address) external view returns (uint256);

  function burn(uint) external;
}

// ROUTER Interface
interface iROUTER {
  function depositWithExpiry(
    address,
    address,
    uint,
    string calldata,
    uint
  ) external;
}

// THORChain_Router is managed by THORChain Vaults
contract THORChain_Router {
  struct Coin {
    address asset;
    uint amount;
  }

  // Used in V5 function implementation
  struct TransferOutData {
    address payable to;
    address asset;
    uint amount;
    string memo;
  }

  // Used in V5 function implementation
  struct TransferOutAndCallData {
    address payable target;
    address fromAsset;
    uint256 fromAmount;
    address toAsset;
    address recipient;
    uint256 amountOutMin;
    string memo;
    bytes payload;
    string originAddress;
  }

  // Vault allowance for each asset
  mapping(address => mapping(address => uint)) private _vaultAllowance;

  uint256 private constant _NOT_ENTERED = 1;
  uint256 private constant _ENTERED = 2;
  uint256 private _status;

  // Emitted for all deposits, the memo distinguishes for swap, add, remove, donate etc
  event Deposit(
    address indexed to,
    address indexed asset,
    uint amount,
    string memo
  );

  // Emitted for all outgoing transfers, the vault dictates who sent it, memo used to track.
  event TransferOut(
    address indexed vault,
    address indexed to,
    address asset,
    uint amount,
    string memo
  );

  // Used in original V4 methods
  // Emitted for all outgoing transferAndCalls, the vault dictates who sent it, memo used to track.
  event TransferOutAndCall(
    address indexed vault,
    address target,
    uint amount,
    address finalAsset,
    address to,
    uint256 amountOutMin,
    string memo
  );

  // Emitted for all outgoing transferAndCalls, the vault dictates who sent it, memo used to track.
  event TransferOutAndCallV5(
    address indexed vault,
    address target,
    uint amount,
    address finalAsset,
    address indexed to,
    uint256 amountOutMin,
    string memo,
    bytes payload,
    string originAddress
  );

  // Changes the spend allowance between vaults
  event TransferAllowance(
    address indexed oldVault,
    address indexed newVault,
    address asset,
    uint amount,
    string memo
  );

  // Specifically used to batch send the entire vault assets
  event VaultTransfer(
    address indexed oldVault,
    address indexed newVault,
    Coin[] coins,
    string memo
  );

  modifier nonReentrant() {
    require(_status != _ENTERED, "ReentrancyGuard: reentrant call");
    _status = _ENTERED;
    _;
    _status = _NOT_ENTERED;
  }

  constructor() {
    _status = _NOT_ENTERED;
  }

  // Deposit with Expiry (preferred)
  function depositWithExpiry(
    address payable vault,
    address asset,
    uint amount,
    string memory memo,
    uint expiration
  ) external payable {
    require(block.timestamp < expiration, "THORChain_Router: expired");
    _deposit(vault, asset, amount, memo);
  }

  // Deposit an asset with a memo. ETH is forwarded, ERC-20 stays in ROUTER
  function _deposit(
    address payable vault,
    address asset,
    uint amount,
    string memory memo
  ) private nonReentrant {
    uint safeAmount;
    if (asset == address(0)) {
      safeAmount = msg.value;
      bool success = vault.send(safeAmount);
      require(success);
    } else {
      require(msg.value == 0, "unexpected eth"); // protect user from accidentally locking up eth
      safeAmount = safeTransferFrom(asset, amount); // Transfer asset
      _vaultAllowance[vault][asset] += safeAmount; // Credit to chosen vault
    }
    emit Deposit(vault, asset, safeAmount, memo);
  }

  //############################## ALLOWANCE TRANSFERS ##############################

  // Use for "moving" assets between vaults (asgard<>ygg), as well "churning" to a new Asgard
  function transferAllowance(
    address router,
    address newVault,
    address asset,
    uint amount,
    string memory memo
  ) external nonReentrant {
    if (router == address(this)) {
      _adjustAllowances(newVault, asset, amount);
      emit TransferAllowance(msg.sender, newVault, asset, amount, memo);
    } else {
      _routerDeposit(router, newVault, asset, amount, memo);
    }
  }

  //############################## ASSET TRANSFERS ##############################

  // V4 transferOut kept intact in V5 for backwards compatibility with previous code
  // Any vault calls to transfer any asset to any recipient.
  // Note: Contract recipients of ETH are only given 2300 Gas to complete execution.
  function transferOut(
    address payable to,
    address asset,
    uint amount,
    string memory memo
  ) public payable nonReentrant {
    uint safeAmount;
    if (asset == address(0)) {
      safeAmount = msg.value;
      bool success = to.send(safeAmount); // Send ETH.
      if (!success) {
        payable(address(msg.sender)).transfer(safeAmount); // For failure, bounce back to vault & continue.
      }
    } else {
      _vaultAllowance[msg.sender][asset] -= amount; // Reduce allowance
      (bool success, bytes memory data) = asset.call(
        abi.encodeWithSignature("transfer(address,uint256)", to, amount)
      );
      require(success && (data.length == 0 || abi.decode(data, (bool))));
      safeAmount = amount;
    }
    emit TransferOut(msg.sender, to, asset, safeAmount, memo);
  }

  function _transferOutV5(TransferOutData memory transferOutPayload) private {
    if (transferOutPayload.asset == address(0)) {
      bool success = transferOutPayload.to.send(transferOutPayload.amount); // Send ETH.
      if (!success) {
        payable(address(msg.sender)).transfer(transferOutPayload.amount); // For failure, bounce back to vault & continue.
      }
    } else {
      _vaultAllowance[msg.sender][
        transferOutPayload.asset
      ] -= transferOutPayload.amount; // Reduce allowance

      (bool success, bytes memory data) = transferOutPayload.asset.call(
        abi.encodeWithSignature(
          "transfer(address,uint256)",
          transferOutPayload.to,
          transferOutPayload.amount
        )
      );

      require(success && (data.length == 0 || abi.decode(data, (bool))));
    }

    emit TransferOut(
      msg.sender,
      transferOutPayload.to,
      transferOutPayload.asset,
      transferOutPayload.amount,
      transferOutPayload.memo
    );
  }

  function transferOutV5(
    TransferOutData calldata transferOutPayload
  ) public payable nonReentrant {
    _transferOutV5(transferOutPayload);
  }

  // bifrost to budget gas limits, no more than 50% than L1 gas limit
  function batchTransferOutV5(
    TransferOutData[] calldata transferOutPayload
  ) external payable nonReentrant {
    for (uint i = 0; i < transferOutPayload.length; ++i) {
      _transferOutV5(transferOutPayload[i]);
    }
  }

  // V4 transferOutAndCall kept intact in V5 for backwards compatibility with previous aggregator contracts
  // Any vault calls to transferAndCall on a target contract that conforms with "swapOut(address,address,uint256)"
  // Example Memo: "~1b3:ETH.0xFinalToken:0xTo:"
  // Aggregator is matched to the last three digits of whitelisted aggregators
  // FinalToken, To, amountOutMin come from originating memo
  // Memo passed in here is the "OUT:HASH" type
  function transferOutAndCall(
    address payable aggregator,
    address finalToken,
    address to,
    uint256 amountOutMin,
    string memory memo
  ) public payable nonReentrant {
    uint256 _safeAmount = msg.value;
    (bool erc20Success, ) = aggregator.call{value: _safeAmount}(
      abi.encodeWithSignature(
        "swapOut(address,address,uint256)",
        finalToken,
        to,
        amountOutMin
      )
    );
    if (!erc20Success) {
      bool ethSuccess = payable(to).send(_safeAmount); // If can't swap, just send the recipient the ETH
      if (!ethSuccess) {
        payable(address(msg.sender)).transfer(_safeAmount); // For failure, bounce back to vault & continue.
      }
    }

    emit TransferOutAndCall(
      msg.sender,
      aggregator,
      _safeAmount,
      finalToken,
      to,
      amountOutMin,
      memo
    );
  }

  // Any vault calls to transferAndCall on a target contract that conforms with "swapOut(address,uint256,address,address,uint256,bytes)"
  // Example Memo: "~1b3:ETH.0xFinalToken:0xTo:"
  // Target is fuzzy-matched to the last three digits of whitelisted aggregators
  // toAsset, recipient, amountOutMin come from originating memo
  // Memo passed in here is the "OUT:HASH" type
  // RouterV5: transferOutAndCall can be used with ERC20 tokens as well as Ether
  // RouterV5: payload field can be passed from the originating memo like this thorchainMemo|0xpayload
  // RouterV5: originAddress field is what nodes observed + had consensus on
  // RouterV5: use struct TransferOutAndCallData to pass in the data to reduce local variables
  function _transferOutAndCallV5(
    TransferOutAndCallData calldata aggregationPayload
  ) private {
    if (aggregationPayload.fromAsset == address(0)) {
      // call swapOutV5 with ether
      (bool swapOutSuccess, ) = aggregationPayload.target.call{
        value: msg.value
      }(
        abi.encodeWithSignature(
          "swapOutV5(address,uint256,address,address,uint256,bytes,string)",
          aggregationPayload.fromAsset,
          aggregationPayload.fromAmount,
          aggregationPayload.toAsset,
          aggregationPayload.recipient,
          aggregationPayload.amountOutMin,
          aggregationPayload.payload,
          aggregationPayload.originAddress
        )
      );
      if (!swapOutSuccess) {
        bool sendSuccess = payable(aggregationPayload.target).send(msg.value); // If can't swap, just send the recipient the gas asset
        if (!sendSuccess) {
          payable(address(msg.sender)).transfer(msg.value); // For failure, bounce back to vault & continue.
        }
      }

      emit TransferOutAndCallV5(
        msg.sender,
        aggregationPayload.target,
        msg.value,
        aggregationPayload.toAsset,
        aggregationPayload.recipient,
        aggregationPayload.amountOutMin,
        aggregationPayload.memo,
        aggregationPayload.payload,
        aggregationPayload.originAddress
      );
    } else {
      _vaultAllowance[msg.sender][
        aggregationPayload.fromAsset
      ] -= aggregationPayload.fromAmount; // Reduce allowance

      // send ERC20 to aggregator contract so it can do its thing
      (bool transferSuccess, bytes memory data) = aggregationPayload
        .fromAsset
        .call(
          abi.encodeWithSignature(
            "transfer(address,uint256)",
            aggregationPayload.target,
            aggregationPayload.fromAmount
          )
        );

      require(
        transferSuccess && (data.length == 0 || abi.decode(data, (bool))),
        "Failed to transfer token before dex agg call"
      );

      // add test case if aggregator fails, it should not revert the whole transaction (transferOutAndCallV5 call succeeds)
      // call swapOutV5 with erc20. if the aggregator fails, the transaction should not revert
      (bool _dexAggSuccess, ) = aggregationPayload.target.call{value: 0}(
        abi.encodeWithSignature(
          "swapOutV5(address,uint256,address,address,uint256,bytes,string)",
          aggregationPayload.fromAsset,
          aggregationPayload.fromAmount,
          aggregationPayload.toAsset,
          aggregationPayload.recipient,
          aggregationPayload.amountOutMin,
          aggregationPayload.payload,
          aggregationPayload.originAddress
        )
      );

      emit TransferOutAndCallV5(
        msg.sender,
        aggregationPayload.target,
        aggregationPayload.fromAmount,
        aggregationPayload.toAsset,
        aggregationPayload.recipient,
        aggregationPayload.amountOutMin,
        aggregationPayload.memo,
        aggregationPayload.payload,
        aggregationPayload.originAddress
      );
    }
  }

  function transferOutAndCallV5(
    TransferOutAndCallData calldata aggregationPayload
  ) external payable nonReentrant {
    _transferOutAndCallV5(aggregationPayload);
  }

  function batchTransferOutAndCallV5(
    TransferOutAndCallData[] calldata aggregationPayloads
  ) external payable nonReentrant {
    for (uint i = 0; i < aggregationPayloads.length; ++i) {
      _transferOutAndCallV5(aggregationPayloads[i]);
    }
  }
  
  //############################## VAULT MANAGEMENT ##############################

  // A vault can call to "return" all assets to an asgard, including ETH.
  function returnVaultAssets(
    address router,
    address payable asgard,
    Coin[] memory coins,
    string memory memo
  ) external payable nonReentrant {
    if (router == address(this)) {
      for (uint i = 0; i < coins.length; i++) {
        _adjustAllowances(asgard, coins[i].asset, coins[i].amount);
      }
      emit VaultTransfer(msg.sender, asgard, coins, memo); // Does not include ETH.
    } else {
      for (uint i = 0; i < coins.length; i++) {
        _routerDeposit(router, asgard, coins[i].asset, coins[i].amount, memo);
      }
    }
    bool success = asgard.send(msg.value);
    require(success);
  }

  //############################## HELPERS ##############################

  function vaultAllowance(
    address vault,
    address token
  ) public view returns (uint amount) {
    return _vaultAllowance[vault][token];
  }

  // Safe transferFrom in case asset charges transfer fees
  function safeTransferFrom(
    address _asset,
    uint _amount
  ) internal returns (uint amount) {
    uint _startBal = iERC20(_asset).balanceOf(address(this));
    (bool success, bytes memory data) = _asset.call(
      abi.encodeWithSignature(
        "transferFrom(address,address,uint256)",
        msg.sender,
        address(this),
        _amount
      )
    );
    require(success && (data.length == 0 || abi.decode(data, (bool))));
    return (iERC20(_asset).balanceOf(address(this)) - _startBal);
  }

  // Decrements and Increments Allowances between two vaults
  function _adjustAllowances(
    address _newVault,
    address _asset,
    uint _amount
  ) internal {
    _vaultAllowance[msg.sender][_asset] -= _amount;
    _vaultAllowance[_newVault][_asset] += _amount;
  }

  // Adjust allowance and forwards funds to new router, credits allowance to desired vault
  function _routerDeposit(
    address _router,
    address _vault,
    address _asset,
    uint _amount,
    string memory _memo
  ) internal {
    _vaultAllowance[msg.sender][_asset] -= _amount;
    safeApprove(_asset, _router, _amount);

    iROUTER(_router).depositWithExpiry(
      _vault,
      _asset,
      _amount,
      _memo,
      type(uint).max
    ); // Transfer by depositing
  }

  function safeApprove(
    address _asset,
    address _address,
    uint _amount
  ) internal {
    (bool success, ) = _asset.call(
      abi.encodeWithSignature("approve(address,uint256)", _address, _amount)
    ); // Approve to transfer
    require(success);
  }
}
