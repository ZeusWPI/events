-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event (
  id SERIAL PRIMARY KEY,
  url VARCHAR(255) NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  start_time TIMESTAMPTZ NOT NULL,
  end_time TIMESTAMPTZ NOT NULL,
  academic_year VARCHAR(255) NOT NULL,
  location TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event;
-- +goose StatementEnd
