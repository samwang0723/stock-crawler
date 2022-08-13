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
	"strconv"

	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/app/parser"
	"github.com/samwang0723/stock-crawler/internal/app/pipeline"

	"github.com/sirupsen/logrus"
)

type textExtractor struct {
	parser   parser.Parser
	logger   *logrus.Entry
	memCache map[string][]*entity.StakeConcentration
}

func newTextExtractor(parser parser.Parser, logger *logrus.Entry) *textExtractor {
	return &textExtractor{
		parser:   parser,
		logger:   logger,
		memCache: make(map[string][]*entity.StakeConcentration),
	}
}

func (te *textExtractor) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	payload := p.(*crawlerPayload)
	te.parser.SetStrategy(payload.Strategy, payload.Date)
	err := te.parser.Execute(payload.RawContent, payload.URL)
	if err != nil {
		return nil, err
	}

	err = te.broadcast(ctx, payload.Strategy, te.parser.Flush())
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (te *textExtractor) broadcast(ctx context.Context, strategy convert.Source, data *[]interface{}) error {
	switch strategy {
	case convert.TwseStockList, convert.TpexStockList:
	case convert.TwseDailyClose, convert.TpexDailyClose:
	case convert.TwseThreePrimary, convert.TpexThreePrimary:
	case convert.StakeConcentration:
		for _, v := range *data {
			if val, ok := v.(*entity.StakeConcentration); ok {
				te.memCache[val.StockID] = append(te.memCache[val.StockID], val)
				if len(te.memCache[val.StockID]) == 5 {
					te.logger.Infof("Count: %d, Process: %+v", len(te.memCache[val.StockID]), te.mapReduce(ctx, te.memCache[val.StockID]))
					//TODO: clear in-memory cache
				}
			}
		}
	}

	return nil
}

func (te *textExtractor) mapReduce(ctx context.Context, objs []*entity.StakeConcentration) *entity.StakeConcentration {
	volumeDiff := []int32{0, 0, 0, 0, 0}
	var res *entity.StakeConcentration

	for _, val := range objs {
		idx, _ := strconv.Atoi(val.HiddenField)
		// make sure to cover latest concentration data
		if idx == 0 {
			res = val
		}
		volumeDiff[idx] = int32(val.SumBuyShares - val.SumSellShares)
	}

	res.Diff = volumeDiff

	return res
}
