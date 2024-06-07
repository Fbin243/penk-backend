package users

import (
	"context"
	"log"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"

	"tenkhours/pkg/db/usersdb"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func registerAccount(params graphql.ResolveParams) (interface{}, error) {
	authProfile, ok := params.Context.Value(auth.ProfileContextKey).(auth.Profile)
	if !ok {
		return nil, auth.ErrorProfileNotFound
	}

	user := usersdb.User{
		ID:          primitive.NewObjectID(),
		Name:        authProfile.Name,
		Email:       authProfile.Email,
		FirebaseUID: authProfile.UID,
		ImageURL:    "avatar.jpg", // default avatar URL
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := db.GetDB().Collection("users").InsertOne(ctx, user)
	if err != nil {
		log.Printf("Failed to insert user: %v\n", err)
		return nil, err
	}

	return user, nil
}
