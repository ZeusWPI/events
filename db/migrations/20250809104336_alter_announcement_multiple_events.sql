-- +goose Up
-- +goose StatementBegin
CREATE TABLE announcement_event (
  id SERIAL PRIMARY KEY,
  event_id INTEGER NOT NULL REFERENCES event (id) ON DELETE CASCADE,
  announcement_id INTEGER NOT NULL REFERENCES announcement (id) ON DELETE CASCADE
);

INSERT INTO announcement_event (event_id, announcement_id)
SELECT event_id, id
FROM announcement;

ALTER TABLE announcement
DROP COLUMN event_id;

CREATE FUNCTION check_announcement_event_year_match()
RETURNS TRIGGER AS $$
BEGIN
  IF (
    SELECT a.year_id
    FROM announcement a
    WHERE a.id = NEW.announcement_id
  ) IS DISTINCT FROM (
    SELECT e.year_id
    FROM event e
    WHERE e.id = NEW.event_id
  ) THEN
    RAISE EXCEPTION 'Year ID mismatch between announcement % and event %', NEW.announcement_id, NEW.event_id;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER announcement_event_year_check
BEFORE INSERT OR UPDATE ON announcement_event
FOR EACH ROW
EXECUTE FUNCTION check_announcement_event_year_match();
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TRIGGER announcement_event_year_check ON announcement_event;
DROP FUNCTION check_announcement_event_year_match();

ALTER TABLE announcement
ADD COLUMN event_id INTEGER REFERENCES event (id) ON DELETE CASCADE;

UPDATE announcement a
SET event_id = (
  SELECT event_id
  FROM announcement_event a_e
  WHERE a_e.announcement_id = a.id
  LIMIT 1
);

ALTER TABLE announcement
ALTER COLUMN event_id SET NOT NULL;

DROP TABLE announcement_event;
-- +goose StatementEnd
