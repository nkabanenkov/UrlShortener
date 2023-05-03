package storage_test

import (
	"context"
	"testing"
	"urlshortener/internal/urlshortener/storage"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreateContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15.2-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "urlshortener",
		},
		HostConfigModifier: func(hc *container.HostConfig) {
			hc.NetworkMode = "host"
		},
		WaitingFor: wait.ForLog("PostgreSQL init process complete; ready for start up"),
		SkipReaper: true,
	}
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func Connect(ctx context.Context, cont testcontainers.Container) (storage.Storage, error) {
	ip, _ := cont.Host(ctx)
	port, _ := cont.MappedPort(ctx, "5432")
	return storage.NewPgStorage(testEnc, storage.DbConfig{
		ip,
		uint(port.Int()),
		"postgres",
		"postgres",
		"urlshortener",
	})
}

func TestPgNotFound(t *testing.T) {
	ctx := context.Background()
	cont, err := CreateContainer(ctx)
	if err != nil {
		panic(err)
	}
	defer cont.Terminate(ctx)

	stor, err := Connect(ctx, cont)
	if err != nil {
		panic(err)
	}
	defer stor.Close()

	testNotFound(stor, t)
}

func TestPgBadEncoding(t *testing.T) {
	ctx := context.Background()
	cont, err := CreateContainer(ctx)
	if err != nil {
		panic(err)
	}
	defer cont.Terminate(ctx)

	stor, err := Connect(ctx, cont)
	if err != nil {
		panic(err)
	}
	defer stor.Close()

	testBadEncoding(stor, t)
}

func TestPgCreateGet(t *testing.T) {
	ctx := context.Background()
	cont, err := CreateContainer(ctx)
	if err != nil {
		panic(err)
	}
	defer cont.Terminate(ctx)

	stor, err := Connect(ctx, cont)
	if err != nil {
		panic(err)
	}
	defer stor.Close()

	testCreateGet(stor, t)
}
