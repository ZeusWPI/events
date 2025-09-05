-- name: CheckGetAll :many
SELECT *
FROM "check";

-- name: CheckCreate :one
INSERT INTO "check" (description, deadline)
VALUES ($1, $2)
RETURNING id;

