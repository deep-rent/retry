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

func clock(now time.Time) backoff.Clock {
	return func() time.Time { return now }
}

func TestTimeoutBelow(t *testing.T) {
	d1 := time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
	d2 := time.Date(0, 0, 0, 0, 0, 1, 0, time.Local)

	s := backoff.Timeout(backoff.Constant(1*time.Second), 2*time.Second, clock(d2))
	act := s.Delay(1, d1)

	const exp = 1 * time.Second

	if act != exp {
		t.Errorf("delay was %s, want %s", act, exp)
	}
}

func TestTimeoutAbove(t *testing.T) {
	d1 := time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)
	d2 := time.Date(0, 0, 0, 0, 0, 2, 0, time.Local)

	s := backoff.Timeout(backoff.Constant(1*time.Second), 1*time.Second, clock(d2))
	act := s.Delay(1, d1)

	exp := backoff.Exit

	if act != exp {
		t.Errorf("delay was %s, want %s", act, exp)
	}
}

func TestTimeoutZero(t *testing.T) {
	d := time.Date(0, 0, 0, 0, 0, 0, 0, time.Local)

	s := backoff.Timeout(backoff.Constant(1*time.Second), 0, clock(d))
	act := s.Delay(1, d)

	const exp = 1 * time.Second

	if act != exp {
		t.Errorf("delay was %s, want %s", act, exp)
	}
}
