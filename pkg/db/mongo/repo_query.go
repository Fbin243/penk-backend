package mongodb

import (
	"context"

	"tenkhours/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Count returns the number of documents in the collection that match the filter.
func (r *BaseRepo[M, N]) Count(ctx context.Context, filter any) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	count, err := r.Collection.CountDocuments(ctx, filter)
	return int(count), err
}

// Exists checks if a document exists in the collection based on the provided filter.
func (r *BaseRepo[M, N]) Exists(ctx context.Context, filter any) error {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	count, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.ErrMongoNotFound
	}

	return nil
}

// FindByID retrieves a document by its ID from the collection.
func (r *BaseRepo[M, N]) FindByID(ctx context.Context, id string) (*M, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	var m M
	err := r.Collection.FindOne(ctx, bson.M{"_id": ToObjectID(id)}).Decode(&m)
	return &m, err
}

// FindByIDs retrieves multiple documents by their IDs from the collection.
func (r *BaseRepo[M, N]) FindByIDs(ctx context.Context, ids []string) ([]M, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	return r.FindMany(ctx, bson.M{"_id": bson.M{"$in": ToObjectIDs(ids)}})
}

// FindOne retrieves a single document from the collection based on the provided filter.
func (r *BaseRepo[M, N]) FindOne(ctx context.Context, filter any) (*M, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	var m M
	err := r.Collection.FindOne(ctx, filter).Decode(&m)
	return &m, err
}

// FindMany retrieves multiple documents from the collection based on the provided filter.
func (r *BaseRepo[M, N]) FindMany(ctx context.Context, filter any, opts ...*options.FindOptions) ([]M, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	cursor, err := r.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ms []M
	err = cursor.All(ctx, &ms)
	return ms, err
}

// AggregateQuery performs an aggregation operation on the collection using the provided pipeline.
func (r *BaseRepo[M, N]) AggregateQuery(ctx context.Context, pipeline any) ([]M, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var ms []M
	err = cursor.All(ctx, &ms)
	return ms, err
}

// AggregateCount performs an aggregation operation on the collection using the provided pipeline and returns the count.
func (r *BaseRepo[M, N]) AggregateCount(ctx context.Context, pipeline any) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	cursor, err := r.Collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		Count int `bson:"count"`
	}
	err = cursor.All(ctx, &result)
	if err != nil {
		return 0, err
	}

	if len(result) == 0 {
		return 0, nil
	}

	return result[0].Count, nil
}
