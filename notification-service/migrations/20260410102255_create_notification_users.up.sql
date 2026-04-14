CREATE TABLE notification_users (
    user_id BIGINT PRIMARY KEY,
    email TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);