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
package crawler

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/samwang0723/stock-crawler/internal/app/dto"
	"github.com/samwang0723/stock-crawler/internal/app/graph"
	"github.com/samwang0723/stock-crawler/internal/helper"
	log "github.com/samwang0723/stock-crawler/internal/logger"
)

const (
	TwseDailyClose    = "https://www.twse.com.tw/exchangeReport/MI_INDEX?response=csv&date=%s&type=ALLBUT0999"
	TwseThreePrimary  = "http://www.tse.com.tw/fund/T86?response=csv&date=%s&selectType=ALLBUT0999"
	OperatingDays     = "https://www.twse.com.tw/holidaySchedule/holidaySchedule?response=csv&queryYear=%d"
	TpexDailyClose    = "http://www.tpex.org.tw/web/stock/aftertrading/daily_close_quotes/stk_quote_download.php?l=zh-tw&d=%s&s=0,asc,0"
	TpexThreePrimary  = "https://www.tpex.org.tw/web/stock/3insti/daily_trade/3itrade_hedge_result.php?l=zh-tw&o=csv&se=EW&t=D&d=%s"
	TWSEStocks        = "https://isin.twse.com.tw/isin/C_public.jsp?strMode=2"
	TPEXStocks        = "https://isin.twse.com.tw/isin/C_public.jsp?strMode=4"
	ConcentrationDays = "https://stockchannelnew.sinotrade.com.tw/z/zc/zco/zco_%s_%d.djhtm"
)

type Crawler interface {
	Fetch(ctx context.Context, sink chan<- dto.Payload, errs chan<- error)
}

// URLGetter is implemented by objects that can perform HTTP GET requests.
type URLGetter interface {
	Do(req *http.Request) (*http.Response, error)
}

type crawlerImpl struct {
	client URLGetter
	proxy  *proxy
	it     graph.LinkIterator
}

func New(links []*graph.Link) Crawler {
	res := &crawlerImpl{
		client: &http.Client{
			Timeout: time.Second * 60,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		proxy: &proxy{Type: WebScraping},
		it:    &linkIterator{links: links},
	}
	return res
}

func (c *crawlerImpl) Fetch(ctx context.Context, sink chan<- dto.Payload, errs chan<- error) {
	for c.it.Next() {
		link := c.it.Link()
		uri := link.URL
		if link.UseProxy {
			uri = fmt.Sprintf("%s&url=%s", c.proxy.URI(), url.QueryEscape(link.URL))
		}

		req, err := http.NewRequest("GET", uri, nil)
		if err != nil {
			errs <- FetchError("NewRequest initialized failed", uri, err)
		}
		req.Header = http.Header{
			"Content-Type": []string{"text/csv;charset=ms950"},
			// It is important to close the connection otherwise fd count will overhead
			"Connection": []string{"close"},
		}
		req = req.WithContext(ctx)
		log.Debugf("download started: %s", uri)

		resp, err := c.client.Do(req)
		if err != nil {
			errs <- FetchError("client.Do", uri, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			errs <- FetchError(fmt.Sprintf("Status = %d", resp.StatusCode), uri, err)
		}

		p := payloadPool.Get().(*crawlerPayload)
		p.URL = link.URL
		p.RetrievedAt = time.Now()

		// copy stream from response body, although it consumes memory but
		// better helps on concurrent handling in goroutine.
		_, err = io.Copy(&p.RawContent, resp.Body)
		if err != nil {
			errs <- FetchError("Unable to io.Copy", link.URL, err)
		}
		log.Debugf("download completed (%s), URL: %s", helper.GetReadableSize(p.RawContent.Len(), 2), link.URL)

		sink <- p
	}
}
