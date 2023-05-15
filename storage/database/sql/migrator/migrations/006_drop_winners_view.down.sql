-- View: campaigns_winners
CREATE OR REPLACE VIEW campaigns_winners AS
    SELECT wc.campaign_id AS campaign_id,
           string_agg(wc.account_address, ',') AS winners
      FROM (
               SELECT w.account_address,
                      w.campaign_id AS campaign_id
                 FROM participants w
                WHERE w.position IS NOT NULL
                ORDER BY w.position ASC
           )
           AS wc
     GROUP BY wc.campaign_id;
