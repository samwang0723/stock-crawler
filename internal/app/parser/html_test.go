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

package parser

import (
	"strings"
	"testing"
)

func Test_parseHtml(t *testing.T) {
	wrongDoc := "<html><body></body></html>"
	correctDoc := `
	<link rel="stylesheet" href="http://isin.twse.com.tw/isin/style1.css" type="text/css">
	<body><table  align=center><h2><strong><font class='h1'>本國上市證券國際證券辨識號碼一覽表</font></strong>
	</h2><h2><strong><font class='h1'><center>最近更新日期:2021/12/03  </center> </font></strong></h2><h2>
	<font color='red'><center>掛牌日以正式公告為準</center></font></h2>
	</table>
	<TABLE class='h4' align=center cellSpacing=3 cellPadding=2 width=750 border=0>
	<tr align=center>
		<td bgcolor=#D5FFD5>有價證券代號及名稱 </td><td bgcolor=#D5FFD5>國際證券辨識號碼(ISIN Code)</td>
		<td bgcolor=#D5FFD5>上市日</td><td bgcolor=#D5FFD5>市場別</td>
		<td bgcolor=#D5FFD5>產業別</td>
		<td bgcolor=#D5FFD5>CFICode</td><td bgcolor=#D5FFD5>備註</td>
	</tr>
	<tr><td bgcolor=#FAFAD2 colspan=7 ><B> 股票 <B> </td></tr>
	<tr><td bgcolor=#FAFAD2>1101　台泥</td><td bgcolor=#FAFAD2>TW0001101004</td><td bgcolor=#FAFAD2>1962/02/09</td><td bgcolor=#FAFAD2>上市</td><td bgcolor=#FAFAD2>水泥工業</td><td bgcolor=#FAFAD2>ESVUFR</td><td bgcolor=#FAFAD2></td></tr>
	<tr><td bgcolor=#FAFAD2>1102　亞泥</td><td bgcolor=#FAFAD2>TW0001102002</td><td bgcolor=#FAFAD2>1962/06/08</td><td bgcolor=#FAFAD2>上市</td><td bgcolor=#FAFAD2>水泥工業</td><td bgcolor=#FAFAD2>ESVUFR</td><td bgcolor=#FAFAD2></td></tr>
	<tr><td bgcolor=#FAFAD2>1103　嘉泥</td><td bgcolor=#FAFAD2>TW0001103000</td><td bgcolor=#FAFAD2>1969/11/14</td><td bgcolor=#FAFAD2>上市</td><td bgcolor=#FAFAD2>水泥工業</td><td bgcolor=#FAFAD2>ESVUFR</td><td bgcolor=#FAFAD2></td></tr>
	<tr><td bgcolor=#FAFAD2>1104　環泥</td><td bgcolor=#FAFAD2>TW0001104008</td><td bgcolor=#FAFAD2>1971/02/01</td><td bgcolor=#FAFAD2>上市</td><td bgcolor=#FAFAD2>水泥工業</td><td bgcolor=#FAFAD2>ESVUFR</td><td bgcolor=#FAFAD2></td></tr>
	<tr><td bgcolor=#FAFAD2 colspan=7 ><B> 特別股 <B> </td></tr>
	<tr><td bgcolor=#FAFAD2>1101B　台泥乙特</td><td bgcolor=#FAFAD2>TW0001101B05</td><td bgcolor=#FAFAD2>2019/01/29</td><td bgcolor=#FAFAD2>上市</td><td bgcolor=#FAFAD2></td><td bgcolor=#FAFAD2>EPNRAR</td><td bgcolor=#FAFAD2></td></tr>
	<tr><td bgcolor=#FAFAD2>1312A　國喬特</td><td bgcolor=#FAFAD2>TW0001312A01</td><td bgcolor=#FAFAD2>1988/12/21</td><td bgcolor=#FAFAD2>上市</td><td bgcolor=#FAFAD2>塑膠工業</td><td bgcolor=#FAFAD2>EPNRQR</td><td bgcolor=#FAFAD2></td></tr>
	</table>
	`

	tests := []struct {
		name    string
		content string
		want    int
	}{
		{
			name:    "normal html",
			content: correctDoc,
			want:    4,
		},
		{
			name:    "wrong html",
			content: wrongDoc,
			want:    0,
		},
	}

	for _, tt := range tests {
		tt := tt
		conf := Config{
			Capacity: 6,
			Type:     TwseStockList,
		}

		t.Run(tt.name, func(t *testing.T) {
			//t.Parallel()
			in := strings.NewReader(tt.content)
			res := &parserImpl{
				result: &[]interface{}{},
			}
			res.parseHtml(conf, in)

			if got := len(*res.result); got != tt.want {
				t.Errorf("len(parser.result) = %v, want %v", got, tt.want)
			}
		})
	}

}
