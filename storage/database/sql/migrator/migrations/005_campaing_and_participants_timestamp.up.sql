ALTER TABLE IF EXISTS campaigns
    ADD COLUMN created_at timestamp without time zone DEFAULT current_timestamp;

ALTER TABLE IF EXISTS campaigns
    ADD COLUMN updated_at timestamp without time zone DEFAULT current_timestamp;

ALTER TABLE IF EXISTS participants
    ADD COLUMN enrolled_at timestamp without time zone DEFAULT current_timestamp;
