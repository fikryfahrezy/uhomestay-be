package homestay

import (
	"context"
	"time"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type HomestayImageRepository struct {
	PostgreDb *pgxpool.Pool
}

func NewHomestayImageRepository(postgreDb *pgxpool.Pool) *HomestayImageRepository {
	return &HomestayImageRepository{
		PostgreDb: postgreDb,
	}
}

type (
	HomestayImageExecutor           func(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	HomestayHomestayImageQuerierRow func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	HomestayImageQuerier            func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
)

func (r *HomestayImageRepository) Save(ctx context.Context, m HomestayImageModel) (nm HomestayImageModel, err error) {
	sqlQuery := `
		INSERT INTO homestay_images (
			name,
			alphnum_name,
			url,
			member_homestay_id,
			created_at,
			deleted_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var queryRow HomestayHomestayImageQuerierRow
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
		m.MemberHomestayId,
		t,
		nil,
	).Scan(&lastInsertId)

	if err != nil {
		return HomestayImageModel{}, err
	}

	m.Id = lastInsertId
	m.CreatedAt = t

	return m, nil
}

func (r *HomestayImageRepository) FindById(ctx context.Context, id uint64) (m HomestayImageModel, err error) {
	querystr := `
		SELECT
			id,
			name,
			alphnum_name,
			url,
			member_homestay_id,
			created_at,
			deleted_at
		FROM homestay_images 
		WHERE deleted_at IS NULL
		AND id = $1
	`

	var query HomestayImageQuerier
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
		return HomestayImageModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return HomestayImageModel{}, err
	}

	return m, nil
}

func (r *HomestayImageRepository) FindByMemberHomestayId(ctx context.Context, id uint64) ([]HomestayImageModel, error) {
	sqlQuery := `
		SELECT
			id,
			name,
			alphnum_name,
			url,
			member_homestay_id,
			created_at,
			deleted_at
		FROM homestay_images 
		WHERE deleted_at IS NULL
			AND member_homestay_id = $1
		ORDER BY id ASC
	`

	rows, _ := r.PostgreDb.Query(
		context.Background(),
		sqlQuery,
		id,
	)
	defer rows.Close()

	var mps []*HomestayImageModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []HomestayImageModel{}, err
	}

	ms := make([]HomestayImageModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *HomestayImageRepository) DeleteById(ctx context.Context, id uint64) error {
	sqlQuery := `
		UPDATE homestay_images
		SET deleted_at = $1
		WHERE id = $2
	`

	var exec HomestayImageExecutor
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

func (r *HomestayImageRepository) QueryInId(ctx context.Context, ids []uint64) (ms []HomestayImageModel, err error) {
	sqlQuery := `
		SELECT 
			id,
			name,
			alphnum_name,
			url,
			member_homestay_id,
			created_at,
			deleted_at
		FROM homestay_images
		WHERE deleted_at IS NULL
		AND id = ANY($1)
		ORDER BY id ASC
	`

	var query HomestayImageQuerier
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

	var mps []*HomestayImageModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []HomestayImageModel{}, err
	}

	ms = make([]HomestayImageModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *HomestayImageRepository) UpdateHomestayIdInId(ctx context.Context, id uint64, ids []uint64) error {
	sqlQuery := `
	UPDATE homestay_images
	SET member_homestay_id = $1
	WHERE id = ANY($2)
`

	var exec HomestayImageExecutor
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		exec = tx.Exec
	} else {
		exec = r.PostgreDb.Exec
	}

	_, err := exec(
		context.Background(),
		sqlQuery,
		id,
		ids,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *HomestayImageRepository) DeleteInId(ctx context.Context, ids []uint64) error {
	sqlQuery := `
	UPDATE homestay_images
	SET deleted_at = $1
	WHERE id = ANY($2)
`

	var exec HomestayImageExecutor
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		exec = tx.Exec
	} else {
		exec = r.PostgreDb.Exec
	}

	_, err := exec(
		context.Background(),
		sqlQuery,
		time.Now(),
		ids,
	)
	if err != nil {
		return err
	}

	return nil
}
