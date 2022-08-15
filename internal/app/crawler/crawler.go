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
package crawler

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"

	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/app/graph"
	"github.com/samwang0723/stock-crawler/internal/app/parser"
	"github.com/samwang0723/stock-crawler/internal/app/pipeline"
	"github.com/sirupsen/logrus"
)

const (
	TwseDailyClose    = "https://www.twse.com.tw/exchangeReport/MI_INDEX?response=csv&date=%s&type=ALLBUT0999"
	TwseThreePrimary  = "http://www.tse.com.tw/fund/T86?response=csv&date=%s&selectType=ALLBUT0999"
	TpexDailyClose    = "http://www.tpex.org.tw/web/stock/aftertrading/daily_close_quotes/stk_quote_download.php?l=zh-tw&d=%s&s=0,asc,0"
	TpexThreePrimary  = "https://www.tpex.org.tw/web/stock/3insti/daily_trade/3itrade_hedge_result.php?l=zh-tw&o=csv&se=EW&t=D&d=%s"
	TWSEStocks        = "https://isin.twse.com.tw/isin/C_public.jsp?strMode=2"
	TPEXStocks        = "https://isin.twse.com.tw/isin/C_public.jsp?strMode=4"
	ConcentrationDays = "https://stockchannelnew.sinotrade.com.tw/z/zc/zco/zco_%s_%d.djhtm"
)

var (
	DefaultHttpClient = &http.Client{
		Timeout: time.Second * 60,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	TypeLinkMapping = map[string]string{
		convert.TpexStockList.String():      TPEXStocks,
		convert.TwseStockList.String():      TWSEStocks,
		convert.TwseDailyClose.String():     TwseDailyClose,
		convert.TpexDailyClose.String():     TpexDailyClose,
		convert.TwseThreePrimary.String():   TwseThreePrimary,
		convert.TpexThreePrimary.String():   TpexThreePrimary,
		convert.StakeConcentration.String(): ConcentrationDays,
	}
)

type Crawler interface {
	Crawl(ctx context.Context, linkIt graph.LinkIterator, interceptChan ...chan convert.InterceptData) (int, error)
}

// URLGetter is implemented by objects that can perform HTTP GET requests.
type URLGetter interface {
	Do(req *http.Request) (*http.Response, error)
}

type Config struct {
	// A URLGetter instance to fetch links.
	URLGetter URLGetter

	// A Proxy instance for avoiding remote rate limiting
	Proxy *Proxy

	// The number of concurrent workers used for retrieving links.
	FetchWorkers int

	// Rate limit interval to prevent remote site blocking
	RateLimitInterval int

	Logger *logrus.Entry
}

// crawlerImpl implements a stock information crawling pipeline consisting of following stages:
//
// - Given an URL, retrieve content from remote server
// - Extract useful trading information from retrieved pages
type crawlerImpl struct {
	cfg       Config
	extractor *textExtractor
	pipe      *pipeline.Pipeline
}

func New(cfg Config) Crawler {
	extractor := newTextExtractor(
		parser.New(parser.Config{Logger: cfg.Logger}),
		cfg.Logger,
	)
	return &crawlerImpl{
		cfg:       cfg,
		extractor: extractor,
		pipe:      assembleCrawlerPipeline(cfg, extractor),
	}
}

// assembleCrawlerPipeline creates the various stages of a crawler pipeline
// using the options in cfg and assembles them into a pipeline instance.
func assembleCrawlerPipeline(cfg Config, extractor *textExtractor) *pipeline.Pipeline {
	return pipeline.New(
		pipeline.DynamicWorkerPool(
			newLinkFetcher(cfg.URLGetter, cfg.Proxy, cfg.Logger),
			cfg.FetchWorkers,
		),
		pipeline.FIFO(extractor),
	)
}

// Crawl iterates linkIt and sends each link through the crawler pipeline
// returning the total count of links that went through the pipeline. Calls to
// Crawl block until the link iterator is exhausted, an error occurs or the
// context is cancelled.
func (c *crawlerImpl) Crawl(ctx context.Context, linkIt graph.LinkIterator, interceptChan ...chan convert.InterceptData) (int, error) {
	sink := new(countingSink)
	if len(interceptChan) == 1 {
		c.extractor.InterceptData(ctx, interceptChan[0])
	}
	err := c.pipe.Process(ctx, &linkSource{linkIt: linkIt}, sink)

	return sink.getCount(), err
}

type linkSource struct {
	linkIt graph.LinkIterator
}

func (ls *linkSource) Error() error              { return ls.linkIt.Error() }
func (ls *linkSource) Next(context.Context) bool { return ls.linkIt.Next() }
func (ls *linkSource) Payload() pipeline.Payload {
	link := ls.linkIt.Link()
	p := payloadPool.Get().(*crawlerPayload)

	p.URL = link.URL
	p.Strategy = link.Strategy
	p.Date = link.Date
	p.RetrievedAt = time.Now()

	return p
}

// countingSink for calculate total parsed records
type countingSink struct {
	count int
}

func (s *countingSink) Consume(_ context.Context, p pipeline.Payload) error {
	s.count++
	return nil
}

func (s *countingSink) getCount() int {
	return s.count
}
