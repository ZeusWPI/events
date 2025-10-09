-- name: ImageGet :one
SELECT *
FROM image
WHERE id = $1;

-- name: ImageCreate :one
INSERT INTO image (name, file_id)
VALUES ($1, $2)
RETURNING id;
