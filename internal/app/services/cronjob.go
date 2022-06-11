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

	"github.com/sirupsen/logrus"
)

// Config encapsulates the settings for configuring the redis service.
type CronjobConfig struct {
	// The logger to use. If not defined an output-discarding logger will
	// be used instead.
	Logger *logrus.Entry
}

func (cfg *CronjobConfig) validate() error {
	var err error
	if cfg.Logger == nil {
		cfg.Logger = logrus.NewEntry(&logrus.Logger{Out: ioutil.Discard})
	}

	return err
}

func (s *serviceImpl) StartCron() {
	s.cronjob.Start()
}

func (s *serviceImpl) StopCron() {
	s.cronjob.Stop()
}

func (s *serviceImpl) AddJob(ctx context.Context, spec string, job func()) error {
	return s.cronjob.AddJob(ctx, spec, job)
}
