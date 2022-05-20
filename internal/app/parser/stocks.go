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

	"golang.org/x/net/html"
)

type htmlStrategy struct {
	converter convert.IConvert
	capacity  int
	source    convert.Source
}

func (s *htmlStrategy) Parse(in io.Reader, additional ...string) ([]interface{}, error) {
	var output []interface{}
	var records []string
	var isColumn, isBold, startParsing bool

	z := html.NewTokenizer(in)

	for {
		tt := z.Next()
		switch {
		case tt == html.StartTagToken:
			t := z.Token()
			isColumn = t.Data == "td"
			isBold = t.Data == "b"

			if t.Data == "tr" {
				if s.capacity == len(records) {
					// flush the temporary cache into output queue
					output = append(output, s.converter.Execute(&convert.ConvertData{
						Target:  s.source,
						RawData: records,
					}))
				}
				// reset the buffer to parse next row
				records = []string{}
			}
		case tt == html.TextToken:
			t := z.Token()
			content := strings.TrimSpace(t.Data)
			if len(content) == 0 {
				continue
			}
			switch {
			case isColumn:
				if startParsing {
					records = append(records, content)
				}
			case isBold:
				startParsing = content == "股票"
			}
		case tt == html.ErrorToken:
			if len(output) == 0 {
				return nil, NoParseResults
			}
			return output, nil
		}
	}
}
