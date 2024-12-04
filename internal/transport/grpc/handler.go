package grpc

import (
	"context"
	"usdt-rates/internal/service"
	pb "usdt-rates/proto"
)

type RateHandler struct {
	pb.UnimplementedRateServiceServer
}

func (h *RateHandler) GetRates(ctx context.Context, req *pb.Empty) (*pb.RateResponse, error) {
	rate, err := service.FetchRates()
	if err != nil {
		return nil, err
	}

	return &pb.RateResponse{
		Ask:       rate.Ask,
		Bid:       rate.Bid,
		Timestamp: rate.Timestamp,
	}, nil
}

func (h *RateHandler) HealthCheck(ctx context.Context, req *pb.Empty) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Status: true}, nil
}
