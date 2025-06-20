-- +goose Up
-- +goose StatementBegin
CREATE TABLE mail (
  id SERIAL PRIMARY KEY,
  content TEXT NOT NULL,
  send_time TIMESTAMPTZ NOT NULL,
  send BOOLEAN NOT NULL,
  error TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE mail;
-- +goose StatementEnd
