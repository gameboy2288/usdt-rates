package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"usdt-rates/internal/domain"
)

func FetchRates() (*domain.Rates, error) {
	resp, err := http.Get("https://garantex.org/api/v2/depth?market=usdtrub")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var rates domain.Rates
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		fmt.Println("decode err")
		return nil, err
	}

	// askPrice := rates.Asks[0]
	// bidPrice := rates.Bids[0]

	// ask, _ := strconv.ParseFloat(askPrice, 64)
	// bid, _ := strconv.ParseFloat(bidPrice, 64)

	return &rates, nil
}
