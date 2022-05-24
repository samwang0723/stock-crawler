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
package cache

import (
	"context"
	"os"
	"testing"
	"time"

	log "github.com/samwang0723/stock-crawler/internal/logger"
	logtest "github.com/samwang0723/stock-crawler/internal/logger/structured"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	redismock "github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func setup() {
	logger := logtest.NullLogger()
	log.Initialize(logger)
}

func shutdown() {
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func Test_ObtainLock(t *testing.T) {
	ctx := context.TODO()
	duration := 10 * time.Second
	client, mock := redismock.NewClientMock()
	impl := &redisImpl{
		instance: client,
	}

	tests := []struct {
		name     string
		err      error
		obtained bool
	}{
		{
			name:     "Redis distributed lock obtained successfully",
			obtained: true,
			err:      nil,
		},
		{
			name:     "Redis distributed lock obtain failed",
			obtained: false,
			err:      redislock.ErrNotObtained,
		},
		{
			name:     "Redis distributed lock panic",
			obtained: false,
			err:      redis.ErrClosed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.err {
			case redislock.ErrNotObtained:
				mock.Regexp().ExpectSetNX(CronjobLock, `[a-z]+`, duration).SetErr(tt.err)
				lock := impl.ObtainLock(ctx, CronjobLock, duration)
				assert.Equal(t, tt.obtained, lock != nil)
			case redis.ErrClosed:
				mock.Regexp().ExpectSetNX(CronjobLock, `[a-z]+`, duration).SetErr(redis.ErrClosed)
				assert.Panics(t, func() { impl.ObtainLock(ctx, CronjobLock, duration) }, "The code did not panic")
			default:
				mock.Regexp().ExpectSetNX(CronjobLock, `[a-z]+`, duration).SetVal(true)
				lock := impl.ObtainLock(ctx, CronjobLock, duration)
				assert.Equal(t, tt.obtained, lock != nil)
			}
		})
	}
}
