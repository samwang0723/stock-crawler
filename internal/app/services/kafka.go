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

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"
	"golang.org/x/xerrors"
)

// Config encapsulates the settings for configuring the kafka service.
type KafkaConfig struct {
	// Kafka controller DNS hostname
	Controller string

	// The logger to use. If not defined an output-discarding logger will
	// be used instead.
	Logger *zerolog.Logger
}

func (cfg *KafkaConfig) validate() error {
	var err error
	if cfg.Controller == "" {
		err = multierror.Append(err, xerrors.Errorf("invalid value for controller hostname"))
	}

	return err
}

func (s *serviceImpl) sendKafka(ctx context.Context, topic string, message []byte) error {
	if s.producer == nil {
		return xerrors.Errorf("kafka producer is not initialized")
	}

	return s.producer.WriteMessages(ctx, topic, message)
}

func (s *serviceImpl) StopKafka() error {
	if s.producer == nil {
		return xerrors.Errorf("kafka producer is not initialized")
	}

	return s.producer.Close()
}
