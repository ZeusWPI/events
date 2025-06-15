-- +goose Up
-- +goose StatementBegin
CREATE TABLE organizer (
  id SERIAL PRIMARY KEY,
  event INTEGER NOT NULL REFERENCES event (id),
  board INTEGER NOT NULL REFERENCES board (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organizer;
-- +goose StatementEnd
