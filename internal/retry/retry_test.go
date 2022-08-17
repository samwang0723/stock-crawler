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
	"flag"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func Test_Retry(t *testing.T) {
	tests := []struct {
		name  string
		err   error
		count int
	}{
		{
			name:  "Retry with success func",
			err:   nil,
			count: 1,
		},
		{
			name:  "NoRetry with error func",
			err:   NoRetryError(errors.New("No retry")),
			count: 1,
		},
		{
			name:  "Retry with error func",
			err:   errors.New("Need retry"),
			count: 3,
		},
	}

	attempts := new(int)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			*attempts = 0
			Retry(3, 10*time.Millisecond, log.With().Str("test", "retry").Logger(), func() error {
				*attempts++
				return tt.err
			})
			assert.Equal(t, tt.count, *attempts)
		})
	}
}
