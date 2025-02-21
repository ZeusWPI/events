-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS academic_year (
  id SERIAL PRIMARY KEY,
  start_year INTEGER NOT NULL,
  end_year INTEGER NOT NULL
);

ALTER TABLE event
DROP COLUMN academic_year;

ALTER TABLE event
ADD COLUMN academic_year INTEGER NOT NULL REFERENCES academic_year (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE event
DROP COLUMN academic_year;

ALTER TABLE event
ADD COLUMN academic_year VARCHAR(255);

DROP TABLE IF EXISTS academic_year;
-- +goose StatementEnd
