-- name: OrganizerCreate :one 
INSERT INTO organizer (event_id, board_id)
VALUES ($1, $2)
RETURNING id;

-- name: OrganizerDeleteByBoardEvent :exec 
DELETE FROM organizer 
WHERE board_id = $1 AND event_id = $2;
