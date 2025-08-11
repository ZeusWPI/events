-- name: AnnouncementEventCreateBatch :exec
INSERT INTO announcement_event (announcement_id, event_id)
VALUES (
  UNNEST($1::int[]),
  UNNEST($2::int[])
);

-- name: AnnouncementEventDeleteByAnnouncement :exec
DELETE FROM announcement_event
WHERE announcement_id = $1;
