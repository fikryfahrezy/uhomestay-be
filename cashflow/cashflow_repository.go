package cashflow

import (
	"context"
	"time"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type CashflowRepository struct {
	PostgreDb *pgxpool.Pool
}

func NewRepository(postgreDb *pgxpool.Pool) *CashflowRepository {
	return &CashflowRepository{
		PostgreDb: postgreDb,
	}
}

type (
	CasflowExecutor    func(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	CashflowQuerierRow func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	CashflowQuerier    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
)

func (r *CashflowRepository) Save(ctx context.Context, m CashflowModel) (nm CashflowModel, err error) {
	sqlQuery := `
		INSERT INTO cashflows (
			date,
			idr_amount,
			type,
			note,
			prove_file_url,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var queryRow CashflowQuerierRow
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
		m.Type,
		m.Note,
		m.ProveFileUrl,
		t,
		t,
		nil,
	).Scan(&lastInsertId)

	if err != nil {
		return CashflowModel{}, err
	}

	m.Id = lastInsertId
	m.CreatedAt = t
	m.UpdatedAt = t

	return m, nil
}

func (r *CashflowRepository) UpdateById(ctx context.Context, id uint64, m CashflowModel) error {
	sqlQuery := `
		UPDATE cashflows SET (
			date,
			idr_amount,
			type,
			note,
			prove_file_url,
			updated_at
		) = ($1, $2, $3, $4, $5, $6)
		WHERE id = $7
	`

	var exec CasflowExecutor
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
		m.Type,
		m.Note,
		m.ProveFileUrl,
		t,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *CashflowRepository) FindById(ctx context.Context, id uint64) (m CashflowModel, err error) {
	querystr := `
		SELECT
			id,
			date,
			idr_amount,
			type,
			note,
			prove_file_url,
			created_at,
			updated_at,
			deleted_at
		FROM cashflows
		WHERE deleted_at IS NULL
		AND id = $1
	`

	var query CashflowQuerier
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
		return CashflowModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return CashflowModel{}, err
	}

	return m, nil
}

func (r *CashflowRepository) DeleteById(ctx context.Context, id uint64) error {
	sqlQuery := `
		UPDATE cashflows
		SET deleted_at = $1
		WHERE id = $2
	`

	var exec CasflowExecutor
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

func (r *CashflowRepository) Query(ctx context.Context, id, limit int64) ([]CashflowModel, error) {
	fromId := "id > $1"
	if id != 0 {
		fromId = "id < $1"
	}

	sqlQuery := `
		SELECT 
			id,
			date,
			idr_amount,
			type,
			note,
			prove_file_url,
			created_at,
			updated_at,
			deleted_at
		FROM cashflows 
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

	var mps []*CashflowModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []CashflowModel{}, err
	}

	ms := make([]CashflowModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *CashflowRepository) QueryAmtByStatus(ctx context.Context, status CashflowType) ([]string, error) {
	sqlQuery := `
		SELECT 
			idr_amount
		FROM cashflows 
		WHERE deleted_at IS NULL
			AND type = $1
	`

	rows, _ := r.PostgreDb.Query(
		context.Background(),
		sqlQuery,
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
