package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"garantexGRPC/configs"
	"garantexGRPC/internal/garantex_api"
	"garantexGRPC/internal/grpc_server"
	"garantexGRPC/internal/repository/postgres"
	"garantexGRPC/internal/service"
	"garantexGRPC/pkg/logger"
	"garantexGRPC/pkg/tracer"
	garantex_sso_v1_ssov1 "garantexGRPC/protos/gen/go"
	"github.com/golang-migrate/migrate/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var (
	log *zap.SugaredLogger
)

func main() {

	conf, err := configs.New()
	if err != nil {
		panic(err)
	}

	db, err := postgres.NewPostgresDB(conf.Database)
	if err != nil {
		fmt.Println("db err", err)
		panic(err)
	}

	m, err := migrate.New("file://./migrations", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		conf.Database.Username, conf.Database.Password, conf.Database.Host, conf.Database.Port, conf.Database.DBName))
	if err != nil {
		log.Fatal("migrate error", zap.Error(err))
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("migrate error", zap.Error(err))
	}

	trace, err := tracer.NewProvider(conf.Trace)
	if err != nil {
		log.Fatal("tracer error", zap.Error(err))
	}
	defer func() {
		_ = trace.Provider.Shutdown(context.Background())
	}()

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Info("metrics listen and serve", zap.String("port", conf.HTTPServer.Port))
		if err := http.ListenAndServe(":"+conf.HTTPServer.Port, nil); err != nil {
			log.Fatal("metrics listen and serve error", zap.Error(err))
		}
	}()

	logger.BuildLogger(conf.Logger.LogLevel)
	log = logger.Logger().Named("main").Sugar()

	garApi := garantex_api.NewGarantex(conf.Garantex)

	strorage := postgres.NewPostgres(db)

	exc := service.NewExchange(garApi, strorage)

	grpcExchange := grpc_server.NewExchangeGRPC(exc)

	srv := grpc.NewServer()
	garantex_sso_v1_ssov1.RegisterRatesServer(srv, grpcExchange)

	healthcheck := health.NewServer()
	healthgrpc.RegisterHealthServer(srv, healthcheck)

	reflection.Register(srv)

	l, err := net.Listen("tcp", configs.CreateAddr(conf.GRPCServer.Host, conf.GRPCServer.Port))
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}
	defer l.Close()

	go func() {
		log.Info("GRPC server start", zap.String("address", configs.CreateAddr(conf.GRPCServer.Host, conf.GRPCServer.Port)))
		if err := srv.Serve(l); err != nil {
			log.Fatal("failed to serve", zap.Error(err))
		}

	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Info("Received shutdown signal", zap.String("signal", sig.String()))

	srv.GracefulStop()

	log.Info("Server stopped gracefully")
}
