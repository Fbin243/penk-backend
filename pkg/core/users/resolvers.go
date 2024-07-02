package users

import (
	"context"
	"fmt"
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

func NewUsersResolver(usersRepo *coredb.UsersRepo) *UsersResolver {
	return &UsersResolver{
		UsersRepo: usersRepo,
	}
}

func (r *UsersResolver) RegisterAccount(params graphql.ResolveParams) (interface{}, error) {
	authProfile, err := auth.GetProfileByContext(params.Context)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return user.ID.Hex(), nil
}

func (r *UsersResolver) GetUserByToken(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}
