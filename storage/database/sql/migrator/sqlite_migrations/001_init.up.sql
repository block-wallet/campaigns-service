
-- Table: campaigns
CREATE TABLE IF NOT EXISTS campaigns (
    id             TEXT PRIMARY KEY
                        UNIQUE
                        NOT NULL,
    name           TEXT NOT NULL,
    description    TEXT NOT NULL,
    status         TEXT NOT NULL,
    start_date     TEXT NOT NULL,
    end_date       TEXT NOT NULL,
    enroll_message TEXT
);


-- Table: campaigns_supported_chains
CREATE TABLE IF NOT EXISTS campaigns_supported_chains (
    campaign_id TEXT REFERENCES campaigns (id),
    chain_id    TEXT NOT NULL,
    PRIMARY KEY (
        campaign_id,
        chain_id
    )
);


-- Table: campaigns_tags
CREATE TABLE IF NOT EXISTS campaigns_tags (
    campaign_id TEXT REFERENCES campaigns (id) 
                     NOT NULL,
    tag         TEXT NOT NULL,
    PRIMARY KEY (
        campaign_id,
        tag
    )
);


-- Table: participants
CREATE TABLE IF NOT EXISTS participants (
    campaign_id     TEXT    CONSTRAINT fk_campaigns_participants REFERENCES campaigns (id) MATCH SIMPLE
                            NOT NULL,
    account_address TEXT    NOT NULL,
    position        NUMERIC,
    CONSTRAINT pk_campaign_account PRIMARY KEY (
        campaign_id COLLATE NOCASE,
        account_address COLLATE NOCASE
    )
);


-- Table: rewards
CREATE TABLE IF NOT EXISTS rewards (
    reward_id   TEXT NOT NULL,
    campaign_id TEXT NOT NULL,
    token_id    TEXT REFERENCES tokens (token_id) 
                     NOT NULL,
    type        TEXT NOT NULL,
    amounts     TEXT,
    PRIMARY KEY (
        reward_id
    )
);


-- Table: tokens
CREATE TABLE IF NOT EXISTS tokens (
    id       TEXT    PRIMARY KEY,
    name     TEXT    NOT NULL,
    decimals NUMERIC NOT NULL,
    symbol   TEXT    NOT NULL
);


-- Table: tokens_contracts
CREATE TABLE IF NOT EXISTS tokens_contracts (
    token_id TEXT    REFERENCES tokens (id) 
                     NOT NULL,
    chain_id NUMERIC NOT NULL,
    address  TEXT    NOT NULL,
    PRIMARY KEY (
        token_id COLLATE NOCASE,
        chain_id COLLATE BINARY
    )
);


-- View: campaigns_winners
CREATE VIEW IF NOT EXISTS campaigns_winners AS
    SELECT wc.campaign_id AS campaign_id,
           group_concat(wc.account_address) AS winners
      FROM (
               SELECT w.account_address,
                      w.campaign_id AS campaign_id
                 FROM participants w
                WHERE w.position IS NOT NULL
                ORDER BY w.position ASC
           )
           AS wc
     GROUP BY wc.campaign_id;
