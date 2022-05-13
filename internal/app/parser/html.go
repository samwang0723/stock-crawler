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

package parser

import (
	"fmt"
	"io"
	"strings"

	"github.com/go-errors/errors"
	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"golang.org/x/net/html"
)

func getInnerContent(z *html.Tokenizer, start *bool) (string, error) {
	switch z.Next() {
	case html.TextToken:
		text := (string)(z.Text())
		t := strings.TrimSpace(text)
		if len(t) > 0 {
			return t, nil
		}
	case html.StartTagToken:
		tag, _ := z.TagName()
		n := (string)(tag)
		if n == "b" {
			t, err := getInnerContent(z, start)
			if err != nil {
				break
			}
			if t == "股票" {
				*start = true
			} else {
				*start = false
			}
		}
	}
	return "", errors.New("invalid text content")
}

func (p *parserImpl) parseHtml(config Config, in io.Reader) error {
	originLen := len(*p.result)
	updatedLen := originLen

	z := html.NewTokenizer(in)
	buffer := []string{}
	start := false

	// while have not hit the </body> tag
	for z.Token().Data != "body" {
		tt := z.Next()
		if tt == html.StartTagToken {
			t := z.Token()
			switch t.Data {
			case "tr":
				if config.Capacity == len(buffer) {
					s := strings.Split(buffer[0], "　")
					stock := &entity.Stock{
						StockID:  strings.TrimSpace(s[0]),
						Name:     strings.TrimSpace(s[1]),
						Country:  "TW",
						Category: strings.TrimSpace(buffer[4]),
					}
					*p.result = append(*p.result, stock)
					updatedLen++
				}

				// reset the buffer to parse new raw
				buffer = []string{}
			case "td":
				t, err := getInnerContent(z, &start)
				if start && err == nil {
					buffer = append(buffer, t)
				}
			}
		} else if tt == html.ErrorToken {
			break
		}
	}

	if updatedLen <= originLen {
		return fmt.Errorf("empty parsing results")
	}

	return nil
}
