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
	"fmt"
	"os"

	"github.com/samwang0723/stock-crawler/internal/helper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Log struct {
		Level string `yaml:"level"`
	} `yaml:"log"`
	RedisCache struct {
		Master        string   `yaml:"master"`
		SentinelAddrs []string `yaml:"sentinelAddrs"`
	} `yaml:"redis"`
	Kafka struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"kafka"`
	Server struct {
		Name string `yaml:"name"`
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	WorkerPool struct {
		MaxPoolSize  int `yaml:"maxPoolSize"`
		MaxQueueSize int `yaml:"maxQueueSize"`
	} `yaml:"workerpool"`
}

var (
	instance Config
)

func Load(loc ...string) {
	var yamlFile string
	if len(loc) > 0 {
		yamlFile = loc[0]
	} else {
		yamlFile = fmt.Sprintf("./configs/config.%s.yaml", helper.GetCurrentEnv())
	}
	f, err := os.Open(yamlFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&instance)
	if err != nil {
		panic(err)
	}
}

func GetCurrentConfig() *Config {
	return &instance
}
