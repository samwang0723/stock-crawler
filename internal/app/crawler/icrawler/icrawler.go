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
package icrawler

import (
	"context"
)

const (
	All                = "ALL"
	StockOnly          = "ALLBUT0999"
	TwseDailyClose     = "https://www.twse.com.tw/exchangeReport/MI_INDEX?response=csv&date=%s&type=ALLBUT0999"
	TwseThreePrimary   = "http://www.tse.com.tw/fund/T86?response=csv&date=%s&selectType=ALLBUT0999"
	OperatingDays      = "https://www.twse.com.tw/holidaySchedule/holidaySchedule?response=csv&queryYear=%d"
	TpexDailyClose     = "http://www.tpex.org.tw/web/stock/aftertrading/daily_close_quotes/stk_quote_download.php?l=zh-tw&d=%s&s=0,asc,0"
	TpexThreePrimary   = "https://www.tpex.org.tw/web/stock/3insti/daily_trade/3itrade_hedge_result.php?l=zh-tw&o=csv&se=EW&t=D&d=%s"
	TWSEStocks         = "https://isin.twse.com.tw/isin/C_public.jsp?strMode=2"
	TPEXStocks         = "https://isin.twse.com.tw/isin/C_public.jsp?strMode=4"
	StakeConcentration = "https://stockchannelnew.sinotrade.com.tw/z/zc/zco/zco.djhtm?a=%s&e=%s&f=%s"
	ConcentrationDays  = "https://stockchannelnew.sinotrade.com.tw/z/zc/zco/zco_%s_%d.djhtm"
)

type ICrawler interface {
	Fetch(ctx context.Context) (string, []byte, error)
	AppendURL(url string)
	GetURLs() []string
}
