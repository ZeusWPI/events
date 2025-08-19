-- +goose Up
-- +goose StatementBegin
ALTER TABLE poster
ADD COLUMN webp_id VARCHAR(255) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE poster
DROP COLUMN webp_id;
-- +goose StatementEnd
