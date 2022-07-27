CREATE TABLE IF NOT EXISTS documents_x (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(200) DEFAULT '' NOT NULL,
  alphnum_name VARCHAR(200) DEFAULT '' NOT NULL,
  url TEXT DEFAULT '' NOT NULL,
  type doctype DEFAULT 'file' NOT NULL,
  dir_id BIGINT DEFAULT 0 NOT NULL,
  is_private BOOLEAN DEFAULT false NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL,
  textsearchable_index_col tsvector GENERATED ALWAYS AS (
    to_tsvector('english',
      coalesce(alphnum_name, '')
      || ' ' || coalesce(name, '')
      || ' ' || coalesce(url, '')
    )
  ) STORED,
  textrank_index_col tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(alphnum_name, '')), 'A')
    || setweight(to_tsvector('english', coalesce(name, '')), 'B')
    || setweight(to_tsvector('english', coalesce(url, '')), 'C')
  ) STORED
);

INSERT INTO documents_x
  (name, alphnum_name, url, type, dir_id, is_private, created_at, updated_at, deleted_at)
(
    SELECT name, name, url, type, dir_id, is_private, created_at, updated_at, deleted_at
    FROM documents
);

DROP TABLE documents;

ALTER TABLE documents_x RENAME TO documents;
ALTER SEQUENCE documents_x_id_seq RENAME TO documents_id_seq;

CREATE TABLE IF NOT EXISTS images_x (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(200) DEFAULT '' NOT NULL,
  alphnum_name VARCHAR(200) DEFAULT '' NOT NULL,
  url TEXT DEFAULT '' NOT NULL,
  description TEXT DEFAULT '' NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL
);


INSERT INTO images_x
  (name, alphnum_name, url, description, created_at, deleted_at)
(
    SELECT name, alphnum_name, url, description, created_at, deleted_at
    FROM images
);

DROP TABLE images;

ALTER TABLE images_x RENAME TO images;
ALTER SEQUENCE images_x_id_seq RENAME TO images_id_seq;
