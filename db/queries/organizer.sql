-- name: OrganizerGetByYearWithBoard :many 
SELECT * FROM organizer o 
INNER JOIN board b ON b.id = o.board_id
INNER JOIN member m ON m.id = b.member_id
WHERE b.year_id = $1;

-- name: OrganizerCreate :one 
INSERT INTO organizer (event_id, board_id)
VALUES ($1, $2)
RETURNING id;

-- name: OrganizerDeleteByBoardEvent :exec 
DELETE FROM organizer 
WHERE board_id = $1 AND event_id = $2;
