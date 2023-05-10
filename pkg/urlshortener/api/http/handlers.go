package http

import (
	"log"
	"net/http"
	"urlshortener/internal/urlshortener/encoder"
	"urlshortener/internal/urlshortener/storage"
	"urlshortener/internal/urlshortener/validator"
	"urlshortener/pkg/urlshortener"

	"github.com/gin-gonic/gin"
)

type urlShortenerHttpHandlers struct {
	app *urlshortener.UrlShortener
}

func MakeUrlShortenerHttpHandlers(app *urlshortener.UrlShortener) urlShortenerHttpHandlers {
	return urlShortenerHttpHandlers{app}
}

func (h *urlShortenerHttpHandlers) GetHandler(ctx *gin.Context) {
	url := ctx.Param("url")
	if url == "" {
		ctx.JSON(http.StatusBadRequest, UnshortenResponse{"empty URL"})
		return
	}

	decodedUrl, err := h.app.Unshorten(ctx, url)
	if err != nil {
		switch err.(type) {
		case encoder.DecodingError:
			ctx.JSON(http.StatusBadRequest, UnshortenResponse{err.Error()})
			return
		case storage.DatabaseError:
			log.Println("Database error occured: " + err.Error())
			ctx.JSON(http.StatusInternalServerError, UnshortenResponse{err.Error()})
			return
		case storage.UrlNotFoundError:
			ctx.JSON(http.StatusNotFound, UnshortenResponse{err.Error()})
			return
		default:
			panic(err.Error())
		}
	}

	ctx.JSON(http.StatusOK, UnshortenResponse{decodedUrl})
}

func (h *urlShortenerHttpHandlers) PostHandler(ctx *gin.Context) {
	var req ShortenRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ShortenResponse{"invalid body request format"})
		return
	}

	encodedUrl, err := h.app.Shorten(ctx, req.Url)
	if err != nil {
		switch err.(type) {
		case validator.InvalidUrlError:
			ctx.JSON(http.StatusBadRequest, ShortenResponse{err.Error()})
			return
		case storage.DatabaseError:
			log.Println("Database error occured: " + err.Error())
			ctx.JSON(http.StatusInternalServerError, ShortenResponse{err.Error()})
			return
		case encoder.EncodingOverflowError:
			log.Fatalln("Encoding overflow")
		default:
			panic(err.Error())
		}
	}

	ctx.JSON(http.StatusOK, ShortenResponse{encodedUrl})
}
