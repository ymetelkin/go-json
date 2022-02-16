package v2

import "fmt"

type boolValue struct {
	value bool
}

//Bool creates BoolType value
func Bool(v bool) Value {
	return &boolValue{
		value: v,
	}
}

func (v *boolValue) Value() interface{} {
	return v.value
}

func (v *boolValue) Type() ValueType {
	return BoolType
}

func (v *boolValue) String() string {
	if v.value {
		return "true"
	}
	return "false"
}

func (v *boolValue) Copy() Value {
	return Bool(v.value)
}

func (v *boolValue) IsEmpty() bool {
	return v == nil || !v.value
}

//used when a first byte is 't'
func (p *byteParser) ParseTrue() (Value, error) {
	if p.Index+4 > p.Size {
		return nil, fmt.Errorf("parsing [ true ] at %d: EOF", p.Index)
	}

	var (
		expected = []byte{'r', 'u', 'e'}
		idx      = p.Index + 1
	)

	for i := 0; i < 3; i++ {
		if p.Bytes[idx+i] != expected[i] {
			return nil, fmt.Errorf("parsing [ true ] at %d: expecting [ %s ], found [ %s ]", idx+i, string(expected[i]), string(p.Bytes[idx+i]))
		}
	}

	p.Index += 3
	return &boolValue{value: true}, p.SkipWS()
}

//used when a first byte is 'f'
func (p *byteParser) ParseFalse() (Value, error) {
	if p.Index+5 > p.Size {
		return nil, fmt.Errorf("parsing [ false ] at %d: EOF", p.Index)
	}

	var (
		expected = []byte{'a', 'l', 's', 'e'}
		idx      = p.Index + 1
	)

	for i := 0; i < 4; i++ {
		if p.Bytes[idx+i] != expected[i] {
			return nil, fmt.Errorf("parsing [ false ] at %d: expecting [ %s ], found [ %s ]", idx+i, string(expected[i]), string(p.Bytes[idx+i]))
		}
	}

	p.Index += 4
	return &boolValue{}, p.SkipWS()
}
