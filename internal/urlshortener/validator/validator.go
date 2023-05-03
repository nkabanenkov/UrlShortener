package validator

type Validator interface {
	Valid(string) error
}

type InvalidUrlError struct {
	message string
}

func (e InvalidUrlError) Error() string {
	return e.message
}
