package urlshortener

import (
	"urlshortener/internal/urlshortener/storage"
	"urlshortener/internal/urlshortener/validator"
)

type UrlShortener struct {
	storage    storage.Storage
	validators []validator.Validator
}

func NewUrlShortener(stor storage.Storage) *UrlShortener {
	return &UrlShortener{stor, nil}
}

func (u *UrlShortener) AddValidator(validators ...validator.Validator) {
	u.validators = append(u.validators, validators...)
}

func (u *UrlShortener) Create(url string) (string, error) {
	for _, v := range u.validators {
		if err := v.Valid(url); err != nil {
			return "", err
		}
	}

	return u.storage.Create(url)
}

func (u *UrlShortener) Get(url string) (string, error) {
	return u.storage.Get(url)
}
