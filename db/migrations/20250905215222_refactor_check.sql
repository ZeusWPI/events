-- +goose Up
-- +goose StatementBegin
DROP TABLE "check"; -- Da's pech, data weg

CREATE TYPE check_status AS ENUM ('done', 'done_late', 'todo', 'todo_late', 'warning');
CREATE TYPE check_type AS ENUM ('manual', 'automatic');

CREATE TABLE "check" (
  uid VARCHAR(255) NOT NULL PRIMARY KEY,
  description VARCHAR(255) NOT NULL,
  deadline BIGINT,
  active BOOLEAN NOT NULL,
  "type" CHECK_TYPE NOT NULL,
  creator_id INTEGER REFERENCES board (id) ON DELETE CASCADE
);

CREATE TABLE check_event (
  id SERIAL PRIMARY KEY,
  check_uid VARCHAR(255) NOT NULL REFERENCES "check" (uid) ON DELETE CASCADE,
  event_id INTEGER NOT NULL REFERENCES event (id) ON DELETE CASCADE,
  status CHECK_STATUS NOT NULL,
  message VARCHAR(255),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  UNIQUE (event_id, check_uid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE check_event;
DROP TABLE "check";

DROP TYPE check_status;
DROP TYPE check_type;

CREATE TABLE "check" (
  id SERIAL PRIMARY KEY,
  event_id INTEGER NOT NULL REFERENCES event (id) ON DELETE CASCADE,
  description TEXT NOT NULL,
  done BOOLEAN NOT NULL,

  UNIQUE (event_id, description)
);
-- +goose StatementEnd
