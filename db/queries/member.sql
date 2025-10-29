-- name: MemberGetAll :many 
SELECT * FROM member;

-- name: MemberGetByID :one 
SELECT * FROM member 
WHERE id = $1;

-- name: MemberGetByName :one
SELECT * FROM member 
WHERE name ILIKE $1;

-- name: MemberCreate :one 
INSERT INTO member (name, username, mattermost, zauth_id)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: MemberUpdate :exec 
UPDATE member 
SET name = $2, username = $3, mattermost = $4, zauth_id = $5
WHERE id = $1;

-- name: MemberDelete :exec 
DELETE FROM member 
WHERE id = $1;
