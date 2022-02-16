package v2

import (
	"fmt"
	"testing"
)

func TestNumberParse(t *testing.T) {
	s := "{ i: -123, f: 2.3E-2 }"
	p := newParser([]byte(s))
	p.SkipWS()        //{
	p.SkipWS()        //i
	p.Read()          //:
	err := p.SkipWS() //-
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '-' {
		t.Error("Failed to parse [ - ]")
	}
	v, err := p.ParseNumber()
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("%v\t%v\t%s\n", v.Type(), v.Value(), v.String())
	if p.Byte != ',' {
		t.Error("Failed to parse [ , ]")
	}

	p.SkipWS()       //f
	p.Read()         //:
	err = p.SkipWS() //2
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '2' {
		t.Error("Failed to parse [ 2 ]")
	}
	v, err = p.ParseNumber()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '}' {
		t.Error("Failed to parse [ } ]")
	}
	fmt.Printf("%v\t%v\t%s\n", v.Type(), v.Value(), v.String())

	s = `{"id":9223372036854776000}`
	jo, err := ParseObject([]byte(s))
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(jo.String())
}
