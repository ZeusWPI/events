-- name: OrganizerCreate :one 
INSERT INTO organizer (event_id, board_id)
VALUES ($1, $2)
RETURNING id;

-- name: OrganizerCreateBatch :exec
INSERT INTO organizer (event_id, board_id)
VALUES (
  UNNEST($1::int[]),
  UNNEST($2::int[])
);

-- name: OrganizerDeleteByBoardEvent :exec 
DELETE FROM organizer 
WHERE board_id = $1 AND event_id = $2;

-- name: OrganizerDeleteByEvent :exec
DELETE FROM organizer
WHERE event_id = $1;
