-- +goose Up
-- +goose StatementBegin
CREATE TABLE poster (
  id SERIAL PRIMARY KEY,
  event_id INTEGER NOT NULL REFERENCES event (id) ON DELETE CASCADE,
  file_id VARCHAR(255) NOT NULL,
  scc BOOLEAN NOT NULL,

  UNIQUE (event_id, scc)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE poster;
-- +goose StatementEnd
