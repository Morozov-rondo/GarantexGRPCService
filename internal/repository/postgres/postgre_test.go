package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"garantexGRPC/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

var (
	ctx = context.Background()
)

func TestStorage_SaveRates(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	storage := NewPostgres(sqlxDB)

	rates := models.Rates{
		Timestamp: time.Time{},
		Market:    "usdt",
		AskPrice:  decimal.NewFromFloat(100.200),
		BidPrice:  decimal.NewFromFloat(300.400),
	}

	mock.ExpectExec(`^INSERT INTO rates`).WithArgs(rates.Timestamp, rates.Market, rates.AskPrice, rates.BidPrice).WillReturnResult(sqlmock.NewResult(0, 1))

	err = storage.SaveRates(ctx, rates)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestStorage_SaveRates_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	storage := NewPostgres(sqlxDB)

	rates := models.Rates{
		Timestamp: time.Time{},
		Market:    "usdt",
		AskPrice:  decimal.NewFromFloat(100.200),
		BidPrice:  decimal.NewFromFloat(300.400),
	}

	mock.ExpectExec(`^INSERT INTO rates`).WithArgs(rates.Timestamp, rates.Market, rates.AskPrice, rates.BidPrice).WillReturnError(errors.New("insert error"))

	err = storage.SaveRates(ctx, rates)

	if err == nil {
		t.Error("expected error, got none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
