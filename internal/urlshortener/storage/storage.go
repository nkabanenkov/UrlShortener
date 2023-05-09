package storage

type Storage interface {
	Shorten(string) (string, error)
	Unshorten(string) (string, error)
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
