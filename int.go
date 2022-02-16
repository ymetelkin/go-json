package v2

import "strconv"

type intValue struct {
	value int
	text  string
}

//Int creates IntType value
func Int(v int) Value {
	return &intValue{
		value: v,
	}
}

func (v *intValue) Value() interface{} {
	return v.value
}

func (v *intValue) Type() ValueType {
	return IntType
}

func (v *intValue) String() string {
	if v.text == "" {
		v.text = strconv.Itoa(v.value)
	}
	return v.text
}

func (v *intValue) IsEmpty() bool {
	return v == nil || v.value == 0
}

func (p *byteParser) ParseInt() (Value, error) {
	v, err := p.ParseNumber()
	if v == nil || v.Type() == IntType {
		return v, err
	}
	if v.Type() == FloatType {
		f, ok := v.Value().(float64)
		if ok && f > 0 {
			return &intValue{
				value: int(f),
			}, nil
		}
	}
	return nil, err
}
