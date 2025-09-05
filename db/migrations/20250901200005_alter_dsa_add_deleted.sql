-- +goose Up
-- +goose StatementBegin
ALTER TABLE dsa
ADD COLUMN deleted boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE dsa
DROP COLUMN deleted;
-- +goose StatementEnd
