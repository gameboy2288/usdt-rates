package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	myGrpc "usdt-rates/internal/transport/grpc"
	pb "usdt-rates/proto"

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

	server := grpc.NewServer()

	// Регистрация обработчиков gRPC
	pb.RegisterRateServiceServer(server, &myGrpc.RateHandler{})
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
