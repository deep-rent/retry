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
	"math"
	"time"
)

type exponential struct {
	d time.Duration // initial delay
	m float64       // exponential multiplier
}

func (exp *exponential) Delay(n int, start time.Time) time.Duration {
	return time.Duration(float64(exp.d) * math.Pow(exp.m, float64(n-1)))
}

// Exponential returns a backoff Strategy producing delays that exponentially
// grow (m > 1), or shrink (m < 1) by the factor m, starting from the
// specified initial delay d. The function panics if d or m are negative.
func Exponential(d time.Duration, m float64) Strategy {
	switch {
	case d < 0:
		panic(fmt.Sprintf("d = %s, must be >= 0", d))
	case m < 0:
		panic(fmt.Sprintf("m = %f, must be >= 0", m))
	case d == 0 || m == 0:
		return Constant(0)
	case m == 1:
		return Constant(d)
	default:
		return &exponential{
			d: d,
			m: m,
		}
	}
}
