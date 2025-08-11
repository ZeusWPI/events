-- +goose Up
-- +goose StatementBegin
ALTER TABLE announcement
ADD COLUMN author_id INTEGER REFERENCES board (id);

UPDATE announcement a
SET author_id = (
  SELECT b.id
  FROM board b
  WHERE b.year_id = a.year_id
  LIMIT 1
);

ALTER TABLE announcement
ALTER COLUMN author_id SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE announcement
DROP COLUMN author_id;
-- +goose StatementEnd
