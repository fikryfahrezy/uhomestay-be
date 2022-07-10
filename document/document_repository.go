package document

import (
	"context"
	"time"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DocumentRepository struct {
	PostgreDb *pgxpool.Pool
}

func NewRepository(postgreDb *pgxpool.Pool) *DocumentRepository {
	return &DocumentRepository{
		PostgreDb: postgreDb,
	}
}

type (
	DocumentExecutor   func(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	DocumentQuerierRow func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	DocumentQuerier    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
)

func (r *DocumentRepository) Save(ctx context.Context, m DocumentModel) (nm DocumentModel, err error) {
	sqlQuery := `
		INSERT INTO documents (
			name,
			url,
			type,
			dir_id,
			is_private,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	var queryRow DocumentQuerierRow
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
		m.Url,
		m.Type,
		m.DirId,
		m.IsPrivate,
		t,
		t,
		nil,
	).Scan(&lastInsertId)

	if err != nil {
		return DocumentModel{}, err
	}

	m.Id = lastInsertId
	m.CreatedAt = t
	m.UpdatedAt = t

	return m, nil
}

func (r *DocumentRepository) UpdateById(ctx context.Context, id uint64, m DocumentModel) error {
	sqlQuery := `
		UPDATE documents SET (
			name,
			url,
			is_private,
			updated_at
		) = ($1, $2, $3, $4)
		WHERE id = $5
	`

	var exec DocumentExecutor
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
		m.Url,
		m.IsPrivate,
		t,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *DocumentRepository) FindById(ctx context.Context, id uint64) (m DocumentModel, err error) {
	querystr := `
		SELECT
			id,
			name,
			url,
			type,
			dir_id,
			is_private,
			created_at,
			updated_at,
			deleted_at
		FROM documents 
		WHERE deleted_at IS NULL
		AND id = $1
	`

	var query DocumentQuerier
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
		return DocumentModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return DocumentModel{}, err
	}

	return m, nil
}

func (r *DocumentRepository) FindDirById(ctx context.Context, id uint64) (m DocumentModel, err error) {
	querystr := `
		SELECT
			id,
			name,
			url,
			type,
			dir_id,
			is_private,
			created_at,
			updated_at,
			deleted_at
		FROM documents 
		WHERE deleted_at IS NULL
			AND type = 'dir'
			AND id = $1
	`

	var query DocumentQuerier
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
		return DocumentModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return DocumentModel{}, err
	}

	return m, nil
}

func (r *DocumentRepository) DeleteInId(ctx context.Context, ids []uint64) error {
	sqlQuery := `
		UPDATE documents
		SET deleted_at = $1
		WHERE id = ANY($2)
	`

	var exec DocumentExecutor
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
		ids,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *DocumentRepository) Query(ctx context.Context, q string, id, limit int64) ([]DocumentModel, error) {
	fromId := "id > $1"
	if id != 0 {
		fromId = "id < $1"
	}

	like := "id > $2"
	order := "id"
	if q != "" {
		q = q + ":*"
		like = "textsearchable_index_col @@ websearch_to_tsquery($2)"
		order = "textrank_index_col"
	}

	if q == "" {
		q = "0"
	}

	sqlQuery := `
		SELECT 
			id,
			name,
			url,
			type,
			dir_id,
			is_private,
			created_at,
			updated_at,
			deleted_at
		FROM documents 
		WHERE deleted_at IS NULL
			AND ` + fromId + `
			AND ` + like + `
		ORDER BY ` + order + ` DESC
		LIMIT $3
	`

	rows, _ := r.PostgreDb.Query(
		context.Background(),
		sqlQuery,
		id,
		q,
		limit,
	)
	defer rows.Close()

	var mps []*DocumentModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []DocumentModel{}, err
	}

	ms := make([]DocumentModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *DocumentRepository) FindChildren(ctx context.Context, dirId uint64, q string, id, limit int64) ([]DocumentModel, error) {
	fromId := "id > $2"
	if id != 0 {
		fromId = "id < $2"
	}

	like := "id > $3"
	order := "id"
	if q != "" {
		q = q + ":*"
		like = "textsearchable_index_col @@ websearch_to_tsquery($3)"
		order = "textrank_index_col"
	}

	if q == "" {
		q = "0"
	}

	sqlQuery := `
		SELECT 
			id,
			name,
			url,
			type,
			dir_id,
			is_private,
			created_at,
			updated_at,
			deleted_at
		FROM documents 
		WHERE deleted_at IS NULL
			AND dir_id = $1
			AND ` + fromId + `
			AND ` + like + `
		ORDER BY ` + order + ` DESC
		LIMIT $4
	`

	rows, _ := r.PostgreDb.Query(
		context.Background(),
		sqlQuery,
		dirId,
		id,
		q,
		limit,
	)
	defer rows.Close()

	var mps []*DocumentModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []DocumentModel{}, err
	}

	ms := make([]DocumentModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *DocumentRepository) FindAllChildren(ctx context.Context, dirId uint64) ([]DocumentModel, error) {
	sqlQuery := `
		SELECT 
			id,
			name,
			url,
			type,
			dir_id,
			is_private,
			created_at,
			updated_at,
			deleted_at
		FROM documents 
		WHERE deleted_at IS NULL
			AND dir_id = $1
		ORDER BY id DESC
	`

	rows, _ := r.PostgreDb.Query(
		context.Background(),
		sqlQuery,
		dirId,
	)
	defer rows.Close()

	var mps []*DocumentModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []DocumentModel{}, err
	}

	ms := make([]DocumentModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}
