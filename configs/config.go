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
	"fmt"
	"os"

	"github.com/samwang0723/stock-crawler/internal/helper"
	yaml "gopkg.in/yaml.v3"
)

type SystemConfig struct {
	RedisCache struct {
		Master        string   `yaml:"master"`
		SentinelAddrs []string `yaml:"sentinelAddrs"`
	} `yaml:"redis"`
	Kafka struct {
		Controller string   `yaml:"controller"`
		GroupID    string   `yaml:"groupId"`
		Brokers    []string `yaml:"brokers"`
		Topics     []string `yaml:"topics"`
	} `yaml:"kafka"`
	Server struct {
		Name         string `yaml:"name"`
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		MaxGoroutine int    `yaml:"maxGoroutine"`
		DNSLatency   int64  `yaml:"dnsLatency"`
	} `yaml:"server"`
	Crawler struct {
		FetchWorkers int   `yaml:"fetchWorkers"`
		RateLimit    int64 `yaml:"rateLimit"`
	} `yaml:"crawler"`
}

//nolint:nolintlint, gochecknoglobals
var instance SystemConfig

func Load(loc ...string) {
	var yamlFile string

	if len(loc) > 0 {
		yamlFile = loc[0]
	} else {
		yamlFile = fmt.Sprintf("./configs/config.%s.yaml", helper.GetCurrentEnv())
	}

	file, err := os.Open(yamlFile)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&instance)

	if err != nil {
		panic(err)
	}
}

func GetCurrentConfig() *SystemConfig {
	return &instance
}
