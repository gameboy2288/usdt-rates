CREATE TABLE rates (
    id SERIAL PRIMARY KEY,
    ask FLOAT NOT NULL,
    bid FLOAT NOT NULL,
    timestamp BIGINT NOT NULL
);
