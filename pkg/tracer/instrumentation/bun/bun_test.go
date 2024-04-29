//go:build unit

package tracerbun

import (
	"context"
	"database/sql"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"

	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"

	"github.com/dalmarcogd/test-defi/pkg/tracer"
)

type user struct {
	bun.BaseModel `bun:"table:myusers,alias:u"`
	Name          string `bun:"myname"`
}

func TestQueryHook(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	trc, err := tracer.New(
		tracer.Config{
			Endpoint:    "localhost:8080",
			Env:         "dev",
			Version:     "v0.0.1",
			ServiceName: "my-service",
		},
	)
	assert.NoError(t, err)
	defer func(trc tracer.Tracer, ctx context.Context) {
		err := trc.Stop(ctx)
		if err != nil {
			panic(err)
		}
	}(trc, ctx)

	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	assert.NoError(t, err)

	db := bun.NewDB(sqldb, sqlitedialect.New())

	_, err = db.NewCreateTable().Model(&user{}).Exec(ctx)
	assert.NoError(t, err)

	var i int
	err = db.NewRaw(queryPingToIgnore).Scan(ctx, &i)
	assert.NoError(t, err)
	assert.Equal(t, 1, i)

	err = db.NewRaw("select tod").Scan(ctx, &i)
	assert.EqualError(t, err, "SQL logic error: no such column: tod (1)")
}
