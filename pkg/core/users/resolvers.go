package users

import (
	"context"
	"log"
	"time"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db"
	"tenkhours/pkg/db/coredb"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func registerAccount(params graphql.ResolveParams) (interface{}, error) {
	authProfile, err := auth.GetProfileByContext(params.Context)
	if err != nil {
		return nil, err
	}

	user := coredb.User{
		ID:          primitive.NewObjectID(),
		Name:        authProfile.Name,
		Email:       authProfile.Email,
		FirebaseUID: authProfile.UID,
		ImageURL:    "", // default avatar URL
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.GetUsersCollection().InsertOne(ctx, user)
	if err != nil {
		log.Printf("Failed to insert user: %v\n", err)
		return nil, err
	}

	return user, nil
}

func getUserByEmail(params graphql.ResolveParams) (interface{}, error) {
	email, ok := params.Args["email"].(string)
	if !ok {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user coredb.User
	err := db.GetUsersCollection().FindOne(ctx, map[string]string{
		"email": email,
	}).Decode(&user)
	if err != nil {
		log.Printf("Failed to find user: %v\n", err)
		return nil, err
	}

	return user, nil
}
