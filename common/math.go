package common

import (
	"errors"
	"sort"

	"gitlab.com/thorchain/thornode/common/cosmos"
)

func GetMedianUint(vals []cosmos.Uint) cosmos.Uint {
	if len(vals) == 0 {
		return cosmos.ZeroUint()
	}

	sort.SliceStable(vals, func(i, j int) bool {
		return vals[i].LT(vals[j])
	})

	// calculate median
	var median cosmos.Uint
	if len(vals)%2 > 0 {
		// odd number of figures in our slice. Take the middle figure. Since
		// slices start with an index of zero, just need to length divide by two.
		medianSpot := len(vals) / 2
		median = vals[medianSpot]
	} else {
		// even number of figures in our slice. Average the middle two figures.
		pt1 := vals[len(vals)/2-1]
		pt2 := vals[len(vals)/2]
		median = pt1.Add(pt2).QuoUint64(2)
	}
	return median
}

func GetMedianInt64(vals []int64) int64 {
	switch len(vals) {
	case 0:
		return 0
	case 1:
		return vals[0]
	}

	sort.SliceStable(vals, func(i, j int) bool {
		return vals[i] < vals[j]
	})

	// calculate median
	var median int64
	if len(vals)%2 > 0 {
		// odd number of figures in our slice. Take the middle figure. Since
		// slices start with an index of zero, just need to length divide by two.
		medianSpot := len(vals) / 2
		median = vals[medianSpot]
	} else {
		// even number of figures in our slice. Average the middle two figures.
		pt1 := vals[len(vals)/2-1]
		pt2 := vals[len(vals)/2]
		median = (pt1 + pt2) / 2
	}
	return median
}

// WeightedMean calculates the weighted mean of a set of values and their weights.
func WeightedMean(vals, weights []cosmos.Uint) (cosmos.Uint, error) {
	totalWeight := cosmos.Sum(weights)

	// if total weight is zero, return an error
	if totalWeight.IsZero() {
		return cosmos.ZeroUint(), errors.New("total weight is zero")
	}

	// assert that the number of values and weights are the same
	if len(vals) != len(weights) {
		panic("number of values and weights do not match")
	}

	// calculate the weight in basis points for each anchor
	weightedTotal := cosmos.ZeroUint()
	for i, val := range vals {
		weightedTotal = weightedTotal.Add(val.Mul(weights[i]))
	}

	return weightedTotal.Quo(totalWeight), nil
}
