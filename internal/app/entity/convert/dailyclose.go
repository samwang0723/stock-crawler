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

package convert

import (
	"fmt"
	"strings"

	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/helper"
)

type dailyCloseImpl struct{}

func DailyClose() IConvert {
	return &dailyCloseImpl{}
}

func (c *dailyCloseImpl) Execute(data *Data) any {
	var output *entity.DailyClose
	if data == nil || len(data.RawData) < 17 {
		return output
	}

	//nolint:nolintlint, exhaustive
	switch data.Target {
	case TpexDailyClose:
		output = &entity.DailyClose{
			StockID:      data.RawData[0],
			Date:         data.ParseDate,
			TradedShares: helper.ToUint64(strings.ReplaceAll(data.RawData[7], ",", "")),
			Transactions: helper.ToUint64(strings.ReplaceAll(data.RawData[9], ",", "")),
			Turnover:     helper.ToUint64(strings.ReplaceAll(data.RawData[8], ",", "")),
			Open:         helper.ToFloat32(strings.ReplaceAll(data.RawData[4], ",", "")),
			High:         helper.ToFloat32(strings.ReplaceAll(data.RawData[5], ",", "")),
			Low:          helper.ToFloat32(strings.ReplaceAll(data.RawData[6], ",", "")),
			Close:        helper.ToFloat32(strings.ReplaceAll(data.RawData[2], ",", "")),
			PriceDiff:    helper.ToFloat32(strings.TrimSpace(data.RawData[3])),
		}
	case TwseDailyClose:
		output = &entity.DailyClose{
			StockID:      data.RawData[0],
			Date:         data.ParseDate,
			TradedShares: helper.ToUint64(strings.ReplaceAll(data.RawData[2], ",", "")),
			Transactions: helper.ToUint64(strings.ReplaceAll(data.RawData[3], ",", "")),
			Turnover:     helper.ToUint64(strings.ReplaceAll(data.RawData[4], ",", "")),
			Open:         helper.ToFloat32(strings.ReplaceAll(data.RawData[5], ",", "")),
			High:         helper.ToFloat32(strings.ReplaceAll(data.RawData[6], ",", "")),
			Low:          helper.ToFloat32(strings.ReplaceAll(data.RawData[7], ",", "")),
			Close:        helper.ToFloat32(strings.ReplaceAll(data.RawData[8], ",", "")),
			PriceDiff: helper.ToFloat32(
				fmt.Sprintf(
					"%s%s",
					strings.TrimSpace(data.RawData[9]),
					strings.TrimSpace(data.RawData[10]),
				),
			),
		}
	}

	return output
}
