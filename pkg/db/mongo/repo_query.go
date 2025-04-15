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
