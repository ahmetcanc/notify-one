-- +goose Up
-- +goose StatementBegin
CREATE TYPE notification_channel AS ENUM ('sms', 'email', 'push');
CREATE TYPE notification_status AS ENUM ('pending', 'processing', 'sent', 'failed', 'cancelled');
CREATE TYPE notification_priority AS ENUM ('low', 'normal', 'high');

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY,
    batch_id UUID,
    recipient TEXT NOT NULL,
    channel notification_channel NOT NULL,
    content TEXT NOT NULL,
    priority notification_priority DEFAULT 'normal',
    status notification_status DEFAULT 'pending',
    idempotency_key TEXT UNIQUE,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_batch_id ON notifications(batch_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notifications;
DROP TYPE IF EXISTS notification_status;
DROP TYPE IF EXISTS notification_channel;
DROP TYPE IF EXISTS notification_priority;
-- +goose StatementEnd