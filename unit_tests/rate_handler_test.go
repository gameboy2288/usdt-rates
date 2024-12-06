package unit_tests

import (
	"context"
	"errors"
	"testing"
	"usdt-rates/internal/repository/mocks" // Моки, сгенерированные через MockGen
	pb "usdt-rates/proto"

	myGrpc "usdt-rates/internal/transport/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateHandler_GetRates_Success(t *testing.T) {
	// Подготовка
	ctx := context.Background()

	mockRepo := new(mocks.RepositoryInterface)
	// rates := &domain.Rates{
	// 	Timestamp: time.Now().Unix(),
	// 	Asks:      []domain.Ask{{Price: "100.08", Volume: "55000.0", Amount: "5504400.0", Factor: "0.0", Type: "limit"}},
	// 	Bids:      []domain.Ask{{Price: "95.0", Volume: "5.0", Amount: "", Factor: "", Type: ""}},
	// }
	// mockRepo.On("FetchRates").Return(rates, nil)
	mockRepo.On("SaveRate", mock.AnythingOfType("int64")).Return(int64(1), nil)
	mockRepo.On("SaveAsk", int64(1), mock.AnythingOfType("domain.Ask")).Return(nil)
	mockRepo.On("SaveBid", int64(1), mock.AnythingOfType("domain.Ask")).Return(nil)

	handler := myGrpc.NewRateHandler(mockRepo)

	// Тестируемый метод
	req := &pb.Empty{}
	resp, err := handler.GetRates(ctx, req)

	// Проверка
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Убедиться, что все ожидаемые вызовы выполнены
	mockRepo.AssertExpectations(t)
}

func TestRateHandler_GetRates_Error(t *testing.T) {
	// Подготовка
	ctx := context.Background()

	mockRepo := new(mocks.RepositoryInterface)
	mockRepo.On("SaveRate", mock.AnythingOfType("int64")).Return(int64(0), errors.New("database error"))

	handler := myGrpc.NewRateHandler(mockRepo)

	// Тестируемый метод
	req := &pb.Empty{}
	resp, err := handler.GetRates(ctx, req)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, resp)

	// Убедиться, что все ожидаемые вызовы выполнены
	mockRepo.AssertExpectations(t)
}
