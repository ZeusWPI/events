-- name: PosterGetAll :many
SELECT *
FROM poster;

-- name: PosterGetByEvents :many
SELECT *
FROM poster
WHERE event_id = ANY($1::int[]);

-- name: PosterGet :one
SELECT *
FROM poster
WHERE id = $1;

-- name: PosterCreate :one
INSERT INTO poster (event_id, file_id, scc)
VALUES ($1, $2, $3)
RETURNING id;

-- name: PosterUpdate :exec
UPDATE poster
SET event_id = $1, file_id = $2, scc = $3
WHERE id = $4;

-- name: PosterDelete :exec
DELETE FROM poster
WHERE id = $1;
