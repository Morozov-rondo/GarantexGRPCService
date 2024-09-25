package service

import (
	"context"

	"garantexGRPC/internal/garantex_api"
	"garantexGRPC/internal/repository"
	"garantexGRPC/models"
	"garantexGRPC/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var (
	log = logger.Logger().Named("exchange_service").Sugar()
)
var (
	tracer = otel.Tracer("exchange_service")
)

//go:generate mockgen -source=exchange.go -destination=./mocks/mock_exchange_interface.go -package=mocks
type Exchanger interface {
	GetRates(ctx context.Context, market string) (models.Rates, error)
}

type Exchange struct {
	garantex garantex_api.Garantexer
	db       repository.Storager
}

func NewExchange(garantex garantex_api.Garantexer, db repository.Storager) *Exchange {
	return &Exchange{
		garantex: garantex,
		db:       db,
	}
}

func (e *Exchange) GetRates(ctx context.Context, market string) (models.Rates, error) {

	ctx, span := tracer.Start(ctx, "GetRates")
	defer span.End()

	rates, err := e.garantex.GetRates(ctx, market)
	if err != nil {
		log.Error("garantex get rates error", zap.Error(err))
		return models.Rates{}, err
	}

	err = e.db.SaveRates(ctx, rates)
	if err != nil {
		log.Error("postgres save rates error", zap.Error(err))
		return models.Rates{}, err
	}

	return rates, nil
}
