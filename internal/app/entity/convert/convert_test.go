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
	"testing"

	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/stretchr/testify/assert"
)

func Test_Concentration(t *testing.T) {
	tests := []struct {
		name string
		val  *ConvertData
		exp  *entity.StakeConcentration
	}{
		{
			name: "convert StakeConcentration",
			val: &ConvertData{
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
			name: "empty ConvertData",
			val:  nil,
			exp:  nil,
		},
		{
			name: "missing elements in ConvertData",
			val: &ConvertData{
				RawData: []string{"1", "20220525", "2330", "1000"},
				Target:  StakeConcentration,
			},
			exp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
