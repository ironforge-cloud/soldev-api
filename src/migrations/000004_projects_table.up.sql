CREATE TABLE IF NOT EXISTS projects
(
  id           SERIAL PRIMARY KEY                     NOT NULL,
  category_id  INT REFERENCES project_categories (id) NOT NULL,
  project_name VARCHAR(100)                           NOT NULL,
  logo         VARCHAR(300)                           NOT NULL,
  description  text                                   NOT NULL,
  status       VARCHAR(10)                            NOT NULL,
  project_url  jsonb                                  NOT NULL DEFAULT '{
    "data": []
  }'::jsonb,
  created_at   timestamp                              NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at   timestamp                              NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at   timestamp
)
