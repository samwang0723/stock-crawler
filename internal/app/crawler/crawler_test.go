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
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/samwang0723/stock-crawler/internal/helper"
	log "github.com/samwang0723/stock-crawler/internal/logger"
	logtest "github.com/samwang0723/stock-crawler/internal/logger/structured"
)

func setup() {
	logger := logtest.NullLogger()
	log.Initialize(logger)
}

func shutdown() {
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func Test_Fetch(t *testing.T) {
	setup()
	tests := []struct {
		server *httptest.Server
		name   string
		want   bool
	}{
		{
			name: "Regular http fetch",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write(helper.String2Bytes("Success"))
			})),
			want: false,
		},
		{
			name: "error fetching from server",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(500)
			})),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer tt.server.Close()
			c := &crawlerImpl{
				urls:   []string{tt.server.URL},
				client: tt.server.Client(),
			}
			_, _, err := c.Fetch(context.TODO())
			if (err != nil) != tt.want {
				t.Errorf("Fetch() = %v, want %v", err != nil, tt.want)
			}
		})
	}
}
