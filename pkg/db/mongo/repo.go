package mongodb

import (
	"context"
	"time"

	"tenkhours/pkg/db/base"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo base repo implements IBaseRepo
type BaseRepo[M base.IBaseEntity, N any] struct {
	*mongo.Collection
	*Mapper[M, N]
}

func NewBaseRepo[M base.IBaseEntity, N any](collection *mongo.Collection, mapper *Mapper[M, N]) *BaseRepo[M, N] {
	return &BaseRepo[M, N]{collection, mapper}
}

func (r *BaseRepo[M, N]) CountAll(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.Collection.CountDocuments(ctx, bson.M{})
}

func (r *BaseRepo[M, N]) InsertMany(ctx context.Context, ms []M) ([]M, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	docs := make([]interface{}, len(ms))
	for i, m := range ms {
		addMissingFields(m)
		docs[i] = r.ToMongoEntity(&m)
	}

	_, err := r.Collection.InsertMany(ctx, docs)
	return ms, err
}

func (r *BaseRepo[M, N]) InsertOne(ctx context.Context, m *M) (*M, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	addMissingFields(*m)
	_, err := r.Collection.InsertOne(ctx, r.ToMongoEntity(m))
	return m, err
}

func (r *BaseRepo[M, N]) FindAll(ctx context.Context) ([]M, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var ms []M
	err = cursor.All(ctx, &ms)
	return ms, err
}

func (r *BaseRepo[M, N]) FindByID(ctx context.Context, id string) (*M, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var m M
	err := r.Collection.FindOne(ctx, bson.M{"_id": ToObjectID(id)}).Decode(&m)
	return &m, err
}

func (r *BaseRepo[M, N]) UpdateByID(ctx context.Context, id string, m *M) (*M, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	lo.FromPtr(m).SetUpdatedAtByNow()

	err := r.Collection.FindOneAndUpdate(ctx,
		bson.M{"_id": ToObjectID(id)},
		bson.M{"$set": r.ToMongoEntity(m)},
		FindOneAndUpdateOptions).Decode(&m)
	return m, err
}

func (r *BaseRepo[M, N]) DeleteByID(ctx context.Context, id string) (*M, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var m M
	err := r.Collection.FindOneAndDelete(ctx, bson.M{"_id": ToObjectID(id)}).Decode(&m)
	return &m, err
}
