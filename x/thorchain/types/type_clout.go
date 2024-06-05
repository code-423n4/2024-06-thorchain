package types

import (
	"fmt"

	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
)

func NewSwapperClout(addr common.Address) SwapperClout {
	return SwapperClout{
		Address:   addr,
		Score:     cosmos.ZeroUint(),
		Reclaimed: cosmos.ZeroUint(),
		Spent:     cosmos.ZeroUint(),
	}
}

func (c SwapperClout) Valid() error {
	if ok := c.Address.IsEmpty(); !ok {
		return fmt.Errorf("invalid swapper clout address: %s", c.Address)
	}
	return nil
}

// calculate the available clout to spend
func (c SwapperClout) Available() cosmos.Uint {
	spent := common.SafeSub(c.Spent, c.Reclaimed)
	return common.SafeSub(c.Score, spent)
}

// calculate the available clout to reclaim
func (c SwapperClout) Claimable() cosmos.Uint {
	return common.SafeSub(c.Spent, c.Reclaimed)
}

func (c *SwapperClout) Reclaim(value cosmos.Uint) {
	c.Reclaimed = c.Reclaimed.Add(value)
	// reclaim should never exceed spent
	if c.Reclaimed.GT(c.Spent) {
		c.Reclaimed = c.Spent
	}
}

// if last spent occurred more than the limit, then reset the reclaim to equal
// spent. This means that clout is restored to 100% clout available after
// "limit" blocks have occurred without a swap. This is helpful if there is an
// accounting bug, it is corrected automatically after "limit" (typically 1
// hour)
func (c *SwapperClout) Restore(height, limit int64) {
	if c.LastSpentHeight+limit < height {
		c.Reclaimed = c.Spent
	}
}
