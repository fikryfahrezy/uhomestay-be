package image

import (
	"context"
	"time"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ImageRepository struct {
	PostgreDb *pgxpool.Pool
}

func NewRepository(postgreDb *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{
		PostgreDb: postgreDb,
	}
}

type (
	ImageExecutor   func(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	ImageQuerierRow func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	ImageQuerier    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
)

func (r *ImageRepository) Save(ctx context.Context, m ImageModel) (nm ImageModel, err error) {
	sqlQuery := `
		INSERT INTO images (
			name,
			alphnum_name,
			url,
			description,
			created_at,
			deleted_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var queryRow ImageQuerierRow
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
		m.AlphnumName,
		m.Url,
		m.Description,
		t,
		nil,
	).Scan(&lastInsertId)

	if err != nil {
		return ImageModel{}, err
	}

	m.Id = lastInsertId
	m.CreatedAt = t

	return m, nil
}

func (r *ImageRepository) UpdateById(ctx context.Context, id uint64, m ImageModel) error {
	sqlQuery := `
		UPDATE images SET (
			name,
			alphnum_name,
			url,
			description,
			updated_at
		) = ($1, $2, $3, $4, $5)
		WHERE id = $6
	`

	var exec ImageExecutor
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
		m.AlphnumName,
		m.Url,
		m.Description,
		t,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ImageRepository) FindById(ctx context.Context, id uint64) (m ImageModel, err error) {
	querystr := `
		SELECT
			id,
			name,
			alphnum_name,
			url,
			description,
			created_at,
			deleted_at
		FROM images 
		WHERE deleted_at IS NULL
		AND id = $1
	`

	var query ImageQuerier
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
		return ImageModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return ImageModel{}, err
	}

	return m, nil
}

func (r *ImageRepository) Query(ctx context.Context, id, limit int64) ([]ImageModel, error) {
	fromId := "id > $1"
	if id != 0 {
		fromId = "id < $1"
	}

	sqlQuery := `
		SELECT
			id,
			name,
			alphnum_name,
			url,
			description,
			created_at,
			deleted_at
		FROM images 
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

	var mps []*ImageModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []ImageModel{}, err
	}

	ms := make([]ImageModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *ImageRepository) DeleteById(ctx context.Context, id uint64) error {
	sqlQuery := `
		UPDATE images
		SET deleted_at = $1
		WHERE id = $2
	`

	var exec ImageExecutor
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

func (r *ImageRepository) CountImage(ctx context.Context) (n int64, err error) {
	sqlQuery := `
		SELECT COUNT(id) AS n
		FROM images
		WHERE deleted_at IS NULL
	`

	var queryRow ImageQuerierRow
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
