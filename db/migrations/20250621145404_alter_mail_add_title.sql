-- +goose Up
-- +goose StatementBegin
ALTER TABLE mail
ADD COLUMN title VARCHAR(255) NOT NULL DEFAULT '[Zeus WPI] Mail';

ALTER TABLE mail
ALTER COLUMN title DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE mail
DROP COLUMN title;
-- +goose StatementEnd
