-- +goose Up
-- +goose StatementBegin
CREATE TABLE announcement (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    time TIMESTAMPTZ NOT NULL,
    target VARCHAR(20),
    sent BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    event INT NOT NULL REFERENCES event (id),
    member INT REFERENCES member (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE announcement;
-- +goose StatementEnd
