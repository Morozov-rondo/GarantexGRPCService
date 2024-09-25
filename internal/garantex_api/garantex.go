package garantex_api

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"garantexGRPC/configs"
	"garantexGRPC/models"
	"garantexGRPC/pkg/logger"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var (
	log = logger.Logger().Named("garantex_api").Sugar()
)

var (
	getRatesDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Subsystem: "garantex_api",
		Name:      "get_rates_duration_seconds",
		Help:      "Garantex get rates duration in seconds",
		Buckets:   prometheus.DefBuckets,
	})

	getRatesCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: "garantex_api",
		Name:      "get_rates_count",
		Help:      "Garantex get rates count",
	}, []string{"market"})
)

var (
	tracer = otel.Tracer("garantex_api")
)

//go:generate mockgen -source=garantex.go -destination=./mocks/mock_garantex_interface.go -package=mocks
type Garantexer interface {
	GetRates(ctx context.Context, market string) (models.Rates, error)
}

type Garantex struct {
	client *http.Client
	url    string
}

func NewGarantex(config configs.GarantexConfig) *Garantex {
	return &Garantex{

		client: &http.Client{
			Timeout: config.Timeout,
		},
		url: config.URL,
	}
}

func (g *Garantex) GetRates(ctx context.Context, market string) (models.Rates, error) {

	ctx, span := tracer.Start(ctx, "GetRates")
	defer span.End()

	start := time.Now()
	defer func() {
		getRatesDuration.Observe(time.Since(start).Seconds())
	}()
	market = strings.ToLower(market)

	req, err := http.NewRequestWithContext(ctx, "GET", g.url+"/depth?market="+market, nil)
	if err != nil {
		log.Error("garantex new request error: %v", zap.Error(err))
		return models.Rates{}, err
	}
	resp, err := g.client.Do(req)
	if err != nil {
		log.Error("garantex do request error", zap.Error(err))
		return models.Rates{}, err
	}
	defer resp.Body.Close()

	getRatesCount.WithLabelValues(market).Inc()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("garantex read body error", zap.Error(err))
		return models.Rates{}, err
	}
	depth, err := models.UnmarshalGrntxDepth(body)
	if err != nil {
		log.Error("garantex unmarshal body error", zap.Error(err))
		return models.Rates{}, err
	}
	rates, err := depth.ToDomain()
	if err != nil {
		log.Error("garantex to domain error", zap.Error(err))
		return models.Rates{}, err
	}
	rates.Market = market
	return rates, nil
}
