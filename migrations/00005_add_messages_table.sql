-- +goose Up
-- +goose StatementBegin
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    text TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_messages_event_id ON messages(event_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_messages_event_id;
DROP INDEX IF EXISTS idx_messages_created_at
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd
