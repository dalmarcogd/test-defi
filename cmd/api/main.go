package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/redirect"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/dalmarcogd/test-defi/cmd/api/envs"
	"github.com/dalmarcogd/test-defi/pkg/database"
	"github.com/dalmarcogd/test-defi/pkg/http/recovermdw"
	"github.com/dalmarcogd/test-defi/pkg/tracer"
	tracerbun "github.com/dalmarcogd/test-defi/pkg/tracer/instrumentation/bun"
	tracerfiber "github.com/dalmarcogd/test-defi/pkg/tracer/instrumentation/fiber"
	"github.com/dalmarcogd/test-defi/pkg/validators"
	"github.com/dalmarcogd/test-defi/pkg/zapctx"
)

func main() {
	err := zapctx.StartZapctx()
	if err != nil {
		log.Fatal(err)
	}

	module := fx.Options(
		// Infra
		fx.Provide(
			validators.Setup,
			envs.New,
			setupTracer,
			setupDatabase,
		),
		// Chats
		fx.Provide(),
		fx.Invoke(runHTTPServer),
	)

	app := fx.New(module, fx.NopLogger)
	err = app.Err()
	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}

func setupTracer(lc fx.Lifecycle, env envs.Envs) (tracer.Tracer, error) {
	trace, err := tracer.New(tracer.Config{
		Endpoint:    env.OtelCollectorHost,
		ServiceName: env.Service,
		Env:         env.Environment,
		Version:     env.Version,
	})
	if err != nil {
		return trace, err
	}

	lc.Append(fx.Hook{
		OnStop: trace.Stop,
	})

	return trace, nil
}

func setupDatabase(lc fx.Lifecycle, env envs.Envs, trc tracer.Tracer) (
	database.Database[database.PostgresConfig, database.Master],
	database.Database[database.PostgresConfig, database.ReadReplica],
	error,
) {
	masterDatabase, err := database.New[database.PostgresConfig, database.Master](
		database.PostgresConfig{
			URI:   env.MasterDatabaseURL,
			Hooks: []database.Hook{tracerbun.NewTracingHook(trc)},
		},
		database.MasterInstanceType,
	)
	if err != nil {
		return nil, nil, err
	}

	lc.Append(fx.Hook{
		OnStop: masterDatabase.Stop,
	})

	readReplicaDatabase, err := database.New[database.PostgresConfig, database.ReadReplica](
		database.PostgresConfig{
			URI:   env.ReadReplicaDatabaseURL,
			Hooks: []database.Hook{tracerbun.NewTracingHook(trc)},
		},
		database.ReadReplicaInstanceType,
	)
	if err != nil {
		return nil, nil, err
	}

	lc.Append(fx.Hook{
		OnStop: readReplicaDatabase.Stop,
	})

	return masterDatabase, readReplicaDatabase, nil
}

func runHTTPServer(
	lc fx.Lifecycle,
	env envs.Envs,
	trc tracer.Tracer,
) error {
	app := fiber.New(fiber.Config{
		AppName:               env.Service,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		DisableStartupMessage: true,
	})
	app.Use(
		recovermdw.NewMiddleware(),
		cors.New(cors.ConfigDefault),
		swagger.New(swagger.Config{
			BasePath: "/api",
			FilePath: "./docs/swagger.yaml",
			Path:     "swagger",
			Title:    "Rules Engine API Docs",
		}),
		redirect.New(redirect.Config{
			Rules: map[string]string{
				"/": "/api/swagger",
			},
			StatusCode: http.StatusMovedPermanently,
		}),
		tracerfiber.NewMiddleware(trc),
	)

	if env.DebugPprof {
		app.Use(pprof.New())
	}

	// Match request starting with /api
	api := app.Group("api")

	_ = api.Group(
		"/v1",
		func(ctx *fiber.Ctx) error {
			ctx.Set("Version", "v1")
			return ctx.Next()
		},
	)

	//appV1.Post("/authenticate", authenticateHandler)
	//appV1.Post("/tenants", createTenantsHandler)
	//appV1.Get("/tenants/:id", getTenantsHandler)
	//appV1.Get("/tenants", listTenantsHandler)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				host := fmt.Sprintf(":%s", env.HTTPPort)
				zap.L().Info(
					"http_server_up",
					zap.String("description", "up and running api server"),
					zap.String("host", host),
				)

				err := app.Listen(host)
				if err != nil {
					zapctx.L(ctx).Fatal("http_server_up", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.ShutdownWithContext(ctx)
		},
	})

	return nil
}
