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
	"flag"
	"os"
	"testing"

	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/helper"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestParseConcentration(t *testing.T) {
	t.Parallel()

	wrongDoc := "<html><body><table><tr><td>WRONG</td></tr></table></body></html>"
	correctDoc, _ := helper.ReadFromFile(".testfiles/concentration_fubon.html")

	b, _ := helper.EncodeBig5([]byte(correctDoc))
	tests := []struct {
		name    string
		content string
		stockID string
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
			shares:  []uint64{5610, 2180},
			price:   []float32{38.19, 38.09},
			stockID: "3704",
			hidden:  "0",
			date:    "20230110",
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
			t.Parallel()

			res := &parserImpl{
				result: &[]interface{}{},
			}
			res.SetStrategy(convert.StakeConcentration, "2023-01-10")
			res.Execute(*bytes.NewBuffer([]byte(tt.content)), "https://stockchannelnew.sinotrade.com.tw/z/zc/zco/zco_3704_1.djhtm")

			respLen := len(*res.result)
			if respLen > 0 != tt.want {
				t.Errorf("len(parser.result) = %v, want %v", respLen, tt.want)
			} else if respLen > 0 {
				c, ok := (*res.result)[0].(*entity.StakeConcentration)
				if ok && c.StockID != tt.stockID ||
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
