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

type linear struct {
	d time.Duration // initial delay
	k time.Duration // slope
}

func (lin *linear) Delay(n int, start time.Time) time.Duration {
	delay := lin.k*time.Duration(n-1) + lin.d
	if delay < 0 {
		return 0
	}
	return delay
}

// Linear returns a backoff Strategy producing delays that grow linearly in
// in steps of k, starting from the specified initial delay d. If k is negative,
// the delay shrinks to 0 and then stops decreasing. The function panics if
// d is negative.
func Linear(d time.Duration, k time.Duration) Strategy {
	switch {
	case d < 0:
		panic(fmt.Sprintf("d = %s, must be >= 0", d))
	case k == 0:
		return Constant(d)
	default:
		return &linear{
			d: d,
			k: k,
		}
	}
}
