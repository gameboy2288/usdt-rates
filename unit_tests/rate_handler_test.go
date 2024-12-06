package unit_tests

import (
	"context"
	"errors"
	"testing"
	"usdt-rates/internal/repository"
	"usdt-rates/internal/repository/mocks" // Моки, сгенерированные через MockGen
	pb "usdt-rates/proto"

	myGrpc "usdt-rates/internal/transport/grpc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateHandler_GetRates_Success(t *testing.T) {
	// Подготовка
	ctx := context.Background()

	mockRepo := new(mocks.Repository)
	mockRepo.On("FetchRates", mock.Anything).Return([]repository.Rate{
		{Timestamp: 1234567890, Price: "100.0", Volume: "10.0"},
	}, nil)

	handler := myGrpc.NewRateHandler(mockRepo)

	// Тестируемый метод
	req := &pb.Empty{}
	resp, err := handler.GetRates(ctx, req)

	// Проверка
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Asks, 1)
	assert.Equal(t, resp.Timestamp, int64(1234567890))
	assert.Equal(t, resp.Asks[0].Price, "100.0")
	assert.Equal(t, resp.Asks[0].Volume, "10.0")

	// Убедиться, что все ожидаемые вызовы выполнены
	mockRepo.AssertExpectations(t)
}

func TestRateHandler_GetRates_Error(t *testing.T) {
	// Подготовка
	ctx := context.Background()

	mockRepo := new(mocks.Repository)
	mockRepo.On("FetchRates", mock.Anything).Return(nil, errors.New("database error"))

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
