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
	"io"

	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"

	"github.com/rs/zerolog"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type Parser interface {
	SetStrategy(source convert.Source, additional ...string)
	Execute(in bytes.Buffer, additional ...string) error
	Flush() *[]interface{}
}

type Strategy interface {
	Parse(in io.Reader, additional ...string) ([]interface{}, error)
}

type Config struct {
	Logger zerolog.Logger
}

type parserImpl struct {
	cfg      Config
	strategy Strategy
	result   *[]interface{}
}

func New(cfg Config) Parser {
	res := &parserImpl{
		cfg:    cfg,
		result: &[]interface{}{},
	}
	return res
}

func (p *parserImpl) SetStrategy(source convert.Source, additional ...string) {
	switch source {
	case convert.TpexStockList, convert.TwseStockList:
		p.strategy = &htmlStrategy{
			capacity:  6,
			source:    source,
			converter: convert.Stock(),
		}
	case convert.TwseDailyClose, convert.TpexDailyClose:
		p.strategy = &csvStrategy{
			capacity:  17,
			source:    source,
			converter: convert.DailyClose(),
			date:      additional[0],
		}
	case convert.TwseThreePrimary:
		p.strategy = &csvStrategy{
			capacity:  19,
			source:    source,
			converter: convert.ThreePrimary(),
			date:      additional[0],
		}
	case convert.TpexThreePrimary:
		p.strategy = &csvStrategy{
			capacity:  24,
			source:    source,
			converter: convert.ThreePrimary(),
			date:      additional[0],
		}
	case convert.StakeConcentration:
		p.strategy = &concentrationStrategy{
			capacity:  7,
			date:      additional[0],
			converter: convert.Concentration(),
		}
	}
}

func (p *parserImpl) Execute(in bytes.Buffer, additional ...string) error {
	reader := transform.NewReader(&in, traditionalchinese.Big5.NewDecoder())
	res, err := p.strategy.Parse(reader, additional...)
	if err != nil {
		return err
	}

	*p.result = append(*p.result, res...)
	return nil
}

func (p *parserImpl) Flush() *[]interface{} {
	res := *p.result
	p.result = &[]interface{}{}
	return &res
}
