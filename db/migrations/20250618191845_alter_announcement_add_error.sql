-- +goose Up
-- +goose StatementBegin
ALTER TABLE announcement
ADD COLUMN error TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE announcement
DROP COLUMN error TEXT;
-- +goose StatementEnd
