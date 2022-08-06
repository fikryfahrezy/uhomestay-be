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

type MemberHomestayRepository struct {
	PostgreDb *pgxpool.Pool
}

func NewMemberHomestayRepository(postgreDb *pgxpool.Pool) *MemberHomestayRepository {
	return &MemberHomestayRepository{
		PostgreDb: postgreDb,
	}
}

type (
	MemberHomestayExecutor   func(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	MemberHomestayQuerierRow func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	MemberHomestayQuerier    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
)

func (r *MemberHomestayRepository) Save(ctx context.Context, m MemberHomestayModel) (nm MemberHomestayModel, err error) {
	sqlQuery := `
		INSERT INTO member_homestays (
			name,
			address,
			latitude,
			longitude,
			thumbnail_url,
			member_id,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	var queryRow MemberHomestayQuerierRow
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
		m.Address,
		m.Latitude,
		m.Longitude,
		m.ThumbnailUrl,
		m.MemberId,
		t,
		t,
		nil,
	).Scan(&lastInsertId)

	if err != nil {
		return MemberHomestayModel{}, err
	}

	m.Id = lastInsertId
	m.CreatedAt = t

	return m, nil
}

func (r *MemberHomestayRepository) UpdateById(ctx context.Context, uid string, id uint64, m MemberHomestayModel) error {
	sqlQuery := `
		UPDATE member_homestays SET (
			name,
			address,
			latitude,
			longitude,
			thumbnail_url,
			member_id,
			created_at,
			updated_at
		) = ($1, $2, $3, $4, $5, $6, $7, $8)
		WHERE member_id = $9 AND id = $10
	`

	var exec MemberHomestayExecutor
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
		m.Address,
		m.Latitude,
		m.Longitude,
		m.ThumbnailUrl,
		m.MemberId,
		t,
		t,
		uid,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *MemberHomestayRepository) FindById(ctx context.Context, uid string, id uint64) (m MemberHomestayModel, err error) {
	querystr := `
		SELECT
			id,
			name,
			address,
			latitude,
			longitude,
			thumbnail_url,
			member_id,
			created_at,
			updated_at,
			deleted_at
		FROM member_homestays 
		WHERE deleted_at IS NULL
		AND member_id = $1 AND id = $2
	`

	var query MemberHomestayQuerier
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
		uid,
		id,
	)

	if err != nil {
		return MemberHomestayModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return MemberHomestayModel{}, err
	}

	return m, nil
}

func (r *MemberHomestayRepository) Query(ctx context.Context, uid string, id, limit int64) ([]MemberHomestayModel, error) {
	fromId := "id > $1"
	if id != 0 {
		fromId = "id < $1"
	}

	sqlQuery := `
		SELECT
			id,
			name,
			address,
			latitude,
			longitude,
			thumbnail_url,
			member_id,
			created_at,
			updated_at,
			deleted_at
		FROM member_homestays 
		WHERE deleted_at IS NULL
			AND ` + fromId + `
			AND member_id = $2
		ORDER BY id DESC
		LIMIT $3
	`

	rows, _ := r.PostgreDb.Query(
		context.Background(),
		sqlQuery,
		id,
		uid,
		limit,
	)
	defer rows.Close()

	var mps []*MemberHomestayModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []MemberHomestayModel{}, err
	}

	ms := make([]MemberHomestayModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *MemberHomestayRepository) DeleteById(ctx context.Context, uid string, id uint64) error {
	sqlQuery := `
		UPDATE member_homestays
		SET deleted_at = $1
		WHERE id = $2
		AND member_id = $3
	`

	var exec MemberHomestayExecutor
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
		uid,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *MemberHomestayRepository) CountMemberHomestay(ctx context.Context, uid string) (n int64, err error) {
	sqlQuery := `
		SELECT COUNT(id) AS n
		FROM member_homestays
		WHERE deleted_at IS NULL
		AND member_id = $1
	`

	var queryRow MemberHomestayQuerierRow
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		queryRow = tx.QueryRow
	} else {
		queryRow = r.PostgreDb.QueryRow
	}

	err = queryRow(
		context.Background(),
		sqlQuery,
		uid,
	).Scan(&n)

	if err != nil {
		return 0, err
	}

	return n, nil
}
