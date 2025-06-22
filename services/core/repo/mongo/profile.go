package mongorepo

import (
	"context"
	"log"

	"tenkhours/services/core/entity"
	mongomodel "tenkhours/services/core/repo/mongo/model"

	mongodb "tenkhours/pkg/db/mongo"
	"tenkhours/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProfileRepo struct {
	*mongodb.BaseRepo[entity.Profile, mongomodel.Profile]
}

func NewProfileRepo(db *mongo.Database) *ProfileRepo {
	profilesCollection := db.Collection(mongodb.ProfilesCollection)
	_, err := profilesCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "firebase_uid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		log.Printf("failed to create indexes for %s collection\n", mongodb.ProfilesCollection)
		return nil
	}

	return &ProfileRepo{
		mongodb.NewBaseRepo[entity.Profile, mongomodel.Profile](
			profilesCollection,
			true),
	}
}

func (r *ProfileRepo) GetProfileByFirebaseUID(ctx context.Context, firebaseUID string) (*entity.Profile, error) {
	profile, err := r.FindOne(ctx, bson.M{"firebase_uid": firebaseUID})
	if err == mongo.ErrNoDocuments {
		return nil, errors.ErrMongoNotFound
	}

	return profile, err
}

func (r *ProfileRepo) DeleteProfileByFirebaseUID(ctx context.Context, firebaseUID string) error {
	return r.DeleteOne(ctx, bson.M{"firebase_uid": firebaseUID})
}
