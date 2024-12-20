package grpc

import (
	"context"
	"strconv"
	"usdt-rates/internal/domain"
	"usdt-rates/internal/metrics"
	"usdt-rates/internal/repository"
	"usdt-rates/internal/service"
	pb "usdt-rates/proto"
)

type RateHandler struct {
	repo repository.RepositoryInterface

	pb.UnimplementedRateServiceServer
}

func NewRateHandler(repo repository.RepositoryInterface) *RateHandler {
	return &RateHandler{
		repo: repo,
	}
}

func (h *RateHandler) GetRates(ctx context.Context, req *pb.Empty) (*pb.RateResponse, error) {
	rate, err := service.FetchRates()
	if err != nil {
		metrics.RecordRequest("500")
		return nil, err
	}

	rateID, err := h.repo.SaveRate(rate.Timestamp)
	if err != nil {
		metrics.RecordRequest("500")
		return nil, err
	}

	err = h.repo.SaveAsk(rateID, rate.Asks[0])
	if err != nil {
		metrics.RecordRequest("500")
		return nil, err
	}

	err = h.repo.SaveBid(rateID, rate.Bids[0])
	if err != nil {
		metrics.RecordRequest("500")
		return nil, err
	}

	convertedAsks, err := convertAskToPb(rate.Asks)
	if err != nil {
		metrics.RecordRequest("500")
		return nil, err
	}
	convertedBids, err := convertAskToPb(rate.Bids)
	if err != nil {
		metrics.RecordRequest("500")
		return nil, err
	}
	metrics.RecordRequest("200")
	return &pb.RateResponse{
		Timestamp: rate.Timestamp,
		Asks:      convertedAsks,
		Bids:      convertedBids,
	}, nil
}

func convertAskToPb(asks []domain.Ask) ([]*pb.Ask, error) {
	converted := make([]*pb.Ask, len(asks))
	for i, ask := range asks {
		price, err := strconv.ParseFloat(ask.Price, 32)
		if err != nil {
			return nil, err
		}
		volume, err := strconv.ParseFloat(ask.Volume, 32)
		if err != nil {
			return nil, err
		}
		amount, err := strconv.ParseFloat(ask.Amount, 32)
		if err != nil {
			return nil, err
		}
		factor, err := strconv.ParseFloat(ask.Factor, 32)
		if err != nil {
			return nil, err
		}

		converted[i] = &pb.Ask{
			Price:  float32(price),
			Volume: float32(volume),
			Amount: float32(amount),
			Factor: float32(factor),
			Type:   string(ask.Type),
		}
	}
	return converted, nil
}

func (h *RateHandler) HealthCheck(ctx context.Context, req *pb.Empty) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Status: true}, nil
}
