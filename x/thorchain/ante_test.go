package thorchain

import (
	"github.com/stretchr/testify/suite"

	. "gopkg.in/check.v1"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"gitlab.com/thorchain/thornode/common/cosmos"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

type AnteTestSuite struct {
	suite.Suite
}

var _ = Suite(&AnteTestSuite{})

func (s *AnteTestSuite) TestRejectMutlipleDepositMsgs(c *C) {
	ctx, k := setupKeeperForTest(c)

	ad := AnteDecorator{
		keeper: k,
	}

	msgs := []cosmos.Msg{
		&types.MsgSend{},
		&types.MsgBan{},
	}

	// no deposit msgs is ok
	err := ad.rejectMultipleDepositMsgs(ctx, msgs)
	c.Assert(err, IsNil)

	// one deposit msgs is ok
	msgs = append(msgs, &types.MsgDeposit{})
	err = ad.rejectMultipleDepositMsgs(ctx, msgs)
	c.Assert(err, IsNil)

	// two deposit msgs is not ok
	msgs = append(msgs, &types.MsgDeposit{})
	err = ad.rejectMultipleDepositMsgs(ctx, msgs)
	c.Assert(err, NotNil)
}

func (s *AnteTestSuite) TestAnteHandleMessage(c *C) {
	ctx, k := setupKeeperForTest(c)
	version := GetCurrentVersion()

	ad := AnteDecorator{
		keeper: k,
	}

	fromAddr := GetRandomBech32Addr()
	toAddr := GetRandomBech32Addr()

	// fund an addr so it can pass the fee deduction ante
	funds, err := common.NewCoin(common.RuneNative, cosmos.NewUint(200*common.One)).Native()
	c.Assert(err, IsNil)
	err = k.AddCoins(ctx, fromAddr, cosmos.NewCoins(funds))
	c.Assert(err, IsNil)
	coin, err := common.NewCoin(common.RuneNative, cosmos.NewUint(1*common.One)).Native()
	c.Assert(err, IsNil)

	goodMsg := types.MsgSend{
		FromAddress: fromAddr,
		ToAddress:   toAddr,
		Amount:      cosmos.NewCoins(coin),
	}
	err = ad.anteHandleMessage(ctx, version, &goodMsg)
	c.Assert(err, IsNil)

	// non-thorchain msgs should be rejected
	badMsg := banktypes.MsgSend{}
	err = ad.anteHandleMessage(ctx, version, &badMsg)
	c.Assert(err, NotNil)
}
