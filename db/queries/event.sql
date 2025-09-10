-- name: EventGet :one 
SELECT sqlc.embed(e), sqlc.embed(y)
FROM event e
LEFT JOIN year y ON y.id = e.year_id
WHERE e.id = $1;

-- name: EventGetByIds :many
SELECT sqlc.embed(e), sqlc.embed(y)
FROM event e
LEFT JOIN year y ON y.id = e.year_id
WHERE e.id = ANY($1::int[])
ORDER BY e.start_time;

-- name: EventGetAll :many 
SELECT sqlc.embed(e), sqlc.embed(y)
FROM event e
LEFT JOIN year y ON y.id = e.year_id
WHERE NOT deleted
ORDER BY e.start_time;

-- name: EventGetByYear :many
SELECT sqlc.embed(e), sqlc.embed(y)
FROM event e
LEFT JOIN year y ON y.id = e.year_id
WHERE e.year_id = $1 AND NOT deleted
ORDER BY e.start_time;

-- name: EventGetFuture :many
SELECT sqlc.embed(e), sqlc.embed(y)
FROM event e
LEFT JOIN year y ON e.year_id = y.id
WHERE e.start_time > NOW() AND NOT deleted
ORDER BY e.start_time;

-- name: EventGetNext :one
SELECT sqlc.embed(e), sqlc.embed(y)
FROM event e
INNER JOIN year y ON e.year_id = y.id
WHERE e.end_time > NOW() AND NOT deleted
ORDER BY e.start_time
LIMIT 1;

-- name: EventCreate :one 
INSERT INTO event (file_name, name, description, start_time, end_time, year_id, location, deleted)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: EventUpdate :exec
UPDATE event 
SET name = $2, description = $3, start_time = $4, end_time = $5, year_id = $6, location = $7, deleted = $8
WHERE id = $1;

-- name: EventDelete :exec
UPDATE event
SET deleted = true
WHERE id = $1;
