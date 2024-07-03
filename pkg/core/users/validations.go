package users

import (
	"tenkhours/pkg/db/coredb"

	"github.com/go-playground/validator/v10"
)

func ValidateUser(user coredb.User) error {
	validate := validator.New()

	return validate.Struct(user)
}
