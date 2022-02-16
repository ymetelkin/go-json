package v2

import (
	"fmt"
	"testing"
)

func TestIntParse(t *testing.T) {
	s := "314"
	p := newParser([]byte(s))
	err := p.SkipWS()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '3' {
		t.Error("Failed to parse [ 3 ]")
	}
	v, err := p.ParseInt()
	if v.Value() != 314 || err != errEOF {
		t.Error("Failed to parse [ 314 ]")
	}
	fmt.Println(v)

	s = "2.001"
	p = newParser([]byte(s))
	err = p.SkipWS()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '2' {
		t.Error("Failed to parse [ 2 ]")
	}
	v, err = p.ParseInt()
	if v.Value() != 2 || err != errEOF {
		t.Error("Failed to parse [ 2 ]")
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
	v, err = p.ParseInt()
	if v.Value() != -123 || err != errEOF {
		t.Error("Failed to parse [ -123 ]")
	}
	fmt.Println(v)

	s = "{ a: -123, b: 2.3E2 }"
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
	v, err = p.ParseInt()
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
	v, err = p.ParseInt()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '}' {
		t.Error("Failed to parse [ } ]")
	}
	fmt.Println(v)
}
