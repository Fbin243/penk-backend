package mongodb

import (
	"context"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// InsertMany inserts multiple documents into the collection.
func (r *BaseRepo[M, N]) InsertMany(ctx context.Context, ms []M) ([]M, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	docs := make([]any, len(ms))
	for i, m := range ms {
		addMissingFields(m, r.WithTimestamp)
		docs[i] = ToMongoEntity[M, N](&m)
	}

	var err error
	if len(docs) > 0 {
		_, err = r.Collection.InsertMany(ctx, docs)
	}

	return ms, err
}

// InsertOne inserts a single document into the collection.
func (r *BaseRepo[M, N]) InsertOne(ctx context.Context, m *M) (*M, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	addMissingFields(*m, r.WithTimestamp)
	_, err := r.Collection.InsertOne(ctx, ToMongoEntity[M, N](m))
	return m, err
}

// FindAndUpdateByID finds a document by ID and updates it with the provided data.
func (r *BaseRepo[M, N]) FindAndUpdateByID(ctx context.Context, id string, m *M) (*M, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	if r.WithTimestamp {
		lo.FromPtr(m).SetUpdatedAtByNow()
	}

	err := r.Collection.FindOneAndUpdate(ctx,
		bson.M{"_id": ToObjectID(id)},
		bson.M{"$set": ToMongoEntity[M, N](m)},
		FindOneAndUpdateOptions).Decode(&m)
	return m, err
}

// FindAndUpdateByIDs finds multiple documents by their IDs and updates them with the provided data.
func (r *BaseRepo[M, N]) FindAndUpdateByIDs(ctx context.Context, ms []M) ([]M, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	models := []mongo.WriteModel{}
	for _, m := range ms {
		if r.WithTimestamp {
			m.SetUpdatedAtByNow()
		}

		models = append(models, mongo.NewUpdateOneModel().
			SetFilter(bson.M{"_id": ToObjectID(m.GetID())}).
			SetUpdate(bson.M{"$set": ToMongoEntity[M, N](&m)}))
	}

	var err error
	if len(models) > 0 {
		_, err = r.Collection.BulkWrite(ctx, models)
	}

	return ms, err
}

// UpdateByID updates a single document in the collection based on the provided ID.
func (r *BaseRepo[M, N]) UpdateByID(ctx context.Context, id string, m *M) error {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	if r.WithTimestamp {
		lo.FromPtr(m).SetUpdatedAtByNow()
	}

	_, err := r.Collection.UpdateOne(ctx,
		bson.M{"_id": ToObjectID(id)},
		bson.M{"$set": ToMongoEntity[M, N](m)},
	)
	return err
}

// UpdateOne updates a single document in the collection based on the provided filter.
func (r *BaseRepo[M, N]) UpdateOne(ctx context.Context, filter, update any) error {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

// UpdateMany updates multiple documents in the collection based on the provided filter.
func (r *BaseRepo[M, N]) UpdateMany(ctx context.Context, filter, update any) error {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	_, err := r.Collection.UpdateMany(ctx, filter, update)
	return err
}

// FindAndDeleteByID finds a document by ID and deletes it.
func (r *BaseRepo[M, N]) FindAndDeleteByID(ctx context.Context, id string) (*M, error) {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	var m M
	err := r.Collection.FindOneAndDelete(ctx, bson.M{"_id": ToObjectID(id)}).Decode(&m)
	return &m, err
}

// DeleteByID deletes a document by ID.
func (r *BaseRepo[M, N]) DeleteByID(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	_, err := r.Collection.DeleteOne(ctx, bson.M{"_id": ToObjectID(id)})
	return err
}

// DeleteOne deletes a single document from the collection based on the provided filter.
func (r *BaseRepo[M, N]) DeleteOne(ctx context.Context, filter any) error {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	_, err := r.Collection.DeleteOne(ctx, filter)
	return err
}

// DeleteMany deletes multiple documents from the collection based on the provided filter.
func (r *BaseRepo[M, N]) DeleteMany(ctx context.Context, filter any) error {
	ctx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	_, err := r.Collection.DeleteMany(ctx, filter)
	return err
}
