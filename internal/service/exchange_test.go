package service

import (
	"context"
	"errors"
	"testing"
	"time"

	gapimock "garantexGRPC/internal/garantex_api/mocks"
	repomock "garantexGRPC/internal/repository/mocks"
	"garantexGRPC/models"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
)

type webmockBehavior func(m *gapimock.MockGarantexer)
type storagemockBehaivor func(m *repomock.MockStorager)

var (
	ctx = context.Background()
)

func TestExchange_GetAndSaveRates(t *testing.T) {
	market := "usdt"
	rates := models.Rates{
		Timestamp: time.Unix(100000, 0),
		AskPrice:  decimal.NewFromFloat(100.200),
		BidPrice:  decimal.NewFromFloat(300.400),
	}
	tests := []struct {
		name        string
		webMock     webmockBehavior
		storageMock storagemockBehaivor
		want        models.Rates
		wantErr     bool
	}{
		{
			name: "test ok",
			webMock: func(m *gapimock.MockGarantexer) {
				m.EXPECT().GetRates(gomock.Any(), market).Return(rates, nil)
			},
			storageMock: func(m *repomock.MockStorager) {
				m.EXPECT().SaveRates(gomock.Any(), rates).Return(nil)
			},
			want:    rates,
			wantErr: false,
		},
		{
			name: "test get rates error",
			webMock: func(m *gapimock.MockGarantexer) {
				m.EXPECT().GetRates(gomock.Any(), market).Return(models.Rates{}, errors.New("get rates from web error"))
			},
			storageMock: func(m *repomock.MockStorager) {},
			want:        models.Rates{},
			wantErr:     true,
		},
		{
			name: "test save rates error",
			webMock: func(m *gapimock.MockGarantexer) {
				m.EXPECT().GetRates(gomock.Any(), market).Return(rates, nil)
			},
			storageMock: func(m *repomock.MockStorager) {
				m.EXPECT().SaveRates(gomock.Any(), rates).Return(errors.New("save rates to db error"))
			},
			want:    models.Rates{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()
			mockWebCli := gapimock.NewMockGarantexer(c)
			mockStorage := repomock.NewMockStorager(c)

			e := NewExchange(mockWebCli, mockStorage)
			tt.webMock(mockWebCli)
			tt.storageMock(mockStorage)

			got, err := e.GetRates(ctx, market)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAndSaveRates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAndSaveRates() got = %v, want %v", got, tt.want)
			}
		})
	}
}
