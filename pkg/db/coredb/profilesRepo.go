package coredb

import (
	"context"
	"log"
	"time"

	"tenkhours/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProfilesRepo struct {
	*mongo.Collection
}

func NewProfilesRepo(mongodb *mongo.Database) *ProfilesRepo {
	profilesCollection := mongodb.Collection(db.ProfileCollection)
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
		log.Println("failed to create indexes for profiles collection")
		return nil
	}

	return &ProfilesRepo{profilesCollection}
}

func (r *ProfilesRepo) GetProfileByFirebaseUID(firebaseUID string) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var profile Profile
	err := r.FindOne(ctx, bson.M{"firebase_uid": firebaseUID}).Decode(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *ProfilesRepo) GetProfileByEmail(email string) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var profile Profile
	err := r.FindOne(ctx, bson.M{"email": email}).Decode(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *ProfilesRepo) CreateNewProfile(profile *Profile) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.InsertOne(ctx, profile)

	return profile, err
}

func (r *ProfilesRepo) UpdateProfile(profile *Profile) (*Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.FindOneAndUpdate(ctx, bson.M{"_id": profile.ID}, bson.M{"$set": profile}, db.FindOneAndUpdateOptions).Decode(profile)

	return profile, err
}

func (r *ProfilesRepo) DeleteProfileByFirebaseUID(firebaseUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.DeleteOne(ctx, bson.M{"firebase_uid": firebaseUID})
	return err
}
