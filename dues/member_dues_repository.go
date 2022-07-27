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

type MemberDuesRepository struct {
	PostgreDb *pgxpool.Pool
}

func NewMemberDeusRepository(postgreDb *pgxpool.Pool) *MemberDuesRepository {
	return &MemberDuesRepository{
		PostgreDb: postgreDb,
	}
}

type (
	MemberDuesExecutor         func(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	MemberMemberDuesQuerierRow func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	MemberDuesQuerier          func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
)

func (r *MemberDuesRepository) QueryMDVByUid(ctx context.Context, uid string, id, limit int64) ([]MemberDuesViewModel, int64, error) {
	from := `
		FROM member_dues md
			LEFT JOIN dues d ON d.id = md.dues_id
		WHERE d.deleted_at IS NULL
			AND md.deleted_at IS NULL
			AND md.member_id = $1
	`

	countSqlQuery := `
		SELECT COUNT(md.id) AS n
		` + from + `
	`

	var n int64
	err := r.PostgreDb.QueryRow(
		context.Background(),
		countSqlQuery,
		uid,
	).Scan(&n)
	if err != nil {
		return []MemberDuesViewModel{}, 0, err
	}

	fromId := "md.id > $2"
	if id != 0 {
		fromId = "md.id < $2"
	}

	selectSqlQuery := `
		SELECT 
			md.id,
			md.dues_id,
			d.date,
			md.status,
			d.idr_amount,
			md.prove_file_url,
			md.pay_date
		` + from + `
			AND ` + fromId + `
		ORDER BY md.id DESC
		LIMIT $3
	`

	rows, _ := r.PostgreDb.Query(
		context.Background(),
		selectSqlQuery,
		uid,
		id,
		limit,
	)
	defer rows.Close()

	var mps []*MemberDuesViewModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []MemberDuesViewModel{}, 0, err
	}

	ms := make([]MemberDuesViewModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, n, nil
}

func (r *MemberDuesRepository) QueryDMVByDuesId(ctx context.Context, duesId uint64, id, limit int64, startDate, endDate time.Time) ([]DuesMemberViewModel, error) {
	dateFilter := ""
	selectQueryParams := []interface{}{duesId, id, limit}

	if !startDate.IsZero() {
		dateFilter = dateFilter + `
			AND md.pay_date >= $4::timestamp
		`
		selectQueryParams = append(selectQueryParams, startDate.Format(time.RFC3339))
	}

	if !endDate.IsZero() {
		dateFilter = dateFilter + `
		AND md.pay_date <= $5::timestamp
		`
		selectQueryParams = append(selectQueryParams, endDate.Format(time.RFC3339))
	}

	fromUid := "md.id > $2"
	if id != 0 {
		fromUid = "md.id < $2"
	}

	sqlSelectQuery := `
		SELECT
			md.id,
			md.member_id,
			md.status,
			md.created_at,
			md.pay_date,
			m.name,
			m.profile_pic_url
		FROM member_dues md
			LEFT JOIN members m ON m.id = md.member_id
		WHERE md.deleted_at IS NULL
			AND md.dues_id = $1
			AND m.deleted_at IS NULL
			AND ` + fromUid + `
			` + dateFilter + `
		ORDER BY md.id DESC
		LIMIT $3
	`

	rows, _ := r.PostgreDb.Query(
		context.Background(),
		sqlSelectQuery,
		selectQueryParams...,
	)
	defer rows.Close()

	var mps []*DuesMemberViewModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []DuesMemberViewModel{}, err
	}

	ms := make([]DuesMemberViewModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *MemberDuesRepository) FindUnpaidById(ctx context.Context, id uint64) (m MemberDuesModel, err error) {
	querystr := `
		SELECT
			id,
			member_id,
			dues_id,
			status,
			prove_file_url,
			created_at,
			updated_at,
			deleted_at
		FROM member_dues 
		WHERE deleted_at IS NULL
			AND id = $1
			AND status != 'paid'
	`

	var query MemberDuesQuerier
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
		return MemberDuesModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return MemberDuesModel{}, err
	}

	return m, nil
}

func (r *MemberDuesRepository) FindUnpaidByIdAndMemberId(ctx context.Context, id uint64, uid string) (m MemberDuesModel, err error) {
	querystr := `
		SELECT
			id,
			member_id,
			dues_id,
			status,
			prove_file_url,
			created_at,
			updated_at,
			deleted_at
		FROM member_dues 
		WHERE deleted_at IS NULL
			AND id = $1
			AND member_id = $2
			AND status = 'unpaid'
	`

	var query MemberDuesQuerier
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
		uid,
	)

	if err != nil {
		return MemberDuesModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return MemberDuesModel{}, err
	}

	return m, nil
}

