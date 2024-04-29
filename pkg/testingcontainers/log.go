package testingcontainers

import (
	"testing"
)

type Logging interface {
	Printf(format string, v ...interface{})
}

type logger struct {
	t *testing.T
}

func NewLogging(t *testing.T) Logging {
	return logger{t: t}
}

func (l logger) Printf(format string, v ...interface{}) {
	l.t.Logf(format, v...)
}
