package grpc

import (
	"context"

	pb "usdt-rates/proto"
)

// HealthHandler - обработчик для HealthCheck
type HealthHandler struct {
	pb.UnimplementedHealthServer
}

// HealthCheck проверяет состояние сервиса
func (h *HealthHandler) HealthCheck(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Status: "SERVING",
	}, nil
}
