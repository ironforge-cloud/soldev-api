CREATE TABLE IF NOT EXISTS project_categories
(
  id            SERIAL PRIMARY KEY NOT NULL,
  category_name VARCHAR(100)       NOT NULL,
  description   text               NOT NULL,
  status        VARCHAR(10)        NOT NULL,
  created_at    timestamp          NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    timestamp          NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at    timestamp
)
