package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var Sgtin validator.Func = func(fl validator.FieldLevel) bool {
	result, _ := regexp.MatchString("\\d{14}[\\da-zA-Z]{13}", fl.Field().String())
	return result
}
