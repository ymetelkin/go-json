package v2

import "fmt"

type nullValue struct {
}

//Null creates NullType value
func Null() Value {
	return new(nullValue)
}

func (v *nullValue) Value() interface{} {
	return nil
}

func (v *nullValue) Type() ValueType {
	return NullType
}

func (v *nullValue) String() string {
	return "null"
}

func (v *nullValue) IsEmpty() bool {
	return true
}

//used when a first byte is 'n'
func (p *byteParser) ParseNull() (Value, error) {
	if p.Index+4 > p.Size {
		return nil, fmt.Errorf("parsing [ null ] at %d: EOF", p.Index)
	}

	var (
		expected = []byte{'u', 'l', 'l'}
		idx      = p.Index + 1
	)

	for i := 0; i < 3; i++ {
		if p.Bytes[idx+i] != expected[i] {
			return nil, fmt.Errorf("parsing [ null ] at %d: expecting [ %s ], found [ %s ]", idx+i, string(expected[i]), string(p.Bytes[idx+i]))
		}
	}

	p.Index += 3
	return &nullValue{}, p.SkipWS()
}
