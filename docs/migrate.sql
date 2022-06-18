CREATE TABLE IF NOT EXISTS historiesx (
  id BIGSERIAL PRIMARY KEY,
  content jsonb DEFAULT '{}'::jsonb NOT NULL,
  content_text TEXT DEFAULT '' NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
INSERT INTO historiesx
  (id, content, created_at)
(
    SELECT id, content, created_at
    FROM histories
);

CREATE TABLE IF NOT EXISTS goalsx (
  id BIGSERIAL PRIMARY KEY,
  vision jsonb DEFAULT '{}'::jsonb NOT NULL,
  vision_text TEXT DEFAULT '' NOT NULL,
  mission jsonb DEFAULT '{}'::jsonb NOT NULL,
  mission_text TEXT DEFAULT '' NOT NULL,
  org_period_id BIGINT NOT NULL REFERENCES org_periods(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
INSERT INTO goalsx
  (id, vision, mission, org_period_id, created_at)
(
    SELECT id, vision, mission, org_period_id, created_at
    FROM goals
);

CREATE TABLE IF NOT EXISTS blogsx (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(200) DEFAULT '' NOT NULL,
  short_desc VARCHAR(200) DEFAULT '' NOT NULL,
  thumbnail_url TEXT DEFAULT '' NOT NULL,
  content jsonb DEFAULT '{}'::jsonb NOT NULL,
  content_text TEXT DEFAULT '' NOT NULL,
  slug VARCHAR(200) DEFAULT '' NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL,
  textsearchable_index_col tsvector GENERATED ALWAYS AS (
    to_tsvector('english',
      coalesce(title, '') 
      || ' ' || coalesce(short_desc, '')
      || ' ' || coalesce(content_text, '')
      || ' ' || coalesce(slug, '')
    )
  ) STORED,
  textrank_index_col tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(title, '')), 'A')
    || setweight(to_tsvector('english', coalesce(short_desc, '')), 'B')
    || setweight(to_tsvector('english', coalesce(content_text, '')), 'C')
    || setweight(to_tsvector('english', coalesce(slug, '')), 'D')
  ) STORED
);
INSERT INTO blogsx
  (id, title, short_desc, thumbnail_url, content, slug, created_at, updated_at, deleted_at)
(
    SELECT id, title, short_desc, thumbnail_url, content, slug, created_at, updated_at, deleted_at
    FROM blogs
);

CREATE TABLE IF NOT EXISTS documentsx (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(200) DEFAULT '' NOT NULL,
  url VARCHAR(200) DEFAULT '' NOT NULL,
  type doctype DEFAULT 'file' NOT NULL,
  dir_id BIGINT DEFAULT 0 NOT NULL,
  is_private BOOLEAN DEFAULT false NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL,
  textsearchable_index_col tsvector GENERATED ALWAYS AS (
    to_tsvector('english',
      coalesce(name, '') 
      || ' ' || coalesce(url, '')
    )
  ) STORED,
  textrank_index_col tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(name, '')), 'A')
    || setweight(to_tsvector('english', coalesce(url, '')), 'B')
  ) STORED
);
INSERT INTO documentsx
  (id, name, url, type, dir_id, is_private, created_at, updated_at, deleted_at)
(
    SELECT id, name, url, type, dir_id, is_private, created_at, updated_at, deleted_at
    FROM documents
);

CREATE TABLE IF NOT EXISTS membersx (
  id UUID PRIMARY KEY,
  name VARCHAR(100) DEFAULT '' NOT NULL,
  wa_phone VARCHAR(50) DEFAULT '' NOT NULL UNIQUE,
  other_phone VARCHAR(50) DEFAULT '' NOT NULL UNIQUE,
  homestay_name VARCHAR(100) DEFAULT '' NOT NULL,
  homestay_address VARCHAR(200) DEFAULT '' NOT NULL,
  homestay_latitude VARCHAR(50) DEFAULT '' NOT NULL,
  homestay_longitude VARCHAR(50) DEFAULT '' NOT NULL,
  profile_pic_url TEXT DEFAULT '' NOT NULL,
  username VARCHAR(50) DEFAULT '' NOT NULL UNIQUE,
  password VARCHAR(200) DEFAULT '' NOT NULL,
  is_admin BOOLEAN DEFAULT false NOT NULL,
  is_approved BOOLEAN DEFAULT false NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL,
  textsearchable_index_col tsvector GENERATED ALWAYS AS (
    to_tsvector('english',
      coalesce(name, '') 
      || ' ' || coalesce(wa_phone, '')
      || ' ' || coalesce(other_phone, '')
      || ' ' || coalesce(homestay_name, '')
      || ' ' || coalesce(homestay_address, '')
      || ' ' || coalesce(username, '')
    )
  ) STORED,
  textrank_index_col tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(name, '')), 'A')
    || setweight(to_tsvector('english', coalesce(wa_phone, '')), 'B')
    || setweight(to_tsvector('english', coalesce(homestay_name, '')), 'C')
    || setweight(to_tsvector('english', coalesce(homestay_address, '')), 'D')
  ) STORED
);
INSERT INTO membersx
  (id, name, wa_phone, other_phone, homestay_name, homestay_address, homestay_latitude, homestay_longitude, profile_pic_url, username, password, is_admin, is_approved, created_at, updated_at, deleted_at)
(
    SELECT id, name, wa_phone, other_phone, homestay_name, homestay_address, homestay_latitude, homestay_longitude, profile_pic_url, username, password, is_admin, is_approved, created_at, updated_at, deleted_at
    FROM members
);

ALTER TABLE member_dues DROP CONSTRAINT member_dues_member_id_fkey RESTRICT;
ALTER TABLE org_structures DROP CONSTRAINT org_structures_member_id_fkey RESTRICT;

DROP TABLE members;
DROP TABLE histories;
DROP TABLE goals;
DROP TABLE blogs;
DROP TABLE documents;

ALTER TABLE membersx RENAME TO members;
ALTER TABLE historiesx RENAME TO histories;
ALTER TABLE goalsx RENAME TO goals;
ALTER TABLE blogsx RENAME TO blogs;
ALTER TABLE documentsx RENAME TO documents;

ALTER TABLE member_dues ADD CONSTRAINT member_dues_member_id_fkey FOREIGN KEY (member_id) REFERENCES members(id);
ALTER TABLE org_structures ADD CONSTRAINT org_structures_member_id_fkey FOREIGN KEY (member_id) REFERENCES members(id);

CREATE INDEX members_textsearch_idx ON members USING GIN (textsearchable_index_col);
CREATE INDEX members_textrank_idx ON members USING GIN (textrank_index_col);

CREATE INDEX documents_textsearch_idx ON documents USING GIN (textsearchable_index_col);
CREATE INDEX documents_textrank_idx ON documents USING GIN (textrank_index_col);

CREATE INDEX blogs_textsearch_idx ON blogs USING GIN (textsearchable_index_col);
CREATE INDEX blogs_textrank_idx ON blogs USING GIN (textrank_index_col);

ALTER SEQUENCE blogsx_id_seq RENAME TO blogs_id_seq;
ALTER SEQUENCE historiesx_id_seq RENAME TO histories_id_seq;
ALTER SEQUENCE goalsx_id_seq RENAME TO goals_id_seq;
ALTER SEQUENCE documentsx_id_seq RENAME TO documents_id_seq;
