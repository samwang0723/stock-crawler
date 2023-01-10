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

package cronjob

import (
	"context"
	"time"

	"github.com/samwang0723/stock-crawler/internal/helper"

	cron "github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"golang.org/x/xerrors"
)

//go:generate mockgen -source=cronjob.go -destination=mocks/cronjobs.go -package=cronjob
type Cronjob interface {
	Start()
	Stop()
	AddJob(ctx context.Context, spec string, job func()) error
}

type Config struct {
	Logger *zerolog.Logger
}

type cronjobImpl struct {
	instance *cron.Cron
}

func New(cfg Config) Cronjob {
	// load location with Taipei timezone
	location, err := time.LoadLocation(helper.TimeZone)
	if err != nil {
		cfg.Logger.Fatal().Err(err).Msg("cronjob.New: failed")

		return nil
	}

	job := &cronjobImpl{
		instance: cron.New(
			cron.WithLocation(location),
			cron.WithLogger(cfg),
		),
	}

	return job
}

func (c Config) Info(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) > 0 {
		c.Logger.Info().Msgf("cronjob.Run: success, data=%+v;", keysAndValues)
	}
}

func (c Config) Error(err error, msg string, keysAndValues ...interface{}) {
	if len(keysAndValues) > 0 {
		c.Logger.Warn().Msgf("cronjob.Run: failed, data=%+v;", keysAndValues)
	}
}

func (c *cronjobImpl) AddJob(ctx context.Context, spec string, job func()) error {
	_, err := c.instance.AddFunc(spec, job)
	if err != nil {
		return xerrors.Errorf("cronjob.AddJob: failed, spec=%s, err=%w", spec, err)
	}

	return nil
}

func (c *cronjobImpl) Start() {
	c.instance.Start()
}

func (c *cronjobImpl) Stop() {
	c.instance.Stop()
}
