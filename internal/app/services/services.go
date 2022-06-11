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
	"time"

	"github.com/bsm/redislock"
	"github.com/samwang0723/stock-crawler/internal/app/crawler"
	"github.com/samwang0723/stock-crawler/internal/app/graph"
	"github.com/samwang0723/stock-crawler/internal/cache"
	"github.com/samwang0723/stock-crawler/internal/cronjob"
	"github.com/samwang0723/stock-crawler/internal/kafka"
)

type IService interface {
	StartCron()
	StopCron()
	AddJob(ctx context.Context, spec string, job func()) error
	DailyCloseThroughKafka(ctx context.Context, objs *[]interface{}) error
	StockThroughKafka(ctx context.Context, objs *[]interface{}) error
	ThreePrimaryThroughKafka(ctx context.Context, objs *[]interface{}) error
	StakeConcentrationThroughKafka(ctx context.Context, objs *[]interface{}) error
	ObtainLock(ctx context.Context, key string, expire time.Duration) *redislock.Lock
	StopRedis() error
	StopKafka() error
	ListCrawlingConcentrationURLs(ctx context.Context, date string) ([]string, error)
	Crawl(ctx context.Context, linkIt graph.LinkIterator) (int, error)
}

type serviceImpl struct {
	cronjob  cronjob.Cronjob
	producer kafka.Kafka
	cache    cache.Redis
	crawler  crawler.Crawler
}

func New(opts ...Option) IService {
	impl := &serviceImpl{}
	for _, opt := range opts {
		opt(impl)
	}
	return impl
}
