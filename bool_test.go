package v2

import (
	"fmt"
	"testing"
)

func TestBoolParse(t *testing.T) {
	s := "true"
	p := newParser([]byte(s))
	err := p.SkipWS()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != 't' {
		t.Error("Failed to parse [ t ]")
	}
	v, err := p.ParseTrue()
	if err != errEOF {
		t.Error("Failed to parse [ true ]")
	}
	fmt.Println(v)

	s = "{ t: true, f: false }"
	p = newParser([]byte(s))
	p.SkipWS() //{
	fmt.Println(string(p.Byte))
	p.SkipWS() //t
	fmt.Println(string(p.Byte))
	p.Read() //:
	fmt.Println(string(p.Byte))
	err = p.SkipWS() //t
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != 't' {
		t.Error("Failed to parse [ t ]")
	}
	v, err = p.ParseTrue()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(v)
	if p.Byte != ',' {
		t.Error("Failed to parse [ , ]")
	}

	p.SkipWS() //f
	fmt.Println(string(p.Byte))
	p.Read() //:
	fmt.Println(string(p.Byte))
	err = p.SkipWS() //f
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != 'f' {
		t.Error("Failed to parse [ f ]")
	}
	v, err = p.ParseFalse()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '}' {
		t.Error("Failed to parse [ } ]")
	}
	fmt.Println(v)
}
