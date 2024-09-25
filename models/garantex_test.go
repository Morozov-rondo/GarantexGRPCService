package models

import (
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

func TestGrntxDepth_ToDomain(t *testing.T) {
	type fields struct {
		Timestamp int64
		Asks      []GarntxRate
		Bids      []GarntxRate
	}
	tests := []struct {
		name    string
		fields  fields
		want    Rates
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Timestamp: 100,
				Asks: []GarntxRate{{
					Price: "100.200",
				}},
				Bids: []GarntxRate{{
					Price: "300.400",
				}},
			},
			want: Rates{
				Timestamp: time.Unix(100, 0),
				AskPrice:  decimal.RequireFromString("100.200"),
				BidPrice:  decimal.RequireFromString("300.400"),
			},
			wantErr: false,
		},
		{
			name:    "valid error",
			fields:  fields{},
			want:    Rates{},
			wantErr: true,
		},
		{
			name: "ask error",
			fields: fields{
				Timestamp: 100,
				Asks: []GarntxRate{{
					Price: "invalid",
				}},
				Bids: []GarntxRate{{
					Price: "300.400",
				}},
			},
			want:    Rates{},
			wantErr: true,
		},
		{
			name: "bid error",
			fields: fields{
				Timestamp: 100,
				Asks: []GarntxRate{{
					Price: "100.200",
				}},
				Bids: []GarntxRate{{
					Price: "invalid",
				}},
			},
			want:    Rates{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &GarantxDepth{
				Timestamp: tt.fields.Timestamp,
				Asks:      tt.fields.Asks,
				Bids:      tt.fields.Bids,
			}
			got, err := d.ToDomain()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToDomain() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGrntxDepth_Valid(t *testing.T) {
	type fields struct {
		Timestamp int64
		Asks      []GarntxRate
		Bids      []GarntxRate
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "valid",
			fields: fields{
				Timestamp: 1,
				Asks:      make([]GarntxRate, 1),
				Bids:      make([]GarntxRate, 1),
			},
			want: true,
		},
		{
			name: "invalid timestamp",
			fields: fields{
				Timestamp: -1,
			},
			want: false,
		},
		{
			name: "invalid asks",
			fields: fields{
				Timestamp: 1,
				Asks:      nil,
			},
			want: false,
		},
		{
			name: "invalid bids",
			fields: fields{
				Timestamp: 1,
				Asks:      make([]GarntxRate, 1),
				Bids:      nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &GarantxDepth{
				Timestamp: tt.fields.Timestamp,
				Asks:      tt.fields.Asks,
				Bids:      tt.fields.Bids,
			}
			if got := d.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnmarshalGrntxDepth(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "test error",
			data:    []byte{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UnmarshalGrntxDepth(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalGrntxDepth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
