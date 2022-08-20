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
package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/samwang0723/stock-crawler/internal/app/crawler"
	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/helper"
	"github.com/samwang0723/stock-crawler/internal/kafka"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/xerrors"
)

const (
	sourceStockList    = "./configs/stock_ids.json"
	defaultCacheExpire = 6 * time.Hour
)

//nolint:nolintlint, gochecknoglobals
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (s *serviceImpl) DailyCloseThroughKafka(ctx context.Context, objs *[]interface{}) error {
	for _, val := range *objs {
		if res, ok := val.(*entity.DailyClose); ok {
			b, err := json.Marshal(res)
			if err != nil {
				return xerrors.Errorf("DailyCloseThroughKafka: json.Marshal failed: %w", err)
			}

			err = s.sendKafka(ctx, kafka.DailyClosesV1, b)
			if err != nil {
				return xerrors.Errorf("DailyCloseThroughKafka: sendKafka failed: %w", err)
			}
		} else {
			return xerrors.Errorf("Cannot cast interface to *dto.DailyClose: %v", reflect.TypeOf(val).Elem())
		}
	}

	return nil
}

func (s *serviceImpl) StockThroughKafka(ctx context.Context, objs *[]interface{}) error {
	for _, val := range *objs {
		if res, ok := val.(*entity.Stock); ok {
			b, err := json.Marshal(res)
			if err != nil {
				return xerrors.Errorf("StockThroughKafka: json.Marshal failed: %w", err)
			}

			err = s.sendKafka(ctx, kafka.StocksV1, b)
			if err != nil {
				return xerrors.Errorf("StockThroughKafka: sendKafka failed: %w", err)
			}
		} else {
			return xerrors.Errorf("Cannot cast interface to *dto.Stock: %v", reflect.TypeOf(val).Elem())
		}
	}

	return nil
}

func (s *serviceImpl) ThreePrimaryThroughKafka(ctx context.Context, objs *[]interface{}) error {
	for _, val := range *objs {
		if res, ok := val.(*entity.ThreePrimary); ok {
			b, err := json.Marshal(res)
			if err != nil {
				return xerrors.Errorf("ThreePrimaryThroughKafka: json.Marshal failed: %w", err)
			}

			err = s.sendKafka(ctx, kafka.ThreePrimaryV1, b)
			if err != nil {
				return xerrors.Errorf("ThreePrimaryThroughKafka: sendKafka failed: %w", err)
			}
		} else {
			return xerrors.Errorf("Cannot cast interface to *dto.ThreePrimary: %v", reflect.TypeOf(val).Elem())
		}
	}

	return nil
}

func (s *serviceImpl) StakeConcentrationThroughKafka(ctx context.Context, objs *[]interface{}) error {
	for _, val := range *objs {
		if res, ok := val.(*entity.StakeConcentration); ok {
			b, err := json.Marshal(res)
			if err != nil {
				return xerrors.Errorf("StakeConcentrationThroughKafka: json.Marshal failed: %w", err)
			}

			err = s.sendKafka(ctx, kafka.StakeConcentrationV1, b)
			if err != nil {
				return xerrors.Errorf("StakeConcentrationThroughKafka: sendKafka failed: %w", err)
			}

			// record parsed records to prevent duplicate parsing, default expire the key after 6 hours
			err = s.cacheParsedConcentration(ctx, res.Date, res.StockID)
			if err != nil {
				return xerrors.Errorf("StakeConcentrationThroughKafka: cacheParsedConcentration failed: %w", err)
			}

			res.Recycle()
		} else {
			return xerrors.Errorf("Cannot cast interface to *dto.StakeConcentration: %v", reflect.TypeOf(val).Elem())
		}
	}

	return nil
}

func (s *serviceImpl) cacheParsedConcentration(ctx context.Context, date, stockID string) error {
	key := strings.ReplaceAll(date, "-", "")

	err := s.cache.SAdd(ctx, key, stockID)
	if err != nil {
		return xerrors.Errorf("cacheParsedConcentration: cache.SAdd failed: %w", err)
	}

	err = s.cache.SetExpire(ctx, key, time.Now().Add(defaultCacheExpire))
	if err != nil {
		return xerrors.Errorf("cacheParsedConcentration: cache.SetExpire failed: %w", err)
	}

	return nil
}

func (s *serviceImpl) ListCrawlingConcentrationURLs(ctx context.Context, date string) ([]string, error) {
	defaultList, err := listStocks()
	if err != nil {
		return nil, xerrors.Errorf("ListCrawlingConcentrationURLs: listStocks failed: %w", err)
	}

	res, err := s.cache.SMembers(ctx, date)
	if err != nil {
		return nil, xerrors.Errorf("ListCrawlingConcentrationURLs: redis(SMembers) failed: %w", err)
	}

	var urls []string

	stockIds := helper.Diff(res, defaultList)

	for _, sid := range stockIds {
		// in order to get accurate data, we must query each page
		// https://stockchannelnew.sinotrade.com.tw/z/zc/zco/zco_6598_6.djhtm
		// as the top 15 brokers may different from day to day and not possible to store all detailed daily data
		indexes := []int{1, 2, 3, 4, 6}
		for _, idx := range indexes {
			urls = append(urls, fmt.Sprintf(crawler.ConcentrationDays, sid, idx))
		}
	}

	return urls, nil
}

func listStocks() ([]string, error) {
	// Open stock list jsonFile
	jsonFile, err := os.Open(sourceStockList)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, xerrors.Errorf("listStocks: os.Open failed: %w", err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, xerrors.Errorf("listStocks: io.ReadAll failed: %w", err)
	}

	// decode json stock list
	var list struct {
		StockIds []string `json:"stockIds"`
	}

	err = json.Unmarshal(byteValue, &list)
	if err != nil {
		return nil, xerrors.Errorf("listStocks: json.Unmarshal failed: %w", err)
	}

	return list.StockIds, nil
}
