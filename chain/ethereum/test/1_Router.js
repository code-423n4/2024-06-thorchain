/*
 * This test deploys the ROUTER1 & TOKEN from User address
 * It then updates Asgard to an Asgard address
 * It then tests for User to Deposit in Assets, which gets forwarded to Asgard
 * It then tests for Asgard to send out those assets back to User
 */

const Router = artifacts.require("THORChain_Router");
const Token = artifacts.require("ERC20Token");
const EvilToken = artifacts.require("EvilERC20Token");
const RevertingContract = artifacts.require("RevertingContract");
const EvilCallback = artifacts.require("EvilCallback");
const Rune = artifacts.require("ETH_RUNE");
const USDT = artifacts.require("TetherToken");

const BigNumber = require("bignumber.js");
const { expect } = require("chai");
const truffleAssert = require("truffle-assertions");

function BN2Str(BN) {
  return new BigNumber(BN).toFixed();
}
function getBN(BN) {
  return new BigNumber(BN);
}

var ROUTER1;
var ROUTER2;
var ROUTER3;
var RUNE, TOKEN, EVILTOKEN;
var ETH = "0x0000000000000000000000000000000000000000";
var usdt;
var REVERTINGCONTRACT, EVILCALLBACK;
var USER1, USER2;
var ASGARD1, ASGARD2;
var YGGDRASIL1, YGGDRASIL2;

const _1 = "1000000000000000000";
const _10 = "10000000000000000000";
const _20 = "20000000000000000000";
const _300 = "300000000000000000000";
const _400 = "400000000000000000000";
const _1000 = "1000000000000000000000";
const _5000 = "5000000000000000000000";
const _50k = "50000000000000000000000";
const _100k = "100000000000000000000000";
const _250k = "250000000000000000000000";
const _500k = "500000000000000000000000";
const _1m = "1000000000000000000000000";
const _9m = "9000000000000000000000000";

const maxUint = require("ethers").constants.MaxUint256;

const currentTime = Math.floor(Date.now() / 1000 + 15 * 60); // time plus 15 mins

