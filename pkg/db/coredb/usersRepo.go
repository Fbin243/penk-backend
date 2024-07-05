package coredb

import (
	"context"
	"log"
	"time"

	"tenkhours/pkg/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UsersRepo struct {
	*mongo.Collection
}

func NewUsersRepo() *UsersRepo {
	usersCollection := db.GetUsersCollection()
	_, err := usersCollection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
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
		log.Println("failed to create indexes for users collection")
		return nil
	}

	return &UsersRepo{usersCollection}
}

func (r *UsersRepo) GetUserByFirebaseUID(firebaseUID string) (*User, error) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.FindOne(ctx, bson.M{"firebase_uid": firebaseUID}).Decode(&user)

	return &user, err
}

func (r *UsersRepo) CreateNewUser(user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.InsertOne(ctx, user)

	return user, err
}

func (r *UsersRepo) UpdateUserByID(id primitive.ObjectID, user *User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.FindOneAndUpdate(ctx, bson.M{"_id": id}, bson.M{"$set": user}, db.FindOneAndUpdateOptions).Decode(user)

	return user, err
}
