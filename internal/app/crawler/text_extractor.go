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

	"github.com/samwang0723/stock-crawler/internal/app/parser"
	"github.com/samwang0723/stock-crawler/internal/app/pipeline"

	"golang.org/x/xerrors"
)

type textExtractor struct {
	parser parser.Parser
}

func newTextExtractor(cfg Config) *textExtractor {
	return &textExtractor{
		parser: parser.New(parser.Config{Logger: cfg.Logger}),
	}
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

	payload.ParsedContent = te.parser.Flush()

	return payload, nil
}
