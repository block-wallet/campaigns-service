ALTER TABLE IF EXISTS campaigns DROP COLUMN created_at;

ALTER TABLE IF EXISTS campaigns DROP COLUMN updated_at;

ALTER TABLE IF EXISTS participants DROP COLUMN enrolled_at;
