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
	"bytes"
	"fmt"

	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type Source int

//go:generate stringer -type=Source
const (
	TwseDailyClose Source = iota
	TwseThreePrimary
	TpexDailyClose
	TpexThreePrimary
	TwseStockList
	TpexStockList
	StakeConcentration
)

type IParser interface {
	Parse(config Config, in []byte) error
	Flush() *[]interface{}
}

type parserImpl struct {
	result *[]interface{}
}

type Config struct {
	ParseDay  *string
	SourceURL string
	Capacity  int
	Type      Source
}

func New() IParser {
	res := &parserImpl{
		result: &[]interface{}{},
	}
	return res
}

func (p *parserImpl) Parse(config Config, in []byte) error {
	if p.result == nil {
		return fmt.Errorf("didn't initialized the result map\n")
	}

	raw := bytes.NewBuffer(in)
	reader := transform.NewReader(raw, traditionalchinese.Big5.NewDecoder())

	switch config.Type {
	case TwseStockList, TpexStockList:
		return p.parseHtml(config, reader)
	case TwseDailyClose, TpexDailyClose, TwseThreePrimary, TpexThreePrimary:
		return p.parseCsv(config, reader)
	case StakeConcentration:
		return p.parseConcentration(config, reader)
	}
	return nil
}

func (p *parserImpl) Flush() *[]interface{} {
	res := *p.result
	p.result = &[]interface{}{}
	return &res
}
