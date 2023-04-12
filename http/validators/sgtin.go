package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

const SgtinPattern = "\\d{14}[\\da-zA-Z]{13}"

var Sgtin validator.Func = func(fl validator.FieldLevel) bool {
	result, _ := regexp.MatchString(SgtinPattern, fl.Field().String())
	return result
}
