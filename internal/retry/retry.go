// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package retry

import (
	"errors"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	// DefaultRetryCount is the default retries times
	defaultRetryCount = 2
)

// Retry mechanism
func Retry(attempts int, sleep time.Duration, fnc func() error) error {
	if err := fnc(); err != nil {
		if errors.As(err, &Stop{}) {
			return err
		}

		attempts--
		if attempts > 0 {
			log.Warn().Msgf("retry func error: %s. attempts #%d after %s.", err.Error(), attempts, sleep)
			time.Sleep(sleep)
			// if continue to fail on retry, double the interval
			return Retry(attempts, defaultRetryCount*sleep, fnc)
		}

		return err
	}

	return nil
}

type Stop struct {
	error
}

// If don't want to retry, pass this instead of error
func NoRetryError(err error) Stop {
	return Stop{err}
}
