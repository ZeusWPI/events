-- name: BoardGetAllWithMemberYear :many
SELECT * FROM board b 
INNER JOIN member m ON b.member = m.id 
INNER JOIN year y ON b.year = y.id;

-- name: BoardGetByYearWithMemberYear :many 
SELECT * FROM board b 
INNER JOIN member m ON b.member = m.id 
INNER JOIN year y ON b.year = y.id
WHERE b.year = $1;

-- name: BoardGetByMemberYear :one 
SELECT * FROM board b 
INNER JOIN member m ON b.member = m.id 
INNER JOIN year y ON b.year = y.id
WHERE m.id = $1 AND y.id = $2;

-- name: BoardCreate :one
INSERT INTO board (member, year, role)
VALUES ($1, $2, $3)
RETURNING id;

