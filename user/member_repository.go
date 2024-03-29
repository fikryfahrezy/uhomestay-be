package user

import (
	"context"
	"time"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type MemberRepository struct {
	PostgreDb *pgxpool.Pool
}

func NewMemberRepository(postgreDb *pgxpool.Pool) *MemberRepository {
	return &MemberRepository{
		PostgreDb: postgreDb,
	}
}

type (
	MemberExecutor   func(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	MemberQuerierRow func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	MemberQuerier    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
)

func (r *MemberRepository) Save(ctx context.Context, m MemberModel) error {
	sqlQuery := `
		INSERT INTO members(
			id,
			name,
			other_phone,
			wa_phone,
			profile_pic_url,
			id_card_url,
			username,
			password,
			is_admin,
			is_approved,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	var exec MemberExecutor
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
		m.Id,
		m.Name,
		m.OtherPhone,
		m.WaPhone,
		m.ProfilePicUrl,
		m.IdCardUrl,
		m.Username,
		m.Password,
		m.IsAdmin,
		m.IsApproved,
		t,
		t,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *MemberRepository) CheckUniqueField(ctx context.Context, m MemberModel) (em MemberModel, err error) {
	sqlQuery := `
		SELECT id
		FROM members
		WHERE (
			username = $1
			OR other_phone = $2
			OR wa_phone = $3
		) AND deleted_at IS NULL
		LIMIT 1
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
		m.Username,
		m.OtherPhone,
		m.WaPhone,
	).Scan(&em.Id)

	if err != nil {
		return MemberModel{}, err
	}

	return em, nil
}

func (r *MemberRepository) FindByUsername(username string) (m MemberModel, err error) {
	sqlQuery := `
		SELECT
			id,
			name,
			other_phone,
			wa_phone,
			profile_pic_url,
			id_card_url,
			username,
			password,
			is_admin,
			is_approved,
			created_at,
			updated_at,
			deleted_at
		FROM members
		WHERE deleted_at IS NULL
		AND username = $1
	`

	rows, _ := r.PostgreDb.Query(
		context.Background(),
		sqlQuery,
		username,
	)

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return MemberModel{}, err
	}

	return m, nil
}

func (r *MemberRepository) CheckOtherUniqueField(ctx context.Context, uid string, m MemberModel) (em MemberModel, err error) {
	sqlQuery := `
		SELECT id
		FROM members
		WHERE (
			username = $1
			OR other_phone = $2
			OR wa_phone = $3
			OR other_phone = $4
			OR wa_phone = $5
		) AND deleted_at IS NULL
		AND id != $6 LIMIT 1
	`

	var queryRow MemberQuerierRow
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		queryRow = tx.QueryRow
	} else {
		queryRow = r.PostgreDb.QueryRow
	}

	var xid pgtypeuuid.UUID
	err = queryRow(
		context.Background(),
		sqlQuery,
		m.Username,
		m.OtherPhone,
		m.WaPhone,
		m.WaPhone,
		m.OtherPhone,
		uid,
	).Scan(&xid)

	if err != nil {
		return MemberModel{}, err
	}

	return m, nil
}

func (r *MemberRepository) Update(ctx context.Context, id string, m MemberModel) error {
	sqlQuery := `
		UPDATE members SET (
			name,
			other_phone,
			wa_phone,
			profile_pic_url,
			id_card_url,
			username,
			password,
			is_admin,
			is_approved,
			updated_at
		) = ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		WHERE id = $11
	`

	var exec MemberExecutor
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
		m.OtherPhone,
		m.WaPhone,
		m.ProfilePicUrl,
		m.IdCardUrl,
		m.Username,
		m.Password,
		m.IsAdmin,
		m.IsApproved,
		t,
		id,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *MemberRepository) FindById(ctx context.Context, uid string) (m MemberModel, err error) {
	sqlQuery := `
		SELECT
			id,
			name,
			other_phone,
			wa_phone,
			profile_pic_url,
			id_card_url,
			username,
			password,
			is_admin,
			is_approved,
			created_at,
			updated_at,
			deleted_at
		FROM members
		WHERE deleted_at IS NULL
		AND id = $1
	`

	var query MemberQuerier
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
		uid,
	)

	if err != nil {
		return MemberModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return MemberModel{}, err
	}

	return m, nil
}

func (r *MemberRepository) DeleteById(ctx context.Context, uid string) error {
	sqlQuery := `
		UPDATE members
		SET
			username = CONCAT(username, $1::text),
			wa_phone = CONCAT(wa_phone, $2::text),
			other_phone = CONCAT(other_phone, $3::text),
			deleted_at = $4
		WHERE id = $5
	`

	var exec MemberExecutor
	tx, ok := ctx.Value(arbitary.TrxX{}).(pgx.Tx)
	if ok {
		exec = tx.Exec
	} else {
		exec = r.PostgreDb.Exec
	}

	var err error
	t := time.Now()
	idFraction := "-" + uid[:8]

	_, err = exec(
		context.Background(),
		sqlQuery,
		idFraction,
		idFraction,
		idFraction,
		t,
		uid,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *MemberRepository) QueryInId(ctx context.Context, uids []string) (ms []MemberModel, err error) {
	sqlQuery := `
		SELECT
			id,
			name,
			other_phone,
			wa_phone,
			profile_pic_url,
			id_card_url,
			username,
			password,
			is_admin,
			is_approved,
			created_at,
			updated_at,
			deleted_at
		FROM members
		WHERE deleted_at IS NULL
		AND id = ANY($1)
	`

	var query MemberQuerier
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
		uids,
	)

	defer rows.Close()

	var mps []*MemberModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []MemberModel{}, err
	}

	ms = make([]MemberModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *MemberRepository) Query(ctx context.Context, uid pgtypeuuid.UUID, q string, t time.Time, limit int64) (ms []MemberModel, err error) {
	fromUid := "id > $1"
	if !uid.UUID.IsNil() {
		fromUid = "id < $1"
	}

	created := "created_at >= $2::timestamp"
	if !t.IsZero() {
		created = "created_at <= $2::timestamp"
	}

	like := "username LIKE $3"
	order := "id"
	if q != "" {
		q = q + ":*"
		like = "textsearchable_index_col @@ websearch_to_tsquery($3)"
		order = "textrank_index_col"
	}

	if q == "" {
		q = "%" + q + "%"
	}

	sqlQuery := `
		SELECT
			id,
			name,
			other_phone,
			wa_phone,
			profile_pic_url,
			id_card_url,
			username,
			password,
			is_admin,
			is_approved,
			created_at,
			updated_at,
			deleted_at
		FROM members
		WHERE deleted_at IS NULL
			AND ` + fromUid + `
			AND ` + created + `
			AND ` + like + `
		ORDER BY ` + order + ` DESC
		LIMIT $4
	`

	rows, _ := r.PostgreDb.Query(
		context.Background(),
		sqlQuery,
		uid.UUID.String(),
		t.Format(time.RFC3339),
		q,
		limit,
	)
	defer rows.Close()

	var mps []*MemberModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []MemberModel{}, err
	}

	ms = make([]MemberModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *MemberRepository) CountMember(ctx context.Context) (n int64, err error) {
	sqlQuery := `
		SELECT COUNT(id) AS n
		FROM members
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

func (r *MemberRepository) SaveMemberHomestay(ctx context.Context, m MemberHomestayModel) (nm MemberHomestayModel, err error) {
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

	var queryRow MemberQuerierRow
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

func (r *MemberRepository) SaveHomestayPhoto(ctx context.Context, m HomestayImageModel) (nm HomestayImageModel, err error) {
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

	var queryRow MemberQuerierRow
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
