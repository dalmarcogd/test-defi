//go:build unit

package tracer

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	noop2 "go.opentelemetry.io/otel/trace/noop"
)

func TestSpans(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	serviceImpl, err := New(Config{
		Endpoint:    "localhost:8126",
		ServiceName: "",
		Env:         "",
		Version:     "",
	})
	assert.NoError(t, err)
	otel.SetTracerProvider(noop2.NewTracerProvider())

	_, s := serviceImpl.Span(context.Background())
	s.End()
}
