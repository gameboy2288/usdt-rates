package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"usdt-rates/internal/transport/grpc"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
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

	// Запуск gRPC сервера
	go func() {
		if err := grpc.StartServer(); err != nil {
			log.Fatalf("failed to start gRPC server: %v", err)
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
