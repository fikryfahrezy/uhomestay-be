package history

import (
	"context"
	"time"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type HistoryRepository struct {
	PostgreDb *pgxpool.Pool
}

func NewRepository(
	postgreDb *pgxpool.Pool,
) *HistoryRepository {
	return &HistoryRepository{
		PostgreDb: postgreDb,
	}
}

type (
	HistoryExecutor   func(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	HistoryQuerierRow func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	HistoryQuerier    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
)

func (r *HistoryRepository) Save(ctx context.Context, m HistoryModel) (nm HistoryModel, err error) {
	sqlQuery := `
		INSERT INTO histories (
			content_text,
			content,
			created_at
		)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var queryRow HistoryQuerierRow
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
		m.ContentText,
		m.Content,
		t,
	).Scan(&lastInsertId)

	if err != nil {
		return HistoryModel{}, err
	}

	m.Id = lastInsertId
	m.CreatedAt = t

	return m, nil
}

func (r *HistoryRepository) FindLatest(ctx context.Context) (m HistoryModel, err error) {
	querystr := `
		SELECT
			id,
			content_text,
			created_at,
			content
		FROM histories
		ORDER BY id DESC
		LIMIT 1
	`

	var query HistoryQuerier
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
		return HistoryModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return HistoryModel{}, err
	}

	return m, nil
}
