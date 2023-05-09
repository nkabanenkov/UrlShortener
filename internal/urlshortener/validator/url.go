package validator

type UrlValidator struct {
	validCharacters map[rune]struct{}
}

func MakeUrlValidator() UrlValidator {
	val := make(map[rune]struct{})
	for i := 'a'; i <= 'z'; i++ {
		val[rune(i)] = struct{}{}
	}
	for i := 'A'; i <= 'Z'; i++ {
		val[rune(i)] = struct{}{}
	}
	for i := '0'; i <= '9'; i++ {
		val[rune(i)] = struct{}{}
	}
	for _, ch := range "-._~:/?#[]@!$&%'()*+,;%=" {
		val[ch] = struct{}{}
	}

	return UrlValidator{val}
}

// Returns `nil` if all the characters in `url` are valid in RFC 3986 or `InvalidUrlError` otherwise.
// https://www.ietf.org/rfc/rfc3986.txt
func (uv UrlValidator) Valid(url string) error {
	for _, ch := range url {
		if _, found := uv.validCharacters[ch]; !found {
			return InvalidUrlError{"invalid character in URL"}
		}
	}
	return nil
}
