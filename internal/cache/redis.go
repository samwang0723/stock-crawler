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
	"time"

	config "github.com/samwang0723/stock-crawler/configs"
	log "github.com/samwang0723/stock-crawler/internal/logger"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

const (
	CronjobLock = "cronjob-lock"
)

type Redis interface {
	SetExpire(ctx context.Context, key string, expired time.Time) error
	SAdd(ctx context.Context, key string, value string) error
	SMembers(ctx context.Context, key string) ([]string, error)
	Close() error
	ObtainLock(ctx context.Context, key string, expire time.Duration) *redislock.Lock
}

type redisImpl struct {
	instance *redis.Client
}

func New(cfg *config.Config) Redis {
	impl := &redisImpl{
		instance: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    cfg.RedisCache.Master,
			SentinelAddrs: cfg.RedisCache.SentinelAddrs,
		}),
	}

	return impl
}

func (r *redisImpl) SetExpire(ctx context.Context, key string, expired time.Time) error {
	expire, err := r.instance.ExpireAt(ctx, key, expired).Result()
	if err != nil {
		return err
	}
	log.Infof("Redis:SetExpire: key: %s expiredAt: %s", key, expire)
	return nil
}

func (r *redisImpl) SAdd(ctx context.Context, key string, value string) error {
	err := r.instance.SAdd(ctx, key, value).Err()
	log.Infof("Redis:SAdd: key: %s, value: %s, err: %w", key, value, err)
	return err
}

func (r *redisImpl) SMembers(ctx context.Context, key string) ([]string, error) {
	res, err := r.instance.SMembers(ctx, key).Result()
	log.Infof("Redis:SMembers: res: %+v, err: %w", res, err)
	return res, err
}

func (r *redisImpl) Close() error {
	return r.instance.Close()
}

func (r *redisImpl) ObtainLock(ctx context.Context, key string, expire time.Duration) *redislock.Lock {
	// Create a new lock client.
	locker := redislock.New(r.instance)

	// Try to obtain lock.
	lock, err := locker.Obtain(ctx, key, expire, nil)
	if err == redislock.ErrNotObtained {
		log.Errorf("Redis:ObtainLock: Could not obtain lock! reason: %w", err)
		return nil
	} else if err != nil {
		log.Panic(err)
	}

	log.Debugf("Redis:ObtainLock: (%s) lock obtained successfully!", key)
	return lock
}
