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
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/helper"
)

type csvStrategy struct {
	converter convert.IConvert
	date      string
	source    convert.Source
	capacity  int
}

func (s *csvStrategy) Parse(in io.Reader, additional ...string) ([]interface{}, error) {
	if len(s.date) == 0 {
		return nil, fmt.Errorf("parse day missing")
	}

	var output []interface{}
	updatedLen := 0

	reader := csv.NewReader(in)
	reader.Comma = ','
	reader.FieldsPerRecord = -1

	//override to standarize date string (20211123)
	date := helper.UnifiedDateFormatToTwse(s.date)

	for {
		records, err := reader.Read()
		if err == io.EOF {
			break
		} else if len(records) == 0 || s.capacity > len(records) {
			continue
		}

		// make sure only parse recognized stock_id
		records[0] = strings.TrimSpace(records[0])
		if len(records[0]) > 0 && len(records[0]) < 6 && helper.IsInteger(records[0][0:2]) {
			output = append(
				output,
				s.converter.Execute(&convert.ConvertData{
					ParseDate: date,
					RawData:   records,
					Target:    s.source,
				}),
			)
			updatedLen++
		}
	}
	if updatedLen == 0 {
		return nil, NoParseResults
	}

	return output, nil
}
