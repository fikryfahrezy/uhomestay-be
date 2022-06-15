package dues

import (
	"context"
	"time"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DuesRepository struct {
	PostgreDb *pgxpool.Pool
}

func NewDeusRepository(postgreDb *pgxpool.Pool) *DuesRepository {
	return &DuesRepository{
		PostgreDb: postgreDb,
	}
}

type (
	DuesExecutor   func(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	DuesQuerierRow func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	DuesQuerier    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
)

func (r *DuesRepository) FindOtherByYYYYMM(ctx context.Context, id uint64, date time.Time) (m DuesModel, err error) {
	// Ref: YYYY-MM column type in PostgreSQL
	// https://stackoverflow.com/a/43657553/12976234
	querystr := `
		SELECT
			id,
			date,
			idr_amount,
			created_at,
			updated_at,
			deleted_at
		FROM dues 
		WHERE deleted_at IS NULL
			AND date = date_trunc('month', $1::timestamp)
			AND id != $2
	`

	var query DuesQuerier
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		query = tx.Query
	} else {
		query = r.PostgreDb.Query
	}

	var rows pgx.Rows
	rows, err = query(
		context.Background(),
		querystr,
		date.Format(time.RFC3339),
		id,
	)

	if err != nil {
		return DuesModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return DuesModel{}, err
	}

	return m, nil
}

func (r *DuesRepository) Save(ctx context.Context, m DuesModel) (nm DuesModel, err error) {
	sqlQuery := `
		INSERT INTO dues (
			date,
			idr_amount,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES (date_trunc('month', $1::timestamp), $2, $3, $4, $5)
		RETURNING id
	`

	var queryRow DuesQuerierRow
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
		m.Date,
		m.IdrAmount,
		t,
		t,
		nil,
	).Scan(&lastInsertId)

	if err != nil {
		return DuesModel{}, err
	}

	m.Id = lastInsertId
	m.CreatedAt = t
	m.UpdatedAt = t

	return m, nil
}

func (r *DuesRepository) UpdateById(ctx context.Context, id uint64, m DuesModel) error {
	sqlQuery := `
		UPDATE dues SET (
			date,
			idr_amount,
			updated_at
		) = (date_trunc('month', $1::timestamp), $2, $3)
		WHERE id = $4
	`

	var exec DuesExecutor
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
		m.Date,
		m.IdrAmount,
		t,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DuesRepository) FindById(ctx context.Context, id uint64) (m DuesModel, err error) {
	querystr := `
		SELECT
			id,
			date,
			idr_amount,
			created_at,
			updated_at,
			deleted_at
		FROM dues 
		WHERE deleted_at IS NULL
		AND id = $1
	`

	var query DuesQuerier
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		query = tx.Query
	} else {
		query = r.PostgreDb.Query
	}

	var rows pgx.Rows
	rows, err = query(
		context.Background(),
		querystr,
		id,
	)

	if err != nil {
		return DuesModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return DuesModel{}, err
	}

	return m, nil
}

func (r *DuesRepository) DeleteById(ctx context.Context, id uint64) error {
	sqlQuery := `
		UPDATE dues
		SET deleted_at = $1
		WHERE id = $2
	`

	var exec DuesExecutor
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

func (r *DuesRepository) Query(ctx context.Context, id, limit int64) ([]DuesModel, error) {
	fromId := "id > $1"
	if id != 0 {
		fromId = "id < $1"
	}

	sqlQuery := `
		SELECT 
			id,
			date,
			idr_amount,
			created_at,
			updated_at,
			deleted_at
		FROM dues 
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

	var mps []*DuesModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []DuesModel{}, err
	}

	ms := make([]DuesModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *DuesRepository) QueryAmtByUidStatus(ctx context.Context, uid string, status DuesStatus) ([]string, error) {
	sqlQuery := `
		SELECT 
			dues.idr_amount
		FROM dues 
		RIGHT JOIN member_dues md ON md.dues_id = dues.id
		WHERE dues.deleted_at IS NULL
			AND md.member_id = $1
			AND md.deleted_at IS NULL
			AND md.status = $2
	`

	rows, _ := r.PostgreDb.Query(
		context.Background(),
		sqlQuery,
		uid,
		status.String,
	)
	defer rows.Close()

	var mps []*string
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []string{}, err
	}

	ms := make([]string, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *DuesRepository) Latest(ctx context.Context) (m DuesModel, err error) {
	querystr := `
		SELECT
			id,
			date,
			idr_amount,
			created_at,
			updated_at,
			deleted_at
		FROM dues 
		WHERE deleted_at IS NULL
		ORDER BY id DESC
		LIMIT 1
	`

	var query DuesQuerier
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		query = tx.Query
	} else {
		query = r.PostgreDb.Query
	}

	var rows pgx.Rows
	rows, err = query(
		context.Background(),
		querystr,
	)

	if err != nil {
		return DuesModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return DuesModel{}, err
	}

	return m, nil
}
