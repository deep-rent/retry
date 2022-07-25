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

package backoff_test

import (
	"testing"
	"time"

	"github.com/deep-rent/retry/backoff"
)

func TestExponentialIncrease(t *testing.T) {
	s := backoff.Exponential(1*time.Second, 2.0)

	d := time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
	for i, exp := range []time.Duration{
		1 * time.Second,
		2 * time.Second,
		4 * time.Second,
		8 * time.Second,
	} {
		n := i + 1
		act := s.Delay(n, d)

		if act != exp {
			t.Errorf("delay #%d was %s, want %s", n, act, exp)
		}
	}
}

func TestExponentialDecrease(t *testing.T) {
	s := backoff.Exponential(8*time.Second, 0.5)

	d := time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
	for i, exp := range []time.Duration{
		8 * time.Second,
		4 * time.Second,
		2 * time.Second,
		1 * time.Second,
	} {
		n := i + 1
		act := s.Delay(n, d)

		if act != exp {
			t.Errorf("delay #%d was %s, want %s", n, act, exp)
		}
	}
}
