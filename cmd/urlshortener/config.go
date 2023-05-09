package main

import (
	"log"
	"os"
	"strconv"
	"urlshortener/pkg/urlshortener/config"
)

func readConfig() config.Config {
	c := config.Config{}

	env, found := os.LookupEnv("URLSHORTENER_ALPHABET")
	if !found {
		log.Fatalln("URLSHORTENER_ALPHABET is not set")
	}
	c.Alphabet = []rune(env)

	env, found = os.LookupEnv("URLSHORTENER_WIDTH")
	if !found {
		log.Fatalln("URLSHORTENER_WIDTH is not set")
	}
	width, err := strconv.Atoi(env)
	if err != nil || width < 0 {
		log.Fatalln("Invalid URLSHORTENER_WIDTH")
	}
	c.Width = uint(width)

	env, found = os.LookupEnv("URLSHORTENER_INMEMORY")
	if !found {
		c.InMemory = false
	} else {
		c.InMemory, err = strconv.ParseBool(env)
		if err != nil {
			log.Fatalln("Invalid URLSHORTENER_INMEMORY")
		}
	}

	if !c.InMemory {
		env, found := os.LookupEnv("DB_HOST")
		if !found {
			log.Fatalln("DB_HOST must be set")
		}
		c.DbHost = env

		env, found = os.LookupEnv("DB_PORT")
		if !found {
			log.Fatalln("DB_PORT is not set")
		}
		port, err := strconv.Atoi(env)
		if err != nil || port < 0 {
			log.Fatalln("Invalid POSTGRES_PORT")
		}
		c.DbPort = uint(port)

		env, found = os.LookupEnv("POSTGRES_USER")
		if !found {
			log.Fatalln("POSTGRES_USER must be set")
		}
		c.DbUser = env

		env, found = os.LookupEnv("POSTGRES_PASSWORD")
		if !found {
			log.Fatalln("POSTGRES_PASSWORD must be set")
		}
		c.DbPassword = env

		env, found = os.LookupEnv("POSTGRES_DB")
		if !found {
			log.Fatalln("POSTGRES_DB must be set")
		}
		c.DbName = env
	}

	return c
}
