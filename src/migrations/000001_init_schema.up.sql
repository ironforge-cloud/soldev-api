CREATE TABLE IF NOT EXISTS companies
(
  id          serial PRIMARY KEY NOT NULL,
  name        VARCHAR(50)        NOT NULL,
  description VARCHAR(300)       NOT NULL,
  logo        VARCHAR(100)       NOT NULL,
  bg_color    VARCHAR(10)        NOT NULL,
  created_at  timestamp          NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at  timestamp          NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at  timestamp
);

CREATE TABLE IF NOT EXISTS bounties
(
  id            serial PRIMARY KEY            NOT NULL,
  company_id    INT REFERENCES companies (id) NOT NULL,
  title         VARCHAR(50)                   NOT NULL,
  description   VARCHAR(300)                  NOT NULL,
  reward        INT                           NOT NULL,
  reward_asset  VARCHAR(50)                   NOT NULL,
  status        VARCHAR(10)                   NOT NULL,
  delivery_date VARCHAR(50)                   NOT NULL,
  created_at    timestamp                     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    timestamp                     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at    timestamp
);
