package types

import (
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

////////////////////////////////////////////////////////////////////////////////////////
// Watcher
////////////////////////////////////////////////////////////////////////////////////////

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

// Execute will execute the watcher.
func (w *Watcher) Execute(c *OpConfig, l zerolog.Logger) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	for ; ; time.Sleep(w.Interval) {
		err := w.Fn(c)
		if err != nil {
			return err
		}
	}
}
