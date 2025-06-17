-- +goose Up
-- +goose StatementBegin
CREATE TABLE dsa (
  id SERIAL PRIMARY KEY,
  event_id INTEGER NOT NULL REFERENCES event (id) ON DELETE CASCADE,
  entry BOOLEAN NOT NULL,

  UNIQUE(event_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dsa;
-- +goose StatementEnd
