-- +goose Up
-- +goose StatementBegin
ALTER TABLE dsa
ADD COLUMN dsa_id INTEGER;

ALTER TABLE dsa
ADD UNIQUE (dsa_id);

ALTER TABLE dsa
DROP COLUMN entry;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE dsa
DROP COLUMN dsa_id;

ALTER TABLE dsa
ADD COLUMN entry BOOLEAN NOT NULL DEFAULT true;
-- +goose StatementEnd
