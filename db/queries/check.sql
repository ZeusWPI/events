-- name: CheckGetByEvents :many
SELECT *
FROM "check"
WHERE event_id = ANY($1::int[]);

-- name: CheckCreate :one 
INSERT INTO "check" (event_id, description, done)
VALUES ($1, $2, $3)
RETURNING Id;

-- name: CheckToggle :exec 
UPDATE "check"
SET done = NOT done
WHERE id = $1;

-- name: CheckDelete :exec 
DELETE FROM "check"
WHERE id = $1;
