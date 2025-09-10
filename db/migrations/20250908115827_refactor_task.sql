-- +goose Up
-- +goose StatementBegin
DROP TABLE task; -- Ja kijk, jammer.

CREATE TYPE task_type AS ENUM ('recurring', 'once');

CREATE TABLE task (
  uid VARCHAR(255) NOT NULL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  active BOOLEAN NOT NULL,
  "type" TASK_TYPE NOT NULL
);

CREATE TABLE task_run (
  id SERIAL PRIMARY KEY,
  task_uid VARCHAR(255) NOT NULL REFERENCES task (uid) ON DELETE CASCADE,
  run_at TIMESTAMPTZ NOT NULL,
  result TASK_RESULT NOT NULL,
  error TEXT,
  duration BIGINT NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE task_run;
DROP TABLE task;

DROP TYPE task_type;

CREATE TABLE task (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  result TASK_RESULT NOT NULL,
  run_at TIMESTAMPTZ NOT NULL,
  error TEXT,
  recurring BOOLEAN NOT NULL,
duration INTERVAL NOT NULL DEFAULT '0 seconds'
);
-- +goose StatementEnd
