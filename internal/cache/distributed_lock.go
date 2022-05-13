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
	"fmt"
	"time"

	log "github.com/samwang0723/stock-crawler/internal/logger"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	config "github.com/samwang0723/stock-crawler/configs"
)

const (
	CronjobLock = "cronjob-lock"
)

func ObtainLock(key string, expire time.Duration) *redislock.Lock {
	cfg := config.GetCurrentConfig().RedisCache
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	})
	defer client.Close()

	// Create a new lock client.
	locker := redislock.New(client)

	// Try to obtain lock.
	ctx := context.Background()
	lock, err := locker.Obtain(ctx, key, expire, nil)
	if err == redislock.ErrNotObtained {
		log.Errorf("Could not obtain lock! reason: %s", err)
		return nil
	} else if err != nil {
		log.Fatal(err)
	}

	log.Debugf("(%s) redis lock obtained successfully!", key)
	return lock
}
