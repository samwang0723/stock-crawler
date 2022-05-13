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

	"github.com/samwang0723/stock-crawler/internal/cronjob/icronjob"
	"github.com/samwang0723/stock-crawler/internal/helper"
	structuredlog "github.com/samwang0723/stock-crawler/internal/logger/structured"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

type cronjobImpl struct {
	instance *cron.Cron
}

type cronLog struct {
	clog structuredlog.ILogger
}

func (l *cronLog) Info(msg string, keysAndValues ...interface{}) {
	l.clog.RawLogger().WithFields(logrus.Fields{
		"data": keysAndValues,
	}).Info(msg)
}

func (l *cronLog) Error(err error, msg string, keysAndValues ...interface{}) {
	l.clog.RawLogger().WithFields(logrus.Fields{
		"msg":  msg,
		"data": keysAndValues,
	}).Warn(msg)
}

func New(l structuredlog.ILogger) icronjob.ICronJob {
	logger := &cronLog{clog: l}
	logger.clog.RawLogger().SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// load location with Taipei timezone
	location, _ := time.LoadLocation(helper.TimeZone)
	job := &cronjobImpl{
		instance: cron.New(
			cron.WithLocation(location),
			cron.WithLogger(logger),
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
