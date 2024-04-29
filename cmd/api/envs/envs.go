package envs

import "github.com/gosidekick/goconfig"

// Envs this object keep the all command environment variables.
type Envs struct {
	// Application
	Environment string `cfg:"ENVIRONMENT" cfgRequired:"true" cfgDefault:"local"`
	Service     string `cfg:"SERVICE" cfgRequired:"true" cfgDefault:"one-traive-api-command"`
	Version     string `cfg:"VERSION" cfgRequired:"true" cfgDefault:"local-snapshot"`
	HTTPHost    string `cfg:"HTTP_HOST" cfgRequired:"true" cfgDefault:"localhost"`
	HTTPPort    string `cfg:"HTTP_PORT" cfgRequired:"true" cfgDefault:"8080"`
	DebugPprof  bool   `cfg:"DEBUG_PPROF" cfgDefault:"true"`
	// Database
	MasterDatabaseURL      string `cfg:"POSTGRES_DATABASE_PRIMARY_URL" cfgRequired:"true"`
	ReadReplicaDatabaseURL string `cfg:"POSTGRES_DATABASE_REPLICA_URL" cfgRequired:"true"`
	// Open Telemetry
	OtelCollectorHost string `cfg:"OTEL_COLLECTOR_HOST" cfgRequired:"true"`
}

func New() (Envs, error) {
	var env Envs
	err := goconfig.Parse(&env)
	return env, err
}
