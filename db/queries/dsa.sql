-- name: DsaGetByEvents :many
SELECT *
FROM dsa 
WHERE event_id = ANY($1::int[]);

-- name: DsaCreate :one
INSERT INTO dsa (event_id, entry)
VALUES ($1, $2)
RETURNING id;

-- name: DsaDelete :exec
DELETE FROM dsa 
WHERE id = $1;
