-- name: MailGetByID :many
SELECT *
FROM mail m
LEFT JOIN mail_event m_e ON m_e.mail_id = m.id
WHERE m.id = $1;

-- name: MailGetByYear :many
SELECT *
FROM mail m
LEFT JOIN mail_event m_e ON m_e.mail_id = m.id
WHERE m.year_id = $1
ORDER BY m.send_time;

-- name: MailGetByEvents :many
SELECT *
FROM mail m
LEFT JOIN mail_event m_e ON m_e.mail_id = m.id
WHERE m_e.event_id = ANY($1::int[])
ORDER BY send_time;

-- name: MailGetUnsend :many
SELECT *
FROM mail m
LEFT JOIN mail_event m_e ON m_e.mail_id = m.id
WHERE NOT send AND error IS NULL;

-- name: MailCreate :one
INSERT INTO mail (year_id, title, content, send_time, send, error)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id;

-- name: MailUpdate :exec
UPDATE mail
SET title = $1, content = $2, send_time = $3
WHERE id = $4 AND NOT send AND error IS NULL;

-- name: MailSend :exec
UPDATE mail
SET send = true
WHERE id = $1;

-- name: MailError :exec
UPDATE mail
SET error = $1
WHERE id = $2;

-- name: MailDelete :exec
DELETE FROM mail
WHERE id = $1
AND NOT send AND error IS NULL;
