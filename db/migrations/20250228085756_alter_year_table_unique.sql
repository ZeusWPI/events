-- +goose Up
-- +goose StatementBegin
ALTER TABLE year
ADD UNIQUE (start_year, end_year);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE year 
DROP CONSTRAINT year_start_year_end_year_key;
-- +goose StatementEnd
