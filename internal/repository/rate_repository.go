package repository

import (
	"database/sql"
	"usdt-rates/internal/domain"

	_ "github.com/lib/pq"
)

type RateRepository struct {
	db *sql.DB
}

func NewRateRepository(db *sql.DB) *RateRepository {
	return &RateRepository{db: db}
}

func (r *RateRepository) SaveRate(rate *domain.Rate) error {
	_, err := r.db.Exec("INSERT INTO rates (ask, bid, timestamp) VALUES ($1, $2, $3)", rate.Ask, rate.Bid, rate.Timestamp)
	return err
}
