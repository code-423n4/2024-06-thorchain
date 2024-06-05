//go:build regtest
// +build regtest

package thorchain

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"gitlab.com/thorchain/thornode/common/cosmos"
	q "gitlab.com/thorchain/thornode/x/thorchain/query"
)

func init() {
	initManager = func(mgr *Mgrs, ctx cosmos.Context) {
		_ = mgr.BeginBlock(ctx)
	}

	optionalQuery = func(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
		switch path[0] {
		case q.QueryExport.Key:
			return queryExport(ctx, path[1:], req, mgr)
		default:
			return nil, cosmos.ErrUnknownRequest(
				fmt.Sprintf("unknown thorchain query endpoint: %s", path[0]),
			)
		}
	}
}

func queryExport(ctx cosmos.Context, path []string, req abci.RequestQuery, mgr *Mgrs) ([]byte, error) {
	return jsonify(ctx, ExportGenesis(ctx, mgr.Keeper()))
}
