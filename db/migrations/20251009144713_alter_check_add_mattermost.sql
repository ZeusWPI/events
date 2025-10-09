-- +goose Up
-- +goose StatementBegin
ALTER TABLE check_event
ADD COLUMN mattermost BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE check_event
DROP COLUMN mattermost;
-- +goose StatementEnd
