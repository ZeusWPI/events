-- name: AcademicYearGetAll :many 
SELECT * FROM academic_year
ORDER BY start_year DESC;  

-- name: AcademicYearCreate :one 
INSERT INTO academic_year (start_year, end_year)
VALUES ($1, $2)
RETURNING id;

-- name: AcademicYearUpdate :exec
UPDATE academic_year
SET start_year = $1, end_year = $2 
WHERE id = $3;

