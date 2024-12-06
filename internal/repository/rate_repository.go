package repository

import (
	"database/sql"

	"usdt-rates/internal/domain"
)

type RepositoryInterface interface {
	SaveRate(timestamp int64) (int64, error)
	SaveAsk(rateID int64, ask domain.Ask) error
	SaveBid(rateID int64, bid domain.Ask) error
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveRate(timestamp int64) (int64, error) {
	query := `INSERT INTO rates (timestamp) VALUES ($1) RETURNING id`
	var rateID int64
	err := r.db.QueryRow(query, timestamp).Scan(&rateID)
	if err != nil {
		return 0, err
	}
	return rateID, nil
}

func (r *Repository) SaveAsk(rateID int64, ask domain.Ask) error {
	query := `INSERT INTO asks (rate_id, price, volume, amount, factor, type) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, rateID, ask.Price, ask.Volume, ask.Amount, ask.Factor, string(ask.Type))
	return err
}

func (r *Repository) SaveBid(rateID int64, bid domain.Ask) error {
	query := `INSERT INTO bids (rate_id, price, volume, amount, factor, type) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(query, rateID, bid.Price, bid.Volume, bid.Amount, bid.Factor, string(bid.Type))
	return err
}
