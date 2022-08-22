package crawler

import (
	"context"

	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
	"github.com/samwang0723/stock-crawler/internal/app/pipeline"

	"golang.org/x/xerrors"
)

const (
	stakeConcentrationTotalCount = 5
)

type broadcastor struct {
	interceptChan chan convert.InterceptData
	memCache      map[string][]*entity.StakeConcentration
}

func newBroadcastor() *broadcastor {
	return &broadcastor{
		memCache: make(map[string][]*entity.StakeConcentration),
	}
}

func (b *broadcastor) InterceptData(ctx context.Context, interceptChan chan convert.InterceptData) {
	b.interceptChan = interceptChan
}

func (b *broadcastor) Process(ctx context.Context, pipe pipeline.Payload) (pipeline.Payload, error) {
	payload, ok := pipe.(*crawlerPayload)
	if !ok {
		return nil, xerrors.New("invalid payload")
	}

	intercept := convert.InterceptData{}

	if payload.Strategy == convert.StakeConcentration {
		if st := b.cacheInMemory(payload.ParsedContent); st != nil {
			intercept = convert.InterceptData{
				Data: &[]interface{}{st},
				Type: payload.Strategy,
			}
		}
	} else {
		intercept = convert.InterceptData{
			Data: payload.ParsedContent,
			Type: payload.Strategy,
		}
	}

	if b.interceptChan != nil && intercept.Data != nil {
		b.interceptChan <- intercept
	}

	return pipe, nil
}

func (b *broadcastor) cacheInMemory(data *[]interface{}) *entity.StakeConcentration {
	for _, v := range *data {
		if val, ok := v.(*entity.StakeConcentration); ok {
			b.memCache[val.StockID] = append(b.memCache[val.StockID], val)

			if len(b.memCache[val.StockID]) == stakeConcentrationTotalCount {
				output := entity.MapReduceStakeConcentration(b.memCache[val.StockID])
				delete(b.memCache, val.StockID)

				return output
			}
		}
	}

	return nil
}
