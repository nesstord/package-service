package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

const GtinPattern = "\\d{14}"

var GtinRegexp = regexp.MustCompile(GtinPattern)

var Gtin validator.Func = func(fl validator.FieldLevel) bool {
	result, _ := regexp.MatchString(GtinPattern, fl.Field().String())
	return result
}
