-- +goose Up
-- +goose StatementBegin
CREATE TABLE announcement (
  id SERIAL PRIMARY KEY,
  event_id INTEGER NOT NULL REFERENCES event (id) ON DELETE CASCADE,
  content TEXT NOT NULL,
  send_time TIMESTAMPTZ NOT NULL,
  send BOOLEAN NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE announcement;
-- +goose StatementEnd
