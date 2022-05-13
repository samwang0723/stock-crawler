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
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/samwang0723/stock-crawler/internal/app/crawler/icrawler"
	"github.com/samwang0723/stock-crawler/internal/app/crawler/proxy"
	"github.com/samwang0723/stock-crawler/internal/helper"
	log "github.com/samwang0723/stock-crawler/internal/logger"
)

type crawlerImpl struct {
	client *http.Client
	proxy  *proxy.Proxy
	urls   []string
}

func New(p *proxy.Proxy) icrawler.ICrawler {
	res := &crawlerImpl{
		client: &http.Client{
			Timeout: time.Second * 60,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		proxy: p,
	}
	return res
}

func (c *crawlerImpl) AppendURL(url string) {
	c.urls = append(c.urls, url)
}

func (c *crawlerImpl) GetURLs() []string {
	return c.urls
}

func (c *crawlerImpl) Fetch(ctx context.Context) (string, []byte, error) {
	if len(c.urls) <= 0 {
		return "", nil, fmt.Errorf("no url to parse")
	}
	source := c.urls[0]
	uri := source
	if c.proxy != nil {
		uri = fmt.Sprintf("%s&url=%s", c.proxy.URI(), url.QueryEscape(source))
	}

	res, err := download(ctx, c.client, uri)
	if err != nil {
		return source, nil, err
	}

	// dequeue the first element
	if len(c.urls) > 0 {
		c.urls = c.urls[1:]
	}
	return source, res, nil
}

func download(ctx context.Context, client *http.Client, uri string) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("new fetch request initialize error: %v", err)
	}
	req.Header = http.Header{
		"Content-Type": []string{"text/csv;charset=ms950"},
		// It is important to close the connection otherwise fd count will overhead
		"Connection": []string{"close"},
	}
	req = req.WithContext(ctx)
	log.Debugf("download started: %s", uri)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch request error: %v, url: %s", err, uri)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch status error: %v, url: %s", resp.StatusCode, uri)
	}

	// copy stream from response body, although it consumes memory but
	// better helps on concurrent handling in goroutine.
	f, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fetch unable to read body: %v, url: %s", err, uri)
	}

	log.Debugf("download completed (%s), URL: %s", helper.GetReadableSize(len(f), 2), uri)
	return f, nil
}
