package tracerbun

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/dalmarcogd/test-defi/pkg/tracer"
)

// queryPingToIgnore is a query to ping database in health check processes because of this we are
// ignoring it at the tracing level (e.g. Datadog).
const queryPingToIgnore = "SELECT 1"

type tracingHook struct {
	tracer tracer.Tracer
}

// NewTracingHook returns a new implementation QueryHook to span queries.
func NewTracingHook(tracer tracer.Tracer) bun.QueryHook {
	return tracingHook{tracer: tracer}
}

// BeforeQuery implements bun.QueryHook.
func (t tracingHook) BeforeQuery(ctx context.Context, qe *bun.QueryEvent) context.Context {
	if isPingQuery(qe) {
		return ctx
	}

	query := qe.Query
	if query == "" {
		query = "unknown"
	}

	opts := []tracer.Attributes{
		semconv.DBSystemPostgreSQL,
		semconv.DBStatementKey.String(query),
		semconv.DBNameKey.String(qe.DB.String()),
	}

	ctx, _ = t.tracer.SpanName(ctx, fmt.Sprintf("%s.query", qe.DB.Dialect().Name()), opts...)
	return ctx
}

// AfterQuery implements bun.QueryHook.
func (t tracingHook) AfterQuery(ctx context.Context, qe *bun.QueryEvent) {
	if isPingQuery(qe) {
		return
	}

	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		defer span.End()
		span.End()

		if qe.Err != nil && !errors.Is(qe.Err, sql.ErrNoRows) {
			span.RecordError(qe.Err)
		}
	}
}

// Helpers

// isPingQuery returns if bun.QueryEvent present a query to ping database.
func isPingQuery(qe *bun.QueryEvent) bool {
	return qe.Query == queryPingToIgnore
}
