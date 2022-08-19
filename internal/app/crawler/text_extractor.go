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
package crawler

import (
	"context"

	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/app/parser"
	"github.com/samwang0723/stock-crawler/internal/app/pipeline"
	"golang.org/x/xerrors"
)

const (
	stakeConcentrationTotalCount = 5
)

type textExtractor struct {
	parser        parser.Parser
	memCache      map[string][]*entity.StakeConcentration
	interceptChan chan convert.InterceptData
}

func newTextExtractor(pa parser.Parser) *textExtractor {
	return &textExtractor{
		parser:   pa,
		memCache: make(map[string][]*entity.StakeConcentration),
	}
}

func (te *textExtractor) InterceptData(ctx context.Context, interceptChan chan convert.InterceptData) {
	te.interceptChan = interceptChan
}

func (te *textExtractor) Process(ctx context.Context, raw pipeline.Payload) (pipeline.Payload, error) {
	payload, ok := raw.(*crawlerPayload)
	if !ok {
		return nil, xerrors.New("invalid payload")
	}

	te.parser.SetStrategy(payload.Strategy, payload.Date)

	err := te.parser.Execute(payload.RawContent, payload.URL)
	if err != nil {
		return nil, xerrors.Errorf("parse error: %w", err)
	}

	te.broadcast(payload.Strategy, te.parser.Flush())

	return raw, nil
}

func (te *textExtractor) broadcast(strategy convert.Source, data *[]interface{}) {
	intercept := convert.InterceptData{}

	if strategy == convert.StakeConcentration {
		if st := te.cacheInMemory(data); st != nil {
			intercept = convert.InterceptData{
				Data: &[]interface{}{st},
				Type: strategy,
			}
		}
	} else {
		intercept = convert.InterceptData{
			Data: data,
			Type: strategy,
		}
	}

	if te.interceptChan != nil && intercept.Data != nil {
		te.interceptChan <- intercept
	}
}

func (te *textExtractor) cacheInMemory(data *[]interface{}) *entity.StakeConcentration {
	for _, v := range *data {
		if val, ok := v.(*entity.StakeConcentration); ok {
			te.memCache[val.StockID] = append(te.memCache[val.StockID], val)

			if len(te.memCache[val.StockID]) == stakeConcentrationTotalCount {
				output := entity.MapReduceStakeConcentration(te.memCache[val.StockID])
				delete(te.memCache, val.StockID)

				return output
			}
		}
	}

	return nil
}
