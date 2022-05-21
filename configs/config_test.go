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
package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_ConfigLoad(t *testing.T) {
	tests := []struct {
		name string
		want Config
	}{
		{
			name: "load configuration",
			want: Config{
				Log: struct {
					Level string "yaml:\"level\""
				}{
					Level: "DEBUG",
				},
				RedisCache: struct {
					Master string   "yaml:\"master\""
					Ports  []string "yaml:\"ports\""
				}{
					Master: "redis-master",
					Ports:  []string{":26379", ":26380", ":26381"},
				},
				Kafka: struct {
					Host string "yaml:\"host\""
					Port int    "yaml:\"port\""
				}{
					Host: "kafka",
					Port: 9092,
				},
				Server: struct {
					Name string "yaml:\"name\""
					Host string "yaml:\"host\""
					Port int    "yaml:\"port\""
				}{
					Name: "stock-crawler",
					Host: "0.0.0.0",
					Port: 8086,
				},
				WorkerPool: struct {
					MaxPoolSize  int "yaml:\"maxPoolSize\""
					MaxQueueSize int "yaml:\"maxQueueSize\""
				}{
					MaxPoolSize:  20,
					MaxQueueSize: 40,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Load("config.dev.yaml")
			cfg := GetCurrentConfig()
			if cmp.Equal(*cfg, tt.want) == false {
				t.Errorf("config.Load() = %+v, want %+v", *cfg, tt.want)
			}
		})
	}
}
