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
package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/helper"
	log "github.com/samwang0723/stock-crawler/internal/logger"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (s *serviceImpl) DailyCloseThroughKafka(ctx context.Context, objs *[]interface{}) error {
	for _, v := range *objs {
		if val, ok := v.(*entity.DailyClose); ok {
			b, err := json.Marshal(val)
			if err != nil {
				return fmt.Errorf("DeliverStocks: json.Marshal failed: %w", err)
			}
			s.sendKafka(ctx, b)
		} else {
			return fmt.Errorf("Cannot cast interface to *dto.Stock: %v\n", reflect.TypeOf(v).Elem())
		}
	}
	return nil

}

func (s *serviceImpl) StockThroughKafka(ctx context.Context, objs *[]interface{}) error {
	for _, v := range *objs {
		if val, ok := v.(*entity.Stock); ok {
			b, err := json.Marshal(val)
			if err != nil {
				return fmt.Errorf("DeliverStocks: json.Marshal failed: %w", err)
			}
			s.sendKafka(ctx, b)
		} else {
			return fmt.Errorf("Cannot cast interface to *dto.Stock: %v\n", reflect.TypeOf(v).Elem())
		}
	}
	return nil
}

func (s *serviceImpl) ThreePrimaryThroughKafka(ctx context.Context, objs *[]interface{}) error {
	for _, v := range *objs {
		if val, ok := v.(*entity.ThreePrimary); ok {
			b, err := json.Marshal(val)
			if err != nil {
				return fmt.Errorf("DeliverStocks: json.Marshal failed: %w", err)
			}
			s.sendKafka(ctx, b)
		} else {
			return fmt.Errorf("Cannot cast interface to *dto.Stock: %v\n", reflect.TypeOf(v).Elem())
		}
	}
	return nil
}

func (s *serviceImpl) StakeConcentrationThroughKafka(ctx context.Context, objs *[]interface{}) error {
	var redisKey string
	for _, v := range *objs {
		if val, ok := v.(*entity.StakeConcentration); ok {
			b, err := json.Marshal(val)
			if err != nil {
				return fmt.Errorf("DeliverStocks: json.Marshal failed: %w", err)
			}
			s.sendKafka(ctx, b)

			if len(redisKey) == 0 {
				redisKey = strings.ReplaceAll(val.Date, "-", "")
			}
			// record parsed records to prevent duplicate parsing
			s.cache.LPush(ctx, redisKey, val.StockID)
		} else {
			return fmt.Errorf("Cannot cast interface to *dto.Stock: %v\n", reflect.TypeOf(v).Elem())
		}
	}
	return nil
}

func (s *serviceImpl) ListBackfillStakeConcentrationStockIds(ctx context.Context, date string) ([]string, error) {
	defaultList, err := loadStockList()
	if err != nil {
		return nil, err
	}
	res, err := s.cache.LRange(ctx, date)
	if err != nil {
		return nil, err
	}

	return helper.Difference(res, defaultList), nil
}

func loadStockList() ([]string, error) {
	// Open stock list jsonFile
	jsonFile, err := os.Open("./configs/stock_ids.test.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	log.Info("Successfully Opened stock_ids.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// decode json stock list
	var list struct {
		StockIds []string `json:"stockIds"`
	}
	err = json.Unmarshal(byteValue, &list)
	if err != nil {
		return nil, err
	}
	return list.StockIds, nil
}
