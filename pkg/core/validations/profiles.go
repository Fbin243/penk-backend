package validations

import (
	"tenkhours/pkg/db/coredb"

	"github.com/go-playground/validator/v10"
)

func ValidateProfile(profile coredb.Profile) error {
	validate := validator.New()

	return validate.Struct(profile)
}
