-- name: BoardGetAllPopulated :many
SELECT * FROM board b 
INNER JOIN member m ON b.member_id = m.id 
INNER JOIN year y ON b.year_id = y.id;

-- name: BoardGetByYearPopulated :many 
SELECT * FROM board b 
INNER JOIN member m ON b.member_id = m.id 
INNER JOIN year y ON b.year_id = y.id
WHERE b.year_id = $1;

-- name: BoardGetByMemberYear :one 
SELECT * FROM board b 
INNER JOIN member m ON b.member_id = m.id 
INNER JOIN year y ON b.year_id = y.id
WHERE m.id = $1 AND y.id = $2;

-- name: BoardCreate :one
INSERT INTO board (role, member_id, year_id)
VALUES ($1, $2, $3)
RETURNING id;

