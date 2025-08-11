-- +goose Up
-- +goose StatementBegin
ALTER TABLE mail
ADD COLUMN year_id INTEGER NOT NULL REFERENCES year (id) ON DELETE CASCADE;

CREATE FUNCTION check_mail_event_year_match()
RETURNS TRIGGER AS $$
BEGIN
  IF (
    SELECT m.year_id
    FROM mail m
    WHERE m.id = NEW.mail_id
  ) IS DISTINCT FROM (
    SELECT e.year_id
    FROM event e
    WHERE e.id = NEW.event_id
  ) THEN
    RAISE EXCEPTION 'Year ID mismatch between mail % and event %', NEW.mail_id, NEW.event_id;
  END IF;

  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER mail_event_year_check
BEFORE INSERT OR UPDATE ON mail_event
FOR EACH ROW
EXECUTE FUNCTION check_mail_event_year_match();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER mail_event_year_check ON mail_event;
DROP FUNCTION check_mail_event_year_match();

ALTER TABLE mail
DROP COLUMN year_id;
-- +goose StatementEnd
