CREATE TABLE ledgers (
    id BIGSERIAL PRIMARY KEY,
    account_number TEXT NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    idempotent TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);