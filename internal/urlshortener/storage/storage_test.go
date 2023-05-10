package storage_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"urlshortener/internal/urlshortener/encoder"
	"urlshortener/internal/urlshortener/storage"
)

var testAlphabet = []rune{'a', 'b', 'c', 'd'}
var testWidth = uint(5)
var testPower = powInt64(uint64(len(testAlphabet)), testWidth)
var testEnc = encoder.NewBaseEncoder(testAlphabet, testWidth)

func testNotFound(stor storage.Storage, t *testing.T) {
	builder := strings.Builder{}
	for i := 0; i < int(testWidth); i++ {
		builder.WriteRune(testAlphabet[0])
	}
	_, err := stor.Unshorten(context.Background(), builder.String())
	if _, ok := err.(storage.UrlNotFoundError); !ok {
		t.Error("Expected to recieve an error")
	}
}

func testBadEncoding(stor storage.Storage, t *testing.T) {
	_, err := stor.Unshorten(context.Background(), "!!!")
	if _, ok := err.(encoder.DecodingError); !ok {
		t.Error("Expected to recieve an error")
	}
}

func testCreateGet(stor storage.Storage, t *testing.T) {
	uniq := make(map[string]string)
	pattern := "https://example.com/?id=%d"

	for i := 0; i < int(testPower); i++ {
		encoded, err := stor.Shorten(context.Background(), fmt.Sprintf(pattern, i))
		if err != nil {
			t.Errorf("Failed to encode %d-th string", i)
		}
		if _, found := uniq[encoded]; found {
			t.Errorf("Encoded %d-th string is not unique", i)
		}

		decoded, err := stor.Unshorten(context.Background(), encoded)
		if err != nil {
			t.Errorf("Failed to decode %d-th string", i)
		}
		if decoded != fmt.Sprintf(pattern, i) {
			t.Errorf("Create-Get failed for %d-th string", i)
		}
	}
}
