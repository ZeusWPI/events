-- name: MailGetAll :many
SELECT *
FROM mail
ORDER BY send_time;

-- name: MailGetUnsend :many
SELECT *
FROM mail
WHERE NOT send AND error IS NULL;

-- name: MailCreate :one
INSERT INTO mail (content, send_time, send, error)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: MailUpdate :exec
UPDATE mail
SET content = $1, send_time = $2
WHERE id = $3 AND NOT send;

-- name: MailSend :exec
UPDATE mail
SET send = true
WHERE id = $1;

-- name: MailError :exec
UPDATE mail
SET error = $1
WHERE id = $2;
