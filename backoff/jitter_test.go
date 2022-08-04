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

func random(r float64) backoff.Random {
	return func() float64 { return r }
}

func TestJitterSmall(t *testing.T) {
	s := backoff.Jitter(backoff.Constant(1*time.Second), 0.25, random(0.25))
	act := s.Delay(1, time.Date(0, 0, 0, 0, 0, 0, 0, time.Local))

	const exp = 875 * time.Millisecond

	if act != exp {
		t.Errorf("delay was %s, want %s", act, exp)
	}
}

func TestJitterLarge(t *testing.T) {
	s := backoff.Jitter(backoff.Constant(1*time.Second), 0.75, random(0.75))
	act := s.Delay(1, time.Date(0, 0, 0, 0, 0, 0, 0, time.Local))

	const exp = 1375 * time.Millisecond

	if act != exp {
		t.Errorf("delay was %s, want %s", act, exp)
	}
}

func TestJitterMean(t *testing.T) {
	s := backoff.Jitter(backoff.Constant(1*time.Second), 0.5, random(0.5))
	act := s.Delay(1, time.Date(0, 0, 0, 0, 0, 0, 0, time.Local))

	const exp = 1 * time.Second

	if act != exp {
		t.Errorf("delay was %s, want %s", act, exp)
	}
}

func TestJitterExit(t *testing.T) {
	s := backoff.Jitter(backoff.Once, 0.5, random(0.5))
	act := s.Delay(1, time.Date(0, 0, 0, 0, 0, 0, 0, time.Local))

	exp := backoff.Exit

	if act != exp {
		t.Errorf("delay was %s, want %s", act, exp)
	}
}

func TestJitterZero(t *testing.T) {
	s := backoff.Jitter(backoff.Constant(1*time.Second), 0, random(0.5))
	act := s.Delay(1, time.Date(0, 0, 0, 0, 0, 0, 0, time.Local))

	const exp = 1 * time.Second

	if act != exp {
		t.Errorf("delay was %s, want %s", act, exp)
	}
}
