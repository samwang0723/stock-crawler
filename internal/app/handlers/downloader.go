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

func (h *handlerImpl) CronDownload(ctx context.Context, req *dto.StartCronjobRequest) error {
	return h.dataService.AddJob(ctx, req.Schedule, func() {
		// since we will have multiple daemonSet in nodes, need to make sure same cronjob
		// only running once at a time, here we use distrubted lock through Redis.
		if h.dataService.ObtainLock(ctx, cache.CronjobLock, 5*time.Minute) != nil {
			h.batchingDownload(ctx, 0, req.Types)
		}
	})
}

func (h *handlerImpl) Download(ctx context.Context, req *dto.StartCronjobRequest) {
	h.batchingDownload(ctx, 0, req.Types)
}

// batching download all the historical stock data
func (h *handlerImpl) batchingDownload(ctx context.Context, rewind int32, types []convert.Source) {
	var links []*graph.Link

	for _, t := range types {
		d := formatQueryDate(rewind, t)
		urls := h.generateURLs(ctx, d, t)

		for _, l := range urls {
			links = append(links, &graph.Link{
				URL:      l,
				Date:     d,
				Strategy: t,
			})
		}
	}
	h.dataService.Crawl(ctx, &linkIterator{links: links})
}

func (h *handlerImpl) generateURLs(ctx context.Context, d string, t convert.Source) []string {
	var urls []string
	switch t {
	case convert.StakeConcentration:
		// align the date format to be 20220107, but remains the query date as 2022-01-07
		unifiedDate := strings.ReplaceAll(d, "-", "")
		urls, _ = h.dataService.ListCrawlingConcentrationURLs(ctx, unifiedDate)
	case convert.TwseStockList, convert.TpexStockList:
		urls = append(urls, crawler.TypeLinkMapping[t.String()])
	default:
		urls = append(urls, fmt.Sprintf(crawler.TypeLinkMapping[t.String()], d))
	}
	return urls
}

func formatQueryDate(rewind int32, t convert.Source) string {
	var date string
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
