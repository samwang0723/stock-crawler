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
	"io"
	"strings"

	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/helper"
	"golang.org/x/net/html"
)

type concentrationStrategy struct {
	converter convert.IConvert
	date      string
	url       string
	capacity  int
}

func (s *concentrationStrategy) Parse(in io.Reader, additional ...string) ([]interface{}, error) {
	var records []string
	var output []interface{}
	var isColumn, isTitle, startParsing bool

	s.url = additional[0]
	l := len(s.url)
	if l > 7 && strings.HasSuffix(s.url, ".djhtm") {
		records = append(records, s.url[l-7:l-6])                      // 0: hidden field
		records = append(records, strings.ReplaceAll(s.date, "-", "")) // 1: date
	}

	z := html.NewTokenizer(in)

	for {
		tt := z.Next()
		switch {
		case tt == html.StartTagToken:
			t := z.Token()
			if t.Data == "tr" && len(t.Attr) == 1 {
				a := t.Attr[0]
				if a.Key == "id" && a.Val == "oScrollFoot" {
					startParsing = true
				}
			}

			isColumn = t.Data == "td"
			isTitle = t.Data == "title"
		case tt == html.TextToken:
			t := z.Token()
			content := strings.TrimSpace(t.Data)
			if len(content) == 0 {
				continue
			}

			switch {
			case isTitle:
				header := strings.Split(content, "-")
				if len(header) <= 1 {
					return nil, WrongConcentrationTitle
				}
				records = append(records, strings.TrimSpace(header[1])) // 2: stockId
			case isColumn:
				if startParsing {
					// make sure records are storing correct numbers
					n := strings.Replace(content, ",", "", -1)
					if helper.ToUint64(n) > 0 || helper.ToFloat32(n) > 0 {
						records = append(records, n) // 3,4,5,6: BuySum, SellSum, AvgBuy, AvgSell
					}

					// reached the row capacity limit, flush the cache data
					if s.capacity == len(records) {
						output = append(output, s.converter.Execute(&convert.ConvertData{
							RawData: records,
							Target:  convert.StakeConcentration,
						}))

						// truncate the temporary cache
						records = []string{}
					}
				}
			}
		case tt == html.ErrorToken:
			if len(output) == 0 {
				return nil, NoParseResults
			}
			return output, nil
		}
	}
}
