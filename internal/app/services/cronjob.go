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

func (s *serviceImpl) StartCron() {
	if s.cronjob != nil {
		s.cronjob.Start()
	}
}

func (s *serviceImpl) StopCron() {
	if s.cronjob != nil {
		s.cronjob.Stop()
	}
}

func (s *serviceImpl) AddJob(ctx context.Context, spec string, job func()) error {
	if s.cronjob == nil {
		return fmt.Errorf("No cronjob instance being initialized: %+v", s.cronjob)
	}
	return s.cronjob.AddJob(ctx, spec, job)
}
