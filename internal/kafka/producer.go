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

	"github.com/samwang0723/stock-crawler/internal/helper"
	"golang.org/x/xerrors"

	"github.com/rs/zerolog"
	kafkago "github.com/segmentio/kafka-go"
)

const (
	DailyClosesV1        = "dailycloses-v1"
	StocksV1             = "stocks-v1"
	ThreePrimaryV1       = "threeprimary-v1"
	StakeConcentrationV1 = "stakeconcentration-v1"
)

//go:generate mockgen -source=producer.go -destination=mocks/kafka.go -package=kafka
type Kafka interface {
	Close() error
	WriteMessages(ctx context.Context, topic string, message []byte) error
}

// Config encapsulates the settings for configuring the kafka service.
type Config struct {
	// Kafka controller DNS hostname
	Controller string

	// The logger to use. If not defined an output-discarding logger will
	// be used instead.
	Logger *zerolog.Logger
}

type kafkaImpl struct {
	cfg      Config
	instance *kafkago.Writer
}

// NewKafka creates a new kafka service.
//
//nolint:nolintlint, gomnd
func New(cfg Config) Kafka {
	return &kafkaImpl{
		cfg: cfg,
		instance: &kafkago.Writer{
			Addr:         kafkago.TCP(cfg.Controller),
			Balancer:     &kafkago.LeastBytes{},
			BatchSize:    100,
			BatchTimeout: 100 * time.Millisecond,
		},
	}
}

func (k *kafkaImpl) WriteMessages(ctx context.Context, topic string, message []byte) error {
	msg := kafkago.Message{
		Topic: topic,
		Value: message,
	}

	err := k.instance.WriteMessages(ctx, msg)
	if err != nil {
		return xerrors.Errorf("kafka writeMessages(): %w", err)
	}

	k.cfg.Logger.Info().Msgf(
		"kafka writeMessages(): written bytes: %d, topic: %s, data: %s",
		len(message),
		topic,
		helper.Bytes2String(message),
	)

	return nil
}

func (k *kafkaImpl) Close() error {
	k.cfg.Logger.Info().Msg("kafka close()")

	if err := k.instance.Close(); err != nil {
		k.cfg.Logger.Error().Err(err).Msg("close failed")

		return xerrors.Errorf("kafka close(): %w", err)
	}

	return nil
}
