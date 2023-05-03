package encoder

import "testing"

func TestBadAlphabet(t *testing.T) {
	enc := NewBaseEncoder([]rune{}, 0)
	if enc != nil {
		t.Error("Empty alphabet should not be accepted")
	}
}

func TestNull(t *testing.T) {
	enc := NewBaseEncoder([]rune{'a'}, 0)
	if a, err := enc.Encode(0); a != "a" || err != nil {
		t.Error("Bad encoding of 0")
	}

	enc = NewBaseEncoder([]rune{'a'}, 5)
	if a, err := enc.Encode(0); a != "aaaaa" || err != nil {
		t.Error("Bad encoding of 0 with padding")
	}
}

func TestEncodingBinary(t *testing.T) {
	enc := NewBaseEncoder([]rune{'a', 'b'}, 0)
	if a, err := enc.Encode(3); a != "bb" || err != nil {
		t.Error("Bad encoding of 3")
	}
	if a, err := enc.Encode(4); a != "aab" || err != nil {
		t.Error("Bad encoding of 4")
	}
}

func TestEncodingOverflow(t *testing.T) {
	enc := NewBaseEncoder([]rune{'a', 'b'}, 3)
	if a, err := enc.Encode(7); a != "bbb" || err != nil {
		t.Error("Bad encoding of 7")
	}
	_, err := enc.Encode(8)
	if _, ok := err.(EncodingError); !ok {
		t.Error("Should've returned an EnvalidEncoding error")
	}
}

func TestDecodingBinary(t *testing.T) {
	enc := NewBaseEncoder([]rune{'a', 'b'}, 0)
	if n, err := enc.Decode("aaab"); n != 8 || err != nil {
		t.Error("Wrong decoding of 8")
	}
	if n, err := enc.Decode("baab"); n != 9 || err != nil {
		t.Error("Wrong decoding of 8")
	}
}

func TestEncodeDecode(t *testing.T) {
	enc := NewBaseEncoder([]rune{'a', 'b', 'c'}, 4)
	for i := 0; i < 3*3*3*3; i++ {
		a, err := enc.Encode(uint64(i))
		if err != nil {
			t.Errorf("Failed to encode %d", i)
		}

		n, err := enc.Decode(a)
		if err != nil {
			t.Errorf("Failed to decode %d", i)
		}

		if n != uint64(i) {
			t.Errorf("Failed to encode and then decode %d", i)
		}
	}
}

func TestDecodeInvalid(t *testing.T) {
	enc := NewBaseEncoder([]rune{'a', 'b', 'c'}, 4)
	_, err := enc.Decode("!!!")
	if _, ok := err.(DecodingError); !ok {
		t.Error("Decoding error expected")
	}

	_, err = enc.Decode("aa")
	if _, ok := err.(DecodingError); !ok {
		t.Error("Decoding error expected")
	}
}

func TestUnique(t *testing.T) {
	enc := NewBaseEncoder([]rune{'a', 'b', 'c'}, 4)
	m := make(map[string]struct{})

	for i := 0; i < 3*3*3*3; i++ {
		a, _ := enc.Encode(uint64(i))
		if _, found := m[a]; found {
			t.Errorf("Encoding of %d is not unique", i)
		}
		m[a] = struct{}{}
	}
}
