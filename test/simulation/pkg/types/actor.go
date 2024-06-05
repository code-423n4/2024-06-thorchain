package types

import (
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

////////////////////////////////////////////////////////////////////////////////////////
// Actor
////////////////////////////////////////////////////////////////////////////////////////

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
	Children map[*Actor]bool `json:"-"`

	// -------------------- internal --------------------

	log      zerolog.Logger
	parents  map[*Actor]bool
	started  *atomic.Bool
	finished *atomic.Bool
}

// NewActor will create a new actor with the provided name.
func NewActor(name string) *Actor {
	return &Actor{
		Name:     name,
		Children: map[*Actor]bool{},
		parents:  map[*Actor]bool{},
	}
}

// Start will mark the actor as started.
func (a *Actor) Start() {
	a.started.Store(true)
}

// Started returns true if the actor has started executing.
func (a *Actor) Started() bool {
	return a.started.Load()
}

// Finished returns true if the actor has finished executing.
func (a *Actor) Finished() bool {
	return a.finished.Load()
}

// Parents returns the parents of the actor.
func (a *Actor) Parents() map[*Actor]bool {
	return a.parents
}

// WalkDepthFirst will walk the actor tree depth first and execute the provided function
// on each actor. If the function returns false, the walk will stop.
func (a *Actor) WalkDepthFirst(f func(*Actor) bool) (cont bool) {
	if !f(a) {
		return false
	}
	for child := range a.Children {
		child.WalkDepthFirst(f)
	}
	return true
}

// InitRoot will initialize the entire actor tree.
func (a *Actor) InitRoot() {
	a.WalkDepthFirst(func(b *Actor) bool {
		b.started = &atomic.Bool{}
		b.finished = &atomic.Bool{}
		b.log = log.Logger.With().Str("actor", b.Name).Logger()
		if b.Timeout == 0 {
			b.Timeout = time.Minute
		}
		if b.Interval == 0 {
			b.Interval = time.Second
		}
		for child := range b.Children {
			child.parents[b] = true
		}
		return true
	})

	// mark root as started and finished
	a.started.Store(true)
	a.finished.Store(true)
}

// Append will append the provided actor - all descendants of the actor will be parents
// of the provided actor (the provided actor and descendants will not be executed until
// all parents have completed).
func (a *Actor) Append(b *Actor) {
	// gather all final descendants of a
	descendants := map[*Actor]bool{}
	a.WalkDepthFirst(func(actor *Actor) bool {
		if len(actor.Children) == 0 {
			descendants[actor] = true
		}
		return true
	})

	// set b as a child of all descendants
	for descendant := range descendants {
		descendant.Children[b] = true
	}

	// set all descendants as parents of b
	b.parents = descendants
}

// SetLogger will set the logger for the actor.
func (a *Actor) SetLogger(l zerolog.Logger) {
	a.log = l
}

// Log will return the logger for the actor.
func (a *Actor) Log() *zerolog.Logger {
	return &a.log
}

// Execute will execute the actor.
func (a *Actor) Execute(c *OpConfig) (err error) {
	start := time.Now()

	// mark finished on return
	defer func() {
		a.finished.Store(true)
	}()

	for _, op := range a.Ops {
		for { // run each op until continue or finished
			result := op(c)
			if result.Error != nil {
				a.Log().Err(result.Error).Msg("op failed")
			}
			if result.Finish {
				return result.Error
			}
			if result.Continue {
				break
			}
			if a.Timeout > 0 && time.Since(start) > a.Timeout {
				return ErrTimeout
			}
			time.Sleep(a.Interval)
		}
	}

	return nil
}
