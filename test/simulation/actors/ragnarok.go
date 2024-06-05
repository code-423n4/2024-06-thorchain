package actors

import (
	"fmt"
	"time"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/test/simulation/pkg/thornode"
	. "gitlab.com/thorchain/thornode/test/simulation/pkg/types"
	"gitlab.com/thorchain/thornode/x/thorchain/types"
)

////////////////////////////////////////////////////////////////////////////////////////
// RagnarokPoolActor
////////////////////////////////////////////////////////////////////////////////////////

type RagnarokPoolActor struct {
	Actor

	asset common.Asset
}

func NewRagnarokPoolActor(asset common.Asset) *Actor {
	a := &RagnarokPoolActor{
		Actor: *NewActor(fmt.Sprintf("Ragnarok-%s", asset)),
		asset: asset,
	}
	a.Timeout = 5 * time.Minute

	// TODO: get all LPs for the asset and store in state

	// send ragnarok mimir from admin
	a.Ops = append(a.Ops, a.sendMimir)

	// TODO: verify l1 balances
	// TODO: verify rune balances

	// verify pool removal
	a.Ops = append(a.Ops, a.verifyPoolRemoval)

	return &a.Actor
}

////////////////////////////////////////////////////////////////////////////////////////
// Ops
////////////////////////////////////////////////////////////////////////////////////////

func (a *RagnarokPoolActor) sendMimir(config *OpConfig) OpResult {
	accAddr, err := config.AdminUser.PubKey().GetThorAddress()
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to get thor address")
		return OpResult{
			Continue: false,
		}
	}
	mimir := types.NewMsgMimir(fmt.Sprintf("RAGNAROK-%s", a.asset.MimirString()), 1, accAddr)
	txid, err := config.AdminUser.Thorchain.Broadcast(mimir)
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to broadcast mimir")
		return OpResult{
			Continue: false,
		}
	}
	a.Log().Info().Str("txid", txid.String()).Msg("broadcasted mimir")
	return OpResult{
		Continue: true,
	}
}

func (a *RagnarokPoolActor) verifyPoolRemoval(config *OpConfig) OpResult {
	// fetch pools
	pools, err := thornode.GetPools()
	if err != nil {
		a.Log().Error().Err(err).Msg("failed to get pools")
		return OpResult{
			Continue: false,
		}
	}

	// verify pool removal
	found := false
	for _, pool := range pools {
		if pool.Asset == a.asset.String() {
			found = true
			break
		}
	}

	if found {
		return OpResult{
			Continue: false,
		}
	}

	a.Log().Info().Msg("pool removed")
	return OpResult{
		Finish: true,
	}
}
