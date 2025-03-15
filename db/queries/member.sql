-- name: MemberGetAll :many 
SELECT * FROM member;

-- name: MemberGetByID :one 
SELECT * FROM member 
WHERE id = $1;

-- name: MemberGetByName :one
SELECT * FROM member 
WHERE name ILIKE $1;

-- name: MemberCreate :one 
INSERT INTO member (name, username, zauth_id)
VALUES ($1, $2, $3)
RETURNING id;

-- name: MemberUpdate :exec 
UPDATE member 
SET name = $1, username = $2, zauth_id = $3
WHERE id = $4;

