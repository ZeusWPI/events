-- +goose Up
-- +goose StatementBegin
ALTER TABLE task
ADD COLUMN duration INTERVAL NOT NULL DEFAULT '0 seconds';

ALTER TABLE task
ALTER COLUMN duration DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE task
DROP COLUMN duration;
-- +goose StatementEnd
