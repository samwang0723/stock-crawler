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

	config "github.com/samwang0723/stock-crawler/configs"
	"github.com/samwang0723/stock-crawler/internal/kafka/ikafka"
	log "github.com/samwang0723/stock-crawler/internal/logger"
	"github.com/segmentio/kafka-go"
)

type kafkaImpl struct {
	topic    string
	instance *kafka.Conn
}

func New(cfg *config.Config, topic string) ikafka.IKafka {
	conn, err := kafka.DialLeader(context.Background(), "tcp", fmt.Sprintf("%s:%d", cfg.Kafka.Host, cfg.Kafka.Port), topic, 0)
	if err != nil {
		log.Fatal("Failed to dial kafka leader:", err)
	}

	return &kafkaImpl{
		topic:    topic,
		instance: conn,
	}
}

func (k *kafkaImpl) GetTopic() string {
	return k.topic
}

func (k *kafkaImpl) WriteMessages(ctx context.Context, message []byte) error {
	msg := kafka.Message{
		Value: message,
	}

	writtenBytes, err := k.instance.WriteMessages(msg)
	log.Infof("WriteMessages: written bytes: %d", writtenBytes)

	return err
}

func (k *kafkaImpl) Close() error {
	log.Infof("Close: topic: %s", k.topic)
	err := k.instance.Close()
	if err != nil {
		log.Errorf("Close failed: %w", err)
	}
	return err
}
