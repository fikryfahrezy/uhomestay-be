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

type GoalRepository struct {
	PostgreDb *pgxpool.Pool
}

func NewGoalRepository(postgreDb *pgxpool.Pool) *GoalRepository {
	return &GoalRepository{
		PostgreDb: postgreDb,
	}
}

type (
	GoalExecutor   func(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	GoalQuerierRow func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	GoalQuerier    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
)

func (r *GoalRepository) Save(ctx context.Context, g GoalModel) (gm GoalModel, err error) {
	sqlQuery := `
		INSERT INTO goals (
			vision,
			mission,
			org_period_id,
			created_at
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var queryRow GoalQuerierRow
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
		g.Vision,
		g.Mission,
		g.OrgPeriodId,
		t,
	).Scan(&lastInsertId)

	if err != nil {
		return GoalModel{}, err
	}

	g.Id = lastInsertId
	g.CreatedAt = t

	return g, nil
}

func (r *GoalRepository) FindByOrgPeriodId(ctx context.Context, orgPeriodId uint64) (gm GoalModel, err error) {
	sqlQuery := `
		SELECT 
			id,
			vision,
			mission,
			org_period_id,
			created_at
		FROM goals
		WHERE org_period_id = $1
		ORDER BY id DESC
		LIMIT 1
	`

	var query GoalQuerier
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
		orgPeriodId,
	)

	if err != nil {
		return GoalModel{}, err
	}

	if err = pgxscan.ScanOne(&gm, rows); err != nil {
		return GoalModel{}, err
	}

	return gm, nil
}
