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
package services

import (
	"context"

	"github.com/samwang0723/stock-crawler/internal/app/crawler"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/app/graph"

	"github.com/hashicorp/go-multierror"
	"github.com/rs/zerolog"
	"golang.org/x/xerrors"
)

// Config encapsulates the settings for configuring the web-crawler service.
type CrawlerConfig struct {
	// An API for performing HTTP requests. If not specified,
	// crawler.DefaultHttpClient will be used instead.
	URLGetter crawler.URLGetter

	// The number of concurrent workers used for retrieving links.
	FetchWorkers int

	// The time between subsequent crawler passes.
	RateLimitInterval int

	// Proxy for preventing remote site's rate limiting
	Proxy *crawler.Proxy

	// The logger to use. If not defined an output-discarding logger will
	// be used instead.
	Logger zerolog.Logger
}

func (cfg *CrawlerConfig) validate() error {
	var err error
	if cfg.URLGetter == nil {
		cfg.URLGetter = crawler.DefaultHttpClient
	}

	if cfg.FetchWorkers <= 0 {
		err = multierror.Append(err, xerrors.Errorf("invalid value for fetch workers"))
	}
	if cfg.RateLimitInterval == 0 {
		err = multierror.Append(err, xerrors.Errorf("invalid value for rate limit interval"))
	}

	return err
}

func (s *serviceImpl) Crawl(ctx context.Context, linkIt graph.LinkIterator, interceptChan ...chan convert.InterceptData) (int, error) {
	return s.crawler.Crawl(ctx, linkIt, interceptChan...)
}
