package pv

import "reflect"

func Valuer[T any](p *T) T {
	var t T
	valueOf := reflect.ValueOf(p)
	if valueOf.IsNil() {
		return t
	}
	if !valueOf.Elem().IsZero() {
		t = *p
	}

	return t
}
