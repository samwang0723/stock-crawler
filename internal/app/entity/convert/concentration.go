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
	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/helper"
)

type concentrationImpl struct {
}

func Concentration() IConvert {
	return &concentrationImpl{}
}

func (c *concentrationImpl) Execute(data *ConvertData) interface{} {
	return &entity.StakeConcentration{
		HiddenField:   data.RawData[0],
		Date:          data.RawData[1],
		StockID:       data.RawData[2],
		SumBuyShares:  helper.ToUint64(data.RawData[3]),
		SumSellShares: helper.ToUint64(data.RawData[4]),
		AvgBuyPrice:   helper.ToFloat32(data.RawData[5]),
		AvgSellPrice:  helper.ToFloat32(data.RawData[6]),
	}
}
