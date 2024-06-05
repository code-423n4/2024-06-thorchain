package keeperv1

import (
	"gitlab.com/thorchain/thornode/common"
	"gitlab.com/thorchain/thornode/common/cosmos"
	"gitlab.com/thorchain/thornode/constants"
)

func (k KVStore) GetOutboundTxFee(ctx cosmos.Context) cosmos.Uint {
	if k.usdFeesEnabled(ctx) {
		return k.DollarConfigInRune(ctx, constants.NativeOutboundFeeUSD)
	}
	fee := k.GetConfigInt64(ctx, constants.OutboundTransactionFee)
	return cosmos.NewUint(uint64(fee))
}

// GetOutboundFeeWithheldRune - record of RUNE collected by the Reserve for an Asset's outbound fees
func (k KVStore) GetOutboundFeeWithheldRune(ctx cosmos.Context, outAsset common.Asset) (cosmos.Uint, error) {
	var record uint64
	_, err := k.getUint64(ctx, k.GetKey(ctx, prefixOutboundFeeWithheldRune, outAsset.String()), &record)
	return cosmos.NewUint(record), err
}

// AddToOutboundFeeWithheldRune - add to record of RUNE collected by the Reserve for an Asset's outbound fees
func (k KVStore) AddToOutboundFeeWithheldRune(ctx cosmos.Context, outAsset common.Asset, withheld cosmos.Uint) error {
	outboundFeeWithheldRune, err := k.GetOutboundFeeWithheldRune(ctx, outAsset)
	if err != nil {
		return err
	}

	outboundFeeWithheldRune = outboundFeeWithheldRune.Add(withheld)
	k.setUint64(ctx, k.GetKey(ctx, prefixOutboundFeeWithheldRune, outAsset.String()), outboundFeeWithheldRune.Uint64())
	return nil
}

// GetOutboundFeeWithheldRuneIterator to iterate through all Assets' OutboundFeeWithheldRune
// (e.g. for hard-fork GenesisState export)
func (k KVStore) GetOutboundFeeWithheldRuneIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixOutboundFeeWithheldRune)
}

// GetOutboundFeeSpentRune - record of RUNE spent by the Reserve for an Asset's outbounds' gas costs
func (k KVStore) GetOutboundFeeSpentRune(ctx cosmos.Context, outAsset common.Asset) (cosmos.Uint, error) {
	var record uint64
	_, err := k.getUint64(ctx, k.GetKey(ctx, prefixOutboundFeeSpentRune, outAsset.String()), &record)
	return cosmos.NewUint(record), err
}

// AddToOutboundFeeSpentRune - add to record of RUNE spent by the Reserve for an Asset's outbounds' gas costs
func (k KVStore) AddToOutboundFeeSpentRune(ctx cosmos.Context, outAsset common.Asset, spent cosmos.Uint) error {
	outboundFeeSpentRune, err := k.GetOutboundFeeSpentRune(ctx, outAsset)
	if err != nil {
		return err
	}

	outboundFeeSpentRune = outboundFeeSpentRune.Add(spent)
	k.setUint64(ctx, k.GetKey(ctx, prefixOutboundFeeSpentRune, outAsset.String()), outboundFeeSpentRune.Uint64())
	return nil
}

// GetOutboundFeeSpentRuneIterator to iterate through all Assets' OutboundFeeSpentRune
// (e.g. for hard-fork GenesisState export)
func (k KVStore) GetOutboundFeeSpentRuneIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixOutboundFeeSpentRune)
}
