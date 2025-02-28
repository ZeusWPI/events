-- name: OrganizerGetByYearWithBoard :many 
SELECT * FROM organizer o 
INNER JOIN board b ON b.id = o.board 
INNER JOIN member m ON m.id = b.member 
WHERE b.year = $1;

-- name: OrganizerCreate :one 
INSERT INTO organizer (event, board)
VALUES ($1, $2)
RETURNING id;

-- name: OrganizerDelete :exec 
DELETE FROM organizer 
WHERE id = $1;
