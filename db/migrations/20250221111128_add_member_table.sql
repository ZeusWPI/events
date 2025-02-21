-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS member (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  username VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS member;
-- +goose StatementEnd
