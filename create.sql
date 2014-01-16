DROP TABLE IF EXISTS cards;

CREATE TABLE cards ( 
  id          SERIAL PRIMARY KEY, 
  data        VARCHAR(200) NOT NULL UNIQUE,
  # question    VARCHAR(200) NOT NULL UNIQUE,
  # answer      VARCHAR(200) NOT NULL UNIQUE,
  # explanation VARCHAR(200) NOT NULL UNIQUE,
  updated_at  DATETIME NOT NULL,
  created_at  DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

####

DROP TABLE IF EXISTS tags;

CREATE TABLE tags ( 
  id          SERIAL PRIMARY KEY, 
  data        VARCHAR(200) NOT NULL UNIQUE,
  updated_at  DATETIME NOT NULL,
  created_at  DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

####

DROP TABLE IF EXISTS users;

CREATE TABLE users ( 
  id          SERIAL PRIMARY KEY, 
  name        VARCHAR(200) NOT NULL,
	incorrect   INT DEFAULT 0,
  updated_at  DATETIME NOT NULL,
  created_at  DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

####

DROP TABLE IF EXISTS stats;

CREATE TABLE stats ( 
  id          SERIAL PRIMARY KEY, 
  user_id     INT NOT NULL,
  card_id     INT NOT NULL,
	correct     INT DEFAULT 0,
	incorrect   INT DEFAULT 0,
  updated_at  DATETIME NOT NULL,
  created_at  DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

####

DROP TABLE IF EXISTS cards_tags;

CREATE TABLE cards_tags (
  card_id     BIGINT UNSIGNED NOT NULL, 
  tag_id      BIGINT UNSIGNED NOT NULL, 
  updated_at  DATETIME NOT NULL,
  created_at  DATETIME NOT NULL,
  PRIMARY KEY (card_id, tag_id),
  KEY         fk_course_id (card_id),
  CONSTRAINT  fk_card_id FOREIGN KEY (card_id) REFERENCES cards(id),
  CONSTRAINT  fk_tag_id  FOREIGN KEY (tag_id)  REFERENCES tags(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
