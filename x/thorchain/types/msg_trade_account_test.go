package types

import (
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	. "gopkg.in/check.v1"
)

type MsgTradeAccountSuite struct{}

var _ = Suite(&MsgTradeAccountSuite{})

func (MsgTradeAccountSuite) TestDeposit(c *C) {
	asset := common.BNBAsset
	amt := cosmos.NewUint(100)
	signer := GetRandomBech32Addr()
	dummyTx := common.Tx{ID: "test"}

	m := NewMsgTradeAccountDeposit(asset, amt, signer, signer, dummyTx)
	EnsureMsgBasicCorrect(m, c)
	c.Check(m.Type(), Equals, "set_trade_account_deposit")

	m = NewMsgTradeAccountDeposit(common.EmptyAsset, amt, signer, signer, dummyTx)
	c.Check(m.ValidateBasic(), NotNil)

	m = NewMsgTradeAccountDeposit(common.RuneAsset(), amt, signer, signer, dummyTx)
	c.Check(m.ValidateBasic(), NotNil)

	m = NewMsgTradeAccountDeposit(asset, cosmos.ZeroUint(), signer, signer, dummyTx)
	c.Check(m.ValidateBasic(), NotNil)
}

func (MsgTradeAccountSuite) TestWithdrawal(c *C) {
	asset := common.BNBAsset.GetTradeAsset()
	amt := cosmos.NewUint(100)
	bnbAddr := GetRandomBNBAddress()
	signer := GetRandomBech32Addr()
	dummyTx := common.Tx{ID: "test"}

	m := NewMsgTradeAccountWithdrawal(asset, amt, bnbAddr, signer, dummyTx)
	EnsureMsgBasicCorrect(m, c)
	c.Check(m.Type(), Equals, "set_trade_account_withdrawal")

	m = NewMsgTradeAccountWithdrawal(common.EmptyAsset, amt, bnbAddr, signer, dummyTx)
	c.Check(m.ValidateBasic(), NotNil)

	m = NewMsgTradeAccountWithdrawal(common.RuneAsset(), amt, bnbAddr, signer, dummyTx)
	c.Check(m.ValidateBasic(), NotNil)

	m = NewMsgTradeAccountWithdrawal(asset, cosmos.ZeroUint(), bnbAddr, signer, dummyTx)
	c.Check(m.ValidateBasic(), NotNil)

	m = NewMsgTradeAccountWithdrawal(asset, cosmos.ZeroUint(), GetRandomTHORAddress(), signer, dummyTx)
	c.Check(m.ValidateBasic(), NotNil)
}
