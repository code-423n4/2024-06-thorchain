package thorchain

import (
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	. "gopkg.in/check.v1"
)

type HandlerTradeAccountWithdrawal struct{}

func (HandlerTradeAccountWithdrawal) TestTradeAccountWithdrawal_Run(c *C) {
	ctx, mgr := setupManagerForTest(c)
	h := NewTradeAccountWithdrawalHandler(mgr)
	asset := common.BTCAsset.GetTradeAsset()
	addr := GetRandomBech32Addr()
	bc1Addr := GetRandomBTCAddress()
	dummyTx := common.Tx{ID: "test"}

	_, err := mgr.TradeAccountManager().Deposit(ctx, asset, cosmos.NewUint(500), addr, common.NoAddress, dummyTx.ID)
	c.Assert(err, IsNil)

	msg := NewMsgTradeAccountWithdrawal(asset, cosmos.NewUint(350), bc1Addr, addr, dummyTx)

	_, err = h.Run(ctx, msg)
	c.Assert(err, IsNil)

	bal := mgr.TradeAccountManager().BalanceOf(ctx, asset, addr)
	c.Check(bal.String(), Equals, "150")

	items, err := mgr.TxOutStore().GetOutboundItems(ctx)
	c.Assert(err, IsNil)
	c.Assert(items, HasLen, 1)
	c.Check(items[0].Coin.String(), Equals, "350 BTC~BTC")
	c.Check(items[0].ToAddress.String(), Equals, bc1Addr.String())
}
