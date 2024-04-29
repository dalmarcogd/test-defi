package validators

import (
	"bytes"
	"reflect"
	"regexp"
	"strconv"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var CPFRegexp = regexp.MustCompile(`^\d{3}\.?\d{3}\.?\d{3}-?\d{2}$`)

const (
	cpfSize = 9
	cpfPos  = 10
)

func CPF(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		return IsCPF(field.String())
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

// IsCPF verifies if the given string is a valid CPF document.
func IsCPF(doc string) bool {
	if !CPFRegexp.MatchString(doc) {
		return false
	}

	doc = cleanNonDigits(doc)

	// Invalidates documents with all digits equal.
	if allEq(doc) {
		return false
	}

	d := doc[:cpfSize]
	digit := calculateDigit(d, cpfPos)

	d += digit
	digit = calculateDigit(d, cpfPos+1)

	return doc == d+digit
}

// cleanNonDigits removes every rune that is not a digit.
func cleanNonDigits(doc string) string {
	buf := bytes.NewBufferString("")
	for _, r := range doc {
		if unicode.IsDigit(r) {
			buf.WriteRune(r)
		}
	}

	return buf.String()
}

// allEq checks if every rune in a given string is equal.
func allEq(doc string) bool {
	base := doc[0]
	for i := 1; i < len(doc); i++ {
		if base != doc[i] {
			return false
		}
	}

	return true
}

// calculateDigit calculates the next digit for the given document.
func calculateDigit(doc string, position int) string {
	var sum int
	for _, r := range doc {
		sum += int(r-'0') * position
		position--

		if position < 2 {
			position = 9
		}
	}

	sum %= 11
	if sum < 2 {
		return "0"
	}

	return strconv.Itoa(11 - sum)
}
