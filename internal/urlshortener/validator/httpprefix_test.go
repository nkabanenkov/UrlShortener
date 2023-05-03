package validator

import "testing"

func TestPrefix(t *testing.T) {
	val := NewHttpPrefixValidator()
	if _, ok := (val.Valid("abc://example.com")).(InvalidUrlError); !ok {
		t.Error("Expected an error")
	}
	if val.Valid("http://example.com") != nil || val.Valid("https://example.com") != nil {
		t.Error("Valid URLS are invalidated by the validator")
	}
}
