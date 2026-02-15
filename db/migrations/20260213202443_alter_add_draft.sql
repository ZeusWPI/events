-- +goose Up
-- +goose StatementBegin
ALTER TABLE mail
ALTER COLUMN send_time DROP NOT NULL;

ALTER TABLE mail
ADD COLUMN draft BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE mail
ALTER COLUMN draft DROP DEFAULT;


ALTER TABLE announcement
ALTER COLUMN send_time DROP NOT NULL;

ALTER TABLE announcement
ADD COLUMN draft BOOLEAN NOT NULL DEFAULT false;

ALTER TABLE announcement
ALTER COLUMN draft DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE announcement
DROP COLUMN draft;

ALTER TABLE announcement
ALTER COLUMN send_time SET NOT NULL;

ALTER TABLE mail
DROP COLUMN draft;

ALTER TABLE mail
ALTER COLUMN send_time SET NOT NULL;
-- +goose StatementEnd
