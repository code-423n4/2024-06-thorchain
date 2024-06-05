package mimir

import (
	"fmt"
	"strings"

	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
	"gitlab.com/thorchain/thornode/x/thorchain/keeper"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

func (m *mimir) fetchValueV124(ctx cosmos.Context, keeper keeper.Keeper) (value int64) {
	var (
		err    error
		active types.NodeAccounts
		key    string
	)
	active, err = keeper.ListActiveValidators(ctx)
	if err != nil {
		ctx.Logger().Error("failed to get active validator set", "error", err)
	}

	key = fmt.Sprintf("%d-%s", m.id, strings.ToUpper(m.reference))
	var mimirs types.NodeMimirs
	mimirs, err = keeper.GetNodeMimirsV2(ctx, key)
	if err != nil {
		ctx.Logger().Error("failed to get node mimir v2", "error", err)
	}
	value = int64(-1)
	switch m.Type() {
	case EconomicMimir:
		value = mimirs.ValueOfEconomic(key, active.GetNodeAddresses())
		if value < 0 {
			// no value, fallback to last economic value (if present)
			value, err = keeper.GetMimirV2(ctx, key)
			if err != nil {
				ctx.Logger().Error("failed to get mimir v2", "error", err)
			}
			if value >= 0 {
				return value
			}
		} else {
			// value reached, save to persist it beyond losing 2/3rds
			keeper.SetMimirV2(ctx, key, value)
		}
	case OperationalMimir:
		value = mimirs.ValueOfOperational(key, constants.MinMimirV2Vote, active.GetNodeAddresses())
	}
	return
}
