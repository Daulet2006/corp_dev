package security

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
)

// Sanitizer: use an allowlist policy for text fields

var Sanitizer = bluemonday.StrictPolicy()

// RegisterCommonValidators registers custom validators like strongpass.

func RegisterCommonValidators(v *validator.Validate) {
	// strong password: min 8, at least one upper, one lower, one digit, one special
	v.RegisterValidation("strongpass", func(fl validator.FieldLevel) bool {
		pw := fl.Field().String()
		if len(pw) < 8 {
			return false
		}
		var (
			hasUpper = regexp.MustCompile(`[A-Z]`).MatchString
			hasLower = regexp.MustCompile(`[a-z]`).MatchString
			hasDigit = regexp.MustCompile(`\d`).MatchString
			hasSpec  = regexp.MustCompile(`[\W_]`).MatchString
		)
		return hasUpper(pw) && hasLower(pw) && hasDigit(pw) && hasSpec(pw)
	})

	// optional: normalize emails before validation
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}
