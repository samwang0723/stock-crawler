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
package graph

import (
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
)

// Iterator is implemented by graph objects that can be iterated.
type Iterator interface {
	// Next advances the iterator. If no more items are available or an
	// error occurs, calls to Next() return bool.
	Next() bool

	// Error returns the last error encountered by the iterator.
	Error() error
}

// LinkIterator is implemented by objects that can iterate the graph links.
type LinkIterator interface {
	Iterator

	// Link returns the currently fetched link object.
	Link() *Link
}

// Link encapsulates all necessary setup used as source.
type Link struct {
	// The link target.
	URL string

	// Parsing target date
	Date string

	// Use which strategy for parsing
	Strategy convert.Source
}
