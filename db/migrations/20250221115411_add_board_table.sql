-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS board (
  id SERIAL PRIMARY KEY,
  member INTEGER NOT NULL REFERENCES member (id),
  academic_year INTEGER NOT NULL REFERENCES academic_year (id),
  role VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS board;
-- +goose StatementEnd
