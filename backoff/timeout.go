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

package backoff

import "time"

// A Clock is used to determine the reference time in time-based logic.
type Clock interface {
	// Time returns the current time.
	Time() time.Time
}

// A ClockFunc is the functional implementation of the [Clock] interface.
type ClockFunc func() time.Time

func (f ClockFunc) Time() time.Time { return f() }

type timeout struct {
	strategy Strategy      // wrapped strategy
	clock    Clock         // determines the reference time
	limit    time.Duration // maximum execution time
}

func (t *timeout) Delay(n int, start time.Time) time.Duration {
	if t.clock.Time().Sub(start) >= t.limit {
		return Exit
	}
	return t.strategy.Delay(n, start)
}

// Timeout wraps a backoff [Strategy] to exit the retry cycle after the given
// duration has passed. The elapsed time is measured relative to the time
// supplied by clock. If limit <= 0, no timeout will be applied.
func Timeout(strategy Strategy, limit time.Duration, clock Clock) Strategy {
	if limit <= 0 {
		return strategy
	}
	return &timeout{
		strategy: strategy,
		limit:    limit,
		clock:    clock,
	}
}
