package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"usdt-rates/internal/config"
	"usdt-rates/internal/logger"
	"usdt-rates/internal/metrics"
	"usdt-rates/internal/repository"
	myGrpc "usdt-rates/internal/transport/grpc"
	pb "usdt-rates/proto"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

func main() {
	logger.InitLogger()

	metrics.InitMetrics()

	go metrics.StartMetricsServer()

	ctx := context.Background()

	// Инициализация OpenTelemetry
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		logger.Log.Fatal("failed to initialize exporter: ", zap.Error(err))
	}

	tp := trace.NewTracerProvider(trace.WithBatcher(exp))
	otel.SetTracerProvider(tp)

	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Log.Fatal("failed to shutdown TracerProvider: ", zap.Error(err))
		}
	}()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatal("failed to load config: ", zap.Error(err))
	}

	db, err := sql.Open("postgres", fmt.Sprintf("host=postgres port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbPort,
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbName))
	if err != nil {
		logger.Log.Fatal("failed to connect to database: ", zap.Error(err))
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		logger.Log.Fatal("failed to ping database: ", zap.Error(err))
	}
	err = goose.Up(db, "db/migrations")
	if err != nil {
		logger.Log.Fatal("failed to run migrations: ", zap.Error(err))
	}

	repo := repository.NewRepository(db)
	rateHandler := myGrpc.NewRateHandler(repo)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	// Регистрация обработчиков gRPC
	pb.RegisterRateServiceServer(server, rateHandler)
	pb.RegisterHealthServer(server, &myGrpc.HealthHandler{})

	// Запуск gRPC сервера
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log.Fatal("failed to listen: ", zap.Error(err))
	}
	go func() {
		logger.Log.Info("Starting gRPC server on port 50051...")
		if err := server.Serve(listener); err != nil {
			logger.Log.Fatal("failed to serve: ", zap.Error(err))
		}
	}()

	// Обработка Graceful Shutdown
	gracefulShutdown(ctx, tp)
}

func gracefulShutdown(ctx context.Context, tp *trace.TracerProvider) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Остановка провайдера трассировки
	if err := tp.Shutdown(ctx); err != nil {
		logger.Log.Info("Error shutting down tracer provider: ", zap.Error(err))
	}
}
