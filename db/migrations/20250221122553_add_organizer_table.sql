-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS organizer (
  id SERIAL PRIMARY KEY,
  event INTEGER NOT NULL REFERENCES event (id),
  board INTEGER NOT NULL REFERENCES board (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS organizer;
-- +goose StatementEnd
