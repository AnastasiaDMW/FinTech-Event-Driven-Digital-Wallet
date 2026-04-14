CREATE TABLE users_profile (
    user_id BIGINT PRIMARY KEY,
    phone TEXT UNIQUE,
    birth_date DATE,
    first_name TEXT,
    last_name TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    account_number TEXT NOT NULL UNIQUE,
    currency TEXT NOT NULL DEFAULT 'RUB',
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_accounts_user
        FOREIGN KEY (user_id)
        REFERENCES users_profile(user_id)
        ON DELETE CASCADE
);

CREATE INDEX idx_accounts_user_id ON accounts(user_id);