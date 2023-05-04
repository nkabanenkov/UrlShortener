package validator_test

import (
	"testing"
	"urlshortener/internal/urlshortener/validator"
)

func TestBadUrls(t *testing.T) {
	val := validator.MakeUrlValidator()
	if _, ok := (val.Valid("ftp://space in url.com")).(validator.InvalidUrlError); !ok {
		t.Error("Expected an error")
	}
	if _, ok := (val.Valid("ftp://\"example.com|")).(validator.InvalidUrlError); !ok {
		t.Error("Expected an error")
	}
	if _, ok := (val.Valid("ftp://example.com/1<>2")).(validator.InvalidUrlError); !ok {
		t.Error("Expected an error")
	}
}

func TestOkUrls(t *testing.T) {
	val := validator.MakeUrlValidator()
	if val.Valid("http://example09.com:8080/index.php?id=1%202&date_test=~01.01.2023;10.10.2023$") != nil {
		t.Error("Valid URLS are invalidated by the validator")
	}
}
