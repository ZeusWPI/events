-- name: BoardGetByMemberYear :one 
SELECT sqlc.embed(b), sqlc.embed(m), sqlc.embed(y)
FROM board b
LEFT JOIN member m ON b.member_id = m.id
LEFT JOIN year y ON b.year_id = y.id
WHERE m.id = $1 AND y.id = $2;

-- name: BoardGetAll :many
SELECT sqlc.embed(b), sqlc.embed(m), sqlc.embed(y)
FROM board b
LEFT JOIN member m ON b.member_id = m.id 
LEFT JOIN year y ON b.year_id = y.id;

-- name: BoardGetByIds :many
SELECT sqlc.embed(b), sqlc.embed(m), sqlc.embed(y)
FROM board b
LEFT JOIN member m ON b.member_id = m.id 
LEFT JOIN year y ON b.year_id = y.id
WHERE b.id = ANY($1::INT[]);

-- name: BoardGetByYear :many 
SELECT sqlc.embed(b), sqlc.embed(m), sqlc.embed(y)
FROM board b
LEFT JOIN member m ON b.member_id = m.id
LEFT JOIN year y ON b.year_id = y.id
WHERE b.year_id = $1;

-- name: BoardGetByMember :many 
SELECT sqlc.embed(b), sqlc.embed(m), sqlc.embed(y)
FROM board b
LEFT JOIN member m ON b.member_id = m.id
LEFT JOIN year y ON b.year_id = y.id
WHERE m.id = $1;

-- name: BoardCreate :one
INSERT INTO board (role, member_id, year_id, is_organizer)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: BoardUpdate :exec
UPDATE board
SET role = $1, member_id = $2, year_id = $3, is_organizer = $4
WHERE id = $5;

-- name: BoardDelete :exec
DELETE FROM board
WHERE id = $1;
