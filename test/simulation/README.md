# THORNode Simulation Testing Framework

While we have existing test frameworks (unit, regression, smoke), none are able to fully explore the boundaries of the system and the dynamic interactions between all components. Unit tests are narrowly focused with limited usefulness on the majority of logic which relies on complex interactions between components. Regression tests are one layer removed to perform discrete and deterministic tests against the Thornode state machine - they have proven useful, but remain insufficient to cover the full path from external L1s, through Bifrost, and into Thornode. Smoke tests were intended to cover this end to end testing, but they fail to capture the dynamic nature of the underlying chain states, making them brittle and difficult to extend. Smoke tests have additional barriers that have prevented ongoing extension to their coverage - namely the static and inter-dependent test definitions, as well as the system being built in Python which requires protobuf changes, mock additions, etcetera to add any new chains.

The goal of the simulation test framework is to provide a revised approach to the smoke test problem. Instead of attempting to replicate the state machine and compare it to a mocknet state, we define independent test sequences similar to regression tests. The simulation tests are not an extension or replacement of regression tests, but rather they are complementary - the delineation being simulation test operations are dynamic instead of discrete numerical checks (ex: get a swap quote for half the account BTC balance -> ETH, send the deposit, ensure the received amount is within some expected range based on the quote). The full test suite will be defined as an execution DAG of tests - operations in each test are executed sequentially, but there are a configurable number of tests run concurrently and the global ordering is intentionally random and interleaved.

The initial implementation will leverage a static actor DAG that runs to completion - like current smoke tests. Future work can extend these static tests with fuzzing to generate a random DAG of arbitrary length to further probe boundaries of the system over a longer duration.

## Running Simulation Tests

Simulation tests can be run with the following command from the repo root:

```bash
make test-simulation
```

This will internally build the simulation framework image, start mocknet, and run the tests from a Docker container with host mode networking to reach the local mocknet.

## Design

The simulation test framework will accept a test suite defined as a DAG of actors. Each actor defines operations that are executed sequentially, but may be randomly delayed and interleaved with operations from other concurrently running actors. Unlike regression tests where operations are structs, in the dynamic world of simulation tests every operation is a function which receives config and a mutable state, and returns an error, a bool whether to continue to the next operation in the test, and a bool to mark the test as finished.

Separate from the test actors, we allow definition of watchers. These watchers run asynchrously to check things like invariants and trigger test failure independent of any actor execution.

### Structure

```none
actors/   # one file per actor definition
cmd/      # the executable entry point for the simulation framework
pkg/      # internal packages
suites/   # suites of actors
watchers/ # one file per watcher
```

### Actors

```golang
// Actor is a node in the actor tree, which is processed as a DAG. Each actor has a set
// of operations to execute and a set of children. The actor will execute each operation
// in order until the operation returns Continue = false. If the operation returns
// Finish = true, the actor will stop executing and return the error. Once an actor has
// finished, the children of the actor will be executed. Actors may be appended to other
// actors, which results in the final descendants of the actor being parents of the
// appended actor - blocking the execution of the appended actor until all parents have
// completed, allowing for "fan in" of the execution DAG.
type Actor struct {
	// Name is the name of the actor.
	Name string

	// Timeout is the maximum amount of time the actor is allowed to execute.
	Timeout time.Duration

	// Interval is the amount of time to wait between operations.
	Interval time.Duration

	// Ops is the set of operations to execute.
	Ops []Op `json:"-"`

	// Children is the set of children to execute after the operations have completed.
	Children []*Actor
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
```

### Watchers

```golang
// Watcher wraps a function that will be executed on some interval in the background for
// the duration of the simulation. The function is passed the same OpConfig as the
// operations in the actor tree so it is able to access the clients and accounts used by
// the actors. If it returns an error or panics, the simulation will be aborted.
type Watcher struct {
	// Name is the name of the watcher.
	Name string

	// Interval is the interval at which the watcher will be executed.
	Interval time.Duration

	// Fn is the function to execute.
	Fn func(config *OpConfig) error
}
```

## Mocknet Master Accounts

The following accounts derive from the `master` mnemonic:

```bash
$ go run -tags mocknet tools/pubkey2address/pubkey2address.go -p tthorpub1addwnpepqwutw9cpacdkgnduh7e6cgd8ar7v5rgqkemxffuxdauzw3nlfq7sxtymlzs

BSC Address: 0xee4eaa642b992412f628ff4cec1c96cf2fd0ea4d
ETH Address: 0xee4eaa642b992412f628ff4cec1c96cf2fd0ea4d
BTC Address: bcrt1qf4l5dlqhaujgkxxqmug4stfvmvt58vx2h44c39
LTC Address: rltc1qf4l5dlqhaujgkxxqmug4stfvmvt58vx2fc03xm
BCH Address: qpxh73huzlhjfzcccr03zkpd9nd3wsasegmrreet72
DOGE Address: mnaioCtEGdw6bd6rWJ13Mbre1kN5rPa2Mo
THOR Address: tthor1f4l5dlqhaujgkxxqmug4stfvmvt58vx2tspx4g
GAIA Address: cosmos1f4l5dlqhaujgkxxqmug4stfvmvt58vx2fqfdej
AVAX Address: 0xee4eaa642b992412f628ff4cec1c96cf2fd0ea4d
```
