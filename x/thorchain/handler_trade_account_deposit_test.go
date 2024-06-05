package thorchain

import (
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	. "gopkg.in/check.v1"
)

type HandlerTradeAccountDeposit struct{}

func (HandlerTradeAccountDeposit) TestTradeAccountDeposit_Run(c *C) {
	ctx, mgr := setupManagerForTest(c)
	h := NewTradeAccountDepositHandler(mgr)
	asset := common.BTCAsset.GetTradeAsset()
	addr := GetRandomBech32Addr()
	dummyTx := common.Tx{ID: "test"}

	msg := NewMsgTradeAccountDeposit(asset, cosmos.NewUint(350), addr, addr, dummyTx)

	_, err := h.Run(ctx, msg)
	c.Assert(err, IsNil)

	bal := mgr.TradeAccountManager().BalanceOf(ctx, asset, addr)
	c.Check(bal.String(), Equals, "350")
}
