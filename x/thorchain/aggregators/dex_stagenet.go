//go:build stagenet
// +build stagenet

package aggregators

import (
	"github.com/blang/semver"
)

// If the contract whitelist is not (as in stagenet),
// use a default max gas and fall through to the suffix
// that is passed in. This should help dex agg contract devs test
// their work without having to run a mocknet or stagenet.
func DexAggregators(version semver.Version) []Aggregator {
	if version.GTE(semver.MustParse("0.1.0")) {
		return []Aggregator{}
	}
	return nil
}
