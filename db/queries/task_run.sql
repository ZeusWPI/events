-- name: TaskRunCreate :one
INSERT INTO task_run (task_uid, run_at, result, error, duration)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: TaskRunResolve :exec
UPDATE task_run
SET result = 'resolved'
WHERE id = $1;
