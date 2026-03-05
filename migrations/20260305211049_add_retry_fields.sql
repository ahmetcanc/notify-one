-- +goose Up
-- +goose StatementBegin
ALTER TABLE notifications ADD COLUMN retry_count INT DEFAULT 0;
ALTER TABLE notifications ADD COLUMN next_retry_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE notifications DROP COLUMN retry_count;
ALTER TABLE notifications DROP COLUMN next_retry_at;
-- +goose StatementEnd