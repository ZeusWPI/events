-- name: BoardGetAllWithMemberYear :many
SELECT * FROM board b 
INNER JOIN member m ON b.member = m.id 
INNER JOIN year a_y ON b.year = a_y.id;

-- name: BoardGetByYearWithMemberYear :many 
SELECT * FROM board b 
INNER JOIN member m ON b.member = m.id 
INNER JOIN year a_y ON b.year = a_y.id
WHERE b.year = $1;

-- name: BoardCreate :one
INSERT INTO board (member, year, role)
VALUES ($1, $2, $3)
RETURNING id;

