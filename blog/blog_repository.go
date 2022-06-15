package blog

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BlogRepository struct {
	MongoDbName     string
	MongoCollection string
	ImgCacheName    string
	RedisCl         *redis.Client
	MongoDb         *mongo.Client
}

func NewRepository(
	mongoDbName string,
	mongoCollection string,
	imgCacheName string,
	redisCl *redis.Client,
	mongoDb *mongo.Client,
) *BlogRepository {
	return &BlogRepository{
		MongoDbName:     mongoDbName,
		MongoCollection: mongoCollection,
		ImgCacheName:    imgCacheName,
		RedisCl:         redisCl,
		MongoDb:         mongoDb,
	}
}

func (r *BlogRepository) Save(ctx context.Context, b BlogModel) (id string, err error) {
	t := time.Now()
	b.CreatedAt = t
	b.UpdatedAt = t

	coll := r.MongoDb.Database(r.MongoDbName).Collection(r.MongoCollection)
	result, err := coll.InsertOne(ctx, b)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", primitive.ErrInvalidHex
	}

	return oid.Hex(), nil
}

func (r *BlogRepository) Query(ctx context.Context, idHex string, limit int64) (bs []BlogModel, err error) {
	coll := r.MongoDb.Database(r.MongoDbName).Collection(r.MongoCollection)

	filter := bson.M{"deleted_at": nil}
	if idHex != "" {
		id, _ := primitive.ObjectIDFromHex(idHex)
		filter["_id"] = bson.M{
			"$lt": id,
		}
	}

	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(limit)

	cursor, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return []BlogModel{}, err
	}

	if err = cursor.All(ctx, &bs); err != nil {
		return []BlogModel{}, err
	}

	return bs, nil
}

func (r *BlogRepository) FindUndeletedById(ctx context.Context, idHex string) (b BlogModel, err error) {
	coll := r.MongoDb.Database(r.MongoDbName).Collection(r.MongoCollection)
	id, _ := primitive.ObjectIDFromHex(idHex)

	filter := bson.M{"_id": id, "deleted_at": nil}
	err = coll.FindOne(ctx, filter).Decode(&b)
	if err != nil {
		return BlogModel{}, err
	}

	return b, nil
}

func (r *BlogRepository) UpdateById(ctx context.Context, idHex string, b BlogModel) (err error) {
	coll := r.MongoDb.Database(r.MongoDbName).Collection(r.MongoCollection)
	id, _ := primitive.ObjectIDFromHex(idHex)

	update := bson.M{
		"$set": bson.M{
			"title":         b.Title,
			"short_desc":    b.ShortDesc,
			"content":       b.Content,
			"thumbnail_url": b.ThumbnailUrl,
		},
	}
	filter := bson.M{"_id": id, "deleted_at": nil}
	_, err = coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *BlogRepository) DeleteById(ctx context.Context, idHex string) (err error) {
	coll := r.MongoDb.Database(r.MongoDbName).Collection(r.MongoCollection)
	id, _ := primitive.ObjectIDFromHex(idHex)

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}
	filter := bson.M{"_id": id, "deleted_at": nil}
	_, err = coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *BlogRepository) SetImgUrlCache(ctx context.Context, imgId, imgUrl string) (err error) {
	_, err = r.RedisCl.HSet(ctx, r.ImgCacheName, map[string]interface{}{imgId: imgUrl}).Result()
	if err != nil {
		return err
	}

	return nil
}

func (r *BlogRepository) GetImgUrlsCache(ctx context.Context) (res map[string]string, err error) {
	vals, err := r.RedisCl.HGetAll(ctx, r.ImgCacheName).Result()
	if err != nil {
		return nil, err
	}

	return vals, nil
}

func (r *BlogRepository) DelImgUrlCache(ctx context.Context) (err error) {
	_, err = r.RedisCl.Del(ctx, r.ImgCacheName).Result()
	if err != nil {
		return err
	}

	return nil
}
