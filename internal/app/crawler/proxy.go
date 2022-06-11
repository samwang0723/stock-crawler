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
	"fmt"
	"net/url"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

const (
	WebScraping    = "WEB_SCRAPING"
	WebScrapingUrl = "https://api.webscrapingapi.com/v1?api_key=%s"
	ProxyCrawl     = "PROXY_CRAWL"
	ProxyCrawlUrl  = "https://api.proxycrawl.com/?token=%s"
)

type Proxy struct {
	Type string
}

func (p *Proxy) URI(source string) string {
	var prefix string
	switch p.Type {
	case ProxyCrawl:
		prefix = ProxyCrawlUrl
	case WebScraping:
		prefix = WebScrapingUrl
	}
	token := os.Getenv(p.Type)
	return fmt.Sprintf("%s&url=%s", fmt.Sprintf(prefix, token), url.QueryEscape(source))
}
