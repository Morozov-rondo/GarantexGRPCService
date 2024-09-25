package garantex_api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"garantexGRPC/configs"
	"garantexGRPC/models"
	"github.com/shopspring/decimal"
)

const (
	testResponse = `{
  "timestamp": 1648047486,
  "asks": [
    {
      "price": "156612.04",
      "volume": "0.23625259",
      "amount": "37000.0",
      "factor": "0.008",
      "type": "limit"      
    }
  ],
  "bids": [
    {
      "price": "153590.1",
      "volume": "0.04143581",
      "amount": "6364.13",
      "factor": "-0.011",
      "type": "factor"
    }
  ]
}`
	invalidResponse = `{
  "timestamp": 1648047486,
  "asks": [
],
  "bids": [
    {
      "price": "153590.1",
      "volume": "0.04143581",
      "amount": "6364.13",
      "factor": "-0.011",
      "type": "factor"
    }
  ]
}`
)

var (
	ctx = context.Background()
)

func TestGarantex_GetRatesFromDepth(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Invalid http request method got: %s, want: %s", r.Method, http.MethodGet)
		}
		fmt.Fprintln(w, testResponse)
	}))
	defer mockServer.Close()

	ask, _ := decimal.NewFromString("156612.04")
	bid, _ := decimal.NewFromString("153590.1")

	tests := []struct {
		name    string
		market  string
		url     string
		want    models.Rates
		wantErr bool
	}{
		{
			name:   "200 OK",
			market: "usdtrub",
			url:    mockServer.URL,
			want: models.Rates{
				ID:        0,
				Market:    "usdtrub",
				Timestamp: time.Unix(1648047486, 0),
				AskPrice:  ask,
				BidPrice:  bid,
			},
			wantErr: false,
		},
		{
			name:    "new request error",
			market:  "usdtrub",
			url:     string([]byte{127}), // parse err: invalid control character in URL
			want:    models.Rates{},
			wantErr: true,
		},
		{
			name:    "do request error",
			market:  "usdtrub",
			url:     "some wrong URL", // error: unsupported protocol scheme
			want:    models.Rates{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			g := NewGarantex(configs.GarantexConfig{
				URL: tt.url,
			})

			got, err := g.GetRates(ctx, tt.market)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRatesFromDepth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRatesFromDepth() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGarantex_GetRatesFromDepth_ReadAllError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "1") //error: unexpected EOF
	}))
	defer mockServer.Close()

	tests := []struct {
		name    string
		query   string
		url     string
		want    models.Rates
		wantErr bool
	}{
		{
			name:    "read all error",
			query:   "usdtrub",
			url:     mockServer.URL,
			want:    models.Rates{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			g := NewGarantex(configs.GarantexConfig{
				URL: tt.url,
			})

			got, err := g.GetRates(ctx, tt.query)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRatesFromDepth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRatesFromDepth() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGarantex_GetRatesFromDepth_UnmarshalError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})) // error: unexpected end of JSON input
	defer mockServer.Close()

	tests := []struct {
		name    string
		query   string
		url     string
		want    models.Rates
		wantErr bool
	}{
		{
			name:    "unmarshal error",
			query:   "usdtrub",
			url:     mockServer.URL,
			want:    models.Rates{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			g := NewGarantex(configs.GarantexConfig{
				URL: tt.url,
			})

			got, err := g.GetRates(ctx, tt.query)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRatesFromDepth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRatesFromDepth() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGarantex_GetRatesFromDepth_ToDomainError(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, invalidResponse)
	}))
	defer mockServer.Close()

	tests := []struct {
		name    string
		market  string
		url     string
		want    models.Rates
		wantErr bool
	}{
		{
			name:    "to domain error",
			market:  "usdtrub",
			url:     mockServer.URL,
			want:    models.Rates{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			g := NewGarantex(configs.GarantexConfig{
				URL: tt.url,
			})

			got, err := g.GetRates(ctx, tt.market)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRatesFromDepth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRatesFromDepth() got = %v, want %v", got, tt.want)
			}
		})
	}
}
