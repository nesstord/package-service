package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var Gtin validator.Func = func(fl validator.FieldLevel) bool {
	result, _ := regexp.MatchString("\\d{14}", fl.Field().String())
	return result
}
