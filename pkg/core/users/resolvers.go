package users

import (
	"log"

	"tenkhours/pkg/auth"
	"tenkhours/pkg/db/coredb"
	"tenkhours/pkg/utils"

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

func (r *UsersResolver) GetUserByToken(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	return user, nil
}

func (r *UsersResolver) UpdateAccount(params graphql.ResolveParams) (interface{}, error) {
	user, ok := params.Context.Value(auth.UserKey).(coredb.User)
	if !ok {
		return nil, auth.ErrorUnauthorized
	}

	input := params.Args["input"].(map[string]interface{})
	if name, ok := input["name"].(string); ok {
		user.Name = name
	}

	if imageURL, ok := input["imageURL"].(string); ok {
		user.ImageURL = imageURL
	}

	if currentCharacterID, ok := input["currentCharacterID"].(string); ok {
		currentCharacterOID, err := primitive.ObjectIDFromHex(currentCharacterID)
		if err != nil {
			return nil, err
		}

		user.CurrentCharacterID = currentCharacterOID
	}

	if autoSnapshot, ok := input["autoSnapshot"].(bool); ok {
		user.AutoSnapshot = autoSnapshot
	}

	user.UpdatedAt = utils.Now()

	err := ValidateUser(user)
	if err != nil {
		return nil, err
	}

	updatedUser, err := r.UsersRepo.UpdateUser(&user)
	if err != nil {
		log.Printf("failed to update user: %v\n", err)
		return nil, err
	}

	return *updatedUser, nil
}
