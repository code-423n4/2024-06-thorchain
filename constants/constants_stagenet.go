//go:build stagenet
// +build stagenet

package constants

func init() {
	int64Overrides = map[ConstantName]int64{
		ChurnInterval:               432000,
		MinRunePoolDepth:            1_00000000,
		MinimumBondInRune:           200_000_00000000,
		PoolCycle:                   720,
		EmissionCurve:               8,
		StopFundYggdrasil:           1, // TODO remove on hard fork
		YggFundLimit:                0, // TODO remove on hard fork
		NumberOfNewNodesPerChurn:    4,
		MintSynths:                  1,
		BurnSynths:                  1,
		KillSwitchStart:             1,
		KillSwitchDuration:          1,
		MaxRuneSupply:               500_000_000_00000000,
		FullImpLossProtectionBlocks: 0, // TODO remove on hard fork
	}
}
