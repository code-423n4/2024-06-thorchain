const Router = artifacts.require("THORChain_Router");
const Aggregator = artifacts.require("THORChain_Aggregator");
const SushiRouter = artifacts.require("SushiRouterSmol");
const RevertingContract = artifacts.require("RevertingContract");
const Token = artifacts.require("ERC20Token");
const Rune = artifacts.require("ETH_RUNE");
const Usdt = artifacts.require("TetherToken");
const Weth = artifacts.require("WETH");
const BigNumber = require("bignumber.js");
const { expect } = require("chai");
const truffleAssert = require("truffle-assertions");
function BN2Str(BN) {
  return new BigNumber(BN).toFixed();
}
function getBN(BN) {
  return new BigNumber(BN);
}

var ROUTER;
var ASGARD;
var AGG;
var WETH;
var SUSHIROUTER;
var REVERTINGCONTRACT;
var RUNE, TOKEN;
var USDT;
var WETH;
var ETH = "0x0000000000000000000000000000000000000000";
var USER1;
var USER2;

const _1 = "1000000000000000000";
const _2 = "2000000000000000000";
const _3 = "3000000000000000000";
const _10 = "10000000000000000000";
const _20 = "20000000000000000000";
const _1m = "1000000000000000000000000";

const currentTime = Math.floor(Date.now() / 1000 + 15 * 60); // time plus 15 mins

