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
package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

const (
	DailyClose    = "DAILYCLOSE_PROXY"
	Concentration = "CONCENTRATION_PROXY"
)

type Proxy struct {
	Type          string
	RequireClient bool
}

func (p *Proxy) Client() *http.Client {
	proxyURL, _ := url.Parse("http://109.127.82.34:8080")
	return &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
}

func (p *Proxy) URI() string {
	token := os.Getenv(p.Type)
	return fmt.Sprintf("https://api.proxycrawl.com/?token=%s", token)
}
