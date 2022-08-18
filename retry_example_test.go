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
	"fmt"
	"time"

	"github.com/deep-rent/retry"
	"github.com/deep-rent/retry/backoff"
)

// This example uses exponential backoff to retry a dummy function.
func ExampleCycler() {
	exp := backoff.Exponential(2*time.Millisecond, 2.0)

	cycler := retry.NewCycler(exp)
	cycler.Cap(2 * time.Second)      // cap backoff delay at 2 seconds
	cycler.Timeout(15 * time.Second) // stop retrying after 15 seconds
	// cycler.Jitter(0.5)            // introduce 50% jitter

	// register an error handler
	cycler.OnError(func(n int, delay time.Duration, err error) {
		ms := delay.Milliseconds()
		fmt.Printf("attempt #%d: %v => wait %2d ms\n", n, err, ms)
	})

	const N = 5 // number of tries

	// start retry cycle
	err := cycler.Try(func(n int) error {
		if n == N {
			// succeed after 5 attempts
			return nil
		} else {
			// force retry
			return errors.New("failed")
		}
	})

	if err != nil {
		fmt.Printf("failed after retries: %v", err)
	} else {
		fmt.Printf("attempt #%d: succeeded", N)
	}

	// Output:
	// attempt #1: failed => wait  2 ms
	// attempt #2: failed => wait  4 ms
	// attempt #3: failed => wait  8 ms
	// attempt #4: failed => wait 16 ms
	// attempt #5: succeeded
}

// This example uses linear backoff to retry a dummy function.
func ExampleCycler_Try() {
	lin := backoff.Linear(5*time.Millisecond, 5*time.Millisecond)

	cycler := retry.NewCycler(lin)
	cycler.Limit(10) // stop retrying after 10 attempts

	// register an error handler
	cycler.OnError(func(n int, delay time.Duration, err error) {
		ms := delay.Milliseconds()
		fmt.Printf("attempt #%d: %v => wait %2d ms\n", n, err, ms)
	})

	const N = 5 // number of tries

	// start retry cycle
	err := cycler.Try(func(n int) error {
		if n == N {
			// succeed after 5 attempts
			return nil
		} else {
			// force retry
			return errors.New("failed")
		}
	})

	if err != nil {
		fmt.Printf("failed after retries: %v", err)
	} else {
		fmt.Printf("attempt #%d: succeeded", N)
	}

	// Output:
	// attempt #1: failed => wait  5 ms
	// attempt #2: failed => wait 10 ms
	// attempt #3: failed => wait 15 ms
	// attempt #4: failed => wait 20 ms
	// attempt #5: succeeded
}

// This example uses a cancellable context to stop a retry cycle.
func ExampleCycler_TryWithContext() {
	con := backoff.Constant(10 * time.Millisecond)

	cycler := retry.NewCycler(con)
	ctx, cancel := context.WithCancel(context.Background())

	// register an error handler
	cycler.OnError(func(n int, delay time.Duration, err error) {
		ms := delay.Milliseconds()
		fmt.Printf("attempt #%d: %v => wait %2d ms\n", n, err, ms)
	})

	const N = 5 // number of tries

	// start retry cycle
	err := cycler.TryWithContext(ctx, func(n int) error {
		if n == N {
			cancel()
			// succeed after 5 attempts
			return nil
		}
		// force retry
		return errors.New("failed")
	})

	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Printf("failed after retries: %v", err)
	} else {
		fmt.Printf("attempt #%d: succeeded", N)
	}

	// Output:
	// attempt #1: failed => wait 10 ms
	// attempt #2: failed => wait 10 ms
	// attempt #3: failed => wait 10 ms
	// attempt #4: failed => wait 10 ms
	// attempt #5: succeeded
}
