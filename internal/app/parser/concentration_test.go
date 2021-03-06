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

	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/helper"
)

func Test_parseConcentration(t *testing.T) {
	wrongDoc := "<html><body><table><tr><td>WRONG</td></tr></table></body></html>"
	correctDoc, _ := helper.ReadFromFile("testfiles/concentration.html")

	b, _ := helper.EncodeBig5([]byte(correctDoc))
	tests := []struct {
		name    string
		content string
		stockId string
		hidden  string
		date    string
		shares  []uint64
		price   []float32
		want    bool
	}{
		{
			name:    "normal concentration html",
			content: string(b),
			want:    true,
			shares:  []uint64{12449, 40221},
			price:   []float32{63.45, 63.53},
			stockId: "6727",
			hidden:  "0",
			date:    "20211029",
		},
		{
			name:    "wrong concentration html",
			content: wrongDoc,
			want:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			//t.Parallel()
			res := &parserImpl{
				result: &[]interface{}{},
			}
			res.SetStrategy(convert.StakeConcentration, "2021-10-29")
			res.Execute([]byte(tt.content), "https://stockchannelnew.sinotrade.com.tw/z/zc/zco/zco_2330_1.djhtm")

			respLen := len(*res.result)
			if respLen > 0 != tt.want {
				t.Errorf("len(parser.result) = %v, want %v", respLen, tt.want)
			} else if respLen > 0 {
				c := (*res.result)[0].(*entity.StakeConcentration)
				if c.StockID != tt.stockId ||
					c.HiddenField != tt.hidden ||
					c.Date != tt.date ||
					c.SumBuyShares != tt.shares[0] ||
					c.SumSellShares != tt.shares[1] ||
					c.AvgBuyPrice != tt.price[0] ||
					c.AvgSellPrice != tt.price[1] {
					t.Errorf("response details = %+v is not fit the expectation", c)
				}
			}
		})
	}
}
