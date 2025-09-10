-- name: OrganizerGetByEvents :many
SELECT sqlc.embed(o), sqlc.embed(b), sqlc.embed(e), sqlc.embed(m), sqlc.embed(y)
FROM organizer o
LEFT JOIN event e ON e.id = o.event_id
LEFT JOIN year y ON y.id = e.year_id
LEFT JOIN board b ON b.id = o.board_id
LEFT JOIN member m ON m.id = b.member_id
WHERE e.id = ANY($1::int[]);

-- name: OrganizerCreateBatch :exec
INSERT INTO organizer (event_id, board_id)
VALUES (
  UNNEST($1::int[]),
  UNNEST($2::int[])
);

-- name: OrganizerDeleteByEvent :exec
DELETE FROM organizer
WHERE event_id = $1;
