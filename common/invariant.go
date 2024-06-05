package common

import (
	"gitlab.com/thorchain/thornode/common/cosmos"
)

// borrowed and modified from cosmos types
// https://github.com/cosmos/cosmos-sdk/blob/v0.45.1/types/invariant.go
// https://github.com/cosmos/cosmos-sdk/blob/v0.45.1/x/crisis/types/route.go
//
// differences are:
// - array of messages for readability purposes
// - no module in routes, we use a single thorchain module

// An Invariant is a function which tests a particular invariant.
// The invariant returns a descriptive message about what happened
// and a boolean indicating whether the invariant has been broken.
type Invariant func(ctx cosmos.Context) (msg []string, broken bool)

// invariant route
type InvariantRoute struct {
	Route     string
	Invariant Invariant
}

// NewInvariantRoute - create an InvariantRoute object
func NewInvariantRoute(route string, invariant Invariant) InvariantRoute {
	return InvariantRoute{
		Route:     route,
		Invariant: invariant,
	}
}
