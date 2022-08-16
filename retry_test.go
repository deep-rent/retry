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

package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/deep-rent/retry"
	"github.com/deep-rent/retry/backoff"
)

var ErrTest = errors.New("test")

func TestCycler_Try(t *testing.T) {
	cycler := retry.NewCycler(backoff.Constant(1 * time.Millisecond))

	const N = 3
	err := cycler.Try(func(n int) error {
		switch {
		case n < N:
			return ErrTest
		case n > N:
			t.Fatalf("too many attempts: n > %d", N)
			return nil
		default:
			return nil
		}
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestCycler_TryWithContext(t *testing.T) {
	cycler := retry.NewCycler(backoff.Constant(1 * time.Millisecond))

	ctx, cancel := context.WithCancel(context.Background())

	const N = 3
	err := cycler.TryWithContext(ctx, func(n int) error {
		switch {
		case n < N:
			return ErrTest
		case n > N:
			t.Fatalf("too many attempts: n > %d", N)
			return nil
		default:
			cancel()
			return ErrTest
		}
	})

	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("unexpected error: %#v", err)
	}
}

func TestCycler_OnError(t *testing.T) {
	const D = 1 * time.Millisecond
	cycler := retry.NewCycler(backoff.Constant(D))

	const N = 3
	i := 1

	cycler.OnError(func(n int, delay time.Duration, err error) {
		if n > N {
			t.Fatalf("too many attempts: n > %d", N)
		} else if err == nil {
			t.Errorf("expected an error, got nil")
		} else if err != ErrTest {
			t.Errorf("unexpected error: %#v", err)
		} else if n != i {
			t.Fatalf("n = %d, want %d", n, i)
		} else if delay != D {
			t.Errorf("delay = %s, want %s", delay, D)
		} else {
			i++
		}
	})

	_ = cycler.Try(func(n int) error {
		if n == N {
			return nil
		} else {
			return ErrTest
		}
	})

	if i != N {
		t.Fatalf("i = %d, want %d", i, N)
	}
}

func TestCycler_Try_ExitError(t *testing.T) {
	cycler := retry.NewCycler(backoff.Constant(1 * time.Millisecond))

	const N = 3
	err := cycler.Try(func(n int) error {
		switch {
		case n < N:
			return ErrTest
		case n > N:
			t.Fatalf("too many attempts: n > %d", N)
			return nil
		default:
			return retry.ForceExit(ErrTest)
		}
	})

	if err == nil {
		t.Fatalf("expected an error, got nil")
	}

	if err != ErrTest {
		t.Errorf("unexpected error: %#v", err)
	}
}
