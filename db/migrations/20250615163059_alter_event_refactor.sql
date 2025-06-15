-- +goose Up
-- +goose StatementBegin
ALTER TABLE event RENAME COLUMN url TO file_name;
ALTER TABLE event RENAME COLUMN year TO year_id;

ALTER TABLE event DROP CONSTRAINT event_year_fkey;
ALTER TABLE event
ADD CONSTRAINT fk_event_year
FOREIGN KEY (year_id)
REFERENCES year (id)
ON DELETE CASCADE;

ALTER TABLE event DROP COLUMN created_at;
ALTER TABLE event DROP COLUMN updated_at;
ALTER TABLE event DROP COLUMN deleted_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE event ADD COLUMN deleted_at TIMESTAMPTZ DEFAULT NULL;
ALTER TABLE event ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE event ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE event DROP CONSTRAINT event_year_fkey;
ALTER TABLE event
ADD CONSTRAINT fk_event_year
FOREIGN KEY (year_id)
REFERENCES year (id);

ALTER TABLE event RENAME COLUMN year_id TO year;
ALTER TABLE event RENAME COLUMN file_name TO url;
-- +goose StatementEnd