describe("Router contract", function () {
  let accounts;

  before(async function () {
    accounts = await web3.eth.getAccounts();
    RUNE = await Rune.new();
    ROUTER1 = await Router.new(RUNE.address);
    ROUTER2 = await Router.new(RUNE.address);
    ROUTER3 = await Router.new(RUNE.address);
    TOKEN = await Token.new(); // User gets 1m TOKENS during construction
    EVILTOKEN = await EvilToken.new();
    REVERTINGCONTRACT = await RevertingContract.new();
    EVILCALLBACK = await EvilCallback.new(ROUTER1.address);
    usdt = await USDT.new(_1m, "Tether", "USDT", 6);
    USER1 = accounts[0];
    USER2 = accounts[1];
    ASGARD1 = accounts[3];
    ASGARD2 = accounts[4];
    YGGDRASIL1 = accounts[7];
    YGGDRASIL2 = accounts[8];
    ASGARD3 = accounts[9];
  });

  describe("User Deposit Assets", function () {
    it("Should Deposit Ether To Asgard1", async function () {
      let startBal = getBN(await web3.eth.getBalance(ASGARD1));
      let tx = await ROUTER1.depositWithExpiry(
        ASGARD1,
        ETH,
        _1000,
        "SWAP:THOR.RUNE",
        currentTime,
        { from: USER1, value: _1000 },
      );
      // console.log(tx.logs[0].args)

      expect(tx.logs[0].event).to.equal("Deposit");
      expect(tx.logs[0].args.asset).to.equal(ETH);
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_1000);
      expect(tx.logs[0].args.memo).to.equal("SWAP:THOR.RUNE");

      let endBal = getBN(await web3.eth.getBalance(ASGARD1));
      let changeBal = BN2Str(endBal.minus(startBal));
      expect(changeBal).to.equal(_1000);
    });

    it("Should revert Deposit Ether To Asgard1", async function () {
      await truffleAssert.reverts(
        ROUTER1.depositWithExpiry(
          ASGARD1,
          ETH,
          _1000,
          "SWAP:THOR.RUNE",
          getBN(0),
          { from: USER1, value: _1000 },
        ),
        "THORChain_Router: expired",
      );
    });

    it("Should Deposit RUNE to Asgard1", async function () {
      expect(BN2Str(await TOKEN.totalSupply())).to.equal(_1m);
      expect(BN2Str(await TOKEN.balanceOf(USER1))).to.equal(_1m);

      // tx
      let tx = await ROUTER1.depositWithExpiry(
        ASGARD1,
        RUNE.address,
        _1m,
        "SWITCH:THOR.RUNE",
        currentTime,
      );
      // console.log(tx.logs)
      expect(tx.logs[0].event).to.equal("Deposit");
      expect(tx.logs[0].args.asset).to.equal(RUNE.address);
      expect(tx.logs[0].args.to).to.equal(ASGARD1);
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_1m);
      expect(tx.logs[0].args.memo).to.equal("SWITCH:THOR.RUNE");

      expect(BN2Str(await RUNE.totalSupply())).to.equal(_9m);

      expect(BN2Str(await RUNE.balanceOf(ROUTER1.address))).to.equal("0");
      expect(BN2Str(await RUNE.balanceOf(USER1))).to.equal(_9m);
      expect(
        BN2Str(await ROUTER1.vaultAllowance(ASGARD1, RUNE.address)),
      ).to.equal("0");
    });

    it("Should revert Deposit RUNE to Asgard1", async function () {
      await truffleAssert.reverts(
        ROUTER1.depositWithExpiry(
          ASGARD1,
          RUNE.address,
          _1m,
          "SWITCH:THOR.RUNE",
          getBN(0),
        ),
        "THORChain_Router: expired",
      );
    });

    it("Should Deposit Token to Asgard1", async function () {
      expect(BN2Str(await TOKEN.totalSupply())).to.equal(_1m);
      expect(BN2Str(await TOKEN.balanceOf(USER1))).to.equal(_1m);

      // Approval - we are approving the ROUTER1 to spend all of USER1 funds
      await TOKEN.approve(ROUTER1.address, _1m, { from: USER1 });
      expect(BN2Str(await TOKEN.allowance(USER1, ROUTER1.address))).to.equal(
        _1m,
      );

      // tx
      let tx = await ROUTER1.depositWithExpiry(
        ASGARD1,
        TOKEN.address,
        _500k,
        "SWAP:THOR.RUNE",
        currentTime,
      );
      // console.log(tx.logs)
      expect(tx.logs[0].event).to.equal("Deposit");
      expect(tx.logs[0].args.asset).to.equal(TOKEN.address);
      expect(tx.logs[0].args.to).to.equal(ASGARD1);
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_500k);
      expect(tx.logs[0].args.memo).to.equal("SWAP:THOR.RUNE");

      expect(BN2Str(await TOKEN.balanceOf(ROUTER1.address))).to.equal(_500k);
      expect(BN2Str(await TOKEN.balanceOf(USER1))).to.equal(_500k);
      expect(
        BN2Str(await ROUTER1.vaultAllowance(ASGARD1, TOKEN.address)),
      ).to.equal(_500k);
    });

    it("Should revert Deposit Token to Asgard1", async function () {
      await truffleAssert.reverts(
        ROUTER1.depositWithExpiry(
          ASGARD1,
          TOKEN.address,
          _500k,
          "SWAP:THOR.RUNE",
          getBN(0),
        ),
        "THORChain_Router: expired",
      );
    });

    it("Should revert when ETH sent during ERC20 Deposit", async function () {
      await truffleAssert.reverts(
        ROUTER1.depositWithExpiry(
          ASGARD1,
          TOKEN.address,
          _1000,
          "MEMO",
          currentTime,
          { from: USER1, value: _1 },
        ),
        "unexpected eth",
      );
    });

    it("Should Deposit USDT to Asgard1", async function () {
      // await usdt.issue(_1m, { from: USER1 });
      expect(BN2Str(await usdt.totalSupply())).to.equal(_1m);
      expect(BN2Str(await usdt.balanceOf(USER1))).to.equal(_1m);

      // Approval - we are approving the ROUTER1 to spend all of USER1 funds
      await usdt.approve(ROUTER1.address, _1m, { from: USER1 });
      expect(BN2Str(await usdt.allowance(USER1, ROUTER1.address))).to.equal(
        _1m,
      );

      // tx
      let tx = await ROUTER1.depositWithExpiry(
        ASGARD1,
        usdt.address,
        _500k,
        "SWAP:THOR.RUNE",
        currentTime,
      );
      // console.log(tx.logs)
      expect(tx.logs[0].event).to.equal("Deposit");
      expect(tx.logs[0].args.asset).to.equal(usdt.address);
      expect(tx.logs[0].args.to).to.equal(ASGARD1);
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_500k);
      expect(tx.logs[0].args.memo).to.equal("SWAP:THOR.RUNE");

      expect(BN2Str(await usdt.balanceOf(ROUTER1.address))).to.equal(_500k);
      expect(BN2Str(await usdt.balanceOf(USER1))).to.equal(_500k);
      expect(
        BN2Str(await ROUTER1.vaultAllowance(ASGARD1, usdt.address)),
      ).to.equal(_500k);
    });
  });

  describe("Fund Yggdrasil, Yggdrasil Transfer Out", function () {
    it("Should fund yggdrasil ETH", async function () {
      let startBal = getBN(await web3.eth.getBalance(YGGDRASIL1));
      let tx = await ROUTER1.transferOut(YGGDRASIL1, ETH, _300, "ygg+:123", {
        from: ASGARD1,
        value: _400,
      });

      expect(tx.logs[0].event).to.equal("TransferOut");
      expect(tx.logs[0].args.asset).to.equal(ETH);
      expect(tx.logs[0].args.vault).to.equal(ASGARD1);
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_400);
      expect(tx.logs[0].args.memo).to.equal("ygg+:123");

      let endBal = getBN(await web3.eth.getBalance(YGGDRASIL1));
      let changeBal = BN2Str(endBal.minus(startBal));
      expect(changeBal).to.equal(_400);
    });

    it("Should fund yggdrasil tokens", async function () {
      let tx = await ROUTER1.transferAllowance(
        ROUTER1.address,
        YGGDRASIL1,
        TOKEN.address,
        _500k,
        "yggdrasil+:1234",
        { from: ASGARD1 },
      );
      expect(tx.logs[0].event).to.equal("TransferAllowance");
      expect(tx.logs[0].args.newVault).to.equal(YGGDRASIL1);
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_500k);

      expect(BN2Str(await TOKEN.balanceOf(ROUTER1.address))).to.equal(_500k);
      expect(
        BN2Str(await ROUTER1.vaultAllowance(YGGDRASIL1, TOKEN.address)),
      ).to.equal(_500k);
      expect(
        BN2Str(await ROUTER1.vaultAllowance(ASGARD1, TOKEN.address)),
      ).to.equal("0");
    });

    it("Should transfer ETH to USER2", async function () {
      let startBal = getBN(await web3.eth.getBalance(USER2));
      let tx = await ROUTER1.transferOut(USER2, ETH, _10, "OUT:", {
        from: YGGDRASIL1,
        value: _10,
      });
      expect(tx.logs[0].event).to.equal("TransferOut");
      expect(tx.logs[0].args.to).to.equal(USER2);
      expect(tx.logs[0].args.asset).to.equal(ETH);
      expect(tx.logs[0].args.memo).to.equal("OUT:");
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_10);

      let endBal = getBN(await web3.eth.getBalance(USER2));
      let changeBal = BN2Str(endBal.minus(startBal));
      expect(changeBal).to.equal(_10);
    });

    it("Should take ETH amount from the amount in transaction, instead of the amount parameter", async function () {
      let startBal = getBN(await web3.eth.getBalance(USER2));
      let tx = await ROUTER1.transferOut(USER2, ETH, _20, "OUT:", {
        from: YGGDRASIL1,
        value: _10,
      });
      expect(tx.logs[0].event).to.equal("TransferOut");
      expect(tx.logs[0].args.to).to.equal(USER2);
      expect(tx.logs[0].args.asset).to.equal(ETH);
      expect(tx.logs[0].args.memo).to.equal("OUT:");
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_10);

      let endBal = getBN(await web3.eth.getBalance(USER2));
      let changeBal = BN2Str(endBal.minus(startBal));
      expect(changeBal).to.equal(_10);
    });

    it("Should transfer tokens to USER2", async function () {
      let tx = await ROUTER1.transferOut(USER2, TOKEN.address, _250k, "OUT:", {
        from: YGGDRASIL1,
      });
      expect(tx.logs[0].event).to.equal("TransferOut");
      expect(tx.logs[0].args.to).to.equal(USER2);
      expect(tx.logs[0].args.asset).to.equal(TOKEN.address);
      expect(tx.logs[0].args.memo).to.equal("OUT:");
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_250k);

      expect(BN2Str(await TOKEN.balanceOf(ROUTER1.address))).to.equal(_250k);
      expect(
        BN2Str(await ROUTER1.vaultAllowance(YGGDRASIL1, TOKEN.address)),
      ).to.equal(_250k);
    });
    it("Should transfer USDT to USER2", async function () {
      let tx = await ROUTER1.transferOut(USER2, usdt.address, _500k, "OUT:", {
        from: ASGARD1,
      });
      expect(tx.logs[0].event).to.equal("TransferOut");
      expect(tx.logs[0].args.to).to.equal(USER2);
      expect(tx.logs[0].args.asset).to.equal(usdt.address);
      expect(tx.logs[0].args.memo).to.equal("OUT:");
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_500k);

      expect(BN2Str(await usdt.balanceOf(ROUTER1.address))).to.equal("0");
      expect(
        BN2Str(await ROUTER1.vaultAllowance(YGGDRASIL1, usdt.address)),
      ).to.equal("0");
    });
  });

  describe("Yggdrasil Returns Funds, Asgard Churns, Old Vaults can't spend", function () {
    it("Ygg returns", async function () {
      let ethBal = _20;
      let coins = {
        asset: TOKEN.address,
        amount: _250k,
      };
      let tx = await ROUTER1.returnVaultAssets(
        ROUTER1.address,
        ASGARD1,
        [coins],
        "yggdrasil-:1234",
        { from: YGGDRASIL1, value: ethBal },
      );
      expect(tx.logs[0].event).to.equal("VaultTransfer");
      expect(tx.logs[0].args.coins[0].asset).to.equal(TOKEN.address);
      expect(BN2Str(tx.logs[0].args.coins[0].amount)).to.equal(_250k);
      expect(tx.logs[0].args.memo).to.equal("yggdrasil-:1234");

      expect(BN2Str(await TOKEN.balanceOf(ROUTER1.address))).to.equal(_250k);
      expect(
        BN2Str(await ROUTER1.vaultAllowance(YGGDRASIL1, TOKEN.address)),
      ).to.equal("0");
      expect(
        BN2Str(await ROUTER1.vaultAllowance(ASGARD1, TOKEN.address)),
      ).to.equal(_250k);
    });
    it("Asgard Churns", async function () {
      let tx = await ROUTER1.transferAllowance(
        ROUTER1.address,
        ASGARD2,
        TOKEN.address,
        _250k,
        "migrate:1234",
        { from: ASGARD1 },
      );
      expect(tx.logs[0].event).to.equal("TransferAllowance");
      expect(tx.logs[0].args.asset).to.equal(TOKEN.address);
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_250k);

      expect(BN2Str(await TOKEN.balanceOf(ROUTER1.address))).to.equal(_250k);
      expect(
        BN2Str(await ROUTER1.vaultAllowance(ASGARD1, TOKEN.address)),
      ).to.equal("0");
      expect(
        BN2Str(await ROUTER1.vaultAllowance(ASGARD2, TOKEN.address)),
      ).to.equal(_250k);
    });
    it("Should fail to when old Asgard interacts", async function () {
      await truffleAssert.fails(
        ROUTER1.transferAllowance(
          ROUTER1.address,
          ASGARD2,
          TOKEN.address,
          _50k,
          "migrate:1234",
          { from: ASGARD1 },
        ),
        truffleAssert.ErrorType.REVERT,
      );
      await truffleAssert.fails(
        ROUTER1.transferOut(USER2, TOKEN.address, _50k, "OUT:", {
          from: ASGARD1,
        }),
        truffleAssert.ErrorType.REVERT,
      );
    });
    it("Should fail to when old Yggdrasil interacts", async function () {
      await truffleAssert.fails(
        ROUTER1.transferAllowance(
          ROUTER1.address,
          ASGARD2,
          TOKEN.address,
          _50k,
          "migrate:1234",
          { from: YGGDRASIL1 },
        ),
        truffleAssert.ErrorType.REVERT,
      );
      await truffleAssert.fails(
        ROUTER1.transferOut(USER2, TOKEN.address, _50k, "OUT:", {
          from: YGGDRASIL1,
        }),
        truffleAssert.ErrorType.REVERT,
      );
    });
  });

  describe("Upgrade contract", function () {
    it("should return vault assets to new router", async function () {
      let asgard1StartBalance = getBN(await web3.eth.getBalance(YGGDRASIL1));
      await ROUTER1.depositWithExpiry(
        YGGDRASIL1,
        TOKEN.address,
        _50k,
        "SEED",
        currentTime,
        { from: USER1 },
      );
      await ROUTER1.depositWithExpiry(
        YGGDRASIL1,
        usdt.address,
        _50k,
        "SEED",
        currentTime,
        { from: USER1 },
      );
      await ROUTER1.depositWithExpiry(
        YGGDRASIL1,
        ETH,
        "0",
        "SEED ETH",
        currentTime,
        { from: accounts[10], value: _1 },
      );

      let ethBal = "1000000000000";
      // console.log(ethBal)
      // migrate _50k from asgard1 to asgard3 , to new Router3 contract

      let coin1 = {
        asset: TOKEN.address,
        amount: _50k,
      };
      let coin2 = {
        asset: usdt.address,
        amount: _50k,
      };

      let tx = await ROUTER1.returnVaultAssets(
        ROUTER3.address,
        ASGARD3,
        [coin1, coin2],
        "yggdrasil-:1234",
        { from: YGGDRASIL1, value: ethBal },
      );

      //console.log(tx.logs);
      expect(tx.logs[0].event).to.equal("Deposit");
      expect(tx.logs[0].args.to).to.equal(ASGARD3);
      expect(tx.logs[0].args.asset).to.equal(TOKEN.address);
      expect(tx.logs[0].args.memo).to.equal("yggdrasil-:1234");
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_50k);

      // make sure the token had been transfer to ASGARD3 and Router3
      expect(BN2Str(await TOKEN.balanceOf(ROUTER3.address))).to.equal(_50k);
      expect(
        BN2Str(await ROUTER3.vaultAllowance(ASGARD3, TOKEN.address)),
      ).to.equal(_50k);
      expect(
        BN2Str(await ROUTER1.vaultAllowance(ASGARD1, TOKEN.address)),
      ).to.equal("0");
    });

    it("should transfer all token and allowance to new contract", async function () {
      let asgard1StartBalance = getBN(await web3.eth.getBalance(ASGARD1));
      await ROUTER1.depositWithExpiry(
        ASGARD1,
        TOKEN.address,
        _50k,
        "SEED",
        currentTime,
        { from: USER1 },
      );
      await ROUTER1.depositWithExpiry(
        ASGARD1,
        usdt.address,
        _50k,
        "SEED",
        currentTime,
        { from: USER1 },
      );
      await ROUTER1.depositWithExpiry(
        ASGARD1,
        ETH,
        "0",
        "SEED ETH",
        currentTime,
        { from: accounts[10], value: _1 },
      );

      let asgard1EndBalance = getBN(await web3.eth.getBalance(ASGARD1));
      expect(BN2Str(asgard1EndBalance.minus(asgard1StartBalance))).to.equal(_1);
      // migrate _50k from asgard1 to asgard3 , to new Router3 contract
      let tx = await ROUTER1.transferAllowance(
        ROUTER3.address,
        ASGARD3,
        TOKEN.address,
        _50k,
        "MIGRATE:1",
        { from: ASGARD1 },
      );
      //console.log(tx.logs);
      expect(tx.logs[0].event).to.equal("Deposit");
      expect(tx.logs[0].args.to).to.equal(ASGARD3);
      expect(tx.logs[0].args.asset).to.equal(TOKEN.address);
      expect(tx.logs[0].args.memo).to.equal("MIGRATE:1");
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_50k);

      // make sure the token had been transfer to ASGARD3 and Router3
      expect(BN2Str(await TOKEN.balanceOf(ROUTER3.address))).to.equal(_100k);
      expect(
        BN2Str(await ROUTER3.vaultAllowance(ASGARD3, TOKEN.address)),
      ).to.equal(_100k);
      expect(
        BN2Str(await ROUTER1.vaultAllowance(ASGARD1, TOKEN.address)),
      ).to.equal("0");

      let tx2 = await ROUTER1.transferAllowance(
        ROUTER3.address,
        ASGARD3,
        usdt.address,
        _50k,
        "MIGRATE:1",
        { from: ASGARD1 },
      );
      expect(tx2.logs[0].event).to.equal("Deposit");
      expect(tx2.logs[0].args.to).to.equal(ASGARD3);
      expect(tx2.logs[0].args.asset).to.equal(usdt.address);
      expect(tx2.logs[0].args.memo).to.equal("MIGRATE:1");
      expect(BN2Str(tx2.logs[0].args.amount)).to.equal(_50k);

      // make sure the token had been transfer to ASGARD3 and Router3
      expect(BN2Str(await usdt.balanceOf(ROUTER3.address))).to.equal(_100k);
      expect(
        BN2Str(await ROUTER3.vaultAllowance(ASGARD3, usdt.address)),
      ).to.equal(_100k);
      expect(
        BN2Str(await ROUTER1.vaultAllowance(ASGARD1, usdt.address)),
      ).to.equal("0");

      let asgard3StartBalance = getBN(await web3.eth.getBalance(ASGARD3));
      // this ignore the gas cost on ASGARD1
      // transfer out ETH.ETH
      let tx1 = await ROUTER1.transferOut(ASGARD3, ETH, "0", "MIGRATE:1", {
        from: ASGARD1,
        value: _5000,
      });
      // console.log(tx1.logs)
      expect(tx1.logs[0].event).to.equal("TransferOut");
      expect(tx1.logs[0].args.vault).to.equal(ASGARD1);
      expect(tx1.logs[0].args.to).to.equal(ASGARD3);
      expect(tx1.logs[0].args.asset).to.equal(ETH);
      expect(tx1.logs[0].args.memo).to.equal("MIGRATE:1");

      let asgard3EndBalance = getBN(await web3.eth.getBalance(ASGARD3));
      expect(BN2Str(asgard3EndBalance.minus(asgard3StartBalance))).to.equal(
        _5000,
      );
    });
  });

  describe("Evil callbacks", function () {
    it("should not give more allowance than tokens transferred", async function () {
      let startAllowance = getBN(
        await ROUTER1.vaultAllowance(USER1, EVILTOKEN.address),
      );
      await EVILTOKEN.approve(ROUTER1.address, _10, { from: USER1 });

      await truffleAssert.reverts(
        ROUTER1.depositWithExpiry(
          USER1,
          EVILTOKEN.address,
          _10,
          "",
          currentTime,
          { from: USER1 },
        ),
      );
    });
    it("Test transferOut reverting contract", async function () {
      let startEvilBal = await web3.eth.getBalance(REVERTINGCONTRACT.address);
      let tx = await ROUTER1.transferOut(
        REVERTINGCONTRACT.address,
        ETH,
        _10,
        "OUT:",
        { from: YGGDRASIL1, value: _10 },
      );
      // Should emit a valid log ("fire and forget")
      expect(tx.logs[0].event).to.equal("TransferOut");
      expect(tx.logs[0].args.to).to.equal(REVERTINGCONTRACT.address);
      expect(tx.logs[0].args.asset).to.equal(ETH);
      expect(tx.logs[0].args.memo).to.equal("OUT:");
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_10);
      // The REVERTINGCONTRACT will have no ETH
      let endEvilBal = await web3.eth.getBalance(REVERTINGCONTRACT.address);
      expect(startEvilBal).to.equal(endEvilBal);
    });
    it("Test transferOut 'to' recipient tries re-entrancy", async function () {
      let startEvilBal = await web3.eth.getBalance(EVILCALLBACK.address);
      let tx = await ROUTER1.transferOut(
        EVILCALLBACK.address,
        ETH,
        _10,
        "OUT:",
        { from: YGGDRASIL1, value: _10 },
      );
      // Should emit a valid log ("fire and forget")
      expect(tx.logs[0].event).to.equal("TransferOut");
      expect(tx.logs[0].args.to).to.equal(EVILCALLBACK.address);
      expect(tx.logs[0].args.asset).to.equal(ETH);
      expect(tx.logs[0].args.memo).to.equal("OUT:");
      expect(BN2Str(tx.logs[0].args.amount)).to.equal(_10);
      // The EVILCALLBACK will try to re-entrant but run out of gas, resulting in balance unchanged.
      let endEvilBal = await web3.eth.getBalance(EVILCALLBACK.address);
      expect(startEvilBal).to.equal(endEvilBal);
    });
    it("Test re-entrancy protection (generic)", async function () {
      await truffleAssert.fails(
        ROUTER1.transferAllowance(
          EVILCALLBACK.address,
          ASGARD1,
          ETH,
          0,
          "EvilReentrancy",
          { from: USER1 },
        ),
        truffleAssert.ErrorType.REVERT,
      );
    });
  });
});
