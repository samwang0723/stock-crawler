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
	"errors"
	"io"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"golang.org/x/xerrors"
)

const (
	StockCap            = 5
	DailyCloseCap       = 17
	TwseThreePrimaryCap = 19
	TpexThreePrimaryCap = 24
	ConcentrationCap    = 7
)

type Parser interface {
	SetStrategy(source convert.Source, additional ...string)
	Execute(in bytes.Buffer, additional ...string) error
	Flush() *[]any
}

type Strategy interface {
	Parse(in io.Reader, additional ...string) ([]any, error)
}

type Config struct {
	Logger *zerolog.Logger
}

type parserImpl struct {
	cfg      Config
	strategy Strategy
	result   *[]any
}

func New(cfg Config) Parser {
	res := &parserImpl{
		cfg:    cfg,
		result: &[]any{},
	}

	return res
}

func (p *parserImpl) SetStrategy(source convert.Source, additional ...string) {
	switch source {
	case convert.TpexStockList, convert.TwseStockList:
		p.strategy = &htmlStrategy{
			capacity:  StockCap,
			source:    source,
			converter: convert.Stock(),
		}
	case convert.TwseDailyClose, convert.TpexDailyClose:
		p.strategy = &csvStrategy{
			capacity:  DailyCloseCap,
			source:    source,
			converter: convert.DailyClose(),
			date:      additional[0],
		}
	case convert.TwseThreePrimary:
		p.strategy = &csvStrategy{
			capacity:  TwseThreePrimaryCap,
			source:    source,
			converter: convert.ThreePrimary(),
			date:      additional[0],
		}
	case convert.TpexThreePrimary:
		p.strategy = &csvStrategy{
			capacity:  TpexThreePrimaryCap,
			source:    source,
			converter: convert.ThreePrimary(),
			date:      additional[0],
		}
	case convert.StakeConcentration:
		p.strategy = &concentrationStrategy{
			capacity:  ConcentrationCap,
			date:      additional[0],
			converter: convert.Concentration(),
		}
	}
}

func (p *parserImpl) Execute(in bytes.Buffer, additional ...string) error {
	reader := transform.NewReader(&in, traditionalchinese.Big5.NewDecoder())

	res, err := p.strategy.Parse(reader, additional...)
	if err != nil {
		// here we treat empty content as not an error, still continue
		if errors.Is(err, ErrNoParseResults) || errors.Is(err, ErrWrongConcentrationTitle) {
			log.Error().Err(err).Msg("parser.Execute: failed + continue")

			return nil
		}

		return xerrors.Errorf("parser.Execute: failed, err=%w;", err)
	}

	*p.result = append(*p.result, res...)

	return nil
}

func (p *parserImpl) Flush() *[]any {
	res := *p.result
	p.result = &[]any{}

	return &res
}
