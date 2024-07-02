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
package handlers

import (
	"sync"

	"github.com/samwang0723/stock-crawler/internal/app/graph"
)

type linkIterator struct {
	links    []*graph.Link
	curIndex int
	mu       sync.RWMutex
}

// Next implements graph.LinkIterator.
func (i *linkIterator) Next() bool {
	if i.curIndex >= len(i.links) {
		return false
	}
	i.curIndex++

	return true
}

// Error implements graph.LinkIterator.
func (i *linkIterator) Error() error {
	return nil
}

// Link implements graph.LinkIterator.
func (i *linkIterator) Link() *graph.Link {
	// The link pointer contents may be overwritten by external update; to
	// avoid data-races we acquire the read lock first and clone the link
	i.mu.RLock()
	link := new(graph.Link)
	*link = *i.links[i.curIndex-1]
	i.mu.RUnlock()

	return link
}
