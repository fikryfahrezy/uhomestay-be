CREATE TABLE IF NOT EXISTS org_structures_x (
  id BIGSERIAL PRIMARY KEY,
  position_name VARCHAR(200) DEFAULT '' NOT NULL,
  position_level SMALLINT DEFAULT 0 NOT NULL,
  member_id UUID NOT NULL REFERENCES members(id),
  position_id BIGINT NOT NULL REFERENCES positions(id),
  org_period_id BIGINT NOT NULL REFERENCES org_periods(id),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP DEFAULT NULL
);

INSERT INTO org_structures_x
  (position_name, position_level, member_id, position_id, org_period_id, created_at, updated_at, deleted_at)
(
    SELECT position_name, position_level, member_id, position_id, org_period_id, created_at, created_at, deleted_at
    FROM org_structures ogs
);

DROP TABLE org_structures;

ALTER TABLE org_structures_x RENAME TO org_structures;
ALTER SEQUENCE org_structures_x_id_seq RENAME TO org_structures_id_seq;
