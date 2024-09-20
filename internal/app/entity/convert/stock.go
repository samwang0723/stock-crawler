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
)

type stockImpl struct{}

const maxLength = 5

func Stock() IConvert {
	return &stockImpl{}
}

func (c *stockImpl) Execute(data *Data) any {
	var output *entity.Stock
	if data == nil || len(data.RawData) < maxLength {
		return output
	}

	str := strings.Split(data.RawData[0], "　")
	t := strings.TrimSpace(data.RawData[3])

	market := "tse"
	if strings.Contains(t, "上櫃") {
		market = "otc"
	}

	if len(data.RawData) == maxLength {
		data.RawData[4] = "臺灣存託憑證(TDR)"
	}

	output = &entity.Stock{
		StockID:  strings.TrimSpace(str[0]),
		Name:     strings.TrimSpace(str[1]),
		Country:  "TW",
		Market:   market,
		Category: strings.TrimSpace(data.RawData[4]),
	}

	return output
}
