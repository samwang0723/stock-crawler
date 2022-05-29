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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/helper"
	"github.com/samwang0723/stock-crawler/internal/kafka/ikafka"
	log "github.com/samwang0723/stock-crawler/internal/logger"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (s *serviceImpl) DailyCloseThroughKafka(ctx context.Context, objs *[]interface{}) error {
	for _, v := range *objs {
		if val, ok := v.(*entity.DailyClose); ok {
			b, err := json.Marshal(val)
			if err != nil {
				return fmt.Errorf("DailyCloseThroughKafka: json.Marshal failed: %w", err)
			}
			err = s.sendKafka(ctx, ikafka.DailyClosesV1, b)
			if err != nil {
				return fmt.Errorf("DailyCloseThroughKafka: sendKafka failed: %w", err)
			}
		} else {
			return fmt.Errorf("Cannot cast interface to *dto.DailyClose: %v\n", reflect.TypeOf(v).Elem())
		}
	}
	return nil

}

func (s *serviceImpl) StockThroughKafka(ctx context.Context, objs *[]interface{}) error {
	for _, v := range *objs {
		if val, ok := v.(*entity.Stock); ok {
			b, err := json.Marshal(val)
			if err != nil {
				return fmt.Errorf("StockThroughKafka: json.Marshal failed: %w", err)
			}
			err = s.sendKafka(ctx, ikafka.StocksV1, b)
			if err != nil {
				return fmt.Errorf("StockThroughKafka: sendKafka failed: %w", err)
			}
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
				return fmt.Errorf("ThreePrimaryThroughKafka: json.Marshal failed: %w", err)
			}
			err = s.sendKafka(ctx, ikafka.ThreePrimaryV1, b)
			if err != nil {
				return fmt.Errorf("ThreePrimaryThroughKafka: sendKafka failed: %w", err)
			}
		} else {
			return fmt.Errorf("Cannot cast interface to *dto.ThreePrimary: %v\n", reflect.TypeOf(v).Elem())
		}
	}
	return nil
}

func (s *serviceImpl) StakeConcentrationThroughKafka(ctx context.Context, objs *[]interface{}) error {
	for _, v := range *objs {
		if val, ok := v.(*entity.StakeConcentration); ok {
			b, err := json.Marshal(val)
			if err != nil {
				return fmt.Errorf("StakeConcentrationThroughKafka: json.Marshal failed: %w", err)
			}
			err = s.sendKafka(ctx, ikafka.StakeConcentrationV1, b)
			if err != nil {
				return fmt.Errorf("StakeConcentrationThroughKafka: sendKafka failed: %w", err)
			}
			// record parsed records to prevent duplicate parsing, default expire the key after 6 hours
			key := strings.ReplaceAll(val.Date, "-", "")
			err = s.cache.SAdd(ctx, key, val.StockID)
			if err != nil {
				return fmt.Errorf("StakeConcentrationThroughKafka: redis(SAdd) failed: %w", err)
			}
			err = s.cache.SetExpire(ctx, key, time.Now().Add(6*time.Hour))
			if err != nil {
				return fmt.Errorf("StakeConcentrationThroughKafka: redis(SetExpure) failed: %w", err)
			}
		} else {
			return fmt.Errorf("Cannot cast interface to *dto.StakeConcentration: %v\n", reflect.TypeOf(v).Elem())
		}
	}
	return nil
}

func (s *serviceImpl) ListBackfillStakeConcentrationStockIds(ctx context.Context, date string) ([]string, error) {
	defaultList, err := loadStockList()
	if err != nil {
		return nil, err
	}
	res, err := s.cache.SMembers(ctx, date)
	if err != nil {
		return nil, err
	}

	return helper.Difference(res, defaultList), nil
}

func loadStockList() ([]string, error) {
	loc := "./configs/stock_ids.json"
	if helper.IsTesting() {
		loc = "./configs/stock_ids.test.json"
	}
	// Open stock list jsonFile
	jsonFile, err := os.Open(loc)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	log.Infof("Successfully Opened %s", loc)
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
