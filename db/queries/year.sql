-- name: YearGetAll :many 
SELECT * 
FROM year
ORDER BY year_start DESC;  

-- name: YearGetLast :one 
SELECT * 
FROM year 
ORDER BY year_start DESC
LIMIT 1;

-- name: YearCreate :one 
INSERT INTO year (year_start, year_end)
VALUES ($1, $2)
RETURNING id;
