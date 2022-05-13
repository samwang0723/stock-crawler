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
	"github.com/samwang0723/stock-crawler/internal/cache/icache"
	"github.com/samwang0723/stock-crawler/internal/cronjob/icronjob"
	"github.com/samwang0723/stock-crawler/internal/kafka/ikafka"
)

type Option func(o *serviceImpl)

func WithCronJob(cronjob icronjob.ICronJob) Option {
	return func(i *serviceImpl) {
		i.cronjob = cronjob
	}
}

func WithKafka(producer ikafka.IKafka) Option {
	return func(i *serviceImpl) {
		if i.producers == nil {
			i.producers = make(map[string]ikafka.IKafka)
		}
		i.producers[producer.GetTopic()] = producer
	}
}

func WithRedis(redis icache.IRedis) Option {
	return func(i *serviceImpl) {
		i.cache = redis
	}
}
