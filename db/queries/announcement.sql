-- name: AnnouncementGetByID :many
SELECT *
FROM announcement a
LEFT JOIN announcement_event a_e ON a_e.announcement_id = a.id
WHERE a.id = $1;

-- name: AnnouncementGetAll :many
SELECT *
FROM announcement a
LEFT JOIN announcement_event a_e ON a_e.announcement_id = a.id;

-- name: AnnouncmentGetByYear :many
SELECT *
FROM announcement a
LEFT JOIN announcement_event a_e ON a_e.announcement_id = a.id
WHERE a.year_id = $1
ORDER BY a.send_time;

-- name: AnnouncementGetByEvents :many
SELECT * 
FROM announcement a
LEFT JOIN announcement_event a_e ON a_e.announcement_id = a.id
WHERE a_e.event_id = ANY($1::int[])
ORDER BY send_time;

-- name: AnnouncementGetUnsend :many
SELECT *
FROM announcement a
LEFT JOIN announcement_event a_e ON a_e.announcement_id = a.id
WHERE NOT send AND error IS NULL;

-- name: AnnouncementCreate :one 
INSERT INTO announcement (year_id, author_id, content, send_time, draft, send, error)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: AnnouncementUpdate :exec
UPDATE announcement
SET content = $2, send_time = $3, draft = $4, error = $5
WHERE id = $1 AND NOT send;

-- name: AnnouncementSend :exec 
UPDATE announcement
SET send = true
WHERE id = $1;

-- name: AnnouncementDelete :exec
DELETE FROM announcement
WHERE id = $1
AND NOT send AND error IS NULL;
