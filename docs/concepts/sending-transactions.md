# Sending Transactions

Confirm you have:

- [ ] Connected to Midgard or THORNode
- [ ] Located the latest vault (and router) for the chain
- [ ] Prepared the transaction details (and memo)
- [ ] Checked the network is not halted for your transaction

You are ready to make the transaction and swap via THORChain.

## UTXO Chains

> ⚠️ THORChain does NOT currently support BTC Taproot. User funds will be lost if sent to or from a taproot address!

- [ ] Ensure the [address type](./querying-thorchain.md#supported-address-formats) is supported
- [ ] Send the transaction with Asgard vault as VOUT0
- [ ] Pass all change back to the VIN0 address in a subsequent VOUT e.g. VOUT1
- [ ] Include the memo as an OP_RETURN in a subsequent VOUT e.g. VOUT2
- [ ] Use a high enough `gas_rate` to be included
- [ ] Do not send below the dust threshold (10k Sats BTC, BCH, LTC, 1m DOGE), exhaustive values can be found on the [Inbound Addresses](https://thornode.ninerealms.com/thorchain/inbound_addresses) endpoint
- [ ] Do not send funds that are part of a transaction with more than 10 outputs

```admonish warning
Inbound transactions should not be delayed for any reason else there is risk funds will be sent to an unreachable address. Use standard transactions, check the [`Inbound_Address`](querying-thorchain.md#getting-the-asgard-vault) before sending and use the recommended [`gas rate`](querying-thorchain.md#getting-the-asgard-vault) to ensure transactions are confirmed in the next block to the latest `Inbound_Address`.
```

```admonish info
Memo limited to 80 bytes on BTC, BCH, LTC and DOGE. Use abbreviated options and [THORNames](https://docs.thorchain.org/network/thorchain-name-service) where possible.
```

```admonish warning
Do not use HD wallets that forward the change to a new address, because THORChain IDs the user as the address in VIN0. The user must keep their VIN0 address funded for refunds.
```

```admonish danger
Override randomised VOUT ordering; THORChain requires specific output ordering. Funds using wrong ordering are very likely to be lost.
```

### EVM Chains

{{#embed https://gitlab.com/thorchain/ethereum/eth-router/-/blob/master/contracts/THORChain_Router.sol#L66 }}

```go
depositWithExpiry(vault, asset, amount, memo, expiry)
```

- [ ] If ERC20, approve the router to spend an allowance of the token first
- [ ] Send the transaction as a `depositWithExpiry()` on the router
- [ ] Vault is the Asgard vault address, asset is the token address to swap, memo as a string
- [ ] Use an expiry which is +60mins on the current time (if the tx is delayed, it will get refunded). The timestamp is in seconds (Solidity's `block.timestamp`).
- [ ] Use a high enough `gas_rate` to be included, otherwise the tx will get stuck

```admonish info
ETH is `0x0000000000000000000000000000000000000000`
```

```admonish danger
ETH is sent and received as an internal transaction. Your wallet may not be set to read internal balances and transactions.
```

```admonish danger
Do not send assets from a smart contract (including smart contract wallets) without adding your contract to the whitelist. As a security measure, Thorchain ignores transactions coming from unknown smart contracts, resulting in a loss of funds.
```

### BFT Chains

- [ ] Send the transaction to the Asgard vault
- [ ] Include the memo
- [ ] Only use the base asset as the choice for gas asset

## THORChain

To initiate a $RUNE -> $ASSET swap a `MsgDeposit` must be broadcasted to the THORChain blockchain. The `MsgDeposit` does not have a destination address, and has the following properties. The full definition can be found [here](https://gitlab.com/thorchain/thornode/-/blob/develop/x/thorchain/types/msg_deposit.go).

```go
MsgDeposit{
    Coins:  coins,
    Memo:   memo,
    Signer: signer,
}
```

If you are using Javascript, [CosmJS](https://github.com/cosmos/cosmjs) is the recommended package to build and broadcast custom message types. [Here is a walkthrough](https://github.com/cosmos/cosmjs/blob/main/packages/stargate/CUSTOM_PROTOBUF_CODECS.md).

### Code Examples (Javascript)

1. **Generate codec files.** To build/broadcast native transactions in Javascript/Typescript, the protobuf files need to be generated into js types. The below script uses `pbjs` and `pbts` to generate the types using the relevant files from the THORNode repo. Alternatively, the .`js` and `.d.ts` files can be downloaded directly from the [XChainJS repo](https://github.com/xchainjs/xchainjs-lib/tree/master/packages/xchain-thorchain/src/types/proto).

   ```bash
   #!/bin/bash

   # this script checks out thornode master and generates the proto3 typescript buindings for MsgDeposit and MsgSend

   MSG_COMPILED_OUTPUTFILE=src/types/proto/MsgCompiled.js
   MSG_COMPILED_TYPES_OUTPUTFILE=src/types/proto/MsgCompiled.d.ts

   TMP_DIR=$(mktemp -d)

   tput setaf 2; echo "Checking out https://gitlab.com/thorchain/thornode  to $TMP_DIR";tput sgr0
   (cd $TMP_DIR && git clone https://gitlab.com/thorchain/thornode)

   # Generate msgs
   tput setaf 2; echo "Generating $MSG_COMPILED_OUTPUTFILE";tput sgr0
   yarn run pbjs -w commonjs  -t static-module $TMP_DIR/thornode/proto/thorchain/v1/common/common.proto $TMP_DIR/thornode/proto/thorchain/v1/x/thorchain/types/msg_deposit.proto $TMP_DIR/thornode/proto/thorchain/v1/x/thorchain/types/msg_send.proto $TMP_DIR/thornode/third_party/proto/cosmos/base/v1beta1/coin.proto -o $MSG_COMPILED_OUTPUTFILE

   tput setaf 2; echo "Generating $MSG_COMPILED_TYPES_OUTPUTFILE";tput sgr0
   yarn run pbts  $MSG_COMPILED_OUTPUTFILE -o $MSG_COMPILED_TYPES_OUTPUTFILE

   tput setaf 2; echo "Removing $TMP_DIR/thornode";tput sgr0
   rm -rf $TMP_DIR
   ```

2. **Using @cosmjs build/broadcast the TX.**

   ```javascript
   const {
     DirectSecp256k1HdWallet,
     Registry,
   } = require("@cosmjs/proto-signing");
   const {
     defaultRegistryTypes: defaultStargateTypes,
     SigningStargateClient,
   } = require("@cosmjs/stargate");
   const { stringToPath } = require("@cosmjs/crypto");
   const bech32 = require("bech32-buffer");

   const { MsgDeposit } = require("./types/MsgCompiled").types;

   async function main() {
     const myRegistry = new Registry(defaultStargateTypes);
     myRegistry.register("/types.MsgDeposit", MsgDeposit);

     const signerMnemonic = "mnemonic here";
     const signerAddr = "thor1...";

     const signer = await DirectSecp256k1HdWallet.fromMnemonic(signerMnemonic, {
       prefix: "thor", // THORChain prefix
       hdPaths: [stringToPath("m/44'/931'/0'/0/0")], // THORChain HD Path
     });

     const client = await SigningStargateClient.connectWithSigner(
       "https://rpc.ninerealms.com/",
       signer,
       { registry: myRegistry },
     );

     const memo = `=:BNB/BNB:${signerAddr}`; // THORChain memo

     const msg = {
       coins: [
         {
           asset: {
             chain: "THOR",
             symbol: "RUNE",
             ticker: "RUNE",
           },
           amount: "100000000", // Value in 1e8 (100000000 = 1 RUNE)
         },
       ],
       memo: memo,
       signer: bech32.decode(signerAddr).data,
     };

     const depositMsg = {
       typeUrl: "/types.MsgDeposit",
       value: MsgDeposit.fromObject(msg),
     };

     const fee = {
       amount: [],
       gas: "50000000", // Set arbitrarily high gas limit; this is not actually deducted from user account.
     };

     const response = await client.signAndBroadcast(
       signerAddr,
       [depositMsg],
       fee,
       memo,
     );
     console.log("response: ", response);

     if (response.code !== 0) {
       console.log("Error: ", response.rawLog);
     } else {
       console.log("Success!");
     }
   }

   main();
   ```

### Native Transaction Fee

As of [ADR-009](https://gitlab.com/thorchain/thornode/-/blob/develop/docs/architecture/adr-009-reserve-income-fee-overhaul.md), the native transaction fee for $RUNE transfers or inbound swaps is USD-denominated, but ultimately paid in $RUNE, which means the fee is dynamic. Interfaces should pull the native transaction fee from THORNode before each new transaction is built/broadcasted.

**THORNode Network Endpoint**: [/thorchain/network](https://thornode.ninerealms.com/thorchain/network)

```json
{
  ...
  "native_outbound_fee_rune": "2000000", // (1e8) Outbound fee for $Asset -> $RUNE swaps
  "native_tx_fee_rune": "2000000", // (1e8) Fee for $RUNE transfers or $RUNE -> $Asset swaps
  ...
  "rune_price_in_tor": "354518918", // (1e8) Current $RUNE price in USD
  ...
}
```

The native transaction fee is automatically deducted from the user's account for $RUNE transfers and inbound swaps. Ensure the user's balance exceeds `tx amount + native_tx_fee_rune` before broadcasting the transaction.
