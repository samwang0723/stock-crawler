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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

const (
	TB                       = 1000000000000
	GB                       = 1000000000
	MB                       = 1000000
	KB                       = 1000
	TimeZone                 = "Asia/Taipei"
	TwseDateFormat           = "20060102"
	TpexDateFormat           = "2006/01/02"
	StakeConcentrationFormat = "2006-01-02"
	Signature                = `
 _____ _             _                                  _           
/  ___| |           | |                                | |          
\ '--.| |_ ___   ___| | ________ ___ _ __ __ ___      _| | ___ _ __ 
 '--. \ __/ _ \ / __| |/ /______/ __| '__/ _' \ \ /\ / / |/ _ \ '__|
/\__/ / || (_) | (__|   <      | (__| | | (_| |\ V  V /| |  __/ |   
\____/ \__\___/ \___|_|\_\      \___|_|  \__,_| \_/\_/ |_|\___|_|

                                                        Version (%s)
Stand-alone stock data crawling service
Environment (%s)
_______________________________________________
`
)

func ReadFromFile(fileName string) (string, error) {
	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func EncodeBig5(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func IsInteger(v string) bool {
	if _, err := strconv.Atoi(v); err == nil {
		return true
	}
	return false
}

func ToInt64(v string) int64 {
	if i, err := strconv.ParseInt(v, 10, 64); err == nil {
		return i
	}

	return 0
}

func ToUint64(v string) uint64 {
	if i, err := strconv.ParseUint(v, 10, 64); err == nil {
		return i
	}

	return 0
}

func ToFloat32(v string) float32 {
	if f, err := strconv.ParseFloat(v, 32); err == nil {
		return float32(f)
	}
	return 0
}

func FormalizeValidTimeWithLocation(input time.Time, offset ...int32) *time.Time {
	l, _ := time.LoadLocation(TimeZone)
	t := input.In(l)
	if len(offset) > 0 {
		t = t.AddDate(0, 0, int(offset[0]))
	}

	if IncludeWeekend() == false {
		// only within workday will be valid
		wkDay := t.Weekday()
		if wkDay == time.Saturday || wkDay == time.Sunday {
			return nil
		}
	}

	return &t
}

func GetDateFromUTC(timestamp string, format string) string {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return ""
	}
	t := FormalizeValidTimeWithLocation(time.Unix(i, 0))
	if t == nil {
		return ""
	}

	// Twse format: 20190213
	s := t.Format(format)

	switch format {
	case TpexDateFormat:
		// Tpex format: 108/02/06
		s = UnifiedDateFormatToTpex(s)
	}
	return s
}

func GetDateFromOffset(offset int32, format string, input ...time.Time) string {
	current := time.Now()
	if len(input) > 0 {
		current = input[0]
	}
	t := FormalizeValidTimeWithLocation(current, offset)
	if t == nil {
		return ""
	}
	timestamp := strconv.FormatInt(t.Unix(), 10)
	return GetDateFromUTC(timestamp, format)
}

func UnifiedDateFormatToTpex(input string) string {
	if strings.Contains(input, "/") {
		res := strings.Split(input, "/")
		year, _ := strconv.Atoi(res[0])
		return fmt.Sprintf("%d/%s/%s", year-1911, res[1], res[2])
	}
	return input
}

func UnifiedDateFormatToTwse(input string) string {
	if strings.Contains(input, "/") {
		i := strings.Split(input, "/")
		year, _ := strconv.Atoi(i[0])
		res := fmt.Sprintf("%d%s%s", year+1911, i[1], i[2])
		return res
	}
	return input
}

func GetReadableSize(length int, decimals int) (out string) {
	var unit string
	var i int
	var remainder int

	// Get whole number, and the remainder for decimals
	if length > TB {
		unit = "TB"
		i = length / TB
		remainder = length - (i * TB)
	} else if length > GB {
		unit = "GB"
		i = length / GB
		remainder = length - (i * GB)
	} else if length > MB {
		unit = "MB"
		i = length / MB
		remainder = length - (i * MB)
	} else if length > KB {
		unit = "KB"
		i = length / KB
		remainder = length - (i * KB)
	} else {
		return strconv.Itoa(length) + " B"
	}

	if decimals == 0 {
		return strconv.Itoa(i) + " " + unit
	}

	// This is to calculate missing leading zeroes
	width := 0
	if remainder > GB {
		width = 12
	} else if remainder > MB {
		width = 9
	} else if remainder > KB {
		width = 6
	} else {
		width = 3
	}

	// Insert missing leading zeroes
	remainderString := strconv.Itoa(remainder)
	for iter := len(remainderString); iter < width; iter++ {
		remainderString = "0" + remainderString
	}
	if decimals > len(remainderString) {
		decimals = len(remainderString)
	}

	return fmt.Sprintf("%d.%s %s", i, remainderString[:decimals], unit)
}

func GetCurrentEnv() string {
	env := os.Getenv("ENVIRONMENT")
	output := "dev"
	switch env {
	case "development":
		output = "dev"
	case "staging":
		output = "staging"
	case "production":
		output = "prod"
	}
	return output
}

func IncludeWeekend() bool {
	includeWeekend := os.Getenv("INCLUDE_WEEKEND")
	isTest, err := strconv.ParseBool(includeWeekend)
	if err != nil {
		return false
	}
	return isTest
}

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Diff(slice1 []string, slice2 []string) []string {
	diffStr := []string{}
	m := map[string]int{}

	for _, s1Val := range slice1 {
		m[s1Val] = 1
	}
	for _, s2Val := range slice2 {
		m[s2Val] = m[s2Val] + 1
	}

	for mKey, mVal := range m {
		if mVal == 1 {
			diffStr = append(diffStr, mKey)
		}
	}

	return diffStr
}
