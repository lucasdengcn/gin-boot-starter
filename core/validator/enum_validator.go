package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func EnumValidator(fl validator.FieldLevel) bool {
	// validation on status
	// expect values. eg. enum=on off
	matches := strings.Split(fl.Param(), " ")
	// value of the field
	str := fl.Field().String()
	// log.Debug().Msgf("Field value is: %v", str)
	//
	for _, s := range matches {
		if s == str {
			return true
		}
	}
	//
	return false
}
