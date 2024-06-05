package common

import (
	"errors"

	"gitlab.com/thorchain/thornode/common/cosmos"
	. "gopkg.in/check.v1"
)

type MathSuite struct{}

var _ = Suite(&MathSuite{})

func (s *MathSuite) TestWeightedMean(c *C) {
	vals := []cosmos.Uint{cosmos.NewUint(10), cosmos.NewUint(20), cosmos.NewUint(30)}
	weights := []cosmos.Uint{cosmos.NewUint(1), cosmos.NewUint(2), cosmos.NewUint(3)}
	expectedMean := cosmos.NewUint(140 / 6) // (10*1 + 20*2 + 30*3) / (1+2+3)
	mean, err := WeightedMean(vals, weights)
	c.Assert(err, IsNil)
	c.Assert(mean.String(), Equals, expectedMean.String())
}

func (s *MathSuite) TestWeightedMeanErrors(c *C) {
	// mismatched values and weights
	valsMismatch := []cosmos.Uint{cosmos.NewUint(10), cosmos.NewUint(20)}
	weightsMismatch := []cosmos.Uint{cosmos.NewUint(1), cosmos.NewUint(2), cosmos.NewUint(3)}
	testFn := func() { _, _ = WeightedMean(valsMismatch, weightsMismatch) }
	c.Assert(testFn, PanicMatches, "number of values and weights do not match")

	// zero total weight
	valsZero := []cosmos.Uint{cosmos.NewUint(10), cosmos.NewUint(20), cosmos.NewUint(30)}
	weightsZero := []cosmos.Uint{cosmos.ZeroUint(), cosmos.ZeroUint(), cosmos.ZeroUint()}
	_, errZero := WeightedMean(valsZero, weightsZero)
	c.Assert(errZero, DeepEquals, errors.New("total weight is zero"))
}
