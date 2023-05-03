package api

import (
	"net/http"
	"urlshortener/internal/urlshortener"
	"urlshortener/internal/urlshortener/encoder"
	"urlshortener/internal/urlshortener/storage"
	"urlshortener/internal/urlshortener/validator"

	"github.com/gin-gonic/gin"
)

func MakeUrlsGetHandler(app *urlshortener.UrlShortener, paramName string) func(*gin.Context) {
	return func(c *gin.Context) {
		url := c.Param(paramName)
		if url == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "empty URL"})
			return
		}

		decodedUrl, err := app.Get(url)
		if err != nil {
			if _, ok := err.(encoder.DecodingError); ok {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			} else if _, ok := err.(storage.DatabaseError); ok {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			} else if _, ok := err.(storage.UrlNotFoundError); ok {
				c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
				return
			} else {
				panic(err.Error())
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": decodedUrl})
	}
}

func MakeUrlsPostHandler(app *urlshortener.UrlShortener) func(*gin.Context) {
	return func(c *gin.Context) {
		var body EncodingRequestBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "invalid body request format"})
			return
		}

		encodedUrl, err := app.Create(body.Url)
		if err != nil {
			if _, ok := err.(validator.InvalidUrlError); ok {
				c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			} else if _, ok := err.(storage.DatabaseError); ok {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			} else {
				panic(err.Error())
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": encodedUrl})
	}
}
