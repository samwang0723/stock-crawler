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
package crawler

import (
	"context"
	"fmt"

	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/app/parser"
	"github.com/samwang0723/stock-crawler/internal/app/pipeline"
)

type textExtractor struct {
	parser   parser.Parser
	memCache map[string][]*entity.StakeConcentration
}

func newTextExtractor(parser parser.Parser) *textExtractor {
	return &textExtractor{
		parser:   parser,
		memCache: make(map[string][]*entity.StakeConcentration),
	}
}

func (te *textExtractor) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	payload := p.(*crawlerPayload)
	te.parser.SetStrategy(payload.Strategy, payload.Date)
	// Bypass the parsing error
	te.parser.Execute(payload.RawContent, payload.URL)

	//TODO: do switch case to parse different types
	objs := te.parser.Flush()
	for _, v := range *objs {
		if val, ok := v.(*entity.StakeConcentration); ok {
			te.memCache[val.StockID] = append(te.memCache[val.StockID], val)
			fmt.Printf("Count: %d, Process: %+v\n", len(te.memCache[val.StockID]), val)
			if len(te.memCache[val.StockID]) == 5 {
				//TODO: send through Kafka channel
			}
		}
	}
	return p, nil
}
