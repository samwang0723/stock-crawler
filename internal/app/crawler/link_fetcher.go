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
	"io"
	"net/http"

	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/app/pipeline"
	"github.com/samwang0723/stock-crawler/internal/helper"
	"golang.org/x/xerrors"

	"github.com/rs/zerolog"
)

var _ pipeline.Processor = (*linkFetcher)(nil)

// linkFetcher uses the configured URLGetter and Proxy to retrieve the remote content
// proxy could be nil if not necessary
type linkFetcher struct {
	urlGetter URLGetter
	proxy     *Proxy
	logger    *zerolog.Logger
}

func newLinkFetcher(urlGetter URLGetter, proxy *Proxy, logger *zerolog.Logger) *linkFetcher {
	return &linkFetcher{
		urlGetter: urlGetter,
		proxy:     proxy,
		logger:    logger,
	}
}

func (lf *linkFetcher) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	payload, ok := p.(*crawlerPayload)
	if !ok {
		return nil, xerrors.Errorf("linkFetcher.Process: failed, payload_type=%T;", p)
	}

	uri := payload.URL
	if lf.proxy != nil && payload.Strategy == convert.StakeConcentration {
		uri = lf.proxy.URI(payload.URL)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", uri, http.NoBody)
	if err != nil {
		return nil, xerrors.Errorf("linkFetcher.Process: failed, err=%w;", err)
	}

	req.Header = http.Header{
		"Content-Type": []string{"text/csv;charset=ms950"},
		// It is important to close the connection otherwise fd count will overhead
		"Connection": []string{"close"},
	}

	lf.logger.Info().Msgf("linkFetcher.Process: success, reason: download started; url=%s;", uri)

	resp, err := lf.urlGetter.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("urlGetter.Do(): %w", err)
	}

	// Skip payloads for invalid http status codes.
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, xerrors.Errorf("linkFetcher.Process: failed, http_status_code=%d;", resp.StatusCode)
	}

	// copy stream from response body, although it consumes memory but
	// better helps on concurrent handling in goroutine.
	_, err = io.Copy(&payload.RawContent, resp.Body)
	resp.Body.Close()

	if err != nil {
		return nil, xerrors.Errorf("linkFetcher.Process: failed, err=%w;", err)
	}

	//nolint:nolintlint, gomnd
	lf.logger.Info().Msgf("linkFetcher.Process: success, reason: download completed; size=%s; url=%s;",
		helper.GetReadableSize(payload.RawContent.Len(), 2), uri)

	return payload, nil
}
