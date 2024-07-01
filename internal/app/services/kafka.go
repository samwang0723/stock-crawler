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
package services

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/samwang0723/stock-crawler/internal/app/dto"
	"github.com/samwang0723/stock-crawler/internal/kafka"
	"golang.org/x/xerrors"
)

// Config encapsulates the settings for configuring the kafka service.
type KafkaConfig struct {
	// Kafka controller DNS hostname
	Controller string

	GroupID string
	Brokers []string
	Topics  []string

	// The logger to use. If not defined an output-discarding logger will
	// be used instead.
	Logger *zerolog.Logger
}

func (cfg *KafkaConfig) validate() error {
	if cfg.Controller == "" {
		return xerrors.Errorf(
			"service.kafka.validate: failed, reason: invalid kafka config value for controller hostname",
		)
	}

	return nil
}

//
//nolint:nolintlint, cyclop
func (s *serviceImpl) ListeningDownloadRequest(
	ctx context.Context,
	downloadChan chan *dto.StartCronjobRequest,
) {
	go func() {
		for {
			msg, err := s.producer.ReadMessage(ctx)
			if err != nil {
				continue
			}

			request, err := unmarshalMessage(msg)
			if err == nil {
				downloadChan <- request
			}

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()
}

func (s *serviceImpl) sendKafka(ctx context.Context, topic string, message []byte) error {
	if s.producer == nil {
		return xerrors.Errorf("service.sendKafka: failed, reason: producer is not initialized")
	}

	err := s.producer.WriteMessages(ctx, topic, message)
	if err != nil {
		return xerrors.Errorf("service.sendKafka: failed, reason: cannot write message %w", err)
	}

	return nil
}

func (s *serviceImpl) StopKafka() error {
	if s.producer == nil {
		return xerrors.Errorf("service.stopKafka: failed, reason: producer is not initialized")
	}

	if err := s.producer.Close(); err != nil {
		return xerrors.Errorf("service.stopKafka: failed, reason: cannot close producer %w", err)
	}

	return nil
}

func unmarshalMessage(msg *kafka.ReceivedMessage) (*dto.StartCronjobRequest, error) {
	var err error

	var output dto.StartCronjobRequest

	if msg.Topic == kafka.DownloadV1 {
		err = jsoni.Unmarshal(msg.Message, &output)
	}

	if err != nil {
		return nil, xerrors.Errorf("unmarshalMessage: failed, reason: %w", err)
	}

	return &output, nil
}
