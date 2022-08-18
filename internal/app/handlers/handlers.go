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
package handlers

import (
	"context"

	"github.com/samwang0723/stock-crawler/internal/app/dto"
	"github.com/samwang0723/stock-crawler/internal/app/services"

	"github.com/rs/zerolog"
)

type IHandler interface {
	CronDownload(ctx context.Context, req *dto.StartCronjobRequest) error
	Download(ctx context.Context, req *dto.StartCronjobRequest)
}

type handlerImpl struct {
	logger      *zerolog.Logger
	dataService services.IService
}

func New(dataService services.IService, logger *zerolog.Logger) IHandler {
	res := &handlerImpl{
		logger:      logger,
		dataService: dataService,
	}

	return res
}
