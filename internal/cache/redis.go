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

	config "github.com/samwang0723/stock-crawler/configs"
	"github.com/samwang0723/stock-crawler/internal/cache/icache"
	log "github.com/samwang0723/stock-crawler/internal/logger"

	"github.com/go-redis/redis/v8"
)

type redisImpl struct {
	instance *redis.Client
}

func New(cfg *config.Config) icache.IRedis {
	impl := &redisImpl{
		instance: redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    fmt.Sprintf("%s:%d", cfg.RedisCache.Host, cfg.RedisCache.Port),
		}),
	}
	impl.ping()

	return impl
}

func (r *redisImpl) ping() {
	pong, err := r.instance.Ping(context.Background()).Result()
	log.Infof("Ping redis instance: %s", pong)
	if err != nil {
		log.Fatal(err)
	}
}

func (r *redisImpl) SetExpire(ctx context.Context, key string, expired time.Time) error {
	expire, err := r.instance.ExpireAt(ctx, key, expired).Result()
	if err != nil {
		return err
	}
	log.Infof("Redis Key: %s ExpiredAt: %s", key, expire)
	return nil
}

func (r *redisImpl) LPush(ctx context.Context, key string, value string) error {
	return r.instance.LPush(ctx, key, value).Err()
}

func (r *redisImpl) LRange(ctx context.Context, key string) ([]string, error) {
	return r.instance.LRange(ctx, key, 0, -1).Result()
}

func (r *redisImpl) Close() error {
	return r.instance.Close()
}
