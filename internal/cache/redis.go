// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cache

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
	redis "github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"golang.org/x/xerrors"
)

const (
	CronjobLock = "cronjob-lock"
)

//go:generate mockgen -source=redis.go -destination=mocks/redis.go -package=cache
type Redis interface {
	SetExpire(ctx context.Context, key string, expired time.Time) error
	SAdd(ctx context.Context, key string, value string) error
	SMembers(ctx context.Context, key string) ([]string, error)
	Close() error
	ObtainLock(ctx context.Context, key string, expire time.Duration) *redislock.Lock
}

// Config encapsulates the settings for configuring the redis service.
type Config struct {
	// Redis master node DNS hostname
	Master string

	// Redis sentinel addresses
	SentinelAddrs []string

	// The logger to use. If not defined an output-discarding logger will
	// be used instead.
	Logger *zerolog.Logger
}

type redisImpl struct {
	cfg      Config
	instance *redis.Client
}

func New(cfg Config) Redis {
	impl := &redisImpl{
		cfg: cfg,
		instance: redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    cfg.Master,
			SentinelAddrs: cfg.SentinelAddrs,
		}),
	}

	return impl
}

func (r *redisImpl) SetExpire(ctx context.Context, key string, expired time.Time) error {
	expire, err := r.instance.ExpireAt(ctx, key, expired).Result()
	if err != nil {
		return xerrors.Errorf("redis SetExpire(): key: %s, expired: %s, err: %w", key, expired, err)
	}

	r.cfg.Logger.Info().Msgf("redis SetExpire(): key: %s expired: %t", key, expire)

	return nil
}

func (r *redisImpl) SAdd(ctx context.Context, key, value string) error {
	err := r.instance.SAdd(ctx, key, value).Err()
	if err != nil {
		r.cfg.Logger.Error().Err(err).Msgf("redis SAdd(): key: %s, value: %s", key, value)

		return xerrors.Errorf("redis SAdd(): key: %s, value: %s, err: %w", key, value, err)
	}

	r.cfg.Logger.Info().Msgf("redis SAdd(): key: %s, value: %s", key, value)

	return nil
}

func (r *redisImpl) SMembers(ctx context.Context, key string) ([]string, error) {
	res, err := r.instance.SMembers(ctx, key).Result()
	if err != nil {
		r.cfg.Logger.Error().Err(err).Msgf("redis SMembers(): res: %+v", res)

		return res, xerrors.Errorf("redis SMembers(): key: %s, err: %w", key, err)
	}

	r.cfg.Logger.Info().Msgf("redis SMembers(): res: %+v", res)

	return res, nil
}

func (r *redisImpl) Close() error {
	if err := r.instance.Close(); err != nil {
		return xerrors.Errorf("redis Close(): err: %w", err)
	}

	return nil
}

func (r *redisImpl) ObtainLock(ctx context.Context, key string, expire time.Duration) *redislock.Lock {
	// Create a new lock client.
	locker := redislock.New(r.instance)

	// Try to obtain lock.
	lock, err := locker.Obtain(ctx, key, expire, nil)
	if errors.Is(err, redislock.ErrNotObtained) {
		r.cfg.Logger.Error().Err(err).Msg("redis ObtainLock(): Could not obtain lock!")

		return nil
	} else if err != nil {
		panic(err)
	}

	r.cfg.Logger.Debug().Msgf("redis ObtainLock(): (%s) lock obtained successfully!", key)

	return lock
}
