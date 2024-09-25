package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type Rates struct {
	ID        int
	Market    string
	Timestamp time.Time
	AskPrice  decimal.Decimal
	BidPrice  decimal.Decimal
}


type GarantxDepth struct {
	Timestamp int64       `json:"timestamp"`
	Asks      []GarntxRate `json:"asks"`
	Bids      []GarntxRate `json:"bids"`
}

type GarntxRate struct {
	Price  string `json:"price"`
	Volume string `json:"volume"`
	Amount string `json:"amount"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

func UnmarshalGrntxDepth(data []byte) (GarantxDepth, error) {
	var r GarantxDepth
	err := json.Unmarshal(data, &r)
	return r, err
}

func (d *GarantxDepth) Valid() bool {
	switch {
	case d.Timestamp <= 0:
		return false
	case len(d.Asks) < 1:
		return false
	case len(d.Bids) < 1:
		return false
	default:
		return true
	}
}

func (d *GarantxDepth) ToDomain() (Rates, error) {
	if !d.Valid() {
		return Rates{}, fmt.Errorf("invalid data")
	}
	ask, err := decimal.NewFromString(d.Asks[0].Price)
	if err != nil {
		return Rates{}, fmt.Errorf("invalid data")
	}
	bid, err := decimal.NewFromString(d.Bids[0].Price)
	if err != nil {
		return Rates{}, fmt.Errorf("invalid data")
	}
	return Rates{
		Timestamp: time.Unix(d.Timestamp, 0),
		AskPrice:  ask,
		BidPrice:  bid,
	}, nil
}
