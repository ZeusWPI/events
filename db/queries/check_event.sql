-- name: CheckEventCreate :one
INSERT INTO check_event (check_uid, event_id, status, message)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: CheckEventCreateBatch :exec
INSERT INTO check_event (check_uid, event_id, status, message)
VALUES (
  UNNEST($1::varchar[]),
  UNNEST($2::int[]),
  UNNEST($3::text[])::check_status,
  UNNEST($4::varchar[])

);

-- name: CheckEventUpdate :exec
UPDATE check_event
SET status = $2, message = $3, updated_at = NOW()
WHERE id = $1;
