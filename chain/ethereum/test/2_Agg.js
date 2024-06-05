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

const _1 = "1000000000000000000";
const _10 = "10000000000000000000";
const _20 = "20000000000000000000";
const _1m = "1000000000000000000000000";

describe("Aggregator contract", function () {
  let accounts;

  before(async function () {
    accounts = await web3.eth.getAccounts();
    RUNE = await Rune.new();
    ROUTER = await Router.new(RUNE.address);
    TOKEN = await Token.new(); // User gets 1m TOKENS during construction
    USDT = await Usdt.new(_1m, "Tether", "USDT", 6);
    USER1 = accounts[0];
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

      expect(BN2Str(await TOKEN.balanceOf(SUSHIROUTER.address))).to.equal(_10);
      expect(BN2Str(await USDT.balanceOf(SUSHIROUTER.address))).to.equal(_10);
      expect(BN2Str(await web3.eth.getBalance(SUSHIROUTER.address))).to.equal(
        _10,
      );
      expect(BN2Str(await WETH.balanceOf(SUSHIROUTER.address))).to.equal(_10);
    });

    it("Should Swap In Token using Aggregator", async function () {
      expect(BN2Str(await web3.eth.getBalance(ASGARD))).to.equal(
        "5620999504914401621842",
      ); // 1000 ETH

      // Approval - we are approving the AGG to spend all of USER1 funds
      await TOKEN.approve(AGG.address, _1m, { from: USER1 });

      let deadline = ~~(Date.now() / 1000) + 100;

      // Send 10 token to agg, which sends it to Sushi for 1 WETH,
      // Then unwraps to 1 ETH, then sends 1 ETH to Asgard vault
      await AGG.swapIn(
        ASGARD,
        ROUTER.address,
        "SWAP:BTC.BTC:bc1Address:",
        TOKEN.address,
        _10,
        0,
        BN2Str(deadline),
      );

      expect(BN2Str(await TOKEN.balanceOf(USER1))).to.equal(
        "999980000000000000000000",
      ); // Less 10
      expect(BN2Str(await TOKEN.balanceOf(SUSHIROUTER.address))).to.equal(_20); // Add 10
      expect(BN2Str(await WETH.balanceOf(SUSHIROUTER.address))).to.equal(
        "9000000000000000000",
      ); // Less 1 WETH
      expect(BN2Str(await web3.eth.getBalance(SUSHIROUTER.address))).to.equal(
        _10,
      ); // No change
      expect(BN2Str(await web3.eth.getBalance(ASGARD))).to.equal(
        "5621999504914401621842",
      ); // Add 1 ETH
    });

    it("Should Swap In USDT using Aggregator", async function () {
      // Approval - we are approving the AGG to spend all of USER1 funds
      await USDT.approve(AGG.address, _1m, { from: USER1 });

      let deadline = ~~(Date.now() / 1000) + 100;

      // Send 10 token to agg, which sends it to Sushi for 1 WETH,
      // Then unwraps to 1 ETH, then sends 1 ETH to Asgard vault
      await AGG.swapIn(
        ASGARD,
        ROUTER.address,
        "SWAP:BTC.BTC:bc1Address:",
        USDT.address,
        _10,
        0,
        BN2Str(deadline),
      );

      expect(BN2Str(await USDT.balanceOf(USER1))).to.equal(
        "999980000000000000000000",
      ); // Less 10
      expect(BN2Str(await USDT.balanceOf(SUSHIROUTER.address))).to.equal(_20); // Add 10
      expect(BN2Str(await WETH.balanceOf(SUSHIROUTER.address))).to.equal(
        "8000000000000000000",
      ); // Less 1 WETH
      expect(BN2Str(await web3.eth.getBalance(SUSHIROUTER.address))).to.equal(
        _10,
      ); // No change
      expect(BN2Str(await web3.eth.getBalance(ASGARD))).to.equal(
        "5622999504914401621842",
      ); // Add 1 ETH
    });

    it("Should Swap Out using Aggregator", async function () {
      // Asgard transferOutAndCall() 1 ETH
      // Send 1 ETH to router, forward to agg, forward to Sushi
      // Wraps to 1 WETH, then sends 1 token to user

      await ROUTER.transferOutAndCall(
        AGG.address,
        TOKEN.address,
        USER1,
        "0",
        "OUT:HASH",
        { from: ASGARD, value: _1 },
      );

      expect(BN2Str(await web3.eth.getBalance(ASGARD))).to.equal(
        "5621999408399091952910",
      ); // Less 1 ETH
      expect(BN2Str(await web3.eth.getBalance(SUSHIROUTER.address))).to.equal(
        _10,
      ); // No change
      expect(BN2Str(await WETH.balanceOf(SUSHIROUTER.address))).to.equal(
        "9000000000000000000",
      ); // Add 1 WETH
      expect(BN2Str(await TOKEN.balanceOf(SUSHIROUTER.address))).to.equal(
        "19000000000000000000",
      ); // Less 1
      expect(BN2Str(await TOKEN.balanceOf(USER1))).to.equal(
        "999981000000000000000000",
      ); // Add 1
    });

    it("Should Fail Swap Out using Aggregator", async function () {
      expect(BN2Str(await web3.eth.getBalance(USER1))).to.equal(
        "8979971387836268055422",
      ); // Start bal
      expect(BN2Str(await USDT.balanceOf(SUSHIROUTER.address))).to.equal(
        "20000000000000000000",
      ); // Start bal
      expect(BN2Str(await USDT.balanceOf(USER1))).to.equal(
        "999980000000000000000000",
      ); // Start bal
      // Asgard transferOutAndCall() 1 ETH
      // Send 1 ETH to router, forward to agg, forward to Sushi
      // Fail due price check, send back to asgard, sent ETH to user
      await ROUTER.transferOutAndCall(
        AGG.address,
        USDT.address,
        USER1,
        "99999999999999999999999999999999999",
        "OUT:HASH",
        { from: ASGARD, value: _1 },
      );

      expect(BN2Str(await web3.eth.getBalance(ASGARD))).to.equal(
        "5620999333547258975226",
      ); // Less 1 ETH
      expect(BN2Str(await web3.eth.getBalance(SUSHIROUTER.address))).to.equal(
        _10,
      ); // No change
      expect(BN2Str(await web3.eth.getBalance(USER1))).to.equal(
        "8980971387836268055422",
      ); // +1 ETH (from failed swap)
      expect(BN2Str(await WETH.balanceOf(SUSHIROUTER.address))).to.equal(
        "9000000000000000000",
      ); // No change
      expect(BN2Str(await USDT.balanceOf(SUSHIROUTER.address))).to.equal(
        "20000000000000000000",
      ); // No change
      expect(BN2Str(await USDT.balanceOf(USER1))).to.equal(
        "999980000000000000000000",
      ); // No change
    });

    it("Should Fail Swap Out and ETH using Aggregator", async function () {
      expect(BN2Str(await USDT.balanceOf(SUSHIROUTER.address))).to.equal(
        "20000000000000000000",
      ); // Start bal
      expect(BN2Str(await web3.eth.getBalance(ASGARD))).to.equal(
        "5620999333547258975226",
      ); // Start bal

      // Asgard transferOutAndCall() 1 ETH
      // Send 1 ETH to router, forward to agg, forward to Sushi
      // Fail due price check. Try send ETH to user. Fails due to reverting contract. Send ETH back to Ygg.
      await ROUTER.transferOutAndCall(
        AGG.address,
        USDT.address,
        REVERTINGCONTRACT.address,
        "99999999999999999999999999999999999",
        "OUT:HASH",
        { from: ASGARD, value: _1 },
      );

      expect(
        BN2Str(await web3.eth.getBalance(REVERTINGCONTRACT.address)),
      ).to.equal("0"); // Zero
      expect(BN2Str(await web3.eth.getBalance(ASGARD))).to.equal(
        "5620999251766111414896",
      ); // Less 1 ETH
      expect(BN2Str(await web3.eth.getBalance(SUSHIROUTER.address))).to.equal(
        _10,
      ); // No change
      expect(BN2Str(await WETH.balanceOf(SUSHIROUTER.address))).to.equal(
        "9000000000000000000",
      ); // No change
      expect(BN2Str(await USDT.balanceOf(SUSHIROUTER.address))).to.equal(
        "20000000000000000000",
      ); // No change
    });
  });
});
