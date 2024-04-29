package validators

import (
	"reflect"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func NotSpace(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		for _, v := range field.String() {
			if unicode.IsSpace(v) {
				return false
			}
		}

		return true
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}
