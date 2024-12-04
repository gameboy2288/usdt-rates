package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"usdt-rates/internal/domain"
)

func FetchRates() (*domain.Rate, error) {
	resp, err := http.Get("https://garantex.org/api/v2/depth?market=usdtrub")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Asks [][]string `json:"asks"`
		Bids [][]string `json:"bids"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	askPrice := result.Asks[0][0]
	bidPrice := result.Bids[0][0]

	ask, _ := strconv.ParseFloat(askPrice, 64)
	bid, _ := strconv.ParseFloat(bidPrice, 64)

	return &domain.Rate{
		Ask:       float32(ask),
		Bid:       float32(bid),
		Timestamp: time.Now().Unix(),
	}, nil
}
