/*
Copyright (c) 2022 deep.rent GmbH (https://deep.rent)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package retry implements a retry mechanism based on configurable backoff
// strategies.
//
// The fundamental structure is called [Cycler]. A cycler can be obtained by
// passing an appropriate [backoff.Strategy] to [NewCycler]. Any function whose
// signature matches [AttemptFunc] can then be retried using either [Cycler.Try]
// or [Cycler.TryWithContext].
package retry

import (
	"context"
	"math/rand"
	"time"

	"github.com/deep-rent/retry/backoff"
)

type (
	// An AttemptFunc can be scheduled in a retry cycle. The function will be
	// retried if it returns an error, while returning nil indicates successful
	// completion. The argument n is the current attempt count, starting at
	// n = 1.
	AttemptFunc func(n int) error

	// An ErrorHandlerFunc is invoked when the n-th execution of an
	// [AttemptFunc] failed with err, and the next retry is pending after delay
	// has passed. Note that the initial execution corresponds to n = 1.
	ErrorHandlerFunc func(n int, delay time.Duration, err error)
)

// seed a new pseudo-random number generator
var rd *rand.Rand = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

// random is the default implementation of [backoff.Random].
func random() float64 {
	return rd.Float64()
}

// An ExitError signals that an [AttemptFunc] should no longer be retried. Use
// [ForceExit] to wrap an error such that it forces the current retry cycle to
// exit. This is useful when an error is encountered that the program cannot
// possibly recover after additional retries.
type ExitError struct {
	Cause error
}

func (e *ExitError) Error() string { return e.Cause.Error() }

// ForceExit wraps err in an [ExitError].
func ForceExit(err error) error {
	return &ExitError{Cause: err}
}

// now is the default implementation of [backoff.Clock].
var now backoff.Clock = backoff.ClockFunc(func() time.Time {
	return time.Now()
})

// A Cycler is used to schedule retry cycles in which an [AttemptFunc] is
// repeatedly executed until it succeeds. Once configured, the same cycler can
// be used to schedule any number of retry cycles.
type Cycler struct {
	strategy backoff.Strategy
	handlers []ErrorHandlerFunc
	Clock    backoff.Clock // used to track the execution time of retry cycles
}

// NewCycler creates a new retry [Cycler]. The specified [backoff.Strategy]
// determines the backoff delay between consecutive attempts. A cycler is meant
// to be reused; recreating the same cycler should be avoided.
func NewCycler(strategy backoff.Strategy) *Cycler {
	return &Cycler{
		strategy: strategy,
		Clock:    now,
	}
}

// OnError registers a callback to be invoked when a failed [AttemptFunc] needs
// to be retried. Typically, these callbacks are used to log intermediate errors
// that would otherwise remain unhandled.
func (c *Cycler) OnError(handler ErrorHandlerFunc) {
	c.handlers = append(c.handlers, handler)
}

// Cap sets the maximum delay between consecutive attempts. If max <= 0, no
// limit will be applied.
func (c *Cycler) Cap(max time.Duration) {
	c.strategy = backoff.Cap(c.strategy, max)
}

// Jitter randomly spreads delays between consecutive attempts around in time.
// The spread factor determines the relative range in which delays are
// scattered. It must fall in the half-open interval [0,1). For example, a
// spread of 0.5 results in delays ranging between 50% above and 50% below the
// values produced by the underlying backoff strategy. If spread = 0, no jitter
// will be applied.
func (c *Cycler) Jitter(spread float64) {
	c.strategy = backoff.Jitter(c.strategy, spread, random)
}

// Limit sets the maximum number of attempts in a retry cycle. A retry cycle
// will stop after the n-th attempt. If n < 1, no limit will be applied.
func (c *Cycler) Limit(n int) {
	c.strategy = backoff.Limit(c.strategy, n)
}

// Timeout sets the maximum duration of retry cycles. A retry cycle will stop
// after the time elapsed since it was scheduled goes past the maximum. If
// limit <= 0, no timeout will be applied.
func (c *Cycler) Timeout(limit time.Duration) {
	c.strategy = backoff.Timeout(c.strategy, limit, c.Clock)
}

// Try calls [TryWithContext] using [context.Background].
func (c *Cycler) Try(attempt AttemptFunc) error {
	return c.TryWithContext(context.Background(), attempt)
}

// TryWithContext schedules a retry cycle in which attempt is repeatedly
// executed until it returns nil. The cycle stops early if
//
//  1. some limit is exceeded,
//  2. ctx is cancelled, or
//  3. an [ExitError] occurs.
//
// When an invocation of attempt returns nil before the cycle stops, this method
// also returns nil. Otherwise, this method returns the last error returned by
// attempt. If ctx contains an error, this error will be returned instead.
//
// In any case, attempt is guaranteed to be executed at least once. Be aware
// that retry cycles with neither [Cycler.Limit] nor [Cycler.Timeout] set will
// run forever if attempt keeps failing.
func (c *Cycler) TryWithContext(
	ctx context.Context,
	attempt AttemptFunc,
) error {
	var t *time.Timer
	defer func() {
		if t != nil {
			t.Stop()
		}
	}()

	n := 0                  // number of attempts
	start := c.Clock.Time() // current time

	// retry loop
	for {
		// increase attempt count
		n++

		err := attempt(n)
		if err == nil {
			// success
			return nil
		}

		// unrecoverable error
		if e, ok := err.(*ExitError); ok {
			return e.Cause
		}

		delay := c.strategy.Delay(n, start)

		if delay == backoff.Exit {
			e := ctx.Err()
			if e != nil {
				err = e
			}
			// exit early
			return err
		}

		// notify error handlers
		if c.handlers != nil {
			for _, h := range c.handlers {
				h(n, delay, err)
			}
		}

		if t == nil {
			t = time.NewTimer(delay)
		} else {
			t.Reset(delay)
		}

		select {
		case <-ctx.Done():
			// exit early
			return ctx.Err()
		case <-t.C:
			// wait for delay to elapse
		}
	}
}
