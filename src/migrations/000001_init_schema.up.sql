CREATE TABLE IF NOT EXISTS companies
(
  id          serial PRIMARY KEY NOT NULL,
  name        VARCHAR(50)        NOT NULL,
  description VARCHAR(300)       NOT NULL,
  logo        VARCHAR(100)       NOT NULL,
  bg_color    VARCHAR(10)        NOT NULL,
  status      VARCHAR(10)        NOT NULL DEFAULT 'inactive',
  created_at  timestamp          NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  timestamp          NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at  timestamp
);

CREATE TABLE IF NOT EXISTS bounties
(
  id            serial PRIMARY KEY            NOT NULL,
  company_id    INT REFERENCES companies (id) NOT NULL,
  title         VARCHAR(100)                  NOT NULL,
  description   text                          NOT NULL,
  reward        INT                           NOT NULL,
  reward_asset  VARCHAR(50)                   NOT NULL,
  url           VARCHAR(300)                  NOT NULL,
  tags          jsonb                         not null default '{
    "names": []
  }'::jsonb,
  status        VARCHAR(10)                   NOT NULL DEFAULT 'active',
  delivery_date VARCHAR(50)                   NOT NULL,
  created_at    timestamp                     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    timestamp                     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at    timestamp
);
