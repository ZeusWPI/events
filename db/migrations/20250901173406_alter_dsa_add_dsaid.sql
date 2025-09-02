-- +goose Up
-- +goose StatementBegin
ALTER TABLE dsa
ADD COLUMN dsa_id integer;
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
ADD COLUMN entry boolean NOT NULL DEFAULT true;
-- +goose StatementEnd
