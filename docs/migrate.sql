CREATE TABLE IF NOT EXISTS members_x (
  id UUID PRIMARY KEY,
  name VARCHAR(100) DEFAULT '' NOT NULL,
  other_phone VARCHAR(50) DEFAULT '' NOT NULL UNIQUE,
  wa_phone VARCHAR(50) DEFAULT '' NOT NULL UNIQUE,
  profile_pic_url TEXT DEFAULT '' NOT NULL,
  id_card_url TEXT DEFAULT '' NOT NULL,
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
      || ' ' || coalesce(other_phone, '')
      || ' ' || coalesce(wa_phone, '')
      || ' ' || coalesce(username, '')
    )
  ) STORED,
  textrank_index_col tsvector GENERATED ALWAYS AS (
    setweight(to_tsvector('english', coalesce(name, '')), 'A')
    || setweight(to_tsvector('english', coalesce(wa_phone, '')), 'B')
    || setweight(to_tsvector('english', coalesce(other_phone, '')), 'C')
  ) STORED
);

INSERT INTO members_x
  (id, name, other_phone, wa_phone, profile_pic_url, id_card_url, username, password, is_admin, is_approved, created_at, updated_at, deleted_at)
(
    SELECT id, name, other_phone, wa_phone, profile_pic_url, profile_pic_url, username, password, is_admin, is_approved, created_at, updated_at, deleted_at
    FROM members
);

ALTER TABLE members RENAME TO members_xx;
ALTER TABLE members_x RENAME TO members;

ALTER TABLE org_structures DROP CONSTRAINT org_structures_member_id_fkey1;
ALTER TABLE org_structures ADD CONSTRAINT org_structures_member_id_fkey1 FOREIGN KEY (member_id) REFERENCES members(id);

ALTER TABLE member_dues DROP CONSTRAINT member_dues_x_member_id_fkey;
ALTER TABLE member_dues ADD CONSTRAINT member_dues_x_member_id_fkey FOREIGN KEY (member_id) REFERENCES members(id);

CREATE TABLE IF NOT EXISTS member_homestays (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(200) DEFAULT '' NOT NULL,
  address VARCHAR(200) DEFAULT '' NOT NULL,
  latitude VARCHAR(50) DEFAULT '' NOT NULL,
  longitude VARCHAR(50) DEFAULT '' NOT NULL,
  thumbnail_url TEXT DEFAULT '' NOT NULL,
  member_id UUID NOT NULL REFERENCES members(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL
);

INSERT INTO member_homestays (name, address, latitude, longitude, member_id)
(
	SELECT homestay_name, homestay_address, homestay_latitude, homestay_longitude, id
	FROM members_xx
);

CREATE TABLE IF NOT EXISTS homestay_images (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(200) DEFAULT '' NOT NULL,
  alphnum_name VARCHAR(200) DEFAULT '' NOT NULL,
  url TEXT DEFAULT '' NOT NULL,
  member_homestay_id BIGINT DEFAULT NULL REFERENCES member_homestays(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL
);

DROP TABLE members_xx;

-- CREATE TABLE IF NOT EXISTS articles (
--   id BIGSERIAL PRIMARY KEY,
--   title VARCHAR(200) DEFAULT '' NOT NULL,
--   short_desc VARCHAR(200) DEFAULT '' NOT NULL,
--   thumbnail_url TEXT DEFAULT '' NOT NULL,
--   content jsonb DEFAULT '{}'::jsonb NOT NULL,
--   content_text TEXT DEFAULT '' NOT NULL,
--   slug VARCHAR(200) DEFAULT '' NOT NULL,
--   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
--   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
--   deleted_at TIMESTAMP DEFAULT NULL,
--   textsearchable_index_col tsvector GENERATED ALWAYS AS (
--     to_tsvector('english',
--       coalesce(title, '') 
--       || ' ' || coalesce(short_desc, '')
--       || ' ' || coalesce(content_text, '')
--       || ' ' || coalesce(slug, '')
--     )
--   ) STORED,
--   textrank_index_col tsvector GENERATED ALWAYS AS (
--     setweight(to_tsvector('english', coalesce(title, '')), 'A')
--     || setweight(to_tsvector('english', coalesce(short_desc, '')), 'B')
--     || setweight(to_tsvector('english', coalesce(content_text, '')), 'C')
--     || setweight(to_tsvector('english', coalesce(slug, '')), 'D')
--   ) STORED
-- );

-- CREATE INDEX articles_textsearch_idx ON articles USING GIN (textsearchable_index_col);

-- CREATE INDEX articles_textrank_idx ON articles USING GIN (textrank_index_col);

-- INSERT INTO articles (title, short_desc, thumbnail_url, content, content_text, slug, created_at, updated_at, deleted_at) 
-- (
-- 	SELECT title, short_desc, thumbnail_url, content, content_text, slug, created_at, updated_at, deleted_at FROM blogs
-- );