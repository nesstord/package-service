package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var sscc validator.Func = func(fl validator.FieldLevel) bool {
	result, _ := regexp.MatchString("\\d{18}", fl.Field().String())
	return result
}
