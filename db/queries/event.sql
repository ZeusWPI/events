-- CRUD 

-- name: EventGetAll :many 
SELECT * FROM event;

-- name: EventCreate :one 
INSERT INTO event (url, name, description, start_time, end_time, academic_year, location)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: EventUpdate :one 
UPDATE event 
SET url = $1, name = $2, description = $3, start_time = $4, end_time = $5, academic_year = $6, location = $7, updated_at = CURRENT_TIMESTAMP, deleted_at = NULL
WHERE id = $8
RETURNING *;

-- name: EventDelete :one
UPDATE event 
SET deleted_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
