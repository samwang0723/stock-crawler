package crawler

import (
	"context"
	"testing"

	"github.com/samwang0723/stock-crawler/internal/app/entity/convert"
)

func TestInterceptData(t *testing.T) {
	t.Parallel()

	type args struct {
		mockChan chan convert.InterceptData
	}

	tests := []struct {
		name string
		args args
		isEq bool
	}{
		{
			name: "configure intercept data channel",
			args: args{
				mockChan: make(chan convert.InterceptData),
			},
			isEq: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := newBroadcastor()
			b.InterceptData(context.Background(), tt.args.mockChan)
			if (b.interceptChan == tt.args.mockChan) != tt.isEq {
				t.Errorf("InterceptData() set channel faield: %v != %v", b.interceptChan, tt.args.mockChan)
			}
		})
	}
}
