package main

import (
	"os"
	"strconv"
	"urlshortener/internal/urlshortener/storage"
)

type config struct {
	alphabet []rune
	width    uint
	inMemory bool
}

func readConfig() config {
	c := config{}

	env, found := os.LookupEnv("URLSHORTENER_ALPHABET")
	if !found {
		panic("URLSHORTENER_ALPHABET is not set")
	}
	c.alphabet = []rune(env)

	env, found = os.LookupEnv("URLSHORTENER_WIDTH")
	if !found {
		panic("URLSHORTENER_WIDTH is not set")
	}
	width, err := strconv.Atoi(env)
	if err != nil || width < 0 {
		panic("Invalid URLSHORTENER_WIDTH")
	}
	c.width = uint(width)

	env, found = os.LookupEnv("URLSHORTENER_INMEMORY")
	if !found {
		c.inMemory = false
	} else {
		inMem, err := strconv.ParseBool(env)
		if err != nil {
			panic("Invalid URLSHORTENER_INMEMORY")
		}
		c.inMemory = inMem
	}

	return c
}

func readDbConfig() storage.DbConfig {
	c := storage.DbConfig{}

	env, found := os.LookupEnv("DB_HOST")
	if !found {
		panic("DB_HOST must be set")
	}
	c.Hostname = env

	env, found = os.LookupEnv("DB_PORT")
	if !found {
		panic("DB_PORT is not set")
	}
	port, err := strconv.Atoi(env)
	if err != nil || port < 0 {
		panic("Invalid POSTGRES_PORT")
	}
	c.Port = uint(port)

	env, found = os.LookupEnv("POSTGRES_USER")
	if !found {
		panic("POSTGRES_USER must be set")
	}
	c.User = env

	env, found = os.LookupEnv("POSTGRES_PASSWORD")
	if !found {
		panic("POSTGRES_PASSWORD must be set")
	}
	c.Password = env

	env, found = os.LookupEnv("POSTGRES_DB")
	if !found {
		panic("POSTGRES_DB must be set")
	}
	c.DbName = env

	return c
}
