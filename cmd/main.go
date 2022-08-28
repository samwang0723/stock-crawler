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
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/samwang0723/stock-crawler/internal/app/server"

	"github.com/rs/zerolog/log"
)

//nolint:nolintlint, gochecknoglobals
var (
	appName = "stock-crawler"
)

func main() {
	logger := log.With().Str("app", appName).Logger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-quit:
			logger.Info().Msg("shutdown: started, reason: signal")
			cancel()
		case <-ctx.Done():
		}
	}()

	if err := server.Serve(ctx, &logger); err != nil {
		logger.Error().Err(err).Msg("server.Serve: failed")
	}

	logger.Info().Msg("shutdown: completed")
}
