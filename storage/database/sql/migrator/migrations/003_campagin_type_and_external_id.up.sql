ALTER TABLE IF EXISTS campaigns ADD COLUMN campaign_type text NOT NULL DEFAULT 'PARTNER_OFFERS';

ALTER TABLE IF EXISTS campaigns ADD COLUMN external_campaign_id text;