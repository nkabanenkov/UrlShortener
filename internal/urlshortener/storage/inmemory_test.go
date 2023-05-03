package storage

import (
	"fmt"
	"testing"
	"urlshortener/internal/urlshortener/encoder"
)

var alphabet = []rune{'a', 'b', 'c', 'd'}
var padding = 10
var enc = encoder.NewBaseEncoder(alphabet, uint(padding))

func TestNotFound(t *testing.T) {
	b := NewInMemoryStorage(enc)
	_, err := b.Get("aaabcdabcd")
	if _, ok := err.(UrlNotFoundError); !ok {
		t.Error("Expected to recieve an error")
	}
}

func TestCreateGet(t *testing.T) {
	b := NewInMemoryStorage(enc)
	uniq := make(map[string]string)
	pattern := "https://example.com/?id=%d"

	for i := 0; i < len(alphabet); i++ {
		encoded, err := b.Create(fmt.Sprintf(pattern, i))
		if err != nil {
			t.Errorf("Failed to encode %d-th string", i)
		}
		if _, found := uniq[encoded]; found {
			t.Errorf("Encoded %d-th string is not unique", i)
		}

		decoded, err := b.Get(encoded)
		if err != nil {
			t.Errorf("Failed to decode %d-th string", i)
		}
		if decoded != fmt.Sprintf(pattern, i) {
			t.Errorf("Create-Get failed for %d-th string", i)
		}
	}
}
