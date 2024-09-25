package postgres

import (
	"context"
	"time"

	"garantexGRPC/models"
	"garantexGRPC/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var (
	repoDuretion = prometheus.NewHistogram(prometheus.HistogramOpts{
		Subsystem: "postgreSQL",
		Name:      "save_rates_to_db_duration_seconds",
		Help:      "Database query duration in seconds",
		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 10),
	})
)

var (
	log = logger.Logger().Named("postgres").Sugar()
)
var (
	tracer = otel.Tracer("postgres")
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(db *sqlx.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) SaveRates(ctx context.Context, rates models.Rates) error {

	ctx, span := tracer.Start(ctx, "SaveRates")
	defer span.End()

	start := time.Now()
	defer func() {
		repoDuretion.Observe(time.Since(start).Seconds())
	}()

	q := `INSERT INTO rates (timestamp, market, ask, bid) VALUES ($1, $2, $3, $4)`

	_, err := p.db.ExecContext(ctx, q, rates.Timestamp, rates.Market, rates.AskPrice, rates.BidPrice)
	if err != nil {
		log.Error("postgres insert rates error", zap.Error(err))
		return err
	}

	return nil
}
