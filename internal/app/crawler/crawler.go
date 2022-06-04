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

	"github.com/samwang0723/stock-crawler/internal/app/crawler/icrawler"
	"github.com/samwang0723/stock-crawler/internal/app/crawler/proxy"
	"github.com/samwang0723/stock-crawler/internal/app/dto"
	"github.com/samwang0723/stock-crawler/internal/app/graph"
	"github.com/samwang0723/stock-crawler/internal/helper"
	log "github.com/samwang0723/stock-crawler/internal/logger"
)

type crawlerImpl struct {
	client icrawler.URLGetter
	proxy  *proxy.Proxy
	it     graph.LinkIterator
}

func New(links []*graph.Link) icrawler.ICrawler {
	res := &crawlerImpl{
		client: &http.Client{
			Timeout: time.Second * 60,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		proxy: *proxy.Proxy{Type: proxy.WebScraping},
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
