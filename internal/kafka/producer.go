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
package kafka

import (
	"context"
	"fmt"
	"time"

	config "github.com/samwang0723/stock-crawler/configs"
	"github.com/samwang0723/stock-crawler/internal/helper"
	"github.com/samwang0723/stock-crawler/internal/kafka/ikafka"
	log "github.com/samwang0723/stock-crawler/internal/logger"
	"github.com/segmentio/kafka-go"
)

type kafkaImpl struct {
	instance *kafka.Writer
}

func New(cfg *config.Config) ikafka.IKafka {
	return &kafkaImpl{
		instance: &kafka.Writer{
			Addr:         kafka.TCP(fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port)),
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    100,
			BatchTimeout: 100 * time.Millisecond,
		},
	}
}

func (k *kafkaImpl) WriteMessages(ctx context.Context, topic string, message []byte) error {
	msg := kafka.Message{
		Topic: topic,
		Value: message,
	}
	err := k.instance.WriteMessages(ctx, msg)
	log.Infof("Kafka:WriteMessages: written bytes: %d, topic: %s, data: %s, err: %s", len(message), topic, helper.Bytes2String(message), err)

	return err
}

func (k *kafkaImpl) Close() error {
	log.Info("Kafka:Close")
	err := k.instance.Close()
	if err != nil {
		log.Errorf("Close failed: %w", err)
	}
	return err
}
