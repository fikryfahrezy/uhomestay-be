package history

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HistoryRepository struct {
	MongoDbName     string
	MongoCollection string
	MongoDb         *mongo.Client
}

func NewRepository(
	mongoDbName string,
	mongoCollection string,
	mongoDb *mongo.Client,
) *HistoryRepository {
	return &HistoryRepository{
		MongoDbName:     mongoDbName,
		MongoCollection: mongoCollection,
		MongoDb:         mongoDb,
	}
}

func (r *HistoryRepository) Save(ctx context.Context, h HistoryModel) (id string, err error) {
	h.CreatedAt = time.Now()

	coll := r.MongoDb.Database(r.MongoDbName).Collection(r.MongoCollection)
	result, err := coll.InsertOne(ctx, h)
	if err != nil {
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", primitive.ErrInvalidHex
	}

	return oid.Hex(), nil
}

func (r *HistoryRepository) FindLatest(ctx context.Context) (h HistoryModel, err error) {
	coll := r.MongoDb.Database(r.MongoDbName).Collection(r.MongoCollection)

	filter := bson.M{}
	opts := options.FindOne().SetSort(bson.M{"created_at": -1})

	err = coll.FindOne(ctx, filter, opts).Decode(&h)
	if err != nil {
		return HistoryModel{}, err
	}

	return h, nil
}
