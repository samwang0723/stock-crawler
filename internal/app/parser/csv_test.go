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
	"testing"

	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/helper"
)

func Test_parseCsv(t *testing.T) {
	wrongCsv, _ := helper.ReadFromFile("testfiles/wrong.csv")
	correctCsv, _ := helper.ReadFromFile("testfiles/correct.csv")
	threePrimaryCsv, _ := helper.ReadFromFile("testfiles/twse_threeprimary.csv")
	tpexThreePrimaryCsv, _ := helper.ReadFromFile("testfiles/tpex_threeprimary.csv")

	b1, _ := helper.EncodeBig5([]byte(correctCsv))
	b2, _ := helper.EncodeBig5([]byte(wrongCsv))
	b3, _ := helper.EncodeBig5([]byte(threePrimaryCsv))
	b4, _ := helper.EncodeBig5([]byte(tpexThreePrimaryCsv))

	tests := []struct {
		name    string
		content string
		want    int
		target  convert.Source
	}{
		{
			name:    "normal dailyclose csv",
			content: string(b1),
			want:    29,
			target:  convert.TwseDailyClose,
		},
		{
			name:    "wrong dailyclose csv",
			content: string(b2),
			want:    0,
			target:  convert.TwseDailyClose,
		},
		{
			name:    "twse three primary csv",
			content: string(b3),
			want:    22,
			target:  convert.TwseThreePrimary,
		},
		{
			name:    "tpex three primary csv",
			content: string(b4),
			want:    63,
			target:  convert.TpexThreePrimary,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res := &parserImpl{
				result: &[]interface{}{},
			}
			res.SetStrategy(tt.target, "20211130")
			res.Execute([]byte(tt.content))

			if got := len(*res.result); got != tt.want {
				t.Errorf("len(parser.result) = %v, want %v", got, tt.want)
			}
		})
	}

}
