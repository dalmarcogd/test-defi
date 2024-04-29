//go:build unit

package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func TestDatabaseLogger(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	var beforeCalled bool
	before := func(ctx context.Context, q *bun.QueryEvent) {
		beforeCalled = true
	}

	var afterCalled bool
	after := func(ctx context.Context, q *bun.QueryEvent) {
		afterCalled = true
	}

	databaseLogger := NewDatabaseLogger(before, after)

	_ = databaseLogger.BeforeQuery(ctx, &bun.QueryEvent{})
	databaseLogger.AfterQuery(ctx, &bun.QueryEvent{})

	assert.True(t, beforeCalled, "beforeFunction was not called")
	assert.True(t, afterCalled, "afterFunction was not called")
}
