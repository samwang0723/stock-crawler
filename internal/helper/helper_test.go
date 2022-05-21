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

package helper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_IsInteger(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want bool
	}{
		{
			name: "integer string",
			val:  "4968",
			want: true,
		},
		{
			name: "non-integer string",
			val:  "ABCD",
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := IsInteger(tt.val); got != tt.want {
				t.Errorf("IsInteger(t.val) = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ToInt64(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want int64
	}{
		{
			name: "integer string",
			val:  "-4968",
			want: -4968,
		},
		{
			name: "non-integer string",
			val:  "ABCD",
			want: 0,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ToInt64(tt.val); got != tt.want {
				t.Errorf("ToInt64(t.val) = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ToUint64(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want uint64
	}{
		{
			name: "integer string",
			val:  "4968",
			want: 4968,
		},
		{
			name: "signed-integer string",
			val:  "-4968",
			want: 0,
		},
		{
			name: "non-integer string",
			val:  "ABCD",
			want: 0,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ToUint64(tt.val); got != tt.want {
				t.Errorf("ToUint64(t.val) = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ToFloat32(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want float32
	}{
		{
			name: "float string",
			val:  "320.32",
			want: 320.32,
		},
		{
			name: "non-float string",
			val:  "ABCD",
			want: 0.0,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ToFloat32(tt.val); got != tt.want {
				t.Errorf("ToFloat32(t.val) = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_FormalizeValidTimeWithLocation(t *testing.T) {
	l, _ := time.LoadLocation(TimeZone)
	testTime := time.Date(2021, time.December, 23, 12, 0, 0, 0, l)
	expTime := time.Date(2021, time.December, 22, 12, 0, 0, 0, l)
	tests := []struct {
		name   string
		ti     time.Time
		offset int
		want   *time.Time
	}{
		{
			name:   "regular time",
			ti:     testTime,
			offset: -1,
			want:   &expTime,
		},
		{
			name:   "weekend time",
			ti:     testTime,
			offset: -4,
			want:   nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := FormalizeValidTimeWithLocation(tt.ti, int32(tt.offset))
			//t.Errorf("FormalizeValidTimeWithLocation(tt.ti, tt.offset) = %v, want %v", got, tt.want)
			if got != nil {
				assert.True(t, got.Equal(*tt.want))
			} else {
				assert.Equal(t, got, tt.want)
			}
		})
	}
}

func Test_GetDateFromUTC(t *testing.T) {
	tests := []struct {
		name   string
		val    string
		format string
		want   string
	}{
		{
			name:   "convert timestamp to Twse",
			val:    "1641323046",
			format: TwseDateFormat,
			want:   "20220105",
		},
		{
			name:   "convert timestamp to Tpex",
			val:    "1641323046",
			format: TpexDateFormat,
			want:   "111/01/05",
		},
		{
			name:   "wrong timestamp",
			val:    "",
			format: TpexDateFormat,
			want:   "",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := GetDateFromUTC(tt.val, tt.format); got != tt.want {
				t.Errorf("GetDateFromUTC(t.val, tt.format) = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetDateFromOffset(t *testing.T) {
	l, _ := time.LoadLocation(TimeZone)
	expTime := time.Date(2021, time.December, 22, 12, 0, 0, 0, l)

	tests := []struct {
		name   string
		format string
		want   string
		offset int
	}{
		{
			name:   "convert current timestamp to Twse",
			offset: 0,
			format: TwseDateFormat,
			want:   expTime.Format(TwseDateFormat),
		},
		{
			name:   "convert current timestamp to Tpex",
			offset: 0,
			format: TpexDateFormat,
			want:   UnifiedDateFormatToTpex(expTime.Format(TpexDateFormat)),
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := GetDateFromOffset(int32(tt.offset), tt.format, expTime); got != tt.want {
				t.Errorf("GetDateFromOffset(t.offset, tt.format) = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetReadableSize(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		length  int
		decimal int
	}{
		{
			name:    "convert bytes to MB",
			length:  1024678,
			decimal: 2,
			want:    "1.02 MB",
		},
		{
			name:    "convert bytes to GB",
			length:  20010246700,
			decimal: 0,
			want:    "20 GB",
		},
		{
			name:    "convert bytes to KB",
			length:  340456,
			decimal: 3,
			want:    "340.456 KB",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := GetReadableSize(tt.length, tt.decimal); got != tt.want {
				t.Errorf("GetReadableSize(t.length, tt.decimal) = %v, want %v", got, tt.want)
			}
		})
	}
}
