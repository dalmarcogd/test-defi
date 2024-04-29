package validators

import (
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var CNPJRegexp = regexp.MustCompile(`^\d{2}\.?\d{3}\.?\d{3}\/?(:?\d{3}[1-9]|\d{2}[1-9]\d|\d[1-9]\d{2}|[1-9]\d{3})-?\d{2}$`)

const (
	cnpjSize = 12
	cnpjPos  = 5
)

func CNPJ(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		return IsCNPJ(field.String())
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

// IsCNPJ verifies if the given string is a valid CPF document.
func IsCNPJ(doc string) bool {
	if !CNPJRegexp.MatchString(doc) {
		return false
	}

	doc = cleanNonDigits(doc)

	// Invalidates documents with all digits equal.
	if allEq(doc) {
		return false
	}

	d := doc[:cnpjSize]
	digit := calculateDigit(d, cnpjPos)

	d += digit
	digit = calculateDigit(d, cnpjPos+1)

	return doc == d+digit
}
