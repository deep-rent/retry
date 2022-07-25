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

type constant struct {
	d time.Duration
}

func (con *constant) Delay(n int, start time.Time) time.Duration {
	return con.d
}

// Constant returns a backoff Strategy that always returns delay d. The function
// panics if d < 0.
func Constant(d time.Duration) Strategy {
	if d < 0 {
		panic(fmt.Sprintf("d = %s, must be >= 0", d))
	}
	return &constant{d: d}
}

// Once is a backoff Strategy that always returns Exit. Mostly useful for
// testing purposes.
var Once Strategy = &constant{Exit}
