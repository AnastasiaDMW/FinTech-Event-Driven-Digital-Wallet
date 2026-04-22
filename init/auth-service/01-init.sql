CREATE TABLE users (
    id bigserial not null primary key,
    email varchar not null unique,
    encrypted_password varchar not null,
    role VARCHAR NOT NULL DEFAULT 'user',
    email_verified boolean DEFAULT FALSE,
    created_at timestamp not null DEFAULT NOW(),
    updated_at timestamp not null DEFAULT NOW()
);