describe("Batch Outbounds", function () {
  let accounts;

  before(async function () {
    accounts = await web3.eth.getAccounts();
    RUNE = await Rune.new();
    ROUTER = await Router.new();
    TOKEN = await Token.new(); // User gets 1m TOKENS during construction
    USDT = await Usdt.new(_1m, "Tether", "USDT", 6);
    USER1 = accounts[0];
    USER2 = accounts[1];
    ASGARD = accounts[3];

    WETH = await Weth.new();
    SUSHIROUTER = await SushiRouter.new(WETH.address);
    AGG = await Aggregator.new(WETH.address, SUSHIROUTER.address);
    REVERTINGCONTRACT = await RevertingContract.new();
  });

  describe("Swap In and Out", function () {
    it("Should Deposit Assets to Router", async function () {
      await TOKEN.transfer(SUSHIROUTER.address, _10);
      await USDT.transfer(SUSHIROUTER.address, _10);
      await web3.eth.sendTransaction({
        to: SUSHIROUTER.address,
        from: USER1,
        value: _10,
      });
      await web3.eth.sendTransaction({
        to: WETH.address,
        from: USER1,
        value: _10,
      });
      await WETH.transfer(SUSHIROUTER.address, _10);

      await TOKEN.approve(ROUTER.address, _1m, { from: USER1 });
      await ROUTER.depositWithExpiry(
        ASGARD,
        TOKEN.address,
        _20,
        "SWAP:THOR.RUNE",
        currentTime,
        {
          from: USER1,
        },
      );

      expect(BN2Str(await TOKEN.balanceOf(SUSHIROUTER.address))).to.equal(_10);
      expect(BN2Str(await USDT.balanceOf(SUSHIROUTER.address))).to.equal(_10);
      expect(BN2Str(await web3.eth.getBalance(SUSHIROUTER.address))).to.equal(
        _10,
      );
      expect(BN2Str(await WETH.balanceOf(SUSHIROUTER.address))).to.equal(_10);
      expect(BN2Str(await TOKEN.balanceOf(ROUTER.address))).to.equal(_20);
    });

    it("Should send 4 outbound transactions", async function () {
      expect(BN2Str(await web3.eth.getBalance(ASGARD))).to.equal(
        "5618997418429237272913",
      );
      expect(BN2Str(await web3.eth.getBalance(USER1))).to.equal(
        "8962897962844698362208",
      );
      expect(BN2Str(await web3.eth.getBalance(USER2))).to.equal(
        "10020000000000000000000",
      );
      expect(BN2Str(await TOKEN.balanceOf(USER1))).to.equal(
        "999970000000000000000000",
      );
      expect(BN2Str(await TOKEN.balanceOf(USER2))).to.equal("0");

      await ROUTER.batchTransferOutV5(
        [
          [USER1, ETH, _1, "OUT:ETH-USER1"],
          [USER1, ETH, _1, "OUT:ETH-USER1-AGAIN"],
          [USER2, ETH, _1, "OUT:ETH-USER2"],
          [USER2, TOKEN.address, _1, "OUT:TOKEN-USER2"],
        ],
        { from: ASGARD, value: _3 },
      );

      expect(BN2Str(await web3.eth.getBalance(ASGARD))).to.equal(
        "5615997136338575602609",
      ); // Less 3 eth and some gas

      expect(BN2Str(await web3.eth.getBalance(USER1))).to.equal(
        "8964897962844698362208",
      ); // Add 2 eth

      expect(BN2Str(await web3.eth.getBalance(USER2))).to.equal(
        "10021000000000000000000",
      ); // Add 1 eth

      expect(BN2Str(await TOKEN.balanceOf(USER1))).to.equal(
        "999970000000000000000000",
      ); // no change
      expect(BN2Str(await TOKEN.balanceOf(USER2))).to.equal(
        "1000000000000000000",
      ); // Add 1 TOKEN
    });

    it("Should fail if not enough ether is sent", async function () {
      await truffleAssert.reverts(
        ROUTER.batchTransferOutV5(
          [
            [USER1, ETH, _1, "OUT:ETH-USER1"],
            [USER1, ETH, _1, "OUT:ETH-USER1-AGAIN"],
            [USER2, ETH, _1, "OUT:ETH-USER2"],
            [USER2, TOKEN.address, _1, "OUT:TOKEN-USER2"],
          ],
          { from: ASGARD, value: _2 },
        ),
        "Transaction reverted: function call failed to execute",
      );
    });

    it("Should send 2 dex agg outbounds", async function () {
      expect(BN2Str(await web3.eth.getBalance(ASGARD))).to.equal(
        "5615996948578189793361",
      ); // asgard starting eth balance

      expect(BN2Str(await web3.eth.getBalance(USER1))).to.equal(
        "8964897962844698362208",
      ); // user1 starting eth balance
      expect(BN2Str(await TOKEN.balanceOf(USER1))).to.equal(
        "999970000000000000000000",
      ); // user1 starting token balance

      expect(BN2Str(await web3.eth.getBalance(USER2))).to.equal(
        "10021000000000000000000",
      ); // user2 starting eth balance
      expect(BN2Str(await TOKEN.balanceOf(USER2))).to.equal(
        "1000000000000000000",
      ); // user2 starting token balance

      expect(BN2Str(await TOKEN.balanceOf(ROUTER.address))).to.equal(
        "19000000000000000000",
      ); // router starting token balance

      await ROUTER.batchTransferOutAndCallV5(
        [
          [
            AGG.address,
            TOKEN.address,
            _1,
            ETH,
            USER1,
            "0",
            "OUT:HASH",
            "0x", // empty payload
            "bc123", // dummy address
          ],
          [
            AGG.address,
            ETH,
            _1,
            TOKEN.address,
            USER2,
            "0",
            "OUT:HASH",
            "0x", // empty payload
            "bc123", // dummy address
          ],
        ],
        { from: ASGARD, value: "0" },
      );

      expect(BN2Str(await web3.eth.getBalance(ASGARD))).to.equal(
        "5615996530297437056673",
      ); // less some gas

      expect(BN2Str(await web3.eth.getBalance(USER1))).to.equal(
        "8965897962844698362208",
      ); // Add 1 eth

      expect(BN2Str(await TOKEN.balanceOf(USER1))).to.equal(
        "999970000000000000000000",
      ); // No change

      expect(BN2Str(await web3.eth.getBalance(USER2))).to.equal(
        "10021000000000000000000",
      ); // No change

      expect(BN2Str(await TOKEN.balanceOf(USER2))).to.equal(
        "2000000000000000000",
      ); // Add 1 TOKEN
    });
  });
});
