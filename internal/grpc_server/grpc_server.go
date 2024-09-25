package grpc_server

import (
	"context"
	"time"

	"garantexGRPC/internal/service"
	"garantexGRPC/pkg/logger"
	garantex_sso_v1_ssov1 "garantexGRPC/protos/gen/go"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var (
	log = logger.Logger().Named("grpc_server").Sugar()
)

var (
	getRequest = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "get_request",
		Help:      "Get request counter",
		Subsystem: "grpc_server",
	})

	getDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:      "get_duration",
		Help:      "Get duration histogram",
		Subsystem: "grpc_server",
		Buckets:   prometheus.DefBuckets,
	})
)

var (
	tracer = otel.Tracer("grpc_server")
)

type ExchangeGRPC struct {
	Exchange service.Exchanger
	garantex_sso_v1_ssov1.UnimplementedRatesServer
}

func NewExchangeGRPC(exchange service.Exchanger) *ExchangeGRPC {

	return &ExchangeGRPC{
		Exchange:                 exchange,
		UnimplementedRatesServer: garantex_sso_v1_ssov1.UnimplementedRatesServer{},
	}
}

func (e *ExchangeGRPC) Get(ctx context.Context, request *garantex_sso_v1_ssov1.GetRequest) (*garantex_sso_v1_ssov1.GetResponse, error) {

	ctx, span := tracer.Start(ctx, "Get")
	defer span.End()

	getRequest.Inc()
	start := time.Now()
	defer func() {
		getDuration.Observe(time.Since(start).Seconds())
	}()

	rates, err := e.Exchange.GetRates(ctx, request.GetMarket().String())
	if err != nil {
		log.Error("exchange get rates error", zap.Error(err))
		return nil, err
	}
	return &garantex_sso_v1_ssov1.GetResponse{
		Timestamp: rates.Timestamp.Unix(),
		Market:    request.GetMarket(),
		Ask:       rates.AskPrice.InexactFloat64(),
		Bid:       rates.BidPrice.InexactFloat64(),
	}, nil
}
