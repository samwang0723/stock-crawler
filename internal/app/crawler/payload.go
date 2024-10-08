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
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/app/pipeline"
)

//nolint:nolintlint, gochecknoglobals
var (
	_ pipeline.Payload = (*crawlerPayload)(nil)

	payloadPool = sync.Pool{
		New: func() any { return new(crawlerPayload) },
	}
)

type crawlerPayload struct {
	RetrievedAt   time.Time
	ParsedContent *[]any
	URL           string
	Date          string
	RawContent    bytes.Buffer
	Strategy      convert.Source
}

func (p *crawlerPayload) Clone() pipeline.Payload {
	newP, ok := payloadPool.Get().(*crawlerPayload)
	if !ok {
		return nil
	}

	newP.URL = p.URL
	newP.Strategy = p.Strategy
	newP.Date = p.Date
	newP.RetrievedAt = p.RetrievedAt
	newP.ParsedContent = p.ParsedContent

	_, err := io.Copy(&newP.RawContent, &p.RawContent)
	if err != nil {
		log.Error().Err(err).Msg("payload.Clone: failed, reason: copy raw content failed")
		panic(fmt.Sprintf("payload.Clone: failed, err=%v;", err))
	}

	return newP
}

func (p *crawlerPayload) MarkAsProcessed() {
	p.URL = p.URL[:0]
	p.Date = p.Date[:0]
	p.Strategy = -1
	p.ParsedContent = nil
	p.RawContent.Reset()

	payloadPool.Put(p)
}
