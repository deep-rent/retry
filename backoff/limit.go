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

type limit struct {
	strategy Strategy // wrapped strategy
	n        int      // maximum number of attempts
}

func (lim *limit) Delay(n int, start time.Time) time.Duration {
	if n >= lim.n {
		return Exit
	}
	return lim.strategy.Delay(n, start)
}

// Limit wraps a backoff [Strategy] to end the retry cycle after n attempts. If
// n < 1, no limit will be applied.
func Limit(strategy Strategy, n int) Strategy {
	if n < 1 {
		return strategy
	}
	return &limit{
		strategy: strategy,
		n:        n,
	}
}
