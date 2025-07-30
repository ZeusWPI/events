-- +goose Up
-- +goose StatementBegin
ALTER TABLE board
ADD COLUMN is_organizer boolean NOT NULL DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE
DROP COLUMN is_organizer;
-- +goose StatementEnd
