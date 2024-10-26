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

func NewValidateBuilder() *ValidateBuilder {
	return &ValidateBuilder{
		validateString: "",
	}
}

func (v *ValidateBuilder) Required() *ValidateBuilder {
	v.validateString += "required,"
	return v
}

func (v *ValidateBuilder) OmitEmpty() *ValidateBuilder {
	v.validateString += "omitempty,"
	return v
}

func (v *ValidateBuilder) Min(min int) *ValidateBuilder {
	v.validateString += fmt.Sprintf("min=%d,", min)
	return v
}

func (v *ValidateBuilder) Max(max int) *ValidateBuilder {
	v.validateString += fmt.Sprintf("max=%d,", max)
	return v
}

func (v *ValidateBuilder) Email() *ValidateBuilder {
	v.validateString += "email,"
	return v
}

func (v *ValidateBuilder) Color() *ValidateBuilder {
	v.validateString += "hexcolor,"
	return v
}

func (v *ValidateBuilder) AddCustomValidate(tag string, fn validator.Func) *ValidateBuilder {
	validate.RegisterValidation(tag, fn)
	return v
}

func (v *ValidateBuilder) Custom(tag string) *ValidateBuilder {
	v.validateString += tag + ","
	return v
}

func (v *ValidateBuilder) Validate(field interface{}) error {
	v.validateString = v.validateString[:len(v.validateString)-1]
	return validate.Var(field, v.validateString)
}
