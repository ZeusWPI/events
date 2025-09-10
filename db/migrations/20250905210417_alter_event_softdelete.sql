-- +goose Up
-- +goose StatementBegin
ALTER TABLE event
ADD COLUMN deleted BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE event
ALTER COLUMN deleted DROP DEFAULT;

ALTER TABLE event 
DROP CONSTRAINT event_file_name_year_id_key;

CREATE UNIQUE INDEX unique_file_name_year_id_not_deleted
ON event (file_name, year_id)
WHERE deleted = false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX unique_file_name_year_id_not_deleted;

ALTER TABLE event
DROP COLUMN deleted;

ALTER TABLE event 
ADD UNIQUE (file_name, year_id);
-- +goose StatementEnd
