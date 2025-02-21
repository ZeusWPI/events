-- name: OrganizerGetAllByEvent :many 
SELECT * FROM organizer 
WHERE event = $1;

-- name: OrganizerCreate :one 
INSERT INTO organizer (event, board)
VALUES ($1, $2)
RETURNING id;

-- name: OrganizerDelete :exec 
DELETE FROM organizer 
WHERE id = $1;
