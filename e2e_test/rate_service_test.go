package e2e_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	pb "usdt-rates/proto"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestRateService_GetRate(t *testing.T) {
	// Устанавливаем соединение с gRPC-сервером
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewRateServiceClient(conn)

	// Выполняем запрос
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := client.GetRates(ctx, &pb.Empty{})
	fmt.Println(resp)
	assert.NoError(t, err)
	assert.Greater(t, resp.Bids[3].Price, float32(0.0), "Rate should be greater than 0")
}
