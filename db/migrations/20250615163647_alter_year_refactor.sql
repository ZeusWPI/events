-- +goose Up
-- +goose StatementBegin
ALTER TABLE year RENAME COLUMN start_year TO year_start;
ALTER TABLE year RENAME COLUMN end_year TO year_end;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE year RENAME COLUMN year_end TO end_year;
ALTER TABLE year RENAME COLUMN year_start TO start_year;
-- +goose StatementEnd
