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

package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heptiolabs/healthcheck"
	config "github.com/samwang0723/stock-crawler/configs"
	"github.com/samwang0723/stock-crawler/internal/app/dto"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/app/handlers"
	"github.com/samwang0723/stock-crawler/internal/app/services"
	"github.com/samwang0723/stock-crawler/internal/cache"
	"github.com/samwang0723/stock-crawler/internal/cronjob"
	"github.com/samwang0723/stock-crawler/internal/helper"
	"github.com/samwang0723/stock-crawler/internal/kafka"
	log "github.com/samwang0723/stock-crawler/internal/logger"
	structuredlog "github.com/samwang0723/stock-crawler/internal/logger/structured"
)

const (
	gracefulShutdownPeriod = 5 * time.Second
)

type IServer interface {
	Name() string
	Logger() structuredlog.ILogger
	Handler() handlers.IHandler
	Config() *config.SystemConfig
	Run(context.Context) error
	Start(context.Context) error
	Stop() error
}

type server struct {
	opts Options
}

func Serve() {
	config.Load()
	cfg := config.GetCurrentConfig()
	logger := structuredlog.Logger(cfg)
	// bind DAL layer with service
	dataService := services.New(
		services.WithCronJob(cronjob.New(logger)),
		services.WithKafka(kafka.New(cfg)),
		services.WithRedis(cache.New(cfg)),
		services.WithCrawler(services.CrawlerConfig{
			FetchWorkers:      10,
			RateLimitInterval: 3000,
		}),
	)
	// associate service with handler
	handler := handlers.New(dataService)

	//health check
	health := healthcheck.NewHandler()
	// Our app is not happy if we've got more than 100 goroutines running.
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(10000))
	// Our app is not ready if we can't resolve our upstream dependency in DNS.
	health.AddReadinessCheck(
		"upstream-redis-dns",
		healthcheck.DNSResolveCheck(cfg.RedisCache.Master, 200*time.Millisecond))
	health.AddReadinessCheck(
		"upstream-kafka-dns",
		healthcheck.DNSResolveCheck(cfg.Kafka.Host, 200*time.Millisecond))
	healthServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler: health,
	}

	s := newServer(
		Name(cfg.Server.Name),
		Config(cfg),
		Logger(logger),
		Handler(handler),
		HealthCheck(healthServer),
		BeforeStart(func() error {
			dataService.StartCron()
			return nil
		}),
		BeforeStop(func() error {
			dataService.StopCron()
			dataService.StopRedis()
			dataService.StopKafka()
			return nil
		}),
	)

	log.Initialize(s.Logger())
	err := s.Run(context.Background())
	if err != nil && s.Logger() != nil {
		log.Errorf("error returned by service.Run(): %s\n", err.Error())
	}
}

func newServer(opts ...Option) IServer {
	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &server{
		opts: o,
	}
}

func (s *server) Start(ctx context.Context) error {
	var err error
	for _, fn := range s.opts.BeforeStart {
		if err = fn(); err != nil {
			return err
		}
	}

	signature := `
 _____ _             _                                  _           
/  ___| |           | |                                | |          
\ '--.| |_ ___   ___| | ________ ___ _ __ __ ___      _| | ___ _ __ 
 '--. \ __/ _ \ / __| |/ /______/ __| '__/ _' \ \ /\ / / |/ _ \ '__|
/\__/ / || (_) | (__|   <      | (__| | | (_| |\ V  V /| |  __/ |   
\____/ \__\___/ \___|_|\_\      \___|_|  \__,_| \_/\_/ |_|\___|_|

                                                        Version (%s)
Stand-alone stock data crawling service
Environment (%s)
_______________________________________________
`
	signatureOut := fmt.Sprintf(signature, "v1.1.0", helper.GetCurrentEnv())
	fmt.Println(signatureOut)

	// by default starting cronjob for regular daily updates pulling
	// cronjob using redis distrubted lock to prevent multiple instances
	// pulling same content
	s.Handler().CronDownload(ctx, &dto.StartCronjobRequest{
		Schedule: "30 16 * * 1-5",
		Types: []convert.Source{
			convert.TwseDailyClose,
			convert.TpexDailyClose,
			convert.TwseThreePrimary,
			convert.TpexThreePrimary,
		},
	})
	s.Handler().CronDownload(ctx, &dto.StartCronjobRequest{
		Schedule: "30 18 * * 1-5",
		Types:    []convert.Source{convert.StakeConcentration},
	})
	// backfill failed concentration records
	s.Handler().CronDownload(ctx, &dto.StartCronjobRequest{
		Schedule: "30 19 * * 1-5",
		Types:    []convert.Source{convert.StakeConcentration},
	})

	// start healthcheck specific server
	go func() {
		err = s.HealthCheck().ListenAndServe()
	}()

	return err
}

func (s *server) Stop() error {
	var err error
	for _, fn := range s.opts.BeforeStop {
		if err = fn(); err != nil {
			break
		}
	}

	// shutdown healthcheck server
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownPeriod)
	defer cancel()
	err = s.HealthCheck().Shutdown(ctx)

	log.Warn("server being gracefully shuted down")

	return err
}

// Run starts the server and shut down gracefully afterwards
func (s *server) Run(ctx context.Context) error {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := s.Start(childCtx); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		log.Warn("singal interrupt")
		cancel()
	case <-childCtx.Done():
		log.Warn("main context being cancelled")
	}
	return s.Stop()
}

func (s *server) Logger() structuredlog.ILogger {
	return s.opts.Logger
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
