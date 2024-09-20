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
	"errors"
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

//nolint:nolintlint, cyclop
func (s *csvStrategy) Parse(input io.Reader, _ ...string) ([]any, error) {
	if s.date == "" {
		return nil, ErrParseDayMissing
	}

	var output []any

	reader := csv.NewReader(input)
	reader.Comma = ','
	reader.FieldsPerRecord = -1

	// override to standarize date string (20211123)
	date := helper.UnifiedDateFormatToTwse(s.date)

	for {
		records, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if len(records) == 0 || s.capacity > len(records) {
			continue
		}

		// make sure only parse recognized stock_id
		records[0] = strings.TrimSpace(records[0])
		if len(records[0]) > 0 && len(records[0]) < 6 && helper.IsInteger(records[0][0:2]) {
			res := s.converter.Execute(&convert.Data{
				ParseDate: date,
				RawData:   records,
				Target:    s.source,
			})
			if res != nil {
				output = append(output, res)
			}
		}
	}

	if len(output) == 0 {
		return nil, ErrNoParseResults
	}

	return output, nil
}
