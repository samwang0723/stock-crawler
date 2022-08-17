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
	"github.com/samwang0723/stock-crawler/internal/app/crawler"
	"github.com/samwang0723/stock-crawler/internal/cache"
	"github.com/samwang0723/stock-crawler/internal/cronjob"
	"github.com/samwang0723/stock-crawler/internal/kafka"
)

type Option func(o *serviceImpl)

func WithCronJob(cfg CronjobConfig) Option {
	return func(i *serviceImpl) {
		i.cronjob = cronjob.New(cronjob.Config{
			Logger: cfg.Logger,
		})
	}
}

func WithKafka(cfg KafkaConfig) Option {
	return func(i *serviceImpl) {
		cfg.validate()
		i.producer = kafka.New(kafka.Config{
			Controller: cfg.Controller,
			Logger:     cfg.Logger,
		})
	}
}

func WithRedis(cfg RedisConfig) Option {
	return func(i *serviceImpl) {
		cfg.validate()
		i.cache = cache.New(cache.Config{
			Master:        cfg.Master,
			SentinelAddrs: cfg.SentinelAddrs,
			Logger:        cfg.Logger,
		})
	}
}

func WithCrawler(cfg CrawlerConfig) Option {
	return func(i *serviceImpl) {
		cfg.validate()
		i.crawler = crawler.New(crawler.Config{
			URLGetter:         cfg.URLGetter,
			FetchWorkers:      cfg.FetchWorkers,
			RateLimitInterval: cfg.RateLimitInterval,
			Proxy:             cfg.Proxy,
			Logger:            cfg.Logger,
		})
	}
}
