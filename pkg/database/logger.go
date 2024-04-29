package database

import (
	"context"

	"github.com/uptrace/bun"
)

type (
	logFunction func(ctx context.Context, q *bun.QueryEvent)

	dbLogger struct {
		before logFunction
		after  logFunction
	}
)

// NewDatabaseLogger returns a new implementation QueryHook to log SQL queries on console.
func NewDatabaseLogger(before, after logFunction) bun.QueryHook {
	return dbLogger{
		before: before,
		after:  after,
	}
}

func (d dbLogger) BeforeQuery(ctx context.Context, q *bun.QueryEvent) context.Context {
	if d.before != nil {
		d.before(ctx, q)
	}

	return ctx
}

func (d dbLogger) AfterQuery(ctx context.Context, q *bun.QueryEvent) {
	if d.after != nil {
		d.after(ctx, q)
	}
}
