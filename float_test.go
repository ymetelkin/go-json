package v2

import (
	"fmt"
	"testing"
)

func TestFloatParse(t *testing.T) {
	s := "3.14"
	p := newParser([]byte(s))
	err := p.SkipWS()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '3' {
		t.Error("Failed to parse [ 3 ]")
	}
	v, err := p.ParseFloat()
	if v.Value() != 3.14 || err != errEOF {
		t.Error("Failed to parse [ 3.14 ]")
	}
	fmt.Println(v)

	s = "0.001"
	p = newParser([]byte(s))
	err = p.SkipWS()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '0' {
		t.Error("Failed to parse [ 0 ]")
	}
	v, err = p.ParseFloat()
	if v.Value() != 0.001 || err != errEOF {
		t.Error("Failed to parse [ 0.001 ]")
	}
	fmt.Println(v)

	s = "-123"
	p = newParser([]byte(s))
	err = p.SkipWS()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '-' {
		t.Error("Failed to parse [ - ]")
	}
	v, err = p.ParseFloat()
	if v.Value() != float64(-123) || err != errEOF {
		t.Error("Failed to parse [ -123 ]")
	}
	fmt.Println(v)

	s = "3.14E-2"
	p = newParser([]byte(s))
	err = p.SkipWS()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '3' {
		t.Error("Failed to parse [ 3 ]")
	}
	v, err = p.ParseFloat()
	if v.Value() != 3.14e-2 || err != errEOF {
		t.Error("Failed to parse [ 3.14e-2 ]")
	}
	fmt.Println(v)
	fmt.Println(v.Value())

	s = "{ a: -0.123, b: 2.3E2 }"
	p = newParser([]byte(s))
	p.SkipWS() //{
	fmt.Println(string(p.Byte))
	p.SkipWS() //a
	fmt.Println(string(p.Byte))
	p.Read() //:
	fmt.Println(string(p.Byte))
	err = p.SkipWS() //-
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '-' {
		t.Error("Failed to parse [ - ]")
	}
	v, err = p.ParseFloat()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(v)
	if p.Byte != ',' {
		t.Error("Failed to parse [ , ]")
	}

	p.SkipWS() //b
	fmt.Println(string(p.Byte))
	p.Read() //:
	fmt.Println(string(p.Byte))
	err = p.SkipWS() //2
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '2' {
		t.Error("Failed to parse [ 2 ]")
	}
	v, err = p.ParseFloat()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '}' {
		t.Error("Failed to parse [ } ]")
	}
	fmt.Println(v)
}
