-- name: CheckStatusGetByEvent :many
SELECT *
FROM check_status
WHERE event_id = $1;

-- name: CheckStatusCreate :one
INSERT INTO check_status (event_id, check_id, status, message)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: CheckStatusUpdate :exec
UPDATE check_status
SET status = $2, message = $3
WHERE id = $1;
