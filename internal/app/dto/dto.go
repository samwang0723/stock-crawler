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
package dto

type DownloadType int32

//go:generate stringer -type=DownloadType
const (
	DailyClose DownloadType = iota
	ThreePrimary
	Concentration
	StockList
)

// Payload is implemented by values that can be sent through a pipeline.
type Payload interface {
	// Clone returns a new Payload that is a deep-copy of the original.
	Clone() Payload

	// MarkAsProcessed is invoked by the consumer when the Payload either
	// reaches the parser or it gets discarded by one of the
	// consumer.
	MarkAsProcessed()
}

type DownloadRequest struct {
	UTCTimestamp string         `json:"utcTimestamp"`
	Types        []DownloadType `json:"types"`
	RewindLimit  int32          `json:"rewindLimit"`
	RateLimit    int32          `json:"rateLimit"`
}

type StartCronjobRequest struct {
	Schedule string         `json:"schedule"`
	Types    []DownloadType `json:"types"`
}
