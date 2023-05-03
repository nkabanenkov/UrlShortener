package encoder

type Encoder interface {
	Encode(uint64) (string, error)
	Decode(string) (uint64, error)
}

type EncodingError struct {
	message string
}

func (e EncodingError) Error() string {
	return e.message
}

type DecodingError struct {
	message string
}

func (e DecodingError) Error() string {
	return e.message
}
