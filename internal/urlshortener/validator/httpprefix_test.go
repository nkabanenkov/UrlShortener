package validator_test

import (
	"testing"
	"urlshortener/internal/urlshortener/validator"
)

func TestPrefix(t *testing.T) {
	val := validator.MakeHttpPrefixValidator()
	if _, ok := (val.Valid("abc://example.com")).(validator.InvalidUrlError); !ok {
		t.Error("Expected an error")
	}
	if val.Valid("http://example.com") != nil || val.Valid("https://example.com") != nil {
		t.Error("Valid URLS are invalidated by the validator")
	}
}
