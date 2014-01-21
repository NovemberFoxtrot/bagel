DROP TABLE IF EXISTS cards_tags;
DROP TABLE IF EXISTS cards_users;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS users;

CREATE TABLE cards ( 
  id          SERIAL PRIMARY KEY, 

  question    VARCHAR(200) NOT NULL,
  answer      VARCHAR(200) NOT NULL,
  explanation VARCHAR(200) NOT NULL,

  updated_at  DATETIME NOT NULL,
  created_at  DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE tags ( 
  id          SERIAL PRIMARY KEY, 

  data        VARCHAR(200) NOT NULL,

  updated_at  DATETIME NOT NULL,
  created_at  DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE users ( 
  id          SERIAL PRIMARY KEY, 

  name        VARCHAR(200) NOT NULL,

  updated_at  DATETIME NOT NULL,
  created_at  DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE cards_users ( 
  card_id     BIGINT UNSIGNED NOT NULL,
  user_id     BIGINT UNSIGNED NOT NULL,

	correct     INT DEFAULT 0,
	incorrect   INT DEFAULT 0,

  updated_at  DATETIME NOT NULL,
  created_at  DATETIME NOT NULL,

	PRIMARY KEY (card_id, user_id),
  CONSTRAINT  fk_card_id FOREIGN KEY (card_id) REFERENCES cards(id),
  CONSTRAINT  fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE cards_tags (
  card_id     BIGINT UNSIGNED NOT NULL, 
  tag_id      BIGINT UNSIGNED NOT NULL, 

  updated_at  DATETIME NOT NULL,
  created_at  DATETIME NOT NULL,

  PRIMARY KEY (card_id, tag_id),
  CONSTRAINT  fk_card_tags_card_id FOREIGN KEY (card_id) REFERENCES cards(id),
  CONSTRAINT  fk_card_tags_tag_id  FOREIGN KEY (tag_id)  REFERENCES tags(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
