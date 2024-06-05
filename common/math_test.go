package common

import (
	"errors"

	"gitlab.com/thorchain/thornode/common/cosmos"
	. "gopkg.in/check.v1"
)

type MathSuite struct{}

var _ = Suite(&MathSuite{})

func (s *MathSuite) TestMax(c *C) {
	c.Assert(Max(1, 2), Equals, 2)
	c.Assert(Max(2, 1), Equals, 2)
	c.Assert(Max(1, 1), Equals, 1)

	c.Assert(Max(int64(1), int64(2)), Equals, int64(2))
	c.Assert(Max(int64(2), int64(1)), Equals, int64(2))
	c.Assert(Max(int64(1), int64(1)), Equals, int64(1))

	c.Assert(Max(uint(1), uint(2)), Equals, uint(2))
	c.Assert(Max(uint(2), uint(1)), Equals, uint(2))
	c.Assert(Max(uint(1), uint(1)), Equals, uint(1))
}

func (s *MathSuite) TestMin(c *C) {
	c.Assert(Min(1, 2), Equals, 1)
	c.Assert(Min(2, 1), Equals, 1)
	c.Assert(Min(1, 1), Equals, 1)

	c.Assert(Min(int64(1), int64(2)), Equals, int64(1))
	c.Assert(Min(int64(2), int64(1)), Equals, int64(1))
	c.Assert(Min(int64(1), int64(1)), Equals, int64(1))

	c.Assert(Min(uint(1), uint(2)), Equals, uint(1))
	c.Assert(Min(uint(2), uint(1)), Equals, uint(1))
	c.Assert(Min(uint(1), uint(1)), Equals, uint(1))
}

func (s *MathSuite) TestAbs(c *C) {
	c.Assert(Abs(1), Equals, 1)
	c.Assert(Abs(-1), Equals, 1)
	c.Assert(Abs(0), Equals, 0)

	c.Assert(Abs(int64(1)), Equals, int64(1))
	c.Assert(Abs(int64(-1)), Equals, int64(1))
	c.Assert(Abs(int64(0)), Equals, int64(0))
}

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
