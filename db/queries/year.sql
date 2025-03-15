-- name: YearGetAll :many 
SELECT * FROM year
ORDER BY start_year DESC;  

-- name: YearGetLatest :one 
SELECT * FROM year 
ORDER BY start_year DESC
LIMIT 1;

-- name: YearCreate :one 
INSERT INTO year (start_year, end_year)
VALUES ($1, $2)
RETURNING id;

-- name: YearUpdate :exec
UPDATE year
SET start_year = $1, end_year = $2 
WHERE id = $3;

