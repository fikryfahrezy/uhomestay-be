CREATE TYPE public.cashflowtype AS ENUM (
    'income',
    'outcome'
);

CREATE TYPE public.doctype AS ENUM (
    'dir',
    'file'
);

CREATE TYPE public.duesstatus AS ENUM (
    'unpaid',
    'waiting',
    'paid'
);

CREATE TABLE public.articles (
    id bigint NOT NULL,
    title character varying(200) DEFAULT ''::character varying NOT NULL,
    short_desc character varying(200) DEFAULT ''::character varying NOT NULL,
    thumbnail_url text DEFAULT ''::text NOT NULL,
    content jsonb DEFAULT '{}'::jsonb NOT NULL,
    content_text text DEFAULT ''::text NOT NULL,
    slug character varying(200) DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    textsearchable_index_col tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, (((((((COALESCE(title, ''::character varying))::text || ' '::text) || (COALESCE(short_desc, ''::character varying))::text) || ' '::text) || COALESCE(content_text, ''::text)) || ' '::text) || (COALESCE(slug, ''::character varying))::text))) STORED,
    textrank_index_col tsvector GENERATED ALWAYS AS ((((setweight(to_tsvector('english'::regconfig, (COALESCE(title, ''::character varying))::text), 'A'::"char") || setweight(to_tsvector('english'::regconfig, (COALESCE(short_desc, ''::character varying))::text), 'B'::"char")) || setweight(to_tsvector('english'::regconfig, COALESCE(content_text, ''::text)), 'C'::"char")) || setweight(to_tsvector('english'::regconfig, (COALESCE(slug, ''::character varying))::text), 'D'::"char"))) STORED
);

CREATE SEQUENCE public.articles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.articles_id_seq OWNED BY public.articles.id;

CREATE TABLE public.cashflows (
    id bigint NOT NULL,
    date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    idr_amount character varying(200) DEFAULT ''::character varying NOT NULL,
    type public.cashflowtype DEFAULT 'income'::public.cashflowtype NOT NULL,
    note text DEFAULT ''::text NOT NULL,
    prove_file_url text DEFAULT ''::text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);

CREATE SEQUENCE public.cashflows_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.cashflows_id_seq OWNED BY public.cashflows.id;

CREATE TABLE public.documents (
    id bigint NOT NULL,
    name character varying(200) DEFAULT ''::character varying NOT NULL,
    alphnum_name character varying(200) DEFAULT ''::character varying NOT NULL,
    url text DEFAULT ''::text NOT NULL,
    type public.doctype DEFAULT 'file'::public.doctype NOT NULL,
    dir_id bigint DEFAULT 0 NOT NULL,
    is_private boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    textsearchable_index_col tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, (((((COALESCE(alphnum_name, ''::character varying))::text || ' '::text) || (COALESCE(name, ''::character varying))::text) || ' '::text) || COALESCE(url, ''::text)))) STORED,
    textrank_index_col tsvector GENERATED ALWAYS AS (((setweight(to_tsvector('english'::regconfig, (COALESCE(alphnum_name, ''::character varying))::text), 'A'::"char") || setweight(to_tsvector('english'::regconfig, (COALESCE(name, ''::character varying))::text), 'B'::"char")) || setweight(to_tsvector('english'::regconfig, COALESCE(url, ''::text)), 'C'::"char"))) STORED
);

CREATE SEQUENCE public.documents_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.documents_id_seq OWNED BY public.documents.id;

CREATE TABLE public.dues (
    id bigint NOT NULL,
    date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    idr_amount character varying(200) DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);

CREATE SEQUENCE public.dues_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.dues_id_seq OWNED BY public.dues.id;

