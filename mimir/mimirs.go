package mimir

import (
	"gitlab.com/thorchain/thornode/constants"
)

func getRef(refs []string) (reference string) {
	for _, ref := range refs {
		if len(ref) > 0 {
			reference = ref
		}
	}
	return
}

func NewAffiliateFeeBasisPointsMax(refs ...string) Mimir {
	id := AffiliateFeeBasisPointsMax
	return &mimir{
		id:           id,
		name:         id.String(),
		defaultValue: int64(constants.MaxBasisPts),
		mimirType:    EconomicMimir,
		reference:    getRef(refs),
		tags:         []string{"economic", "affiliate fee"},
		description:  "Maximum fee to allow affiliates to set",
		legacyMimirKey: func(_ string) string {
			return "MaxAffiliateFeeBasisPoints"
		},
	}
}

func NewBondPause(refs ...string) Mimir {
	id := BondPause
	return &mimir{
		id:           id,
		name:         id.String(),
		defaultValue: 0,
		reference:    getRef(refs),
		mimirType:    OperationalMimir,
		tags:         []string{"operational", "bond"},
		description:  "Pauses bonding (unbonding is still allowed)",
		legacyMimirKey: func(_ string) string {
			return "PauseBond"
		},
	}
}

func NewConfBasisPointValue(refs ...string) Mimir {
	id := ConfMultiplierBasisPoints
	return &mimir{
		id:           id,
		name:         id.String(),
		defaultValue: int64(constants.MaxBasisPts),
		reference:    getRef(refs),
		mimirType:    EconomicMimir,
		tags:         []string{"economic", "chain-client"},
		description:  "adjusts confirmation multiplier for chain client",
		legacyMimirKey: func(_ string) string {
			return ""
		},
	}
}

func NewMaxConfValue(refs ...string) Mimir {
	id := MaxConfirmations
	return &mimir{
		id:           id,
		name:         id.String(),
		defaultValue: 0,
		reference:    getRef(refs),
		mimirType:    EconomicMimir,
		tags:         []string{"economic", "chain-client"},
		description:  "max confirmations for chain client",
		legacyMimirKey: func(_ string) string {
			return ""
		},
	}
}

func NewSwapperCloutLimit(refs ...string) Mimir {
	id := CloutSwapperLimit
	return &mimir{
		id:           id,
		name:         id.String(),
		defaultValue: 0,
		mimirType:    EconomicMimir,
		reference:    getRef(refs),
		tags:         []string{"economic", "clout"},
		description:  "Maximum clout applicable to an outbound txn",
		legacyMimirKey: func(_ string) string {
			return "CloutLimit"
		},
	}
}

func NewSwapperCloutReset(refs ...string) Mimir {
	id := CloutSwapperReset
	return &mimir{
		id:           id,
		name:         id.String(),
		defaultValue: 720,
		mimirType:    EconomicMimir,
		reference:    getRef(refs),
		tags:         []string{"economic", "clout"},
		description:  "Amount of blocks before pending clout spent is reset",
		legacyMimirKey: func(_ string) string {
			return "CloutReset"
		},
	}
}

func NewSwapSlipBasisPointsMin(refs ...string) Mimir {
	id := SwapSlipBasisPointsMin
	return &mimir{
		id:           id,
		name:         id.String(),
		defaultValue: 0,
		reference:    getRef(refs),
		mimirType:    EconomicMimir,
		tags:         []string{"economic", "swap min slip"},
		description:  "Min slip on swap",
		legacyMimirKey: func(_ string) string {
			return ""
		},
	}
}

func NewTradeAccountsEnabled(refs ...string) Mimir {
	id := TradeAccountEnabled
	return &mimir{
		id:          id,
		name:        id.String(),
		mimirType:   OperationalMimir,
		reference:   getRef(refs),
		tags:        []string{"operational", "trade"},
		description: "Enable or disable trade accounts",
		legacyMimirKey: func(_ string) string {
			return "TradeAccountsEnabled"
		},
	}
}
