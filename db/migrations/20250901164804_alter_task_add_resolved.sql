-- +goose Up
-- +goose StatementBegin
CREATE TYPE task_result AS ENUM ('success', 'failed', 'resolved');

ALTER TABLE task
ADD COLUMN result_new TASK_RESULT NOT NULL DEFAULT 'failed';

ALTER TABLE task
ALTER COLUMN result_new DROP DEFAULT;

UPDATE task
SET result_new = result::task_result;

ALTER TABLE task DROP COLUMN result;

ALTER TABLE task RENAME COLUMN result_new TO result;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE task
ADD COLUMN result_old VARCHAR(255);

UPDATE task
SET result_old = result::text;

UPDATE task
SET result_old = result::text;

ALTER TABLE task DROP COLUMN result;

ALTER TABLE task RENAME COLUMN result_old TO result;

DROP TYPE task_result;
-- +goose StatementEnd
