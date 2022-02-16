package v2

import "strconv"

type uintValue struct {
	value uint64
	text  string
}

//Int creates UIntType value
func UInt(v uint64) Value {
	return &uintValue{
		value: v,
	}
}

func (v *uintValue) Value() interface{} {
	return v.value
}

func (v *uintValue) Type() ValueType {
	return IntType
}

func (v *uintValue) String() string {
	if v.text == "" {
		v.text = strconv.FormatUint(v.value, 10)
	}
	return v.text
}

func (v *uintValue) IsEmpty() bool {
	return v == nil || v.value == 0
}

func (p *byteParser) ParseUInt() (Value, error) {
	v, err := p.ParseNumber()
	if v == nil || v.Type() == UIntType {
		return v, err
	}
	if v.Type() == FloatType {
		f, ok := v.Value().(float64)
		if ok && f > 0 {
			return &uintValue{
				value: uint64(f),
			}, nil
		}
	}
	return nil, err
}
