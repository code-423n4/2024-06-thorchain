package gaia

import (
	"context"
	"os"

	"github.com/cometbft/cometbft/libs/json"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
)

type TendermintRPC interface {
	Block(ctx context.Context, height *int64) (*ctypes.ResultBlock, error)
	BlockResults(ctx context.Context, height *int64) (*ctypes.ResultBlockResults, error)
}

type mockTendermintRPC struct{}

func (m *mockTendermintRPC) Block(ctx context.Context, height *int64) (*ctypes.ResultBlock, error) {
	out := new(ctypes.ResultBlock)

	path := "./test-data/latest_block.json"
	if height != nil {
		path = "./test-data/block_by_height.json"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, out)

	return out, err
}

func (m *mockTendermintRPC) BlockResults(ctx context.Context, height *int64) (*ctypes.ResultBlockResults, error) {
	out := new(ctypes.ResultBlockResults)
	data, err := os.ReadFile("./test-data/tx_results_by_height.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, out)

	return out, err
}
