package testingcontainers

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/network"

	"github.com/testcontainers/testcontainers-go"
)

type PortMapping struct {
	ContainerPort int
	HostPort      int
	ContainerURL  string
	HostURL       string
}

type Network struct {
	network *testcontainers.DockerNetwork
	Name    string
}

func CreateNetwork(ctx context.Context) (Network, error) {
	dockerNetwork, err := network.New(ctx, network.WithCheckDuplicate())
	if err != nil {
		return Network{}, err
	}

	return Network{network: dockerNetwork, Name: dockerNetwork.Name}, err
}

func (n Network) Cleanup(ctx context.Context) error {
	return n.network.Remove(ctx)
}

func getRandomOpenPort(base int, maxOffset int64) (int, error) {
	var err error
	var port int
	for {
		port, err = getRandomPort(base, maxOffset)
		if err != nil {
			break
		}

		portOpen, err := isPortOpen(port)
		if err != nil {
			break
		}

		if portOpen {
			break
		}
	}

	return port, err
}

func getRandomPort(base int, maxOffset int64) (int, error) {
	maxPortOffset := big.NewInt(maxOffset)
	offset, err := rand.Int(rand.Reader, maxPortOffset)
	if err != nil {
		return -1, err
	}

	return base + int(offset.Int64()), nil
}

func isPortOpen(port int) (bool, error) {
	timeout := time.Second
	host := net.JoinHostPort("localhost", strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return false, err
	}

	err = conn.Close()
	if err != nil {
		return false, err
	}
	return true, nil
}

func resolvePortMapping(binding nat.PortBinding, containerHost string, containerPort int) (PortMapping, error) {
	hostPort, err := strconv.Atoi(binding.HostPort)
	if err != nil {
		return PortMapping{}, err
	}
	return PortMapping{
		ContainerPort: containerPort,
		HostPort:      hostPort,
		ContainerURL:  fmt.Sprintf("http://%s:%d", containerHost, containerPort),
		HostURL:       fmt.Sprintf("http://%s:%d", binding.HostIP, hostPort),
	}, nil
}

func resolvePortMappings(ctx context.Context, container testcontainers.Container, containerHost string) (map[int]PortMapping, error) {
	portMap, err := container.Ports(ctx)
	if err != nil {
		return nil, err
	}

	portMappings := make(map[int]PortMapping)
	for k, v := range portMap {
		if len(v) != 0 {
			portMapping, err := resolvePortMapping(v[0], containerHost, k.Int())
			if err != nil {
				return nil, err
			}
			portMappings[k.Int()] = portMapping
		}
	}

	return portMappings, nil
}
