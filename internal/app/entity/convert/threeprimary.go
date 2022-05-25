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
	"strings"

	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/helper"
)

type threePrimaryImpl struct {
}

func ThreePrimary() IConvert {
	return &threePrimaryImpl{}
}

func (c *threePrimaryImpl) Execute(data *ConvertData) interface{} {
	var output *entity.ThreePrimary
	if data == nil || len(data.RawData) < 19 {
		return output
	}
	switch data.Target {
	case TpexThreePrimary:
		output = &entity.ThreePrimary{
			StockID:            data.RawData[0],
			Date:               data.ParseDate,
			ForeignTradeShares: helper.ToInt64(strings.Replace(data.RawData[10], ",", "", -1)),
			TrustTradeShares:   helper.ToInt64(strings.Replace(data.RawData[13], ",", "", -1)),
			DealerTradeShares:  helper.ToInt64(strings.Replace(data.RawData[16], ",", "", -1)),
			HedgingTradeShares: helper.ToInt64(strings.Replace(data.RawData[19], ",", "", -1)),
		}
	case TwseThreePrimary:
		output = &entity.ThreePrimary{
			StockID:            data.RawData[0],
			Date:               data.ParseDate,
			ForeignTradeShares: helper.ToInt64(strings.Replace(data.RawData[4], ",", "", -1)),
			TrustTradeShares:   helper.ToInt64(strings.Replace(data.RawData[10], ",", "", -1)),
			DealerTradeShares:  helper.ToInt64(strings.Replace(data.RawData[14], ",", "", -1)),
			HedgingTradeShares: helper.ToInt64(strings.Replace(data.RawData[17], ",", "", -1)),
		}
	}
	return output
}
