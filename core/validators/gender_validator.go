package validators

import (
	"gin001/core/enums"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func GenderValidator(fl validator.FieldLevel) bool {
	// value of the field
	str := fl.Field().String()
	// log.Debug().Msgf("Gender is: %v", str)
	num, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	switch num {
	case int(enums.Male):
		return true
	case int(enums.Female):
		return true
	}
	return false
}
