-- name: TaskGetByUID :one
SELECT *
FROM task
WHERE uid = $1;

-- name: TaskRunGet :one
SELECT sqlc.embed(t), sqlc.embed(r)
FROM task_run r
LEFT JOIN task t ON t.uid = r.task_uid
WHERE r.id = $1;

-- name: TaskGetFiltered :many
SELECT sqlc.embed(t), sqlc.embed(r)
FROM task_run r
LEFT JOIN task t ON t.uid = r.task_uid
WHERE
  (t.uid = $1 OR NOT @filter_task_uid) AND
  (r.result = $2 OR NOT @filter_result) AND
  t.active
ORDER BY r.run_at DESC
LIMIT $3 OFFSET $4;

-- name: TaskCreate :exec
INSERT INTO task (uid, name, active, "type")
VALUES ($1, $2, $3, $4);

-- name: TaskUpdate :exec
UPDATE task
SET name = $2, active = $3
WHERE uid = $1;

-- name: TaskSetInactiveRecurring :exec
UPDATE task
SET active = false
WHERE "type" = 'recurring';
