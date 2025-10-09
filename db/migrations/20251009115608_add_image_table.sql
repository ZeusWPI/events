-- +goose Up
-- +goose StatementBegin
CREATE TABLE image (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  file_id VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE image;
-- +goose StatementEnd
