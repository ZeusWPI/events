-- name: TaskGet :one
SELECT *
FROM task
WHERE id = $1;

-- name: TaskGetAll :many
SELECT * 
FROM task
ORDER BY run_at DESC;

-- name: TaskCreate :one 
INSERT INTO task (name, result, run_at, error, recurring, duration)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: TaskUpdateResult :exec
UPDATE task
SET result = $1
WHERE id = $2;
