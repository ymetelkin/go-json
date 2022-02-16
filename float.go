package v2

import (
	"fmt"
)

type floatValue struct {
	value float64
	text  string
}

//Float creates FloatType value
func Float(v float64) Value {
	return &floatValue{
		value: v,
	}
}

func (v *floatValue) Value() interface{} {
	return v.value
}

func (v *floatValue) Type() ValueType {
	return FloatType
}

func (v *floatValue) String() string {
	if v.text == "" {
		v.text = fmt.Sprintf("%v", v.value)
	}
	return v.text
}

func (v *floatValue) Copy() Value {
	return Float(v.value)
}

func (v *floatValue) IsEmpty() bool {
	return v == nil || v.value == 0
}

func (p *byteParser) ParseFloat() (Value, error) {
	v, err := p.ParseNumber()
	if v == nil || v.Type() == FloatType {
		return v, err
	}
	if v.Type() == IntType {
		return &floatValue{
			value: float64(v.Value().(int)),
		}, err
	}
	return nil, err
}
