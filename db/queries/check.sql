-- name: CheckGetByEvents :many
SELECT sqlc.embed(e), sqlc.embed(c)
FROM check_event e
LEFT JOIN "check" c ON c.uid = e.check_uid
WHERE c.active AND event_id = ANY($1::int[]);

-- name: CheckGetByCheckUID :one
SELECT *
FROM "check"
WHERE uid = $1;

-- name: CheckGetByCheckEventID :one
SELECT sqlc.embed(e), sqlc.embed(c)
FROM check_event e
LEFT JOIN "check" c ON c.uid = e.check_uid
WHERE e.id = $1
LIMIT 1;

-- name: CheckGetByCheckUIDAll :many
SELECT sqlc.embed(e), sqlc.embed(c)
FROM check_event e
LEFT JOIN "check" c ON c.uid = e.check_uid
WHERE c.uid = $1;

-- name: CheckGetByCheckEvent :one
SELECT sqlc.embed(e), sqlc.embed(c)
FROM check_event e
LEFT JOIN "check" c ON c.uid = e.check_uid
WHERE c.uid = $1 AND e.event_id = $2;

-- name: CheckCreate :exec
INSERT INTO "check" (uid, description, deadline, active, "type", creator_id)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: CheckUpdate :exec
UPDATE "check"
SET description = $2, deadline = $3, active = $4, "type" = $5
WHERE uid = $1;

-- name: CheckSetInactiveAutomatic :exec
UPDATE "check"
SET active = false
WHERE "type" = 'automatic';

-- name: CheckDelete :exec
DELETE FROM "check"
WHERE uid = $1;
