CREATE TABLE IF NOT EXISTS tags (
  tag_id      SERIAL PRIMARY KEY,
  data        VARCHAR(200) NOT NULL UNIQUE,
  updated_at  TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  created_at  DATETIME DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
