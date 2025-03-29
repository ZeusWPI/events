-- name: TaskGet :many 
SELECT * FROM task 
WHERE name ILIKE $1
ORDER BY run_at DESC;

-- name: TaskGetAll :many
SELECT * FROM task
ORDER BY run_at DESC;

-- name: TaskCreate :one 
INSERT INTO task (name, result, run_at, error, recurring)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;
