-- +goose Up
-- +goose StatementBegin
CREATE TABLE year (
  id SERIAL PRIMARY KEY,
  start_year INTEGER NOT NULL,
  end_year INTEGER NOT NULL
);

ALTER TABLE event
DROP COLUMN year;

ALTER TABLE event
ADD COLUMN year INTEGER NOT NULL REFERENCES year (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE event
DROP COLUMN year;

ALTER TABLE event
ADD COLUMN year VARCHAR(255);

DROP TABLE year;
-- +goose StatementEnd
