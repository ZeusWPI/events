-- name: MailEventGetByEvents :many
SELECT *
FROM mail_event
WHERE event_id = ANY($1::int[]);

-- name: MailEventCreateBatch :exec
INSERT INTO mail_event (mail_id, event_id)
VALUES (
  UNNEST($1::int[]),
  UNNEST($2::int[])
);

-- name: MailEventDeleteByMail :exec
DELETE FROM mail_event
WHERE mail_id = $1;

