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
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/samwang0723/stock-crawler/internal/app/dto"
)

var (
	_ dto.Payload = (*crawlerPayload)(nil)

	payloadPool = sync.Pool{
		New: func() interface{} { return new(crawlerPayload) },
	}
)

type crawlerPayload struct {
	URL         string
	RetrievedAt time.Time
	RawContent  bytes.Buffer
}

func (p *crawlerPayload) Clone() dto.Payload {
	newP := payloadPool.Get().(*crawlerPayload)
	newP.URL = p.URL
	newP.RetrievedAt = p.RetrievedAt

	_, err := io.Copy(&newP.RawContent, &p.RawContent)
	if err != nil {
		panic(fmt.Sprintf("error cloning payload raw content: %v", err))
	}
	return newP
}

func (p *crawlerPayload) MarkAsProcessed() {
	p.URL = p.URL[:0]
	p.RawContent.Reset()

	payloadPool.Put(p)
}
