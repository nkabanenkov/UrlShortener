package storage

import "context"

type Storage interface {
	Shorten(context.Context, string) (string, error)
	Unshorten(context.Context, string) (string, error)
}

type DatabaseError struct {
	message string
}

func (e DatabaseError) Error() string {
	return e.message
}

type UrlNotFoundError struct {
	message string
}

func (e UrlNotFoundError) Error() string {
	return e.message
}
