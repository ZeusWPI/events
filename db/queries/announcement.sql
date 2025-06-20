-- name: AnnouncementGetByEvents :many
SELECT * 
FROM announcement
WHERE event_id = ANY($1::int[])
ORDER BY send_time;

-- name: AnnouncementGetUnsend :many
SELECT *
FROM announcement
WHERE NOT send AND error IS NULL;

-- name: AnnouncementCreate :one 
INSERT INTO announcement (event_id, content, send_time, send, error)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: AnnouncementUpdate :exec
UPDATE announcement
SET content = $1, send_time = $2
WHERE id = $3 AND NOT send;

-- name: AnnouncementSend :exec 
UPDATE announcement
SET send = true
WHERE id = $1;

-- name: AnnouncementError :exec
UPDATE announcement
SET error = $1
WHERE id = $2;
