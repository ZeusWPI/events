-- +goose Up
-- +goose StatementBegin
ALTER TABLE member 
ADD COLUMN zauth_id INTEGER;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE member 
DROP COLUMN zauth_id;
-- +goose StatementEnd
