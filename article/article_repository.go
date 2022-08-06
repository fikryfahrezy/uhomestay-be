package article

import (
	"context"
	"time"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ArticleRepository struct {
	ImgCacheName string
	RedisCl      *redis.Client
	PostgreDb    *pgxpool.Pool
}

func NewRepository(
	imgCacheName string,
	redisCl *redis.Client,
	postgreDb *pgxpool.Pool,
) *ArticleRepository {
	return &ArticleRepository{
		ImgCacheName: imgCacheName,
		RedisCl:      redisCl,
		PostgreDb:    postgreDb,
	}
}

type (
	ArticleExecutor   func(ctx context.Context, sql string, arguments ...interface{}) (commandTag pgconn.CommandTag, err error)
	ArticleQuerierRow func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	ArticleQuerier    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
)

func (r *ArticleRepository) Save(ctx context.Context, m ArticleModel) (nm ArticleModel, err error) {
	sqlQuery := `
		INSERT INTO articles (
			title,
			short_desc,
			thumbnail_url,
			content,
			content_text,
			slug,
			created_at,
			updated_at,
			deleted_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	var queryRow ArticleQuerierRow
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
		m.Title,
		m.ShortDesc,
		m.ThumbnailUrl,
		m.Content,
		m.ContentText,
		m.Slug,
		t,
		t,
		nil,
	).Scan(&lastInsertId)

	if err != nil {
		return ArticleModel{}, err
	}

	m.Id = lastInsertId
	m.CreatedAt = t
	m.UpdatedAt = t

	return m, nil
}

func (r *ArticleRepository) Query(ctx context.Context, q string, id, limit int64) ([]ArticleModel, error) {
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
			title,
			short_desc,
			thumbnail_url,
			content,
			content_text,
			slug,
			created_at,
			updated_at,
			deleted_at
		FROM articles 
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

	var mps []*ArticleModel
	if err := pgxscan.ScanAll(&mps, rows); err != nil {
		return []ArticleModel{}, err
	}

	ms := make([]ArticleModel, len(mps))
	for i, m := range mps {
		ms[i] = *m
	}

	return ms, nil
}

func (r *ArticleRepository) FindUndeletedById(ctx context.Context, id uint64) (m ArticleModel, err error) {
	querystr := `
		SELECT
			id,
			title,
			thumbnail_url,
			content,
			content_text,
			slug,
			created_at,
			updated_at,
			deleted_at
		FROM articles
		WHERE deleted_at IS NULL
		AND id = $1
	`

	var query ArticleQuerier
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
		return ArticleModel{}, err
	}

	if err = pgxscan.ScanOne(&m, rows); err != nil {
		return ArticleModel{}, err
	}

	return m, nil
}

func (r *ArticleRepository) UpdateById(ctx context.Context, id uint64, m ArticleModel) error {
	sqlQuery := `
		UPDATE articles SET (
			title,
			short_desc,
			content,
			content_text,
			thumbnail_url,
			updated_at
		) = ($1, $2, $3, $4, $5, $6)
		WHERE id = $7
	`

	var exec ArticleExecutor
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
		m.Title,
		m.ShortDesc,
		m.Content,
		m.ContentText,
		m.ThumbnailUrl,
		t,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ArticleRepository) DeleteById(ctx context.Context, id uint64) error {
	sqlQuery := `
		UPDATE articles
		SET deleted_at = $1
		WHERE id = $2
	`

	var exec ArticleExecutor
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

func (r *ArticleRepository) SetImgUrlCache(ctx context.Context, imgId, imgUrl string) (err error) {
	_, err = r.RedisCl.HSet(ctx, r.ImgCacheName, map[string]interface{}{imgId: imgUrl}).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *ArticleRepository) GetImgUrlsCache(ctx context.Context) (res map[string]string, err error) {
	vals, err := r.RedisCl.HGetAll(ctx, r.ImgCacheName).Result()
	if err != nil {
		return nil, err
	}

	return vals, nil
}

func (r *ArticleRepository) DelImgUrlCache(ctx context.Context) (err error) {
	_, err = r.RedisCl.Del(ctx, r.ImgCacheName).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *ArticleRepository) CountArticle(ctx context.Context) (n int64, err error) {
	sqlQuery := `
		SELECT COUNT(id) AS n
		FROM articles 
		WHERE deleted_at IS NULL
	`

	var queryRow ArticleQuerierRow
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
