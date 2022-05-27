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

type DailyClose struct {
	StockID      string  `json:"stockId"`
	Date         string  `json:"date"`
	TradedShares uint64  `json:"tradeShares"`  // Total volumes of shares being traded.
	Transactions uint64  `json:"transactions"` // Total numbers of transaction.
	Turnover     uint64  `json:"turnover"`     // Total traded dollar volume
	Open         float32 `json:"open"`
	Close        float32 `json:"close"`
	High         float32 `json:"high"`
	Low          float32 `json:"low"`
	PriceDiff    float32 `json:"priceDiff"`
}
