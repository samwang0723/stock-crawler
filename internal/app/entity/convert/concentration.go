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
	"strconv"

	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/helper"
)

type concentrationImpl struct{}

func Concentration() IConvert {
	return &concentrationImpl{}
}

func (c *concentrationImpl) Execute(data *Data) interface{} {
	var output *entity.StakeConcentration
	if data == nil || len(data.RawData) < 7 {
		return output
	}

	return &entity.StakeConcentration{
		HiddenField:   c.convertHiddenIndex(data.RawData[0]),
		Date:          data.RawData[1],
		StockID:       data.RawData[2],
		SumBuyShares:  helper.ToUint64(data.RawData[3]),
		SumSellShares: helper.ToUint64(data.RawData[4]),
		AvgBuyPrice:   helper.ToFloat32(data.RawData[5]),
		AvgSellPrice:  helper.ToFloat32(data.RawData[6]),
	}
}

func (c *concentrationImpl) convertHiddenIndex(idx string) string {
	// original ids from parser represents 3rd party concentration index
	// need to convert back to match our concentration model diff index
	// 1 = concentration within 1 day = Diff[0]
	// 2 = concentration within 5 day = Diff[1]
	// 3 = concentration within 10 day = Diff[2]
	// 4 = concentration within 20 day = Diff[3]
	// 6 = concentration within 60 day = Diff[4]
	switch idx {
	case "6":
		return "4"
	default:
		newIdx, err := strconv.Atoi(idx)
		if err != nil {
			return ""
		}

		return strconv.Itoa(newIdx - 1)
	}
}
