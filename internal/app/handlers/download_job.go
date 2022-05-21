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
	"strings"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/samwang0723/stock-crawler/internal/app/crawler"
	"github.com/samwang0723/stock-crawler/internal/app/crawler/icrawler"
	"github.com/samwang0723/stock-crawler/internal/app/crawler/proxy"
	"github.com/samwang0723/stock-crawler/internal/app/dto"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/app/parser"
	"github.com/samwang0723/stock-crawler/internal/concurrent"
	"github.com/samwang0723/stock-crawler/internal/helper"
	log "github.com/samwang0723/stock-crawler/internal/logger"
)

// download job to run in workerpool
type downloadJob struct {
	date      string
	stockId   string
	respChan  chan *[]interface{}
	origin    convert.Source
	ctx       context.Context
	rateLimit int32
}

func (job *downloadJob) Do() error {
	var c icrawler.ICrawler
	p := parser.New()

	switch job.origin {
	case convert.TwseDailyClose:
		c = crawler.New(&proxy.Proxy{Type: proxy.DailyClose})
		url := fmt.Sprintf(icrawler.TwseDailyClose, job.date)
		c.AppendURL(url)
		p.SetStrategy(job.origin, job.date)

	case convert.TpexDailyClose:
		c = crawler.New(&proxy.Proxy{Type: proxy.DailyClose})
		url := fmt.Sprintf(icrawler.TpexDailyClose, job.date)
		c.AppendURL(url)

		p.SetStrategy(job.origin, job.date)

	case convert.TwseThreePrimary:
		c = crawler.New(&proxy.Proxy{Type: proxy.DailyClose})
		url := fmt.Sprintf(icrawler.TwseThreePrimary, job.date)
		c.AppendURL(url)
		p.SetStrategy(job.origin, job.date)

	case convert.TpexThreePrimary:
		c = crawler.New(&proxy.Proxy{Type: proxy.DailyClose})
		url := fmt.Sprintf(icrawler.TpexThreePrimary, job.date)
		c.AppendURL(url)
		p.SetStrategy(job.origin, job.date)

	case convert.TwseStockList:
		c = crawler.New(nil)
		c.AppendURL(icrawler.TWSEStocks)
		p.SetStrategy(job.origin)

	case convert.TpexStockList:
		c = crawler.New(nil)
		c.AppendURL(icrawler.TPEXStocks)
		p.SetStrategy(job.origin)

	case convert.StakeConcentration:
		c = crawler.New(&proxy.Proxy{Type: proxy.Concentration})
		// in order to get accurate data, we must query each page https://stockchannelnew.sinotrade.com.tw/z/zc/zco/zco_6598_6.djhtm
		// as the top 15 brokers may different from day to day and not possible to store all detailed daily data
		indexes := []int{1, 2, 3, 4, 6}
		for _, idx := range indexes {
			c.AppendURL(fmt.Sprintf(icrawler.ConcentrationDays, job.stockId, idx))
		}
		p.SetStrategy(job.origin, job.date)

	default:
		return fmt.Errorf("no recognized job source being specified: %s", job.origin)
	}

	// looping to download all URLs
	for {
		urls := c.GetURLs()
		if len(urls) == 0 {
			break
		}

		sourceURL, bytes, err := c.Fetch(job.ctx)
		if err != nil {
			sentry.CaptureException(err)
			return fmt.Errorf("(%s/%s): %+v", job.origin, job.date, err)
		}
		err = p.Execute(bytes, sourceURL)
		if err != nil {
			sentry.CaptureException(err)
			return fmt.Errorf("(%s/%s): %+v", job.origin, job.date, err)
		}
	}

	job.respChan <- p.Flush()

	// rate limit protection and context.cancel
	select {
	case <-time.After(time.Duration(job.rateLimit) * time.Millisecond):
	case <-job.ctx.Done():
		//log.Warn("(downloadJob) - context cancelled!")
	}

	return nil
}

func (h *handlerImpl) generateJob(ctx context.Context, origin convert.Source, req *dto.DownloadRequest, respChan chan *[]interface{}) {
	for _, d := range getParseDates(origin, req) {
		if err := h.queueJobs(ctx, origin, d, req.RateLimit, respChan); err != nil {
			log.Error(err)
			continue
		}
	}
	log.Debug("(BatchingDownload): all download jobs sent!")
}

func getParseDates(origin convert.Source, req *dto.DownloadRequest) []string {
	var dates []string
	if len(req.UTCTimestamp) > 0 {
		// parsing for specific date
		var date string
		switch origin {
		case convert.TwseDailyClose, convert.TwseThreePrimary:
			date = helper.GetDateFromUTC(req.UTCTimestamp, helper.TwseDateFormat)
		case convert.TpexDailyClose, convert.TpexThreePrimary:
			date = helper.GetDateFromUTC(req.UTCTimestamp, helper.TpexDateFormat)
		case convert.StakeConcentration:
			date = helper.GetDateFromUTC(req.UTCTimestamp, helper.StakeConcentrationFormat)
		}
		dates = append(dates, date)
	} else {
		// parsing for sequence dates using date rewind number
		for i := req.RewindLimit * -1; i <= 0; i++ {
			var date string
			switch origin {
			case convert.TwseDailyClose, convert.TwseThreePrimary:
				date = helper.GetDateFromOffset(i, helper.TwseDateFormat)
			case convert.TpexDailyClose, convert.TpexThreePrimary:
				date = helper.GetDateFromOffset(i, helper.TpexDateFormat)
			case convert.StakeConcentration:
				date = helper.GetDateFromOffset(i, helper.StakeConcentrationFormat)
			}
			dates = append(dates, date)
		}
	}
	return dates
}

func (h *handlerImpl) queueJobs(ctx context.Context, origin convert.Source, date string, rateLimit int32, respChan chan *[]interface{}) error {
	var jobs []*downloadJob
	if len(date) > 0 && origin == convert.StakeConcentration {
		// align the date format to be 20220107, but remains the query date as 2022-01-07
		//TODO: use redis to replace database query
		res, err := h.dataService.ListBackfillStakeConcentrationStockIds(ctx, strings.ReplaceAll(date, "-", ""))
		if err != nil {
			return fmt.Errorf("ListBackfillStakeConcentrationStockIds error: %+v", err)
		}
		for _, id := range res {
			job := &downloadJob{
				ctx:       ctx,
				date:      date,
				stockId:   id,
				respChan:  respChan,
				rateLimit: rateLimit,
				origin:    origin,
			}
			jobs = append(jobs, job)
		}
	} else {
		job := &downloadJob{
			ctx:       ctx,
			date:      date, // it could be empty date but doesn't matter, will identify use origin
			respChan:  respChan,
			rateLimit: rateLimit,
			origin:    origin,
		}
		jobs = append(jobs, job)
	}

	// batch put into job queue
	for _, job := range jobs {
		concurrent.JobQueue <- job
	}
	return nil
}
