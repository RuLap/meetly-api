-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    birth_date DATE,
    email VARCHAR(100) UNIQUE NOT NULL,
    provider VARCHAR(50) DEFAULT 'local',
    provider_id VARCHAR(255),
    password TEXT,
    avatar_url TEXT,
    email_confirmed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_users_email_provider ON users(email, provider);
CREATE INDEX idx_users_created_at ON users(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_email_provider;
DROP INDEX IF EXISTS idx_users_created_at;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
