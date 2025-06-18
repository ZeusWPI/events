-- name: AnnouncementGetByEvents :many
SELECT * 
FROM announcement
WHERE event_id = ANY($1::int[]);

-- name: AnnouncementCreate :one 
INSERT INTO announcement (event_id, content, send_time, send)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: AnnouncementUpdate :exec
UPDATE announcement
SET content = $1, send_time = $2
WHERE id = $3 AND NOT send;

-- name: AnnouncementSend :exec 
UPDATE announcement
SET send = $1
WHERE id = $2;
