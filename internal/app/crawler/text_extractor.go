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

	"github.com/samwang0723/stock-crawler/internal/app/parser"
	"github.com/samwang0723/stock-crawler/internal/app/pipeline"
)

type textExtractor struct {
	parser parser.Parser
}

func newTextExtractor(parser parser.Parser) *textExtractor {
	return &textExtractor{
		parser: parser,
	}
}

func (te *textExtractor) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	payload := p.(*crawlerPayload)
	te.parser.SetStrategy(payload.Strategy, payload.Date)
	// Bypass the parsing error
	te.parser.Execute(payload.RawContent, payload.URL)

	return p, nil
}
