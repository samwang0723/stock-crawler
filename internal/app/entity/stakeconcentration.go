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

package entity

import (
	"strconv"
	"sync"
)

var (
	//nolint:nolintlint, gochecknoglobals
	concentrationPool = sync.Pool{
		New: func() interface{} { return new(StakeConcentration) },
	}
)

type StakeConcentration struct {
	StockID       string  `json:"stockId"`
	Date          string  `json:"exchangeDate"`
	HiddenField   string  `json:"-"` // this field is use to identify which period the SumBuyShares/SumSellShares are
	Diff          []int32 `json:"diff"`
	SumBuyShares  uint64  `json:"sumBuyShares"`
	SumSellShares uint64  `json:"sumSellShares"`
	AvgBuyPrice   float32 `json:"avgBuyPrice"`
	AvgSellPrice  float32 `json:"avgSellPrice"`
}

func MapReduceStakeConcentration(objs []*StakeConcentration) *StakeConcentration {
	volumeDiff := []int32{0, 0, 0, 0, 0}

	var res *StakeConcentration

	for _, val := range objs {
		idx, err := strconv.Atoi(val.HiddenField)
		// make sure to cover latest source of truth date's concentration data
		if err == nil && idx == 0 {
			res = val
		}

		volumeDiff[idx] = int32(val.SumBuyShares - val.SumSellShares)
	}

	res.Diff = volumeDiff

	return res.Clone()
}

func (sc *StakeConcentration) Clone() *StakeConcentration {
	newSc, ok := concentrationPool.Get().(*StakeConcentration)
	if !ok {
		return nil
	}

	newSc.StockID = sc.StockID
	newSc.Date = sc.Date
	newSc.Diff = sc.Diff
	newSc.SumBuyShares = sc.SumBuyShares
	newSc.SumSellShares = sc.SumSellShares
	newSc.AvgBuyPrice = sc.AvgBuyPrice
	newSc.AvgSellPrice = sc.AvgSellPrice

	return newSc
}

func (sc *StakeConcentration) Recycle() {
	sc.StockID = ""
	sc.Date = ""
	sc.Diff = sc.Diff[:0]
	sc.SumBuyShares = 0
	sc.SumSellShares = 0
	sc.AvgBuyPrice = 0.0
	sc.AvgSellPrice = 0.0

	concentrationPool.Put(sc)
}
