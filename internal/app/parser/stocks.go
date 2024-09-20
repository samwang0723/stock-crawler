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

//nolint:nolintlint, cyclop, gocognit
func (s *htmlStrategy) Parse(input io.Reader, _ ...string) ([]any, error) {
	var output []any

	var records []string

	var isColumn, isBold, startParsing bool

	tokenizer := html.NewTokenizer(input)

	for {
		next := tokenizer.Next()

		//nolint:nolintlint,exhaustive // ignore rest of the TokenType
		switch next {
		case html.StartTagToken:
			t := tokenizer.Token()
			isColumn = t.Data == "td"
			isBold = t.Data == "b"
		case html.TextToken:
			t := tokenizer.Token()
			content := strings.TrimSpace(t.Data)

			if content == "" {
				continue
			}

			switch {
			case isColumn:
				if startParsing {
					records = append(records, content)
				}
			case isBold:
				startParsing = tagToStart(content)
			}
		case html.ErrorToken:
			if len(output) == 0 {
				return nil, ErrNoParseResults
			}

			return output, nil
		case html.EndTagToken:
			t := tokenizer.Token()
			if t.Data == "tr" {
				if s.capacity <= len(records) {
					// flush the temporary cache into output queue
					res := s.converter.Execute(&convert.Data{
						Target:  s.source,
						RawData: records,
					})

					if res != nil {
						output = append(output, res)
					}
				}
				// reset the buffer to parse next row
				records = []string{}
			}
		}
	}
}

func tagToStart(content string) bool {
	return content == "股票" || content == "臺灣存託憑證(TDR)"
}
