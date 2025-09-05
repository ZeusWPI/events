-- name: CheckCustomGet :one
SELECT *
FROM check_custom
WHERE id = $1;

-- name: CheckCustomGetByEvent :many
SELECT *
FROM check_custom
WHERE event_id = $1;

-- name: CheckCustomCreate :one
INSERT INTO check_custom (event_id, description, status, creator_id)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: CheckCustomUpdate :exec
UPDATE check_custom
SET description = $2, status = $3
WHERE id = $1;
