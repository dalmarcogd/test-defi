package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"

	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/schema"
)

type (
	Hook                                bun.QueryHook
	Database[C Configs, I InstanceType] interface {
		Client() *bun.DB
		InstanceType() I
		Config() C
		Stop(ctx context.Context) error
	}

	database[C Configs, I InstanceType] struct {
		db  *bun.DB
		cfg C
		it  I
	}
)

type (
	Configs interface {
		any | PostgresConfig | SQLiteConfig
	}

	PostgresConfig struct {
		// URI e.g: postgres://pass:user@localhost:5432/my-database?sslmode=disable
		URI   string
		Hooks []Hook
	}

	SQLiteConfig struct {
		// URI e.g: file::memory:?cache=shared
		URI   string
		Hooks []Hook
	}

	InstanceType interface {
		any | Master | ReadReplica
	}

	Master      struct{}
	ReadReplica struct{}
)

var (
	MasterInstanceType      Master
	ReadReplicaInstanceType ReadReplica

	Nil = database[any, any]{}
)

// New this function initialize a wrapper client over the database connection client.
//
// We are currently working with github.com/uptrace/bun and it can support:
//   - PostgreSQL,
//   - MySQL,
//   - MSSQL,
//   - SQLite.
func New[T Configs, I InstanceType](cfg T, it I) (Database[T, I], error) {
	var newDB *sql.DB
	var dialect schema.Dialect
	var hooks []Hook

	switch c := any(cfg).(type) {
	case PostgresConfig:
		connector, err := pgdriver.NewDriver().OpenConnector(c.URI)
		if err != nil {
			return nil, err
		}

		newDB = sql.OpenDB(connector)
		newDB.SetMaxOpenConns(1)
		dialect = pgdialect.New()
		hooks = c.Hooks
	case SQLiteConfig:
		d, err := sql.Open(sqliteshim.ShimName, c.URI)
		if err != nil {
			return nil, err
		}

		newDB = d
		dialect = sqlitedialect.New()
		hooks = c.Hooks
	default:
		return nil, errors.New("invalid structure configuration")
	}

	db := bun.NewDB(newDB, dialect, bun.WithDiscardUnknownColumns())
	for _, hook := range hooks {
		db.AddQueryHook(hook)
	}

	return database[T, I]{db: db, cfg: cfg, it: it}, nil
}

func (m database[C, I]) Client() *bun.DB {
	return m.db
}

func (m database[C, I]) Config() C {
	return m.cfg
}

func (m database[C, I]) InstanceType() I {
	return m.it
}

func (m database[C, I]) Stop(_ context.Context) error {
	return m.db.Close()
}
