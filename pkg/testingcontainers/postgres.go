package testingcontainers

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	// Load driver to read in file system the migrations. See: https://github.com/golang-migrate/migrate/tree/master/source/file
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type (
	PostgresDBContainer struct {
		testcontainers.Container
		URI      string
		UserName string
		Password string
		Port     string
		Config   PostgresConfig
		network  Network
		log      Logging
	}

	PostgresConfig struct {
		HostName     string
		DatabaseName string
		Network      Network
		Port         string
	}
)

func (container *PostgresDBContainer) Cleanup(ctx context.Context) error {
	err := container.Terminate(ctx)
	if err != nil {
		container.log.Printf("%s", err.Error())
	}

	if container.Config.Network.Name == "" {
		return container.network.Cleanup(ctx)
	}

	return nil
}

// NewPostgresContainer returns a postgres container for use in integration tests or other things
//
// Guidelines:
//
//   - container.Cleanup should always be called at the end of execution, this will shutdown the container and remove what was done inside it.
func NewPostgresContainer(
	ctx context.Context,
	log Logging,
	config PostgresConfig,
) (*PostgresDBContainer, error) {
	network := config.Network
	if network.Name == "" {
		var err error
		network, err = CreateNetwork(ctx)
		if err != nil {
			return nil, err
		}
	}

	templateURL := "postgres://%s:%s@%s:%s/%s?sslmode=disable"
	username := gofakeit.Username()
	password := gofakeit.Password(false, false, false, false, false, 15)
	portName := fmt.Sprintf("%s/tcp", config.Port)
	req := testcontainers.ContainerRequest{
		Image:        "postgis/postgis:14-3.3-alpine",
		Hostname:     config.HostName,
		ExposedPorts: []string{portName},
		Env: map[string]string{
			"POSTGRES_DB":       config.DatabaseName,
			"POSTGRES_USER":     username,
			"POSTGRES_PASSWORD": password,
			"PORT":              config.Port,
			"POSTGRES_SSL_MODE": "disable",
		},
		Cmd: []string{
			"postgres", "-c", "log_statement=all", "-c", "log_destination=stderr",
		},
		WaitingFor: wait.ForSQL(
			nat.Port(portName),
			"postgres",
			func(host string, port nat.Port) string {
				return fmt.Sprintf(templateURL, username, password, host, port.Port(), config.DatabaseName)
			},
		).WithStartupTimeout(time.Second * 20),
		Networks: []string{network.Name},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           log,
	})
	if err != nil {
		return nil, err
	}

	// Find port of container
	ports, err := container.Ports(ctx)
	if err != nil {
		return nil, err
	}

	portBinding := ports[nat.Port(portName)][0]
	driverURL := fmt.Sprintf(templateURL, username, password, portBinding.HostIP, portBinding.HostPort, config.DatabaseName)

	return &PostgresDBContainer{
		URI:       driverURL,
		Container: container,
		UserName:  username,
		Password:  password,
		Port:      portBinding.HostPort,
		Config:    config,
		network:   network,
		log:       log,
	}, nil
}
