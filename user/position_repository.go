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

type PositionRepository struct {
	PostgreDb *pgxpool.Pool
}

func NewPositionRepository(postgreDb *pgxpool.Pool) *PositionRepository {
	return &PositionRepository{
		PostgreDb: postgreDb,
	}
}

type (
	PositionExecutor   func(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	PositionQuerierRow func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	PositionQuerier    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
)

func (r *PositionRepository) Save(ctx context.Context, m PositionModel) (nm PositionModel, err error) {
	sqlQuery := `
		INSERT INTO positions (
			name,
			level,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var queryRow PositionQuerierRow
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
		m.Name,
		m.Level,
		t,
		t,
		nil,
	).Scan(&lastInsertId)

	if err != nil {
		return PositionModel{}, err
	}

	m.Id = lastInsertId
	m.CreatedAt = t
	m.UpdatedAt = t

	return m, nil
}

func (r *PositionRepository) UpdateById(ctx context.Context, id uint64, m PositionModel) error {
	sqlQuery := `
		UPDATE positions SET (
			name,
			level,
			updated_at
		) = ($1, $2, $3)
		WHERE id = $4
	`

	var exec PositionExecutor
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
		m.Name,
		m.Level,
		t,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *PositionRepository) FindUndeletedById(ctx context.Context, id uint64) (m PositionModel, err error) {
	sqlQuery := `
		SELECT 
			id,
			name,
			level,
			created_at,
			updated_at,
			deleted_at
		FROM positions
		WHERE deleted_at IS NULL
		AND id = $1
	`

	var query PositionQuerier
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
		return PositionModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return PositionModel{}, err
	}

	return m, nil
}

func (r *PositionRepository) DeleteById(ctx context.Context, id uint64) error {
	sqlQuery := `
		UPDATE positions
		SET deleted_at = $1 
		WHERE id = $2
	`

	var exec PositionExecutor
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

func (r *PositionRepository) Query(ctx context.Context, id, limit int64) ([]PositionModel, error) {
	fromId := "id > $1"
	if id != 0 {
		fromId = "id < $1"
	}

	sqlQuery := `
		SELECT 
			id,
			name,
			level,
			created_at,
			updated_at,
			deleted_at
		FROM positions
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

	var err error
	var mps []*PositionModel
	if err = pgxscan.ScanAll(&mps, rows); err != nil {
		return []PositionModel{}, err
	}

	ms := make([]PositionModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *PositionRepository) MaxLevel(ctx context.Context) (level int64, err error) {
	sqlQuery := `
		SELECT 
			MAX(level)
		FROM positions
		WHERE deleted_at IS NULL
	`
	var queryRow PositionQuerierRow
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		queryRow = tx.QueryRow
	} else {
		queryRow = r.PostgreDb.QueryRow
	}

	err = queryRow(context.Background(), sqlQuery).Scan(&level)

	return level, nil
}

func (r *PositionRepository) QueryUndeletedInId(ctx context.Context, ids []uint64) (ms []PositionModel, err error) {
	sqlQuery := `
		SELECT 
			id,
			name,
			level,
			created_at,
			updated_at,
			deleted_at
		FROM positions
		WHERE id = ANY($1)
	`

	var query PositionQuerier
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
		ids,
	)
	defer rows.Close()

	var mps []*PositionModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []PositionModel{}, err
	}

	ms = make([]PositionModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *PositionRepository) QueryInId(ctx context.Context, ids []uint64) (ms []PositionModel, err error) {
	sqlQuery := `
		SELECT 
			id,
			name,
			level,
			created_at,
			updated_at,
			deleted_at
		FROM positions
		WHERE deleted_at IS NULL
		AND id = ANY($1)
	`

	var query PositionQuerier
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
		ids,
	)
	defer rows.Close()

	var mps []*PositionModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []PositionModel{}, err
	}

	ms = make([]PositionModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *PositionRepository) FindById(ctx context.Context, id uint64) (m PositionModel, err error) {
	sqlQuery := `
		SELECT 
			id,
			name,
			level,
			created_at,
			updated_at,
			deleted_at
		FROM positions
		WHERE id = $1
	`
	var query PositionQuerier
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
		return PositionModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return PositionModel{}, err
	}

	return m, nil
}

func (r *PositionRepository) CountPosition(ctx context.Context) (n int64, err error) {
	sqlQuery := `
		SELECT COUNT(id) AS n
		FROM positions
		WHERE deleted_at IS NULL
	`

	var queryRow MemberQuerierRow
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		queryRow = tx.QueryRow
	} else {
		queryRow = r.PostgreDb.QueryRow
	}

	err = queryRow(
		context.Background(),
		sqlQuery,
	).Scan(&n)

	if err != nil {
		return 0, err
	}

	return n, nil
}
