-- +goose Up
-- +goose StatementBegin
ALTER TABLE board
ADD COLUMN mattermost VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE board
DROP COLUMN mattermost;
-- +goose StatementEnd
