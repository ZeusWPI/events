-- +goose Up
-- +goose StatementBegin
CREATE TABLE member (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  username VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE member;
-- +goose StatementEnd
