package user

import (
	"context"
	"time"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OrgStructureRepository struct {
	PostgreDb *pgxpool.Pool
}

func NewOrgStructureRepository(postgreDb *pgxpool.Pool) *OrgStructureRepository {
	return &OrgStructureRepository{
		PostgreDb: postgreDb,
	}
}

type (
	OrgStructureExecutor   func(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	OrgStructureQuerierRow func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	OrgStructureQuerier    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	OrgStructureCopierFrom func(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
)

func (r *OrgStructureRepository) Save(ctx context.Context, m OrgStructureModel) error {
	sqlQuery := `
		INSERT INTO org_structures (
			position_name,
			position_level,
			member_id,
			position_id,
			org_period_id,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	var exec OrgStructureExecutor
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		exec = tx.Exec
	} else {
		exec = r.PostgreDb.Exec
	}

	var err error
	t := time.Now()

	_, err = exec(
		context.Background(),
		sqlQuery,
		m.PositionName,
		m.PositionLevel,
		m.MemberId,
		m.PositionId,
		m.OrgPeriodId,
		t,
		t,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrgStructureRepository) FindLatestByMemberId(ctx context.Context, uid string) (m OrgStructureModel, err error) {
	sqlQuery := `
		SELECT
			id,
			position_name,
			position_level,
			member_id,
			position_id,
			org_period_id,
			created_at,
			updated_at,
			deleted_at
		FROM org_structures
		WHERE member_id = $1
		ORDER BY id DESC
		LIMIT 1
	`

	var query OrgStructureQuerier
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		query = tx.Query
	} else {
		query = r.PostgreDb.Query
	}

	var rows pgx.Rows

	rows, err = query(
		context.Background(),
		sqlQuery,
		uid,
	)

	if err != nil {
		return OrgStructureModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return OrgStructureModel{}, err
	}

	return m, nil
}

func (r *OrgStructureRepository) BulkSave(ctx context.Context, ms []OrgStructureModel) error {
	var copyFrom OrgStructureCopierFrom
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		copyFrom = tx.CopyFrom
	} else {
		copyFrom = r.PostgreDb.CopyFrom
	}
	var err error
	t := time.Now()

	_, err = copyFrom(
		context.Background(),
		pgx.Identifier{"org_structures"},
		[]string{"position_name", "position_level", "member_id", "position_id", "org_period_id", "created_at", "updated_at", "deleted_at"},
		pgx.CopyFromSlice(len(ms), func(i int) ([]interface{}, error) {
			memberId := pgtypeuuid.UUID{
				Status: pgtype.Null,
			}

			if ms[i].MemberId != "" {
				err := memberId.Scan(ms[i].MemberId)
				if err != nil {
					return []interface{}{}, err
				}
			}

			return []interface{}{ms[i].PositionName, ms[i].PositionLevel, memberId, ms[i].PositionId, ms[i].OrgPeriodId, t, t, nil}, nil
		}),
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrgStructureRepository) FindByPeriodId(ctx context.Context, periodId uint64) (ms []OrgStructureModel, err error) {
	sqlQuery := `
		SELECT
			id,
			position_name,
			position_level,
			member_id,
			position_id,
			org_period_id,
			created_at,
			updated_at,
			deleted_at
		FROM org_structures
		WHERE org_period_id = $1
	`

	var query OrgStructureQuerier
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		query = tx.Query
	} else {
		query = r.PostgreDb.Query
	}

	var rows pgx.Rows

	rows, _ = query(
		context.Background(),
		sqlQuery,
		periodId,
	)

	defer rows.Close()

	var mps []*OrgStructureModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []OrgStructureModel{}, err
	}

	ms = make([]OrgStructureModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *OrgStructureRepository) DeleteByPeriodId(ctx context.Context, periodId uint64) error {
	sqlQuery := `
		UPDATE org_structures 
		SET (updated_at, deleted_at) = ($1, $2)
		WHERE org_period_id = $3
	`

	var exec OrgStructureExecutor
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		exec = tx.Exec
	} else {
		exec = r.PostgreDb.Exec
	}

	var err error
	t := time.Now()

	_, err = exec(
		context.Background(),
		sqlQuery,
		t,
		t,
		periodId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrgStructureRepository) FindByOrgIdAndMemberId(ctx context.Context, orgId uint64, uid string) (m []OrgStructureModel, err error) {
	sqlQuery := `
		SELECT
			id,
			position_name,
			position_level,
			member_id,
			position_id,
			org_period_id,
			created_at,
			updated_at,
			deleted_at
		FROM org_structures
		WHERE org_period_id = $1
			AND member_id = $2
		ORDER BY id DESC
	`

	var query OrgStructureQuerier
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		query = tx.Query
	} else {
		query = r.PostgreDb.Query
	}

	rows, _ := query(
		context.Background(),
		sqlQuery,
		orgId,
		uid,
	)

	var mps []*OrgStructureModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []OrgStructureModel{}, err
	}

	ms := make([]OrgStructureModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *OrgStructureRepository) DeleteByOrgIdAndMemberId(ctx context.Context, orgId uint64, uid string) error {
	sqlQuery := `
	DELETE FROM org_structures
	WHERE org_period_id = $1 AND member_id = $2
`

	var exec OrgStructureExecutor
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		exec = tx.Exec
	} else {
		exec = r.PostgreDb.Exec
	}

	var err error
	_, err = exec(
		context.Background(),
		sqlQuery,
		orgId,
		uid,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrgStructureRepository) FindByPosIdAndActivePeriod(ctx context.Context, posId uint64) (m OrgStructureModel, err error) {
	sqlQuery := `
	SELECT
		orgs.id,
		position_name,
		position_level,
		member_id,
		position_id,
		org_period_id,
		orgs.created_at,
		orgs.updated_at,
		orgs.deleted_at
	FROM org_structures orgs
	LEFT JOIN org_periods op on op.id = orgs.org_period_id 
	WHERE orgs.position_id = $1 AND op.is_active = true AND op.deleted_at IS NULL
	ORDER BY orgs.id DESC
	LIMIT 1
	`

	var query OrgStructureQuerier
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		query = tx.Query
	} else {
		query = r.PostgreDb.Query
	}

	var rows pgx.Rows

	rows, err = query(
		context.Background(),
		sqlQuery,
		posId,
	)

	if err != nil {
		return OrgStructureModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return OrgStructureModel{}, err
	}

	return m, nil
}

func (r *OrgStructureRepository) UpdatePosByPosIdAndOrgId(ctx context.Context, posId, periodId uint64, pos PositionModel) error {
	sqlQuery := `
		UPDATE org_structures
		SET (position_name, position_level, updated_at) = ($1, $2, $3)
		WHERE position_id = $4 AND org_period_id = $5
	`

	var exec OrgStructureExecutor
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		exec = tx.Exec
	} else {
		exec = r.PostgreDb.Exec
	}

	var err error

	_, err = exec(
		context.Background(),
		sqlQuery,
		pos.Name,
		pos.Level,
		time.Now(),
		posId,
		periodId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrgStructureRepository) UndeleteByPeriodId(ctx context.Context, periodId uint64) error {
	sqlQuery := `
		UPDATE org_structures 
		SET (updated_at, deleted_at) = ($1, NULL)
		WHERE org_period_id = $2
	`

	var exec OrgStructureExecutor
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		exec = tx.Exec
	} else {
		exec = r.PostgreDb.Exec
	}

	var err error
	_, err = exec(
		context.Background(),
		sqlQuery,
		time.Now(),
		periodId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrgStructureRepository) FindLatestUndeletedByMemberId(ctx context.Context, uid string) (m OrgStructureModel, err error) {
	sqlQuery := `
		SELECT
			id,
			position_name,
			position_level,
			member_id,
			position_id,
			org_period_id,
			created_at,
			updated_at,
			deleted_at
		FROM org_structures
		WHERE member_id = $1
		ORDER BY updated_at DESC
		LIMIT 1
	`

	var query OrgStructureQuerier
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		query = tx.Query
	} else {
		query = r.PostgreDb.Query
	}

	var rows pgx.Rows

	rows, err = query(
		context.Background(),
		sqlQuery,
		uid,
	)

	if err != nil {
		return OrgStructureModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return OrgStructureModel{}, err
	}

	return m, nil
}
