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
package helper

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"golang.org/x/xerrors"
)

//nolint:nolintlint, varnamelen, dupword
const (
	TB                       = 1000000000000
	GB                       = 1000000000
	MB                       = 1000000
	KB                       = 1000
	TimeZone                 = "Asia/Taipei"
	TwseDateFormat           = "20060102"
	TpexDateFormat           = "2006/01/02"
	StakeConcentrationFormat = "2006-01-02"
)

func ReadFromFile(fileName string) (string, error) {
	bs, err := os.ReadFile(fileName)
	if err != nil {
		return "", xerrors.Errorf("helper.ReadFromFile: failed, file=%s; err=%w;", fileName, err)
	}

	return string(bs), nil
}

func EncodeBig5(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewEncoder())
	data, e := io.ReadAll(O)

	if e != nil {
		return nil, xerrors.Errorf("helper.EncodeBig5: failed, err=%w;", e)
	}

	return data, nil
}

func IsInteger(v string) bool {
	if _, err := strconv.Atoi(v); err == nil {
		return true
	}

	return false
}

//nolint:nolintlint, gomnd
func ToInt64(v string) int64 {
	if i, err := strconv.ParseInt(v, 10, 64); err == nil {
		return i
	}

	return 0
}

//nolint:nolintlint, gomnd
func ToUint64(v string) uint64 {
	if i, err := strconv.ParseUint(v, 10, 64); err == nil {
		return i
	}

	return 0
}

//nolint:nolintlint, gomnd
func ToFloat32(v string) float32 {
	if f, err := strconv.ParseFloat(v, 32); err == nil {
		return float32(f)
	}

	return 0
}

func FormalizeValidTimeWithLocation(input time.Time, offset ...int32) *time.Time {
	loc, err := time.LoadLocation(TimeZone)
	if err != nil {
		return nil
	}

	localTime := input.In(loc)

	if len(offset) > 0 {
		localTime = localTime.AddDate(0, 0, int(offset[0]))
	}

	if !IncludeWeekend() {
		// only within workday will be valid
		wkDay := localTime.Weekday()
		if wkDay == time.Saturday || wkDay == time.Sunday {
			return nil
		}
	}

	return &localTime
}

func GetDateFromUTC(timestamp, format string) string {
	//nolint:nolintlint, gomnd
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return ""
	}

	localTime := FormalizeValidTimeWithLocation(time.Unix(i, 0))

	if localTime == nil {
		return ""
	}

	// Twse format: 20190213
	dataStr := localTime.Format(format)

	if format == TpexDateFormat {
		// Tpex format: 108/02/06
		dataStr = UnifiedDateFormatToTpex(dataStr)
	}

	return dataStr
}

func GetDateFromOffset(offset int32, format string, input ...time.Time) string {
	current := time.Now()

	if len(input) > 0 {
		current = input[0]
	}

	localTime := FormalizeValidTimeWithLocation(current, offset)

	if localTime == nil {
		return ""
	}

	//nolint:nolintlint, gomnd
	timestamp := strconv.FormatInt(localTime.Unix(), 10)

	return GetDateFromUTC(timestamp, format)
}

func UnifiedDateFormatToTpex(input string) string {
	if strings.Contains(input, "/") {
		res := strings.Split(input, "/")

		year, err := strconv.Atoi(res[0])
		if err != nil {
			return ""
		}

		//nolint:nolintlint, gomnd
		return fmt.Sprintf("%d/%s/%s", year-1911, res[1], res[2])
	}

	return input
}

func UnifiedDateFormatToTwse(input string) string {
	if strings.Contains(input, "/") {
		res := strings.Split(input, "/")

		year, err := strconv.Atoi(res[0])
		if err != nil {
			return ""
		}

		//nolint:nolintlint, gomnd
		resNew := fmt.Sprintf("%d%s%s", year+1911, res[1], res[2])

		return resNew
	}

	unifiedDate := strings.ReplaceAll(input, "-", "")

	return unifiedDate
}

//nolint:nolintlint, cyclop
func GetReadableSize(length, decimals int) string {
	var unit string

	var calcVal, width, remainder int

	// Get whole number, and the remainder for decimals
	switch {
	case length > TB:
		unit = "TB"
		calcVal = length / TB
		remainder = length - (calcVal * TB)
	case length > GB:
		unit = "GB"
		calcVal = length / GB
		remainder = length - (calcVal * GB)
	case length > MB:
		unit = "MB"
		calcVal = length / MB
		remainder = length - (calcVal * MB)
	case length > KB:
		unit = "KB"
		calcVal = length / KB
		remainder = length - (calcVal * KB)
	default:
		return strconv.Itoa(length) + " B"
	}

	if decimals == 0 {
		return strconv.Itoa(calcVal) + " " + unit
	}

	// This is to calculate missing leading zeroes
	switch {
	case remainder > GB:
		width = 12
	case remainder > MB:
		width = 9
	case remainder > KB:
		width = 6
	default:
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

	return fmt.Sprintf("%d.%s %s", calcVal, remainderString[:decimals], unit)
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

//nolint:nolintlint, govet
func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	byteHeader := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}

	return *(*[]byte)(unsafe.Pointer(&byteHeader))
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Diff(slice1, slice2 []string) []string {
	diffStr := []string{}
	res := map[string]int{}

	for _, s1Val := range slice1 {
		res[s1Val] = 1
	}

	for _, s2Val := range slice2 {
		res[s2Val]++
	}

	for mKey, mVal := range res {
		if mVal == 1 {
			diffStr = append(diffStr, mKey)
		}
	}

	return diffStr
}
