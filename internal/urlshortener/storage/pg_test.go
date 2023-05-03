package storage_test

import (
	"context"
	"fmt"
	"testing"
	"urlshortener/internal/urlshortener"
	"urlshortener/internal/urlshortener/storage"
	"urlshortener/internal/urlshortener/validator"

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
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
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
	return storage.NewPgStorage(enc, storage.DbConfig{
		ip,
		uint(port.Int()),
		"postgres",
		"postgres",
		"urlshortener",
	})
}

func TestPg(t *testing.T) {
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

	app := urlshortener.NewUrlShortener(stor)
	app.AddValidator(validator.NewHttpPrefixValidator())

	shortnedUrls := make([]string, 0)
	pattern := "https://example.com/i=%d"
	for i := range alphabet {
		url := fmt.Sprintf(pattern, i)
		shortenedUrl, err := app.Create(url)
		if err != nil {
			t.Errorf("Failed to shorten %d-th url", i)
		}
		shortnedUrls = append(shortnedUrls, shortenedUrl)
	}

	for i := range shortnedUrls {
		url, err := app.Get(shortnedUrls[i])
		if err != nil {
			panic(err)
		}
		if fmt.Sprintf(pattern, i) != url {
			t.Error("Encoded and decoded urls are not equal")
		}
	}
}
