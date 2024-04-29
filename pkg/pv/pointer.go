package pv

import (
	"reflect"
)

func Pointer[T any](v T) *T {
	valueOf := reflect.ValueOf(v)
	if valueOf.IsZero() {
		return nil
	}

	return &v
}
