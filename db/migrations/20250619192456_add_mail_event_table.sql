-- +goose Up
-- +goose StatementBegin
CREATE TABLE mail_event (
  id SERIAL PRIMARY KEY,
  mail_id INTEGER NOT NULL REFERENCES mail (id) ON DELETE CASCADE,
  event_id INTEGER NOT NULL REFERENCES event (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE mail_event;
-- +goose StatementEnd
