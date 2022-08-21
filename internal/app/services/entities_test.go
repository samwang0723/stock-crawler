package services

import (
	"context"
	"flag"
	"os"
	"testing"

	"github.com/samwang0723/stock-crawler/internal/app/entity"
	"github.com/samwang0723/stock-crawler/internal/kafka"
	kafkamock "github.com/samwang0723/stock-crawler/internal/kafka/mocks"

	"github.com/golang/mock/gomock"
	"go.uber.org/goleak"
)

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
		data *[]interface{}
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
			},
			wantErr: false,
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
					b, err := json.Marshal(res)
					if err != nil {
						t.Errorf("service DailyCloseThroughKafka: json.Marshal failed: %v", err)
					}
					mockKafka.EXPECT().WriteMessages(ctx, kafka.DailyClosesV1, b).Return(nil).Times(1)
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
