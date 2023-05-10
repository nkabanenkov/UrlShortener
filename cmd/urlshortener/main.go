package main

import (
	"log"
	"urlshortener/pkg/urlshortener"
	"urlshortener/pkg/urlshortener/api/http"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := readConfig()
	app, err := urlshortener.MakeUrlShortener(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}

	server := gin.Default()
	server.GET("/urls/:url", http.MakeGetHandler(&app, "url"))
	server.POST("/urls", http.MakePostHandler(&app))
	if err = server.Run(":8000"); err != nil {
		log.Fatalln(err.Error())
	}
}
