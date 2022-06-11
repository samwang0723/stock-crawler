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
package services

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/bsm/redislock"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

// Config encapsulates the settings for configuring the redis service.
type RedisConfig struct {
	// Redis master node DNS hostname
	Master string

	// Redis sentinel addresses
	SentinelAddrs []string

	// The logger to use. If not defined an output-discarding logger will
	// be used instead.
	Logger *logrus.Entry
}

func (cfg *RedisConfig) validate() error {
	var err error
	if cfg.Master == "" {
		err = multierror.Append(err, xerrors.Errorf("invalid value for master hostname"))
	}
	if len(cfg.SentinelAddrs) == 0 {
		err = multierror.Append(err, xerrors.Errorf("invalid value for sentinel addresses"))
	}
	if cfg.Logger == nil {
		cfg.Logger = logrus.NewEntry(&logrus.Logger{Out: ioutil.Discard})
	}

	return err
}

func (s *serviceImpl) ObtainLock(ctx context.Context, key string, expire time.Duration) *redislock.Lock {
	return s.cache.ObtainLock(ctx, key, expire)
}

func (s *serviceImpl) StopRedis() error {
	return s.cache.Close()
}
