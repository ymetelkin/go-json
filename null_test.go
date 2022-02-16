package v2

import (
	"fmt"
	"testing"
)

func TestNullParse(t *testing.T) {
	s := "null"
	p := newParser([]byte(s))
	err := p.SkipWS()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != 'n' {
		t.Error("Failed to parse [ n ]")
	}
	v, err := p.ParseNull()
	if err != errEOF {
		t.Error("Failed to parse [ null ]")
	}
	fmt.Println(v)

	s = "{ v: null }"
	p = newParser([]byte(s))
	p.SkipWS() //{
	fmt.Println(string(p.Byte))
	p.SkipWS() //v
	fmt.Println(string(p.Byte))
	p.Read() //:
	fmt.Println(string(p.Byte))
	err = p.SkipWS() //n
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != 'n' {
		t.Error("Failed to parse [ n ]")
	}
	v, err = p.ParseNull()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(v)

	if p.Byte != '}' {
		t.Error("Failed to parse [ } ]")
	}
	fmt.Println(v)
}
