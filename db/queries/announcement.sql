-- name: AnnouncementGetAll :many
SELECT * FROM announcement
INNER JOIN member ON announcement.member = member.id;

-- name: AnnouncementGetByEvent :many
SELECT * FROM announcement
INNER JOIN member ON announcement.member = member.id
WHERE event = $1;

-- name: AnnouncementCreate :one
INSERT INTO announcement (content, time, target, event, member)
VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: AnnouncementUpdate :exec
UPDATE announcement
SET content = $1, time = $2, target = $3, sent = $4, event = $5, member = $6, updated_at = CURRENT_TIMESTAMP
WHERE id = $7;

-- name: AnnouncementDelete :exec
DELETE FROM announcement
WHERE id = $1;
