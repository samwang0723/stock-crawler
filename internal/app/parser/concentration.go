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

//
//nolint:nolintlint, cyclop, gocognit, cyclo
func (s *concentrationStrategy) Parse(
	input io.Reader,
	additional ...string,
) ([]any, error) {
	var records []string

	var output []any

	var isColumn, isTitle, startParsing bool

	s.url = additional[0]
	l := len(s.url)

	if l > 7 && strings.HasSuffix(s.url, ".djhtm") {
		// 0: hidden field, 1: date
		records = append(records, s.url[l-7:l-6], strings.ReplaceAll(s.date, "-", ""))
	}

	tokenizer := html.NewTokenizer(input)

	for {
		next := tokenizer.Next()

		switch {
		case next == html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "tr" && len(token.Attr) == 1 {
				a := token.Attr[0]
				if a.Key == "id" && a.Val == "oScrollFoot" {
					startParsing = true
				}
			}

			isColumn = token.Data == "td"
			isTitle = token.Data == "title"
		case next == html.TextToken:
			token := tokenizer.Token()
			content := strings.TrimSpace(token.Data)

			if content == "" {
				continue
			}

			switch {
			case isTitle:
				header := strings.Split(content, "-")
				if len(header) <= 1 {
					return nil, ErrWrongConcentrationTitle
				}

				records = append(records, strings.TrimSpace(header[1])) // 2: stockId
			case isColumn:
				if startParsing {
					// make sure records are storing correct numbers
					n := strings.ReplaceAll(content, ",", "")
					if helper.ToUint64(n) > 0 || helper.ToFloat32(n) > 0 {
						records = append(records, n) // 3,4,5,6: BuySum, SellSum, AvgBuy, AvgSell
					}

					// reached the row capacity limit, flush the cache data
					if s.capacity == len(records) {
						res := s.converter.Execute(&convert.Data{
							RawData: records,
							Target:  convert.StakeConcentration,
						})
						if res != nil {
							output = append(output, res)
						}

						// truncate the temporary cache
						records = []string{}
					}
				}
			}
		case next == html.ErrorToken:
			if len(output) == 0 {
				return nil, ErrNoParseResults
			}

			return output, nil
		}
	}
}
