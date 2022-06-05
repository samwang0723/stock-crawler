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
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samwang0723/stock-crawler/internal/app/server"
	"github.com/samwang0723/stock-crawler/internal/helper"

	"github.com/sirupsen/logrus"
)

var (
	appName = "stock-crawler"
)

func main() {
	host, _ := os.Hostname()
	rootLogger := logrus.New()
	logger := rootLogger.WithFields(logrus.Fields{
		"app":  appName,
		"host": host,
	})

	// manually set time zone, docker image may not have preset timezone
	var err error
	time.Local, err = time.LoadLocation(helper.TimeZone)
	if err != nil {
		logrus.WithField("err", err).Errorf("error loading location '%s': %v\n", helper.TimeZone, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		select {
		case s := <-quit:
			logger.WithField("signal", s.String()).Infof("shutting down due to signal")
			cancel()
		case <-ctx.Done():
		}
	}()

	server.Serve(ctx)
	logger.Info("shutdown complete")
}
