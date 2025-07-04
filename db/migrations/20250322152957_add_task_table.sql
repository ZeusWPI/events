-- +goose Up
-- +goose StatementBegin
CREATE TABLE task (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  result VARCHAR(255) CHECK (result IN ('success', 'failed')),
  run_at TIMESTAMPTZ NOT NULL,
  error TEXT,
  recurring BOOLEAN NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE task;
-- +goose StatementEnd
