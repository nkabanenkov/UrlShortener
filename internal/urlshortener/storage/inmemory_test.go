package storage_test

import (
	"testing"
	"urlshortener/internal/urlshortener/storage"
)

func powInt64(b uint64, e uint) uint64 {
	if e == 0 {
		return 1
	}
	var p uint64 = 1
	for i := uint(0); i < e; i++ {
		p *= b
	}
	return b
}

func TestInMemNotFound(t *testing.T) {
	stor := storage.NewInMemoryStorage(testEnc)
	testNotFound(stor, t)
}

func TestInMemBadEncoding(t *testing.T) {
	stor := storage.NewInMemoryStorage(testEnc)
	testBadEncoding(stor, t)
}

func TestInMemCreateGet(t *testing.T) {
	stor := storage.NewInMemoryStorage(testEnc)
	testCreateGet(stor, t)
}
