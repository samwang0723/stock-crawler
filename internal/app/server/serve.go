// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog"
	config "github.com/samwang0723/stock-crawler/configs"
	"github.com/samwang0723/stock-crawler/internal/app/crawler"
	"github.com/samwang0723/stock-crawler/internal/app/dto"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/app/handlers"
	"github.com/samwang0723/stock-crawler/internal/app/services"
	"github.com/samwang0723/stock-crawler/internal/helper"

	"github.com/heptiolabs/healthcheck"
)

const (
	gracefulShutdownPeriod = 5 * time.Second
)

type IServer interface {
	Name() string
	Handler() handlers.IHandler
	Config() *config.SystemConfig
	Run(context.Context) error
	Start(context.Context) error
	Stop(context.Context) error
}

type server struct {
	opts Options
}

func Serve(ctx context.Context, logger *zerolog.Logger) error {
	config.Load()
	cfg := config.GetCurrentConfig()
	// bind DAL layer with service
	dataService := services.New(
		services.WithCronJob(services.CronjobConfig{
			Logger: logger,
		}),
		services.WithKafka(services.KafkaConfig{
			Controller: cfg.Kafka.Controller,
			Logger:     logger,
		}),
		services.WithRedis(services.RedisConfig{
			Master:        cfg.RedisCache.Master,
			SentinelAddrs: cfg.RedisCache.SentinelAddrs,
			Logger:        logger,
		}),
		services.WithCrawler(services.CrawlerConfig{
			FetchWorkers:      cfg.Crawler.FetchWorkers,
			RateLimitInterval: cfg.Crawler.RateLimit,
			Proxy:             &crawler.Proxy{Type: crawler.WebScraping},
			Logger:            logger,
		}),
	)
	// associate service with handler
	handler := handlers.New(dataService, logger)

	// health check
	health := healthcheck.NewHandler()
	// our app is not happy if we've got more than 10k goroutines running.
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(cfg.Server.MaxGoroutine))
	// our app is not ready if we can't resolve our upstream dependency in DNS.
	health.AddReadinessCheck(
		"upstream-redis-dns",
		healthcheck.DNSResolveCheck(cfg.RedisCache.Master, time.Duration(cfg.Server.DNSLatency)))
	health.AddReadinessCheck(
		"upstream-kafka-dns",
		healthcheck.DNSResolveCheck(cfg.Kafka.Controller, time.Duration(cfg.Server.DNSLatency)))

	healthServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: health,
	}

	svc := newServer(
		Name(cfg.Server.Name),
		Config(cfg),
		Handler(handler),
		HealthCheck(healthServer),
		BeforeStart(func() error {
			dataService.StartCron()

			return nil
		}),
		BeforeStop(func() error {
			dataService.StopCron()
			err := dataService.StopRedis()
			if err != nil {
				return fmt.Errorf("data service stop redis failed: %w", err)
			}

			err = dataService.StopKafka()
			if err != nil {
				return fmt.Errorf("data service stop kafka failed: %w", err)
			}

			return nil
		}),
	)

	if err := svc.Run(ctx); err != nil {
		return fmt.Errorf("server run failed: %w", err)
	}

	return nil
}

func newServer(opts ...Option) IServer {
	option := Options{}
	for _, opt := range opts {
		opt(&option)
	}

	return &server{
		opts: option,
	}
}

func (s *server) Start(ctx context.Context) error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	signatureOut := fmt.Sprintf(helper.Signature, "v2.0.0", helper.GetCurrentEnv())

	fmt.Println(signatureOut)

	go func() {
		// start healthcheck specific server
		s.HealthCheck().ListenAndServe()
	}()

	return nil
}

func (s *server) Stop(ctx context.Context) error {
	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			return fmt.Errorf("server before stop failed: %w", err)
		}
	}

	// shutdown healthcheck server
	sctx, cancel := context.WithTimeout(ctx, gracefulShutdownPeriod)
	defer cancel()

	err := s.HealthCheck().Shutdown(sctx)
	if err != nil {
		return fmt.Errorf("server stop failed: %w", err)
	}

	return nil
}

// Run starts the server and shut down gracefully afterwards
func (s *server) Run(ctx context.Context) error {
	if err := s.Start(ctx); err != nil {
		return err
	}

	// Execute logics
	var waitGroup sync.WaitGroup

	waitGroup.Add(1)

	go func(ctx context.Context, svc *server) {
		defer waitGroup.Done()

		// by default starting cronjob for regular daily updates pulling
		// cronjob using redis distrubted lock to prevent multiple instances
		// pulling same content
		svc.Handler().CronDownload(ctx, &dto.StartCronjobRequest{
			Schedule: "30 16 * * 1-5",
			Types: []convert.Source{
				convert.TwseDailyClose,
				convert.TwseThreePrimary,
				convert.TpexDailyClose,
				convert.TpexThreePrimary,
			},
		})
		svc.Handler().CronDownload(ctx, &dto.StartCronjobRequest{
			Schedule: "30 18 * * 1-5",
			Types:    []convert.Source{convert.StakeConcentration},
		})
		// backfill failed concentration records
		svc.Handler().CronDownload(ctx, &dto.StartCronjobRequest{
			Schedule: "30 19 * * 1-5",
			Types:    []convert.Source{convert.StakeConcentration},
		})

		<-ctx.Done()
	}(ctx, s)
	waitGroup.Wait()

	return s.Stop(ctx)
}

func (s *server) Name() string {
	return s.opts.Name
}

func (s *server) Handler() handlers.IHandler {
	return s.opts.Handler
}

func (s *server) Config() *config.SystemConfig {
	return s.opts.Config
}

func (s *server) HealthCheck() *http.Server {
	return s.opts.HealthCheck
}
