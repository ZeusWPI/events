-- name: MemberGetAll :many 
SELECT * FROM member;

-- name: MemberCreate :one 
INSERT INTO member (name, username)
VALUES ($1, $2)
RETURNING id;

-- name: MemberUpdate :exec 
UPDATE member 
SET name = $1, username = $2
WHERE id = $3;

