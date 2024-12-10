package validator

import (
	"gin-boot-starter/core/enums"

	"github.com/go-playground/validator/v10"
)

func GenderValidator(fl validator.FieldLevel) bool {
	// value of the field
	val := fl.Field().String()
	return enums.ParseGender(val) != enums.Unknown
}
