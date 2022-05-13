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

	"github.com/samwang0723/stock-crawler/internal/app/entity"
)

func Test_parseConcentration(t *testing.T) {
	wrongDoc := "<html><body><table><tr><td>WRONG</td></tr></table></body></html>"
	correctDoc := `<html>
	<head>
		<title>主力賣買超-6727</title>
	</head>
	<body>
	<table id="oMainTable" class="t01" width="612" border="0" cellspacing="1" cellpadding="1">
	<tr id="oScrollHead"><td class="t10" colspan="10"><select name="stk" onchange="changeStkID(this.options[this.selectedIndex].value);">
	<option value="" selected>亞泰金屬(6727)</option><option value="67271">亞泰金屬一(67271)</option>
	</select>主力進出比較圖</td></tr>
	<tr id="oScrollHead"><td colspan="10" class="t3n0" width="100%"></tr>
	</table>
	<table>
		<tr id="oScrollHead"><td class="t10" colspan="10">亞泰金屬(6727)
		主力進出明細
		<div class="t11">單位：張　最後更新日：2021/12/01</div></td></tr>
		<tr id="oScrollHead"><td class="t2" colspan="10">
		<select name="D" onChange="goA();">
		<option value="0" selected>請選擇</option>
		<option value="1">近一日</option>
		<option value="2">近五日</option>
		<option value="3">近十日</option>
		<option value="4">近20日</option>
		<option value="5">近40日</option>
		<option value="6">近60日</option>
		<option value="7">近120日</option>
		<option value="8">近240日</option>
		</select>
		<TR>
		<TD class="t4t1" nowrap><a href="/z/zc/zco/zco0/zco0.djhtm?a=2303&b=0038003800380041&BHID=8880">國泰-館前</a></TD>
		<TD class="t3n1">690</TD>
		<TD class="t3n1">174</TD>
		<TD class="t3n1">516</TD>
		<TD class="t3n1">0.38%</TD>
		<TD class="t4t1" nowrap><a href="/z/zc/zco/zco0/zco0.djhtm?a=2303&b=0039003100380064&BHID=9100">群益金鼎-內湖</a></TD>
		<TD class="t3n1">279</TD>
		<TD class="t3n1">1,610</TD>
		<TD class="t3n1">1,331</TD>
		<TD class="t3n1">0.97%</TD>
		</tr>
		<TR>
		<TD class="t4t1" nowrap><a href="/z/zc/zco/zco0/zco0.djhtm?a=2303&b=1590&BHID=1590">花旗環球</a></TD>
		<TD class="t3n1">728</TD>
		<TD class="t3n1">270</TD>
		<TD class="t3n1">458</TD>
		<TD class="t3n1">0.34%</TD>
		<TD class="t4t1" nowrap><a href="/z/zc/zco/zco0/zco0.djhtm?a=2303&b=0039004100390052&BHID=9A00">永豐金-信義</a></TD>
		<TD class="t3n1">86</TD>
		<TD class="t3n1">1,402</TD>
		<TD class="t3n1">1,316</TD>
		<TD class="t3n1">0.96%</TD>
		</tr>
		<TR id="oScrollFoot">
			<TD class="t4t1" nowrap>合計買超張數</td>
			<td class="t3n1" colspan=4>12,449</td>
			<TD class="t4t1" nowrap>合計賣超張數</td>
			<td class="t3n1" colspan=4>40,221</td>
		</TR>
		<TR id="oScrollFoot">
			<TD class="t4t1" nowrap>平均買超成本</td>
			<td class="t3n1" colspan=4>63.45</td>
			<TD class="t4t1" nowrap>平均賣超成本</td>
			<td class="t3n1" colspan=4>63.53</td>
		</TR>
		<TR id="oScrollFoot">
			<td class="t3t1" colspan=10>
			【註1】上述買賣超個股僅提供排序後的前15名券商，且未計入自營商部份。<BR>
			【註2】合計買超或賣超，為上述家數合計。<BR>
			【註3】平均買超或賣超成本，為上述家數合計買賣超金額/上述家數合計買賣超張數。
			</td>
		</TR>
	    </table>
	</body>
	</html>`

	tests := []struct {
		name    string
		content string
		want    bool
		shares  []uint64
		price   []float32
	}{
		{
			name:    "normal html",
			content: correctDoc,
			want:    true,
			shares:  []uint64{12449, 40221},
			price:   []float32{63.45, 63.53},
		},
		{
			name:    "wrong html",
			content: wrongDoc,
			want:    false,
		},
	}

	for _, tt := range tests {
		tt := tt
		d := "2021-10-29"
		conf := Config{
			ParseDay: &d,
			Capacity: 4,
			Type:     StakeConcentration,
		}

		t.Run(tt.name, func(t *testing.T) {
			//t.Parallel()
			in := strings.NewReader(tt.content)
			res := &parserImpl{
				result: &[]interface{}{},
			}
			res.parseConcentration(conf, in)

			respLen := len(*res.result)
			if respLen > 0 != tt.want {
				t.Errorf("len(parser.result) = %v, want %v", respLen, tt.want)
			} else if respLen > 0 {
				c := (*res.result)[0].(*entity.StakeConcentration)
				if c.SumBuyShares != tt.shares[0] ||
					c.SumSellShares != tt.shares[1] ||
					c.AvgBuyPrice != tt.price[0] ||
					c.AvgSellPrice != tt.price[1] {
					t.Errorf("response details = %+v is not fit the expectation", c)
				}
			}
		})
	}
}
