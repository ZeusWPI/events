-- name: DsaGetByEvents :many
SELECT *
FROM dsa 
WHERE event_id = ANY($1::int[]);

-- name: DsaCreate :one
INSERT INTO dsa (event_id, dsa_id)
VALUES ($1, $2)
RETURNING id;

-- name: DsaDelete :exec
DELETE FROM dsa 
WHERE id = $1;

-- name: DsaUpdate :exec
UPDATE dsa
SET event_id = $1, dsa_id = $2
WHERE id = $3;
