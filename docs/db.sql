CREATE TABLE IF NOT EXISTS members (
  id UUID PRIMARY KEY,
  name VARCHAR(100) DEFAULT '' NOT NULL,
  other_phone VARCHAR(50) DEFAULT '' NOT NULL UNIQUE,
  wa_phone VARCHAR(50) DEFAULT '' NOT NULL UNIQUE,
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
      || ' ' || coalesce(other_phone, '')
      || ' ' || coalesce(wa_phone, '')
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

CREATE INDEX members_textsearch_idx ON members USING GIN (textsearchable_index_col);

CREATE INDEX members_textrank_idx ON members USING GIN (textrank_index_col);

CREATE TABLE IF NOT EXISTS positions (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(200) DEFAULT '' NOT NULL,
  level SMALLINT DEFAULT 0 NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS org_periods (
  id BIGSERIAL PRIMARY KEY,
  start_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  end_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  is_active BOOLEAN DEFAULT true NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS org_structures (
  id BIGSERIAL PRIMARY KEY,
  position_name VARCHAR(200) DEFAULT '' NOT NULL,
  position_level SMALLINT DEFAULT 0 NOT NULL,
  member_id UUID NOT NULL REFERENCES members(id),
  position_id BIGINT NOT NULL REFERENCES positions(id),
  org_period_id BIGINT NOT NULL REFERENCES org_periods(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TYPE doctype AS ENUM ('dir', 'file');

CREATE TABLE IF NOT EXISTS documents (
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

CREATE INDEX documents_textsearch_idx ON documents USING GIN (textsearchable_index_col);

CREATE INDEX documents_textrank_idx ON documents USING GIN (textrank_index_col);

CREATE TABLE IF NOT EXISTS blogs (
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

CREATE INDEX blogs_textsearch_idx ON blogs USING GIN (textsearchable_index_col);

CREATE INDEX blogs_textrank_idx ON blogs USING GIN (textrank_index_col);

CREATE TABLE IF NOT EXISTS histories (
  id BIGSERIAL PRIMARY KEY,
  content jsonb DEFAULT '{}'::jsonb NOT NULL,
  content_text TEXT DEFAULT '' NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS goals (
  id BIGSERIAL PRIMARY KEY,
  vision jsonb DEFAULT '{}'::jsonb NOT NULL,
  vision_text TEXT DEFAULT '' NOT NULL,
  mission jsonb DEFAULT '{}'::jsonb NOT NULL,
  mission_text TEXT DEFAULT '' NOT NULL,
  org_period_id BIGINT NOT NULL REFERENCES org_periods(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TYPE cashflowtype AS ENUM ('income', 'outcome'); 

CREATE TABLE IF NOT EXISTS cashflows (
  id BIGSERIAL PRIMARY KEY,
  date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  idr_amount VARCHAR(200) DEFAULT '' NOT NULL,
  type cashflowtype DEFAULT 'income' NOT NULL,
  note TEXT DEFAULT '' NOT NULL,
  prove_file_url TEXT DEFAULT '' NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL, 
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS dues (
  id BIGSERIAL PRIMARY KEY,
  date TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  idr_amount VARCHAR(200) DEFAULT '' NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL, 
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TYPE duesstatus AS ENUM ('unpaid', 'waiting', 'paid');

CREATE TABLE IF NOT EXISTS member_dues (
  id BIGSERIAL PRIMARY KEY,
  member_id UUID NOT NULL REFERENCES members(id),
  dues_id BIGINT NOT NULL REFERENCES dues(id),
  status duesstatus DEFAULT 'unpaid' NOT NULL,
  prove_file_url TEXT DEFAULT '' NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL, 
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL
);
