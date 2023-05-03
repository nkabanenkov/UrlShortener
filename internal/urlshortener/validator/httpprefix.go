package validator

import "strings"

type httpPrefixValidator struct{}

func NewHttpPrefixValidator() *httpPrefixValidator {
	return &httpPrefixValidator{}
}

// Returns `nil`, if `url` begins with `http://` or `https://` or `InvalidUrlError` otherwise.
func (httpPrefixValidator) Valid(url string) error {
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		return InvalidUrlError{"an URL must begin with either http:// or https://"}
	}
	return nil
}
