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

type UsersResolver struct {
	UsersRepo *coredb.UsersRepo
}

func NewUsersResolver() *UsersResolver {
	return &UsersResolver{
		UsersRepo: coredb.NewUsersRepo(),
	}
}

func (r *UsersResolver) RegisterAccount(params graphql.ResolveParams) (interface{}, error) {
	authProfile, err := auth.GetProfileByContext(params.Context)
	if err != nil {
		return false, err
	}

	user := coredb.User{
		ID:          primitive.NewObjectID(),
		Name:        authProfile.Name,
		Email:       authProfile.Email,
		FirebaseUID: authProfile.UID,
		ImageURL:    "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.GetUsersCollection().InsertOne(ctx, user)
	if err != nil {
		log.Printf("failed to insert user: %v\n", err)
		return false, err
	}

	return true, nil
}

func (r *UsersResolver) GetUserByEmail(params graphql.ResolveParams) (interface{}, error) {
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
		log.Printf("failed to find user: %v\n", err)
		return nil, err
	}

	return user, nil
}
