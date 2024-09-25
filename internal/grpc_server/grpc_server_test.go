package grpc_server

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	exchangemock "garantexGRPC/internal/service/mocks"
	"garantexGRPC/models"
	garantex_sso_v1_ssov1 "garantexGRPC/protos/gen/go"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
)

type exchangeBehavior func(m *exchangemock.MockExchanger)

func TestExchange_GetRates(t *testing.T) {
	market := "usdtrub"
	rates := models.Rates{
		Timestamp: time.Unix(100000, 0),
		AskPrice:  decimal.NewFromFloat(100.200),
		BidPrice:  decimal.NewFromFloat(300.400),
	}
	tests := []struct {
		name         string
		exchangeMock exchangeBehavior
		req          *garantex_sso_v1_ssov1.GetRequest
		want         *garantex_sso_v1_ssov1.GetResponse
		wantErr      bool
	}{
		{
			name: "test ok",
			exchangeMock: func(m *exchangemock.MockExchanger) {
				m.EXPECT().GetRates(gomock.Any(), market).Return(rates, nil)
			},
			req: &garantex_sso_v1_ssov1.GetRequest{
				Market: garantex_sso_v1_ssov1.Market_usdtrub,
			},
			want: &garantex_sso_v1_ssov1.GetResponse{
				Timestamp: rates.Timestamp.Unix(),
				Market:    garantex_sso_v1_ssov1.Market_usdtrub,
				Ask:       rates.AskPrice.InexactFloat64(),
				Bid:       rates.BidPrice.InexactFloat64(),
			},
			wantErr: false,
		},
		{
			name: "test error",
			exchangeMock: func(m *exchangemock.MockExchanger) {
				m.EXPECT().GetRates(gomock.Any(), market).Return(models.Rates{}, errors.New("error get rates"))
			},
			req: &garantex_sso_v1_ssov1.GetRequest{
				Market: garantex_sso_v1_ssov1.Market_usdtrub,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			mockExchanger := exchangemock.NewMockExchanger(c)

			e := NewExchangeGRPC(mockExchanger)

			tt.exchangeMock(mockExchanger)

			got, err := e.Get(context.Background(), tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRates() got = %v, want %v", got, tt.want)
			}
		})
	}
}
