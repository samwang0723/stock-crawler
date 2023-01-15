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
package config

import (
	"flag"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func Test_ConfigLoad(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want SystemConfig
	}{
		{
			name: "load configuration",
			want: SystemConfig{
				RedisCache: struct {
					Master        string   "yaml:\"master\""
					SentinelAddrs []string "yaml:\"sentinelAddrs\""
				}{
					Master:        "redis-master",
					SentinelAddrs: []string{"redis-sentinel-1:26379", "redis-sentinel-2:26380", "redis-sentinel-3:26381"},
				},
				Kafka: struct {
					Controller string "yaml:\"controller\""
				}{
					Controller: "kafka-1:9092",
				},
				Server: struct {
					Name         string "yaml:\"name\""
					Host         string "yaml:\"host\""
					Port         int    "yaml:\"port\""
					MaxGoroutine int    "yaml:\"maxGoroutine\""
					DNSLatency   int64  "yaml:\"dnsLatency\""
				}{
					Name:         "stock-crawler",
					Host:         "0.0.0.0",
					Port:         8086,
					MaxGoroutine: 20000,
					DNSLatency:   200,
				},
				Crawler: struct {
					FetchWorkers int   "yaml:\"fetchWorkers\""
					RateLimit    int64 "yaml:\"rateLimit\""
				}{
					FetchWorkers: 40,
					RateLimit:    500,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			Load("config.dev.yaml")
			cfg := GetCurrentConfig()
			if cmp.Equal(*cfg, tt.want) == false {
				t.Errorf("config.Load() = %+v, want %+v", *cfg, tt.want)
			}
		})
	}
}
