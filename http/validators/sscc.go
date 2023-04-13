package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

const SsccPattern = "\\d{18}"

var SsccRegexp = regexp.MustCompile(SsccPattern)

var Sscc validator.Func = func(fl validator.FieldLevel) bool {
	result, _ := regexp.MatchString(SsccPattern, fl.Field().String())
	return result
}
