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
package entity

type ThreePrimary struct {
	StockID            string `json:"stockId"`
	Date               string `json:"date"`
	ForeignTradeShares int64  `json:"foreignTradeShares"`
	TrustTradeShares   int64  `json:"trustTradeShares"`
	DealerTradeShares  int64  `json:"dealerTradeShares"`
	HedgingTradeShares int64  `json:"hedgingTradeShares"`
}
