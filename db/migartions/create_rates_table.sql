-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS rates (
    id SERIAL PRIMARY KEY,
    timestamp BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS asks (
    id SERIAL PRIMARY KEY,
    rate_id INT REFERENCES rates(id) ON DELETE CASCADE,
    price FLOAT NOT NULL,
    volume FLOAT NOT NULL,
    amount FLOAT NOT NULL,
    factor FLOAT,
    type VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS bids (
    id SERIAL PRIMARY KEY,
    rate_id INT REFERENCES rates(id) ON DELETE CASCADE,
    price FLOAT NOT NULL,
    volume FLOAT NOT NULL,
    amount FLOAT NOT NULL,
    factor FLOAT,
    type VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

-- Индексы для ускорения поиска по таблицам
CREATE INDEX IF NOT EXISTS idx_rates_timestamp ON rates (timestamp);
CREATE INDEX IF NOT EXISTS idx_asks_rate_id ON asks (rate_id);
CREATE INDEX IF NOT EXISTS idx_bids_rate_id ON bids (rate_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_rates_timestamp;
DROP INDEX IF EXISTS idx_asks_rate_id;
DROP INDEX IF EXISTS idx_bids_rate_id;

DROP TABLE IF EXISTS bids;
DROP TABLE IF EXISTS asks;
DROP TABLE IF EXISTS rates;

-- +goose StatementEnd
