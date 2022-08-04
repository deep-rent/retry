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

// Clock returns the current time.
type Clock func() time.Time

type timeout struct {
	strategy Strategy      // wrapped strategy
	max      time.Duration // maximum execution time
	now      Clock         // reference time
}

func (t *timeout) Delay(n int, start time.Time) time.Duration {
	if t.now().Sub(start) >= t.max {
		return Exit
	}
	return t.strategy.Delay(n, start)
}

// Timeout wraps a backoff [Strategy] to exit the retry cycle after the given
// duration has passed. The elapsed time is measured relative to the time
// supplied by now. If max <= 0, no timeout will be applied.
func Timeout(strategy Strategy, max time.Duration, now Clock) Strategy {
	if max <= 0 {
		return strategy
	}
	return &timeout{
		strategy: strategy,
		max:      max,
		now:      now,
	}
}
