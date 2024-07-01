package services

import (
	"context"
	"errors"
	"flag"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	jsoniter "github.com/json-iterator/go"
	"github.com/samwang0723/stock-crawler/internal/app/entity"
	cache "github.com/samwang0723/stock-crawler/internal/cache/mocks"
	"github.com/samwang0723/stock-crawler/internal/kafka"
	kafkamock "github.com/samwang0723/stock-crawler/internal/kafka/mocks"
	"go.uber.org/goleak"
	"golang.org/x/xerrors"
)

//nolint:nolintlint, gochecknoglobals
var jsonTest = jsoniter.ConfigCompatibleWithStandardLibrary
var ErrFailed = errors.New("failed")

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestDailyCloseThroughKafka(t *testing.T) {
	t.Parallel()

	type args struct {
		data         *[]interface{}
		expectReturn error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successfully send daily close data through kafka",
			args: args{
				data: &[]interface{}{
					&entity.DailyClose{
						StockID:      "2330",
						Date:         "20200101",
						TradedShares: 100000,
						Transactions: 100000,
						Turnover:     100000,
						Open:         23.0,
						Close:        24.0,
						High:         25.0,
						Low:          23.0,
						PriceDiff:    0.0,
					},
				},
				expectReturn: nil,
			},
			wantErr: false,
		},
		{
			name: "failed to send correct data format through kafka",
			args: args{
				data: &[]interface{}{
					&entity.StakeConcentration{
						StockID: "2330",
					},
				},
				expectReturn: nil,
			},
			wantErr: true,
		},
		{
			name: "failed to send daily close data due to kafka error",
			args: args{
				data: &[]interface{}{
					&entity.DailyClose{
						StockID: "2330",
					},
				},
				expectReturn: xerrors.Errorf("kafka writeMessages(): %w", ErrFailed),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			mockKafka := kafkamock.NewMockKafka(mockCtl)

			for _, val := range *tt.args.data {
				if res, ok := val.(*entity.DailyClose); ok {
					b, err := jsonTest.Marshal(res)
					if err != nil {
						t.Errorf("service DailyCloseThroughKafka: jsonTest.Marshal failed: %v", err)
					}
					mockKafka.EXPECT().
						WriteMessages(ctx, kafka.DailyClosesV1, b).
						Return(tt.args.expectReturn).
						Times(1)
				}
			}

			svc := &serviceImpl{
				producer: mockKafka,
			}

			err := svc.DailyCloseThroughKafka(ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("service DailyCloseThroughKafka() error = %v", err)
			}
		})
	}
}

func TestStockThroughKafka(t *testing.T) {
	t.Parallel()

	type args struct {
		data         *[]interface{}
		expectReturn error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successfully send stock data through kafka",
			args: args{
				data: &[]interface{}{
					&entity.Stock{
						StockID: "2330",
						Name:    "Test",
					},
				},
				expectReturn: nil,
			},
			wantErr: false,
		},
		{
			name: "failed to send correct data format through kafka",
			args: args{
				data: &[]interface{}{
					&entity.StakeConcentration{
						StockID: "2330",
					},
				},
				expectReturn: nil,
			},
			wantErr: true,
		},
		{
			name: "failed to send stock data due to kafka error",
			args: args{
				data: &[]interface{}{
					&entity.Stock{
						StockID: "2330",
						Name:    "Test",
					},
				},
				expectReturn: xerrors.Errorf("kafka writeMessages(): %w", ErrFailed),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			mockKafka := kafkamock.NewMockKafka(mockCtl)

			for _, val := range *tt.args.data {
				if res, ok := val.(*entity.Stock); ok {
					b, err := jsonTest.Marshal(res)
					if err != nil {
						t.Errorf("service StockThroughKafka: jsonTest.Marshal failed: %v", err)
					}
					mockKafka.EXPECT().
						WriteMessages(ctx, kafka.StocksV1, b).
						Return(tt.args.expectReturn).
						Times(1)
				}
			}

			svc := &serviceImpl{
				producer: mockKafka,
			}

			err := svc.StockThroughKafka(ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("service StockThroughKafka() error = %v", err)
			}
		})
	}
}

