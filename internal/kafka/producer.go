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
package kafka

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/samwang0723/stock-crawler/internal/helper"
	kafkago "github.com/segmentio/kafka-go"
	"golang.org/x/xerrors"
)

const (
	DailyClosesV1        = "dailycloses-v1"
	StocksV1             = "stocks-v1"
	ThreePrimaryV1       = "threeprimary-v1"
	StakeConcentrationV1 = "stakeconcentration-v1"
	DownloadV1           = "download-v1"
	queueCapacity        = 1024
	sessionTimeout       = 10 * time.Second
	rebalanceTimeout     = 5 * time.Second
	maxWait              = 1 * time.Second
	minBytes             = 1    // 1B
	maxBytes             = 10e6 // 10MB
)

//go:generate mockgen -source=producer.go -destination=mocks/kafka.go -package=kafka
type Kafka interface {
	Close() error
	WriteMessages(ctx context.Context, topic string, message []byte) error
	ReadMessage(ctx context.Context) (*ReceivedMessage, error)
}

// Config encapsulates the settings for configuring the kafka service.
type Config struct {
	// Kafka controller DNS hostname
	Controller string

	GroupID string
	Brokers []string
	Topics  []string

	// The logger to use. If not defined an output-discarding logger will
	// be used instead.
	Logger *zerolog.Logger
}

type ReceivedMessage struct {
	Topic   string
	Message []byte
}

type kafkaImpl struct {
	cfg          *Config
	instance     *kafkago.Writer
	readInstance *kafkago.Reader
}

// NewKafka creates a new kafka service.
//
//nolint:nolintlint, gomnd
func New(cfg *Config) Kafka {
	return &kafkaImpl{
		cfg: cfg,
		instance: &kafkago.Writer{
			Addr:         kafkago.TCP(cfg.Controller),
			Balancer:     &kafkago.LeastBytes{},
			BatchSize:    100,
			BatchTimeout: 100 * time.Millisecond,
		},
		readInstance: kafkago.NewReader(kafkago.ReaderConfig{
			Brokers:          cfg.Brokers,
			GroupTopics:      cfg.Topics,
			GroupID:          cfg.GroupID, // having consumer group id to prevent duplication of message consumption
			QueueCapacity:    queueCapacity,
			SessionTimeout:   sessionTimeout,
			RebalanceTimeout: rebalanceTimeout,
			MaxWait:          maxWait,
			MinBytes:         minBytes,
			MaxBytes:         maxBytes,
			Dialer: &kafkago.Dialer{
				Timeout:       10 * time.Second,
				KeepAlive:     30 * time.Second,
				DualStack:     true,
				FallbackDelay: 10 * time.Millisecond,
			},
		}),
	}
}

func (k *kafkaImpl) ReadMessage(ctx context.Context) (*ReceivedMessage, error) {
	msg, err := k.readInstance.ReadMessage(ctx)
	k.cfg.Logger.Info().Msgf("kafka.ReadMessage: read data: %s, err: %s", helper.Bytes2String(msg.Value), err)

	if err != nil {
		return nil, xerrors.Errorf("kafka.ReadMessage: failed, err=%w;", err)
	}

	return &ReceivedMessage{
		Topic:   msg.Topic,
		Message: msg.Value,
	}, nil
}

func (k *kafkaImpl) WriteMessages(ctx context.Context, topic string, message []byte) error {
	msg := kafkago.Message{
		Topic: topic,
		Value: message,
	}

	err := k.instance.WriteMessages(ctx, msg)
	if err != nil {
		return xerrors.Errorf("kafka.WriteMessages: failed, err=%w;", err)
	}

	k.cfg.Logger.Info().Msgf(
		"kafka.WriteMessages: success, bytes=%d; topic=%s; data=%s;",
		len(message),
		topic,
		helper.Bytes2String(message),
	)

	return nil
}

func (k *kafkaImpl) Close() error {
	if err := k.instance.Close(); err != nil {
		return xerrors.Errorf("kafka.Close: failed, err=%w;", err)
	}

	k.cfg.Logger.Info().Msg("kafka.Close: success")

	return nil
}
