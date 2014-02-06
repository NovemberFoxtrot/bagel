DROP TABLE IF EXISTS cards;

CREATE TABLE cards ( 
  id          SERIAL PRIMARY KEY, 

  stack       VARCHAR(2000) NOT NULL,
  question    VARCHAR(2000) NOT NULL,
  answer      VARCHAR(2000) NOT NULL,
	correct     INT DEFAULT 0,
	incorrect   INT DEFAULT 0,

  updated_at  DATETIME NOT NULL,
  created_at  DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
