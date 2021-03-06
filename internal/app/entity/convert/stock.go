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

type stockImpl struct {
}

func Stock() IConvert {
	return &stockImpl{}
}

func (c *stockImpl) Execute(data *ConvertData) interface{} {
	var output *entity.Stock
	if data == nil || len(data.RawData) < 6 {
		return output
	}

	s := strings.Split(data.RawData[0], "　")
	output = &entity.Stock{
		StockID:  strings.TrimSpace(s[0]),
		Name:     strings.TrimSpace(s[1]),
		Country:  "TW",
		Category: strings.TrimSpace(data.RawData[4]),
	}
	return output
}
