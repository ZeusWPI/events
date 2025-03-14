-- name: EventGetAllWithYear :many 
SELECT * FROM event
INNER JOIN year ON event.year = year.id 
WHERE event.deleted_at IS NULL;

-- name: EventGetByYearWithYear :many 
SELECT * FROM event e
INNER JOIN year y ON y.id = e.year
WHERE y.id = $1; 

-- name: EventCreate :one 
INSERT INTO event (url, name, description, start_time, end_time, year, location)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: EventUpdate :exec
UPDATE event 
SET url = $1, name = $2, description = $3, start_time = $4, end_time = $5, year = $6, location = $7, updated_at = CURRENT_TIMESTAMP, deleted_at = NULL
WHERE id = $8;

-- name: EventDelete :exec
UPDATE event 
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1;
