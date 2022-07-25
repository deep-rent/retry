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

// Package backoff provides various backoff strategies.
//
// In particular, the package implements Constant, Linear and Exponential
// backoff strategies as well as some decorators to adjust their behavior.
package backoff

import "time"

// Exit is returned by a backoff Strategy to signal the end of a retry cycle.
var Exit time.Duration = -1

// Strategy determines the delay between consecutive retries in a backoff
// scenario. Implementations of this interface must be stateless.
type Strategy interface {
	// Delay returns the time to wait after the n-th retry of a failing function
	// call. For implementing time-based algorithms, the function also takes the
	// start time of the retry cycle. To stop the cycle after n attempts, the
	// function must return Exit. Note that the initial attempt corresponds to
	// n = 1.
	Delay(n int, start time.Time) time.Duration
}
