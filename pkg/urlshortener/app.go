package urlshortener

import (
	"context"
	"urlshortener/internal/urlshortener/encoder"
	"urlshortener/internal/urlshortener/storage"
	"urlshortener/internal/urlshortener/validator"
	"urlshortener/pkg/urlshortener/config"

	"github.com/pkg/errors"
)

type UrlShortener struct {
	storage    storage.Storage
	validators []validator.Validator
}

func MakeUrlShortener(cfg config.Config) (UrlShortener, error) {
	enc := encoder.NewBaseEncoder(cfg.Alphabet, cfg.Width)
	if enc == nil {
		return UrlShortener{}, errors.New("alphabet is empty")
	}

	var stor storage.Storage
	var err error
	if cfg.InMemory {
		stor = storage.NewInMemoryStorage(enc)
	} else {
		stor, err = storage.NewPgStorage(enc, cfg)
		if err != nil {
			return UrlShortener{}, err
		}
	}

	app := UrlShortener{stor, nil}

	app.addValidator(validator.MakeHttpPrefixValidator())
	app.addValidator(validator.MakeUrlValidator())

	return app, nil
}

func (u *UrlShortener) Shorten(ctx context.Context, url string) (string, error) {
	for _, v := range u.validators {
		if err := v.Valid(url); err != nil {
			return "", err
		}
	}

	return u.storage.Shorten(ctx, url)
}

func (u *UrlShortener) Unshorten(ctx context.Context, url string) (string, error) {
	return u.storage.Unshorten(ctx, url)
}

func (u *UrlShortener) addValidator(validators ...validator.Validator) {
	u.validators = append(u.validators, validators...)
}
