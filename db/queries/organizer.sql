-- name: OrganizerGetAllByEvent :many 
SELECT * FROM organizer 
INNER JOIN event ON organizer.event = event.id
INNER JOIN board ON organizer.board = board.id
WHERE event = $1;

-- name: OrganizerCreate :one 
INSERT INTO organizer (event, board)
VALUES ($1, $2)
RETURNING id;

-- name: OrganizerDelete :exec 
DELETE FROM organizer 
WHERE id = $1;
