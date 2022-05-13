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
	"fmt"
)

func (s *serviceImpl) sendKafka(ctx context.Context, topic string, message []byte) error {
	p := s.producers[topic]
	if p == nil {
		return fmt.Errorf("No kafka instance being initialized: %+v, topic: %s", s.producers, topic)
	}

	p.WriteMessages(ctx, message)
	return nil
}

func (s *serviceImpl) StopKafka() error {
	var err error
	for _, producer := range s.producers {
		if res := producer.Close(); res != nil {
			err = res
		}
	}
	return err
}
