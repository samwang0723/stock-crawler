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

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type Cronjob interface {
	Start()
	Stop()
	AddJob(ctx context.Context, spec string, job func()) error
}

// Config encapsulates the settings for configuring the redis service.
type Config struct {
	// The logger to use. If not defined an output-discarding logger will
	// be used instead.
	Logger *logrus.Entry
}

type cronjobImpl struct {
	cfg      Config
	instance *cron.Cron
}

func New(cfg Config) Cronjob {
	// load location with Taipei timezone
	location, _ := time.LoadLocation(helper.TimeZone)
	job := &cronjobImpl{
		cfg: cfg,
		instance: cron.New(
			cron.WithLocation(location),
			cron.WithLogger(cfg.Logger),
		),
	}
	return job
}

func (c *cronjobImpl) AddJob(ctx context.Context, spec string, job func()) error {
	_, err := c.instance.AddFunc(spec, job)
	return err
}

func (c *cronjobImpl) Start() {
	c.instance.Start()
}

func (c *cronjobImpl) Stop() {
	c.instance.Stop()
}
