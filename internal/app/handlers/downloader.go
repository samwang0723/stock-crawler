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
package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/samwang0723/stock-crawler/internal/app/crawler"
	"github.com/samwang0723/stock-crawler/internal/app/dto"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/app/graph"
	"github.com/samwang0723/stock-crawler/internal/cache"
	"github.com/samwang0723/stock-crawler/internal/helper"
)

const (
	cronLockPeriod     = 5
	cronTerminatedHour = 8
)

func (h *handlerImpl) CronDownload(ctx context.Context, req *dto.StartCronjobRequest) error {
	err := h.dataService.AddJob(ctx, req.Schedule, func() {
		// since we will have multiple daemonSet in nodes, need to make sure same cronjob
		// only running once at a time, here we use distrubted lock through Redis.
		if h.dataService.ObtainLock(ctx, cache.CronjobLock, cronLockPeriod*time.Minute) != nil {
			h.batchingDownload(ctx, 0, req.Types)
		}
	})

	if err != nil {
		return fmt.Errorf("handlers CronDownload(): %w", err)
	}

	return nil
}

func (h *handlerImpl) Download(ctx context.Context, req *dto.StartCronjobRequest) {
	h.batchingDownload(ctx, 0, req.Types)
}

// batching download all the historical stock data
func (h *handlerImpl) batchingDownload(ctx context.Context, rewind int32, types []convert.Source) {
	var links []*graph.Link

	interceptChan := make(chan convert.InterceptData)

	for _, strategy := range types {
		date := formatQueryDate(rewind, strategy)
		urls := h.generateURLs(ctx, date, strategy)

		for _, l := range urls {
			links = append(links, &graph.Link{
				URL:      l,
				Date:     date,
				Strategy: strategy,
			})
		}
	}

	go func() {
		for {
			select {
			// since its hard to predict how many records already been processed,
			// sync.WaitGroup hard to apply in this scenario, use timeout instead
			case <-time.After(cronTerminatedHour * time.Hour):
				h.logger.Warn().Msg("batching download: timeout")

				return
			case <-ctx.Done():
				h.logger.Warn().Msg("batching download: context cancel")

				return
			case obj, ok := <-interceptChan:
				if ok {
					h.processData(ctx, obj)
				}
			}
		}
	}()

	_, err := h.dataService.Crawl(ctx, &linkIterator{links: links}, interceptChan)

	if err != nil {
		h.logger.Error().Err(err).Msg("dataService crawl failed")
	}
}

func (h *handlerImpl) generateURLs(ctx context.Context, date string, source convert.Source) []string {
	var urls []string

	var err error

	//nolint:nolintlint, exhaustive
	switch source {
	case convert.StakeConcentration:
		// align the date format to be 20220107, but remains the query date as 2022-01-07
		unifiedDate := strings.ReplaceAll(date, "-", "")
		urls, err = h.dataService.ListCrawlingConcentrationURLs(ctx, unifiedDate)

		if err != nil {
			h.logger.Error().Err(err).Msg("dataService list crawling concentration urls failed")
		}
	case convert.TwseStockList, convert.TpexStockList:
		urls = append(urls, crawler.TypeLinkMapping[source.String()])
	default:
		urls = append(urls, fmt.Sprintf(crawler.TypeLinkMapping[source.String()], date))
	}

	return urls
}

func formatQueryDate(rewind int32, t convert.Source) string {
	var date string

	//nolint:nolintlint, exhaustive
	switch t {
	case convert.TwseDailyClose, convert.TwseThreePrimary:
		date = helper.GetDateFromOffset(rewind, helper.TwseDateFormat)
	case convert.TpexDailyClose, convert.TpexThreePrimary:
		date = helper.GetDateFromOffset(rewind, helper.TpexDateFormat)
	case convert.StakeConcentration:
		date = helper.GetDateFromOffset(rewind, helper.StakeConcentrationFormat)
	}

	return date
}

func (h *handlerImpl) processData(ctx context.Context, obj convert.InterceptData) {
	var err error

	switch obj.Type {
	case convert.TwseDailyClose, convert.TpexDailyClose:
		err = h.dataService.DailyCloseThroughKafka(ctx, obj.Data)
	case convert.TwseThreePrimary, convert.TpexThreePrimary:
		err = h.dataService.ThreePrimaryThroughKafka(ctx, obj.Data)
	case convert.TwseStockList, convert.TpexStockList:
		err = h.dataService.StockThroughKafka(ctx, obj.Data)
	case convert.StakeConcentration:
		err = h.dataService.StakeConcentrationThroughKafka(ctx, obj.Data)
	}

	if err != nil {
		h.logger.Error().Err(err).Msg(fmt.Sprintf("handler process data: %v", obj.Type))
	}
}
