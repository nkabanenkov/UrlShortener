package storage

import (
	"sync"
	"urlshortener/internal/urlshortener/encoder"
)

type inMemoryStorage struct {
	enc     encoder.Encoder
	counter uint64
	urls    map[uint64]string
	mutex   sync.Mutex
}

func NewInMemoryStorage(enc encoder.Encoder) *inMemoryStorage {
	return &inMemoryStorage{
		enc,
		0,
		make(map[uint64]string),
		sync.Mutex{},
	}
}

func (s *inMemoryStorage) Shorten(url string) (string, error) {
	s.mutex.Lock()

	s.urls[s.counter] = url
	oldCounter := s.counter
	s.counter++

	s.mutex.Unlock()

	return s.enc.Encode(oldCounter)
}

func (s *inMemoryStorage) Unshorten(url string) (string, error) {
	id, err := s.enc.Decode(url)
	if err != nil {
		return "", err
	}

	decodedUrl, exists := s.urls[id]
	if !exists {
		return "", UrlNotFoundError{"the given shortened URL doesn't correspond to any URL"}
	}

	return decodedUrl, nil
}

func (*inMemoryStorage) Close() {}
