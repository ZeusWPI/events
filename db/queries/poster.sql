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
INSERT INTO poster (event_id, file_id, webp_id, scc)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: PosterUpdate :exec
UPDATE poster
SET event_id = $1, file_id = $2, webp_id = $3, scc = $4
WHERE id = $5;

-- name: PosterDelete :exec
DELETE FROM poster
WHERE id = $1;
