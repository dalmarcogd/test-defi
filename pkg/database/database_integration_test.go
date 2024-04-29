//go:build integration

package database

import (
	"context"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"

	"github.com/dalmarcogd/test-defi/pkg/testingcontainers"
)

func TestDatabase(t *testing.T) {
	t.Parallel()

	assert.NoError(t, os.Setenv("SHOW_DATABASE_QUERIES", "true"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	t.Run("Invalid url", func(t *testing.T) {
		_, err := New(PostgresConfig{URI: gofakeit.URL()}, MasterInstanceType)
		assert.Error(t, err)
	})

	t.Run("Ping invalid postgres database", func(t *testing.T) {
		db, err := New(PostgresConfig{URI: gofakeit.URL()}, MasterInstanceType)
		assert.Contains(t, err.Error(), "pgdriver: invalid scheme: ")
		assert.Empty(t, db)
	})

	t.Run("Ping valid postgres database", func(t *testing.T) {
		container, err := testingcontainers.NewPostgresContainer(
			ctx,
			testingcontainers.NewLogging(t),
			testingcontainers.PostgresConfig{
				DatabaseName: "test",
				Port:         "5432",
			},
		)
		defer func() {
			assert.NoError(t, container.Cleanup(ctx))
		}()
		assert.NoError(t, err)

		cfg := PostgresConfig{
			URI: container.URI,
			Hooks: []Hook{
				NewDatabaseLogger(nil, nil),
			},
		}
		db, err := New(cfg, MasterInstanceType)
		assert.NoError(t, err)
		defer func() {
			assert.NoError(t, db.Stop(ctx))
		}()

		assert.NoError(t, db.Client().PingContext(ctx))
		assert.Equal(t, cfg, db.Config())
	})

	t.Run("Ping valid sqlite database", func(t *testing.T) {
		cfg := SQLiteConfig{
			URI: "file::memory:?cache=shared",
			Hooks: []Hook{
				NewDatabaseLogger(nil, nil),
			},
		}
		db, err := New(cfg, MasterInstanceType)
		assert.NoError(t, err)
		defer func() {
			assert.NoError(t, db.Stop(ctx))
		}()

		assert.NoError(t, db.Client().PingContext(ctx))
		assert.Equal(t, cfg, db.Config())
	})
}
