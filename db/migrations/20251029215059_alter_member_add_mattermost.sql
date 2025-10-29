-- +goose Up
-- +goose StatementBegin
ALTER TABLE member
ADD COLUMN mattermost VARCHAR(255);

ALTER TABLE board
DROP COLUMN mattermost;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE board
ADD COLUMN mattermost VARCHAR(255);

ALTER TABLE member
DROP COLUMN mattermost;
-- +goose StatementEnd