func (r *MemberDuesRepository) Save(ctx context.Context, m MemberDuesModel) (nm MemberDuesModel, err error) {
	sqlQuery := `
		INSERT INTO member_dues (
			member_id,
			dues_id,
			status,
			prove_file_url,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var queryRow MemberMemberDuesQuerierRow
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
		m.MemberId,
		m.DuesId,
		m.Status,
		m.ProveFileUrl,
		t,
		t,
		nil,
	).Scan(&lastInsertId)

	if err != nil {
		return MemberDuesModel{}, err
	}

	m.Id = lastInsertId
	m.CreatedAt = t
	m.UpdatedAt = t

	return m, nil
}

func (r *MemberDuesRepository) UpdateById(ctx context.Context, id uint64, m MemberDuesModel) error {
	sqlQuery := `
		UPDATE member_dues SET (
			prove_file_url,
			status,
			pay_date,
			updated_at
		) = ($1, $2, $3, $4)
		WHERE id = $5
	`

	var exec MemberDuesExecutor
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
		m.ProveFileUrl,
		m.Status,
		m.PayDate,
		t,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *MemberDuesRepository) GenerateDues(ctx context.Context, duesId uint64) (err error) {
	// Ref: PostgreSQL: insert from another table
	// https://stackoverflow.com/a/6898775/12976234
	sqlQuery := `
		INSERT INTO member_dues (
			dues_id,
			status, 
			member_id,
			created_at,
			updated_at,
			deleted_at
		) 
		SELECT
			$1,
			'unpaid',
			id,
			$2,
			$3,
			$4
		FROM members 
		WHERE deleted_at IS NULL 
			AND is_approved = true
	`

	var exec MemberDuesExecutor
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		exec = tx.Exec
	} else {
		exec = r.PostgreDb.Exec
	}

	t := time.Now()

	_, err = exec(
		context.Background(),
		sqlQuery,
		duesId,
		t,
		t,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *MemberDuesRepository) CheckSomeonePaid(ctx context.Context, duesId uint64) ([]MemberDuesModel, error) {
	sqlQuery := `
		SELECT
			dues_id,
			status, 
			member_id,
			created_at,
			updated_at,
			deleted_at
		FROM member_dues
		WHERE deleted_at IS NULL 
			AND status != 'unpaid'
			AND dues_id = $1
	`

	rows, _ := r.PostgreDb.Query(
		context.Background(),
		sqlQuery,
		duesId,
	)
	defer rows.Close()

	var mps []*MemberDuesModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []MemberDuesModel{}, err
	}

	ms := make([]MemberDuesModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *MemberDuesRepository) DeleteByDuesId(ctx context.Context, duesId uint64) error {
	sqlQuery := `
		UPDATE member_dues
		SET deleted_at = $1
		WHERE dues_id = $2
	`

	var exec MemberDuesExecutor
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
		duesId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *MemberDuesRepository) QueryAmtByDuesId(ctx context.Context, duesId uint64, startDate, endDate time.Time) ([]MemberDuesAmtViewModel, error) {
	sqlQuery := `
	SELECT
		d.idr_amount,
		md.status
	FROM dues d
		LEFT JOIN member_dues md ON md.dues_id = d.id
	WHERE d.deleted_at IS NULL
		AND md.deleted_at IS NULL
		AND d.id = $1
	`

	queryParams := []interface{}{duesId}
	if !startDate.IsZero() {
		sqlQuery = sqlQuery + `
			AND md.pay_date >= $2::timestamp
		`
		queryParams = append(queryParams, startDate.Format(time.RFC3339))
	}

	if !endDate.IsZero() {
		sqlQuery = sqlQuery + `
		AND md.pay_date <= $3::timestamp
		`
		queryParams = append(queryParams, endDate.Format(time.RFC3339))
	}

	rows, _ := r.PostgreDb.Query(context.Background(), sqlQuery, queryParams...)
	defer rows.Close()

	var mps []*MemberDuesAmtViewModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []MemberDuesAmtViewModel{}, err
	}

	ms := make([]MemberDuesAmtViewModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *MemberDuesRepository) CountDMVByDuesId(ctx context.Context, duesId uint64, startDate, endDate time.Time) (int64, error) {
	sqlCountQuery := `
	SELECT COUNT(md.id) AS n
	FROM member_dues md
			LEFT JOIN members m ON m.id = md.member_id
		WHERE md.deleted_at IS NULL
			AND md.dues_id = $1
			AND m.deleted_at IS NULL
	`
	countQueryParams := []interface{}{duesId}

	if !startDate.IsZero() {
		sqlCountQuery = sqlCountQuery + `
		AND md.pay_date >= $2::timestamp
		`
		countQueryParams = append(countQueryParams, startDate.Format(time.RFC3339))
	}

	if !endDate.IsZero() {
		sqlCountQuery = sqlCountQuery + `
		AND md.pay_date <= $3::timestamp
		`
		countQueryParams = append(countQueryParams, endDate.Format(time.RFC3339))
	}

	var n int64
	err := r.PostgreDb.QueryRow(
		context.Background(),
		sqlCountQuery,
		countQueryParams...,
	).Scan(&n)
	if err != nil {
		return 0, err
	}

	return n, nil
}
