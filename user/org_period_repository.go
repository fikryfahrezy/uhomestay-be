package user

import (
	"context"
	"time"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type OrgPeriodRepository struct {
	PostgreDb *pgxpool.Pool
}

func NewOrgPeriodRepository(postgreDb *pgxpool.Pool) *OrgPeriodRepository {
	return &OrgPeriodRepository{
		PostgreDb: postgreDb,
	}
}

type (
	OrgPeriodExecutor   func(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	OrgPeriodQuerierRow func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	OrgPeriodQuerier    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
)

func (r *OrgPeriodRepository) Save(ctx context.Context, m OrgPeriodModel) (nm OrgPeriodModel, err error) {
	sqlQuery := `
		INSERT INTO org_periods (
			start_date,
			end_date,
			is_active,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var queryRow OrgPeriodQuerierRow
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		queryRow = tx.QueryRow
	} else {
		queryRow = r.PostgreDb.QueryRow
	}

	var lastInsertId uint64
	t := time.Now()

	err = queryRow(
		context.Background(),
		sqlQuery,
		m.StartDate,
		m.EndDate,
		m.IsActive,
		t,
		t,
		nil,
	).Scan(&lastInsertId)

	if err != nil {
		return OrgPeriodModel{}, err
	}

	m.Id = lastInsertId
	m.CreatedAt = t
	m.UpdatedAt = t

	return m, nil
}

func (r *OrgPeriodRepository) UpdateById(ctx context.Context, id uint64, m OrgPeriodModel) error {
	sqlQuery := `
		UPDATE org_periods SET (
			start_date,
			end_date,
			updated_at
		) = ($1, $2, $3)
		WHERE id = $4
	`

	var exec OrgPeriodExecutor
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
		m.StartDate,
		m.EndDate,
		t,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrgPeriodRepository) FindUndeletedById(ctx context.Context, id uint64) (m OrgPeriodModel, err error) {
	sqlQuery := `
		SELECT 
			id,
			start_date,
			end_date,
			is_active,
			created_at,
			updated_at,
			deleted_at
		FROM org_periods
		WHERE deleted_at IS NULL
		AND id = $1
	`

	var query OrgPeriodQuerier
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
		id,
	)

	if err != nil {
		return OrgPeriodModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return OrgPeriodModel{}, err
	}

	return m, nil
}

func (r *OrgPeriodRepository) FindActiveBydId(ctx context.Context, id uint64) (m OrgPeriodModel, err error) {
	sqlQuery := `
		SELECT 
			id,
			start_date,
			end_date,
			is_active,
			created_at,
			updated_at,
			deleted_at
		FROM org_periods
		WHERE deleted_at IS NULL
			AND is_active = true
			AND id = $1
	`

	var query OrgPeriodQuerier
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
		id,
	)

	if err != nil {
		return OrgPeriodModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return OrgPeriodModel{}, err
	}

	return m, nil
}

func (r *OrgPeriodRepository) UpdateStatusById(ctx context.Context, id uint64, m OrgPeriodModel) error {
	sqlQuery := `
		UPDATE org_periods SET is_active = $1 WHERE id = $2
	`

	var exec OrgPeriodExecutor
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
		m.IsActive,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrgPeriodRepository) DeleteById(ctx context.Context, id uint64) error {
	sqlQuery := `
		UPDATE org_periods 
		SET deleted_at = $1 
		WHERE id = $2
	`

	var exec OrgPeriodExecutor
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
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrgPeriodRepository) Query(ctx context.Context, id, limit int64) (ms []OrgPeriodModel, err error) {
	fromId := "id > $1"
	if id != 0 {
		fromId = "id < $1"
	}

	sqlQuery := `
		SELECT 
			id,
			start_date,
			end_date,
			is_active,
			created_at,
			updated_at,
			deleted_at
		FROM org_periods
		WHERE deleted_at IS NULL
			AND ` + fromId + `
		ORDER BY id DESC
		LIMIT $2
	`

	rows, _ := r.PostgreDb.Query(
		context.Background(),
		sqlQuery,
		id,
		limit,
	)

	defer rows.Close()

	var mps []*OrgPeriodModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []OrgPeriodModel{}, err
	}

	ms = make([]OrgPeriodModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *OrgPeriodRepository) DisableAll(ctx context.Context) error {
	sqlQuery := `
		UPDATE org_periods 
		SET is_active = $1
	`

	var exec OrgPeriodExecutor
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
		false,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrgPeriodRepository) EnableOtherLatest(ctx context.Context, id uint64) error {
	sqlQuery := `
		UPDATE org_periods SET is_active = true
		WHERE created_at = (
			SELECT MAX(created_at) FROM org_periods 
			WHERE deleted_at IS NULL AND id != $1
		)
	`

	var exec OrgPeriodExecutor
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
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrgPeriodRepository) QueryActive(ctx context.Context) (pm OrgPeriodModel, err error) {
	sqlQuery := `
		SELECT 
			id,
			start_date,
			end_date,
			is_active,
			created_at,
			updated_at,
			deleted_at
		FROM org_periods
		WHERE deleted_at IS NULL
		AND is_active = true
		LIMIT 1
	`

	rows, _ := r.PostgreDb.Query(
		context.Background(),
		sqlQuery,
	)

	if err != nil {
		return OrgPeriodModel{}, err
	}

	if err = pgxscan.ScanOne(&pm, rows); err != nil {
		return OrgPeriodModel{}, err
	}

	return pm, nil
}

func (r *OrgPeriodRepository) FindOtherLastTx(ctx context.Context) (m OrgPeriodModel, err error) {
	sqlQuery := `
		SELECT
			id,
			start_date,
			end_date,
			is_active,
			created_at,
			updated_at,
			deleted_at
		FROM org_periods
		WHERE deleted_at IS NULL
		ORDER BY id DESC
		LIMIT 1
	`

	var query OrgPeriodQuerier
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
	)

	if err != nil {
		return OrgPeriodModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return OrgPeriodModel{}, err
	}

	return m, nil
}

func (r *OrgPeriodRepository) FindById(ctx context.Context, id uint64) (m OrgPeriodModel, err error) {
	sqlQuery := `
		SELECT 
			id,
			start_date,
			end_date,
			is_active,
			created_at,
			updated_at,
			deleted_at
		FROM org_periods
		WHERE id = $1
	`

	var query OrgPeriodQuerier
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
		id,
	)

	if err != nil {
		return OrgPeriodModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return OrgPeriodModel{}, err
	}

	return m, nil
}
