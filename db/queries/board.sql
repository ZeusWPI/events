-- name: BoardGetAll :many
SELECT * FROM board b 
INNER JOIN member m ON b.member = m.id 
INNER JOIN academic_year a_y ON b.academic_year = a_y.id;

-- name: BoardCreate :one
INSERT INTO board (member, academic_year, role)
VALUES ($1, $2, $3)
RETURNING id;

