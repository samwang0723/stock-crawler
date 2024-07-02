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
	"bytes"
	"context"
	"flag"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/app/graph"
	"github.com/samwang0723/stock-crawler/internal/helper"
	"go.uber.org/goleak"
)

type testLinkIterator struct {
	links    []*graph.Link
	curIndex int
	mu       sync.RWMutex
}

func (i *testLinkIterator) Next() bool {
	if i.curIndex >= len(i.links) {
		return false
	}
	i.curIndex++

	return true
}

func (i *testLinkIterator) Error() error {
	return nil
}

func (i *testLinkIterator) Link() *graph.Link {
	i.mu.RLock()
	link := new(graph.Link)
	*link = *i.links[i.curIndex-1]
	i.mu.RUnlock()

	return link
}

// mock HTTP client
type mockSuccessHTTPClient struct{}

func (ms *mockSuccessHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	correctDoc, err := helper.ReadFromFile("../parser/.testfiles/stocks.html")
	correctBytes, _ := helper.EncodeBig5([]byte(correctDoc))

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(correctBytes)),
	}, err
}

type mockErrorHTTPClient struct{}

func (mf *mockErrorHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	return nil, os.ErrInvalid
}

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestCrawl(t *testing.T) {
	t.Parallel()

	type args struct {
		mockClient URLGetter
		link       *graph.Link
	}

	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			name: "regular http fetch",
			args: args{
				mockClient: &mockSuccessHTTPClient{},
				link: &graph.Link{
					URL:      "http://www.google.com",
					Date:     "20220801",
					Strategy: convert.TwseStockList,
				},
			},
			wantErr: false,
		},
		{
			name: "error fetching from server",
			args: args{
				mockClient: &mockErrorHTTPClient{},
				link: &graph.Link{
					URL:      "https://www.yahoo.com",
					Date:     "20210723",
					Strategy: convert.TwseStockList,
				},
			},
			wantErr: true,
		},
	}

	logger := log.With().Str("test", "crawler").Logger()

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := New(Config{
				URLGetter:         tt.args.mockClient,
				FetchWorkers:      2,
				RateLimitInterval: 1000,
				Logger:            &logger,
			})
			_, err := c.Crawl(context.TODO(), &testLinkIterator{links: []*graph.Link{
				tt.args.link,
			}})
			if (err != nil) != tt.wantErr {
				t.Errorf("Crawl() = %v, want %v, err: %v", err != nil, tt.wantErr, err)
			}
		})
	}
}
