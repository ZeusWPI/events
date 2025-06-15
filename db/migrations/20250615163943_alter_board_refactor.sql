-- +goose Up
-- +goose StatementBegin
ALTER TABLE board RENAME COLUMN member TO member_id;
ALTER TABLE board RENAME COLUMN year TO year_id;

ALTER TABLE board DROP CONSTRAINT board_member_fkey;
ALTER TABLE board
ADD CONSTRAINT fk_board_member
FOREIGN KEY (member_id)
REFERENCES member (id)
ON DELETE CASCADE;

ALTER TABLE board DROP CONSTRAINT board_year_fkey;
ALTER TABLE board
ADD CONSTRAINT fk_board_year
FOREIGN KEY (year_id)
REFERENCES year (id)
ON DELETE CASCADE;

ALTER TABLE board DROP COLUMN created_at;
ALTER TABLE board DROP COLUMN updated_at;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE board ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE board ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE board DROP CONSTRAINT board_year_fkey;
ALTER TABLE board
ADD CONSTRAINT fk_board_year
FOREIGN KEY (year_id)
REFERENCES year (id)

ALTER TABLE board DROP CONSTRAINT board_member_fkey;
ALTER TABLE board
ADD CONSTRAINT fk_board_member
FOREIGN KEY (member_id)
REFERENCES member (id)

ALTER TABLE board RENAME COLUMN year_id TO year;
ALTER TABLE board RENAME COLUMN member_id TO member;
-- +goose StatementEnd
