-- name: EventGetById :one 
SELECT * 
FROM event 
WHERE id = $1;

-- name: EventGetByIds :many
SELECT *
FROM event
WHERE id = ANY($1::int[]);

-- name: EventGetAllWithYear :many 
SELECT * FROM event
INNER JOIN year ON event.year_id = year.id;

-- name: EventGetByYearPopulated :many
SELECT jsonb_build_object(
  'id', e.id,
  'file_name', e.file_name,
  'name', e.name,
  'description', e.description,
  'start_time', e.start_time,
  'end_time', e.end_time,
  'year_id', e.year_id,
  'location', e.location,
  'year', (
    SELECT jsonb_build_object(
      'id', y.id,
      'year_start', y.year_start,
      'year_end', y.year_end
    )
    FROM year y
    WHERE y.id = e.year_id
  ),
  'organizers', (
    SELECT coalesce(json_agg(jsonb_build_object(
      'id', b.id,
      'member_id', b.member_id,
      'year_id', b.year_id,
      'role', b.role,
      'member', (
        SELECT jsonb_build_object(
          'id', m.id,
          'name', m.name,
          'username', m.username,
          'zauth_id', m.zauth_id
        )
        FROM member m 
        WHERE m.id = b.member_id
      ),
      'year', (
        SELECT jsonb_build_object(
          'id', y.id,
          'year_start', y.year_start,
          'year_end', y.year_end
        )
        FROM year y
        WHERE y.id = b.year_id
      )
    )), '[]')
    FROM board b 
    INNER JOIN organizer o ON o.board_id = b.id
    WHERE o.event_id = e.id
  )
)
FROM event e 
WHERE e.year_id = $1;

-- name: EventCreate :one 
INSERT INTO event (file_name, name, description, start_time, end_time, year_id, location)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id;

-- name: EventUpdate :exec
UPDATE event 
SET name = $1, description = $2, start_time = $3, end_time = $4, year_id = $5, location = $6
WHERE id = $7;

-- name: EventDelete :exec
DELETE FROM event 
WHERE id = $1;
