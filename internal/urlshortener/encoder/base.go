package encoder

import (
	"strings"
)

type baseEncoder struct {
	encAlphabet []rune
	decAlphabet map[rune]int
	arity       int
	padding     uint
}

// Takes `encAlphabet` -- an encoding table that maps 0, ..., n-1 to a_0, .., a_{n-1}.
// `padding` is the length of encoded string. It can be 0, which means there encoded string aren't padded with a_0.
func NewBaseEncoder(encAlphabet []rune, padding uint) *baseEncoder {
	arity := len(encAlphabet)
	if arity == 0 {
		return nil
	}

	enc := baseEncoder{
		make([]rune, arity),
		make(map[rune]int),
		arity,
		padding,
	}

	for i := 0; i < enc.arity; i++ {
		ch := encAlphabet[i]
		enc.encAlphabet[i] = ch
		enc.decAlphabet[ch] = i
	}

	return &enc
}

// Encodes `n` with `encAlphabet`. If `padding` is not 0, makes sure the encoded string fits.
// Returns `EncodingError` if failed to fit.
func (e baseEncoder) Encode(n uint64) (string, error) {
	builder := strings.Builder{}
	if e.padding > 0 {
		builder.Grow(int(e.padding))
	} else if n == 0 {
		return string(e.encAlphabet[0]), nil
	}

	for n != 0 {
		builder.WriteRune(e.encAlphabet[n%uint64(e.arity)])
		n /= uint64(e.arity)
	}

	if e.padding > 0 {
		diff := int(e.padding) - builder.Len()
		if diff < 0 {
			return "", EncodingError{"encoding overflow"}
		}
		for i := 0; i < diff; i++ {
			builder.WriteRune(e.encAlphabet[0])
		}
	}

	return builder.String(), nil
}

// Decodes `encoded` with `decAlphabet`.
// Returns `DecodingError` if an unknown character is met or the shortened URL is invalid (when `padding` is not 0).
func (e baseEncoder) Decode(encoded string) (uint64, error) {
	if e.padding != 0 && int(e.padding) != len(encoded) {
		return 0, DecodingError{"malformed shortened url"}
	}

	var decoded uint64
	pow := uint64(1)
	for _, ch := range encoded {
		digit, ok := e.decAlphabet[ch]
		if !ok {
			return 0, DecodingError{"a symbol not from the alphabet occured"}
		}

		decoded += uint64(digit) * pow
		pow *= uint64(e.arity)
	}

	return decoded, nil
}
