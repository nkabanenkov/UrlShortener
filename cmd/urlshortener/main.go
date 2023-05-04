package main

import (
	"urlshortener/internal/urlshortener"
	"urlshortener/internal/urlshortener/api"
	"urlshortener/internal/urlshortener/encoder"
	"urlshortener/internal/urlshortener/storage"
	"urlshortener/internal/urlshortener/validator"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := readConfig()
	var stor storage.Storage

	enc := encoder.NewBaseEncoder(cfg.alphabet, cfg.width)
	if cfg.inMemory {
		stor = storage.NewInMemoryStorage(enc)
	} else {
		var err error
		stor, err = storage.NewPgStorage(enc, readDbConfig())
		if err != nil {
			panic(err)
		}
	}
	defer stor.Close()

	app := urlshortener.NewUrlShortener(stor)
	app.AddValidator(validator.MakeUrlValidator())

	server := gin.Default()
	server.GET("/urls/:url", api.MakeUrlsGetHandler(app, "url"))
	server.POST("/urls", api.MakeUrlsPostHandler(app))
	server.Run(":8000")
}