CREATE TABLE public.goals (
    id bigint NOT NULL,
    vision jsonb DEFAULT '{}'::jsonb NOT NULL,
    vision_text text DEFAULT ''::text NOT NULL,
    mission jsonb DEFAULT '{}'::jsonb NOT NULL,
    mission_text text DEFAULT ''::text NOT NULL,
    org_period_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE SEQUENCE public.goals_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.goals_id_seq OWNED BY public.goals.id;

CREATE TABLE public.histories (
    id bigint NOT NULL,
    content jsonb DEFAULT '{}'::jsonb NOT NULL,
    content_text text DEFAULT ''::text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE SEQUENCE public.histories_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.histories_id_seq OWNED BY public.histories.id;

CREATE TABLE public.homestay_images (
    id bigint NOT NULL,
    name character varying(200) DEFAULT ''::character varying NOT NULL,
    alphnum_name character varying(200) DEFAULT ''::character varying NOT NULL,
    url text DEFAULT ''::text NOT NULL,
    member_homestay_id bigint,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);

CREATE SEQUENCE public.homestay_images_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.homestay_images_id_seq OWNED BY public.homestay_images.id;

CREATE TABLE public.image_caches (
    name character varying DEFAULT ''::character varying NOT NULL,
    image_id character varying DEFAULT ''::character varying NOT NULL,
    image_url character varying DEFAULT ''::character varying NOT NULL
);

CREATE TABLE public.images (
    id bigint NOT NULL,
    name character varying(200) DEFAULT ''::character varying NOT NULL,
    alphnum_name character varying(200) DEFAULT ''::character varying NOT NULL,
    url text DEFAULT ''::text NOT NULL,
    description text DEFAULT ''::text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);

CREATE SEQUENCE public.images_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.images_id_seq OWNED BY public.images.id;

CREATE TABLE public.member_dues (
    id bigint NOT NULL,
    member_id uuid NOT NULL,
    dues_id bigint NOT NULL,
    status public.duesstatus DEFAULT 'unpaid'::public.duesstatus NOT NULL,
    prove_file_url text DEFAULT ''::text NOT NULL,
    pay_date timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);

CREATE SEQUENCE public.member_dues_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.member_dues_id_seq OWNED BY public.member_dues.id;

CREATE TABLE public.member_homestays (
    id bigint NOT NULL,
    name character varying(200) DEFAULT ''::character varying NOT NULL,
    address character varying(200) DEFAULT ''::character varying NOT NULL,
    latitude character varying(50) DEFAULT ''::character varying NOT NULL,
    longitude character varying(50) DEFAULT ''::character varying NOT NULL,
    thumbnail_url text DEFAULT ''::text NOT NULL,
    member_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);

CREATE SEQUENCE public.member_homestays_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.member_homestays_id_seq OWNED BY public.member_homestays.id;

CREATE TABLE public.members (
    id uuid NOT NULL,
    name character varying(100) DEFAULT ''::character varying NOT NULL,
    other_phone character varying(50) DEFAULT ''::character varying NOT NULL,
    wa_phone character varying(50) DEFAULT ''::character varying NOT NULL,
    profile_pic_url text DEFAULT ''::text NOT NULL,
    id_card_url text DEFAULT ''::text NOT NULL,
    username character varying(50) DEFAULT ''::character varying NOT NULL,
    password character varying(200) DEFAULT ''::character varying NOT NULL,
    is_admin boolean DEFAULT false NOT NULL,
    is_approved boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone,
    textsearchable_index_col tsvector GENERATED ALWAYS AS (to_tsvector('english'::regconfig, (((((((COALESCE(name, ''::character varying))::text || ' '::text) || (COALESCE(other_phone, ''::character varying))::text) || ' '::text) || (COALESCE(wa_phone, ''::character varying))::text) || ' '::text) || (COALESCE(username, ''::character varying))::text))) STORED,
    textrank_index_col tsvector GENERATED ALWAYS AS (((setweight(to_tsvector('english'::regconfig, (COALESCE(name, ''::character varying))::text), 'A'::"char") || setweight(to_tsvector('english'::regconfig, (COALESCE(wa_phone, ''::character varying))::text), 'B'::"char")) || setweight(to_tsvector('english'::regconfig, (COALESCE(other_phone, ''::character varying))::text), 'C'::"char"))) STORED
);

CREATE TABLE public.org_periods (
    id bigint NOT NULL,
    start_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    end_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);

CREATE SEQUENCE public.org_periods_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.org_periods_id_seq OWNED BY public.org_periods.id;

CREATE TABLE public.org_structures (
    id bigint NOT NULL,
    position_name character varying(200) DEFAULT ''::character varying NOT NULL,
    position_level smallint DEFAULT 0 NOT NULL,
    member_id uuid NOT NULL,
    position_id bigint NOT NULL,
    org_period_id bigint NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);

CREATE SEQUENCE public.org_structures_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.org_structures_id_seq OWNED BY public.org_structures.id;

CREATE TABLE public.positions (
    id bigint NOT NULL,
    name character varying(200) DEFAULT ''::character varying NOT NULL,
    level smallint DEFAULT 0 NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp without time zone
);

CREATE SEQUENCE public.positions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.positions_id_seq OWNED BY public.positions.id;

ALTER TABLE ONLY public.articles ALTER COLUMN id SET DEFAULT nextval('public.articles_id_seq'::regclass);

ALTER TABLE ONLY public.cashflows ALTER COLUMN id SET DEFAULT nextval('public.cashflows_id_seq'::regclass);

ALTER TABLE ONLY public.documents ALTER COLUMN id SET DEFAULT nextval('public.documents_id_seq'::regclass);

ALTER TABLE ONLY public.dues ALTER COLUMN id SET DEFAULT nextval('public.dues_id_seq'::regclass);

ALTER TABLE ONLY public.goals ALTER COLUMN id SET DEFAULT nextval('public.goals_id_seq'::regclass);

ALTER TABLE ONLY public.histories ALTER COLUMN id SET DEFAULT nextval('public.histories_id_seq'::regclass);

ALTER TABLE ONLY public.homestay_images ALTER COLUMN id SET DEFAULT nextval('public.homestay_images_id_seq'::regclass);

ALTER TABLE ONLY public.images ALTER COLUMN id SET DEFAULT nextval('public.images_id_seq'::regclass);

ALTER TABLE ONLY public.member_dues ALTER COLUMN id SET DEFAULT nextval('public.member_dues_id_seq'::regclass);

ALTER TABLE ONLY public.member_homestays ALTER COLUMN id SET DEFAULT nextval('public.member_homestays_id_seq'::regclass);

ALTER TABLE ONLY public.org_periods ALTER COLUMN id SET DEFAULT nextval('public.org_periods_id_seq'::regclass);

ALTER TABLE ONLY public.org_structures ALTER COLUMN id SET DEFAULT nextval('public.org_structures_id_seq'::regclass);

ALTER TABLE ONLY public.positions ALTER COLUMN id SET DEFAULT nextval('public.positions_id_seq'::regclass);

ALTER TABLE ONLY public.articles
    ADD CONSTRAINT articles_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.cashflows
    ADD CONSTRAINT cashflows_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.documents
    ADD CONSTRAINT documents_x_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.dues
    ADD CONSTRAINT dues_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.goals
    ADD CONSTRAINT goals_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.histories
    ADD CONSTRAINT histories_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.homestay_images
    ADD CONSTRAINT homestay_images_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.images
    ADD CONSTRAINT images_x_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.member_dues
    ADD CONSTRAINT member_dues_x_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.member_homestays
    ADD CONSTRAINT member_homestays_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.members
    ADD CONSTRAINT members_x_other_phone_key UNIQUE (other_phone);

ALTER TABLE ONLY public.members
    ADD CONSTRAINT members_x_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.members
    ADD CONSTRAINT members_x_username_key UNIQUE (username);

ALTER TABLE ONLY public.members
    ADD CONSTRAINT members_x_wa_phone_key UNIQUE (wa_phone);

ALTER TABLE ONLY public.org_periods
    ADD CONSTRAINT org_periods_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.org_structures
    ADD CONSTRAINT org_structures_x_pkey1 PRIMARY KEY (id);

ALTER TABLE ONLY public.positions
    ADD CONSTRAINT positions_pkey PRIMARY KEY (id);

CREATE INDEX articles_textrank_idx ON public.articles USING gin (textrank_index_col);

CREATE INDEX articles_textsearch_idx ON public.articles USING gin (textsearchable_index_col);

ALTER TABLE ONLY public.goals
    ADD CONSTRAINT goals_org_period_id_fkey FOREIGN KEY (org_period_id) REFERENCES public.org_periods(id);

ALTER TABLE ONLY public.homestay_images
    ADD CONSTRAINT homestay_images_member_homestay_id_fkey FOREIGN KEY (member_homestay_id) REFERENCES public.member_homestays(id);

ALTER TABLE ONLY public.member_dues
    ADD CONSTRAINT member_dues_x_dues_id_fkey FOREIGN KEY (dues_id) REFERENCES public.dues(id);

ALTER TABLE ONLY public.member_dues
    ADD CONSTRAINT member_dues_x_member_id_fkey FOREIGN KEY (member_id) REFERENCES public.members(id);

ALTER TABLE ONLY public.member_homestays
    ADD CONSTRAINT member_homestays_member_id_fkey FOREIGN KEY (member_id) REFERENCES public.members(id);

ALTER TABLE ONLY public.org_structures
    ADD CONSTRAINT org_structures_x_member_id_fkey1 FOREIGN KEY (member_id) REFERENCES public.members(id);

ALTER TABLE ONLY public.org_structures
    ADD CONSTRAINT org_structures_x_org_period_id_fkey1 FOREIGN KEY (org_period_id) REFERENCES public.org_periods(id);

ALTER TABLE ONLY public.org_structures
    ADD CONSTRAINT org_structures_x_position_id_fkey1 FOREIGN KEY (position_id) REFERENCES public.positions(id);

