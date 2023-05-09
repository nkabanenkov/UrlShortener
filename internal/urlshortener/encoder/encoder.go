package encoder

type Encoder interface {
	Encode(uint64) (string, error)
	Decode(string) (uint64, error)
}

type EncodingOverflowError struct {
	message string
}

func (e EncodingOverflowError) Error() string {
	return e.message
}

type DecodingError struct {
	message string
}

func (e DecodingError) Error() string {
	return e.message
}
