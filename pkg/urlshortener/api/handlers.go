package api

import (
	"log"
	"net/http"
	"urlshortener/internal/urlshortener/encoder"
	"urlshortener/internal/urlshortener/storage"
	"urlshortener/internal/urlshortener/validator"
	"urlshortener/pkg/urlshortener"

	"github.com/gin-gonic/gin"
)

func MakeGetHandler(app *urlshortener.UrlShortener, paramName string) func(*gin.Context) {
	return func(c *gin.Context) {
		url := c.Param(paramName)
		if url == "" {
			c.JSON(http.StatusBadRequest, UnshortenAnswer{"empty URL"})
			return
		}

		decodedUrl, err := app.Unshorten(url)
		if err != nil {
			if _, ok := err.(encoder.DecodingError); ok {
				c.JSON(http.StatusBadRequest, UnshortenAnswer{err.Error()})
			} else if _, ok := err.(storage.DatabaseError); ok {
				log.Println("Database error occured: " + err.Error())
				c.JSON(http.StatusInternalServerError, UnshortenAnswer{err.Error()})
				return
			} else if _, ok := err.(storage.UrlNotFoundError); ok {
				c.JSON(http.StatusNotFound, UnshortenAnswer{err.Error()})
				return
			} else {
				panic(err.Error())
			}
		}

		c.JSON(http.StatusOK, UnshortenAnswer{decodedUrl})
	}
}

func MakePostHandler(app *urlshortener.UrlShortener) func(*gin.Context) {
	return func(c *gin.Context) {
		var req ShortenRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ShortenAnswer{"invalid body request format"})
			return
		}

		encodedUrl, err := app.Shorten(req.Url)
		if err != nil {
			if _, ok := err.(validator.InvalidUrlError); ok {
				c.JSON(http.StatusBadRequest, ShortenAnswer{err.Error()})
				return
			} else if _, ok := err.(storage.DatabaseError); ok {
				log.Println("Database error occured: " + err.Error())
				c.JSON(http.StatusInternalServerError, ShortenAnswer{err.Error()})
				return
			} else if _, ok := err.(encoder.EncodingOverflowError); ok {
				log.Fatalln("Encoding overflow")
			} else {
				panic(err.Error())
			}
		}

		c.JSON(http.StatusOK, ShortenAnswer{encodedUrl})
	}
}
