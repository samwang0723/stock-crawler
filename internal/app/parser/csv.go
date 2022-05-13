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
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/helper"
)

func (p *parserImpl) parseCsv(config Config, in io.Reader) error {
	if config.ParseDay == nil {
		return fmt.Errorf("parse day missing\n")
	}

	originLen := len(*p.result)
	updatedLen := originLen

	reader := csv.NewReader(in)
	reader.Comma = ','
	reader.FieldsPerRecord = -1

	//override to standarize date string (20211123)
	date := helper.UnifiedDateFormatToTwse(*config.ParseDay)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if len(record) == 0 || config.Capacity > len(record) {
			continue
		}

		// make sure only parse recognized stock_id
		record[0] = strings.TrimSpace(record[0])
		if len(record[0]) > 0 && len(record[0]) < 6 && helper.IsInteger(record[0][0:2]) {
			switch config.Type {
			case TwseDailyClose:
				*p.result = append(*p.result, twseToEntity(date, record))
				updatedLen++
			case TwseThreePrimary:
				*p.result = append(*p.result, twseThreePrimaryToEntity(date, record))
				updatedLen++
			case TpexDailyClose:
				*p.result = append(*p.result, tpexToEntity(date, record))
				updatedLen++
			case TpexThreePrimary:
				*p.result = append(*p.result, tpexThreePrimaryToEntity(date, record))
				updatedLen++
			}
		}
	}
	if updatedLen <= originLen {
		return fmt.Errorf("empty parsing results\n")
	}

	return nil
}

func twseThreePrimaryToEntity(day string, data []string) *entity.ThreePrimary {
	threePrimary := &entity.ThreePrimary{
		StockID:            data[0],
		Date:               day,
		ForeignTradeShares: helper.ToInt64(strings.Replace(data[4], ",", "", -1)),
		TrustTradeShares:   helper.ToInt64(strings.Replace(data[10], ",", "", -1)),
		DealerTradeShares:  helper.ToInt64(strings.Replace(data[14], ",", "", -1)),
		HedgingTradeShares: helper.ToInt64(strings.Replace(data[17], ",", "", -1)),
	}
	return threePrimary
}

func twseToEntity(day string, data []string) *entity.DailyClose {
	dailyClose := &entity.DailyClose{
		StockID:      data[0],
		Date:         day,
		TradedShares: helper.ToUint64(strings.Replace(data[2], ",", "", -1)),
		Transactions: helper.ToUint64(strings.Replace(data[3], ",", "", -1)),
		Turnover:     helper.ToUint64(strings.Replace(data[4], ",", "", -1)),
		Open:         helper.ToFloat32(data[5]),
		High:         helper.ToFloat32(data[6]),
		Low:          helper.ToFloat32(data[7]),
		Close:        helper.ToFloat32(data[8]),
		PriceDiff:    helper.ToFloat32(fmt.Sprintf("%s%s", data[9], data[10])),
	}
	return dailyClose
}

func tpexThreePrimaryToEntity(day string, data []string) *entity.ThreePrimary {
	threePrimary := &entity.ThreePrimary{
		StockID:            data[0],
		Date:               day,
		ForeignTradeShares: helper.ToInt64(strings.Replace(data[10], ",", "", -1)),
		TrustTradeShares:   helper.ToInt64(strings.Replace(data[13], ",", "", -1)),
		DealerTradeShares:  helper.ToInt64(strings.Replace(data[16], ",", "", -1)),
		HedgingTradeShares: helper.ToInt64(strings.Replace(data[19], ",", "", -1)),
	}
	return threePrimary
}

func tpexToEntity(day string, data []string) *entity.DailyClose {
	dailyClose := &entity.DailyClose{
		StockID:      data[0],
		Date:         day,
		TradedShares: helper.ToUint64(strings.Replace(data[8], ",", "", -1)),
		Transactions: helper.ToUint64(strings.Replace(data[10], ",", "", -1)),
		Turnover:     helper.ToUint64(strings.Replace(data[9], ",", "", -1)),
		Open:         helper.ToFloat32(data[4]),
		High:         helper.ToFloat32(data[5]),
		Low:          helper.ToFloat32(data[6]),
		Close:        helper.ToFloat32(data[2]),
		PriceDiff:    helper.ToFloat32(data[3]),
	}
	return dailyClose
}
