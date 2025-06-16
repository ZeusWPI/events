-- +goose Up
-- +goose StatementBegin
ALTER TABLE event 
ADD UNIQUE (file_name, year_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE event 
DROP CONSTRAINT year_file_name_year_id;
-- +goose StatementEnd
