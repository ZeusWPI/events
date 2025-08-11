-- +goose Up
-- +goose StatementBegin
ALTER TABLE mail
ADD COLUMN author_id INTEGER REFERENCES board (id);

UPDATE mail m
SET author_id = (
  SELECT b.id
  FROM board b
  WHERE b.year_id = m.year_id
  LIMIT 1
);

ALTER TABLE mail
ALTER COLUMN author_id SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mail
DROP COLUMN author_id;
-- +goose StatementEnd
