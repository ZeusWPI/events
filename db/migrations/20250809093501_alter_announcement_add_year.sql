-- +goose Up
-- +goose StatementBegin
ALTER TABLE announcement
ADD COLUMN year_id INTEGER REFERENCES year (id) ON DELETE CASCADE;

UPDATE announcement a
SET year_id = (
  SELECT year_id
  FROM event e
  WHERE e.id = a.event_id
);

ALTER TABLE announcement
ALTER COLUMN year_id SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE announcement
DROP COLUMN year_id;
-- +goose StatementEnd
