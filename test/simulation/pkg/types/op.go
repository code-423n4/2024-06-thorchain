package types

////////////////////////////////////////////////////////////////////////////////////////
// Ops
////////////////////////////////////////////////////////////////////////////////////////

// OpConfig is the configuration passed to each operation during execution.
type OpConfig struct {
	// AdminUser is the client for the mimir admin account.
	AdminUser *User

	// NodeUsers is a slice clients for simulation validator keys.
	NodeUsers []*User

	// Users is a slice of clients for simulation user keys.
	Users []*User
}

// OpResult is the result of an operation.
type OpResult struct {
	// Continue indicates that actor should continue to the next operation.
	Continue bool

	// Finish indicates that the actor should stop executing and return the error.
	Finish bool

	// Error is the error returned by the operation.
	Error error
}

// Op is an operation that can be executed by an actor.
type Op func(config *OpConfig) OpResult
