package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"usdt-rates/internal/config"
	"usdt-rates/internal/repository"
	myGrpc "usdt-rates/internal/transport/grpc"
	pb "usdt-rates/proto"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	// Инициализация OpenTelemetry
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	tp := trace.NewTracerProvider(trace.WithBatcher(exp))
	otel.SetTracerProvider(tp)

	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %v", err)
		}
	}()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", fmt.Sprintf("host=postgres port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DbPort,
		cfg.DbUser,
		cfg.DbPassword,
		cfg.DbName))
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	err = goose.Up(db, "db/migrations")
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
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
		log.Fatalf("failed to listen: %v", err)
	}
	go func() {
		log.Println("Starting gRPC server on port 50051...")
		if err := server.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Обработка Graceful Shutdown
	gracefulShutdown(ctx, tp)
}

func gracefulShutdown(ctx context.Context, tp *trace.TracerProvider) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Остановка провайдера трассировки
	if err := tp.Shutdown(ctx); err != nil {
		log.Printf("Error shutting down tracer provider: %v", err)
	}
}
