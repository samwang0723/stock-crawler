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
	"bytes"
	"testing"

	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/helper"
)

func Test_parseHtml(t *testing.T) {
	wrongDoc := "<html><body></body></html>"
	correctDoc, err := helper.ReadFromFile("testfiles/stocks.html")
	if err != nil {
		t.Errorf("failed to load html test file: %s", err)
	}
	correctBytes, _ := helper.EncodeBig5([]byte(correctDoc))

	tests := []struct {
		err     error
		name    string
		content string
		want    int
	}{
		{
			name:    "normal stock list html",
			content: string(correctBytes),
			want:    4,
			err:     nil,
		},
		{
			name:    "wrong stock list html",
			content: wrongDoc,
			want:    0,
			err:     NoParseResults,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := &parserImpl{
				result: &[]interface{}{},
			}
			res.SetStrategy(convert.TwseStockList)
			err := res.Execute(*bytes.NewBuffer([]byte(tt.content)))

			if got := len(*res.result); got != tt.want || err != tt.err {
				t.Errorf("len(parser.result) = %v, want %v", got, tt.want)
			}
		})
	}

}
