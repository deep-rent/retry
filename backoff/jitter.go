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

import (
	"fmt"
	"time"
)

// Random returns a pseudo-random number in the half-open interval [0,1).
type Random func() float64

type jitter struct {
	strategy Strategy // wrapped strategy
	spread   float64  // spread factor
	random   Random   // random number generator
}

func (j *jitter) Delay(n int, start time.Time) (delay time.Duration) {
	delay = j.strategy.Delay(n, start)
	if delay == Exit {
		return
	}
	w := float64(delay) * j.spread
	return time.Duration(float64(delay) - w + (j.random() * (2*w + 1)))
}

// Jitter wraps a backoff Strategy to randomly spread produced delays around in
// time. The spread factor determines the relative range in which delays are
// scattered. It must fall in the half-open interval [0,1). For example, a
// spread of 0.5 results in delays ranging between 50% above and 50% below the
// values produced by the wrapped strategy. If spread = 0, no jitter will be
// applied.
func Jitter(strategy Strategy, spread float64, random Random) Strategy {
	if spread < 0.0 || spread >= 1.0 {
		panic(fmt.Sprintf("spread %f not in [0,1)", spread))
	}
	if spread == 0 {
		return strategy
	}
	return &jitter{
		strategy: strategy,
		spread:   spread,
		random:   random,
	}
}
