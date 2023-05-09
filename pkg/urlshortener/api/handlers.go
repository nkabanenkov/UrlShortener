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
			switch err.(type) {
			case encoder.DecodingError:
				c.JSON(http.StatusBadRequest, UnshortenAnswer{err.Error()})
				return
			case storage.DatabaseError:
				log.Println("Database error occured: " + err.Error())
				c.JSON(http.StatusInternalServerError, UnshortenAnswer{err.Error()})
				return
			case storage.UrlNotFoundError:
				c.JSON(http.StatusNotFound, UnshortenAnswer{err.Error()})
				return
			default:
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
			switch err.(type) {
			case validator.InvalidUrlError:
				c.JSON(http.StatusBadRequest, ShortenAnswer{err.Error()})
				return
			case storage.DatabaseError:
				log.Println("Database error occured: " + err.Error())
				c.JSON(http.StatusInternalServerError, ShortenAnswer{err.Error()})
				return
			case encoder.EncodingOverflowError:
				log.Fatalln("Encoding overflow")
			default:
				panic(err.Error())
			}
		}

		c.JSON(http.StatusOK, ShortenAnswer{encodedUrl})
	}
}
