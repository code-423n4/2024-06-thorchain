package common

import (
	"math/big"

	cosmos "gitlab.com/thorchain/thornode/common/cosmos"
	. "gopkg.in/check.v1"
)

type GasSuite struct{}

var _ = Suite(&GasSuite{})

func (s *GasSuite) TestETHGasFee(c *C) {
	gas := GetEVMGasFee(ETHChain, big.NewInt(20), 4)
	amt := gas[0].Amount
	c.Check(
		amt.Equal(cosmos.NewUint(425440)),
		Equals,
		true,
		Commentf("%d", amt.Uint64()),
	)
	gas = MakeEVMGas(ETHChain, big.NewInt(20), 10000000000) // 10 GWEI
	amt = gas[0].Amount
	c.Check(
		amt.Equal(cosmos.NewUint(20)),
		Equals,
		true,
		Commentf("%d", amt.Uint64()),
	)
	// ETH TxID b89d5eb71765b42117bb1fa30d3a22f6d2bfdba9214da60d26f028bd94bcdb0c example
	gas = MakeEVMGas(ETHChain, big.NewInt(18000803458), 63707)
	// 18000803458 Wei gasPrice, 63,707 gas
	amt = gas[0].Amount
	c.Check(
		amt.Equal(cosmos.NewUint(114678)),
		// Should be rounded up to 114678, not down to 114677,
		// to increase rather than decrease solvency.
		Equals,
		true,
		Commentf("%d", amt.Uint64()),
	)
}

func (s *GasSuite) TestIsEmpty(c *C) {
	gas1 := Gas{
		{Asset: BNBAsset, Amount: cosmos.NewUint(11 * One)},
	}
	c.Check(gas1.IsEmpty(), Equals, false)
	c.Check(Gas{}.IsEmpty(), Equals, true)
}

func (s *GasSuite) TestCombineGas(c *C) {
	gas1 := Gas{
		{Asset: BNBAsset, Amount: cosmos.NewUint(11 * One)},
	}
	gas2 := Gas{
		{Asset: BNBAsset, Amount: cosmos.NewUint(14 * One)},
		{Asset: BTCAsset, Amount: cosmos.NewUint(20 * One)},
	}
	gas := gas1.Add(gas2...)

	// Confirm the slice lengths.
	c.Assert(gas1, HasLen, 1)
	c.Assert(gas2, HasLen, 2)
	c.Assert(gas, HasLen, 2)

	// Check gas contents.
	c.Check(gas[0].Asset.Equals(BNBAsset), Equals, true)
	c.Check(gas[0].Amount.Equal(cosmos.NewUint(25*One)), Equals, true, Commentf("%d", gas[0].Amount.Uint64()))
	c.Check(gas[1].Asset.Equals(BTCAsset), Equals, true)
	c.Check(gas[1].Amount.Equal(cosmos.NewUint(20*One)), Equals, true)

	// Check whether there is any influence on the Amounts of gas1 or gas2 by the combining.
	c.Check(gas1[0].Amount.Equal(cosmos.NewUint(11*One)), Equals, true, Commentf("%d", gas1[0].Amount.Uint64()))
	c.Check(gas2[0].Amount.Equal(cosmos.NewUint(14*One)), Equals, true, Commentf("%d", gas2[0].Amount.Uint64()))
	c.Check(gas2[1].Amount.Equal(cosmos.NewUint(20*One)), Equals, true, Commentf("%d", gas2[1].Amount.Uint64()))
	// gas1 and gas2 are unchanged by the combination, as they should be.

	// Check whether changes to gas can affect gas1 or gas2.
	gas[0].Amount = gas[0].Amount.MulUint64(2)
	gas[1].Amount = gas[1].Amount.MulUint64(2)
	c.Check(gas1[0].Amount.Equal(cosmos.NewUint(11*One)), Equals, true, Commentf("%d", gas1[0].Amount.Uint64()))
	c.Check(gas2[0].Amount.Equal(cosmos.NewUint(14*One)), Equals, true, Commentf("%d", gas2[0].Amount.Uint64()))
	c.Check(gas2[1].Amount.Equal(cosmos.NewUint(20*One)), Equals, true, Commentf("%d", gas2[1].Amount.Uint64()))
	// gas1 and gas2 are unaffected, as they should be.

	// Check whether changes to gas1 or gas2 can affect gas.
	c.Check(gas[0].Amount.Equal(cosmos.NewUint(50*One)), Equals, true, Commentf("%d", gas[0].Amount.Uint64()))
	c.Check(gas[1].Amount.Equal(cosmos.NewUint(40*One)), Equals, true, Commentf("%d", gas[0].Amount.Uint64()))
	gas1[0].Amount = gas1[0].Amount.MulUint64(2)
	gas2[0].Amount = gas2[0].Amount.MulUint64(2)
	gas2[1].Amount = gas2[1].Amount.MulUint64(2)
	c.Check(gas[0].Amount.Equal(cosmos.NewUint(50*One)), Equals, true, Commentf("%d", gas[0].Amount.Uint64()))
	c.Check(gas[1].Amount.Equal(cosmos.NewUint(40*One)), Equals, true, Commentf("%d", gas[0].Amount.Uint64()))
	// gas is unaffected, as it should be.
}

func (s *GasSuite) TestCalcGasPrice(c *C) {
	gasInfo := []cosmos.Uint{cosmos.NewUint(37500), cosmos.NewUint(30000)}
	tx := Tx{
		Coins: Coins{
			NewCoin(BNBAsset, cosmos.NewUint(80808080)),
		},
	}

	gas := CalcBinanceGasPrice(tx, BNBAsset, gasInfo)
	c.Check(gas.Equals(Gas{NewCoin(BNBAsset, cosmos.NewUint(37500))}), Equals, true)

	tx = Tx{
		Coins: Coins{
			NewCoin(BNBAsset, cosmos.NewUint(80808080)),
			NewCoin(BNBAsset, cosmos.NewUint(80808080)),
		},
	}

	gas = CalcBinanceGasPrice(tx, BNBAsset, gasInfo)
	c.Check(gas.Equals(Gas{NewCoin(BNBAsset, cosmos.NewUint(60000))}), Equals, true)
}
