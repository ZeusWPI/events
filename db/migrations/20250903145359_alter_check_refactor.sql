-- +goose Up
-- +goose StatementBegin
DROP TABLE "check"; -- Dat is pech, data weg

CREATE TABLE "check" (
  id SERIAL PRIMARY KEY,
  description VARCHAR(255) NOT NULL,
  deadline BIGINT NOT NULL,

  UNIQUE (id, description)
);

CREATE TYPE check_status_enum AS ENUM ('success', 'failed', 'warning');

CREATE TABLE check_status (
  id SERIAL PRIMARY KEY,
  event_id INTEGER NOT NULL REFERENCES event (id) ON DELETE CASCADE,
  check_id INTEGER NOT NULL REFERENCES "check" (id) ON DELETE CASCADE,
  status CHECK_STATUS_ENUM NOT NULL,
  message VARCHAR(255),

  UNIQUE (event_id, check_id)
);

CREATE TABLE check_custom (
  id SERIAL PRIMARY KEY,
  event_id INTEGER NOT NULL REFERENCES event (id) ON DELETE CASCADE,
  description VARCHAR(255) NOT NULL,
  status CHECK_STATUS_ENUM NOT NULL,
  creator_id INTEGER NOT NULL REFERENCES board (id) ON DELETE CASCADE,

  UNIQUE (event_id, description)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE check_custom;
DROP TABLE check_result;
DROP TYPE check_status;
DROP TABLE "check";

CREATE TABLE "check" (
  id SERIAL PRIMARY KEY,
  event_id INTEGER NOT NULL REFERENCES event (id) ON DELETE CASCADE,
  description TEXT NOT NULL,
  done BOOLEAN NOT NULL,

  UNIQUE (event_id, description)
);
-- +goose StatementEnd
