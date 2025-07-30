-- name: BoardGetAllPopulated :many
SELECT * 
FROM board b
INNER JOIN member m ON b.member_id = m.id 
INNER JOIN year y ON b.year_id = y.id;

-- name: BoardGetByIds :many
SELECT * 
FROM board
WHERE id = ANY($1::int[]);

-- name: BoardGetByYearPopulated :many 
SELECT * 
FROM board b 
INNER JOIN member m ON b.member_id = m.id 
INNER JOIN year y ON b.year_id = y.id
WHERE b.year_id = $1;

-- name: BoardGetByMemberYear :one 
SELECT * 
FROM board b 
INNER JOIN member m ON b.member_id = m.id 
INNER JOIN year y ON b.year_id = y.id
WHERE m.id = $1 AND y.id = $2;

-- name: BoardGetByMemberID :many 
SELECT * 
FROM board
WHERE member_id = $1;

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
