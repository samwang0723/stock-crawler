// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package convert

import (
	"flag"
	"os"
	"testing"

	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/stretchr/testify/assert"
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

func TestConcentration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		val  *Data
		exp  *entity.StakeConcentration
	}{
		{
			name: "convert StakeConcentration",
			val: &Data{
				RawData: []string{"1", "20220525", "2330", "1000", "2000", "523", "518"},
				Target:  StakeConcentration,
			},
			exp: &entity.StakeConcentration{
				StockID:       "2330",
				Date:          "20220525",
				HiddenField:   "0",
				SumBuyShares:  1000,
				SumSellShares: 2000,
				AvgBuyPrice:   523,
				AvgSellPrice:  518,
			},
		},
		{
			name: "convert StakeConcentration hidden index 6 => 4",
			val: &Data{
				RawData: []string{"6", "20220525", "2330", "1000", "2000", "523", "518"},
				Target:  StakeConcentration,
			},
			exp: &entity.StakeConcentration{
				StockID:       "2330",
				Date:          "20220525",
				HiddenField:   "4",
				SumBuyShares:  1000,
				SumSellShares: 2000,
				AvgBuyPrice:   523,
				AvgSellPrice:  518,
			},
		},
		{
			name: "empty ConvertData",
			val:  nil,
			exp:  nil,
		},
		{
			name: "missing elements in ConvertData",
			val: &Data{
				RawData: []string{"1", "20220525", "2330", "1000"},
				Target:  StakeConcentration,
			},
			exp: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := Concentration()
			res := c.Execute(tt.val)
			if val, ok := res.(*entity.StakeConcentration); ok {
				assert.Equal(t, tt.exp, val)
			} else {
				t.Errorf("cannot convert StakeConcentration: %+v", res)
			}
		})
	}
}

func TestDailyClose(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		val  *Data
		exp  *entity.DailyClose
	}{
		{
			name: "convert TwseDailyClose",
			val: &Data{
				ParseDate: "20220525",
				RawData:   []string{"2330", "", "1,000", "1,000", "1,000", "100", "101", "1,005", "98", "-", "12", "", "", "", "", "", ""},
				Target:    TwseDailyClose,
			},
			exp: &entity.DailyClose{
				StockID:      "2330",
				Date:         "20220525",
				TradedShares: 1000,
				Transactions: 1000,
				Turnover:     1000,
				Open:         100,
				Close:        98,
				High:         101,
				Low:          1005,
				PriceDiff:    -12,
			},
		},
		{
			name: "convert TpexDailyClose",
			val: &Data{
				ParseDate: "20220525",
				RawData:   []string{"2330", "", "98", "-12", "100", "101", "105", "", "1,000", "1,000", "1,000", "", "", "", "", "", ""},
				Target:    TpexDailyClose,
			},
			exp: &entity.DailyClose{
				StockID:      "2330",
				Date:         "20220525",
				TradedShares: 1000,
				Transactions: 1000,
				Turnover:     1000,
				Open:         100,
				Close:        98,
				High:         101,
				Low:          105,
				PriceDiff:    -12,
			},
		},
		{
			name: "empty ConvertData",
			val:  nil,
			exp:  nil,
		},
		{
			name: "missing elements in ConvertData",
			val: &Data{
				ParseDate: "20220525",
				RawData:   []string{"", "", "1000", "1000", "1000"},
				Target:    TwseDailyClose,
			},
			exp: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := DailyClose()
			res := c.Execute(tt.val)
			if val, ok := res.(*entity.DailyClose); ok {
				assert.Equal(t, tt.exp, val)
			} else {
				t.Errorf("cannot convert DailyClose: %+v", res)
			}
		})
	}
}

func TestThreePrimary(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		val  *Data
		exp  *entity.ThreePrimary
	}{
		{
			name: "convert TpexThreePrimary",
			val: &Data{
				ParseDate: "20220525",
				RawData:   []string{"2330", "", "", "", "", "", "", "", "", "", "500", "", "", "1,300", "", "", "1,600", "", "", "2,000"},
				Target:    TpexThreePrimary,
			},
			exp: &entity.ThreePrimary{
				StockID:            "2330",
				Date:               "20220525",
				ForeignTradeShares: 500,
				TrustTradeShares:   1300,
				DealerTradeShares:  1600,
				HedgingTradeShares: 2000,
			},
		},
		{
			name: "convert TwseThreePrimary",
			val: &Data{
				ParseDate: "20220525",
				RawData:   []string{"2330", "", "", "", "500", "", "", "", "", "", "1,300", "", "", "", "1,600", "", "", "2,000", "", ""},
				Target:    TwseThreePrimary,
			},
			exp: &entity.ThreePrimary{
				StockID:            "2330",
				Date:               "20220525",
				ForeignTradeShares: 500,
				TrustTradeShares:   1300,
				DealerTradeShares:  1600,
				HedgingTradeShares: 2000,
			},
		},
		{
			name: "empty ConvertData",
			val:  nil,
			exp:  nil,
		},
		{
			name: "missing elements in ConvertData",
			val: &Data{
				ParseDate: "20220525",
				RawData:   []string{"", "", "1000", "1000", "1000"},
				Target:    TwseThreePrimary,
			},
			exp: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := ThreePrimary()
			res := c.Execute(tt.val)
			if val, ok := res.(*entity.ThreePrimary); ok {
				assert.Equal(t, tt.exp, val)
			} else {
				t.Errorf("cannot convert ThreePrimary: %+v", res)
			}
		})
	}
}

func TestStock(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		val  *Data
		exp  *entity.Stock
	}{
		{
			name: "convert Stock",
			val: &Data{
				RawData: []string{"2330　ABC", "", "", "上櫃", "XXX", "", ""},
				Target:  TwseStockList,
			},
			exp: &entity.Stock{
				StockID:  "2330",
				Name:     "ABC",
				Country:  "TW",
				Category: "XXX",
				Market:   "otc",
			},
		},
		{
			name: "empty ConvertData",
			val:  nil,
			exp:  nil,
		},
		{
			name: "missing elements in ConvertData",
			val: &Data{
				RawData: []string{"", "", "1000", "1000"},
				Target:  TwseStockList,
			},
			exp: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := Stock()
			res := c.Execute(tt.val)
			if val, ok := res.(*entity.Stock); ok {
				assert.Equal(t, tt.exp, val)
			} else {
				t.Errorf("cannot convert Stock: %+v", res)
			}
		})
	}
}
