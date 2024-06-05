package common

import (
	. "gopkg.in/check.v1"
)

type RoundSignificantFiguresSuite struct{}

var _ = Suite(&RoundSignificantFiguresSuite{})

func (s *RoundSignificantFiguresSuite) TestRoundSignificantFigures(c *C) {
	testCases := []struct {
		number            uint64
		significantDigits int64
		expectedResult    uint64
	}{
		{123456, 3, 123000},
		{123456, 100, 123456},
		{9876543210, 5, 9876500000},
		{9876543210, 10, 9876543210},
		{0, 2, 0},
		{1, 1, 1},
		{9, 2, 9},
		{999, 1, 900},
		{999, 3, 999},
		{1000, 3, 1000},
		{1000, 2, 1000},
		{1000, 1, 1000},
		{1000, 4, 1000},
		{9999, 3, 9990},
		{1000000000, 1, 1000000000},
		{9999999999, 1, 9000000000},
	}

	for _, testCase := range testCases {
		result := RoundSignificantFigures(testCase.number, testCase.significantDigits)
		if result != testCase.expectedResult {
			c.Logf("case: %+v, result: %d", testCase, result)
			c.Fail()
		}
	}
}
