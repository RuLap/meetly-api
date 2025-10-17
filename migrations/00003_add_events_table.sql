-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    location GEOGRAPHY(Point, 4326) NOT NULL,
    address VARCHAR(255),
    starts_at TIMESTAMP WITH TIME ZONE,
    ends_at TIMESTAMP WITH TIME ZONE,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH ZONE DEFAULT NOW()
);

CREATE INDEX idx_events_location ON events USING GIST (location);
CREATE INDEX idx_events_creator_id ON events(creator_id);
CREATE INDEX idx_events_category_id ON events(category_id);
CREATE INDEX idx_events_starts_at ON events(starts_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_events_location;
DROP INDEX idx_events_creator_id;
DROP INDEX idx_events_category_id;
DROP INDEX idx_events_starts_at;
DROP TABLE events;
-- +goose StatementEnd
