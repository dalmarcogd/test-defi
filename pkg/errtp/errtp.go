package errtp

import (
	"errors"
	"fmt"
)

type (
	ErrorTp struct {
		statusCode int
		err        error
		Message    string   `json:"type"`
		Details    []string `json:"details"`
	}

	Option func(e ErrorTp) ErrorTp
)

func (e ErrorTp) Is(err error) bool {
	return errors.Is(err, e.err)
}

func (e ErrorTp) Unwrap() error {
	return e.err
}

func (e ErrorTp) Error() string {
	return e.err.Error()
}

func Wrap(err error, opts ...Option) ErrorTp {
	var errTp ErrorTp
	if !errors.As(err, &errTp) {
		errTp = ErrorTp{
			statusCode: 0,
			err:        nil,
			Message:    err.Error(),
			Details:    nil,
		}
	}

	for _, opt := range opts {
		errTp = opt(errTp)
	}

	return errTp
}

func WithStatusCode(st int) Option {
	return func(e ErrorTp) ErrorTp {
		e.statusCode = st
		return e
	}
}

func WithStatusCodeIfNoValue(st int) Option {
	return func(e ErrorTp) ErrorTp {
		if e.statusCode == 0 {
			e.statusCode = st
		}

		return e
	}
}

func WithMessage(msg string) Option {
	return func(e ErrorTp) ErrorTp {
		e.Message = msg
		return e
	}
}

func WithDetails(dt ...string) Option {
	return func(e ErrorTp) ErrorTp {
		e.Details = dt
		return e
	}
}

func FieldDetail(f string) string {
	return fmt.Sprintf("field: %s", f)
}

func OriginalError(f error) string {
	return fmt.Sprintf("original error: %s", f)
}

func StatusCode(err error) int {
	var errTp ErrorTp
	if errors.As(err, &errTp) {
		return errTp.statusCode
	}

	return 0
}
