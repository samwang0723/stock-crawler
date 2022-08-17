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
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/app/graph"
	"github.com/samwang0723/stock-crawler/internal/helper"

	"go.uber.org/goleak"
)

type testLinkIterator struct {
	mu       sync.RWMutex
	links    []*graph.Link
	curIndex int
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
		mockServer *httptest.Server
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Regular http fetch",
			args: args{
				mockServer: httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.WriteHeader(http.StatusOK)
						w.Write(helper.String2Bytes("Success"))
					}),
				),
			},
			wantErr: false,
		},
		{
			name: "error fetching from server",
			args: args{
				mockServer: httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						w.WriteHeader(500)
					}),
				),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			defer tt.args.mockServer.Close()

			c := New(Config{
				URLGetter:         tt.args.mockServer.Client(),
				FetchWorkers:      2,
				RateLimitInterval: 1000,
			})
			_, err := c.Crawl(context.TODO(), &testLinkIterator{links: []*graph.Link{
				{
					URL:      tt.args.mockServer.URL,
					Date:     "20220801",
					Strategy: convert.TwseDailyClose,
				},
			}})
			if (err != nil) != tt.wantErr {
				t.Errorf("Crawl() = %v, want %v", err != nil, tt.wantErr)
			}
		})
	}
}
