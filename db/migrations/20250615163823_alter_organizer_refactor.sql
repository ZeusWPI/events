-- +goose Up
-- +goose StatementBegin
ALTER TABLE organizer RENAME COLUMN event TO event_id;
ALTER TABLE organizer RENAME COLUMN board TO board_id;

ALTER TABLE organizer DROP CONSTRAINT organizer_event_fkey;
ALTER TABLE organizer
ADD CONSTRAINT fk_organizer_event
FOREIGN KEY (event_id)
REFERENCES event (id)
ON DELETE CASCADE;

ALTER TABLE organizer DROP CONSTRAINT organizer_board_fkey;
ALTER TABLE organizer
ADD CONSTRAINT fk_organizer_board
FOREIGN KEY (board_id)
REFERENCES board (id)
ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE organizer DROP CONSTRAINT organizer_board_fkey;
ALTER TABLE organizer
ADD CONSTRAINT fk_organizer_board
FOREIGN KEY (board_id)
REFERENCES board (id)

ALTER TABLE organizer DROP CONSTRAINT organizer_event_fkey;
ALTER TABLE organizer
ADD CONSTRAINT fk_organizer_event
FOREIGN KEY (event_id)
REFERENCES event (id)

ALTER TABLE organizer RENAME COLUMN event_id TO event;
ALTER TABLE organizer RENAME COLUMN board_id TO board;
-- +goose StatementEnd
