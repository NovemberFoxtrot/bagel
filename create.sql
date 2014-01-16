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
  card_id INT NOT NULL, 
  tag_id INT NOT NULL, 
  updated_at  DATETIME NOT NULL,
  created_at  DATETIME NOT NULL,
  PRIMARY KEY (card_id, tag_id)
# FOREIGN KEY (card_id) REFERENCES cards(id) ON UPDATE CASCADE, 
# FOREIGN KEY (tag_id)  REFERENCES tag(id)   ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;