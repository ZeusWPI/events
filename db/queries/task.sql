-- name: TaskGet :one
SELECT *
FROM task
WHERE id = $1;

-- name: TaskGetFiltered :many
SELECT *
FROM task
WHERE
  (name = $1 OR NOT @filter_name) AND
  (result = $2 OR NOT @filter_result)
ORDER BY run_at DESC
LIMIT $3 OFFSET $4;

-- name: TaskCreate :one 
INSERT INTO task (name, result, run_at, error, recurring, duration)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: TaskUpdateResult :exec
UPDATE task
SET result = $1
WHERE id = $2;
