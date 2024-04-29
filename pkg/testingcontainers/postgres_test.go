//go:build integration

package testingcontainers

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestPostgresContainer(t *testing.T) {
	ctx := context.Background()
	container, err := NewPostgresContainer(ctx, NewLogging(t), PostgresConfig{
		DatabaseName: "entities",
		Port:         "5432",
	})
	defer func() {
		if e := container.Cleanup(ctx); e != nil {
			assert.NoError(t, e)
		}
	}()
	assert.NoError(t, err)

	db, err := sql.Open("postgres", container.URI)
	assert.NoError(t, err)

	assert.NoError(t, db.PingContext(ctx))
}
