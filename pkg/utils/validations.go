package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ValidateBuilder struct {
	validateString string
}

// NewValidateBuilder creates a new ValidateBuilder
func NewValidateBuilder() *ValidateBuilder {
	return &ValidateBuilder{
		validateString: "",
	}
}

// Required adds the required tag to the validation string
func (v *ValidateBuilder) Required() *ValidateBuilder {
	v.validateString += "required,"
	return v
}

// OmitEmpty adds the omitempty tag to the validation string
func (v *ValidateBuilder) OmitEmpty() *ValidateBuilder {
	v.validateString += "omitempty,"
	return v
}

// Min adds the min tag to the validation string
func (v *ValidateBuilder) Min(min int) *ValidateBuilder {
	v.validateString += fmt.Sprintf("min=%d,", min)
	return v
}

// Max adds the max tag to the validation string
func (v *ValidateBuilder) Max(max int) *ValidateBuilder {
	v.validateString += fmt.Sprintf("max=%d,", max)
	return v
}

// Email adds the email tag to the validation string
func (v *ValidateBuilder) Email() *ValidateBuilder {
	v.validateString += "email,"
	return v
}

// Color adds the hexcolor tag to the validation string
func (v *ValidateBuilder) Color() *ValidateBuilder {
	v.validateString += "hexcolor,"
	return v
}

// Custom adds a custom tag to the validation string and registers the custom function
func (v *ValidateBuilder) Custom(tag string, fn validator.Func) *ValidateBuilder {
	validate.RegisterValidation(tag, fn)
	v.validateString += tag + ","
	return v
}

// Reset resets the validation string
func (v *ValidateBuilder) Reset() *ValidateBuilder {
	v.validateString = ""
	return v
}

// Validate validates the field with the validation string
func (v *ValidateBuilder) Validate(field interface{}) error {
	return validate.Var(field, v.validateString[:len(v.validateString)-1])
}
