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

	"github.com/rs/zerolog"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/app/graph"
	"github.com/samwang0723/stock-crawler/internal/app/pipeline"
	"golang.org/x/xerrors"
)

//nolint:nolintlint, lll
const (
	TwseDailyClose   = "https://www.twse.com.tw/exchangeReport/MI_INDEX?response=csv&date=%s&type=ALLBUT0999"
	TwseThreePrimary = "https://www.twse.com.tw/rwd/zh/fund/T86?response=csv&date=%s&selectType=ALLBUT0999"
	TpexDailyClose   = "https://www.tpex.org.tw/web/stock/aftertrading/daily_close_quotes/stk_quote_download.php?l=zh-tw&d=%s&s=0,asc,0"
	TpexThreePrimary = "https://www.tpex.org.tw/web/stock/3insti/daily_trade/3itrade_hedge_result.php?l=zh-tw&o=csv&se=EW&t=D&d=%s"
	TWSEStocks       = "https://isin.twse.com.tw/isin/C_public.jsp?strMode=2"
	TPEXStocks       = "https://isin.twse.com.tw/isin/C_public.jsp?strMode=4"
	// backup: stockchannelnew.sinotrade.com.tw
	ConcentrationDays = "https://fubon-ebrokerdj.fbs.com.tw/z/zc/zco/zco_%s_%d.djhtm"

	defaultHTTPTimeout = 60 * time.Second
)

//nolint:nolintlint, gochecknoglobals, gosec
var (
	DefaultHTTPClient = &http.Client{
		Timeout: defaultHTTPTimeout,
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
	RateLimitInterval int64

	Logger *zerolog.Logger
}

// crawlerImpl implements a stock information crawling pipeline consisting of following stages:
//
// - Given an URL, retrieve content from remote server
// - Extract useful trading information from retrieved pages
type crawlerImpl struct {
	cfg       Config
	broadcast *broadcastor
	pipe      *pipeline.Pipeline
}

func New(cfg Config) Crawler {
	return &crawlerImpl{
		cfg:       cfg,
		broadcast: newBroadcastor(),
	}
}

// assembleCrawlerPipeline creates the various stages of a crawler pipeline
// using the options in cfg and assembles them into a pipeline instance.
func assembleCrawlerPipeline(cfg Config, broadcastor *broadcastor) *pipeline.Pipeline {
	return pipeline.New(
		pipeline.DynamicWorkerPool(
			newLinkFetcher(cfg.URLGetter, cfg.Proxy, cfg.Logger),
			cfg.FetchWorkers,
			time.Duration(cfg.RateLimitInterval)*time.Millisecond,
		),
		pipeline.FIFO(newTextExtractor(cfg)),
		pipeline.Broadcast(broadcastor),
	)
}

// Crawl iterates linkIt and sends each link through the crawler pipeline
// returning the total count of links that went through the pipeline. Calls to
// Crawl block until the link iterator is exhausted, an error occurs or the
// context is cancelled.
func (c *crawlerImpl) Crawl(
	ctx context.Context,
	linkIt graph.LinkIterator,
	interceptChan ...chan convert.InterceptData,
) (int, error) {
	// reconstruct pipeline every time as previous pipeline may be terminated
	c.pipe = assembleCrawlerPipeline(c.cfg, c.broadcast)

	sink := new(countingSink)

	if len(interceptChan) == 1 {
		c.broadcast.InterceptData(ctx, interceptChan[0])
	}

	err := c.pipe.Process(ctx, &linkSource{linkIt: linkIt}, sink)

	return sink.getCount(), err
}

type linkSource struct {
	linkIt graph.LinkIterator
}

func (ls *linkSource) Error() error {
	if err := ls.linkIt.Error(); err != nil {
		return xerrors.Errorf("linkSource error: %w", err)
	}

	return nil
}

func (ls *linkSource) Next(context.Context) bool { return ls.linkIt.Next() }

func (ls *linkSource) Payload() pipeline.Payload {
	link := ls.linkIt.Link()

	payload, ok := payloadPool.Get().(*crawlerPayload)
	if !ok {
		return nil
	}

	payload.URL = link.URL
	payload.Strategy = link.Strategy
	payload.Date = link.Date
	payload.RetrievedAt = time.Now()

	return payload
}

// countingSink for calculate total parsed records
type countingSink struct {
	count int
}

func (s *countingSink) Consume(_ context.Context, _ pipeline.Payload) error {
	s.count++

	return nil
}

func (s *countingSink) getCount() int {
	return s.count
}
