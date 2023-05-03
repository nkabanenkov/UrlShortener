package storage

type Storage interface {
	Create(string) (string, error)
	Get(string) (string, error)
	Close()
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
