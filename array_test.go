package v2

import (
	"fmt"
	"testing"
)

func TestArrayParse(t *testing.T) {
	s := "[ 1, 2, 3 ]"
	p := newParser([]byte(s))
	err := p.SkipWS()
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '[' {
		t.Error("Failed to parse [")
	}
	v, err := p.ParseArray(false)
	if err != errEOF {
		t.Error("Failed to parse [ 1, 2, 3 ]")
	}
	fmt.Println(v.String())

	s = `{ a: [ "one", "${id?two}", 3, true ] }`
	p = newParser([]byte(s))
	p.SkipWS()       //{
	p.SkipWS()       //a
	p.Read()         //:
	err = p.SkipWS() //[
	if err != nil {
		t.Error(err.Error())
	}
	if p.Byte != '[' {
		t.Error("Failed to parse [")
	}
	v, err = p.ParseArray(true)
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Println(v.String())
}
