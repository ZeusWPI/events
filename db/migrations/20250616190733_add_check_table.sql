-- +goose Up
-- +goose StatementBegin
CREATE TABLE "check" (
  id SERIAL PRIMARY KEY,
  event_id INTEGER NOT NULL REFERENCES event (id) ON DELETE CASCADE,
  description TEXT NOT NULL,
  done BOOLEAN NOT NULL,

  UNIQUE (event_id, description)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "check";
-- +goose StatementEnd