func TestThreePrimaryThroughKafka(t *testing.T) {
	t.Parallel()

	type args struct {
		data         *[]interface{}
		expectReturn error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successfully send three primary data through kafka",
			args: args{
				data: &[]interface{}{
					&entity.ThreePrimary{
						StockID: "2330",
					},
				},
				expectReturn: nil,
			},
			wantErr: false,
		},
		{
			name: "failed to send correct data format through kafka",
			args: args{
				data: &[]interface{}{
					&entity.StakeConcentration{
						StockID: "2330",
					},
				},
				expectReturn: nil,
			},
			wantErr: true,
		},
		{
			name: "failed to send three primary due to kafka error",
			args: args{
				data: &[]interface{}{
					&entity.ThreePrimary{
						StockID: "2330",
					},
				},
				expectReturn: xerrors.Errorf("kafka writeMessages(): %w", ErrFailed),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			mockKafka := kafkamock.NewMockKafka(mockCtl)

			for _, val := range *tt.args.data {
				if res, ok := val.(*entity.ThreePrimary); ok {
					b, err := jsonTest.Marshal(res)
					if err != nil {
						t.Errorf(
							"service ThreePrimaryThroughKafka: jsonTest.Marshal failed: %v",
							err,
						)
					}
					mockKafka.EXPECT().
						WriteMessages(ctx, kafka.ThreePrimaryV1, b).
						Return(tt.args.expectReturn).
						Times(1)
				}
			}

			svc := &serviceImpl{
				producer: mockKafka,
			}

			err := svc.ThreePrimaryThroughKafka(ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("service ThreePrimaryThroughKafka() error = %v", err)
			}
		})
	}
}

func TestStakeConcentrationThroughKafka(t *testing.T) {
	t.Parallel()

	type args struct {
		data         *[]interface{}
		date         string
		expectReturn error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successfully send stake concentration data through kafka",
			args: args{
				data: &[]interface{}{
					&entity.StakeConcentration{
						StockID: "2330",
						Date:    "20220820",
					},
				},
				date:         "20220820",
				expectReturn: nil,
			},
			wantErr: false,
		},
		{
			name: "failed to send correct data format through kafka",
			args: args{
				data: &[]interface{}{
					&entity.Stock{
						StockID: "2330",
					},
				},
				date:         "20220820",
				expectReturn: nil,
			},
			wantErr: true,
		},
		{
			name: "failed to send stake concentration due to kafka error",
			args: args{
				data: &[]interface{}{
					&entity.StakeConcentration{
						StockID: "2330",
						Date:    "20220820",
					},
				},
				date:         "20220820",
				expectReturn: xerrors.Errorf("kafka writeMessages(): %w", ErrFailed),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			mockCtl := gomock.NewController(t)
			defer mockCtl.Finish()

			mockKafka := kafkamock.NewMockKafka(mockCtl)
			mockRedis := cache.NewMockRedis(mockCtl)

			for _, val := range *tt.args.data {
				if res, ok := val.(*entity.StakeConcentration); ok {
					b, err := jsonTest.Marshal(res)
					if err != nil {
						t.Errorf(
							"service StakeConcentrationThroughKafka: jsonTest.Marshal failed: %v",
							err,
						)
					}
					mockKafka.EXPECT().
						WriteMessages(ctx, kafka.StakeConcentrationV1, b).
						Return(tt.args.expectReturn).
						Times(1)
					mockRedis.EXPECT().SAdd(ctx, res.Date, res.StockID).Return(nil).AnyTimes()
					mockRedis.EXPECT().
						SetExpire(ctx, res.Date, gomock.AssignableToTypeOf(time.Now())).
						Return(nil).
						AnyTimes()
				}
			}

			svc := &serviceImpl{
				producer: mockKafka,
				cache:    mockRedis,
			}

			err := svc.StakeConcentrationThroughKafka(ctx, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("service StakeConcentrationThroughKafka() error = %v", err)
			}
		})
	}
}
