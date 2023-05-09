package storage_test

import (
	"context"
	"os"
	"testing"
	"urlshortener/internal/urlshortener/storage"
	"urlshortener/pkg/urlshortener/config"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func CreateContainer(ctx context.Context) (testcontainers.Container, error) {
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	return postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.2-alpine"),
		postgres.WithDatabase("urlshortener"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2)),
		testcontainers.WithHostConfigModifier(func(hc *container.HostConfig) {
			hc.NetworkMode = "host"
		}))
}

func Connect(ctx context.Context, cont testcontainers.Container) (storage.Storage, error) {
	ip, _ := cont.Host(ctx)
	port, _ := cont.MappedPort(ctx, "5432")
	return storage.NewPgStorage(testEnc, config.Config{
		DbHost:     ip,
		DbPort:     uint(port.Int()),
		DbUser:     "postgres",
		DbPassword: "postgres",
		DbName:     "urlshortener",
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

	testCreateGet(stor, t)
}
