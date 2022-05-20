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
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/samwang0723/stock-crawler/internal/app/dto"
	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/cache"
	log "github.com/samwang0723/stock-crawler/internal/logger"
)

const (
	StartCronjob = "START_CRON"
)

func (h *handlerImpl) CronDownload(ctx context.Context, req *dto.StartCronjobRequest) (*dto.StartCronjobResponse, error) {
	envCron := os.Getenv(StartCronjob)
	startCron, err := strconv.ParseBool(envCron)
	if err != nil || !startCron {
		return &dto.StartCronjobResponse{
			Code:     401,
			Error:    "Unauthorized",
			Messages: "Environment not allowed to trigger Cronjob",
		}, err
	}

	// create a separate context since it's not rely on parent grpc.Dial()
	longLiveCtx := context.Background()
	err = h.dataService.AddJob(longLiveCtx, req.Schedule, func() {
		// since we will have multiple daemonSet in nodes, need to make sure same cronjob
		// only running once at a time, here we use distrubted lock through Redis.
		lock := cache.ObtainLock(cache.CronjobLock, 2*time.Minute)
		if lock != nil {
			h.BatchingDownload(longLiveCtx, &dto.DownloadRequest{
				RewindLimit: 0,
				RateLimit:   3000,
				Types:       req.Types,
			})
		} else {
			log.Error("CronDownload: Redis distributed lock obtain failed.")
		}
	})

	if err != nil {
		return &dto.StartCronjobResponse{
			Code:     400,
			Error:    "Bad Request",
			Messages: fmt.Sprintf("Failed to start the schedule: %s with Types: %+v", req.Schedule, req.Types),
		}, err
	}

	return &dto.StartCronjobResponse{
		Code:     200,
		Messages: fmt.Sprintf("Successfully started the schedule: %s with Types: %+v", req.Schedule, req.Types),
	}, nil
}

// batching download all the historical stock data
func (h *handlerImpl) BatchingDownload(ctx context.Context, req *dto.DownloadRequest) {
	dailyCloseChan := make(chan *[]interface{})
	stakeConcentrationChan := make(chan *[]interface{})
	threePrimaryChan := make(chan *[]interface{})
	stockListChan := make(chan *[]interface{})

	for _, t := range req.Types {
		switch t {
		case dto.DailyClose:
			go h.generateJob(ctx, convert.TwseDailyClose, req, dailyCloseChan)
			go h.generateJob(ctx, convert.TpexDailyClose, req, dailyCloseChan)
		case dto.ThreePrimary:
			go h.generateJob(ctx, convert.TwseThreePrimary, req, threePrimaryChan)
			go h.generateJob(ctx, convert.TpexThreePrimary, req, threePrimaryChan)
		case dto.Concentration:
			go h.generateJob(ctx, convert.StakeConcentration, req, stakeConcentrationChan)
		case dto.StockList:
			go h.generateJob(ctx, convert.TwseStockList, req, stockListChan)
			go h.generateJob(ctx, convert.TpexStockList, req, stockListChan)
		}
	}

	go func() {
		for {
			select {
			// since its hard to predict how many records already been processed,
			// sync.WaitGroup hard to apply in this scenario, use timeout instead
			case <-time.After(8 * time.Hour):
				log.Warn("(BatchingDownload): timeout")
				return
			case <-ctx.Done():
				log.Warn("(BatchingDownload): context cancel")
				return
			case objs, ok := <-dailyCloseChan:
				if ok {
					h.dataService.DailyCloseThroughKafka(ctx, objs)
				}
			case objs, ok := <-threePrimaryChan:
				if ok {
					h.dataService.ThreePrimaryThroughKafka(ctx, objs)
				}
			case objs, ok := <-stockListChan:
				if ok {
					h.dataService.StockThroughKafka(ctx, objs)
				}
			case objs, ok := <-stakeConcentrationChan:
				if ok && len(*objs) > 0 {
					diff := []int32{0, 0, 0, 0, 0}
					var sc *entity.StakeConcentration
					for _, obj := range *objs {
						if val, ok := obj.(*entity.StakeConcentration); ok {
							switch val.HiddenField {
							case "1":
								sc = val
								diff[0] = int32(val.SumBuyShares - val.SumSellShares)
							case "2":
								diff[1] = int32(val.SumBuyShares - val.SumSellShares)
							case "3":
								diff[2] = int32(val.SumBuyShares - val.SumSellShares)
							case "4":
								diff[3] = int32(val.SumBuyShares - val.SumSellShares)
							case "6":
								diff[4] = int32(val.SumBuyShares - val.SumSellShares)
							}
						}
					}
					if sc == nil {
						log.Errorf("Failed to assign StakeConcentration: %+v", sc)
						continue
					}
					// refresh the concentration
					sc.Diff = diff
					err := h.dataService.StakeConcentrationThroughKafka(ctx, &[]interface{}{sc})
					if err != nil {
						log.Errorf("Error sending stake_concentration: %w", err)
					}
				}
			}
		}
	}()
}
